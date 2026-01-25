package banner

import (
	"embed"
	"fmt"
	"strings"

	"github.com/ddmoney420/moji/internal/figlet"
)

//go:embed fonts/*.flf
var fontsFS embed.FS

var fontCache = make(map[string]*figlet.Font)

// FontInfo contains font information
type FontInfo struct {
	Name string
	Desc string
}

// ListFonts returns available fonts with descriptions
func ListFonts() []FontInfo {
	return []FontInfo{
		// 3D Styles
		{"3d", "3D block letters"},
		{"3d-ascii", "3D ASCII art"},
		{"isometric1", "3D isometric view"},
		{"isometric2", "3D isometric alt"},
		{"isometric3", "3D isometric blocks"},
		{"henry-3d", "Henry 3D style"},
		{"larry-3d", "Larry 3D effect"},
		{"impossible", "Impossible 3D"},

		// Graffiti/Street
		{"graffiti", "Street art graffiti"},
		{"bloody", "Dripping blood horror"},
		{"poison", "Toxic dripping"},
		{"ghost", "Spooky ghost style"},
		{"swamp-land", "Swampy dripping"},

		// Retro/Gaming
		{"doom", "DOOM game style"},
		{"dos-rebel", "DOS Rebel retro"},
		{"star-wars", "Star Wars crawl"},
		{"electronic", "Electronic/digital"},
		{"speed", "Speed racing"},
		{"sub-zero", "Sub-Zero frozen"},

		// Big & Bold
		{"colossal", "Extra large letters"},
		{"epic", "Epic massive"},
		{"doh", "Homer Simpson DOH!"},
		{"big", "Large block letters"},
		{"banner3-d", "3D banner style"},
		{"bulbhead", "Bulbous letters"},

		// Clean/Classic
		{"standard", "Default FIGlet"},
		{"slant", "Italic/slanted"},
		{"shadow", "Letters with shadow"},
		{"ansi-shadow", "ANSI shadow style"},
		{"ansi-regular", "ANSI regular"},
		{"roman", "Roman style"},
		{"lean", "Lean style"},
		{"univers", "Universal style"},

		// Decorative
		{"fraktur", "Gothic fraktur"},
		{"alligator", "Alligator scales"},
		{"arrows", "Made of arrows"},
		{"blocks", "Block elements"},
		{"rectangles", "Rectangle shapes"},
		{"varsity", "Sports varsity"},
		{"crawford", "Crawford style"},
		{"delta-corps", "Delta Corps military"},

		// Fun/Quirky
		{"weird", "Weird distorted"},
		{"twisted", "Twisted letters"},
		{"crazy", "Crazy style"},
		{"whimsy", "Whimsical"},
		{"puffy", "Puffy clouds"},
		{"chunky", "Chunky blocks"},
		{"pebbles", "Pebble texture"},
		{"tinker-toy", "Tinker toy style"},
		{"rozzo", "Rozzo style"},
		{"fire", "Fire/flames"},
		{"cyberlarge", "Cyberpunk large"},

		// Small/Compact
		{"elite", "Elite/1337 style"},
		{"ogre", "Ogre style"},
		{"peaks", "Mountain peaks"},
		{"nancyj", "Nancy J style"},
		{"calvin", "Calvin style"},
	}
}

// Generate creates ASCII art banner using embedded fonts
func Generate(text, fontName string) (string, error) {
	fontFile := mapFontName(fontName)
	if fontFile == "" {
		fontFile = "Standard.flf" // default
	}

	// Check cache first
	if font, ok := fontCache[fontFile]; ok {
		return font.Render(text), nil
	}

	// Load font from embedded FS
	data, err := fontsFS.ReadFile("fonts/" + fontFile)
	if err != nil {
		// Try lowercase
		data, err = fontsFS.ReadFile("fonts/" + strings.ToLower(fontFile))
		if err != nil {
			return "", fmt.Errorf("font '%s' not found: %v", fontName, err)
		}
	}

	font, err := figlet.ParseFont(string(data))
	if err != nil {
		return "", fmt.Errorf("failed to parse font '%s': %v", fontName, err)
	}

	// Cache for future use
	fontCache[fontFile] = font

	return font.Render(text), nil
}

// mapFontName maps user-friendly names to font files
func mapFontName(name string) string {
	mapping := map[string]string{
		// 3D Styles
		"3d":         "3d.flf",
		"3d-ascii":   "3D-ASCII.flf",
		"isometric":  "Isometric1.flf",
		"isometric1": "Isometric1.flf",
		"isometric2": "Isometric2.flf",
		"isometric3": "Isometric3.flf",
		"henry-3d":   "Henry-3D.flf",
		"henry3d":    "Henry-3D.flf",
		"larry-3d":   "Larry-3D.flf",
		"larry3d":    "larry3d.flf",
		"impossible": "Impossible.flf",

		// Graffiti/Street
		"graffiti":   "Graffiti.flf",
		"bloody":     "bloody.flf",
		"poison":     "Poison.flf",
		"ghost":      "Ghost.flf",
		"swamp-land": "Swamp-Land.flf",
		"swamp":      "Swamp-Land.flf",

		// Retro/Gaming
		"doom":       "doom.flf",
		"dos-rebel":  "DOS-Rebel.flf",
		"dosrebel":   "DOS-Rebel.flf",
		"star-wars":  "Star-Wars.flf",
		"starwars":   "Star-Wars.flf",
		"electronic": "Electronic.flf",
		"speed":      "Speed.flf",
		"sub-zero":   "Sub-Zero.flf",
		"subzero":    "Sub-Zero.flf",

		// Big & Bold
		"colossal":  "colossal.flf",
		"epic":      "Epic.flf",
		"doh":       "Doh.flf",
		"big":       "big.flf",
		"banner3-d": "Banner3-D.flf",
		"banner3d":  "Banner3-D.flf",
		"bulbhead":  "Bulbhead.flf",

		// Clean/Classic
		"standard":     "Standard.flf",
		"slant":        "slant.flf",
		"shadow":       "shadow.flf",
		"ansi-shadow":  "ANSI-Shadow.flf",
		"ansishadow":   "ANSI-Shadow.flf",
		"ansi-regular": "ANSI-Regular.flf",
		"ansi":         "ANSI-Regular.flf",
		"roman":        "Roman.flf",
		"lean":         "Lean.flf",
		"univers":      "Univers.flf",

		// Decorative
		"fraktur":     "Fraktur.flf",
		"alligator":   "Alligator.flf",
		"arrows":      "Arrows.flf",
		"blocks":      "Blocks.flf",
		"rectangles":  "Rectangles.flf",
		"varsity":     "Varsity.flf",
		"crawford":    "Crawford.flf",
		"delta-corps": "Delta-Corps-Priest-1.flf",
		"deltacorps":  "Delta-Corps-Priest-1.flf",

		// Fun/Quirky
		"weird":      "Weird.flf",
		"twisted":    "Twisted.flf",
		"crazy":      "Crazy.flf",
		"whimsy":     "Whimsy.flf",
		"puffy":      "Puffy.flf",
		"chunky":     "Chunky.flf",
		"pebbles":    "Pebbles.flf",
		"tinker-toy": "Tinker-Toy.flf",
		"tinkertoy":  "Tinker-Toy.flf",
		"rozzo":      "Rozzo.flf",
		"fire":       "Fire-Font-s.flf",
		"cyberlarge": "Cyberlarge.flf",
		"cyber":      "cyber.flf",

		// Small/Compact
		"elite":  "elite.flf",
		"ogre":   "Ogre.flf",
		"peaks":  "Peaks.flf",
		"nancyj": "Nancyj.flf",
		"calvin": "calvin.flf",
	}

	if file, ok := mapping[strings.ToLower(name)]; ok {
		return file
	}
	return ""
}
