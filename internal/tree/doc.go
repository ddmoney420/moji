// Package tree provides directory tree visualization.
//
// It generates directory trees with filtering, sorting, and icon support. Trees can be formatted
// with Unicode connectors, depth limiting, and optional file size display.
//
// Example usage:
//
//	tree := tree.Generate(".", opts)
//	output := tree.Format(tree)
//	tree := tree.WithDepth(3)
//	tree := tree.WithIcons()
package tree
