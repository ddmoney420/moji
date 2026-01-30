// Package tabs contains all tab-specific rendering logic for the TUI.
// Each of the 11 feature tabs (Banner, Kaomoji, ArtDB, Filters, Effects, Gradient,
// QR, Patterns, Speech, Calendar, Sysinfo) has dedicated update and render functions.
// This separates preview generation and tab-specific UI from the core Model logic.
package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/ddmoney420/moji/internal/artdb"
	"github.com/ddmoney420/moji/internal/banner"
	"github.com/ddmoney420/moji/internal/calendar"
	"github.com/ddmoney420/moji/internal/effects"
	"github.com/ddmoney420/moji/internal/filters"
	"github.com/ddmoney420/moji/internal/gradient"
	"github.com/ddmoney420/moji/internal/patterns"
	"github.com/ddmoney420/moji/internal/qrcode"
	"github.com/ddmoney420/moji/internal/speech"
	"github.com/ddmoney420/moji/internal/styles"
	"github.com/ddmoney420/moji/internal/sysinfo"
)

// updatePreview generates the preview for the current tab
func (m *Model) updatePreview() {
	text := m.textInput.Value()
	if text == "" {
		text = "MOJI"
	}

	switch m.currentTab {
	case int(tabBanner):
		m.updateBannerPreview(text)
	case int(tabKaomoji):
		m.updateKaomojiPreview()
	case int(tabArtDB):
		m.updateArtDBPreview()
	case int(tabFilters):
		m.updateFiltersPreview(text)
	case int(tabEffects):
		m.updateEffectsPreview(text)
	case int(tabGradient):
		m.updateGradientPreview(text)
	case int(tabQRCode):
		m.updateQRPreview(text)
	case int(tabPatterns):
		m.updatePatternsPreview(text)
	case int(tabSpeech):
		m.updateSpeechPreview(text)
	case int(tabCalendar):
		m.updateCalendarPreview()
	case int(tabSysinfo):
		m.updateSysinfoPreview(text)
	}
}

func (m *Model) updateBannerPreview(text string) {
	if len(m.filteredFonts) > 0 {
		font := m.fonts[m.filteredFonts[m.selectedFont]]
		style := m.colorStyles[m.selectedStyle]
		art, err := banner.Generate(text, font)
		if err == nil {
			m.preview = styles.Apply(art, style)
		} else {
			m.preview = errorStyle.Render(err.Error())
		}
	}
}

func (m *Model) updateKaomojiPreview() {
	if len(m.filteredKaomoji) > 0 {
		k := m.kaomojiList[m.filteredKaomoji[m.selectedKaomoji]]
		m.preview = fmt.Sprintf("%s\n\n%s", selectedStyle.Render(k.Kaomoji), dimStyle.Render(k.Name))
	}
}

func (m *Model) updateArtDBPreview() {
	if len(m.filteredArt) > 0 {
		a := m.artList[m.filteredArt[m.selectedArt]]
		m.preview = a.Art
	}
}

func (m *Model) updateFiltersPreview(text string) {
	filter := m.filterList[m.selectedFilter]
	if f, ok := filters.Get(filter); ok {
		m.preview = f(text)
	}
}

func (m *Model) updateEffectsPreview(text string) {
	effect := m.effectList[m.selectedEffect]
	m.preview = effects.Apply(effect, text)
}

func (m *Model) updateGradientPreview(text string) {
	modes := []string{"horizontal", "vertical", "diagonal"}
	mode := modes[m.gradientMode]
	theme := m.gradientThemes[m.selectedTheme]
	m.preview = gradient.Apply(text, theme, mode)
}

func (m *Model) updateQRPreview(text string) {
	charset := m.qrCharsets[m.selectedQRCharset]
	qr, err := qrcode.Generate(text, qrcode.Options{Charset: charset, Invert: m.qrInvert})
	if err == nil {
		m.preview = qr
	} else {
		m.preview = errorStyle.Render(err.Error())
	}
}

func (m *Model) updatePatternsPreview(text string) {
	if m.patternMode == 0 {
		border := m.patternBorders[m.selectedBorder]
		m.preview = patterns.CreateBorder(text, border, 1)
	} else {
		divider := m.patternDividers[m.selectedDivider]
		m.preview = patterns.CreateDivider(divider, 40)
	}
}

func (m *Model) updateSpeechPreview(text string) {
	style := m.speechStyles[m.selectedSpeech]
	bubble := speech.Wrap(text, style, 40)
	if m.selectedSpeechArt > 0 {
		artName := m.speechArts[m.selectedSpeechArt]
		if a, ok := artdb.Get(artName); ok {
			bubble = speech.Combine(bubble, a.Art)
		}
	}
	m.preview = bubble
}

func (m *Model) updateCalendarPreview() {
	opts := calendar.Options{FirstDayMon: m.calMonday}
	switch m.calendarMode {
	case 0:
		m.preview = calendar.Current(opts)
	case 1:
		m.preview = calendar.WeekView(opts)
	case 2:
		m.preview = calendar.Year(time.Now().Year(), opts)
	}
}

func (m *Model) updateSysinfoPreview(text string) {
	info := sysinfo.Collect()
	art := sysinfo.GetOSLogo()
	if text != "MOJI" {
		art = gradient.Apply(art, text, "diagonal")
	}
	m.preview = info.FormatWithArt(art)
}

// renderTabContent dispatches to tab-specific rendering
func (m Model) renderTabContent() string {
	switch m.currentTab {
	case int(tabBanner):
		return m.renderBannerTab()
	case int(tabKaomoji):
		return m.renderKaomojiTab()
	case int(tabArtDB):
		return m.renderArtDBTab()
	case int(tabFilters):
		return m.renderFiltersTab()
	case int(tabEffects):
		return m.renderEffectsTab()
	case int(tabGradient):
		return m.renderGradientTab()
	case int(tabQRCode):
		return m.renderQRTab()
	case int(tabPatterns):
		return m.renderPatternsTab()
	case int(tabSpeech):
		return m.renderSpeechTab()
	case int(tabCalendar):
		return m.renderCalendarTab()
	case int(tabSysinfo):
		return m.renderSysinfoTab()
	default:
		return ""
	}
}

func (m Model) renderBannerTab() string {
	var b strings.Builder

	// Text input
	b.WriteString(subtitleStyle.Render("Text: "))
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	// Font and style selectors side by side
	fontList := m.renderList("Font", m.filteredFonts, m.selectedFont, func(i int) string {
		return m.fonts[m.filteredFonts[i]]
	}, m.focus == focusInput+1)

	styleList := m.renderList("Style", makeRange(len(m.colorStyles)), m.selectedStyle, func(i int) string {
		return m.colorStyles[i]
	}, m.focus == focusInput+2)

	b.WriteString(fontList)
	b.WriteString("  ")
	b.WriteString(styleList)

	return b.String()
}

func (m Model) renderKaomojiTab() string {
	var b strings.Builder

	// Search
	b.WriteString(subtitleStyle.Render("Search: "))
	b.WriteString(m.searchInput.View())
	b.WriteString(dimStyle.Render(fmt.Sprintf(" (%d results)", len(m.filteredKaomoji))))
	b.WriteString("\n\n")

	// Kaomoji list and categories
	kaomojiList := m.renderList("Kaomoji", m.filteredKaomoji, m.selectedKaomoji, func(i int) string {
		k := m.kaomojiList[m.filteredKaomoji[i]]
		return fmt.Sprintf("%-12s %s", k.Name, k.Kaomoji)
	}, m.focus == focusInput+1)

	catList := m.renderList("Category", makeRange(len(m.kaomojiCats)), m.selectedCat, func(i int) string {
		return m.kaomojiCats[i]
	}, m.focus == focusInput+2)

	b.WriteString(kaomojiList)
	b.WriteString("  ")
	b.WriteString(catList)

	return b.String()
}

func (m Model) renderArtDBTab() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Search: "))
	b.WriteString(m.searchInput.View())
	b.WriteString(dimStyle.Render(fmt.Sprintf(" (%d results)", len(m.filteredArt))))
	b.WriteString("\n\n")

	artList := m.renderList("Art", m.filteredArt, m.selectedArt, func(i int) string {
		a := m.artList[m.filteredArt[i]]
		return fmt.Sprintf("%-15s [%s]", a.Name, a.Category)
	}, m.focus == focusInput+1)

	catList := m.renderList("Category", makeRange(len(m.artCats)), m.selectedArtCat, func(i int) string {
		return m.artCats[i]
	}, m.focus == focusInput+2)

	b.WriteString(artList)
	b.WriteString("  ")
	b.WriteString(catList)

	return b.String()
}

func (m Model) renderFiltersTab() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Text: "))
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	filterList := m.renderList("Filter", makeRange(len(m.filterList)), m.selectedFilter, func(i int) string {
		return m.filterList[i]
	}, m.focus == focusInput+1)

	b.WriteString(filterList)

	return b.String()
}

func (m Model) renderEffectsTab() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Text: "))
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	effectList := m.renderList("Effect", makeRange(len(m.effectList)), m.selectedEffect, func(i int) string {
		return m.effectList[i]
	}, m.focus == focusInput+1)

	b.WriteString(effectList)

	return b.String()
}

func (m Model) renderGradientTab() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Text: "))
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	// Mode selector
	modes := []string{"Horizontal", "Vertical", "Diagonal"}
	b.WriteString(dimStyle.Render("Mode: "))
	for i, mode := range modes {
		if i == m.gradientMode {
			b.WriteString(selectedStyle.Render(fmt.Sprintf("[%s]", mode)))
		} else {
			b.WriteString(dimStyle.Render(fmt.Sprintf(" %s ", mode)))
		}
	}
	b.WriteString(dimStyle.Render("  [←/→ to change]"))
	b.WriteString("\n\n")

	themeList := m.renderList("Theme", makeRange(len(m.gradientThemes)), m.selectedTheme, func(i int) string {
		return m.gradientThemes[i]
	}, m.focus == focusInput+1)

	b.WriteString(themeList)

	return b.String()
}

func (m Model) renderSpeechTab() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Text: "))
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	styleList := m.renderList("Bubble Style", makeRange(len(m.speechStyles)), m.selectedSpeech, func(i int) string {
		return m.speechStyles[i]
	}, m.focus == focusInput+1)

	artList := m.renderList("Character", makeRange(len(m.speechArts)), m.selectedSpeechArt, func(i int) string {
		return m.speechArts[i]
	}, m.focus == focusInput+2)

	b.WriteString(styleList)
	b.WriteString("  ")
	b.WriteString(artList)

	return b.String()
}

func (m Model) renderQRTab() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Text/URL: "))
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	charsetList := m.renderList("Charset", makeRange(len(m.qrCharsets)), m.selectedQRCharset, func(i int) string {
		return m.qrCharsets[i]
	}, m.focus == focusInput+1)

	invertStr := "Normal"
	if m.qrInvert {
		invertStr = "Inverted"
	}
	options := boxStyle.Render(fmt.Sprintf("Invert: %s\n[←/→] to toggle", invertStr))

	b.WriteString(charsetList)
	b.WriteString("  ")
	b.WriteString(options)

	return b.String()
}

func (m Model) renderPatternsTab() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Text: "))
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	modeStr := "Borders"
	if m.patternMode == 1 {
		modeStr = "Dividers"
	}
	b.WriteString(dimStyle.Render(fmt.Sprintf("Mode: %s [←/→ to switch]", modeStr)))
	b.WriteString("\n\n")

	if m.patternMode == 0 {
		borderList := m.renderList("Border", makeRange(len(m.patternBorders)), m.selectedBorder, func(i int) string {
			return m.patternBorders[i]
		}, true)
		b.WriteString(borderList)
	} else {
		dividerList := m.renderList("Divider", makeRange(len(m.patternDividers)), m.selectedDivider, func(i int) string {
			return m.patternDividers[i]
		}, true)
		b.WriteString(dividerList)
	}

	return b.String()
}

func (m Model) renderCalendarTab() string {
	var b strings.Builder

	modes := []string{"Month", "Week", "Year"}
	b.WriteString(subtitleStyle.Render("View: "))
	for i, mode := range modes {
		if i == m.calendarMode {
			b.WriteString(selectedStyle.Render(fmt.Sprintf("[%s]", mode)))
		} else {
			b.WriteString(dimStyle.Render(fmt.Sprintf(" %s ", mode)))
		}
	}
	b.WriteString(dimStyle.Render("  [↑/↓ to change]"))
	b.WriteString("\n\n")

	mondayStr := "Sunday"
	if m.calMonday {
		mondayStr = "Monday"
	}
	b.WriteString(dimStyle.Render(fmt.Sprintf("Week starts: %s", mondayStr)))

	return b.String()
}

func (m Model) renderSysinfoTab() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Gradient theme (optional): "))
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	themes := []string{"rainbow", "neon", "fire", "ocean", "sunset", "forest", "galaxy"}
	b.WriteString(dimStyle.Render("Themes: " + strings.Join(themes, ", ")))

	return b.String()
}
