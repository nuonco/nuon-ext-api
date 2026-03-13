package browser

import (
	"testing"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/nuonco/nuon-ext-api/internal/spec"
)

func TestModelUpdateCopyAction(t *testing.T) {
	route := spec.Route{Method: "GET", Path: "/v1/apps"}
	items := []list.Item{routeItem{route: route}}

	l := list.New(items, list.NewDefaultDelegate(), 80, 24)
	m := model{
		list: l,
		keys: keyMap{
			Open:    key.NewBinding(key.WithKeys("B")),
			Copy:    key.NewBinding(key.WithKeys("c")),
			Execute: key.NewBinding(key.WithKeys("x")),
			Select:  key.NewBinding(key.WithKeys("enter")),
		},
	}

	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'c'}})
	if cmd == nil {
		t.Fatal("expected copy action to return a quit command")
	}

	updated, ok := updatedModel.(model)
	if !ok {
		t.Fatalf("expected model type, got %T", updatedModel)
	}

	if updated.result == nil {
		t.Fatal("expected result to be set")
	}
	if updated.result.Action != ActionCopy {
		t.Fatalf("expected action %v, got %v", ActionCopy, updated.result.Action)
	}
	if updated.result.Route == nil {
		t.Fatal("expected route to be set")
	}
	if updated.result.Route.Path != route.Path {
		t.Fatalf("expected copied path %q, got %q", route.Path, updated.result.Route.Path)
	}
}
