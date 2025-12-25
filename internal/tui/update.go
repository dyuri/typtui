package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// In edit mode, handle special keys first, then forward to inputs
		if m.mode == ModeEdit {
			return m.handleEditModeKeyPress(msg)
		}
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case fileLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			m.mode = ModeError
			return m, nil
		}
		m.typFile = msg.typFile
		m.mode = ModeList
		return m, nil
	}

	return m, nil
}

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		// TODO: Check for unsaved changes
		return m, tea.Quit

	case "?":
		if m.mode == ModeHelp {
			m.mode = ModeList
		} else {
			m.mode = ModeHelp
		}
		return m, nil

	case "tab":
		if m.mode == ModeList {
			m.activeTab = (m.activeTab + 1) % 3
			m.selectedIdx = 0
		}
		return m, nil

	case "up", "k":
		if m.mode == ModeList {
			if m.selectedIdx > 0 {
				m.selectedIdx--
			}
		}
		return m, nil

	case "down", "j":
		if m.mode == ModeList {
			maxIdx := m.getMaxIndex()
			if m.selectedIdx < maxIdx-1 {
				m.selectedIdx++
			}
		}
		return m, nil

	case "enter":
		if m.mode == ModeList && m.typFile != nil && m.getMaxIndex() > 0 {
			m.mode = ModeDetail
		}
		return m, nil

	case "esc":
		if m.mode == ModeDetail {
			m.mode = ModeList
		} else if m.mode == ModeEdit {
			// Cancel editing and return to detail view
			m.mode = ModeDetail
			m.inputs = nil
		}
		return m, nil

	case "e":
		if m.mode == ModeDetail && m.typFile != nil {
			// Enter edit mode
			m.mode = ModeEdit
			// Initialize inputs based on current tab
			switch m.activeTab {
			case TabPoints:
				if m.selectedIdx < len(m.typFile.Points) {
					m.initPointEditInputs(m.typFile.Points[m.selectedIdx])
				}
			case TabLines:
				if m.selectedIdx < len(m.typFile.Lines) {
					m.initLineEditInputs(m.typFile.Lines[m.selectedIdx])
				}
			case TabPolygons:
				if m.selectedIdx < len(m.typFile.Polygons) {
					m.initPolygonEditInputs(m.typFile.Polygons[m.selectedIdx])
				}
			}
		}
		return m, nil

	case "ctrl+s":
		if m.mode == ModeEdit {
			// Save changes
			m.saveEdits()
			m.modified = true
			m.mode = ModeDetail
			m.inputs = nil
		}
		return m, nil
	}

	return m, nil
}

// getMaxIndex returns the maximum index for the current tab
func (m Model) getMaxIndex() int {
	if m.typFile == nil {
		return 0
	}

	switch m.activeTab {
	case TabPoints:
		return len(m.typFile.Points)
	case TabLines:
		return len(m.typFile.Lines)
	case TabPolygons:
		return len(m.typFile.Polygons)
	}
	return 0
}

// handleEditModeKeyPress handles keyboard input in edit mode
func (m Model) handleEditModeKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.String() {
	case "ctrl+s":
		// Save changes
		m.saveEdits()
		m.modified = true
		m.mode = ModeDetail
		m.inputs = nil
		return m, nil

	case "esc":
		// Cancel editing
		m.mode = ModeDetail
		m.inputs = nil
		return m, nil

	case "tab", "shift+tab", "up", "down":
		// Navigate between fields
		if msg.String() == "tab" || msg.String() == "down" {
			m.focusedField++
			if m.focusedField >= len(m.inputs) {
				m.focusedField = 0
			}
		} else {
			m.focusedField--
			if m.focusedField < 0 {
				m.focusedField = len(m.inputs) - 1
			}
		}

		// Update focus
		for i := range m.inputs {
			if i == m.focusedField {
				m.inputs[i].Focus()
			} else {
				m.inputs[i].Blur()
			}
		}
		return m, nil
	}

	// Forward key to focused input
	if m.focusedField < len(m.inputs) {
		m.inputs[m.focusedField], cmd = m.inputs[m.focusedField].Update(msg)
	}

	return m, cmd
}

// saveEdits saves the current edit form values back to the data structure
func (m *Model) saveEdits() {
	if len(m.inputs) < 2 {
		return
	}

	// Get values from inputs
	typeValue := m.inputs[0].Value()
	labelValue := m.inputs[1].Value()

	// Update the appropriate structure
	switch m.activeTab {
	case TabPoints:
		if m.selectedIdx < len(m.typFile.Points) {
			m.typFile.Points[m.selectedIdx].Type = typeValue
			if m.typFile.Points[m.selectedIdx].Labels == nil {
				m.typFile.Points[m.selectedIdx].Labels = make(map[string]string)
			}
			m.typFile.Points[m.selectedIdx].Labels["0x04"] = labelValue
		}

	case TabLines:
		if m.selectedIdx < len(m.typFile.Lines) {
			m.typFile.Lines[m.selectedIdx].Type = typeValue
			if m.typFile.Lines[m.selectedIdx].Labels == nil {
				m.typFile.Lines[m.selectedIdx].Labels = make(map[string]string)
			}
			m.typFile.Lines[m.selectedIdx].Labels["0x04"] = labelValue
		}

	case TabPolygons:
		if m.selectedIdx < len(m.typFile.Polygons) {
			m.typFile.Polygons[m.selectedIdx].Type = typeValue
			if m.typFile.Polygons[m.selectedIdx].Labels == nil {
				m.typFile.Polygons[m.selectedIdx].Labels = make(map[string]string)
			}
			m.typFile.Polygons[m.selectedIdx].Labels["0x04"] = labelValue
		}
	}
}
