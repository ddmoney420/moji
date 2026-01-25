package animate

import (
	"fmt"
	"strings"
	"time"
)

// Frames represents animation frames
type Frames []string

// Animation presets
var Presets = map[string]Frames{
	"spinner": {"|", "/", "-", "\\"},
	"dots":    {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	"bounce":  {"⠁", "⠂", "⠄", "⠂"},
	"grow":    {".", "..", "...", "....", ".....", "......"},
	"arrows":  {"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
	"clock":   {"🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚", "🕛"},
	"moon":    {"🌑", "🌒", "🌓", "🌔", "🌕", "🌖", "🌗", "🌘"},
	"earth":   {"🌍", "🌎", "🌏"},
	"hearts":  {"💗", "💓", "💗", "💖"},
	"bar":     {"[    ]", "[=   ]", "[==  ]", "[=== ]", "[====]", "[ ===]", "[  ==]", "[   =]"},
	"blocks":  {"▏", "▎", "▍", "▌", "▋", "▊", "▉", "█", "▉", "▊", "▋", "▌", "▍", "▎"},
	"wave":    {"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█", "▇", "▆", "▅", "▄", "▃", "▂"},
	"box":     {"◰", "◳", "◲", "◱"},
	"toggle":  {"■", "□"},
	"star":    {"✶", "✷", "✹", "✺"},
	"bounce2": {"⠁", "⠂", "⠄", "⡀", "⢀", "⠠", "⠐", "⠈"},
}

// Play plays an animation n times with given delay
func Play(frames Frames, n int, delayMs int) {
	if len(frames) == 0 {
		return
	}
	delay := time.Duration(delayMs) * time.Millisecond
	for i := 0; i < n; i++ {
		for _, frame := range frames {
			fmt.Printf("\r%s ", frame)
			time.Sleep(delay)
		}
	}
	fmt.Print("\r\033[K") // Clear line
}

// PlayWithText plays animation with text message
func PlayWithText(frames Frames, text string, n int, delayMs int) {
	if len(frames) == 0 {
		return
	}
	delay := time.Duration(delayMs) * time.Millisecond
	for i := 0; i < n; i++ {
		for _, frame := range frames {
			fmt.Printf("\r%s %s ", frame, text)
			time.Sleep(delay)
		}
	}
	fmt.Print("\r\033[K")
}

// Typewriter prints text one character at a time
func Typewriter(text string, delayMs int) {
	delay := time.Duration(delayMs) * time.Millisecond
	for _, r := range text {
		fmt.Print(string(r))
		time.Sleep(delay)
	}
	fmt.Println()
}

// ScrollText scrolls text horizontally
func ScrollText(text string, width int, n int, delayMs int) {
	if width <= 0 {
		width = 40
	}
	delay := time.Duration(delayMs) * time.Millisecond
	padded := strings.Repeat(" ", width) + text + strings.Repeat(" ", width)
	runes := []rune(padded)

	for cycle := 0; cycle < n; cycle++ {
		for i := 0; i <= len(runes)-width; i++ {
			window := string(runes[i : i+width])
			fmt.Printf("\r%s", window)
			time.Sleep(delay)
		}
	}
	fmt.Println()
}

// FadeIn simulates a fade-in effect using block characters
func FadeIn(text string, delayMs int) {
	lines := strings.Split(text, "\n")
	stages := []string{"░", "▒", "▓", "█"}

	for _, stage := range stages[:len(stages)-1] {
		fmt.Print("\033[H\033[2J") // Clear screen
		for _, line := range lines {
			var masked strings.Builder
			for _, r := range line {
				if r != ' ' && r != '\t' && r != '\n' {
					masked.WriteString(stage)
				} else {
					masked.WriteRune(r)
				}
			}
			fmt.Println(masked.String())
		}
		time.Sleep(time.Duration(delayMs) * time.Millisecond)
	}
	// Final stage - show actual text
	fmt.Print("\033[H\033[2J")
	fmt.Println(text)
}

// Blink blinks text n times
func Blink(text string, n int, delayMs int) {
	delay := time.Duration(delayMs) * time.Millisecond
	blank := strings.Repeat(" ", len([]rune(text)))
	for i := 0; i < n; i++ {
		fmt.Printf("\r%s", text)
		time.Sleep(delay)
		fmt.Printf("\r%s", blank)
		time.Sleep(delay)
	}
	fmt.Printf("\r%s\n", text)
}

// Matrix rain effect (simplified single column)
func MatrixRain(width, height int, durationMs int) {
	chars := []rune("ｱｲｳｴｵｶｷｸｹｺｻｼｽｾｿﾀﾁﾂﾃﾄﾅﾆﾇﾈﾉﾊﾋﾌﾍﾎﾏﾐﾑﾒﾓﾔﾕﾖﾗﾘﾙﾚﾛﾜﾝ0123456789")
	columns := make([]int, width)

	endTime := time.Now().Add(time.Duration(durationMs) * time.Millisecond)
	for time.Now().Before(endTime) {
		var line strings.Builder
		for i := 0; i < width; i++ {
			if columns[i] > 0 {
				c := chars[columns[i]%len(chars)]
				line.WriteString(fmt.Sprintf("\033[32m%c\033[0m", c))
				columns[i]++
				if columns[i] > height {
					columns[i] = 0
				}
			} else {
				line.WriteString(" ")
				// Random chance to start a new drop
				if time.Now().UnixNano()%100 < 5 {
					columns[i] = 1
				}
			}
		}
		fmt.Printf("\r%s", line.String())
		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println()
}

// ListPresets returns available animation presets
func ListPresets() []string {
	var names []string
	for k := range Presets {
		names = append(names, k)
	}
	return names
}

// GetPreset returns a preset animation
func GetPreset(name string) (Frames, bool) {
	f, ok := Presets[name]
	return f, ok
}
