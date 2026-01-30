package main

import (
	"encoding/json"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/ddmoney420/moji/internal/batch"
	"github.com/ddmoney420/moji/internal/convert"
	"github.com/ddmoney420/moji/internal/export"
	"github.com/ddmoney420/moji/internal/imgproto"
	"github.com/ddmoney420/moji/internal/ux"
	"github.com/ddmoney420/moji/internal/watch"
	"github.com/spf13/cobra"
)

func newConvertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert [image]",
		Short: "Convert an image to ASCII art",
		Long: `Convert an image file or URL to ASCII art.

Supported formats: PNG, JPEG, GIF, BMP, WebP

Examples:
  moji convert photo.jpg
  moji convert --url https://example.com/cat.png
  moji convert photo.png --width 100 --edge
  moji convert photo.png --color --charset blocks
  moji convert photo.png -o art.txt`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			url, _ := cmd.Flags().GetString("url")
			width, _ := cmd.Flags().GetInt("width")
			height, _ := cmd.Flags().GetInt("height")
			charset, _ := cmd.Flags().GetString("charset")
			edge, _ := cmd.Flags().GetBool("edge")
			colorFlag, _ := cmd.Flags().GetBool("color")
			invert, _ := cmd.Flags().GetBool("invert")
			protocol, _ := cmd.Flags().GetString("protocol")
			watchFlag, _ := cmd.Flags().GetBool("watch")

			if len(args) == 0 && url == "" {
				fmt.Fprintln(os.Stderr, "Error: provide an image file or --url")
				cmd.Help()
				return
			}

			var source string
			if len(args) > 0 {
				source = args[0]
			}

			if watchFlag && source != "" {
				handleConvertWatch(source, url, width, height, charset, edge, colorFlag, invert, protocol)
			} else {
				handleConvert(source, url, width, height, charset, edge, colorFlag, invert, protocol)
			}
		},
	}
	cmd.Flags().String("url", "", "Image URL to convert")
	cmd.Flags().Int("width", 80, "Output width in characters")
	cmd.Flags().Int("height", 0, "Output height (0 = auto based on aspect ratio)")
	cmd.Flags().String("charset", "standard", "Character set: standard, blocks, simple, detailed, binary, dots, ascii, shade")
	cmd.Flags().Bool("edge", false, "Use edge detection for line art style")
	cmd.Flags().Bool("color", false, "Preserve colors using ANSI codes")
	cmd.Flags().Bool("invert", false, "Invert brightness (for light backgrounds)")
	cmd.Flags().String("protocol", "ascii", "Image protocol: ascii, sixel, kitty, iterm2, auto")
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Save to file")
	cmd.Flags().BoolVar(&watchFlag, "watch", false, "Watch for changes and re-render in real-time")
	return cmd
}

func newListCharsetsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-charsets",
		Short: "List available character sets for image conversion",
		Run: func(cmd *cobra.Command, args []string) {
			handleListCharsets()
		},
	}
}

func newBatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "batch [pattern...]",
		Short: "Batch convert multiple images to ASCII art",
		Long: `Convert multiple images to ASCII art at once.

Examples:
  moji batch "*.jpg" "*.png"
  moji batch photos/*.jpg --width 60 --output-dir ascii_art
  moji batch *.png --charset blocks --workers 8`,
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			width, _ := cmd.Flags().GetInt("width")
			height, _ := cmd.Flags().GetInt("height")
			charset, _ := cmd.Flags().GetString("charset")
			outDir, _ := cmd.Flags().GetString("output-dir")
			workers, _ := cmd.Flags().GetInt("workers")
			color, _ := cmd.Flags().GetBool("color")
			handleBatch(args, width, height, charset, outDir, workers, color)
		},
	}
	cmd.Flags().Int("width", 80, "Output width")
	cmd.Flags().Int("height", 0, "Output height (0 = auto)")
	cmd.Flags().String("charset", "standard", "Character set")
	cmd.Flags().String("output-dir", "", "Output directory (prints to stdout if empty)")
	cmd.Flags().Int("workers", 4, "Number of concurrent workers")
	cmd.Flags().Bool("color", false, "Preserve colors")
	return cmd
}

func handleConvertWatch(file, url string, width, height int, charset string, edge, color, invert bool, protocol string) {
	fmt.Println("Watching for changes (Press Ctrl+C to exit)...")
	renderConvert := func() {
		fmt.Print("\033[2J\033[H")
		handleConvert(file, url, width, height, charset, edge, color, invert, protocol)
	}

	err := watch.Watch(file, renderConvert)
	if err != nil && err.Error() != "signal: interrupt" {
		fmt.Fprintf(os.Stderr, "Watch error: %v\n", err)
	}
}

func handleConvert(file, url string, width, height int, charset string, edge, color, invert bool, protocol string) {
	proto := imgproto.ParseProtocol(protocol)

	if proto != imgproto.ASCII {
		var img image.Image
		var err error

		if url != "" {
			spinner := ux.NewSpinner("Fetching image from URL")
			spinner.Start()
			img, err = convert.LoadImageURL(url)
			if err != nil {
				spinner.StopFail("Failed to fetch image")
				ux.Error("Failed to load image: %v", err)
				return
			}
			spinner.StopSuccess("Image fetched")
		} else {
			img, err = convert.LoadImageFile(file)
			if err != nil {
				ux.Error("Failed to load image: %v", err)
				return
			}
		}

		if err := imgproto.WriteToTerminal(img, proto, width); err != nil {
			ux.Error("Failed to render image: %v", err)
		}
		return
	}

	opts := convert.Options{
		Width:      width,
		Height:     height,
		Charset:    convert.GetCharset(charset),
		EdgeDetect: edge,
		Color:      color,
		Invert:     invert,
	}

	var art string
	var err error

	if url != "" {
		spinner := ux.NewSpinner("Fetching image from URL")
		spinner.Start()
		art, err = convert.FromURL(url, opts)
		if err != nil {
			spinner.StopFail("Failed to fetch image")
		} else {
			spinner.StopSuccess("Image fetched")
		}
	} else {
		ux.Debug("Converting file: %s", file)
		art, err = convert.FromFile(file, opts)
	}

	if err != nil {
		ux.Error("Failed to convert image: %v", err)
		return
	}

	if outputFlag != "" {
		if strings.HasSuffix(strings.ToLower(outputFlag), ".png") {
			if err := export.ToPNG(art, outputFlag, bgColorFlag, fgColorFlag); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to export PNG: %v\n", err)
				return
			}
			fmt.Printf("Saved to %s\n", outputFlag)
		} else {
			if err := os.WriteFile(outputFlag, []byte(art), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write file: %v\n", err)
				return
			}
			fmt.Printf("Saved to %s\n", outputFlag)
		}
		return
	}

	if jsonFlag {
		source := file
		if url != "" {
			source = url
		}
		data := map[string]string{"source": source, "art": art}
		json.NewEncoder(os.Stdout).Encode(data)
		return
	}

	if copyFlag {
		if err := clipboard.WriteAll(art); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Println("Copied to clipboard!")
	}
	fmt.Println(art)
}

func handleListCharsets() {
	fmt.Println("Available character sets for image conversion:")
	fmt.Println()
	fmt.Println("=== Basic Sets ===")
	charsets := map[string]string{
		"standard":  "General purpose (default)",
		"simple":    "Minimal 4-character set",
		"ascii":     "Classic ASCII art style",
		"detailed":  "Extended ASCII, maximum detail",
		"binary":    "Just space and solid block",
		"blocks":    "Unicode block elements",
		"shade":     "Block shading with half-blocks",
		"braille":   "Braille patterns (256 levels)",
		"stipple":   "Braille stipple/noise pattern",
		"box":       "Box drawing characters",
		"hatch":     "Hatching/line patterns",
		"geometric": "Geometric shapes (circles, wedges)",
		"shapes":    "Filled squares and circles",
		"dots":      "Dot/circle progression",
		"stars":     "Stars and asterisks",
		"dingbats":  "Dingbats and bullets",
		"symbols":   "Misc symbols (cards, weather, etc)",
		"arrows":    "Arrow characters",
		"music":     "Musical notation symbols",
		"math":      "Mathematical operators",
		"elegant":   "Elegant minimal dots",
		"dense":     "CJK-inspired dense characters",
		"ultra":     "Ultra detailed (ASCII + Unicode mix)",
		"retro":     "Retro terminal style",
	}

	categories := []struct {
		name string
		sets []string
	}{
		{"Basic", []string{"standard", "simple", "ascii", "detailed", "binary"}},
		{"Blocks & Shading", []string{"blocks", "shade", "braille", "stipple"}},
		{"Line Drawing", []string{"box", "hatch"}},
		{"Geometric", []string{"geometric", "shapes", "dots"}},
		{"Decorative", []string{"stars", "dingbats", "symbols", "arrows", "music"}},
		{"Special", []string{"math", "elegant", "dense"}},
		{"Combined", []string{"ultra", "retro"}},
	}

	for _, cat := range categories {
		fmt.Printf("─── %s ───\n", cat.name)
		for _, name := range cat.sets {
			desc := charsets[name]
			cs := convert.GetCharset(name)
			runes := []rune(cs)
			preview := cs
			if len(runes) > 40 {
				preview = string(runes[:40]) + "…"
			}
			fmt.Printf("  %-10s  %s\n", name, desc)
			fmt.Printf("             %s\n", preview)
		}
		fmt.Println()
	}
}

func handleBatch(filePatterns []string, width, height int, charset, outDir string, workers int, color bool) {
	opts := convert.Options{
		Width:   width,
		Height:  height,
		Charset: convert.GetCharset(charset),
		Color:   color,
	}

	spinner := ux.NewSpinner("Processing images")
	spinner.Start()
	results := batch.ConvertImages(filePatterns, opts, workers)
	spinner.Stop()

	if len(results) == 0 {
		ux.Warn("No files found matching the patterns")
		return
	}

	ux.Success("Converted %d images", len(results))

	if outDir != "" {
		if err := batch.SaveResults(results, outDir, "txt"); err != nil {
			ux.Error("Failed to save results: %v", err)
			return
		}
		ux.Success("Saved results to %s/", outDir)
	} else {
		batch.PrintResults(results, true)
	}
}
