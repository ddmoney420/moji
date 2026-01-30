// Package artdb provides a database of ASCII art with categorization and search capabilities.
//
// It manages a collection of ASCII art pieces with metadata including name, category,
// artist attribution, and searchable tags. The package supports retrieving art by name,
// listing all available pieces, searching by criteria, and filtering by category.
//
// Example usage:
//
//	arts := artdb.List()
//	cat := artdb.Get("cat")
//	animals := artdb.ByCategory("animals")
//	results := artdb.Search("cute")
package artdb
