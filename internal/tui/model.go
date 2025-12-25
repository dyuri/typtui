package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dyuri/typtui/internal/parser"
)

// Mode represents the current UI mode
type Mode int

const (
	ModeList Mode = iota
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
