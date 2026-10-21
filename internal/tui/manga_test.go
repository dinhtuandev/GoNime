package tui

import (
	"testing"
	"time"

	"gonime/internal/scraper"

	tea "github.com/charmbracelet/bubbletea"
)

func keyRunes(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func TestMangaChapterNavDirection(t *testing.T) {
	m := MangaModel{
		screen:      mangaScreenReading,
		chapterList: []scraper.Chapter{{Title: "1"}, {Title: "2"}, {Title: "3"}},
	}
	m.currentChapIdx = 1

	next, cmd := m.Update(keyRunes("n"))
	if cmd == nil {
		t.Fatal("n did not request the next chapter")
	}
	nm := next.(MangaModel)
	if nm.currentChapIdx != 2 {
		t.Fatalf("n moved to chapter %d, want 2", nm.currentChapIdx)
	}

	nm.lastPlayTime = time.Time{} // bypass the 300ms debounce
	prev, cmd := nm.Update(keyRunes("p"))
	if cmd == nil {
		t.Fatal("p did not request the previous chapter")
	}
	pm := prev.(MangaModel)
	if pm.currentChapIdx != 1 {
		t.Fatalf("p moved to chapter %d, want 1", pm.currentChapIdx)
	}
}

func TestMangaEscRestoresOrigin(t *testing.T) {
	m := MangaModel{
		screen:       mangaScreenChapters,
		returnScreen: mangaScreenSearch,
		returnIndex:  4,
		mangaList:    []scraper.Manga{{Name: "a"}, {Name: "b"}},
		chapterList:  []scraper.Chapter{{Title: "1"}},
	}

	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	nm := next.(MangaModel)
	if nm.screen != mangaScreenSearch {
		t.Fatalf("esc landed on screen %v, want search", nm.screen)
	}
	if nm.index != 4 {
		t.Fatalf("esc restored index %d, want 4", nm.index)
	}
}

func TestMangaEmptyChapterPageDown(t *testing.T) {
	m := MangaModel{screen: mangaScreenChapters}
	next, _ := m.Update(tea.KeyMsg{Type: tea.KeyPgDown})
	if nm := next.(MangaModel); nm.index != 0 {
		t.Fatalf("pgdown on empty list gave index %d, want 0", nm.index)
	}
}
