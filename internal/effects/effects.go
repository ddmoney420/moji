package effects

import (
	"math/rand"
	"strings"
	"unicode"
)

// Flip map - characters that have upside-down equivalents
var flipMap = map[rune]rune{
	'a': 'ɐ', 'b': 'q', 'c': 'ɔ', 'd': 'p', 'e': 'ǝ', 'f': 'ɟ',
	'g': 'ƃ', 'h': 'ɥ', 'i': 'ᴉ', 'j': 'ɾ', 'k': 'ʞ', 'l': 'l',
	'm': 'ɯ', 'n': 'u', 'o': 'o', 'p': 'd', 'q': 'b', 'r': 'ɹ',
	's': 's', 't': 'ʇ', 'u': 'n', 'v': 'ʌ', 'w': 'ʍ', 'x': 'x',
	'y': 'ʎ', 'z': 'z',
	'A': '∀', 'B': 'B', 'C': 'Ɔ', 'D': 'D', 'E': 'Ǝ', 'F': 'Ⅎ',
	'G': 'פ', 'H': 'H', 'I': 'I', 'J': 'ſ', 'K': 'ʞ', 'L': '˥',
	'M': 'W', 'N': 'N', 'O': 'O', 'P': 'Ԁ', 'Q': 'Q', 'R': 'R',
	'S': 'S', 'T': '┴', 'U': '∩', 'V': 'Λ', 'W': 'M', 'X': 'X',
	'Y': '⅄', 'Z': 'Z',
	'1': 'Ɩ', '2': 'ᄅ', '3': 'Ɛ', '4': 'ㄣ', '5': 'ϛ', '6': '9',
	'7': 'ㄥ', '8': '8', '9': '6', '0': '0',
	'.': '˙', ',': '\'', '\'': ',', '"': '„', '`': ',',
	'?': '¿', '!': '¡', '[': ']', ']': '[', '(': ')', ')': '(',
	'{': '}', '}': '{', '<': '>', '>': '<', '&': '⅋',
	'_': '‾', ';': '؛', '∴': '∵',
}

// Zalgo combining characters
var zalgoUp = []rune{
	'\u030d', '\u030e', '\u0304', '\u0305', '\u033f', '\u0311', '\u0306',
	'\u0310', '\u0352', '\u0357', '\u0351', '\u0307', '\u0308', '\u030a',
	'\u0342', '\u0343', '\u0344', '\u034a', '\u034b', '\u034c', '\u0303',
	'\u0302', '\u030c', '\u0350', '\u0300', '\u0301', '\u030b', '\u030f',
	'\u0312', '\u0313', '\u0314', '\u033d', '\u0309', '\u0363', '\u0364',
	'\u0365', '\u0366', '\u0367', '\u0368', '\u0369', '\u036a', '\u036b',
	'\u036c', '\u036d', '\u036e', '\u036f', '\u0483', '\u0484', '\u0485',
	'\u0486', '\u0487',
}

var zalgoMid = []rune{
	'\u0315', '\u031b', '\u0340', '\u0341', '\u0358', '\u0321', '\u0322',
	'\u0327', '\u0328', '\u0334', '\u0335', '\u0336', '\u034f', '\u035c',
	'\u035d', '\u035e', '\u035f', '\u0360', '\u0362', '\u0338', '\u0337',
	'\u0361', '\u0489',
}

var zalgoDown = []rune{
	'\u0316', '\u0317', '\u0318', '\u0319', '\u031c', '\u031d', '\u031e',
	'\u031f', '\u0320', '\u0324', '\u0325', '\u0326', '\u0329', '\u032a',
	'\u032b', '\u032c', '\u032d', '\u032e', '\u032f', '\u0330', '\u0331',
	'\u0332', '\u0333', '\u0339', '\u033a', '\u033b', '\u033c', '\u0345',
	'\u0347', '\u0348', '\u0349', '\u034d', '\u034e', '\u0353', '\u0354',
	'\u0355', '\u0356', '\u0359', '\u035a', '\u0323',
}

// Apply applies a text effect
func Apply(effect, text string) string {
	switch strings.ToLower(effect) {
	case "flip", "upside-down", "upsidedown":
		return Flip(text)
	case "reverse", "backwards":
		return Reverse(text)
	case "mirror":
		return Mirror(text)
	case "wave", "wavy":
		return Wave(text)
	case "zalgo", "cursed", "glitch":
		return Zalgo(text, 3)
	case "zalgo-mild":
		return Zalgo(text, 1)
	case "zalgo-intense":
		return Zalgo(text, 6)
	case "bubble", "bubbles":
		return Bubble(text)
	case "square", "squares":
		return Square(text)
	case "bold":
		return Bold(text)
	case "italic":
		return Italic(text)
	case "strikethrough", "strike":
		return Strikethrough(text)
	case "underline":
		return Underline(text)
	case "smallcaps", "small-caps":
		return SmallCaps(text)
	case "fullwidth", "wide":
		return Fullwidth(text)
	case "monospace", "mono":
		return Monospace(text)
	case "script", "cursive":
		return Script(text)
	case "fraktur", "gothic":
		return Fraktur(text)
	case "double-struck", "doublestruck", "blackboard":
		return DoubleStruck(text)
	case "sparkle", "sparkles":
		return Sparkle(text)
	default:
		return text
	}
}

// Flip turns text upside down
func Flip(text string) string {
	runes := []rune(text)
	result := make([]rune, len(runes))

	for i, r := range runes {
		if flipped, ok := flipMap[r]; ok {
			result[len(runes)-1-i] = flipped
		} else {
			result[len(runes)-1-i] = r
		}
	}

	return string(result)
}

// Reverse reverses the text
func Reverse(text string) string {
	runes := []rune(text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// Mirror creates a mirrored version (original + reversed)
func Mirror(text string) string {
	return text + " | " + Reverse(text)
}

// Wave applies a wavy effect using superscripts/subscripts
func Wave(text string) string {
	var result strings.Builder
	for i, r := range text {
		if i%2 == 0 {
			result.WriteRune(r)
			result.WriteString("\u0302") // combining circumflex
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Zalgo applies zalgo/glitch text effect
func Zalgo(text string, intensity int) string {
	var result strings.Builder

	for _, r := range text {
		result.WriteRune(r)

		if unicode.IsSpace(r) {
			continue
		}

		// Add combining characters
		for i := 0; i < intensity; i++ {
			result.WriteRune(zalgoUp[rand.Intn(len(zalgoUp))])
		}
		for i := 0; i < intensity/2; i++ {
			result.WriteRune(zalgoMid[rand.Intn(len(zalgoMid))])
		}
		for i := 0; i < intensity; i++ {
			result.WriteRune(zalgoDown[rand.Intn(len(zalgoDown))])
		}
	}

	return result.String()
}

// Bubble converts text to circled letters
func Bubble(text string) string {
	var result strings.Builder
	for _, r := range text {
		if r >= 'a' && r <= 'z' {
			result.WriteRune(rune(0x24D0 + (r - 'a')))
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune(rune(0x24B6 + (r - 'A')))
		} else if r >= '1' && r <= '9' {
			result.WriteRune(rune(0x2460 + (r - '1')))
		} else if r == '0' {
			result.WriteRune(0x24EA)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Square converts text to squared letters
func Square(text string) string {
	var result strings.Builder
	for _, r := range text {
		if r >= 'A' && r <= 'Z' {
			result.WriteRune(rune(0x1F130 + (r - 'A')))
		} else if r >= 'a' && r <= 'z' {
			result.WriteRune(rune(0x1F130 + (r - 'a')))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Bold converts text to mathematical bold
func Bold(text string) string {
	var result strings.Builder
	for _, r := range text {
		if r >= 'a' && r <= 'z' {
			result.WriteRune(rune(0x1D41A + (r - 'a')))
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune(rune(0x1D400 + (r - 'A')))
		} else if r >= '0' && r <= '9' {
			result.WriteRune(rune(0x1D7CE + (r - '0')))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Italic converts text to mathematical italic
func Italic(text string) string {
	var result strings.Builder
	for _, r := range text {
		if r >= 'a' && r <= 'z' {
			if r == 'h' {
				result.WriteRune(0x210E) // special case for h
			} else {
				result.WriteRune(rune(0x1D44E + (r - 'a')))
			}
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune(rune(0x1D434 + (r - 'A')))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Strikethrough adds strikethrough to text
func Strikethrough(text string) string {
	var result strings.Builder
	for _, r := range text {
		result.WriteRune(r)
		if !unicode.IsSpace(r) {
			result.WriteRune('\u0336') // combining long stroke overlay
		}
	}
	return result.String()
}

// Underline adds underline to text
func Underline(text string) string {
	var result strings.Builder
	for _, r := range text {
		result.WriteRune(r)
		if !unicode.IsSpace(r) {
			result.WriteRune('\u0332') // combining low line
		}
	}
	return result.String()
}

// SmallCaps converts lowercase to small capitals
func SmallCaps(text string) string {
	smallCaps := map[rune]rune{
		'a': 'ᴀ', 'b': 'ʙ', 'c': 'ᴄ', 'd': 'ᴅ', 'e': 'ᴇ', 'f': 'ғ',
		'g': 'ɢ', 'h': 'ʜ', 'i': 'ɪ', 'j': 'ᴊ', 'k': 'ᴋ', 'l': 'ʟ',
		'm': 'ᴍ', 'n': 'ɴ', 'o': 'ᴏ', 'p': 'ᴘ', 'q': 'ǫ', 'r': 'ʀ',
		's': 's', 't': 'ᴛ', 'u': 'ᴜ', 'v': 'ᴠ', 'w': 'ᴡ', 'x': 'x',
		'y': 'ʏ', 'z': 'ᴢ',
	}

	var result strings.Builder
	for _, r := range text {
		if mapped, ok := smallCaps[r]; ok {
			result.WriteRune(mapped)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Fullwidth converts to fullwidth characters
func Fullwidth(text string) string {
	var result strings.Builder
	for _, r := range text {
		if r >= '!' && r <= '~' {
			result.WriteRune(rune(0xFF01 + (r - '!')))
		} else if r == ' ' {
			result.WriteRune(0x3000) // ideographic space
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Monospace converts to mathematical monospace
func Monospace(text string) string {
	var result strings.Builder
	for _, r := range text {
		if r >= 'a' && r <= 'z' {
			result.WriteRune(rune(0x1D68A + (r - 'a')))
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune(rune(0x1D670 + (r - 'A')))
		} else if r >= '0' && r <= '9' {
			result.WriteRune(rune(0x1D7F6 + (r - '0')))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Script converts to mathematical script
func Script(text string) string {
	var result strings.Builder
	for _, r := range text {
		if r >= 'a' && r <= 'z' {
			result.WriteRune(rune(0x1D4B6 + (r - 'a')))
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune(rune(0x1D49C + (r - 'A')))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Fraktur converts to mathematical fraktur
func Fraktur(text string) string {
	var result strings.Builder
	for _, r := range text {
		if r >= 'a' && r <= 'z' {
			result.WriteRune(rune(0x1D51E + (r - 'a')))
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune(rune(0x1D504 + (r - 'A')))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// DoubleStruck converts to double-struck (blackboard bold)
func DoubleStruck(text string) string {
	var result strings.Builder
	for _, r := range text {
		if r >= 'a' && r <= 'z' {
			result.WriteRune(rune(0x1D552 + (r - 'a')))
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune(rune(0x1D538 + (r - 'A')))
		} else if r >= '0' && r <= '9' {
			result.WriteRune(rune(0x1D7D8 + (r - '0')))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Sparkle adds sparkles around text
func Sparkle(text string) string {
	sparkles := []string{"✧", "✦", "★", "☆", "✨", "✩", "✪", "✫", "✬", "✭"}
	s1 := sparkles[rand.Intn(len(sparkles))]
	s2 := sparkles[rand.Intn(len(sparkles))]
	s3 := sparkles[rand.Intn(len(sparkles))]
	return s1 + " " + s2 + " " + text + " " + s2 + " " + s3
}

// ListEffects returns all available effects
func ListEffects() []struct{ Name, Desc string } {
	return []struct{ Name, Desc string }{
		{"flip", "Flip text upside down"},
		{"reverse", "Reverse text backwards"},
		{"mirror", "Mirror text with separator"},
		{"wave", "Wavy text effect"},
		{"zalgo", "Zalgo/glitch text (medium)"},
		{"zalgo-mild", "Mild zalgo effect"},
		{"zalgo-intense", "Intense zalgo effect"},
		{"bubble", "Circled letters Ⓐ Ⓑ Ⓒ"},
		{"square", "Squared letters 🅰 🅱 🅲"},
		{"bold", "Mathematical bold 𝐀𝐁𝐂"},
		{"italic", "Mathematical italic 𝐴𝐵𝐶"},
		{"strikethrough", "S̶t̶r̶i̶k̶e̶t̶h̶r̶o̶u̶g̶h̶"},
		{"underline", "U̲n̲d̲e̲r̲l̲i̲n̲e̲"},
		{"smallcaps", "ꜱᴍᴀʟʟ ᴄᴀᴘꜱ"},
		{"fullwidth", "Ｆｕｌｌｗｉｄｔｈ"},
		{"monospace", "𝙼𝚘𝚗𝚘𝚜𝚙𝚊𝚌𝚎"},
		{"script", "𝒮𝒸𝓇𝒾𝓅𝓉"},
		{"fraktur", "𝔉𝔯𝔞𝔨𝔱𝔲𝔯"},
		{"double-struck", "𝔻𝕠𝕦𝕓𝕝𝕖-𝕊𝕥𝕣𝕦𝕔𝕜"},
		{"sparkle", "✧ ★ Sparkle ★ ✧"},
	}
}
