package tui

import (
	"strings"
	"testing"
)

func TestRowAtStride(t *testing.T) {
	var b strings.Builder
	u := newUIState()
	b.WriteString("header\n")
	u.beginRowsStride(&b, 10, 3, 2) // rows at lines 1,3,5
	b.WriteString("row0\nblank\nrow1\nblank\nrow2\nblank\n")

	cases := []struct {
		y      int
		want   int
		wantOK bool
	}{
		{0, 0, false},
		{1, 10, true},
		{2, 10, true},
		{3, 11, true},
		{5, 12, true},
		{7, 0, false},
	}
	for _, tc := range cases {
		got, ok := u.rowAt(0, tc.y)
		if ok != tc.wantOK || (ok && got != tc.want) {
			t.Fatalf("rowAt(y=%d) = (%d,%v), want (%d,%v)", tc.y, got, ok, tc.want, tc.wantOK)
		}
	}
}

func TestRowAtGrid(t *testing.T) {
	var b strings.Builder
	u := newUIState()
	u.beginGrid(&b, 0, 12, 5, 14)
	_, _ = u.rowAt(0, 0)

	if idx, ok := u.rowAt(60, 0); !ok || idx != 4 {
		t.Fatalf("grid click col=4 row=0 = (%d,%v), want (4,true)", idx, ok)
	}
	if idx, ok := u.rowAt(3, 1); !ok || idx != 5 {
		t.Fatalf("grid click col=0 row=1 = (%d,%v), want (5,true)", idx, ok)
	}
	if idx, ok := u.rowAt(0, 2); !ok || idx != 10 {
		t.Fatalf("grid click col=0 row=2 = (%d,%v), want (10,true)", idx, ok)
	}
	if _, ok := u.rowAt(0, 3); ok {
		t.Fatal("grid click past last row reported a hit")
	}
}

func TestShiftRowLine(t *testing.T) {
	var b strings.Builder
	u := newUIState()
	b.WriteString("one\ntwo\n") // next line is 2
	u.beginRows(&b, 0, 2)
	u.shiftRowLine(5)
	if got, ok := u.rowAt(0, 7); !ok || got != 0 {
		t.Fatalf("shifted rowAt = (%d,%v), want (0,true)", got, ok)
	}
}

func TestHintKeyAt(t *testing.T) {
	var b strings.Builder
	u := newUIState()
	u.addKeyHit(20, 24, 4, "2") // tabs are registered before the footer
	writeStatusBar(&b, u, 100, []string{"Esc: Quay lại", "s: Tìm"})

	// The footer line is the last line written.
	y := strings.Count(b.String(), "\n")
	if y > 0 {
		y-- // Count reports the trailing newline as well.
	}
	if _, ok := u.hintKeyAt(0, y); ok {
		t.Fatal("click on border reported a hit")
	}
	if msg, ok := u.hintKeyAt(3, u.hintLine); !ok || msg.String() != "esc" {
		t.Fatalf("hint click = (%q,%v), want esc", msg.String(), ok)
	}

	if msg, ok := u.hintKeyAt(21, 4); !ok || msg.String() != "2" {
		t.Fatalf("tab hit = (%q,%v), want 2", msg.String(), ok)
	}
	if _, ok := u.hintKeyAt(21, 5); ok {
		t.Fatal("hit matched on the wrong line")
	}
}

func TestClampIndex(t *testing.T) {
	cases := []struct {
		index, delta, total, want int
	}{
		{0, 5, 0, 0},
		{0, -1, 10, 0},
		{4, 3, 10, 7},
		{9, 3, 10, 9},
		{9, -10, 10, 0},
	}
	for _, tc := range cases {
		if got := clampIndex(tc.index, tc.delta, tc.total); got != tc.want {
			t.Fatalf("clampIndex(%d,%d,%d) = %d, want %d", tc.index, tc.delta, tc.total, got, tc.want)
		}
	}
}

func TestHintTrigger(t *testing.T) {
	cases := map[string]string{
		"Esc":       "esc",
		"Enter":     "enter",
		"Ctrl+O":    "ctrl+o",
		"↑↓":        "down",
		"←→":        "right",
		"n/p":       "n",
		"s":         "s",
		"1/2":       "1",
		"something": "",
	}
	for label, want := range cases {
		if got := hintTrigger(label); got != want {
			t.Fatalf("hintTrigger(%q) = %q, want %q", label, got, want)
		}
	}
}
