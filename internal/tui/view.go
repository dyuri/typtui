package tui

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/dyuri/typtui/internal/parser"
)

var (
	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170"))

	tabStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240"))

	activeTabStyle = tabStyle.Copy().
			BorderForeground(lipgloss.Color("170")).
			Bold(true)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true)
)

// View renders the UI
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	switch m.mode {
	case ModeError:
		return m.viewError()
	case ModeHelp:
		return m.viewHelp()
	case ModeDetail:
		return m.viewDetail()
	case ModeEdit:
		return m.viewEdit()
	case ModeEditXPM:
		return m.viewEditXPM()
	case ModeConfirmQuit:
		return m.viewConfirmQuit()
	default:
		return m.viewList()
	}
}

// viewError renders the error screen
func (m Model) viewError() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Error"))
	b.WriteString("\n\n")
	b.WriteString(errorStyle.Render(m.err.Error()))
	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("Press q to quit"))

	return b.String()
}

// viewHelp renders the help screen
func (m Model) viewHelp() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("typtui - Help"))
	b.WriteString("\n\n")
	b.WriteString("Keyboard Shortcuts:\n\n")
	b.WriteString("List View:\n")
	b.WriteString("  q, Ctrl+C    Quit\n")
	b.WriteString("  ?            Toggle help\n")
	b.WriteString("  Ctrl+S       Save file to disk\n")
	b.WriteString("  Tab          Switch between tabs (Points/Lines/Polygons)\n")
	b.WriteString("  ↑/k          Move up\n")
	b.WriteString("  ↓/j          Move down\n")
	b.WriteString("  Enter        View details of selected item\n")
	b.WriteString("\n")
	b.WriteString("Detail View:\n")
	b.WriteString("  e            Edit selected item\n")
	b.WriteString("  Ctrl+S       Save file to disk\n")
	b.WriteString("  Esc          Return to list view\n")
	b.WriteString("\n")
	b.WriteString("Edit Mode:\n")
	b.WriteString("  Ctrl+S       Save changes to item and return to detail\n")
	b.WriteString("  Esc          Cancel editing\n")
	b.WriteString("  Tab/↑/↓      Navigate between fields\n")
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press ? to return to the main view"))

	return b.String()
}

// viewList renders the main list view
func (m Model) viewList() string {
	if m.typFile == nil {
		return "No file loaded. Usage: typtui <file.typ>"
	}

	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")

	// Tabs
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	// Content
	b.WriteString(m.renderContent())
	b.WriteString("\n\n")

	// Footer
	b.WriteString(m.renderFooter())

	return b.String()
}

// renderHeader renders the header section
func (m Model) renderHeader() string {
	fileName := m.typFile.FilePath
	if m.modified {
		fileName += " [Modified]"
	}
	return titleStyle.Render("typtui - " + fileName)
}

// renderTabs renders the tab bar
func (m Model) renderTabs() string {
	tabs := []string{"Points", "Lines", "Polygons"}
	var renderedTabs []string

	for i, tab := range tabs {
		style := tabStyle
		if Tab(i) == m.activeTab {
			style = activeTabStyle
		}
		renderedTabs = append(renderedTabs, style.Render(tab))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

// renderContent renders the main content area
func (m Model) renderContent() string {
	var b strings.Builder

	switch m.activeTab {
	case TabPoints:
		if len(m.typFile.Points) == 0 {
			b.WriteString(statusStyle.Render("No points defined"))
		} else {
			for i, point := range m.typFile.Points {
				label := ""
				if len(point.Labels) > 0 {
					// Get first label
					for _, l := range point.Labels {
						label = l
						break
					}
				}

				line := fmt.Sprintf("  %s - %s", point.Type, label)
				if i == m.selectedIdx {
					line = selectedStyle.Render("▸ " + point.Type + " - " + label)
				}
				b.WriteString(line)
				b.WriteString("\n")
			}
		}

	case TabLines:
		if len(m.typFile.Lines) == 0 {
			b.WriteString(statusStyle.Render("No lines defined"))
		} else {
			for i, line := range m.typFile.Lines {
				label := ""
				if len(line.Labels) > 0 {
					for _, l := range line.Labels {
						label = l
						break
					}
				}

				text := fmt.Sprintf("  %s - %s", line.Type, label)
				if i == m.selectedIdx {
					text = selectedStyle.Render("▸ " + line.Type + " - " + label)
				}
				b.WriteString(text)
				b.WriteString("\n")
			}
		}

	case TabPolygons:
		if len(m.typFile.Polygons) == 0 {
			b.WriteString(statusStyle.Render("No polygons defined"))
		} else {
			for i, polygon := range m.typFile.Polygons {
				label := ""
				if len(polygon.Labels) > 0 {
					for _, l := range polygon.Labels {
						label = l
						break
					}
				}

				line := fmt.Sprintf("  %s - %s", polygon.Type, label)
				if i == m.selectedIdx {
					line = selectedStyle.Render("▸ " + polygon.Type + " - " + label)
				}
				b.WriteString(line)
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

// renderFooter renders the footer with help text
func (m Model) renderFooter() string {
	footer := "[Tab] Switch  [↑/↓] Navigate  [Enter] Details  [Ctrl+S] Save  [?] Help  [q] Quit"

	// Show status message if present
	if m.status != "" {
		statusMsg := statusStyle.Render(m.status)
		return statusMsg + "\n" + helpStyle.Render(footer)
	}

	return helpStyle.Render(footer)
}

// viewDetail renders the detail view for a selected item
func (m Model) viewDetail() string {
	if m.typFile == nil {
		return "No file loaded"
	}

	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")

	// Tabs
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	// Detail content based on active tab
	switch m.activeTab {
	case TabPoints:
		if m.selectedIdx < len(m.typFile.Points) {
			b.WriteString(m.renderPointDetail(m.typFile.Points[m.selectedIdx]))
		}
	case TabLines:
		if m.selectedIdx < len(m.typFile.Lines) {
			b.WriteString(m.renderLineDetail(m.typFile.Lines[m.selectedIdx]))
		}
	case TabPolygons:
		if m.selectedIdx < len(m.typFile.Polygons) {
			b.WriteString(m.renderPolygonDetail(m.typFile.Polygons[m.selectedIdx]))
		}
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("[e] Edit  [x] Edit XPM  [Esc] Back  [?] Help  [q] Quit"))

	return b.String()
}

// viewEdit renders the edit form for the selected item
func (m Model) viewEdit() string {
	if m.typFile == nil {
		return "No file loaded"
	}

	var b strings.Builder

	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")

	// Tabs
	b.WriteString(m.renderTabs())
	b.WriteString("\n\n")

	// Edit form title
	var itemType string
	switch m.activeTab {
	case TabPoints:
		itemType = "Point"
	case TabLines:
		itemType = "Line"
	case TabPolygons:
		itemType = "Polygon"
	}
	b.WriteString(titleStyle.Render(fmt.Sprintf("Edit %s", itemType)))
	b.WriteString("\n\n")

	// Render form fields
	for i, input := range m.inputs {
		b.WriteString(input.View())
		b.WriteString("\n")
		if i < len(m.inputs)-1 {
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")
	b.WriteString(helpStyle.Render("[Ctrl+S] Save  [Esc] Cancel  [Tab/↑/↓] Navigate fields"))

	return b.String()
}

// viewConfirmQuit renders the confirmation dialog for quitting with unsaved changes
func (m Model) viewConfirmQuit() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Unsaved Changes"))
	b.WriteString("\n\n")
	b.WriteString("You have unsaved changes. What would you like to do?\n\n")
	b.WriteString("  [Y] Save and quit\n")
	b.WriteString("  [N] Quit without saving\n")
	b.WriteString("  [Esc/C] Cancel and return\n")

	return b.String()
}

// viewEditXPM renders the XPM editor with a scrollable viewport
func (m Model) viewEditXPM() string {
	if m.editingXPM == nil {
		return "No XPM data"
	}

	var b strings.Builder

	// Header (fixed, not scrollable)
	b.WriteString(titleStyle.Render(fmt.Sprintf("Edit %s", m.editingXPMType)))
	b.WriteString("\n\n")

	// Scrollable content area (content already set in viewport via updateXPMViewportContent)
	b.WriteString(m.xpmViewport.View())
	b.WriteString("\n")

	// Footer (fixed, not scrollable)
	if len(m.inputs) > 0 {
		b.WriteString(selectedStyle.Render("Edit Color:"))
		b.WriteString("\n")
		b.WriteString(m.inputs[0].View())
		b.WriteString("\n\n")
		b.WriteString(helpStyle.Render("[Enter] Save Color  [Esc] Cancel"))
	} else {
		b.WriteString(helpStyle.Render("[↑/↓/PgUp/PgDn] Scroll  [Tab] Navigate Colors  [Enter] Edit  [Esc] Back  [Ctrl+S] Save"))
	}

	return b.String()
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// renderPointDetail renders the details of a point type
func (m Model) renderPointDetail(point parser.PointType) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Point Details"))
	b.WriteString("\n\n")

	// Type
	b.WriteString(selectedStyle.Render("Type: "))
	b.WriteString(point.Type)
	b.WriteString("\n\n")

	// Labels
	if len(point.Labels) > 0 {
		b.WriteString(selectedStyle.Render("Labels:"))
		b.WriteString("\n")
		for code, label := range point.Labels {
			langName := getLanguageName(code)
			b.WriteString(fmt.Sprintf("  %s (%s): %s\n", code, langName, label))
		}
		b.WriteString("\n")
	}

	// DayColors
	if len(point.DayColors) > 0 {
		b.WriteString(selectedStyle.Render("Day Colors:"))
		b.WriteString("\n")
		for i, color := range point.DayColors {
			colorDisplay := renderColorWithPreview(color.Hex)
			name := color.Name
			if name == "" {
				name = fmt.Sprintf("Color %d", i+1)
			}
			b.WriteString(fmt.Sprintf("  %s: %s\n", name, colorDisplay))
		}
		b.WriteString("\n")
	}

	// NightColors
	if len(point.NightColors) > 0 {
		b.WriteString(selectedStyle.Render("Night Colors:"))
		b.WriteString("\n")
		for i, color := range point.NightColors {
			colorDisplay := renderColorWithPreview(color.Hex)
			name := color.Name
			if name == "" {
				name = fmt.Sprintf("Color %d", i+1)
			}
			b.WriteString(fmt.Sprintf("  %s: %s\n", name, colorDisplay))
		}
		b.WriteString("\n")
	}

	// DayXpm
	if point.DayXpm != nil {
		b.WriteString(selectedStyle.Render("Day Icon:"))
		b.WriteString("\n")
		b.WriteString(m.renderXPMInfo(point.DayXpm))
		b.WriteString("\n")
	}

	// NightXpm
	if point.NightXpm != nil {
		b.WriteString(selectedStyle.Render("Night Icon:"))
		b.WriteString("\n")
		b.WriteString(m.renderXPMInfo(point.NightXpm))
		b.WriteString("\n")
	}

	return b.String()
}

// renderLineDetail renders the details of a line type
func (m Model) renderLineDetail(line parser.LineType) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Line Details"))
	b.WriteString("\n\n")

	// Type
	b.WriteString(selectedStyle.Render("Type: "))
	b.WriteString(line.Type)
	b.WriteString("\n\n")

	// Labels
	if len(line.Labels) > 0 {
		b.WriteString(selectedStyle.Render("Labels:"))
		b.WriteString("\n")
		for code, label := range line.Labels {
			langName := getLanguageName(code)
			b.WriteString(fmt.Sprintf("  %s (%s): %s\n", code, langName, label))
		}
		b.WriteString("\n")
	}

	// Line properties
	if line.LineWidth > 0 {
		b.WriteString(fmt.Sprintf("Line Width: %d\n", line.LineWidth))
	}
	if line.BorderWidth > 0 {
		b.WriteString(fmt.Sprintf("Border Width: %d\n", line.BorderWidth))
	}
	if line.LineStyle != "" {
		b.WriteString(fmt.Sprintf("Line Style: %s\n", line.LineStyle))
	}
	b.WriteString("\n")

	// DayXpm
	if line.DayXpm != nil {
		b.WriteString(selectedStyle.Render("Day Pattern:"))
		b.WriteString("\n")
		b.WriteString(m.renderXPMInfo(line.DayXpm))
		b.WriteString("\n")
	}

	// NightXpm
	if line.NightXpm != nil {
		b.WriteString(selectedStyle.Render("Night Pattern:"))
		b.WriteString("\n")
		b.WriteString(m.renderXPMInfo(line.NightXpm))
		b.WriteString("\n")
	}

	return b.String()
}

// renderPolygonDetail renders the details of a polygon type
func (m Model) renderPolygonDetail(polygon parser.PolygonType) string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Polygon Details"))
	b.WriteString("\n\n")

	// Type
	b.WriteString(selectedStyle.Render("Type: "))
	b.WriteString(polygon.Type)
	b.WriteString("\n\n")

	// Labels
	if len(polygon.Labels) > 0 {
		b.WriteString(selectedStyle.Render("Labels:"))
		b.WriteString("\n")
		for code, label := range polygon.Labels {
			langName := getLanguageName(code)
			b.WriteString(fmt.Sprintf("  %s (%s): %s\n", code, langName, label))
		}
		b.WriteString("\n")
	}

	// Extended labels
	if polygon.ExtendedLabels {
		b.WriteString("Extended Labels: Yes\n\n")
	}

	// DayXpm
	if polygon.DayXpm != nil {
		b.WriteString(selectedStyle.Render("Day Pattern:"))
		b.WriteString("\n")
		b.WriteString(m.renderXPMInfo(polygon.DayXpm))
		b.WriteString("\n")
	}

	// NightXpm
	if polygon.NightXpm != nil {
		b.WriteString(selectedStyle.Render("Night Pattern:"))
		b.WriteString("\n")
		b.WriteString(m.renderXPMInfo(polygon.NightXpm))
		b.WriteString("\n")
	}

	return b.String()
}

// renderXPMInfo renders information about an XPM icon
func (m Model) renderXPMInfo(xpm *parser.XPMIcon) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("  Size: %dx%d\n", xpm.Width, xpm.Height))
	b.WriteString(fmt.Sprintf("  Colors: %d\n", xpm.Colors))
	b.WriteString(fmt.Sprintf("  Chars per pixel: %d\n", xpm.CharsPerPixel))

	if len(xpm.Palette) > 0 {
		b.WriteString("  Color Palette:\n")

		// Convert to sorted slice for consistent ordering
		type colorEntry struct {
			char  string
			color parser.Color
		}
		var colors []colorEntry
		for char, color := range xpm.Palette {
			colors = append(colors, colorEntry{char, color})
		}

		// Sort alphabetically by character
		sort.Slice(colors, func(i, j int) bool {
			return colors[i].char < colors[j].char
		})

		for _, entry := range colors {
			colorDisplay := renderColorWithPreview(entry.color.Hex)
			b.WriteString(fmt.Sprintf("    %s → %s\n", entry.char, colorDisplay))
		}
	}

	// Render the icon preview with colors
	if len(xpm.Data) > 0 {
		b.WriteString("\n  Icon Preview:\n")
		b.WriteString(renderXPMPreview(xpm))
	}

	return b.String()
}

// renderXPMPreview renders the XPM pixel data with colors applied
func renderXPMPreview(xpm *parser.XPMIcon) string {
	var b strings.Builder

	maxRows := 10 // Show up to 10 rows
	rowsToShow := min(len(xpm.Data), maxRows)

	for i := 0; i < rowsToShow; i++ {
		b.WriteString("  ")
		row := xpm.Data[i]

		// Process each character/pixel in the row
		for _, char := range row {
			charStr := string(char)

			// Look up the color for this character
			if color, ok := xpm.Palette[charStr]; ok {
				b.WriteString(renderPixelWithColor(color.Hex, charStr))
			} else {
				// Unknown character, show as gray
				b.WriteString(renderPixelWithColor("#808080", charStr))
			}
		}

		// Reset color at end of line
		b.WriteString("\x1b[0m\n")
	}

	if len(xpm.Data) > maxRows {
		b.WriteString(fmt.Sprintf("  ... and %d more rows\n", len(xpm.Data)-maxRows))
	}

	return b.String()
}

// renderPixelWithColor renders a single pixel with the given color as background and the marker character
func renderPixelWithColor(hexColor string, char string) string {
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

// renderColorWithPreview renders a color hex code with a visual preview using ANSI codes
func renderColorWithPreview(hexColor string) string {
	if hexColor == "" || hexColor == "none" || hexColor == "transparent" {
		return hexColor
	}

	// Parse hex color (supports #RRGGBB or RRGGBB)
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

// getLanguageName returns a human-readable language name for a language code
func getLanguageName(code string) string {
	// Common language codes used in Garmin TYP files
	languages := map[string]string{
		"0x01": "French",
		"0x02": "German",
		"0x03": "Dutch",
		"0x04": "English",
		"0x05": "Italian",
		"0x06": "Finnish",
		"0x07": "Swedish",
		"0x08": "Spanish",
		"0x09": "Basque",
		"0x0a": "Catalan",
		"0x0b": "Galician",
		"0x0c": "Welsh",
		"0x0d": "Gaelic",
		"0x0e": "Danish",
		"0x0f": "Norwegian",
		"0x10": "Portuguese",
		"0x11": "Slovak",
		"0x12": "Czech",
		"0x13": "Croatian",
		"0x14": "Hungarian",
		"0x15": "Polish",
		"0x16": "Turkish",
		"0x17": "Greek",
		"0x18": "Slovenian",
		"0x19": "Russian",
		"0x1a": "Estonian",
		"0x1b": "Latvian",
		"0x1c": "Romanian",
		"0x1d": "Albanian",
		"0x1e": "Bosnian",
		"0x1f": "Lithuanian",
		"0x20": "Serbian",
		"0x21": "Macedonian",
		"0x22": "Bulgarian",
	}

	if name, ok := languages[code]; ok {
		return name
	}
	return "Unknown"
}
