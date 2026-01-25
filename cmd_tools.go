package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/ddmoney420/moji/internal/animate"
	"github.com/ddmoney420/moji/internal/artdb"
	"github.com/ddmoney420/moji/internal/calendar"
	"github.com/ddmoney420/moji/internal/config"
	"github.com/ddmoney420/moji/internal/demo"
	"github.com/ddmoney420/moji/internal/gradient"
	"github.com/ddmoney420/moji/internal/patterns"
	"github.com/ddmoney420/moji/internal/qrcode"
	"github.com/ddmoney420/moji/internal/sysinfo"
	"github.com/ddmoney420/moji/internal/terminal"
	"github.com/ddmoney420/moji/internal/tree"
	"github.com/ddmoney420/moji/internal/ux"
	"github.com/spf13/cobra"
)

func newQRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "qr [text]",
		Short: "Generate ASCII QR code from text or URL",
		Long: `Generate an ASCII QR code from text or URL.

Examples:
  moji qr "Hello World"
  moji qr "https://example.com" --charset blocks
  moji qr "My text" --compact --invert`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			charset, _ := cmd.Flags().GetString("charset")
			invert, _ := cmd.Flags().GetBool("invert")
			compact, _ := cmd.Flags().GetBool("compact")
			handleQR(args[0], charset, invert, compact)
		},
	}
	cmd.Flags().String("charset", "blocks", "Character set: blocks, shaded, dots, ascii, braille, compact, inverse, minimal, half")
	cmd.Flags().Bool("invert", false, "Invert colors (light on dark)")
	cmd.Flags().Bool("compact", false, "Use compact half-block rendering")
	return cmd
}

func newListQRCharsetsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-qr-charsets",
		Short: "List available QR code character sets",
		Run: func(cmd *cobra.Command, args []string) {
			handleListQRCharsets()
		},
	}
}

func newGradientCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gradient [text]",
		Short: "Apply color gradient to text",
		Long: `Apply a color gradient to text using ANSI true color.

Examples:
  moji gradient "Hello World" --theme rainbow
  moji gradient "$(moji banner Hi)" --theme neon
  cat file.txt | moji gradient --theme fire`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			theme, _ := cmd.Flags().GetString("theme")
			mode, _ := cmd.Flags().GetString("mode")
			perLine, _ := cmd.Flags().GetBool("per-line")

			var text string
			if len(args) > 0 {
				text = args[0]
			} else {
				buf := make([]byte, 1024*1024)
				n, _ := os.Stdin.Read(buf)
				text = string(buf[:n])
			}
			handleGradient(text, theme, mode, perLine)
		},
	}
	cmd.Flags().String("theme", "rainbow", "Color theme")
	cmd.Flags().String("mode", "horizontal", "Gradient mode: horizontal, vertical, diagonal")
	cmd.Flags().Bool("per-line", false, "Reset gradient per line")
	return cmd
}

func newListThemesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-themes",
		Short: "List available color gradient themes",
		Run: func(cmd *cobra.Command, args []string) {
			handleListThemes()
		},
	}
}

func newPatternCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pattern",
		Short: "Generate patterns, borders, and dividers",
		Long: `Generate decorative patterns, borders, and dividers.

Examples:
  moji pattern --divider stars --width 40
  moji pattern --border double --text "Hello"
  moji pattern --preset checker --width 20 --height 5`,
		Run: func(cmd *cobra.Command, args []string) {
			divider, _ := cmd.Flags().GetString("divider")
			border, _ := cmd.Flags().GetString("border")
			preset, _ := cmd.Flags().GetString("preset")
			text, _ := cmd.Flags().GetString("text")
			width, _ := cmd.Flags().GetInt("width")
			height, _ := cmd.Flags().GetInt("height")
			padding, _ := cmd.Flags().GetInt("padding")
			handlePattern(divider, border, preset, text, width, height, padding)
		},
	}
	cmd.Flags().String("divider", "", "Create a divider line (style name)")
	cmd.Flags().String("border", "", "Create a bordered box (style name)")
	cmd.Flags().String("preset", "", "Create a pattern block (preset name)")
	cmd.Flags().String("text", "", "Text to put inside border")
	cmd.Flags().Int("width", 40, "Width of pattern/divider")
	cmd.Flags().Int("height", 3, "Height of pattern block")
	cmd.Flags().Int("padding", 1, "Padding inside border")
	return cmd
}

func newListPatternsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-patterns",
		Short: "List available patterns, borders, and dividers",
		Run: func(cmd *cobra.Command, args []string) {
			handleListPatterns()
		},
	}
}

func newAnimateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "animate [preset]",
		Short: "Display animated ASCII effects",
		Long: `Display animated ASCII effects and spinners.

Examples:
  moji animate spinner --loops 5
  moji animate dots --text "Loading..."
  moji animate --typewriter "Hello, World!"
  moji animate --scroll "This is scrolling text"`,
		Run: func(cmd *cobra.Command, args []string) {
			listFlag, _ := cmd.Flags().GetBool("list")
			loops, _ := cmd.Flags().GetInt("loops")
			delayMs, _ := cmd.Flags().GetInt("delay")
			text, _ := cmd.Flags().GetString("text")
			typewriter, _ := cmd.Flags().GetString("typewriter")
			scroll, _ := cmd.Flags().GetString("scroll")
			width, _ := cmd.Flags().GetInt("width")
			blink, _ := cmd.Flags().GetString("blink")

			if listFlag {
				handleAnimateList()
			} else if typewriter != "" {
				handleTypewriter(typewriter, delayMs)
			} else if scroll != "" {
				handleScroll(scroll, width, loops, delayMs)
			} else if blink != "" {
				handleBlink(blink, loops, delayMs)
			} else if len(args) > 0 {
				handleAnimate(args[0], text, loops, delayMs)
			} else {
				handleAnimateList()
			}
		},
	}
	cmd.Flags().Bool("list", false, "List available animations")
	cmd.Flags().Int("loops", 3, "Number of animation loops")
	cmd.Flags().Int("delay", 100, "Delay between frames (ms)")
	cmd.Flags().String("text", "", "Text to display with animation")
	cmd.Flags().String("typewriter", "", "Text to type out character by character")
	cmd.Flags().String("scroll", "", "Text to scroll horizontally")
	cmd.Flags().String("blink", "", "Text to blink")
	cmd.Flags().Int("width", 40, "Width for scroll animation")
	return cmd
}

func newSysinfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sysinfo",
		Aliases: []string{"info", "neofetch"},
		Short:   "Display system info with ASCII art",
		Run: func(cmd *cobra.Command, args []string) {
			artName, _ := cmd.Flags().GetString("art")
			gradientTheme, _ := cmd.Flags().GetString("gradient")
			handleSysinfo(artName, gradientTheme)
		},
	}
	cmd.Flags().String("art", "", "Custom ASCII art to display (artdb name)")
	cmd.Flags().String("gradient", "", "Apply gradient to ASCII art")
	return cmd
}

func newTreeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tree [path]",
		Short: "Display directory tree as ASCII art",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}
			depth, _ := cmd.Flags().GetInt("depth")
			showHidden, _ := cmd.Flags().GetBool("all")
			dirsOnly, _ := cmd.Flags().GetBool("dirs")
			handleTree(path, depth, showHidden, dirsOnly)
		},
	}
	cmd.Flags().IntP("depth", "d", 3, "Maximum depth")
	cmd.Flags().BoolP("all", "a", false, "Show hidden files")
	cmd.Flags().Bool("dirs", false, "Show directories only")
	return cmd
}

func newCalCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "cal",
		Aliases: []string{"calendar"},
		Short:   "Display ASCII calendar",
		Run: func(cmd *cobra.Command, args []string) {
			year, _ := cmd.Flags().GetBool("year")
			week, _ := cmd.Flags().GetBool("week")
			monday, _ := cmd.Flags().GetBool("monday")
			handleCalendar(year, week, monday)
		},
	}
	cmd.Flags().BoolP("year", "y", false, "Show full year")
	cmd.Flags().BoolP("week", "w", false, "Show week view")
	cmd.Flags().BoolP("monday", "m", false, "Start week on Monday")
	return cmd
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage moji configuration",
		Run: func(cmd *cobra.Command, args []string) {
			initFlag, _ := cmd.Flags().GetBool("init")
			showFlag, _ := cmd.Flags().GetBool("show")
			pathFlag, _ := cmd.Flags().GetBool("path")

			if initFlag {
				handleConfigInit()
			} else if showFlag {
				handleConfigShow()
			} else if pathFlag {
				fmt.Println(config.ConfigPath())
			} else {
				handleConfigShow()
			}
		},
	}
	cmd.Flags().Bool("init", false, "Create default config file")
	cmd.Flags().Bool("show", false, "Show current config")
	cmd.Flags().Bool("path", false, "Show config file path")
	return cmd
}

func newTermCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "term",
		Short: "Show terminal capabilities",
		Run: func(cmd *cobra.Command, args []string) {
			handleTermInfo()
		},
	}
}

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check your environment and moji setup",
		Long: `Run diagnostic checks to verify your environment is set up correctly for moji.

Checks:
  - Terminal type and size
  - Color support level
  - Unicode/UTF-8 locale
  - Clipboard availability
  - Shell type

This is helpful for troubleshooting display or functionality issues.`,
		Run: func(cmd *cobra.Command, args []string) {
			handleDoctor()
		},
	}
}

func newDemoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "demo",
		Short: "Run an interactive demo showcasing all features",
		Long: `Run a comprehensive demo that showcases all moji features.

The demo cycles through kaomoji, banners, filters, gradients, ASCII art,
effects, patterns, QR codes, calendar, and more.

Examples:
  moji demo                     # Run full demo
  moji demo --speed fast        # Run faster
  moji demo --category banners  # Demo only banners
  moji demo --list              # List demo categories`,
		Run: func(cmd *cobra.Command, args []string) {
			listFlag, _ := cmd.Flags().GetBool("list")
			speed, _ := cmd.Flags().GetString("speed")
			category, _ := cmd.Flags().GetString("category")

			if listFlag {
				fmt.Println("Available demo categories:")
				for _, cat := range demo.ListCategories() {
					fmt.Printf("  %s\n", cat)
				}
				return
			}

			opts := demo.Options{
				Speed:    speed,
				Category: category,
			}
			demo.Run(opts)
		},
	}
	cmd.Flags().String("speed", "normal", "Demo speed: slow, normal, fast, instant")
	cmd.Flags().String("category", "", "Run demo for specific category only")
	cmd.Flags().Bool("list", false, "List available demo categories")
	return cmd
}

func newCompletionsCmd(rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "completions [bash|zsh|fish|powershell]",
		Short: "Generate shell completions",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			switch args[0] {
			case "bash":
				rootCmd.GenBashCompletion(os.Stdout)
			case "zsh":
				rootCmd.GenZshCompletion(os.Stdout)
			case "fish":
				rootCmd.GenFishCompletion(os.Stdout, true)
			case "powershell":
				rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
			default:
				fmt.Fprintf(os.Stderr, "Unknown shell: %s\n", args[0])
			}
		},
	}
}

func handleQR(text, charset string, invert, compact bool) {
	var result string
	var err error

	if compact {
		result, err = qrcode.GenerateCompact(text, invert)
	} else {
		opts := qrcode.Options{
			Charset: charset,
			Invert:  invert,
			ANSI:    !noColorFlag,
		}
		result, err = qrcode.Generate(text, opts)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating QR code: %v\n", err)
		return
	}

	if jsonFlag {
		data := map[string]string{"text": text, "qr": result}
		json.NewEncoder(os.Stdout).Encode(data)
		return
	}

	if copyFlag {
		if err := clipboard.WriteAll(result); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Println("Copied QR code to clipboard!")
	}
	fmt.Print(result)
}

func handleListQRCharsets() {
	fmt.Println("Available QR code character sets:")
	for _, name := range qrcode.ListCharsets() {
		cs := qrcode.CharSets[name]
		fmt.Printf("  %-10s  [%s] / [%s]\n", name, cs.Full, cs.Empty)
	}
}

func handleGradient(text, theme, mode string, perLine bool) {
	var result string
	if perLine {
		result = gradient.ApplyPerLine(text, theme)
	} else {
		result = gradient.Apply(text, theme, mode)
	}

	if copyFlag {
		plain := stripANSI(result)
		if err := clipboard.WriteAll(plain); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Println("Copied to clipboard!")
	}
	fmt.Println(result)
}

func handleListThemes() {
	fmt.Println("Available color gradient themes:")
	for _, t := range gradient.ListThemes() {
		fmt.Printf("  %-12s - %s\n", t.Name, t.Desc)
	}
}

func handlePattern(divider, border, preset, text string, width, height, padding int) {
	if divider != "" {
		result := patterns.CreateDivider(divider, width)
		outputResult(result)
		return
	}

	if border != "" {
		if text == "" {
			text = "Hello, World!"
		}
		result := patterns.CreateBorder(text, border, padding)
		outputResult(result)
		return
	}

	if preset != "" {
		result := patterns.GetPreset(preset, width, height)
		outputResult(result)
		return
	}

	fmt.Println("Use --divider, --border, or --preset to generate patterns.")
	fmt.Println("Run 'moji list-patterns' to see available options.")
}

func handleListPatterns() {
	fmt.Println("Border styles (use with --border):")
	for _, name := range patterns.ListBorders() {
		b := patterns.Borders[name]
		fmt.Printf("  %-10s  %s%s%s %s %s%s%s\n", name, b.TL, b.T, b.TR, b.L, b.BL, b.B, b.BR)
	}

	fmt.Println("\nDivider styles (use with --divider):")
	for _, name := range patterns.ListDividers() {
		d := patterns.Dividers[name]
		preview := strings.Repeat(string([]rune(d)[:1]), 10)
		if len([]rune(d)) > 1 {
			preview = patterns.CreateDivider(name, 10)
		}
		fmt.Printf("  %-12s  %s\n", name, preview)
	}

	fmt.Println("\nPattern presets (use with --preset):")
	for _, name := range patterns.ListPatterns() {
		p := patterns.PresetPatterns[name]
		fmt.Printf("  %-12s  %s\n", name, p)
	}
}

func handleAnimateList() {
	fmt.Println("Available animation presets:")
	for _, name := range animate.ListPresets() {
		frames, _ := animate.GetPreset(name)
		preview := strings.Join(frames[:min(4, len(frames))], " ")
		if len(frames) > 4 {
			preview += " ..."
		}
		fmt.Printf("  %-12s  %s\n", name, preview)
	}
	fmt.Println("\nOther effects:")
	fmt.Println("  --typewriter  Type text character by character")
	fmt.Println("  --scroll      Scroll text horizontally")
	fmt.Println("  --blink       Blink text on and off")
}

func handleAnimate(preset, text string, loops, delayMs int) {
	frames, ok := animate.GetPreset(preset)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown animation: '%s'. Use `moji animate --list` to see available.\n", preset)
		return
	}

	if text != "" {
		animate.PlayWithText(frames, text, loops, delayMs)
	} else {
		animate.Play(frames, loops, delayMs)
	}
	fmt.Println()
}

func handleTypewriter(text string, delayMs int) {
	if delayMs <= 0 {
		delayMs = 50
	}
	animate.Typewriter(text, delayMs)
}

func handleScroll(text string, width, loops, delayMs int) {
	if delayMs <= 0 {
		delayMs = 100
	}
	animate.ScrollText(text, width, loops, delayMs)
}

func handleBlink(text string, loops, delayMs int) {
	if delayMs <= 0 {
		delayMs = 500
	}
	animate.Blink(text, loops, delayMs)
}

func handleSysinfo(artName, gradientTheme string) {
	info := sysinfo.Collect()

	var art string
	if artName != "" {
		if a, ok := artdb.Get(artName); ok {
			art = a.Art
		}
	}

	if art == "" {
		art = sysinfo.GetOSLogo()
	}

	if gradientTheme != "" {
		art = gradient.Apply(art, gradientTheme, "diagonal")
	}

	output := info.FormatWithArt(art)
	fmt.Print(output)
}

func handleTree(path string, depth int, showHidden, dirsOnly bool) {
	opts := tree.DefaultOptions()
	opts.MaxDepth = depth
	opts.ShowHidden = showHidden
	opts.DirsOnly = dirsOnly

	entry, err := tree.Generate(path, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return
	}

	result := tree.Format(entry, opts)
	if copyFlag {
		if err := clipboard.WriteAll(result); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Println("Copied to clipboard!")
	}
	fmt.Print(result)
}

func handleCalendar(year, week, monday bool) {
	opts := calendar.Options{
		FirstDayMon: monday,
	}

	var result string
	if year {
		now := time.Now()
		result = calendar.Year(now.Year(), opts)
	} else if week {
		result = calendar.WeekView(opts)
	} else {
		result = calendar.Current(opts)
	}

	if copyFlag {
		plain := stripANSI(result)
		if err := clipboard.WriteAll(plain); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Println("Copied to clipboard!")
	}
	fmt.Print(result)
}

func handleConfigInit() {
	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating config: %v\n", err)
		return
	}
	fmt.Printf("Config file created at: %s\n", config.ConfigPath())
}

func handleConfigShow() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		return
	}

	fmt.Println("Current configuration:")
	fmt.Println()
	fmt.Println("Defaults:")
	fmt.Printf("  Banner font: %s\n", cfg.Defaults.BannerFont)
	fmt.Printf("  Banner style: %s\n", cfg.Defaults.BannerStyle)
	fmt.Printf("  Gradient theme: %s\n", cfg.Defaults.GradientTheme)
	fmt.Printf("  Bubble style: %s\n", cfg.Defaults.BubbleStyle)
	fmt.Printf("  Convert width: %d\n", cfg.Defaults.ConvertWidth)

	if len(cfg.Presets) > 0 {
		fmt.Println("\nPresets:")
		for name, preset := range cfg.Presets {
			fmt.Printf("  %s: %s\n", name, preset.Description)
		}
	}

	fmt.Printf("\nConfig path: %s\n", config.ConfigPath())
}

func handleTermInfo() {
	caps := terminal.Detect()

	fmt.Println("Terminal Capabilities:")
	fmt.Printf("  Terminal: %s\n", caps.Term)
	fmt.Printf("  Size: %dx%d\n", caps.Width, caps.Height)
	fmt.Printf("  Color: %s\n", caps.ColorLevel.String())
	fmt.Printf("  Unicode: %v\n", caps.Unicode)
	fmt.Printf("  Interactive: %v\n", caps.IsInteractive)

	if caps.IsTmux {
		fmt.Println("  Running in: tmux")
	}
	if caps.IsScreen {
		fmt.Println("  Running in: screen")
	}
	if caps.IsSSH {
		fmt.Println("  Running via: SSH")
	}

	fmt.Println("\nGraphics protocols:")
	fmt.Printf("  Sixel: %v\n", caps.Sixel)
	fmt.Printf("  Kitty: %v\n", caps.Kitty)
	fmt.Printf("  iTerm2: %v\n", caps.ITerm2)

	if caps.SupportsTrueColor() {
		fmt.Println("\nColor demo (truecolor):")
		for i := 0; i < 16; i++ {
			r := uint8(i * 16)
			g := uint8((15 - i) * 16)
			b := uint8(128)
			fmt.Printf("\033[48;2;%d;%d;%dm  \033[0m", r, g, b)
		}
		fmt.Println()
	}
}

func handleDoctor() {
	checks := ux.Doctor()
	fmt.Print(ux.FormatDoctorReport(checks))
}
