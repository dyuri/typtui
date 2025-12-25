package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dyuri/typtui/internal/parser"
)

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// In edit mode, handle special keys first, then forward to inputs
		if m.mode == ModeEdit {
			return m.handleEditModeKeyPress(msg)
		}
		// In XPM edit mode, handle XPM-specific navigation
		if m.mode == ModeEditXPM {
			return m.handleXPMEditKeyPress(msg)
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
	// Handle confirm quit mode separately
	if m.mode == ModeConfirmQuit {
		return m.handleConfirmQuit(msg)
	}

	// Clear status message on any key press
	m.status = ""

	switch msg.String() {
	case "q", "ctrl+c":
		// Check for unsaved changes
		if m.modified {
			m.mode = ModeConfirmQuit
			return m, nil
		}
		return m, tea.Quit

	case "ctrl+s":
		// Save file in list or detail mode
		if (m.mode == ModeList || m.mode == ModeDetail) && m.typFile != nil {
			if err := m.saveFile(); err != nil {
				m.status = fmt.Sprintf("Error saving: %v", err)
			}
		}
		return m, nil

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

	case "x":
		if m.mode == ModeDetail && m.typFile != nil {
			// Enter XPM edit mode - default to DayXpm
			switch m.activeTab {
			case TabPoints:
				if m.selectedIdx < len(m.typFile.Points) && m.typFile.Points[m.selectedIdx].DayXpm != nil {
					m.editingXPM = m.typFile.Points[m.selectedIdx].DayXpm
					m.editingXPMType = "DayXpm"
					m.xpmColorIdx = 0
					m.mode = ModeEditXPM
				}
			case TabLines:
				if m.selectedIdx < len(m.typFile.Lines) && m.typFile.Lines[m.selectedIdx].DayXpm != nil {
					m.editingXPM = m.typFile.Lines[m.selectedIdx].DayXpm
					m.editingXPMType = "Xpm"
					m.xpmColorIdx = 0
					m.mode = ModeEditXPM
				}
			case TabPolygons:
				if m.selectedIdx < len(m.typFile.Polygons) && m.typFile.Polygons[m.selectedIdx].DayXpm != nil {
					m.editingXPM = m.typFile.Polygons[m.selectedIdx].DayXpm
					m.editingXPMType = "Xpm"
					m.xpmColorIdx = 0
					m.mode = ModeEditXPM
				}
			}
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

	// Update the appropriate structure
	switch m.activeTab {
	case TabPoints:
		if m.selectedIdx < len(m.typFile.Points) && len(m.inputs) >= 6 {
			// Type (index 0)
			m.typFile.Points[m.selectedIdx].Type = m.inputs[0].Value()

			// SubType (index 1)
			m.typFile.Points[m.selectedIdx].SubType = m.inputs[1].Value()

			// Label (index 2)
			if m.typFile.Points[m.selectedIdx].Labels == nil {
				m.typFile.Points[m.selectedIdx].Labels = make(map[string]string)
			}
			m.typFile.Points[m.selectedIdx].Labels["0x04"] = m.inputs[2].Value()

			// FontStyle (index 3)
			m.typFile.Points[m.selectedIdx].FontStyle = m.inputs[3].Value()

			// Day Color (index 4)
			dayColorValue := m.inputs[4].Value()
			if dayColorValue != "" {
				// Ensure # prefix
				if !strings.HasPrefix(dayColorValue, "#") {
					dayColorValue = "#" + dayColorValue
				}
				// Update or create day color
				if len(m.typFile.Points[m.selectedIdx].DayColors) > 0 {
					m.typFile.Points[m.selectedIdx].DayColors[0].Hex = dayColorValue
				} else {
					m.typFile.Points[m.selectedIdx].DayColors = []parser.Color{
						{Hex: dayColorValue, Day: true},
					}
				}
			}

			// Night Color (index 5)
			nightColorValue := m.inputs[5].Value()
			if nightColorValue != "" {
				// Ensure # prefix
				if !strings.HasPrefix(nightColorValue, "#") {
					nightColorValue = "#" + nightColorValue
				}
				// Update or create night color
				if len(m.typFile.Points[m.selectedIdx].NightColors) > 0 {
					m.typFile.Points[m.selectedIdx].NightColors[0].Hex = nightColorValue
				} else {
					m.typFile.Points[m.selectedIdx].NightColors = []parser.Color{
						{Hex: nightColorValue, Day: false},
					}
				}
			}
		}

	case TabLines:
		if m.selectedIdx < len(m.typFile.Lines) && len(m.inputs) >= 6 {
			// Type (index 0)
			m.typFile.Lines[m.selectedIdx].Type = m.inputs[0].Value()

			// Label (index 1)
			if m.typFile.Lines[m.selectedIdx].Labels == nil {
				m.typFile.Lines[m.selectedIdx].Labels = make(map[string]string)
			}
			m.typFile.Lines[m.selectedIdx].Labels["0x04"] = m.inputs[1].Value()

			// LineWidth (index 2)
			if width, err := strconv.Atoi(m.inputs[2].Value()); err == nil {
				m.typFile.Lines[m.selectedIdx].LineWidth = width
			}

			// BorderWidth (index 3)
			if width, err := strconv.Atoi(m.inputs[3].Value()); err == nil {
				m.typFile.Lines[m.selectedIdx].BorderWidth = width
			}

			// LineStyle (index 4)
			m.typFile.Lines[m.selectedIdx].LineStyle = m.inputs[4].Value()

			// UseOrientation (index 5)
			useOrient := strings.ToUpper(m.inputs[5].Value())
			m.typFile.Lines[m.selectedIdx].UseOrientation = (useOrient == "Y" || useOrient == "YES")
		}

	case TabPolygons:
		if m.selectedIdx < len(m.typFile.Polygons) && len(m.inputs) >= 4 {
			// Type (index 0)
			m.typFile.Polygons[m.selectedIdx].Type = m.inputs[0].Value()

			// Label (index 1)
			if m.typFile.Polygons[m.selectedIdx].Labels == nil {
				m.typFile.Polygons[m.selectedIdx].Labels = make(map[string]string)
			}
			m.typFile.Polygons[m.selectedIdx].Labels["0x04"] = m.inputs[1].Value()

			// ExtendedLabels (index 2)
			extLabels := strings.ToUpper(m.inputs[2].Value())
			m.typFile.Polygons[m.selectedIdx].ExtendedLabels = (extLabels == "Y" || extLabels == "YES")

			// FontStyle (index 3)
			m.typFile.Polygons[m.selectedIdx].FontStyle = m.inputs[3].Value()
		}
	}
}

// handleConfirmQuit handles the confirm quit dialog
func (m Model) handleConfirmQuit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		// Save and quit
		if err := m.saveFile(); err != nil {
			m.status = fmt.Sprintf("Error saving: %v", err)
			m.mode = ModeList
			return m, nil
		}
		return m, tea.Quit

	case "n", "N":
		// Quit without saving
		return m, tea.Quit

	case "esc", "c", "C":
		// Cancel quit and return to list
		m.mode = ModeList
		return m, nil
	}

	return m, nil
}

// handleXPMEditKeyPress handles keyboard input in XPM edit mode
func (m Model) handleXPMEditKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.editingXPM == nil {
		return m, nil
	}

	// If we're editing a color (input is active), handle that differently
	if len(m.inputs) > 0 {
		switch msg.String() {
		case "enter":
			// Save the color change
			newColor := m.inputs[0].Value()
			if !strings.HasPrefix(newColor, "#") {
				newColor = "#" + newColor
			}

			// Find and update the color at xpmColorIdx
			// Convert to sorted slice for consistent ordering
			type colorEntry struct {
				char  string
				color parser.Color
			}
			var colors []colorEntry
			for char, color := range m.editingXPM.Palette {
				colors = append(colors, colorEntry{char, color})
			}

			// Sort alphabetically by character (same as in view)
			sort.Slice(colors, func(i, j int) bool {
				return colors[i].char < colors[j].char
			})

			// Update the selected color
			if m.xpmColorIdx < len(colors) {
				selectedChar := colors[m.xpmColorIdx].char
				color := m.editingXPM.Palette[selectedChar]
				color.Hex = newColor
				m.editingXPM.Palette[selectedChar] = color
			}

			m.inputs = nil
			return m, nil

		case "esc":
			// Cancel color edit
			m.inputs = nil
			return m, nil

		default:
			// Forward to text input
			m.inputs[0], cmd = m.inputs[0].Update(msg)
			return m, cmd
		}
	}

	maxColors := len(m.editingXPM.Palette)

	switch msg.String() {
	case "ctrl+s":
		// Save changes and return to detail view
		m.modified = true
		m.mode = ModeDetail
		m.editingXPM = nil
		m.status = "XPM changes saved"
		return m, nil

	case "esc":
		// Cancel and return to detail view
		m.mode = ModeDetail
		m.editingXPM = nil
		return m, nil

	case "up", "k":
		// Navigate palette colors up
		if m.xpmColorIdx > 0 {
			m.xpmColorIdx--
		}
		return m, nil

	case "down", "j":
		// Navigate palette colors down
		if m.xpmColorIdx < maxColors-1 {
			m.xpmColorIdx++
		}
		return m, nil

	case "enter":
		// Edit the selected color
		return m.enterColorEdit()
	}

	return m, nil
}

// enterColorEdit enters color editing mode for the selected palette entry
func (m Model) enterColorEdit() (tea.Model, tea.Cmd) {
	if m.editingXPM == nil || len(m.editingXPM.Palette) == 0 {
		return m, nil
	}

	// Get the color at the current index
	// Convert to sorted slice for consistent ordering
	type colorEntry struct {
		char  string
		color parser.Color
	}
	var colors []colorEntry
	for char, color := range m.editingXPM.Palette {
		colors = append(colors, colorEntry{char, color})
	}

	// Sort alphabetically by character (same as in view)
	sort.Slice(colors, func(i, j int) bool {
		return colors[i].char < colors[j].char
	})

	// Get the selected color
	if m.xpmColorIdx >= len(colors) {
		return m, nil
	}
	selectedEntry := colors[m.xpmColorIdx]

	// Create a single text input for the color
	input := textinput.New()
	input.Placeholder = "#RRGGBB"
	input.CharLimit = 7
	input.Width = 30
	input.SetValue(selectedEntry.color.Hex)
	input.Focus()
	input.Prompt = fmt.Sprintf("Color for '%s': ", selectedEntry.char)

	m.inputs = []textinput.Model{input}
	m.focusedField = 0

	// Stay in XPM edit mode but now show the input
	return m, nil
}
