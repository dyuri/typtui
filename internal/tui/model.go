package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dyuri/typtui/internal/parser"
)

// Mode represents the current UI mode
type Mode int

const (
	ModeList Mode = iota
	ModeDetail
	ModeEdit
	ModeEditXPM
	ModeHelp
	ModeError
	ModeConfirmQuit
)

// Tab represents the active tab
type Tab int

const (
	TabPoints Tab = iota
	TabLines
	TabPolygons
)

// Model is the main application model
type Model struct {
	// Core data
	typFile  *parser.TYPFile
	modified bool

	// UI state
	mode   Mode
	width  int
	height int

	// Navigation
	activeTab   Tab
	selectedIdx int

	// Components
	list list.Model

	// Edit mode state
	focusedField int
	inputs       []textinput.Model

	// XPM edit mode state
	editingXPM     *parser.XPMIcon
	editingXPMType string // "DayXpm", "NightXpm", etc.
	xpmColorIdx    int    // Currently selected color in palette

	// Messages
	err    error
	status string

	// File path (if loaded from command line)
	filePath string
}

// NewModel creates a new TUI model
func NewModel(filePath string) Model {
	return Model{
		mode:        ModeList,
		activeTab:   TabPoints,
		selectedIdx: 0,
		filePath:    filePath,
	}
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// If a file path was provided, load it
	if m.filePath != "" {
		return loadFileCmd(m.filePath)
	}
	return nil
}

// fileLoadedMsg is sent when a file is loaded
type fileLoadedMsg struct {
	typFile *parser.TYPFile
	err     error
}

// loadFileCmd loads a TYP file
func loadFileCmd(filePath string) tea.Cmd {
	return func() tea.Msg {
		typFile, err := parser.ParseFile(filePath)
		return fileLoadedMsg{typFile: typFile, err: err}
	}
}

// initPointEditInputs initializes text inputs for editing a point type
func (m *Model) initPointEditInputs(point parser.PointType) {
	inputs := make([]textinput.Model, 6)

	// Type field
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "e.g., 0x2f06"
	inputs[0].Focus()
	inputs[0].CharLimit = 10
	inputs[0].Width = 30
	inputs[0].SetValue(point.Type)
	inputs[0].Prompt = "Type: "

	// SubType field
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "e.g., 0x00 (optional)"
	inputs[1].CharLimit = 10
	inputs[1].Width = 30
	inputs[1].SetValue(point.SubType)
	inputs[1].Prompt = "SubType: "

	// English label field (0x04)
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Label"
	inputs[2].CharLimit = 50
	inputs[2].Width = 50
	if label, ok := point.Labels["0x04"]; ok {
		inputs[2].SetValue(label)
	}
	inputs[2].Prompt = "Label (EN): "

	// FontStyle field
	inputs[3] = textinput.New()
	inputs[3].Placeholder = "NoLabel, SmallFont, NormalFont, LargeFont"
	inputs[3].CharLimit = 20
	inputs[3].Width = 40
	inputs[3].SetValue(point.FontStyle)
	inputs[3].Prompt = "FontStyle: "

	// Day Color field
	inputs[4] = textinput.New()
	inputs[4].Placeholder = "#RRGGBB (e.g., #FF0000)"
	inputs[4].CharLimit = 7
	inputs[4].Width = 30
	if len(point.DayColors) > 0 {
		inputs[4].SetValue(point.DayColors[0].Hex)
	}
	inputs[4].Prompt = "Day Color: "

	// Night Color field
	inputs[5] = textinput.New()
	inputs[5].Placeholder = "#RRGGBB (e.g., #0000FF)"
	inputs[5].CharLimit = 7
	inputs[5].Width = 30
	if len(point.NightColors) > 0 {
		inputs[5].SetValue(point.NightColors[0].Hex)
	}
	inputs[5].Prompt = "Night Color: "

	m.inputs = inputs
	m.focusedField = 0
}

// initLineEditInputs initializes text inputs for editing a line type
func (m *Model) initLineEditInputs(line parser.LineType) {
	inputs := make([]textinput.Model, 6)

	// Type field
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "e.g., 0x01"
	inputs[0].Focus()
	inputs[0].CharLimit = 10
	inputs[0].Width = 30
	inputs[0].SetValue(line.Type)
	inputs[0].Prompt = "Type: "

	// English label field
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Label"
	inputs[1].CharLimit = 50
	inputs[1].Width = 50
	if label, ok := line.Labels["0x04"]; ok {
		inputs[1].SetValue(label)
	}
	inputs[1].Prompt = "Label (EN): "

	// LineWidth field
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "e.g., 5"
	inputs[2].CharLimit = 5
	inputs[2].Width = 20
	if line.LineWidth > 0 {
		inputs[2].SetValue(fmt.Sprintf("%d", line.LineWidth))
	}
	inputs[2].Prompt = "LineWidth: "

	// BorderWidth field
	inputs[3] = textinput.New()
	inputs[3].Placeholder = "e.g., 1"
	inputs[3].CharLimit = 5
	inputs[3].Width = 20
	if line.BorderWidth > 0 {
		inputs[3].SetValue(fmt.Sprintf("%d", line.BorderWidth))
	}
	inputs[3].Prompt = "BorderWidth: "

	// LineStyle field
	inputs[4] = textinput.New()
	inputs[4].Placeholder = "solid, dashed, dotted"
	inputs[4].CharLimit = 20
	inputs[4].Width = 30
	inputs[4].SetValue(line.LineStyle)
	inputs[4].Prompt = "LineStyle: "

	// UseOrientation field
	inputs[5] = textinput.New()
	inputs[5].Placeholder = "Y or N"
	inputs[5].CharLimit = 1
	inputs[5].Width = 10
	if line.UseOrientation {
		inputs[5].SetValue("Y")
	} else {
		inputs[5].SetValue("N")
	}
	inputs[5].Prompt = "UseOrientation: "

	m.inputs = inputs
	m.focusedField = 0
}

// initPolygonEditInputs initializes text inputs for editing a polygon type
func (m *Model) initPolygonEditInputs(polygon parser.PolygonType) {
	inputs := make([]textinput.Model, 4)

	// Type field
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "e.g., 0x13"
	inputs[0].Focus()
	inputs[0].CharLimit = 10
	inputs[0].Width = 30
	inputs[0].SetValue(polygon.Type)
	inputs[0].Prompt = "Type: "

	// English label field
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Label"
	inputs[1].CharLimit = 50
	inputs[1].Width = 50
	if label, ok := polygon.Labels["0x04"]; ok {
		inputs[1].SetValue(label)
	}
	inputs[1].Prompt = "Label (EN): "

	// ExtendedLabels field
	inputs[2] = textinput.New()
	inputs[2].Placeholder = "Y or N"
	inputs[2].CharLimit = 1
	inputs[2].Width = 10
	if polygon.ExtendedLabels {
		inputs[2].SetValue("Y")
	} else {
		inputs[2].SetValue("N")
	}
	inputs[2].Prompt = "ExtendedLabels: "

	// FontStyle field
	inputs[3] = textinput.New()
	inputs[3].Placeholder = "NoLabel, SmallFont, NormalFont, LargeFont"
	inputs[3].CharLimit = 20
	inputs[3].Width = 40
	inputs[3].SetValue(polygon.FontStyle)
	inputs[3].Prompt = "FontStyle: "

	m.inputs = inputs
	m.focusedField = 0
}

// saveFile saves the current TYPFile to disk
func (m *Model) saveFile() error {
	if m.typFile == nil {
		return fmt.Errorf("no file loaded")
	}

	if m.filePath == "" {
		return fmt.Errorf("no file path")
	}

	if err := parser.WriteFile(m.typFile, m.filePath); err != nil {
		return err
	}

	m.modified = false
	m.status = "File saved successfully"
	return nil
}
