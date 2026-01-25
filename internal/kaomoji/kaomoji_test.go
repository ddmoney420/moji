package kaomoji

import (
	"testing"
)

func TestListAll(t *testing.T) {
	items := List("", "")
	if len(items) == 0 {
		t.Fatal("List() returned empty")
	}
	for _, item := range items {
		if item.Name == "" {
			t.Error("item has empty Name")
		}
		if item.Kaomoji == "" {
			t.Errorf("item %q has empty Kaomoji", item.Name)
		}
	}
}

func TestListSearch(t *testing.T) {
	items := List("shrug", "")
	if len(items) == 0 {
		t.Fatal("search for 'shrug' returned no results")
	}
	found := false
	for _, item := range items {
		if item.Name == "shrug" {
			found = true
		}
	}
	if !found {
		t.Error("search for 'shrug' did not return a 'shrug' item")
	}
}

func TestListCategory(t *testing.T) {
	cats := ListCategories()
	if len(cats) == 0 {
		t.Fatal("ListCategories() returned empty")
	}

	// Filter by first category
	items := List("", cats[0])
	if len(items) == 0 {
		t.Errorf("filtering by category %q returned no results", cats[0])
	}
	for _, item := range items {
		if item.Category != cats[0] {
			t.Errorf("item %q has category %q, want %q", item.Name, item.Category, cats[0])
		}
	}
}

func TestListSearchNoMatch(t *testing.T) {
	items := List("zzzznonexistent999", "")
	if len(items) != 0 {
		t.Errorf("expected 0 results for nonsense search, got %d", len(items))
	}
}

func TestRandom(t *testing.T) {
	name, kao := Random()
	if name == "" {
		t.Error("Random() returned empty name")
	}
	if kao == "" {
		t.Error("Random() returned empty kaomoji")
	}
}

func TestGet(t *testing.T) {
	kao, ok := Get("shrug")
	if !ok {
		t.Fatal("Get('shrug') returned not ok")
	}
	if kao == "" {
		t.Error("Get('shrug') returned empty string")
	}
}

func TestGetNotFound(t *testing.T) {
	_, ok := Get("zzzznonexistent999")
	if ok {
		t.Error("Get() should return false for nonexistent kaomoji")
	}
}

func TestListCategories(t *testing.T) {
	cats := ListCategories()
	if len(cats) == 0 {
		t.Fatal("ListCategories() returned empty")
	}

	seen := make(map[string]bool)
	for _, c := range cats {
		if c == "" {
			t.Error("empty category name")
		}
		if seen[c] {
			t.Errorf("duplicate category: %q", c)
		}
		seen[c] = true
	}
}
