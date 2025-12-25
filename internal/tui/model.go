package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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
	xpmViewport    viewport.Model

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

// initXPMViewport initializes the viewport for XPM editing
func (m *Model) initXPMViewport() {
	// Reserve space for header (3 lines) and footer (4 lines)
	headerFooterHeight := 7
	viewportHeight := m.height - headerFooterHeight
	if viewportHeight < 10 {
		viewportHeight = 10 // Minimum height
	}

	m.xpmViewport = viewport.New(m.width, viewportHeight)
	m.xpmViewport.YPosition = 3 // Position after header

	// Set initial content
	m.updateXPMViewportContent()
}

// updateXPMViewportContent rebuilds and updates the viewport content
func (m *Model) updateXPMViewportContent() {
	if m.editingXPM == nil {
		return
	}

	var content strings.Builder

	// XPM Info
	content.WriteString(fmt.Sprintf("Size: %dx%d, Colors: %d, Chars/pixel: %d\n\n",
		m.editingXPM.Width, m.editingXPM.Height, m.editingXPM.Colors, m.editingXPM.CharsPerPixel))

	// Color Palette
	content.WriteString("Color Palette\n\n")

	// Convert map to sorted slice for consistent ordering
	type colorEntry struct {
		char  string
		color parser.Color
	}
	var colors []colorEntry
	for char, color := range m.editingXPM.Palette {
		colors = append(colors, colorEntry{char, color})
	}

	// Sort alphabetically by character
	sort.Slice(colors, func(i, j int) bool {
		return colors[i].char < colors[j].char
	})

	for i, entry := range colors {
		prefix := "  "
		if i == m.xpmColorIdx {
			prefix = "▸ "
		}

		// Render color with preview
		colorDisplay := m.renderColorPreview(entry.color.Hex)
		content.WriteString(fmt.Sprintf("%s%s → %s\n", prefix, entry.char, colorDisplay))
	}

	content.WriteString("\n")

	// Icon Preview
	content.WriteString("Icon Preview\n")

	// Render all rows of the XPM with colors
	for i := 0; i < len(m.editingXPM.Data); i++ {
		content.WriteString("  ")
		row := m.editingXPM.Data[i]

		// Process each character/pixel with colors
		for _, char := range row {
			charStr := string(char)

			// Look up the color for this character
			if color, ok := m.editingXPM.Palette[charStr]; ok {
				content.WriteString(m.renderPixelColored(color.Hex, charStr))
			} else {
				// Unknown character, show as gray
				content.WriteString(m.renderPixelColored("#808080", charStr))
			}
		}

		// Reset color at end of line
		content.WriteString("\x1b[0m\n")
	}

	m.xpmViewport.SetContent(content.String())
}

// renderColorPreview renders a color with a visual preview block
func (m Model) renderColorPreview(hexColor string) string {
	if hexColor == "" || hexColor == "none" || hexColor == "transparent" {
		return hexColor
	}

	// Parse hex color
	hex := strings.TrimPrefix(hexColor, "#")
	if len(hex) != 6 {
		return hexColor // Invalid format, just return as-is
	}

	// Convert hex to RGB
	r, err1 := strconv.ParseInt(hex[0:2], 16, 64)
	g, err2 := strconv.ParseInt(hex[2:4], 16, 64)
	b, err3 := strconv.ParseInt(hex[4:6], 16, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		return hexColor // Parse error, return as-is
	}

	// Create ANSI 24-bit true color code for the block character
	colorPreview := fmt.Sprintf("\x1b[38;2;%d;%d;%dm■\x1b[0m", r, g, b)
	return fmt.Sprintf("%s %s", hexColor, colorPreview)
}

// renderPixelColored renders a pixel with colored background and contrasting text
func (m Model) renderPixelColored(hexColor string, char string) string {
	// Handle transparent/none colors
	if hexColor == "" || hexColor == "none" || hexColor == "transparent" {
		// Use a light gray background with the character in black
		return fmt.Sprintf("\x1b[48;2;240;240;240m\x1b[38;2;0;0;0m%s\x1b[0m", char)
	}

	// Parse hex color
	hex := strings.TrimPrefix(hexColor, "#")
	if len(hex) != 6 {
		// Invalid color, use gray with black text
		return fmt.Sprintf("\x1b[48;2;128;128;128m\x1b[38;2;0;0;0m%s\x1b[0m", char)
	}

	// Convert hex to RGB
	r, err1 := strconv.ParseInt(hex[0:2], 16, 64)
	g, err2 := strconv.ParseInt(hex[2:4], 16, 64)
	b, err3 := strconv.ParseInt(hex[4:6], 16, 64)

	if err1 != nil || err2 != nil || err3 != nil {
		// Parse error, use gray with black text
		return fmt.Sprintf("\x1b[48;2;128;128;128m\x1b[38;2;0;0;0m%s\x1b[0m", char)
	}

	// Calculate luminance using the relative luminance formula
	// Y = 0.299*R + 0.587*G + 0.114*B
	luminance := float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114

	// Choose foreground color based on luminance
	// If luminance > 128, use black text; otherwise use white text
	var fgR, fgG, fgB int64
	if luminance > 128 {
		fgR, fgG, fgB = 0, 0, 0 // Black
	} else {
		fgR, fgG, fgB = 255, 255, 255 // White
	}

	// Create ANSI 24-bit background and foreground color with the character
	return fmt.Sprintf("\x1b[48;2;%d;%d;%dm\x1b[38;2;%d;%d;%dm%s\x1b[0m", r, g, b, fgR, fgG, fgB, char)
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
