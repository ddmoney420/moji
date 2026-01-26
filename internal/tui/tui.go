package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ddmoney420/moji/internal/artdb"
	"github.com/ddmoney420/moji/internal/banner"
	"github.com/ddmoney420/moji/internal/calendar"
	"github.com/ddmoney420/moji/internal/effects"
	"github.com/ddmoney420/moji/internal/export"
	"github.com/ddmoney420/moji/internal/filters"
	"github.com/ddmoney420/moji/internal/gradient"
	"github.com/ddmoney420/moji/internal/kaomoji"
	"github.com/ddmoney420/moji/internal/patterns"
	"github.com/ddmoney420/moji/internal/qrcode"
	"github.com/ddmoney420/moji/internal/speech"
	"github.com/ddmoney420/moji/internal/styles"
	"github.com/ddmoney420/moji/internal/sysinfo"
)

// Tabs
const (
	tabBanner = iota
	tabKaomoji
	tabArtDB
	tabFilters
	tabEffects
	tabGradient
	tabQRCode
	tabPatterns
	tabSpeech
	tabCalendar
	tabSysinfo
	tabCount
)

var tabNames = []string{
	"Banner", "Kaomoji", "ArtDB", "Filters", "Effects", "Gradient", "QR", "Patterns", "Speech", "Calendar", "Sysinfo",
}

// Focus states
const (
	focusInput = iota
	focusList1
	focusList2
	focusList3
)

// Export format types
const (
	exportNone = iota
	exportPNG
	exportSVG
	exportHTML
	exportTXT
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

// Styles
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

// Model represents the TUI state
type Model struct {
	// Common
	currentTab  int
	focus       int
	width       int
	height      int
	statusMsg   string
	statusTime  time.Time
	showHelp    bool
	exportModal ExportModal

	// Inputs
	textInput   textinput.Model
	searchInput textinput.Model

	// Banner tab
	fonts         []string
	fontDescs     []string
	filteredFonts []int
	colorStyles   []string
	selectedFont  int
	selectedStyle int

	// Kaomoji tab
	kaomojiList     []kaomoji.KaomojiItem
	filteredKaomoji []int
	selectedKaomoji int
	kaomojiCats     []string
	selectedCat     int

	// ArtDB tab
	artList        []artdb.Art
	filteredArt    []int
	selectedArt    int
	artCats        []string
	selectedArtCat int

	// Filters tab
	filterList     []string
	selectedFilter int

	// Effects tab
	effectList     []string
	selectedEffect int

	// QR tab
	qrCharsets        []string
	selectedQRCharset int
	qrInvert          bool

	// Gradient tab
	gradientThemes []string
	selectedTheme  int
	gradientMode   int // 0=horizontal, 1=vertical, 2=diagonal

	// Patterns tab
	patternBorders  []string
	patternDividers []string
	selectedBorder  int
	selectedDivider int
	patternMode     int // 0=border, 1=divider

	// Speech tab
	speechStyles      []string
	selectedSpeech    int
	speechArts        []string
	selectedSpeechArt int

	// Calendar tab
	calendarMode int // 0=month, 1=week, 2=year
	calMonday    bool

	// Current preview
	preview string
}

// ExportModal handles file path selection
type ExportModal struct {
	active      bool
	format      int
	input       textinput.Model
	currentDir  string
	files       []os.DirEntry
	selectedIdx int
	showBrowser bool
	message     string
}

// NewModel creates a new TUI model
func NewModel() Model {
	// Text input
	ti := textinput.New()
	ti.Placeholder = "Type text here..."
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 40

	// Search input
	si := textinput.New()
	si.Placeholder = "Search..."
	si.CharLimit = 30
	si.Width = 25

	// Get fonts
	fontInfos := banner.ListFonts()
	fonts := make([]string, len(fontInfos))
	fontDescs := make([]string, len(fontInfos))
	filteredFonts := make([]int, len(fontInfos))
	for i, f := range fontInfos {
		fonts[i] = f.Name
		fontDescs[i] = f.Desc
		filteredFonts[i] = i
	}

	// Get styles
	styleInfos := styles.ListStyles()
	colorStyles := make([]string, len(styleInfos))
	for i, s := range styleInfos {
		colorStyles[i] = s.Name
	}

	// Kaomoji
	kaomojiList := kaomoji.List("", "")
	filteredKaomoji := make([]int, len(kaomojiList))
	for i := range kaomojiList {
		filteredKaomoji[i] = i
	}
	kaomojiCats := kaomoji.ListCategories()

	// ArtDB
	artList := artdb.List()
	filteredArt := make([]int, len(artList))
	for i := range artList {
		filteredArt[i] = i
	}
	artCats := artdb.ListCategories()

	// Filters
	filterInfos := filters.ListFilters()
	filterList := make([]string, len(filterInfos))
	for i, f := range filterInfos {
		filterList[i] = f.Name
	}

	// Effects
	effectInfos := effects.ListEffects()
	effectList := make([]string, len(effectInfos))
	for i, e := range effectInfos {
		effectList[i] = e.Name
	}

	// QR charsets
	qrCharsets := qrcode.ListCharsets()

	// Gradient themes
	gradientInfos := gradient.ListThemes()
	gradientThemes := make([]string, len(gradientInfos))
	for i, g := range gradientInfos {
		gradientThemes[i] = g.Name
	}

	// Speech
	speechStyles := []string{"round", "square", "double", "thick", "ascii", "think"}
	speechArts := []string{"(none)"}
	for _, a := range artdb.List() {
		speechArts = append(speechArts, a.Name)
	}

	// Patterns
	patternBorders := patterns.ListBorders()
	patternDividers := patterns.ListDividers()

	// Export modal
	ei := textinput.New()
	ei.Placeholder = "filename.png"
	ei.CharLimit = 100
	ei.Width = 50

	cwd, _ := os.Getwd()

	m := Model{
		textInput:       ti,
		searchInput:     si,
		fonts:           fonts,
		fontDescs:       fontDescs,
		filteredFonts:   filteredFonts,
		colorStyles:     colorStyles,
		kaomojiList:     kaomojiList,
		filteredKaomoji: filteredKaomoji,
		kaomojiCats:     append([]string{"all"}, kaomojiCats...),
		artList:         artList,
		filteredArt:     filteredArt,
		artCats:         append([]string{"all"}, artCats...),
		filterList:      filterList,
		effectList:      effectList,
		qrCharsets:      qrCharsets,
		gradientThemes:  gradientThemes,
		speechStyles:    speechStyles,
		speechArts:      speechArts,
		patternBorders:  patternBorders,
		patternDividers: patternDividers,
		width:           120,
		height:          40,
		exportModal: ExportModal{
			input:      ei,
			currentDir: cwd,
		},
	}

	m.updatePreview()
	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Export modal handling
		if m.exportModal.active {
			return m.handleExportModalKey(msg)
		}

		// Help screen
		if m.showHelp {
			if msg.String() == "?" || msg.String() == "esc" || msg.String() == "q" {
				m.showHelp = false
			}
			return m, nil
		}

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
			if idx < tabCount {
				m.currentTab = idx
				m.focus = focusInput
				m.textInput.Focus()
				m.searchInput.Blur()
				m.updatePreview()
			}
			return m, nil

		case "tab":
			m.currentTab = (m.currentTab + 1) % tabCount
			m.focus = focusInput
			m.textInput.Focus()
			m.updatePreview()
			return m, nil

		case "shift+tab":
			m.currentTab = (m.currentTab + tabCount - 1) % tabCount
			m.focus = focusInput
			m.textInput.Focus()
			m.updatePreview()
			return m, nil

		case "ctrl+e", "e":
			if m.focus != focusInput && m.preview != "" {
				m.openExportModal()
				return m, nil
			}

		case "enter":
			if m.preview != "" {
				plain := stripANSI(m.preview)
				if err := clipboard.WriteAll(plain); err == nil {
					m.statusMsg = "Copied to clipboard!"
					m.statusTime = time.Now()
				}
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

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

func (m *Model) handleUp() {
	switch m.currentTab {
	case tabBanner:
		if m.focus == focusList1 && m.selectedFont > 0 {
			m.selectedFont--
		} else if m.focus == focusList2 && m.selectedStyle > 0 {
			m.selectedStyle--
		}
	case tabKaomoji:
		if m.focus == focusList1 && m.selectedKaomoji > 0 {
			m.selectedKaomoji--
		} else if m.focus == focusList2 && m.selectedCat > 0 {
			m.selectedCat--
			m.filterKaomoji()
		}
	case tabArtDB:
		if m.focus == focusList1 && m.selectedArt > 0 {
			m.selectedArt--
		} else if m.focus == focusList2 && m.selectedArtCat > 0 {
			m.selectedArtCat--
			m.filterArt()
		}
	case tabFilters:
		if m.focus == focusList1 && m.selectedFilter > 0 {
			m.selectedFilter--
		}
	case tabEffects:
		if m.focus == focusList1 && m.selectedEffect > 0 {
			m.selectedEffect--
		}
	case tabGradient:
		if m.focus == focusList1 && m.selectedTheme > 0 {
			m.selectedTheme--
		}
	case tabQRCode:
		if m.focus == focusList1 && m.selectedQRCharset > 0 {
			m.selectedQRCharset--
		}
	case tabPatterns:
		if m.patternMode == 0 && m.selectedBorder > 0 {
			m.selectedBorder--
		} else if m.patternMode == 1 && m.selectedDivider > 0 {
			m.selectedDivider--
		}
	case tabSpeech:
		if m.focus == focusList1 && m.selectedSpeech > 0 {
			m.selectedSpeech--
		} else if m.focus == focusList2 && m.selectedSpeechArt > 0 {
			m.selectedSpeechArt--
		}
	case tabCalendar:
		if m.calendarMode > 0 {
			m.calendarMode--
		}
	}
}

func (m *Model) handleDown() {
	switch m.currentTab {
	case tabBanner:
		if m.focus == focusList1 && m.selectedFont < len(m.filteredFonts)-1 {
			m.selectedFont++
		} else if m.focus == focusList2 && m.selectedStyle < len(m.colorStyles)-1 {
			m.selectedStyle++
		}
	case tabKaomoji:
		if m.focus == focusList1 && m.selectedKaomoji < len(m.filteredKaomoji)-1 {
			m.selectedKaomoji++
		} else if m.focus == focusList2 && m.selectedCat < len(m.kaomojiCats)-1 {
			m.selectedCat++
			m.filterKaomoji()
		}
	case tabArtDB:
		if m.focus == focusList1 && m.selectedArt < len(m.filteredArt)-1 {
			m.selectedArt++
		} else if m.focus == focusList2 && m.selectedArtCat < len(m.artCats)-1 {
			m.selectedArtCat++
			m.filterArt()
		}
	case tabFilters:
		if m.focus == focusList1 && m.selectedFilter < len(m.filterList)-1 {
			m.selectedFilter++
		}
	case tabEffects:
		if m.focus == focusList1 && m.selectedEffect < len(m.effectList)-1 {
			m.selectedEffect++
		}
	case tabGradient:
		if m.focus == focusList1 && m.selectedTheme < len(m.gradientThemes)-1 {
			m.selectedTheme++
		}
	case tabQRCode:
		if m.focus == focusList1 && m.selectedQRCharset < len(m.qrCharsets)-1 {
			m.selectedQRCharset++
		}
	case tabPatterns:
		if m.patternMode == 0 && m.selectedBorder < len(m.patternBorders)-1 {
			m.selectedBorder++
		} else if m.patternMode == 1 && m.selectedDivider < len(m.patternDividers)-1 {
			m.selectedDivider++
		}
	case tabSpeech:
		if m.focus == focusList1 && m.selectedSpeech < len(m.speechStyles)-1 {
			m.selectedSpeech++
		} else if m.focus == focusList2 && m.selectedSpeechArt < len(m.speechArts)-1 {
			m.selectedSpeechArt++
		}
	case tabCalendar:
		if m.calendarMode < 2 {
			m.calendarMode++
		}
	}
}

func (m *Model) handleLeft() {
	if m.focus > focusInput {
		m.focus--
	}
	if m.currentTab == tabPatterns {
		m.patternMode = 0
	}
	if m.currentTab == tabQRCode {
		m.qrInvert = false
	}
	if m.currentTab == tabGradient && m.gradientMode > 0 {
		m.gradientMode--
	}
}

func (m *Model) handleRight() {
	m.focus++
	if m.focus > focusList2 {
		m.focus = focusList2
	}
	if m.currentTab == tabPatterns {
		m.patternMode = 1
	}
	if m.currentTab == tabQRCode {
		m.qrInvert = true
	}
	if m.currentTab == tabGradient && m.gradientMode < 2 {
		m.gradientMode++
	}
}

func (m *Model) filterCurrentList() {
	switch m.currentTab {
	case tabBanner:
		m.filterFonts()
	case tabKaomoji:
		m.filterKaomoji()
	case tabArtDB:
		m.filterArt()
	}
}

func (m *Model) filterFonts() {
	search := strings.ToLower(m.searchInput.Value())
	m.filteredFonts = m.filteredFonts[:0]
	for i, font := range m.fonts {
		if search == "" || strings.Contains(strings.ToLower(font), search) {
			m.filteredFonts = append(m.filteredFonts, i)
		}
	}
	if m.selectedFont >= len(m.filteredFonts) {
		m.selectedFont = 0
	}
}

func (m *Model) filterKaomoji() {
	search := strings.ToLower(m.searchInput.Value())
	cat := ""
	if m.selectedCat > 0 {
		cat = m.kaomojiCats[m.selectedCat]
	}

	m.filteredKaomoji = m.filteredKaomoji[:0]
	for i, k := range m.kaomojiList {
		matchSearch := search == "" || strings.Contains(strings.ToLower(k.Name), search)
		matchCat := cat == "" || k.Category == cat
		if matchSearch && matchCat {
			m.filteredKaomoji = append(m.filteredKaomoji, i)
		}
	}
	if m.selectedKaomoji >= len(m.filteredKaomoji) {
		m.selectedKaomoji = 0
	}
}

func (m *Model) filterArt() {
	search := strings.ToLower(m.searchInput.Value())
	cat := ""
	if m.selectedArtCat > 0 {
		cat = m.artCats[m.selectedArtCat]
	}

	m.filteredArt = m.filteredArt[:0]
	for i, a := range m.artList {
		matchSearch := search == "" || strings.Contains(strings.ToLower(a.Name), search)
		matchCat := cat == "" || a.Category == cat
		if matchSearch && matchCat {
			m.filteredArt = append(m.filteredArt, i)
		}
	}
	if m.selectedArt >= len(m.filteredArt) {
		m.selectedArt = 0
	}
}

func (m *Model) updatePreview() {
	text := m.textInput.Value()
	if text == "" {
		text = "MOJI"
	}

	switch m.currentTab {
	case tabBanner:
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

	case tabKaomoji:
		if len(m.filteredKaomoji) > 0 {
			k := m.kaomojiList[m.filteredKaomoji[m.selectedKaomoji]]
			m.preview = fmt.Sprintf("%s\n\n%s", selectedStyle.Render(k.Kaomoji), dimStyle.Render(k.Name))
		}

	case tabArtDB:
		if len(m.filteredArt) > 0 {
			a := m.artList[m.filteredArt[m.selectedArt]]
			m.preview = a.Art
		}

	case tabFilters:
		filter := m.filterList[m.selectedFilter]
		if f, ok := filters.Get(filter); ok {
			m.preview = f(text)
		}

	case tabEffects:
		effect := m.effectList[m.selectedEffect]
		m.preview = effects.Apply(effect, text)

	case tabGradient:
		modes := []string{"horizontal", "vertical", "diagonal"}
		mode := modes[m.gradientMode]
		theme := m.gradientThemes[m.selectedTheme]
		m.preview = gradient.Apply(text, theme, mode)

	case tabQRCode:
		charset := m.qrCharsets[m.selectedQRCharset]
		qr, err := qrcode.Generate(text, qrcode.Options{Charset: charset, Invert: m.qrInvert})
		if err == nil {
			m.preview = qr
		} else {
			m.preview = errorStyle.Render(err.Error())
		}

	case tabPatterns:
		if m.patternMode == 0 {
			border := m.patternBorders[m.selectedBorder]
			m.preview = patterns.CreateBorder(text, border, 1)
		} else {
			divider := m.patternDividers[m.selectedDivider]
			m.preview = patterns.CreateDivider(divider, 40)
		}

	case tabSpeech:
		style := m.speechStyles[m.selectedSpeech]
		bubble := speech.Wrap(text, style, 40)
		if m.selectedSpeechArt > 0 {
			artName := m.speechArts[m.selectedSpeechArt]
			if a, ok := artdb.Get(artName); ok {
				bubble = speech.Combine(bubble, a.Art)
			}
		}
		m.preview = bubble

	case tabCalendar:
		opts := calendar.Options{FirstDayMon: m.calMonday}
		switch m.calendarMode {
		case 0:
			m.preview = calendar.Current(opts)
		case 1:
			m.preview = calendar.WeekView(opts)
		case 2:
			m.preview = calendar.Year(time.Now().Year(), opts)
		}

	case tabSysinfo:
		info := sysinfo.Collect()
		art := sysinfo.GetOSLogo()
		if text != "MOJI" {
			art = gradient.Apply(art, text, "diagonal")
		}
		m.preview = info.FormatWithArt(art)
	}
}

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

func (m Model) renderTabContent() string {
	switch m.currentTab {
	case tabBanner:
		return m.renderBannerTab()
	case tabKaomoji:
		return m.renderKaomojiTab()
	case tabArtDB:
		return m.renderArtDBTab()
	case tabFilters:
		return m.renderFiltersTab()
	case tabEffects:
		return m.renderEffectsTab()
	case tabGradient:
		return m.renderGradientTab()
	case tabQRCode:
		return m.renderQRTab()
	case tabPatterns:
		return m.renderPatternsTab()
	case tabSpeech:
		return m.renderSpeechTab()
	case tabCalendar:
		return m.renderCalendarTab()
	case tabSysinfo:
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
	}, m.focus == focusList1)

	styleList := m.renderList("Style", makeRange(len(m.colorStyles)), m.selectedStyle, func(i int) string {
		return m.colorStyles[i]
	}, m.focus == focusList2)

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, fontList, "  ", styleList))

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
	}, m.focus == focusList1)

	catList := m.renderList("Category", makeRange(len(m.kaomojiCats)), m.selectedCat, func(i int) string {
		return m.kaomojiCats[i]
	}, m.focus == focusList2)

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, kaomojiList, "  ", catList))

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
	}, m.focus == focusList1)

	catList := m.renderList("Category", makeRange(len(m.artCats)), m.selectedArtCat, func(i int) string {
		return m.artCats[i]
	}, m.focus == focusList2)

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, artList, "  ", catList))

	return b.String()
}

func (m Model) renderFiltersTab() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Text: "))
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	filterList := m.renderList("Filter", makeRange(len(m.filterList)), m.selectedFilter, func(i int) string {
		return m.filterList[i]
	}, m.focus == focusList1)

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
	}, m.focus == focusList1)

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
	}, m.focus == focusList1)

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
	}, m.focus == focusList1)

	artList := m.renderList("Character", makeRange(len(m.speechArts)), m.selectedSpeechArt, func(i int) string {
		return m.speechArts[i]
	}, m.focus == focusList2)

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, styleList, "  ", artList))

	return b.String()
}

func (m Model) renderQRTab() string {
	var b strings.Builder

	b.WriteString(subtitleStyle.Render("Text/URL: "))
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	charsetList := m.renderList("Charset", makeRange(len(m.qrCharsets)), m.selectedQRCharset, func(i int) string {
		return m.qrCharsets[i]
	}, m.focus == focusList1)

	invertStr := "Normal"
	if m.qrInvert {
		invertStr = "Inverted"
	}
	options := boxStyle.Render(fmt.Sprintf("Invert: %s\n[←/→] to toggle", invertStr))

	b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, charsetList, "  ", options))

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

// Export modal methods
func (m *Model) openExportModal() {
	m.exportModal.active = true
	m.exportModal.format = exportPNG
	m.exportModal.input.SetValue(generateFilename("png"))
	m.exportModal.input.Focus()
	m.exportModal.showBrowser = false
	m.loadDirectory()
}

func (m *Model) loadDirectory() {
	files, err := os.ReadDir(m.exportModal.currentDir)
	if err != nil {
		m.exportModal.files = nil
		return
	}
	m.exportModal.files = files
	m.exportModal.selectedIdx = 0
}

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
		if i+1 == m.exportModal.format {
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

// Utility functions
func makeRange(n int) []int {
	r := make([]int, n)
	for i := range r {
		r[i] = i
	}
	return r
}

func generateFilename(ext string) string {
	return fmt.Sprintf("moji_%s.%s", time.Now().Format("20060102_150405"), ext)
}

func changeExtension(filename, newExt string) string {
	ext := filepath.Ext(filename)
	if ext != "" {
		return filename[:len(filename)-len(ext)] + "." + newExt
	}
	return filename + "." + newExt
}

func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

// Run starts the TUI
func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
