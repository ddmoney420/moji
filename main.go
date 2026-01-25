package main

import (
	"os"

	"github.com/ddmoney420/moji/internal/tui"
	"github.com/ddmoney420/moji/internal/ux"
	"github.com/spf13/cobra"
)

var (
	copyFlag    bool
	fontFlag    string
	styleFlag   string
	borderFlag  string
	alignFlag   string
	widthFlag   int
	outputFlag  string
	jsonFlag    bool
	bgColorFlag string
	fgColorFlag string
	quietFlag   bool
	verboseFlag bool
	noColorFlag bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "moji [name]",
		Short: "CLI tool for kaomoji, ASCII banners, emoji, and ASCII art",
		Long: `moji is a CLI tool for looking up kaomoji, generating ASCII banners,
converting smileys to emoji, and displaying ASCII art.

Examples:
  moji shrug                    # Get a kaomoji by name
  moji banner "Hello"           # Create ASCII art banner
  moji artdb cat                # Display ASCII art
  moji convert image.png        # Convert image to ASCII
  moji interactive              # Launch interactive studio

Run 'moji --help' for full command list, or 'moji doctor' to check your setup.`,
		Args:    cobra.MaximumNArgs(1),
		Version: ux.Version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ux.Quiet = quietFlag
			ux.Verbose = verboseFlag
			if noColorFlag {
				ux.NoColor = true
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 1 {
				handleGet(args[0])
			} else {
				cmd.Help()
			}
		},
	}

	rootCmd.SetVersionTemplate(ux.VersionLong() + "\n")

	// Global flags
	rootCmd.PersistentFlags().BoolVar(&copyFlag, "copy", false, "Copy output to clipboard")
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&noColorFlag, "no-color", false, "Disable colored output")

	// Interactive TUI command
	interactiveCmd := &cobra.Command{
		Use:     "interactive",
		Aliases: []string{"i", "ui", "studio"},
		Short:   "Launch interactive ASCII art studio",
		Long:    `Launch an interactive TUI for creating ASCII art with live preview, font selection, and more.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := tui.Run(); err != nil {
				os.Exit(1)
			}
		},
	}

	rootCmd.AddCommand(
		// Kaomoji
		newGetCmd(),
		newListCmd(),
		newRandomCmd(),
		newCategoriesCmd(),
		newEmojiCmd(),
		// Banner
		newBannerCmd(),
		newListFontsCmd(),
		newPreviewCmd(),
		// Effects & Filters
		newEffectsCmd(),
		newListEffectsCmd(),
		newFilterCmd(),
		newLolcatCmd(),
		// Art
		newArtCmd(),
		newArtdbCmd(),
		newFortuneCmd(),
		newSayCmd(),
		// Convert
		newConvertCmd(),
		newListCharsetsCmd(),
		newBatchCmd(),
		// Tools
		newQRCmd(),
		newListQRCharsetsCmd(),
		newGradientCmd(),
		newListThemesCmd(),
		newPatternCmd(),
		newListPatternsCmd(),
		newAnimateCmd(),
		newSysinfoCmd(),
		newTreeCmd(),
		newCalCmd(),
		newConfigCmd(),
		newTermCmd(),
		newDoctorCmd(),
		newDemoCmd(),
		newCompletionsCmd(rootCmd),
		interactiveCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
