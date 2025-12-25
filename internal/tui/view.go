package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
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
	b.WriteString("  q, Ctrl+C    Quit\n")
	b.WriteString("  ?            Toggle help\n")
	b.WriteString("  Tab          Switch between tabs (Points/Lines/Polygons)\n")
	b.WriteString("  ↑/k          Move up\n")
	b.WriteString("  ↓/j          Move down\n")
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
	return helpStyle.Render("[Tab] Switch  [↑/↓] Navigate  [?] Help  [q] Quit")
}
