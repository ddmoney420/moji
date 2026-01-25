package artdb

import (
	"testing"
)

func TestList(t *testing.T) {
	arts := List()
	if len(arts) == 0 {
		t.Fatal("List() returned empty")
	}
	// Verify each art has required fields
	for _, a := range arts {
		if a.Name == "" {
			t.Error("Art with empty name")
		}
		if a.Category == "" {
			t.Errorf("Art %q has empty category", a.Name)
		}
		if a.Art == "" {
			t.Errorf("Art %q has empty art content", a.Name)
		}
	}
}

func TestListCategories(t *testing.T) {
	cats := ListCategories()
	if len(cats) == 0 {
		t.Fatal("ListCategories() returned empty")
	}
	// Check for expected categories
	expected := map[string]bool{"animals": false, "symbols": false, "nature": false, "objects": false}
	for _, c := range cats {
		if _, ok := expected[c]; ok {
			expected[c] = true
		}
	}
	for cat, found := range expected {
		if !found {
			t.Errorf("Expected category %q not found", cat)
		}
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name  string
		found bool
	}{
		{"cat", true},
		{"dog", true},
		{"heart", true},
		{"CAT", true}, // case insensitive
		{"Cat", true},
		{"nonexistent_xyz", false},
	}

	for _, tt := range tests {
		art, ok := Get(tt.name)
		if ok != tt.found {
			t.Errorf("Get(%q) found=%v, want %v", tt.name, ok, tt.found)
		}
		if ok && art.Art == "" {
			t.Errorf("Get(%q) returned empty art content", tt.name)
		}
	}
}

func TestSearch(t *testing.T) {
	// Search by name
	results := Search("cat")
	if len(results) == 0 {
		t.Error("Search('cat') returned no results")
	}
	foundCat := false
	for _, r := range results {
		if r.Name == "cat" {
			foundCat = true
		}
	}
	if !foundCat {
		t.Error("Search('cat') didn't include 'cat' art")
	}

	// Search by tag
	results = Search("pet")
	if len(results) == 0 {
		t.Error("Search('pet') returned no results")
	}

	// Search by category
	results = Search("nature")
	if len(results) == 0 {
		t.Error("Search('nature') returned no results")
	}

	// Empty search
	results = Search("zzzzzzzznotfound")
	if len(results) != 0 {
		t.Errorf("Search for nonexistent returned %d results", len(results))
	}
}

func TestByCategory(t *testing.T) {
	animals := ByCategory("animals")
	if len(animals) == 0 {
		t.Fatal("ByCategory('animals') returned empty")
	}
	for _, a := range animals {
		if a.Category != "animals" {
			t.Errorf("ByCategory('animals') returned art with category %q", a.Category)
		}
	}

	// Case insensitive
	animals2 := ByCategory("Animals")
	if len(animals2) != len(animals) {
		t.Error("ByCategory is not case insensitive")
	}

	// Nonexistent category
	results := ByCategory("nonexistent_category")
	if len(results) != 0 {
		t.Errorf("ByCategory('nonexistent') returned %d results", len(results))
	}
}

func TestRandom(t *testing.T) {
	art := Random()
	if art.Name == "" {
		t.Error("Random() returned art with empty name")
	}
	if art.Art == "" {
		t.Error("Random() returned art with empty content")
	}

	// Call multiple times to ensure no panics
	for i := 0; i < 100; i++ {
		a := Random()
		if a.Name == "" {
			t.Fatalf("Random() returned empty on iteration %d", i)
		}
	}
}
