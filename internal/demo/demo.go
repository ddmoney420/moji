package demo

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ddmoney420/moji/internal/artdb"
	"github.com/ddmoney420/moji/internal/banner"
	"github.com/ddmoney420/moji/internal/calendar"
	"github.com/ddmoney420/moji/internal/effects"
	"github.com/ddmoney420/moji/internal/filters"
	"github.com/ddmoney420/moji/internal/gradient"
	"github.com/ddmoney420/moji/internal/kaomoji"
	"github.com/ddmoney420/moji/internal/patterns"
	"github.com/ddmoney420/moji/internal/progress"
	"github.com/ddmoney420/moji/internal/qrcode"
	"github.com/ddmoney420/moji/internal/styles"
	"github.com/ddmoney420/moji/internal/sysinfo"
	"github.com/ddmoney420/moji/internal/tree"
	"github.com/ddmoney420/moji/internal/ux"
)

// Options for the demo
type Options struct {
	Speed    string // slow, normal, fast
	Category string // specific category to demo
	NoWait   bool   // don't wait between sections
}

var pauseDuration = 2 * time.Second

// Run executes the demo
func Run(opts Options) {
	switch opts.Speed {
	case "slow":
		pauseDuration = 4 * time.Second
	case "fast":
		pauseDuration = 500 * time.Millisecond
	case "instant":
		pauseDuration = 0
		opts.NoWait = true
	default:
		pauseDuration = 2 * time.Second
	}

	clearScreen()

	// Title
	showTitle()
	pause(opts)

	// Run demos based on category or all
	if opts.Category != "" {
		runCategory(opts.Category, opts)
	} else {
		runAllDemos(opts)
	}

	// Finale
	showFinale(opts)
}

func runAllDemos(opts Options) {
	demoKaomoji(opts)
	demoBanners(opts)
	demoFilters(opts)
	demoGradients(opts)
	demoArtDB(opts)
	demoEffects(opts)
	demoPatterns(opts)
	demoQRCode(opts)
	demoCalendar(opts)
	demoTree(opts)
	demoSysinfo(opts)
	demoProgress(opts)
}

func runCategory(cat string, opts Options) {
	switch strings.ToLower(cat) {
	case "kaomoji":
		demoKaomoji(opts)
	case "banner", "banners":
		demoBanners(opts)
	case "filter", "filters":
		demoFilters(opts)
	case "gradient", "gradients":
		demoGradients(opts)
	case "effect", "effects":
		demoEffects(opts)
	case "pattern", "patterns":
		demoPatterns(opts)
	case "qr", "qrcode":
		demoQRCode(opts)
	case "calendar", "cal":
		demoCalendar(opts)
	case "tree":
		demoTree(opts)
	case "sysinfo":
		demoSysinfo(opts)
	case "progress":
		demoProgress(opts)
	case "art", "artdb":
		demoArtDB(opts)
	default:
		fmt.Printf("Unknown category: %s\n", cat)
		fmt.Println("Available: kaomoji, banners, filters, gradients, effects, patterns, qr, calendar, tree, sysinfo, progress, artdb")
	}
}

func clearScreen() {
	fmt.Print("\033[2J\033[H")
}

func pause(opts Options) {
	if !opts.NoWait {
		time.Sleep(pauseDuration)
	}
}

func shortPause() {
	time.Sleep(300 * time.Millisecond)
}

func section(title string) {
	clearScreen()
	divider := strings.Repeat("═", 60)
	fmt.Printf("\033[1;36m%s\033[0m\n", divider)
	fmt.Printf("\033[1;37m  %s\033[0m\n", title)
	fmt.Printf("\033[1;36m%s\033[0m\n\n", divider)
}

func subsection(title string) {
	fmt.Printf("\n\033[1;33m▸ %s\033[0m\n\n", title)
}

func command(cmd string) {
	fmt.Printf("\033[2m$ %s\033[0m\n", cmd)
	shortPause()
}

func showTitle() {
	section("MOJI - CLI Demo")

	// Animated title
	title, _ := banner.Generate("moji", "slant")
	title = gradient.Apply(title, "rainbow", "horizontal")
	fmt.Println(title)

	fmt.Println("\033[2mA CLI tool for kaomoji, ASCII banners, emoji, and more!\033[0m")
	fmt.Println()
	fmt.Println("Press Ctrl+C to exit at any time.")
}

func demoKaomoji(opts Options) {
	section("KAOMOJI")

	subsection("Get kaomoji by name")
	kaomojis := []string{"shrug", "happy", "sad", "love", "angry", "confused", "cool", "tableflip"}

	for _, name := range kaomojis {
		k, ok := kaomoji.Get(name)
		if ok {
			command(fmt.Sprintf("moji %s", name))
			fmt.Printf("  %s\n\n", k)
			shortPause()
		}
	}

	pause(opts)

	subsection("Random kaomoji")
	command("moji random")
	name, k := kaomoji.Random()
	fmt.Printf("  %s (%s)\n", k, name)

	pause(opts)
}

func demoBanners(opts Options) {
	section("ASCII BANNERS")

	subsection("Different fonts")
	fonts := []string{"standard", "slant", "shadow", "small", "big"}

	for _, font := range fonts {
		command(fmt.Sprintf("moji banner \"Hi\" --font %s", font))
		art, _ := banner.Generate("Hi", font)
		fmt.Println(art)
		shortPause()
	}

	pause(opts)

	subsection("With color styles")
	command("moji banner \"Cool\" --style rainbow")
	art, _ := banner.Generate("Cool", "standard")
	styled := styles.Apply(art, "rainbow")
	fmt.Println(styled)

	pause(opts)

	subsection("With borders")
	command("moji banner \"Box\" --border double")
	art, _ = banner.Generate("Box", "small")
	bordered := styles.ApplyBorder(art, "double")
	fmt.Println(bordered)

	pause(opts)
}

func demoFilters(opts Options) {
	section("TEXT FILTERS")

	text := "Hello World"
	filterList := []struct {
		name string
		fn   func(string) string
	}{
		{"rainbow", filters.Rainbow},
		{"metal", filters.Metal},
		{"fire", filters.Fire},
		{"ice", filters.Ice},
		{"neon", filters.Neon},
		{"matrix", filters.Matrix},
		{"glitch", filters.Glitch},
	}

	for _, f := range filterList {
		command(fmt.Sprintf("moji filter %s \"%s\"", f.name, text))
		result := f.fn(text)
		fmt.Println(result)
		fmt.Println()
		shortPause()
	}

	pause(opts)

	subsection("Filter chaining")
	command("moji filter rainbow,border \"Chained\"")
	result := filters.Chain("Chained", []string{"rainbow", "border"})
	fmt.Println(result)

	pause(opts)

	subsection("Filters on banners")
	command("moji banner \"WOW\" --font small | moji filter neon")
	art, _ := banner.Generate("WOW", "small")
	neon := filters.Neon(art)
	fmt.Println(neon)

	pause(opts)
}

func demoGradients(opts Options) {
	section("COLOR GRADIENTS")

	text := "Gradient Text Demo Line"
	themes := []string{"rainbow", "neon", "fire", "ocean", "sunset", "forest", "galaxy"}

	for _, theme := range themes {
		command(fmt.Sprintf("moji gradient \"%s\" --theme %s", text, theme))
		result := gradient.Apply(text, theme, "horizontal")
		fmt.Println(result)
		fmt.Println()
		shortPause()
	}

	pause(opts)

	subsection("Gradient on ASCII art")
	command("moji banner \"GO\" | moji gradient --theme galaxy")
	art, _ := banner.Generate("GO", "slant")
	result := gradient.Apply(art, "galaxy", "diagonal")
	fmt.Println(result)

	pause(opts)
}

func demoEffects(opts Options) {
	section("TEXT EFFECTS")

	text := "Hello World"
	effectsList := []string{"flip", "reverse", "mirror", "bubble", "square", "script"}

	for _, e := range effectsList {
		command(fmt.Sprintf("moji effect %s \"%s\"", e, text))
		result := effects.Apply(e, text)
		fmt.Printf("  %s\n\n", result)
		shortPause()
	}

	pause(opts)

	subsection("Zalgo effect")
	command("moji effect zalgo \"Spooky\"")
	result := effects.Apply("zalgo", "Spooky")
	fmt.Printf("  %s\n", result)

	pause(opts)
}

func demoPatterns(opts Options) {
	section("PATTERNS & BORDERS")

	subsection("Dividers")
	dividers := []string{"double", "stars", "wave", "dots", "heavy"}
	for _, d := range dividers {
		command(fmt.Sprintf("moji pattern --divider %s --width 40", d))
		result := patterns.CreateDivider(d, 40)
		fmt.Println(result)
		shortPause()
	}

	pause(opts)

	subsection("Bordered boxes")
	borders := []string{"single", "double", "round", "heavy"}
	for _, b := range borders {
		command(fmt.Sprintf("moji pattern --border %s --text \"Hello\"", b))
		result := patterns.CreateBorder("Hello", b, 1)
		fmt.Println(result)
		shortPause()
	}

	pause(opts)
}

func demoQRCode(opts Options) {
	section("QR CODES")

	command("moji qr \"https://github.com\"")
	qr, _ := qrcode.GenerateCompact("https://github.com", false)
	fmt.Println(qr)

	pause(opts)

	subsection("Different styles")
	command("moji qr \"Hello\" --charset dots")
	qr, _ = qrcode.Generate("Hello", qrcode.Options{Charset: "dots"})
	fmt.Println(qr)

	pause(opts)
}

func demoCalendar(opts Options) {
	section("CALENDAR")

	subsection("Current month")
	command("moji cal")
	cal := calendar.Current(calendar.Options{})
	fmt.Println(cal)

	pause(opts)

	subsection("Week view")
	command("moji cal --week")
	week := calendar.WeekView(calendar.Options{})
	fmt.Println(week)

	pause(opts)
}

func demoTree(opts Options) {
	section("DIRECTORY TREE")

	command("moji tree . --depth 2")
	opts2 := tree.DefaultOptions()
	opts2.MaxDepth = 2
	entry, err := tree.Generate(".", opts2)
	if err == nil {
		result := tree.Format(entry, opts2)
		// Limit output
		lines := strings.Split(result, "\n")
		if len(lines) > 15 {
			lines = lines[:15]
			lines = append(lines, "  ...")
		}
		fmt.Println(strings.Join(lines, "\n"))
	}

	pause(opts)
}

func demoSysinfo(opts Options) {
	section("SYSTEM INFO")

	command("moji sysinfo")
	info := sysinfo.Collect()
	art := sysinfo.GetOSLogo()
	result := info.FormatWithArt(art)
	fmt.Println(result)

	pause(opts)
}

func demoProgress(opts Options) {
	section("PROGRESS BARS & SPINNERS")

	subsection("Progress bar styles")
	progressStyles := []string{"default", "blocks", "smooth", "dots", "arrows", "shades", "braille", "minimal"}

	for _, style := range progressStyles {
		bar := progress.Bar(75, 100, 35, style)
		fmt.Printf("  %-10s %s\n", style, bar)
		shortPause()
	}

	pause(opts)

	subsection("Animated progress bar")
	fmt.Println()
	styles := []string{"smooth", "blocks", "dots"}
	for _, style := range styles {
		fmt.Printf("  \033[1m%s:\033[0m\n  ", style)
		for i := 0; i <= 100; i += 2 {
			bar := progress.Bar(i, 100, 45, style)
			fmt.Printf("\r  %s", bar)
			time.Sleep(30 * time.Millisecond)
		}
		fmt.Println(" Done!")
		time.Sleep(300 * time.Millisecond)
	}

	pause(opts)

	subsection("Spinners (animated)")
	fmt.Println()
	spinnerStyles := []string{"dots", "line", "bounce", "grow", "circle", "arrow", "clock", "moon"}

	// Show all spinners animating together
	fmt.Println("  \033[2mWatch each spinner animate:\033[0m")
	for _, style := range spinnerStyles {
		fmt.Printf("  %-10s ", style)
		for frame := 0; frame < 20; frame++ {
			s := progress.Spinner(frame, style)
			fmt.Printf("\r  %-10s %s  ", style, s)
			time.Sleep(80 * time.Millisecond)
		}
		fmt.Println()
	}

	pause(opts)

	subsection("Spinner with message")
	messages := []string{"Loading", "Processing", "Almost done", "Finishing up"}
	fmt.Println()
	for _, msg := range messages {
		for frame := 0; frame < 15; frame++ {
			s := progress.Spinner(frame, "dots")
			fmt.Printf("\r  %s %s...   ", s, msg)
			time.Sleep(80 * time.Millisecond)
		}
	}
	fmt.Printf("\r  \033[32m✓\033[0m Complete!        \n")

	pause(opts)

	subsection("Multi-progress bars")
	fmt.Println()
	tasks := []struct {
		name  string
		speed int
	}{
		{"Downloading", 3},
		{"Extracting", 5},
		{"Installing", 2},
	}

	progress1, progress2, progress3 := 0, 0, 0
	for progress1 < 100 || progress2 < 100 || progress3 < 100 {
		// Move cursor up 3 lines and redraw
		if progress1 > 0 {
			fmt.Print("\033[3A")
		}

		progress1 = min(progress1+tasks[0].speed, 100)
		progress2 = min(progress2+tasks[1].speed, 100)
		progress3 = min(progress3+tasks[2].speed, 100)

		bar1 := progress.Bar(progress1, 100, 30, "smooth")
		bar2 := progress.Bar(progress2, 100, 30, "smooth")
		bar3 := progress.Bar(progress3, 100, 30, "smooth")

		status1, status2, status3 := "⏳", "⏳", "⏳"
		if progress1 >= 100 {
			status1 = "\033[32m✓\033[0m"
		}
		if progress2 >= 100 {
			status2 = "\033[32m✓\033[0m"
		}
		if progress3 >= 100 {
			status3 = "\033[32m✓\033[0m"
		}

		fmt.Printf("  %s %-12s %s\n", status1, tasks[0].name, bar1)
		fmt.Printf("  %s %-12s %s\n", status2, tasks[1].name, bar2)
		fmt.Printf("  %s %-12s %s\n", status3, tasks[2].name, bar3)

		time.Sleep(50 * time.Millisecond)
	}
	fmt.Println("\n  \033[1;32mAll tasks complete!\033[0m")

	pause(opts)
}

func demoArtDB(opts Options) {
	section("ASCII ART DATABASE")

	command("moji artdb --categories")
	fmt.Println("Categories:")
	for _, cat := range artdb.ListCategories() {
		count := len(artdb.ByCategory(cat))
		fmt.Printf("  %-12s (%d)\n", cat, count)
	}

	pause(opts)

	subsection("Sample art")
	arts := artdb.List()
	if len(arts) > 3 {
		arts = arts[:3]
	}
	for _, a := range arts {
		command(fmt.Sprintf("moji artdb %s", a.Name))
		fmt.Println(a.Art)
		shortPause()
	}

	pause(opts)
}

func showFinale(opts Options) {
	section("DEMO COMPLETE!")

	// Animated rainbow thank you
	thankYou, _ := banner.Generate("Thanks!", "slant")

	if opts.NoWait {
		colored := gradient.Apply(thankYou, "rainbow", "horizontal")
		fmt.Println(colored)
	} else {
		// Animate the rainbow
		lines := strings.Split(thankYou, "\n")
		for cycle := 0; cycle < 30; cycle++ {
			fmt.Print("\033[H\033[6B") // Move down 6 lines
			for lineIdx, line := range lines {
				for i, r := range line {
					if r == ' ' || r == '\t' {
						fmt.Print(string(r))
						continue
					}
					phase := (float64(i+lineIdx) + float64(cycle)*0.5) * 0.15
					red := uint8(math.Sin(phase)*127 + 128)
					green := uint8(math.Sin(phase+2*math.Pi/3)*127 + 128)
					blue := uint8(math.Sin(phase+4*math.Pi/3)*127 + 128)
					fmt.Printf("\033[38;2;%d;%d;%dm%c\033[0m", red, green, blue, r)
				}
				fmt.Println()
			}
			time.Sleep(50 * time.Millisecond)
		}
	}

	fmt.Println()
	fmt.Println("\033[1mExplore more:\033[0m")
	fmt.Println("  moji --help           Show all commands")
	fmt.Println("  moji doctor           Check your setup")
	fmt.Println("  moji interactive      Launch the TUI studio")
	fmt.Println()
	fmt.Printf("  \033[2m%s\033[0m\n", ux.VersionString())
	fmt.Println()
}

// ListCategories returns available demo categories
func ListCategories() []string {
	return []string{
		"kaomoji",
		"banners",
		"filters",
		"gradients",
		"effects",
		"patterns",
		"qr",
		"calendar",
		"tree",
		"sysinfo",
		"progress",
		"artdb",
	}
}
