package kaomoji

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

type KaomojiEntry struct {
	Name     string
	Kaomoji  string
	Category string
}

var kaomojis []KaomojiEntry
var kaomojiMap map[string]string
var smileys map[string]string
var asciiArt map[string]ArtEntry

type ArtEntry struct {
	Art      string
	Category string
}

func init() {
	rand.Seed(time.Now().UnixNano())

	kaomojis = []KaomojiEntry{
		// Classic expressions
		{"shrug", "¯\\_(ツ)_/¯", "expressions"},
		{"tableflip", "(╯°□°）╯︵ ┻━┻", "expressions"},
		{"unflip", "┬─┬ノ( º _ ºノ)", "expressions"},
		{"lenny", "( ͡° ͜ʖ ͡°)", "expressions"},
		{"disapproval", "ಠ_ಠ", "expressions"},
		{"why", "ლ(ಠ益ಠლ)", "expressions"},
		{"cry", "(╥﹏╥)", "emotions"},
		{"happy", "(◕‿◕)", "emotions"},
		{"sad", "(︶︹︺)", "emotions"},
		{"angry", "(ノಠ益ಠ)ノ", "emotions"},
		{"love", "(♥‿♥)", "emotions"},
		{"wink", "(◕‿-)", "emotions"},
		{"cool", "(⌐■_■)", "emotions"},
		{"confused", "(◎_◎;)", "emotions"},
		{"shocked", "(⊙_⊙)", "emotions"},
		{"sleepy", "(︶.︶✽)", "emotions"},
		{"excited", "(ノ◕ヮ◕)ノ*:・゚✧", "emotions"},
		{"nervous", "(°△°|||)", "emotions"},
		{"dead", "(✖╭╮✖)", "emotions"},
		{"hug", "(づ￣ ³￣)づ", "actions"},
		{"kiss", "(づ ̄ ³ ̄)づ", "actions"},
		{"dance", "♪┏(・o･)┛♪", "actions"},
		{"run", "ε=ε=ε=┌(;*´Д`)ノ", "actions"},
		{"fight", "(ง •̀_•́)ง", "actions"},
		{"flex", "ᕦ(ò_óˇ)ᕤ", "actions"},
		{"wave", "(^_^)/", "actions"},
		{"facepalm", "(－‸ლ)", "expressions"},
		{"thinking", "(￢_￢)", "expressions"},
		{"whatever", "┐(´д`)┌", "expressions"},
		{"sparkle", "(ﾉ◕ヮ◕)ﾉ*:・゚✧", "magic"},
		{"magic", "(ノ◕ヮ◕)ノ*:・゚✧", "magic"},

		// Animals
		{"bear", "ʕ•ᴥ•ʔ", "animals"},
		{"cat", "(=^･ω･^=)", "animals"},
		{"dog", "(◕ᴥ◕)", "animals"},
		{"bunny", "(='.'=)", "animals"},
		{"fish", "><>", "animals"},
		{"bird", "(•ө•)", "animals"},
		{"pig", "(´・ω・`)", "animals"},
		{"spider", "/╲/\\╭(ఠఠ益ఠఠ)╮/\\╱\\", "animals"},
		{"owl", "(ᵔᴥᵔ)", "animals"},
		{"octopus", "(°□°)╯︵ 🐙", "animals"},
		{"crab", "(V)_(°,,,,°)_(V)", "animals"},
		{"butterfly", "ƸӜƷ", "animals"},
		{"mouse", "<:3 )~~~", "animals"},
		{"whale", "~≋≋≋≋(◕⌓◕)≋≋≋≋~", "animals"},
		{"penguin", "(・Θ・)", "animals"},
		{"koala", "ʕ·͡ᴥ·ʔ", "animals"},
		{"panda", "ʕ•̀ω•́ʔ✧", "animals"},
		{"fox", "(^・ω・^)", "animals"},
		{"elephant", "~~_(>_<)_/", "animals"},
		{"snail", "@('.')@", "animals"},
		{"bat", "/|\\(◉◉)/|\\", "animals"},

		// Enhanced Unicode Art
		{"fancy-shrug", "乁༼☯‿☯✿༽ㄏ", "expressions"},
		{"wizard-cat", "(=^･ω･^=)⊃━☆ﾟ.*･｡ﾟ", "magic"},
		{"hadouken", "(つ◕౪◕)つ━☆ﾟ.*･｡ﾟ", "magic"},
		{"kamehameha", "╰(ᵕ◕ᵕ)━☆ﾟ.*・。ﾟ", "magic"},
		{"king-bear", "ʕ♔ᴥ♔ʔ", "animals"},
		{"queen-bear", "ʕ♕ᴥ♕ʔ", "animals"},
		{"mega-sparkle", "✧･ﾟ: *✧･ﾟ:* *:･ﾟ✧*:･ﾟ✧", "decorative"},
		{"shooting", "(☞ﾟヮﾟ)☞ ☆ﾟ.*･｡", "actions"},
		{"thug-life", "( •_•)>⌐■-■ (⌐■_■)", "expressions"},
		{"deal-with-it", "(⌐■_■)", "expressions"},
		{"gamer", "ᕕ( ᐛ )ᕗ🎮", "gaming"},
		{"rage-quit", "(ノಠ益ಠ)ノ彡🎮", "gaming"},
		{"coffee", "(っ˘ڡ˘ς)☕", "food"},
		{"cheers", "(^_^)っ🍺🍻🍺", "food"},
		{"party", "ヽ(>∀<☆)ノ🎉", "celebrations"},
		{"celebrate", "☆*:.｡.o(≧▽≦)o.｡.:*☆", "celebrations"},
		{"fireworks", "✧･ﾟ:*🎆*:･ﾟ✧", "celebrations"},
		{"crown", "♔.•*¨*•.¸¸♕", "decorative"},
		{"sword", "(ง'̀-'́)ง⚔️", "weapons"},
		{"shield", "🛡️(•̀ᴗ•́)و", "weapons"},
		{"wizard", "(∩｀-´)⊃━━☆ﾟ.*･｡ﾟ", "magic"},
		{"ninja", "(⌐▀͡ ̯ʖ▀)︻デ═一", "gaming"},
		{"robot", "[•̀ᴗ•́]🤖", "tech"},
		{"alien", "👽(∩｀-´)⊃━━☆ﾟ", "tech"},
		{"ghost", "👻(´；ω；`)", "spooky"},
		{"skull", "☠️(◣_◢)☠️", "spooky"},
		{"fire", "(◣_◢)🔥🔥🔥", "elements"},
		{"lightning", "⚡(ᗒᗩᗕ)⚡", "elements"},
		{"rainbow", "🌈(◕‿◕)🌈", "elements"},
		{"star-power", "★(≧◡≦)★", "decorative"},
		{"moon", "🌙(￣ρ￣)..zzZZ", "elements"},
		{"sun", "☀️(◕‿◕)☀️", "elements"},
		{"heart-eyes", "(♡ω♡)", "emotions"},
		{"broken-heart", "(◞‸◟；)💔", "emotions"},
		{"music", "♪♫(◕‿◕)♫♪", "music"},
		{"notes", "♪(๑ᴖ◡ᴖ๑)♪", "music"},
		{"success", "✓(≧∇≦)✓", "status"},
		{"fail", "✗(╥﹏╥)✗", "status"},
		{"warning", "⚠️(◎_◎;)⚠️", "status"},
		{"error", "❌(×_×;)❌", "status"},
		{"loading", "◌◌◌(°ロ°)◌◌◌", "tech"},
		{"send", "(ノ°▽°)ノ📨", "tech"},
		{"receive", "📬(◕‿◕)", "tech"},

		// More expressions
		{"blush", "(⁄ ⁄>⁄ ▽ ⁄<⁄ ⁄)", "emotions"},
		{"drool", "(๑´ڡ`๑)", "emotions"},
		{"eyeroll", "(◔_◔)", "expressions"},
		{"smirk", "(¬‿¬)", "expressions"},
		{"crazy", "(⊙_☉)", "expressions"},
		{"derp", "(◐ω◑)", "expressions"},
		{"suspicious", "(¬_¬)", "expressions"},
		{"innocent", "(◕ᴗ◕✿)", "emotions"},
		{"evil", "(◣∀◢)ψ", "expressions"},
		{"devil", "ψ(｀∇´)ψ", "expressions"},
		{"angel", "☆ﾐ(o*･ω･)ﾉ", "expressions"},
		{"pray", "(人´∀`)", "actions"},
		{"bow", "m(_ _)m", "actions"},
		{"clap", "(ノ´ヮ´)ノ*:・゚✧", "actions"},
		{"thumbsup", "(๑•̀ㅂ•́)و✧", "actions"},
		{"highfive", "(っ^▿^)っ", "actions"},
		{"punch", "(ノ-_-)ノ ⌒ ●~*", "actions"},
		{"kick", "(ノ￣皿￣)ノ ⌒== ┫", "actions"},
		{"stare", "(• ε •)", "expressions"},
		{"glare", "(╬ಠ益ಠ)", "expressions"},
		{"peek", "(눈_눈)", "expressions"},
		{"hide", "┬┴┬┴┤(･_├┬┴┬┴", "actions"},
		{"zombie", "[¬º-°]¬", "spooky"},
		{"vampire", "(◕ᗝ◕)", "spooky"},
		{"nosebleed", "(≧◡≦)♡* :", "emotions"},

		// Text decorations
		{"stars", "★☆★☆★", "decorative"},
		{"hearts", "♥♡♥♡♥", "decorative"},
		{"flowers", "✿❀✿❀✿", "decorative"},
		{"diamonds", "◇◆◇◆◇", "decorative"},
		{"arrows", "➳➳➳", "decorative"},
		{"divider", "═══════════", "decorative"},
		{"wave-line", "〰〰〰〰〰", "decorative"},
		{"dots", "•••••", "decorative"},
		{"sparkles-line", "✧･ﾟ:*✧･ﾟ:*", "decorative"},
	}

	// Build map for quick lookup
	kaomojiMap = make(map[string]string)
	for _, k := range kaomojis {
		kaomojiMap[k.Name] = k.Kaomoji
	}

	smileys = map[string]string{
		":)":  "😊",
		":-)": "😊",
		":(":  "😢",
		":-(": "😢",
		":D":  "😃",
		":-D": "😃",
		";)":  "😉",
		";-)": "😉",
		":P":  "😛",
		":-P": "😛",
		":p":  "😛",
		":-p": "😛",
		"<3":  "❤️",
		":O":  "😮",
		":-O": "😮",
		":o":  "😮",
		":-o": "😮",
		"XD":  "😆",
		"xD":  "😆",
		":/":  "😕",
		":-/": "😕",
		":|":  "😐",
		":-|": "😐",
		">:(": "😠",
		":*":  "😘",
		":-*": "😘",
		"B)":  "😎",
		"B-)": "😎",
		":'(": "😢",
		":')": "😂",
		"^_^": "😊",
		"-_-": "😑",
		"o_o": "😳",
		"O_O": "😳",
		">_<": "😣",
		"T_T": "😭",
		":3":  "😺",
		"=)":  "😊",
		"=D":  "😃",
		"D:":  "😧",
		">:)": "😈",
		"0:)": "😇",
	}

	// Initialize ASCII art
	asciiArt = map[string]ArtEntry{
		"skull": {`
    ▄▀▀▀▀▀▀▀▀▀▀▄
   █           █
  █  ▀▄   ▄▀   █
  █   ▀▀▀▀▀    █
  █   ▀▀ ▀▀    █
   █           █
    ▀▄▄▄▄▄▄▄▄▄▀
       █ █
       █ █
      █   █
`, "spooky"},
		"heart": {`
  ████████   ████████
 ██████████ ██████████
█████████████████████████
█████████████████████████
 ███████████████████████
  █████████████████████
   ███████████████████
    █████████████████
     ███████████████
      █████████████
       ███████████
        █████████
         ███████
          █████
           ███
            █
`, "love"},
		"star": {`
        ★
       /|\
      / | \
     /  |  \
    /   |   \
   /____|____\
       /\
      /  \
     /    \
    /      \
`, "decorative"},
		"rocket": {`
        /\
       /  \
      /    \
     |      |
     |  /\  |
     | /  \ |
    /|      |\
   / |      | \
  /__|______|__\
     /      \
    /  /  \  \
   |  |    |  |
   |__|    |__|
`, "tech"},
		"computer": {`
  _______________
 |  ___________  |
 | |           | |
 | |   hello   | |
 | |   world   | |
 | |___________| |
 |_______________|
    /         \
   /___________\
  |  O      O  |
  |_____________|
`, "tech"},
		"coffee-cup": {`
       )  (
      (   ) )
       ) ( (
     _______)_
  .-'---------|
 ( C|/\/\/\/\/|
  '-./\/\/\/\/|
    '_________'
     '-------'
`, "food"},
		"music-note": {`
     ♪♪
    ♪  ♪
   ♪    ♪
  ████████
  █      █
  █      █
  █      █
  ████████
     █
     █
     █
   ████
`, "music"},
		"tree": {`
       ★
      /|\
     / | \
    /  |  \
   /___|___\
      / \
     /   \
    /     \
   /_______\
      |||
      |||
   ____|____
`, "nature"},
		"house": {`
        /\
       /  \
      /    \
     /______\
     |      |
     | []   |
     |      |
     |______|
`, "objects"},
		"cat-face": {`
   /\_____/\
  /  o   o  \
 ( ==  ^  == )
  )         (
 (           )
( (  )   (  ) )
(__(__)___(__)__)
`, "animals"},
		"dog-face": {`
    / \__
   (    @\____
   /         O
  /   (_____/
 /_____/   U
`, "animals"},
		"peace": {`
     _____
    /     \
   /   |   \
  |    |    |
  |   /|\   |
  |  / | \  |
  | /  |  \ |
   \   |   /
    \__|__/
`, "symbols"},
		"yin-yang": {`
      ___
    /     \
   / ●     \
  |   ___   |
  |  /   \  |
  |  \___/  |
   \     ● /
    \_____/
`, "symbols"},
		"diamond": {`
     /\
    /  \
   /    \
  /      \
  \      /
   \    /
    \  /
     \/
`, "decorative"},
		"crown": {`
  __  __  __
 |  ||  ||  |
 |  ||  ||  |
/    \  /    \
|    |  |    |
|____|__|____|
`, "decorative"},
		"sword": {`
       ||
       ||
       ||
       ||
   ___||||___
  |_________|
       ||
       ||
       ||
      /  \
     /    \
`, "weapons"},
		"controller": {`
  _____________
 /   _     _   \
|   (_)   (_)   |
|       _       |
| [_]  (_) [_]  |
|   _       _   |
|  (_)     (_)  |
 \_____________/
`, "gaming"},
		"brain": {
			"     ,---.\n" +
				"    / .   \\\n" +
				"   |  |   |\n" +
				"   |  |   |\n" +
				"    \\ '   /\n" +
				"     '---'\n" +
				"    /_____\\\n" +
				"   |  ___  |\n" +
				"   | |   | |\n" +
				"   | |   | |\n" +
				"   |_|   |_|\n", "tech"},
		"lightning-bolt": {`
    ██████
   ██
  ██████
      ██
  ██████
 ██
██████
`, "elements"},
		"moon-crescent": {`
      @@@@
    @@
   @
  @
  @
   @
    @@
      @@@@
`, "elements"},
		"sun-rays": {`
      \  |  /
       \ | /
    ----☀----
       / | \
      /  |  \
`, "elements"},
		"cloud": {`
       .---.
      (     )
   .-(       )-.
  (             )
   '-.       .-'
      '-----'
`, "elements"},
		"thumbs-up": {`
       _
      | |
     _| |_
    |     |
    |     |
    |     |
     \   /
      \_/
       |
       |
      /_\
`, "actions"},
	}
}

// Get returns a kaomoji by name
func Get(name string) (string, bool) {
	k, ok := kaomojiMap[strings.ToLower(name)]
	return k, ok
}

// KaomojiItem represents a kaomoji for listing
type KaomojiItem struct {
	Name     string `json:"name"`
	Kaomoji  string `json:"kaomoji"`
	Category string `json:"category,omitempty"`
}

// List returns all kaomoji, optionally filtered by search term and category
func List(search, category string) []KaomojiItem {
	var result []KaomojiItem
	searchLower := strings.ToLower(search)
	categoryLower := strings.ToLower(category)

	for _, k := range kaomojis {
		matchesSearch := search == "" || strings.Contains(strings.ToLower(k.Name), searchLower)
		matchesCategory := category == "" || strings.ToLower(k.Category) == categoryLower

		if matchesSearch && matchesCategory {
			result = append(result, KaomojiItem(k))
		}
	}
	return result
}

// Random returns a random kaomoji
func Random() (string, string) {
	idx := rand.Intn(len(kaomojis))
	return kaomojis[idx].Name, kaomojis[idx].Kaomoji
}

// SmileyToEmoji converts an ASCII smiley to emoji
func SmileyToEmoji(smiley string) (string, bool) {
	e, ok := smileys[smiley]
	return e, ok
}

// Suggest returns similar kaomoji names
func Suggest(name string) []string {
	nameLower := strings.ToLower(name)
	var suggestions []string

	for _, k := range kaomojis {
		if strings.Contains(k.Name, nameLower) || strings.Contains(nameLower, k.Name) {
			suggestions = append(suggestions, k.Name)
			if len(suggestions) >= 3 {
				break
			}
		}
	}
	return suggestions
}

// ListCategories returns all unique categories
func ListCategories() []string {
	categorySet := make(map[string]bool)
	for _, k := range kaomojis {
		categorySet[k.Category] = true
	}

	var categories []string
	for cat := range categorySet {
		categories = append(categories, cat)
	}
	sort.Strings(categories)
	return categories
}

// ArtItem represents an ASCII art piece for listing
type ArtItem struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

// ListArt returns all ASCII art names, optionally filtered by category
func ListArt(category string) []ArtItem {
	var result []ArtItem
	categoryLower := strings.ToLower(category)

	for name, entry := range asciiArt {
		if category == "" || strings.ToLower(entry.Category) == categoryLower {
			result = append(result, ArtItem{
				Name:     name,
				Category: entry.Category,
			})
		}
	}

	// Sort by name
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})

	return result
}

// GetArt returns ASCII art by name
func GetArt(name string) (string, bool) {
	entry, ok := asciiArt[strings.ToLower(name)]
	if !ok {
		return "", false
	}
	return entry.Art, true
}
