package tui

import (
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
	ModeHelp
	ModeError
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

	// Messages
	err  error
	info string

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
	inputs := make([]textinput.Model, 2)

	// Type field
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "e.g., 0x2f06"
	inputs[0].Focus()
	inputs[0].CharLimit = 10
	inputs[0].Width = 30
	inputs[0].SetValue(point.Type)
	inputs[0].Prompt = "Type: "

	// English label field (0x04)
	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Label"
	inputs[1].CharLimit = 50
	inputs[1].Width = 50
	if label, ok := point.Labels["0x04"]; ok {
		inputs[1].SetValue(label)
	}
	inputs[1].Prompt = "Label (EN): "

	m.inputs = inputs
	m.focusedField = 0
}

// initLineEditInputs initializes text inputs for editing a line type
func (m *Model) initLineEditInputs(line parser.LineType) {
	inputs := make([]textinput.Model, 2)

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

	m.inputs = inputs
	m.focusedField = 0
}

// initPolygonEditInputs initializes text inputs for editing a polygon type
func (m *Model) initPolygonEditInputs(polygon parser.PolygonType) {
	inputs := make([]textinput.Model, 2)

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

	m.inputs = inputs
	m.focusedField = 0
}
