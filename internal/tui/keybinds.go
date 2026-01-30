// keybinds.go contains all keyboard input handling, including tab navigation,
// focus state transitions, and directional movement through lists.
package tui

import (
	"os"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
)

// handleKeyboard processes key inputs for the main view
func (m *Model) handleKeyboard(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		// Send ourselves SIGINT to ensure clean exit
		if p, err := os.FindProcess(os.Getpid()); err == nil {
			_ = p.Signal(os.Interrupt)
		}
		return m, tea.Quit

	case "q":
		if m.focus != focusInput {
			return m, tea.Quit
		}

	case "?":
		m.showHelp = true
		return m, nil

	case "1", "2", "3", "4", "5", "6", "7", "8", "9":
		idx := int(msg.String()[0] - '1')
		if idx < int(tabCount) {
			m.goToTab(idx)
			m.updatePreview()
		}
		return m, nil

	case "tab":
		m.nextTab()
		m.updatePreview()
		return m, nil

	case "shift+tab":
		m.prevTab()
		m.updatePreview()
		return m, nil

	case "ctrl+e", "e":
		if m.focus != focusInput && m.preview != "" {
			m.openExportModal()
			return m, nil
		}

	case "enter":
		if m.preview != "" {
			return m.copyToClipboard()
		}
		return m, nil

	case "up", "k":
		m.handleUp()
		m.updatePreview()
		return m, nil

	case "down", "j":
		m.handleDown()
		m.updatePreview()
		return m, nil

	case "left", "h":
		m.handleLeft()
		m.updatePreview()
		return m, nil

	case "right", "l":
		m.handleRight()
		m.updatePreview()
		return m, nil

	case "ctrl+f":
		m.focus = focusInput
		m.searchInput.Focus()
		m.textInput.Blur()
		return m, nil
	}

	// Handle text input
	if m.focus == focusInput {
		var inputCmd tea.Cmd
		if m.searchInput.Focused() {
			m.searchInput, inputCmd = m.searchInput.Update(msg)
			m.filterCurrentList()
		} else {
			m.textInput, inputCmd = m.textInput.Update(msg)
		}
		m.updatePreview()
		return m, inputCmd
	}

	return m, nil
}

// handleUp processes up arrow / k key in navigation mode
func (m *Model) handleUp() {
	switch m.currentTab {
	case int(tabBanner):
		if m.focus == focusList1 && m.selectedFont > 0 {
			m.selectedFont--
		} else if m.focus == focusList2 && m.selectedStyle > 0 {
			m.selectedStyle--
		}
	case int(tabKaomoji):
		if m.focus == focusList1 && m.selectedKaomoji > 0 {
			m.selectedKaomoji--
		} else if m.focus == focusList2 && m.selectedCat > 0 {
			m.selectedCat--
			m.filterKaomoji()
		}
	case int(tabArtDB):
		if m.focus == focusList1 && m.selectedArt > 0 {
			m.selectedArt--
		} else if m.focus == focusList2 && m.selectedArtCat > 0 {
			m.selectedArtCat--
			m.filterArt()
		}
	case int(tabFilters):
		if m.focus == focusList1 && m.selectedFilter > 0 {
			m.selectedFilter--
		}
	case int(tabEffects):
		if m.focus == focusList1 && m.selectedEffect > 0 {
			m.selectedEffect--
		}
	case int(tabGradient):
		if m.focus == focusList1 && m.selectedTheme > 0 {
			m.selectedTheme--
		}
	case int(tabQRCode):
		if m.focus == focusList1 && m.selectedQRCharset > 0 {
			m.selectedQRCharset--
		}
	case int(tabPatterns):
		if m.patternMode == 0 && m.selectedBorder > 0 {
			m.selectedBorder--
		} else if m.patternMode == 1 && m.selectedDivider > 0 {
			m.selectedDivider--
		}
	case int(tabSpeech):
		if m.focus == focusList1 && m.selectedSpeech > 0 {
			m.selectedSpeech--
		} else if m.focus == focusList2 && m.selectedSpeechArt > 0 {
			m.selectedSpeechArt--
		}
	case int(tabCalendar):
		if m.calendarMode > 0 {
			m.calendarMode--
		}
	}
}

// handleDown processes down arrow / j key in navigation mode
func (m *Model) handleDown() {
	switch m.currentTab {
	case int(tabBanner):
		if m.focus == focusList1 && m.selectedFont < len(m.filteredFonts)-1 {
			m.selectedFont++
		} else if m.focus == focusList2 && m.selectedStyle < len(m.colorStyles)-1 {
			m.selectedStyle++
		}
	case int(tabKaomoji):
		if m.focus == focusList1 && m.selectedKaomoji < len(m.filteredKaomoji)-1 {
			m.selectedKaomoji++
		} else if m.focus == focusList2 && m.selectedCat < len(m.kaomojiCats)-1 {
			m.selectedCat++
			m.filterKaomoji()
		}
	case int(tabArtDB):
		if m.focus == focusList1 && m.selectedArt < len(m.filteredArt)-1 {
			m.selectedArt++
		} else if m.focus == focusList2 && m.selectedArtCat < len(m.artCats)-1 {
			m.selectedArtCat++
			m.filterArt()
		}
	case int(tabFilters):
		if m.focus == focusList1 && m.selectedFilter < len(m.filterList)-1 {
			m.selectedFilter++
		}
	case int(tabEffects):
		if m.focus == focusList1 && m.selectedEffect < len(m.effectList)-1 {
			m.selectedEffect++
		}
	case int(tabGradient):
		if m.focus == focusList1 && m.selectedTheme < len(m.gradientThemes)-1 {
			m.selectedTheme++
		}
	case int(tabQRCode):
		if m.focus == focusList1 && m.selectedQRCharset < len(m.qrCharsets)-1 {
			m.selectedQRCharset++
		}
	case int(tabPatterns):
		if m.patternMode == 0 && m.selectedBorder < len(m.patternBorders)-1 {
			m.selectedBorder++
		} else if m.patternMode == 1 && m.selectedDivider < len(m.patternDividers)-1 {
			m.selectedDivider++
		}
	case int(tabSpeech):
		if m.focus == focusList1 && m.selectedSpeech < len(m.speechStyles)-1 {
			m.selectedSpeech++
		} else if m.focus == focusList2 && m.selectedSpeechArt < len(m.speechArts)-1 {
			m.selectedSpeechArt++
		}
	case int(tabCalendar):
		if m.calendarMode < 2 {
			m.calendarMode++
		}
	}
}

// handleLeft processes left arrow / h key in navigation mode
func (m *Model) handleLeft() {
	if m.focus > focusInput {
		m.focus--
	}
	if m.currentTab == int(tabPatterns) {
		m.patternMode = 0
	}
	if m.currentTab == int(tabQRCode) {
		m.qrInvert = false
	}
	if m.currentTab == int(tabGradient) && m.gradientMode > 0 {
		m.gradientMode--
	}
}

// handleRight processes right arrow / l key in navigation mode
func (m *Model) handleRight() {
	m.focus++
	if m.focus > focusList2 {
		m.focus = focusList2
	}
	if m.currentTab == int(tabPatterns) {
		m.patternMode = 1
	}
	if m.currentTab == int(tabQRCode) {
		m.qrInvert = true
	}
	if m.currentTab == int(tabGradient) && m.gradientMode < 2 {
		m.gradientMode++
	}
}

// copyToClipboard copies the preview to clipboard
func (m *Model) copyToClipboard() (tea.Model, tea.Cmd) {
	plain := stripANSI(m.preview)
	if err := clipboard.WriteAll(plain); err == nil {
		m.statusMsg = "Copied to clipboard!"
		m.statusTime = time.Now()
	}
	return m, nil
}
