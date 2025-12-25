package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
