// render.go contains the main view rendering logic, including headers, tab lists,
// preview panes, status bars, help screens, and the Dracula color theme styling.
package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Theme colors - Dracula theme
var (
	purple     = lipgloss.Color("#BD93F9")
	pink       = lipgloss.Color("#FF79C6")
	green      = lipgloss.Color("#50FA7B")
	cyan       = lipgloss.Color("#8BE9FD")
	red        = lipgloss.Color("#FF5555")
	yellow     = lipgloss.Color("#F1FA8C")
	comment    = lipgloss.Color("#6272A4")
	foreground = lipgloss.Color("#F8F8F2")
	background = lipgloss.Color("#282A36")
	selection  = lipgloss.Color("#44475A")
)

// Styles used for rendering UI elements
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(pink).
			Background(selection).
			Padding(0, 2)

	tabStyle = lipgloss.NewStyle().
			Foreground(comment).
			Padding(0, 1)

	activeTabStyle = lipgloss.NewStyle().
			Foreground(pink).
			Bold(true).
			Padding(0, 1).
			Underline(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true)

	selectedStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(foreground)

	dimStyle = lipgloss.NewStyle().
			Foreground(comment)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple).
			Padding(0, 1)

	previewBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(pink).
			Padding(1, 2)

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(cyan).
			Padding(1, 2).
			Background(background)

	helpStyle = lipgloss.NewStyle().
			Foreground(comment)

	copiedStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(red).
			Bold(true)

	keyStyle = lipgloss.NewStyle().
			Foreground(yellow).
			Bold(true)

	pathStyle = lipgloss.NewStyle().
			Foreground(cyan)
)

// View renders the main TUI view
func (m Model) View() string {
	if m.showHelp {
		return m.renderHelp()
	}

	if m.exportModal.active {
		return m.renderExportModal()
	}

	var b strings.Builder

	// Title and tabs
	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	// Tab content
	b.WriteString(m.renderTabContent())
	b.WriteString("\n")

	// Preview
	b.WriteString(m.renderPreview())
	b.WriteString("\n")

	// Status bar
	b.WriteString(m.renderStatusBar())

	return b.String()
}

// renderHeader renders the title and tab selector
func (m Model) renderHeader() string {
	title := titleStyle.Render("MOJI Studio")

	var tabs strings.Builder
	for i, name := range tabNames {
		if i == m.currentTab {
			tabs.WriteString(activeTabStyle.Render(fmt.Sprintf("[%d]%s", i+1, name)))
		} else {
			tabs.WriteString(tabStyle.Render(fmt.Sprintf("[%d]%s", i+1, name)))
		}
		tabs.WriteString(" ")
	}

	return lipgloss.JoinVertical(lipgloss.Left, title, tabs.String())
}

// renderList renders a filterable list with selection
func (m Model) renderList(title string, indices []int, selected int, getLabel func(int) string, focused bool) string {
	var b strings.Builder

	var titleStr string
	if focused {
		titleStr = selectedStyle.Render("▸ " + title)
	} else {
		titleStr = dimStyle.Render("  " + title)
	}
	b.WriteString(titleStr)
	b.WriteString("\n")

	visibleCount := 8
	start := selected - visibleCount/2
	if start < 0 {
		start = 0
	}
	end := start + visibleCount
	if end > len(indices) {
		end = len(indices)
		start = end - visibleCount
		if start < 0 {
			start = 0
		}
	}

	for i := start; i < end; i++ {
		label := getLabel(i)
		if i == selected {
			if focused {
				b.WriteString(selectedStyle.Render(fmt.Sprintf("  ▸ %-20s", label)))
			} else {
				b.WriteString(normalStyle.Render(fmt.Sprintf("  ▸ %-20s", label)))
			}
		} else {
			b.WriteString(dimStyle.Render(fmt.Sprintf("    %-20s", label)))
		}
		b.WriteString("\n")
	}

	if len(indices) > visibleCount {
		b.WriteString(dimStyle.Render(fmt.Sprintf("    [%d/%d]", selected+1, len(indices))))
	}

	return b.String()
}

// renderPreview renders the preview section
func (m Model) renderPreview() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Preview"))
	b.WriteString("\n")

	preview := m.preview
	if preview == "" {
		preview = dimStyle.Render("(no preview)")
	}

	// Limit preview height
	lines := strings.Split(preview, "\n")
	maxLines := 15
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, dimStyle.Render("... (truncated)"))
		preview = strings.Join(lines, "\n")
	}

	maxWidth := m.width - 10
	if maxWidth < 60 {
		maxWidth = 60
	}
	if maxWidth > 120 {
		maxWidth = 120
	}

	b.WriteString(previewBoxStyle.MaxWidth(maxWidth).Render(preview))

	return b.String()
}

// renderStatusBar renders the status and help text at the bottom
func (m Model) renderStatusBar() string {
	var status string

	if m.statusMsg != "" && time.Since(m.statusTime) < 3*time.Second {
		status = copiedStyle.Render(m.statusMsg)
	}

	help := helpStyle.Render("[1-9] Tabs  [Tab] Next  [↑↓←→] Navigate  [Enter] Copy  [e] Export  [?] Help  [q] Quit")

	if status != "" {
		return status + "\n" + help
	}
	return help
}

// renderHelp renders the help screen
func (m Model) renderHelp() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("MOJI Studio Help"))
	b.WriteString("\n\n")

	helpItems := []struct{ key, desc string }{
		{"1-9", "Switch to tab by number"},
		{"Tab / Shift+Tab", "Next/previous tab"},
		{"↑/↓ or j/k", "Navigate lists"},
		{"←/→ or h/l", "Switch between lists / toggle options"},
		{"Enter", "Copy to clipboard"},
		{"e or Ctrl+E", "Export (PNG, SVG, HTML, TXT)"},
		{"Ctrl+F", "Focus search/filter input"},
		{"Esc", "Return to main input"},
		{"?", "Toggle help"},
		{"q / Ctrl+C", "Quit"},
	}

	for _, item := range helpItems {
		b.WriteString(fmt.Sprintf("  %s  %s\n",
			keyStyle.Render(fmt.Sprintf("%-18s", item.key)),
			item.desc))
	}

	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render("Tabs"))
	b.WriteString("\n\n")

	tabDescs := []string{
		"Banner - ASCII art text with fonts & styles",
		"Kaomoji - Japanese emoticons",
		"ArtDB - ASCII art database",
		"Filters - Color filters (rainbow, fire, neon...)",
		"Effects - Text effects (flip, zalgo, bubble...)",
		"Gradient - Color gradient themes",
		"QR - QR code generator",
		"Patterns - Borders and dividers",
		"Speech - Speech bubbles with characters",
		"Calendar - Calendar views (Tab to reach)",
		"Sysinfo - System info display (Tab to reach)",
	}

	for i, desc := range tabDescs {
		b.WriteString(fmt.Sprintf("  %s %s\n", keyStyle.Render(fmt.Sprintf("[%d]", i+1)), desc))
	}

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("Press ? or Esc to close"))

	return modalStyle.Render(b.String())
}
