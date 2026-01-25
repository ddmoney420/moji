package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/ddmoney420/moji/internal/artdb"
	"github.com/ddmoney420/moji/internal/fortune"
	"github.com/ddmoney420/moji/internal/gradient"
	"github.com/ddmoney420/moji/internal/kaomoji"
	"github.com/ddmoney420/moji/internal/speech"
	"github.com/spf13/cobra"
)

func newArtCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "art [name]",
		Short: "Display ASCII art",
		Run: func(cmd *cobra.Command, args []string) {
			list, _ := cmd.Flags().GetBool("list")
			category, _ := cmd.Flags().GetString("category")
			if list || len(args) == 0 {
				handleArtList(category)
			} else {
				handleArt(args[0])
			}
		},
	}
	cmd.Flags().Bool("list", false, "List all available art")
	cmd.Flags().StringP("category", "c", "", "Filter by category")
	return cmd
}

func newArtdbCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "artdb [name]",
		Short: "Browse ASCII art database",
		Long: `Browse and display ASCII art from the built-in database.

Examples:
  moji artdb cat
  moji artdb --list
  moji artdb --search heart
  moji artdb --category animals`,
		Run: func(cmd *cobra.Command, args []string) {
			listFlag, _ := cmd.Flags().GetBool("list")
			search, _ := cmd.Flags().GetString("search")
			category, _ := cmd.Flags().GetString("category")
			categories, _ := cmd.Flags().GetBool("categories")
			gradientTheme, _ := cmd.Flags().GetString("gradient")

			if categories {
				handleArtDBCategories()
			} else if listFlag || search != "" || category != "" {
				handleArtDBList(search, category)
			} else if len(args) > 0 {
				handleArtDBGet(args[0], gradientTheme)
			} else {
				handleArtDBList("", "")
			}
		},
	}
	cmd.Flags().Bool("list", false, "List all art")
	cmd.Flags().String("search", "", "Search by name or tags")
	cmd.Flags().String("category", "", "Filter by category")
	cmd.Flags().Bool("categories", false, "List all categories")
	cmd.Flags().String("gradient", "", "Apply gradient theme")
	return cmd
}

func newFortuneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fortune",
		Short: "Display a random fortune with optional character",
		Long: `Display a random fortune, optionally with a Muppet character speaking it.

Examples:
  moji fortune
  moji fortune --joke
  moji fortune --character kermit
  moji fortune --character kermit --bubble think`,
		Run: func(cmd *cobra.Command, args []string) {
			joke, _ := cmd.Flags().GetBool("joke")
			character, _ := cmd.Flags().GetString("character")
			bubbleStyle, _ := cmd.Flags().GetString("bubble")
			handleFortune(joke, character, bubbleStyle)
		},
	}
	cmd.Flags().Bool("joke", false, "Get a programming joke instead")
	cmd.Flags().String("character", "", "Muppet character to speak the fortune")
	cmd.Flags().String("bubble", "round", "Speech bubble style: round, square, double, thick, ascii, think")
	return cmd
}

func newSayCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "say [text]",
		Short: "Display text in a speech bubble with optional character",
		Long: `Display text in a speech bubble, optionally with a Muppet character.

Examples:
  moji say "Hello World"
  moji say "Hello World" --character kermit
  moji say "Hmm..." --character kermit --bubble think
  echo "Hello" | moji say --character piggy`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			character, _ := cmd.Flags().GetString("character")
			bubbleStyle, _ := cmd.Flags().GetString("bubble")
			width, _ := cmd.Flags().GetInt("width")

			var text string
			if len(args) > 0 {
				text = args[0]
			} else {
				buf := make([]byte, 1024*1024)
				n, _ := os.Stdin.Read(buf)
				text = strings.TrimSpace(string(buf[:n]))
			}
			handleSay(text, character, bubbleStyle, width)
		},
	}
	cmd.Flags().String("character", "", "Muppet character to display below bubble")
	cmd.Flags().String("bubble", "round", "Bubble style: round, square, double, thick, ascii, think")
	cmd.Flags().Int("width", 40, "Maximum width of speech bubble")
	return cmd
}

func handleArtList(category string) {
	fmt.Println("Available ASCII art:")
	for _, item := range kaomoji.ListArt(category) {
		fmt.Printf("  %-15s - %s\n", item.Name, item.Category)
	}
}

func handleArt(name string) {
	art, ok := kaomoji.GetArt(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown art: '%s'. Use `moji art --list` to see available.\n", name)
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

func handleArtDBCategories() {
	fmt.Println("Available categories:")
	for _, cat := range artdb.ListCategories() {
		count := len(artdb.ByCategory(cat))
		fmt.Printf("  %-12s (%d)\n", cat, count)
	}
}

func handleArtDBList(search, category string) {
	var arts []artdb.Art

	if search != "" {
		arts = artdb.Search(search)
	} else if category != "" {
		arts = artdb.ByCategory(category)
	} else {
		arts = artdb.List()
	}

	if len(arts) == 0 {
		fmt.Println("No art found.")
		return
	}

	fmt.Println("Available ASCII art:")
	for _, a := range arts {
		tags := strings.Join(a.Tags, ", ")
		fmt.Printf("  %-12s [%s] - %s\n", a.Name, a.Category, tags)
	}
}

func handleArtDBGet(name, gradientTheme string) {
	art, ok := artdb.Get(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "Art '%s' not found. Use `moji artdb --list` to see available.\n", name)
		results := artdb.Search(name)
		if len(results) > 0 {
			fmt.Fprintf(os.Stderr, "Did you mean: ")
			for i, r := range results {
				if i > 0 {
					fmt.Fprintf(os.Stderr, ", ")
				}
				fmt.Fprintf(os.Stderr, "%s", r.Name)
				if i >= 4 {
					break
				}
			}
			fmt.Fprintln(os.Stderr)
		}
		return
	}

	output := art.Art
	if gradientTheme != "" {
		output = gradient.Apply(output, gradientTheme, "diagonal")
	}

	if copyFlag {
		plain := stripANSI(output)
		if err := clipboard.WriteAll(plain); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Println("Copied to clipboard!")
	}
	fmt.Print(output)
}

func handleFortune(joke bool, character, bubbleStyle string) {
	var text string
	if joke {
		text = fortune.GetJoke()
	} else {
		text = fortune.Get()
	}

	if character != "" {
		artEntry, ok := artdb.Get(character)
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown character: '%s'. Using fortune without character.\n", character)
			fmt.Println(text)
			return
		}

		bubble := speech.Wrap(text, bubbleStyle, 40)
		output := speech.Combine(bubble, artEntry.Art)

		if copyFlag {
			if err := clipboard.WriteAll(output); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
				return
			}
			fmt.Println("Copied to clipboard!")
		}
		fmt.Print(output)
	} else {
		if copyFlag {
			if err := clipboard.WriteAll(text); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
				return
			}
			fmt.Printf("Copied to clipboard: %s\n", text)
		} else {
			fmt.Println(text)
		}
	}
}

func handleSay(text, character, bubbleStyle string, width int) {
	bubble := speech.Wrap(text, bubbleStyle, width)

	var result string
	if character != "" {
		artEntry, ok := artdb.Get(character)
		if !ok {
			fmt.Fprintf(os.Stderr, "Unknown character: '%s'. Showing bubble only.\n", character)
			result = bubble
		} else {
			result = speech.Combine(bubble, artEntry.Art)
		}
	} else {
		result = bubble
	}

	if jsonFlag {
		data := map[string]string{"text": text, "output": result}
		json.NewEncoder(os.Stdout).Encode(data)
		return
	}

	if copyFlag {
		if err := clipboard.WriteAll(result); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Println("Copied to clipboard!")
	}
	fmt.Print(result)
}
