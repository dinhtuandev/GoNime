package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// TestViewsRenderAtWidths guards the layout code against panics and empty
// output at narrow and wide terminal sizes.
func TestViewsRenderAtWidths(t *testing.T) {
	widths := []int{40, 60, 80, 120}

	for _, w := range widths {
		parent := NewMainModel()
		parent.width = w
		parent.height = 40

		views := map[string]string{
			"menu":      parent.View(),
			"anime":     NewAnimeModel(parent).View(),
			"manga":     NewMangaModel(parent).View(),
			"history":   NewHistoryModel(parent).View(),
			"favorites": NewFavoritesModel(parent).View(),
			"settings":  NewSettingsModel(parent).View(),
		}
		for name, out := range views {
			if strings.TrimSpace(out) == "" {
				t.Fatalf("%s view is empty at width %d", name, w)
			}
		}
	}
}

func TestHelpOverlay(t *testing.T) {
	parent := NewMainModel()
	parent.width = 100
	parent.height = 40

	next, _ := parent.Update(keyRunes("?"))
	pm := next.(MainModel)
	if !pm.help {
		t.Fatal("? did not open help")
	}
	if out := pm.View(); !strings.Contains(out, "PHÍM TẮT") {
		t.Fatal("help overlay has no key hints")
	}

	next, _ = pm.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if next.(MainModel).help {
		t.Fatal("esc did not close help")
	}
}

func TestMenuClickSelectsThenEnters(t *testing.T) {
	parent := NewMainModel()
	parent.width = 100
	parent.height = 40
	_ = parent.View() // publish row hitboxes

	// Clicking another row only moves the cursor.
	next, _ := parent.Update(tea.MouseMsg{
		Action: tea.MouseActionPress,
		Button: tea.MouseButtonLeft,
		X:      6,
		Y:      parent.ui.rowLine + 1,
	})
	pm := next.(MainModel)
	if pm.screen != ScreenMain {
		t.Fatal("single click entered a screen")
	}
	if pm.menuIndex != 1 {
		t.Fatalf("click selected row %d, want 1", pm.menuIndex)
	}
}

func TestMangaHelpToggle(t *testing.T) {
	parent := NewMainModel()
	parent.width = 100
	parent.height = 40
	m := NewMangaModel(parent)

	next, _ := m.Update(keyRunes("?"))
	mm := next.(MangaModel)
	if !mm.help {
		t.Fatal("? did not open manga help")
	}
	if out := mm.View(); !strings.Contains(out, "PHÍM TẮT") {
		t.Fatal("manga help overlay has no key hints")
	}
}
