package fortune

import (
	"testing"
)

func TestGet(t *testing.T) {
	f := Get()
	if f == "" {
		t.Error("Get() returned empty string")
	}

	// Should return different values sometimes (probabilistic)
	seen := make(map[string]bool)
	for i := 0; i < 50; i++ {
		seen[Get()] = true
	}
	if len(seen) < 2 {
		t.Error("Get() appears to always return the same fortune")
	}
}

func TestGetJoke(t *testing.T) {
	j := GetJoke()
	if j == "" {
		t.Error("GetJoke() returned empty string")
	}

	seen := make(map[string]bool)
	for i := 0; i < 50; i++ {
		seen[GetJoke()] = true
	}
	if len(seen) < 2 {
		t.Error("GetJoke() appears to always return the same joke")
	}
}

func TestGetAll(t *testing.T) {
	all := GetAll()
	if len(all) == 0 {
		t.Fatal("GetAll() returned empty slice")
	}
	for i, f := range all {
		if f == "" {
			t.Errorf("GetAll()[%d] is empty", i)
		}
	}
}

func TestGetCategory(t *testing.T) {
	// Known keyword
	results := GetCategory("bug")
	if len(results) == 0 {
		t.Error("GetCategory('bug') returned empty")
	}

	// Programming keyword
	results = GetCategory("code")
	if len(results) == 0 {
		t.Error("GetCategory('code') returned empty")
	}

	// Nonexistent keyword returns a fallback
	results = GetCategory("zzzzzzzzzzz")
	if len(results) == 0 {
		t.Error("GetCategory with no match should return fallback")
	}
}

func TestContainsIgnoreCase(t *testing.T) {
	tests := []struct {
		s, substr string
		want      bool
	}{
		{"Hello World", "hello", true},
		{"Hello World", "WORLD", true},
		{"Hello World", "xyz", false},
		{"", "", true},
		{"abc", "", true},
		{"", "abc", false},
	}

	for _, tt := range tests {
		got := containsIgnoreCase(tt.s, tt.substr)
		if got != tt.want {
			t.Errorf("containsIgnoreCase(%q, %q) = %v, want %v", tt.s, tt.substr, got, tt.want)
		}
	}
}
