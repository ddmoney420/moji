//go:build js && wasm

package main

import (
	"encoding/json"
	"strings"
	"syscall/js"
	"time"

	"github.com/ddmoney420/moji/internal/artdb"
	"github.com/ddmoney420/moji/internal/banner"
	"github.com/ddmoney420/moji/internal/calendar"
	"github.com/ddmoney420/moji/internal/chain"
	"github.com/ddmoney420/moji/internal/effects"
	"github.com/ddmoney420/moji/internal/filters"
	"github.com/ddmoney420/moji/internal/gradient"
	"github.com/ddmoney420/moji/internal/kaomoji"
	"github.com/ddmoney420/moji/internal/patterns"
	"github.com/ddmoney420/moji/internal/qrcode"
	"github.com/ddmoney420/moji/internal/speech"
	"github.com/ddmoney420/moji/internal/styles"
)

func main() {
	// Banner
	js.Global().Set("mojiBanner", js.FuncOf(mojiBanner))
	js.Global().Set("mojiBannerFonts", js.FuncOf(mojiBannerFonts))

	// Kaomoji
	js.Global().Set("mojiKaomoji", js.FuncOf(mojiKaomojiGet))
	js.Global().Set("mojiKaomojiList", js.FuncOf(mojiKaomojiList))
	js.Global().Set("mojiKaomojiRandom", js.FuncOf(mojiKaomojiRandom))
	js.Global().Set("mojiKaomojiCategories", js.FuncOf(mojiKaomojiCategories))
	js.Global().Set("mojiKaomojiArt", js.FuncOf(mojiKaomojiArt))
	js.Global().Set("mojiKaomojiArtList", js.FuncOf(mojiKaomojiArtList))
	js.Global().Set("mojiKaomojiSuggest", js.FuncOf(mojiKaomojiSuggest))
	js.Global().Set("mojiSmileyToEmoji", js.FuncOf(mojiSmileyToEmoji))

	// Effects
	js.Global().Set("mojiEffect", js.FuncOf(mojiEffect))
	js.Global().Set("mojiEffectList", js.FuncOf(mojiEffectList))

	// Filters
	js.Global().Set("mojiFilter", js.FuncOf(mojiFilter))
	js.Global().Set("mojiFilterList", js.FuncOf(mojiFilterList))
	js.Global().Set("mojiFilterChain", js.FuncOf(mojiFilterChain))

	// Gradient
	js.Global().Set("mojiGradient", js.FuncOf(mojiGradient))
	js.Global().Set("mojiGradientThemes", js.FuncOf(mojiGradientThemes))

	// QR Code
	js.Global().Set("mojiQR", js.FuncOf(mojiQR))
	js.Global().Set("mojiQRCompact", js.FuncOf(mojiQRCompact))

	// Patterns
	js.Global().Set("mojiPattern", js.FuncOf(mojiPatternBorder))
	js.Global().Set("mojiDivider", js.FuncOf(mojiDivider))
	js.Global().Set("mojiPatternCreate", js.FuncOf(mojiPatternCreate))
	js.Global().Set("mojiPatternPreset", js.FuncOf(mojiPatternPreset))
	js.Global().Set("mojiPatternSymmetric", js.FuncOf(mojiPatternSymmetric))
	js.Global().Set("mojiPatternListBorders", js.FuncOf(mojiPatternListBorders))
	js.Global().Set("mojiPatternListDividers", js.FuncOf(mojiPatternListDividers))
	js.Global().Set("mojiPatternListPatterns", js.FuncOf(mojiPatternListPatterns))


	// Styles
	js.Global().Set("mojiStyle", js.FuncOf(mojiStyle))
	js.Global().Set("mojiStyleList", js.FuncOf(mojiStyleList))
	js.Global().Set("mojiStyleBorder", js.FuncOf(mojiStyleBorder))
	js.Global().Set("mojiStyleBorderList", js.FuncOf(mojiStyleBorderList))
	js.Global().Set("mojiStyleAlign", js.FuncOf(mojiStyleAlign))

	// Speech
	js.Global().Set("mojiSpeech", js.FuncOf(mojiSpeech))
	js.Global().Set("mojiSpeechStyles", js.FuncOf(mojiSpeechStyles))
	js.Global().Set("mojiSpeechCombine", js.FuncOf(mojiSpeechCombine))

	// Art DB
	js.Global().Set("mojiArtList", js.FuncOf(mojiArtList))
	js.Global().Set("mojiArtGet", js.FuncOf(mojiArtGet))
	js.Global().Set("mojiArtSearch", js.FuncOf(mojiArtSearch))
	js.Global().Set("mojiArtCategories", js.FuncOf(mojiArtCategories))
	js.Global().Set("mojiArtRandom", js.FuncOf(mojiArtRandom))

	// Chain
	js.Global().Set("mojiChain", js.FuncOf(mojiChain))

	// Calendar
	js.Global().Set("mojiCalendarMonth", js.FuncOf(mojiCalendarMonth))
	js.Global().Set("mojiCalendarYear", js.FuncOf(mojiCalendarYear))
	js.Global().Set("mojiCalendarCurrent", js.FuncOf(mojiCalendarCurrent))
	js.Global().Set("mojiCalendarToday", js.FuncOf(mojiCalendarToday))
	js.Global().Set("mojiCalendarArt", js.FuncOf(mojiCalendarArt))
	js.Global().Set("mojiCalendarWeek", js.FuncOf(mojiCalendarWeek))

	select {}
}

// --- Banner ---

func mojiBanner(_ js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return "error: requires text and font arguments"
	}
	result, err := banner.Generate(args[0].String(), args[1].String())
	if err != nil {
		return "error: " + err.Error()
	}
	return result
}

func mojiBannerFonts(_ js.Value, _ []js.Value) interface{} {
	fonts := banner.ListFonts()
	type item struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	items := make([]item, len(fonts))
	for i, f := range fonts {
		items[i] = item{Name: f.Name, Desc: f.Desc}
	}
	data, _ := json.Marshal(items)
	return string(data)
}

// --- Kaomoji ---

func mojiKaomojiGet(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return ""
	}
	result, ok := kaomoji.Get(args[0].String())
	if !ok {
		return ""
	}
	return result
}

func mojiKaomojiList(_ js.Value, args []js.Value) interface{} {
	search, category := "", ""
	if len(args) > 0 {
		search = args[0].String()
	}
	if len(args) > 1 {
		category = args[1].String()
	}
	data, _ := json.Marshal(kaomoji.List(search, category))
	return string(data)
}

func mojiKaomojiRandom(_ js.Value, _ []js.Value) interface{} {
	name, kao := kaomoji.Random()
	data, _ := json.Marshal(map[string]string{"name": name, "kaomoji": kao})
	return string(data)
}

func mojiKaomojiCategories(_ js.Value, _ []js.Value) interface{} {
	data, _ := json.Marshal(kaomoji.ListCategories())
	return string(data)
}

func mojiKaomojiArt(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return ""
	}
	result, ok := kaomoji.GetArt(args[0].String())
	if !ok {
		return ""
	}
	return strings.Trim(result, "\n")
}

func mojiKaomojiArtList(_ js.Value, args []js.Value) interface{} {
	category := ""
	if len(args) > 0 {
		category = args[0].String()
	}
	data, _ := json.Marshal(kaomoji.ListArt(category))
	return string(data)
}

func mojiKaomojiSuggest(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "[]"
	}
	data, _ := json.Marshal(kaomoji.Suggest(args[0].String()))
	return string(data)
}

func mojiSmileyToEmoji(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return ""
	}
	result, ok := kaomoji.SmileyToEmoji(args[0].String())
	if !ok {
		return ""
	}
	return result
}

// --- Effects ---

func mojiEffect(_ js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return ""
	}
	return effects.Apply(args[0].String(), args[1].String())
}

func mojiEffectList(_ js.Value, _ []js.Value) interface{} {
	list := effects.ListEffects()
	type item struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	items := make([]item, len(list))
	for i, e := range list {
		items[i] = item{Name: e.Name, Desc: e.Desc}
	}
	data, _ := json.Marshal(items)
	return string(data)
}

// --- Filters ---

func mojiFilter(_ js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return ""
	}
	f, ok := filters.Get(args[0].String())
	if !ok {
		return "error: unknown filter: " + args[0].String()
	}
	return f(args[1].String())
}

func mojiFilterList(_ js.Value, _ []js.Value) interface{} {
	list := filters.ListFilters()
	type item struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	items := make([]item, len(list))
	for i, e := range list {
		items[i] = item{Name: e.Name, Desc: e.Desc}
	}
	data, _ := json.Marshal(items)
	return string(data)
}

func mojiFilterChain(_ js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return ""
	}
	text := args[0].String()
	chainSpec := args[1].String()
	filterNames := filters.ParseChain(chainSpec)
	return filters.Chain(text, filterNames)
}

// --- Gradient ---

func mojiGradient(_ js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return ""
	}
	return gradient.Apply(args[0].String(), args[1].String(), args[2].String())
}

func mojiGradientThemes(_ js.Value, _ []js.Value) interface{} {
	list := gradient.ListThemes()
	type item struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	items := make([]item, len(list))
	for i, e := range list {
		items[i] = item{Name: e.Name, Desc: e.Desc}
	}
	data, _ := json.Marshal(items)
	return string(data)
}

// --- QR Code ---

func mojiQR(_ js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return ""
	}
	text := args[0].String()
	charset := args[1].String()
	invert := false
	if len(args) > 2 {
		invert = args[2].Bool()
	}
	opts := qrcode.Options{
		Charset: charset,
		Invert:  invert,
		ANSI:    true,
	}
	result, err := qrcode.Generate(text, opts)
	if err != nil {
		return "error: " + err.Error()
	}
	return result
}

func mojiQRCompact(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return ""
	}
	text := args[0].String()
	invert := false
	if len(args) > 1 {
		invert = args[1].Bool()
	}
	result, err := qrcode.GenerateCompact(text, invert)
	if err != nil {
		return "error: " + err.Error()
	}
	return result
}

// --- Patterns ---

func mojiPatternBorder(_ js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return ""
	}
	padding := 1
	if len(args) > 2 {
		padding = args[2].Int()
	}
	return patterns.CreateBorder(args[0].String(), args[1].String(), padding)
}

func mojiDivider(_ js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return ""
	}
	return patterns.CreateDivider(args[0].String(), args[1].Int())
}

func mojiPatternCreate(_ js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return ""
	}
	return patterns.CreatePattern(args[0].String(), args[1].Int(), args[2].Int())
}

func mojiPatternPreset(_ js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return ""
	}
	return patterns.GetPreset(args[0].String(), args[1].Int(), args[2].Int())
}

func mojiPatternSymmetric(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return ""
	}
	return patterns.CreateSymmetric(args[0].String())
}

func mojiPatternListBorders(_ js.Value, _ []js.Value) interface{} {
	data, _ := json.Marshal(patterns.ListBorders())
	return string(data)
}

func mojiPatternListDividers(_ js.Value, _ []js.Value) interface{} {
	data, _ := json.Marshal(patterns.ListDividers())
	return string(data)
}

func mojiPatternListPatterns(_ js.Value, _ []js.Value) interface{} {
	data, _ := json.Marshal(patterns.ListPatterns())
	return string(data)
}


// --- Styles ---

func mojiStyle(_ js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return ""
	}
	return styles.Apply(args[0].String(), args[1].String())
}

func mojiStyleList(_ js.Value, _ []js.Value) interface{} {
	list := styles.ListStyles()
	type item struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	items := make([]item, len(list))
	for i, s := range list {
		items[i] = item{Name: s.Name, Desc: s.Desc}
	}
	data, _ := json.Marshal(items)
	return string(data)
}

func mojiStyleBorder(_ js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return ""
	}
	return styles.ApplyBorder(args[0].String(), args[1].String())
}

func mojiStyleBorderList(_ js.Value, _ []js.Value) interface{} {
	list := styles.ListBorders()
	type item struct {
		Name string `json:"name"`
		Desc string `json:"desc"`
	}
	items := make([]item, len(list))
	for i, s := range list {
		items[i] = item{Name: s.Name, Desc: s.Desc}
	}
	data, _ := json.Marshal(items)
	return string(data)
}

func mojiStyleAlign(_ js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return ""
	}
	return styles.ApplyAlignment(args[0].String(), args[1].String(), args[2].Int())
}

// --- Speech ---

func mojiSpeech(_ js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		return ""
	}
	return speech.Wrap(args[0].String(), args[1].String(), args[2].Int())
}

func mojiSpeechStyles(_ js.Value, _ []js.Value) interface{} {
	data, _ := json.Marshal(speech.ListStyles())
	return string(data)
}

func mojiSpeechCombine(_ js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return ""
	}
	return speech.Combine(args[0].String(), args[1].String())
}

// --- Art DB ---

func mojiArtList(_ js.Value, args []js.Value) interface{} {
	category := ""
	if len(args) > 0 {
		category = args[0].String()
	}
	type item struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	var arts []artdb.Art
	if category != "" {
		arts = artdb.ByCategory(category)
	} else {
		arts = artdb.List()
	}
	items := make([]item, len(arts))
	for i, a := range arts {
		items[i] = item{Name: a.Name, Category: a.Category}
	}
	data, _ := json.Marshal(items)
	return string(data)
}

func mojiArtGet(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return ""
	}
	art, ok := artdb.Get(args[0].String())
	if !ok {
		return ""
	}
	return strings.Trim(art.Art, "\n")
}

func mojiArtSearch(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return "[]"
	}
	arts := artdb.Search(args[0].String())
	type item struct {
		Name     string `json:"name"`
		Category string `json:"category"`
	}
	items := make([]item, len(arts))
	for i, a := range arts {
		items[i] = item{Name: a.Name, Category: a.Category}
	}
	data, _ := json.Marshal(items)
	return string(data)
}

func mojiArtCategories(_ js.Value, _ []js.Value) interface{} {
	data, _ := json.Marshal(artdb.ListCategories())
	return string(data)
}

func mojiArtRandom(_ js.Value, _ []js.Value) interface{} {
	art := artdb.Random()
	type item struct {
		Name     string `json:"name"`
		Category string `json:"category"`
		Art      string `json:"art"`
	}
	data, _ := json.Marshal(item{Name: art.Name, Category: art.Category, Art: strings.Trim(art.Art, "\n")})
	return string(data)
}

// --- Chain ---

func mojiChain(_ js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return ""
	}
	text := args[0].String()
	opts := chain.Options{}
	if len(args) > 1 {
		optsJSON := args[1].String()
		var raw struct {
			Gradient     string `json:"gradient"`
			GradientMode string `json:"gradientMode"`
			Border       string `json:"border"`
			BorderPad    int    `json:"borderPad"`
			Effect       string `json:"effect"`
			Bubble       string `json:"bubble"`
			BubbleWidth  int    `json:"bubbleWidth"`
		}
		if err := json.Unmarshal([]byte(optsJSON), &raw); err == nil {
			opts.Gradient = raw.Gradient
			opts.GradientMode = raw.GradientMode
			opts.Border = raw.Border
			opts.BorderPad = raw.BorderPad
			opts.Effect = raw.Effect
			opts.Bubble = raw.Bubble
			opts.BubbleWidth = raw.BubbleWidth
		}
	}
	return chain.Apply(text, opts)
}

// --- Calendar ---

func mojiCalendarMonth(_ js.Value, args []js.Value) interface{} {
	year := time.Now().Year()
	month := time.Now().Month()
	if len(args) > 0 {
		year = args[0].Int()
	}
	if len(args) > 1 {
		month = time.Month(args[1].Int())
	}
	opts := calendar.Options{}
	if len(args) > 2 && args[2].Bool() {
		opts.FirstDayMon = true
	}
	return calendar.Month(year, month, opts)
}

func mojiCalendarYear(_ js.Value, args []js.Value) interface{} {
	year := time.Now().Year()
	if len(args) > 0 {
		year = args[0].Int()
	}
	return calendar.Year(year, calendar.Options{})
}

func mojiCalendarCurrent(_ js.Value, _ []js.Value) interface{} {
	return calendar.Current(calendar.Options{})
}

func mojiCalendarToday(_ js.Value, _ []js.Value) interface{} {
	return calendar.Today()
}

func mojiCalendarArt(_ js.Value, _ []js.Value) interface{} {
	return calendar.ASCIIArt()
}

func mojiCalendarWeek(_ js.Value, _ []js.Value) interface{} {
	return calendar.WeekView(calendar.Options{})
}

