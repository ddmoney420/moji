package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/ddmoney420/moji/internal/effects"
	"github.com/ddmoney420/moji/internal/filters"
	"github.com/spf13/cobra"
)

func newEffectsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "effect [effect] [text]",
		Short: "Apply text effects: flip, reverse, mirror, wave, zalgo",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			handleEffect(args[0], args[1])
		},
	}
}

func newListEffectsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list-effects",
		Short: "List all available text effects",
		Run: func(cmd *cobra.Command, args []string) {
			handleListEffects()
		},
	}
}

func newFilterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filter [filters] [text]",
		Short: "Apply filters to text (metal, rainbow, glitch, etc.)",
		Long: `Apply one or more filters to text. Chain filters with commas.

Examples:
  moji filter rainbow "Hello World"
  moji filter metal,border "Text"
  echo "Hello" | moji filter glitch
  moji banner "Hi" | moji filter neon`,
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			listFlag, _ := cmd.Flags().GetBool("list")
			if listFlag {
				handleFilterList()
				return
			}

			if len(args) < 1 {
				handleFilterList()
				return
			}

			filterSpec := args[0]
			var text string
			if len(args) > 1 {
				text = args[1]
			} else {
				buf := make([]byte, 1024*1024)
				n, _ := os.Stdin.Read(buf)
				text = string(buf[:n])
			}
			handleFilter(filterSpec, text)
		},
	}
	cmd.Flags().Bool("list", false, "List available filters")
	return cmd
}

func newLolcatCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lolcat [text]",
		Short: "Rainbow color text (lolcat-style)",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			animateFlag, _ := cmd.Flags().GetBool("animate")
			speed, _ := cmd.Flags().GetFloat64("speed")

			var text string
			if len(args) > 0 {
				text = args[0]
			} else {
				buf := make([]byte, 1024*1024)
				n, _ := os.Stdin.Read(buf)
				text = string(buf[:n])
			}
			handleLolcat(text, animateFlag, speed)
		},
	}
	cmd.Flags().BoolP("animate", "a", false, "Animate the rainbow")
	cmd.Flags().Float64P("speed", "s", 0.1, "Animation speed")
	return cmd
}

func handleEffect(effect, text string) {
	result := effects.Apply(effect, text)

	if jsonFlag {
		data := map[string]string{"effect": effect, "input": text, "output": result}
		json.NewEncoder(os.Stdout).Encode(data)
		return
	}

	if copyFlag {
		if err := clipboard.WriteAll(result); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Printf("Copied to clipboard: %s\n", result)
	} else {
		fmt.Println(result)
	}
}

func handleListEffects() {
	fmt.Println("Available text effects:")
	for _, e := range effects.ListEffects() {
		fmt.Printf("  %-15s - %s\n", e.Name, e.Desc)
	}
}

func handleFilterList() {
	fmt.Println("Available filters:")
	for _, f := range filters.ListFilters() {
		fmt.Printf("  %-12s - %s\n", f.Name, f.Desc)
	}
	fmt.Println("\nChain multiple filters with commas: moji filter rainbow,border \"text\"")
}

func handleFilter(filterSpec, text string) {
	filterNames := filters.ParseChain(filterSpec)
	result := filters.Chain(text, filterNames)

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

func handleLolcat(text string, animate bool, speed float64) {
	if animate {
		lines := strings.Split(text, "\n")
		for cycle := 0; cycle < 50; cycle++ {
			fmt.Print("\033[H")
			for lineIdx, line := range lines {
				for i, r := range line {
					if r == ' ' || r == '\t' {
						fmt.Print(string(r))
						continue
					}
					phase := (float64(i+lineIdx) + float64(cycle)*speed) * 0.1
					red := uint8(math.Sin(phase)*127 + 128)
					green := uint8(math.Sin(phase+2*math.Pi/3)*127 + 128)
					blue := uint8(math.Sin(phase+4*math.Pi/3)*127 + 128)
					fmt.Printf("\033[38;2;%d;%d;%dm%c\033[0m", red, green, blue, r)
				}
				fmt.Println()
			}
			time.Sleep(50 * time.Millisecond)
		}
	} else {
		result := filters.Rainbow(text)
		fmt.Println(result)
	}
}
