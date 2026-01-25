package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestEscapeQuits(t *testing.T) {
	m := NewModel()

	// Simulate pressing escape key
	escMsg := tea.KeyMsg(tea.Key{Type: tea.KeyEscape})
	t.Logf("Escape msg: type=%d string=%q", escMsg.Type, escMsg.String())

	newModel, cmd := m.Update(escMsg)
	if cmd == nil {
		t.Fatal("Expected tea.Quit command, got nil")
	}

	// Execute the command to get the message
	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("Expected QuitMsg, got %T", msg)
	}
	_ = newModel
}

func TestEscapeKeyString(t *testing.T) {
	escMsg := tea.KeyMsg(tea.Key{Type: tea.KeyEscape})
	t.Logf("KeyEscape type value: %d", tea.KeyEscape)
	t.Logf("KeyEsc type value: %d", tea.KeyEsc)
	t.Logf("msg.Type: %d", escMsg.Type)
	t.Logf("msg.String(): %q", escMsg.String())
	t.Logf("Type match KeyEscape: %v", escMsg.Type == tea.KeyEscape)
	t.Logf("Type match KeyEsc: %v", escMsg.Type == tea.KeyEsc)
	t.Logf("String match esc: %v", escMsg.String() == "esc")

	if escMsg.String() != "esc" {
		t.Fatalf("Expected String()=%q, got %q", "esc", escMsg.String())
	}
}

func TestCtrlCQuits(t *testing.T) {
	m := NewModel()

	// ctrl+c sends type 3 (ETX)
	ctrlCMsg := tea.KeyMsg(tea.Key{Type: tea.KeyCtrlC})

	_, cmd := m.Update(ctrlCMsg)
	if cmd == nil {
		t.Fatal("Expected tea.Quit command from ctrl+c, got nil")
	}

	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Fatalf("Expected QuitMsg from ctrl+c, got %T", msg)
	}
}
