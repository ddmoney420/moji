package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/ddmoney420/moji/internal/banner"
	"github.com/ddmoney420/moji/internal/export"
	"github.com/ddmoney420/moji/internal/gradient"
	"github.com/ddmoney420/moji/internal/styles"
	"github.com/ddmoney420/moji/internal/watch"
	"github.com/spf13/cobra"
)

func newBannerCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "banner [text]",
		Short: "Generate ASCII art banner from text",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			gradientTheme, _ := cmd.Flags().GetString("gradient")
			watchFlag, _ := cmd.Flags().GetBool("watch")
			if watchFlag {
				handleBannerWatch(args[0], gradientTheme)
			} else {
				handleBanner(args[0], gradientTheme)
			}
		},
	}
	cmd.Flags().StringVarP(&fontFlag, "font", "f", "standard", "Font name")
	cmd.Flags().StringVarP(&styleFlag, "style", "s", "none", "Color style")
	cmd.Flags().StringVarP(&borderFlag, "border", "b", "none", "Border style: single, double, round, bold, ascii, stars, hash")
	cmd.Flags().StringVarP(&alignFlag, "align", "a", "left", "Text alignment: left, center, right")
	cmd.Flags().IntVarP(&widthFlag, "width", "w", 0, "Output width (0 for auto)")
	cmd.Flags().StringVarP(&outputFlag, "output", "o", "", "Save to file (supports .png, .txt)")
	cmd.Flags().StringVar(&bgColorFlag, "bg", "#282a36", "Background color for PNG export (hex)")
	cmd.Flags().StringVar(&fgColorFlag, "fg", "#f8f8f2", "Foreground color for PNG export (hex)")
	cmd.Flags().String("gradient", "", "Apply color gradient theme (rainbow, neon, fire, etc.)")
	cmd.Flags().BoolVar(&watchFlag, "watch", false, "Watch for changes and re-render in real-time")
	return cmd
}

func newListFontsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-fonts",
		Short: "List available banner fonts and styles",
		Run: func(cmd *cobra.Command, args []string) {
			handleListFonts()
		},
	}
}

func newPreviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preview [text]",
		Short: "Preview text in multiple fonts",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			limit, _ := cmd.Flags().GetInt("limit")
			category, _ := cmd.Flags().GetString("category")
			handlePreview(args[0], limit, category)
		},
	}
	cmd.Flags().IntP("limit", "l", 5, "Number of fonts to preview")
	cmd.Flags().StringP("category", "c", "", "Font category: 3d, graffiti, retro, big, clean, decorative, fun, small")
	return cmd
}

func handleBannerWatch(text string, gradientTheme string) {
	fmt.Println("Watching for changes (Press Ctrl+C to exit)...")
	renderBanner := func() {
		fmt.Print("\033[2J\033[H")
		handleBanner(text, gradientTheme)
	}

	err := watch.Watch(".", renderBanner)
	if err != nil && err.Error() != "signal: interrupt" {
		fmt.Fprintf(os.Stderr, "Watch error: %v\n", err)
	}
}

func handleBanner(text string, gradientTheme string) {
	art, err := banner.Generate(text, fontFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating banner: %v\n", err)
		return
	}

	if widthFlag > 0 {
		art = styles.ApplyAlignment(art, alignFlag, widthFlag)
	}

	art = styles.ApplyBorder(art, borderFlag)

	if gradientTheme != "" {
		art = gradient.Apply(art, gradientTheme, "horizontal")
	}

	styledArt := art
	if gradientTheme == "" {
		styledArt = styles.Apply(art, styleFlag)
	}

	if outputFlag != "" {
		lower := strings.ToLower(outputFlag)
		switch {
		case strings.HasSuffix(lower, ".png"):
			if err := export.ToPNG(art, outputFlag, bgColorFlag, fgColorFlag); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to export PNG: %v\n", err)
				return
			}
		case strings.HasSuffix(lower, ".svg"):
			if err := export.ToSVG(art, outputFlag, bgColorFlag, fgColorFlag); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to export SVG: %v\n", err)
				return
			}
		case strings.HasSuffix(lower, ".html"), strings.HasSuffix(lower, ".htm"):
			if err := export.ToHTML(art, outputFlag, bgColorFlag, fgColorFlag, text); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to export HTML: %v\n", err)
				return
			}
		default:
			if err := os.WriteFile(outputFlag, []byte(art), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to write file: %v\n", err)
				return
			}
		}
		fmt.Printf("Saved banner to %s\n", outputFlag)
		return
	}

	if jsonFlag {
		data := map[string]string{"text": text, "font": fontFlag, "style": styleFlag, "output": art}
		json.NewEncoder(os.Stdout).Encode(data)
		return
	}

	if copyFlag {
		if err := clipboard.WriteAll(art); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Println("Copied banner to clipboard!")
	} else {
		fmt.Print(styledArt)
	}
}

func handleListFonts() {
	fmt.Println("Available fonts:")
	for _, f := range banner.ListFonts() {
		fmt.Printf("  %-12s - %s\n", f.Name, f.Desc)
	}

	fmt.Println("\nColor styles (use with --style):")
	for _, s := range styles.ListStyles() {
		fmt.Printf("  %-12s - %s\n", s.Name, s.Desc)
	}

	fmt.Println("\nBorder styles (use with --border):")
	for _, b := range styles.ListBorders() {
		fmt.Printf("  %-12s - %s\n", b.Name, b.Desc)
	}
}

func handlePreview(text string, limit int, category string) {
	fonts := banner.ListFonts()
	count := 0

	for _, f := range fonts {
		if category != "" && !strings.Contains(strings.ToLower(f.Desc), strings.ToLower(category)) {
			continue
		}

		art, err := banner.Generate(text, f.Name)
		if err != nil {
			continue
		}

		fmt.Printf("\n=== %s (%s) ===\n", f.Name, f.Desc)
		fmt.Println(art)

		count++
		if count >= limit {
			break
		}
	}
}
