package tui

import (
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ddmoney420/moji/internal/artdb"
	"github.com/ddmoney420/moji/internal/banner"
	"github.com/ddmoney420/moji/internal/effects"
	"github.com/ddmoney420/moji/internal/filters"
	"github.com/ddmoney420/moji/internal/gradient"
	"github.com/ddmoney420/moji/internal/kaomoji"
	"github.com/ddmoney420/moji/internal/patterns"
	"github.com/ddmoney420/moji/internal/qrcode"
	"github.com/ddmoney420/moji/internal/styles"
)

// Model represents the TUI state
type Model struct {
	// Common
	currentTab  int
	focus       FocusState
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

// NewModel creates a new TUI model with initialized state
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
		focus:           focusInput,
		exportModal: ExportModal{
			input:      ei,
			currentDir: cwd,
		},
	}

	m.updatePreview()
	return m
}

// Init returns initial commands to run
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and updates model state
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

		// Main key handling
		return m.handleKeyboard(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, cmd
}

// filterCurrentList applies filters based on the current tab
func (m *Model) filterCurrentList() {
	switch m.currentTab {
	case int(tabBanner):
		m.filterFonts()
	case int(tabKaomoji):
		m.filterKaomoji()
	case int(tabArtDB):
		m.filterArt()
	}
}

// filterFonts filters the font list based on search input
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

// filterKaomoji filters the kaomoji list based on search and category
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

// filterArt filters the art list based on search and category
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

// Run starts the TUI program
func Run() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
