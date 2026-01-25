package figlet

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// Font represents a FIGlet font
type Font struct {
	Header     Header
	Characters map[rune][]string
}

// Header contains FIGlet font header information
type Header struct {
	HardBlank    rune
	Height       int
	Baseline     int
	MaxLength    int
	OldLayout    int
	CommentLines int
}

// ParseFont parses a FIGlet font from string content
func ParseFont(content string) (*Font, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))

	// Parse header
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty font file")
	}
	headerLine := scanner.Text()

	header, err := parseHeader(headerLine)
	if err != nil {
		return nil, err
	}

	// Skip comment lines
	for i := 0; i < header.CommentLines; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected end of file in comments")
		}
	}

	// Parse characters
	characters := make(map[rune][]string)

	// Standard ASCII characters 32-126
	for charCode := 32; charCode <= 126; charCode++ {
		lines := make([]string, 0, header.Height)
		for h := 0; h < header.Height; h++ {
			if !scanner.Scan() {
				// Some fonts don't have all characters
				break
			}
			line := scanner.Text()
			// Remove end markers (@ or @@)
			line = strings.TrimRight(line, "@")
			// Replace hardblank with space for display
			line = strings.ReplaceAll(line, string(header.HardBlank), " ")
			lines = append(lines, line)
		}
		if len(lines) == header.Height {
			characters[rune(charCode)] = lines
		}
	}

	return &Font{
		Header:     header,
		Characters: characters,
	}, nil
}

func parseHeader(line string) (Header, error) {
	var header Header

	// FIGlet header format: flf2a[hardblank] height baseline maxlen oldlayout commentlines [printdir hardblanksmush fullwidth]
	if !strings.HasPrefix(line, "flf2a") {
		return header, fmt.Errorf("invalid FIGlet font header")
	}

	// The hardblank is the character right after "flf2a"
	if len(line) < 6 {
		return header, fmt.Errorf("header too short")
	}
	header.HardBlank = rune(line[5])

	// Parse the rest of the header
	parts := strings.Fields(line[6:])
	if len(parts) < 4 {
		return header, fmt.Errorf("incomplete header")
	}

	var err error
	header.Height, err = strconv.Atoi(parts[0])
	if err != nil {
		return header, fmt.Errorf("invalid height: %v", err)
	}

	header.Baseline, err = strconv.Atoi(parts[1])
	if err != nil {
		return header, fmt.Errorf("invalid baseline: %v", err)
	}

	header.MaxLength, err = strconv.Atoi(parts[2])
	if err != nil {
		return header, fmt.Errorf("invalid max length: %v", err)
	}

	header.OldLayout, err = strconv.Atoi(parts[3])
	if err != nil {
		return header, fmt.Errorf("invalid old layout: %v", err)
	}

	if len(parts) >= 5 {
		header.CommentLines, _ = strconv.Atoi(parts[4])
	}

	return header, nil
}

// Render renders text using the font
func (f *Font) Render(text string) string {
	if len(text) == 0 {
		return ""
	}

	// Initialize output lines
	outputLines := make([]string, f.Header.Height)

	for _, char := range text {
		charLines, ok := f.Characters[char]
		if !ok {
			// Use space for unknown characters
			charLines = f.Characters[' ']
			if charLines == nil {
				// Create empty space
				charLines = make([]string, f.Header.Height)
				for i := range charLines {
					charLines[i] = " "
				}
			}
		}

		for i, line := range charLines {
			if i < len(outputLines) {
				outputLines[i] += line
			}
		}
	}

	return strings.Join(outputLines, "\n")
}
