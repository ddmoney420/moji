// Package kaomoji provides a database of kaomoji emoticons organized by category.
//
// It maintains 150+ Japanese-style emoticons with categories (animals, emotions, actions, magic, etc.).
// ASCII art variants are included. The package supports searching, random selection, and suggestions.
//
// Example usage:
//
//	kaomoji := kaomoji.Get("happy")
//	random := kaomoji.Random()
//	animals := kaomoji.ByCategory("animals")
//	suggestions := kaomoji.Suggest("cute")
//	all := kaomoji.List()
package kaomoji
