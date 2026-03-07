package tui

import "testing"

func TestWindowedRange(t *testing.T) {
	tests := []struct {
		name               string
		total, index       int
		maxRows            int
		wantStart, wantEnd int
	}{
		{name: "empty", total: 0, index: 0, maxRows: 10, wantStart: 0, wantEnd: 0},
		{name: "all visible", total: 5, index: 2, maxRows: 10, wantStart: 0, wantEnd: 5},
		{name: "middle window", total: 30, index: 15, maxRows: 10, wantStart: 10, wantEnd: 20},
		{name: "near start", total: 30, index: 1, maxRows: 10, wantStart: 0, wantEnd: 10},
		{name: "near end", total: 30, index: 29, maxRows: 10, wantStart: 20, wantEnd: 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := windowedRange(tt.total, tt.index, tt.maxRows)
			if gotStart != tt.wantStart || gotEnd != tt.wantEnd {
				t.Fatalf("windowedRange() = (%d,%d), want (%d,%d)", gotStart, gotEnd, tt.wantStart, tt.wantEnd)
			}
		})
	}
}

func TestSplitPaneWidths(t *testing.T) {
	left, right := splitPaneWidths(120, 54, 42, 34)
	if left <= 0 || right <= 0 {
		t.Fatalf("expected both panes visible for wide screen, got left=%d right=%d", left, right)
	}

	leftNarrow, rightNarrow := splitPaneWidths(70, 54, 42, 34)
	if rightNarrow != 0 {
		t.Fatalf("expected right pane hidden on narrow width, got left=%d right=%d", leftNarrow, rightNarrow)
	}
	if leftNarrow <= 0 {
		t.Fatalf("expected left pane still visible, got left=%d", leftNarrow)
	}
}
