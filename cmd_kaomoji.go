package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/atotto/clipboard"
	"github.com/ddmoney420/moji/internal/kaomoji"
	"github.com/ddmoney420/moji/internal/ux"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [name]",
		Short: "Look up a kaomoji by name",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handleGet(args[0])
		},
	}
}

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available kaomoji",
		Run: func(cmd *cobra.Command, args []string) {
			search, _ := cmd.Flags().GetString("search")
			category, _ := cmd.Flags().GetString("category")
			handleList(search, category)
		},
	}
	cmd.Flags().StringP("search", "s", "", "Filter by search term")
	cmd.Flags().StringP("category", "c", "", "Filter by category")
	return cmd
}

func newRandomCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "random",
		Short: "Get a random kaomoji",
		Run: func(cmd *cobra.Command, args []string) {
			handleRandom()
		},
	}
}

func newCategoriesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "categories",
		Short: "List all kaomoji categories",
		Run: func(cmd *cobra.Command, args []string) {
			handleCategories()
		},
	}
}

func newEmojiCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "emoji [smiley]",
		Short: "Convert ASCII smiley to emoji",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			handleEmoji(args[0])
		},
	}
}

func handleGet(name string) {
	k, ok := kaomoji.Get(name)
	if !ok {
		suggestions := kaomoji.Suggest(name)
		if len(suggestions) > 0 {
			ux.ErrorWithSuggestion(
				fmt.Sprintf("Unknown kaomoji: '%s'", name),
				ux.DidYouMean(suggestions),
			)
		} else {
			ux.ErrorWithCommand(
				fmt.Sprintf("Unknown kaomoji: '%s'", name),
				"moji list",
			)
		}
		return
	}

	if jsonFlag {
		data := map[string]string{"name": name, "kaomoji": k}
		json.NewEncoder(os.Stdout).Encode(data)
		return
	}

	if copyFlag {
		if err := clipboard.WriteAll(k); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Printf("Copied to clipboard: %s\n", k)
	} else {
		fmt.Println(k)
	}
}

func handleList(search, category string) {
	items := kaomoji.List(search, category)

	if jsonFlag {
		json.NewEncoder(os.Stdout).Encode(items)
		return
	}

	for _, item := range items {
		fmt.Printf("%-15s %s\n", item.Name, item.Kaomoji)
	}
}

func handleRandom() {
	name, k := kaomoji.Random()

	if jsonFlag {
		data := map[string]string{"name": name, "kaomoji": k}
		json.NewEncoder(os.Stdout).Encode(data)
		return
	}

	if copyFlag {
		if err := clipboard.WriteAll(k); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Printf("Copied to clipboard: %s (%s)\n", k, name)
	} else {
		fmt.Printf("%s (%s)\n", k, name)
	}
}

func handleCategories() {
	categories := kaomoji.ListCategories()

	if jsonFlag {
		json.NewEncoder(os.Stdout).Encode(categories)
		return
	}

	fmt.Println("Kaomoji categories:")
	for _, cat := range categories {
		fmt.Printf("  %s\n", cat)
	}
}

func handleEmoji(smiley string) {
	e, ok := kaomoji.SmileyToEmoji(smiley)
	if !ok {
		fmt.Fprintf(os.Stderr, "Unknown smiley: '%s'\n", smiley)
		return
	}

	if jsonFlag {
		data := map[string]string{"smiley": smiley, "emoji": e}
		json.NewEncoder(os.Stdout).Encode(data)
		return
	}

	if copyFlag {
		if err := clipboard.WriteAll(e); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Printf("Copied to clipboard: %s\n", e)
	} else {
		fmt.Println(e)
	}
}
