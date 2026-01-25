package artdb

import (
	"math/rand/v2"
	"strings"
)

// Art represents a piece of ASCII art
type Art struct {
	Name     string
	Category string
	Artist   string
	Art      string
	Tags     []string
}

// Database of ASCII art
var database = []Art{
	// Animals
	{
		Name:     "cat",
		Category: "animals",
		Artist:   "Unknown",
		Tags:     []string{"pet", "feline"},
		Art: `
    /\_____/\
   /  o   o  \
  ( ==  ^  == )
   )         (
  (           )
 ( (  )   (  ) )
(__(__)___(__)__)
`,
	},
	{
		Name:     "dog",
		Category: "animals",
		Artist:   "Unknown",
		Tags:     []string{"pet", "canine"},
		Art: `
   /^-----^\
   V  o o  V
    |  Y  |
    |  ^  |
     \ Q /
     / - \
    |  |  |
    |  |  |
    (__|__)
`,
	},
	{
		Name:     "bunny",
		Category: "animals",
		Artist:   "Unknown",
		Tags:     []string{"pet", "rabbit"},
		Art: `
   () ()
   (. .)
  ( >< )
   / | \
  (_/ \_)
`,
	},
	{
		Name:     "fish",
		Category: "animals",
		Artist:   "Unknown",
		Tags:     []string{"aquatic", "sea"},
		Art: `
     /\
    /  \  /\
   / /\ \/  \
  / /  \ \   \
 / /    \/    \
 \/     /\    /
  \    /  \  /
   \  / /\ \/
    \/\/  \/
`,
	},
	{
		Name:     "bird",
		Category: "animals",
		Artist:   "Unknown",
		Tags:     []string{"flying"},
		Art: `
  .---.
 / o o \
 \  ^  /
  '---'
  /| |\
 (_| |_)
`,
	},
	{
		Name:     "owl",
		Category: "animals",
		Artist:   "Unknown",
		Tags:     []string{"bird", "night"},
		Art: `
   .-.
  {O,O}
  /)_)
 / " " \
| /   \ |
 \|   |/
  '---'
`,
	},
	{
		Name:     "spider",
		Category: "animals",
		Artist:   "Unknown",
		Tags:     []string{"insect", "bug"},
		Art: `
 \  . .  /
  \_| |_/
 (( o.o ))
  /\   /\
 / |   | \
(_/     \_)
`,
	},
	{
		Name:     "snail",
		Category: "animals",
		Artist:   "Unknown",
		Tags:     []string{"slow"},
		Art: `
      @     @
       \   /
     ___\_/___
   .'         '.
  /  .--------.  \
 |  /          \  |
  \ '--------'  /
   '._________.'
`,
	},

	// Symbols & Objects
	{
		Name:     "heart",
		Category: "symbols",
		Artist:   "Unknown",
		Tags:     []string{"love", "romance"},
		Art: `
    .---.   .---.
   /     \_/     \
  |               |
  |               |
   \             /
    \           /
     \         /
      \       /
       \_____/
`,
	},
	{
		Name:     "star",
		Category: "symbols",
		Artist:   "Unknown",
		Tags:     []string{"space", "night"},
		Art: `
         .
        /|\
       / | \
    ,-'  |  '-,
  <'---=\+/=---'>
    '-,  |  ,-'
       \ | /
        \|/
         '
`,
	},
	{
		Name:     "rose",
		Category: "nature",
		Artist:   "Unknown",
		Tags:     []string{"flower", "romantic"},
		Art: `
       _
      (_)
     /| |\
    / | | \
   {  | |  }
    \ | | /
     '-|-'
       |
      -|-
     '-|-'
       |
      -'-
`,
	},
	{
		Name:     "coffee",
		Category: "objects",
		Artist:   "Unknown",
		Tags:     []string{"drink", "morning"},
		Art: `
      (  )   (   )  )
       ) (   )  ( (
      (   ) (    ) )
    .-.__________.-.
   /  .-----------.\
  |  / .---. .---. \|
  | | |     |     | |
  | | '-----'-----' |
   \ \             //
    '-.___________.-'
`,
	},
	{
		Name:     "computer",
		Category: "objects",
		Artist:   "Unknown",
		Tags:     []string{"tech", "work"},
		Art: `
 ____________________
|  ________________  |
| |                | |
| |  C:\>_         | |
| |                | |
| |________________| |
|____________________|
   _[____________]_
  /  ____________  \
 /__________________\
`,
	},
	{
		Name:     "house",
		Category: "objects",
		Artist:   "Unknown",
		Tags:     []string{"home", "building"},
		Art: `
          /\
         /  \
        /    \
       / /\ / \
      / /  Y   \
     / /  / \   \
    /_/  /   \___\
    |   |  O  |  |
    |   |_____|  |
    |___|     |__|
`,
	},
	{
		Name:     "rocket",
		Category: "objects",
		Artist:   "Unknown",
		Tags:     []string{"space", "travel"},
		Art: `
        /\
       /  \
      | /\ |
      |/  \|
      |    |
      |    |
     /| /\ |\
    / |/  \| \
   /  | /\ |  \
   \__|/  \|__/
       \  /
    ----\/----
`,
	},
	{
		Name:     "skull",
		Category: "symbols",
		Artist:   "Unknown",
		Tags:     []string{"death", "halloween"},
		Art: `
    _.--""--._
   /  _    _  \
  | (_\  /_)  |
  |   /  \    |
   \_| -- |__/
    /|_/  \_|\
   / |      | \
    "-"    "-"
`,
	},
	{
		Name:     "music",
		Category: "symbols",
		Artist:   "Unknown",
		Tags:     []string{"notes", "sound"},
		Art: `
  o       o
  |       |
  |   o   |
  |   |   |
 /|   |   |\
| |   |   | |
 \|___|___|/
`,
	},
	{
		Name:     "tree",
		Category: "nature",
		Artist:   "Unknown",
		Tags:     []string{"forest", "plant"},
		Art: `
        ###
       #o###
     #####o###
    ##o#\#|#/###
     ####\|/#o##
      ## }|{  ##
         }|{
        /   \
       /     \
   ___/       \___
`,
	},
	{
		Name:     "sun",
		Category: "nature",
		Artist:   "Unknown",
		Tags:     []string{"sky", "day"},
		Art: `
       \   |   /
        \  |  /
    --- @@@@@@@ ---
       @@@@@@@@@
   --- @@@@@@@@@@@ ---
       @@@@@@@@@
    --- @@@@@@@ ---
        /  |  \
       /   |   \
`,
	},
	{
		Name:     "moon",
		Category: "nature",
		Artist:   "Unknown",
		Tags:     []string{"sky", "night"},
		Art: `
      _..._
    .'     '.
   /    .-""-\
  |   /:. .  |
  |  |::.    |
  |  |::.    |
   \ \::. ./ |
    '.'::..'.'
      '-..-'
`,
	},
	{
		Name:     "cloud",
		Category: "nature",
		Artist:   "Unknown",
		Tags:     []string{"sky", "weather"},
		Art: `
              .-~~~-.
      .- ~ ~-(       )_ _
     /                    ~ -.
    |                         \
     \                      .'
      ~- . _____________ . -~
`,
	},

	// Faces & People
	{
		Name:     "smiley",
		Category: "faces",
		Artist:   "Unknown",
		Tags:     []string{"happy", "emoji"},
		Art: `
   .-----.
  /       \
 |  O   O  |
 |    <    |
 |  \___/  |
  \       /
   '-----'
`,
	},
	{
		Name:     "wizard",
		Category: "people",
		Artist:   "Unknown",
		Tags:     []string{"magic", "fantasy"},
		Art: `
      /\
     /  \
    / ^^ \
   /  ()  \
   |  ||  |
  /|      |\
 / | /--\ | \
(  |/    \|  )
   |      |
  _|______|_
`,
	},
	{
		Name:     "robot",
		Category: "objects",
		Artist:   "Unknown",
		Tags:     []string{"tech", "machine"},
		Art: `
    .---.
   |[o o]|
   |[::/]|
    | : |
   /|   |\
  (_|   |_)
    |   |
   /|   |\
  (_|___|_)
`,
	},
	{
		Name:     "dragon",
		Category: "animals",
		Artist:   "Unknown",
		Tags:     []string{"fantasy", "fire"},
		Art: `
                 __        .-'  _  '.
        \       /  \      / /'-' '-\ \
   --()--\-----/----\----/ /        \ \---
          '.       .' \ \ \        / /
            '-...-'    \ \ '-.  .-' /
                        \_\   ''  /_/
`,
	},
	{
		Name:     "anchor",
		Category: "symbols",
		Artist:   "Unknown",
		Tags:     []string{"nautical", "sea"},
		Art: `
      .-.
     (   )
      '-'
       |
    ---|---
       |
     __|__
    /     \
   |       |
    \_____/
       |
       V
`,
	},

	// Text Art
	{
		Name:     "welcome",
		Category: "text",
		Artist:   "Unknown",
		Tags:     []string{"greeting"},
		Art: `
 __      __       _
 \ \    / /      | |
  \ \  / /___  __| | ___ ___  _ __ ___   ___
   \ \/ // _ \/ _' |/ __/ _ \| '_ ' _ \ / _ \
    \  /|  __/ (_| | (_| (_) | | | | | |  __/
     \/  \___|\__,_|\___\___/|_| |_| |_|\___|
`,
	},
	{
		Name:     "love",
		Category: "text",
		Artist:   "Unknown",
		Tags:     []string{"romantic"},
		Art: `
 _
| |    _____   _____
| |   / _ \ \ / / _ \
| |__| (_) \ V /  __/
|_____\___/ \_/ \___|
`,
	},
}

// List returns all art in the database
func List() []Art {
	return database
}

// ListCategories returns unique categories
func ListCategories() []string {
	seen := make(map[string]bool)
	var cats []string
	for _, a := range database {
		if !seen[a.Category] {
			seen[a.Category] = true
			cats = append(cats, a.Category)
		}
	}
	return cats
}

// Get returns art by name
func Get(name string) (Art, bool) {
	name = strings.ToLower(name)
	for _, a := range database {
		if strings.ToLower(a.Name) == name {
			return a, true
		}
	}
	return Art{}, false
}

// Search searches art by name, category, or tags
func Search(query string) []Art {
	query = strings.ToLower(query)
	var results []Art
	for _, a := range database {
		if strings.Contains(strings.ToLower(a.Name), query) ||
			strings.Contains(strings.ToLower(a.Category), query) {
			results = append(results, a)
			continue
		}
		for _, tag := range a.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				results = append(results, a)
				break
			}
		}
	}
	return results
}

// ByCategory returns all art in a category
func ByCategory(category string) []Art {
	category = strings.ToLower(category)
	var results []Art
	for _, a := range database {
		if strings.ToLower(a.Category) == category {
			results = append(results, a)
		}
	}
	return results
}

// Random returns a random piece of art
func Random() Art {
	if len(database) == 0 {
		return Art{}
	}
	idx := rand.IntN(len(database))
	return database[idx]
}
