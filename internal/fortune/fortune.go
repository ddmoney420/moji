package fortune

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Fortunes - collection of fun quotes
var fortunes = []string{
	// Programming wisdom
	"There are only two hard things in Computer Science: cache invalidation and naming things.",
	"It works on my machine.",
	"99 little bugs in the code, 99 little bugs. Take one down, patch it around, 127 little bugs in the code.",
	"A good programmer is someone who always looks both ways before crossing a one-way street.",
	"Debugging is like being the detective in a crime movie where you are also the murderer.",
	"Code never lies, comments sometimes do.",
	"First, solve the problem. Then, write the code.",
	"The best error message is the one that never shows up.",
	"Weeks of coding can save you hours of planning.",
	"It's not a bug, it's an undocumented feature.",

	// Muppet quotes
	"It's not easy being green. - Kermit",
	"Life's like a movie, write your own ending. - Kermit",
	"There's not a word yet for old friends who've just met. - Gonzo",

	// Fun wisdom
	"The journey of a thousand miles begins with a single step... and usually a segfault.",
	"In a world full of algorithms, be a heuristic.",
	"May your code compile on the first try.",
	"Keep calm and clear your cache.",
	"Have you tried turning it off and on again?",
	"There's no place like 127.0.0.1",
	"Home is where the WiFi connects automatically.",

	// ASCII art themed
	"In a world of GUIs, ASCII art is true art.",
	"Fortune favors the bold (and the monospaced).",
	"Every pixel tells a story, but a character tells a novel.",
	"The terminal is mightier than the mouse.",

	// Classic fortune cookie style
	"You will debug successfully today.",
	"A pleasant surprise is waiting for you.",
	"Your code will run smoothly today.",
	"Good things come to those who commit often.",
	"Today is a good day to refactor.",
	"Trust your instincts. Also trust your linter.",
}

// Programming jokes
var jokes = []string{
	"Why do programmers prefer dark mode? Because light attracts bugs!",
	"Why did the programmer quit his job? Because he didn't get arrays!",
	"A SQL query walks into a bar, walks up to two tables and asks, 'Can I join you?'",
	"There are 10 types of people: those who understand binary and those who don't.",
	"How many programmers does it take to change a light bulb? None, that's a hardware problem.",
	"Why did the developer go broke? Because he used up all his cache.",
	"What's a programmer's favorite hangout place? Foo Bar.",
	"I would tell you a UDP joke, but you might not get it.",
	"What do you call 8 hobbits? A hobbyte.",
}

// Get returns a random fortune
func Get() string {
	return fortunes[rand.Intn(len(fortunes))]
}

// GetJoke returns a random programming joke
func GetJoke() string {
	return jokes[rand.Intn(len(jokes))]
}

// GetAll returns all fortunes
func GetAll() []string {
	return fortunes
}

// GetCategory returns fortunes matching a keyword
func GetCategory(keyword string) []string {
	var matches []string
	for _, f := range fortunes {
		if containsIgnoreCase(f, keyword) {
			matches = append(matches, f)
		}
	}
	if len(matches) == 0 {
		return []string{Get()}
	}
	return matches
}

func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}
