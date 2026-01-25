package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

func stripANSI(s string) string {
	var result strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\033' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		result.WriteRune(r)
	}
	return result.String()
}

func outputResult(s string) {
	if copyFlag {
		if err := clipboard.WriteAll(s); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to copy to clipboard: %v\n", err)
			return
		}
		fmt.Println("Copied to clipboard!")
	}
	fmt.Print(s)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
