package tui

// FocusState represents which UI element has focus in the model.
type FocusState int

const (
	focusInput FocusState = iota
	focusList1
	focusList2
	focusList3
)

// Tab represents a specific TUI tab
type Tab int

const (
	tabBanner Tab = iota
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

// ExportFormat represents supported export formats
type ExportFormat int

const (
	exportNone ExportFormat = iota
	exportPNG
	exportSVG
	exportHTML
	exportTXT
)

var tabNames = []string{
	"Banner", "Kaomoji", "ArtDB", "Filters", "Effects", "Gradient", "QR", "Patterns", "Speech", "Calendar", "Sysinfo",
}

// nextTab moves to the next tab
func (m *Model) nextTab() {
	m.currentTab = (m.currentTab + 1) % int(tabCount)
	m.focus = focusInput
	m.textInput.Focus()
}

// prevTab moves to the previous tab
func (m *Model) prevTab() {
	m.currentTab = (m.currentTab + int(tabCount) - 1) % int(tabCount)
	m.focus = focusInput
	m.textInput.Focus()
}

// goToTab switches to a specific tab by index
func (m *Model) goToTab(idx int) {
	if idx >= 0 && idx < int(tabCount) {
		m.currentTab = idx
		m.focus = focusInput
		m.textInput.Focus()
		m.searchInput.Blur()
	}
}
