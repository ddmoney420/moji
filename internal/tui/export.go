// export.go handles export modal UI, file browser navigation, and export operations
// for PNG, SVG, HTML, and TXT formats.
package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddmoney420/moji/internal/export"
)

// ExportModal handles file path selection and export operations
type ExportModal struct {
	active      bool
	format      ExportFormat
	input       textinput.Model
	currentDir  string
	files       []os.DirEntry
	selectedIdx int
	showBrowser bool
	message     string
}

// openExportModal initializes and displays the export modal
func (m *Model) openExportModal() {
	m.exportModal.active = true
	m.exportModal.format = exportPNG
	m.exportModal.input.SetValue(generateFilename("png"))
	m.exportModal.input.Focus()
	m.exportModal.showBrowser = false
	m.loadDirectory()
}

// loadDirectory loads files from the current directory
func (m *Model) loadDirectory() {
	files, err := os.ReadDir(m.exportModal.currentDir)
	if err != nil {
		m.exportModal.files = nil
		return
	}
	m.exportModal.files = files
	m.exportModal.selectedIdx = 0
}

// handleExportModalKey processes key input in the export modal
func (m Model) handleExportModalKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.exportModal.active = false
		return m, nil

	case "tab":
		m.exportModal.showBrowser = !m.exportModal.showBrowser
		if m.exportModal.showBrowser {
			m.exportModal.input.Blur()
		} else {
			m.exportModal.input.Focus()
		}
		return m, nil

	case "1":
		m.exportModal.format = exportPNG
		m.exportModal.input.SetValue(changeExtension(m.exportModal.input.Value(), "png"))
		return m, nil
	case "2":
		m.exportModal.format = exportSVG
		m.exportModal.input.SetValue(changeExtension(m.exportModal.input.Value(), "svg"))
		return m, nil
	case "3":
		m.exportModal.format = exportHTML
		m.exportModal.input.SetValue(changeExtension(m.exportModal.input.Value(), "html"))
		return m, nil
	case "4":
		m.exportModal.format = exportTXT
		m.exportModal.input.SetValue(changeExtension(m.exportModal.input.Value(), "txt"))
		return m, nil

	case "enter":
		if m.exportModal.showBrowser {
			// Select directory or file
			if len(m.exportModal.files) > 0 && m.exportModal.selectedIdx < len(m.exportModal.files) {
				selected := m.exportModal.files[m.exportModal.selectedIdx]
				if selected.IsDir() {
					m.exportModal.currentDir = filepath.Join(m.exportModal.currentDir, selected.Name())
					m.loadDirectory()
				} else {
					m.exportModal.input.SetValue(filepath.Join(m.exportModal.currentDir, selected.Name()))
					m.exportModal.showBrowser = false
					m.exportModal.input.Focus()
				}
			}
		} else {
			// Do export
			m.doExport()
		}
		return m, nil

	case "backspace":
		if m.exportModal.showBrowser {
			// Go up one directory
			parent := filepath.Dir(m.exportModal.currentDir)
			if parent != m.exportModal.currentDir {
				m.exportModal.currentDir = parent
				m.loadDirectory()
			}
			return m, nil
		}

	case "up", "k":
		if m.exportModal.showBrowser && m.exportModal.selectedIdx > 0 {
			m.exportModal.selectedIdx--
		}
		return m, nil

	case "down", "j":
		if m.exportModal.showBrowser && m.exportModal.selectedIdx < len(m.exportModal.files)-1 {
			m.exportModal.selectedIdx++
		}
		return m, nil
	}

	// Handle text input
	if !m.exportModal.showBrowser {
		var cmd tea.Cmd
		m.exportModal.input, cmd = m.exportModal.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

// doExport performs the actual file export
func (m *Model) doExport() {
	filename := m.exportModal.input.Value()
	if filename == "" {
		m.exportModal.message = "No filename entered"
		return
	}

	// Make absolute if relative
	if !filepath.IsAbs(filename) {
		filename = filepath.Join(m.exportModal.currentDir, filename)
	}

	plain := stripANSI(m.preview)
	var err error

	switch m.exportModal.format {
	case exportPNG:
		err = export.ToPNG(plain, filename, "#282a36", "#f8f8f2")
	case exportSVG:
		err = export.ToSVG(plain, filename, "#282a36", "#f8f8f2")
	case exportHTML:
		err = export.ToHTML(plain, filename, "#282a36", "#f8f8f2", m.textInput.Value())
	case exportTXT:
		err = os.WriteFile(filename, []byte(plain), 0644)
	}

	if err != nil {
		m.statusMsg = fmt.Sprintf("Error: %v", err)
	} else {
		m.statusMsg = fmt.Sprintf("Saved to %s", filename)
	}
	m.statusTime = time.Now()
	m.exportModal.active = false
}

// renderExportModal renders the export modal UI
func (m Model) renderExportModal() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Export"))
	b.WriteString("\n\n")

	// Format selection
	b.WriteString(subtitleStyle.Render("Format: "))
	formats := []struct {
		key  string
		name string
		ext  string
	}{
		{"1", "PNG", "png"},
		{"2", "SVG", "svg"},
		{"3", "HTML", "html"},
		{"4", "TXT", "txt"},
	}
	for i, f := range formats {
		style := dimStyle
		if i == int(m.exportModal.format-1) {
			style = selectedStyle
		}
		b.WriteString(style.Render(fmt.Sprintf("[%s]%s ", f.key, f.name)))
	}
	b.WriteString("\n\n")

	// Path input
	b.WriteString(subtitleStyle.Render("Filename: "))
	b.WriteString(m.exportModal.input.View())
	b.WriteString("\n\n")

	// Current directory
	b.WriteString(dimStyle.Render("Directory: "))
	b.WriteString(pathStyle.Render(m.exportModal.currentDir))
	b.WriteString("\n\n")

	// File browser (toggle with Tab)
	browserLabel := "[Tab] Show browser"
	if m.exportModal.showBrowser {
		browserLabel = "[Tab] Hide browser"
		b.WriteString(subtitleStyle.Render("Files:"))
		b.WriteString(dimStyle.Render(" [Backspace] Go up"))
		b.WriteString("\n")

		visibleCount := 8
		start := m.exportModal.selectedIdx - visibleCount/2
		if start < 0 {
			start = 0
		}
		end := start + visibleCount
		if end > len(m.exportModal.files) {
			end = len(m.exportModal.files)
		}

		for i := start; i < end; i++ {
			f := m.exportModal.files[i]
			icon := "  "
			if f.IsDir() {
				icon = "📁"
			}
			name := f.Name()
			style := dimStyle
			if i == m.exportModal.selectedIdx {
				style = selectedStyle
				name = "▸ " + name
			} else {
				name = "  " + name
			}
			b.WriteString(style.Render(fmt.Sprintf("%s %s", icon, name)))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(dimStyle.Render(browserLabel))
	b.WriteString("\n\n")

	// Preview thumbnail
	b.WriteString(subtitleStyle.Render("Preview:"))
	b.WriteString("\n")
	preview := m.preview
	lines := strings.Split(preview, "\n")
	if len(lines) > 5 {
		lines = lines[:5]
		lines = append(lines, "...")
	}
	b.WriteString(boxStyle.Render(strings.Join(lines, "\n")))
	b.WriteString("\n\n")

	b.WriteString(helpStyle.Render("[Enter] Save  [Esc] Cancel"))

	return modalStyle.Render(b.String())
}
