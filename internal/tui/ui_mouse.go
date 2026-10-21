package tui

import (
	"strings"

	"gonime/internal/ui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ─── Shared UI registry ───────────────────────────────────────────────────────
//
// Bubble Tea calls View on a copy of the model, so clickable regions cannot be
// stored on the model itself. Every model shares a *uiState pointer through
// MainModel: View publishes regions there, Update reads them for mouse events.

type uiState struct {
	rowLine   int // line index of the first visible row
	rowStart  int // data index of the first visible row
	rowCount  int // number of visible rows
	rowStride int // output lines per row (default 1)
	gridCols  int // >0 when rows are rendered in a grid
	colWidth  int
	hintLine  int
	hintHits  []hintHit
	helpShown bool
}

type hintHit struct {
	x0, x1  int
	y       int
	trigger string
}

func newUIState() *uiState {
	return &uiState{rowLine: -1, hintLine: -1}
}

func (u *uiState) reset() {
	if u == nil {
		return
	}
	u.rowLine = -1
	u.rowStart = 0
	u.rowCount = 0
	u.rowStride = 1
	u.gridCols = 0
	u.hintLine = -1
	u.hintHits = nil
}

// beginRows marks the next output line as the first row of a clickable list.
func (u *uiState) beginRows(b *strings.Builder, start, count int) {
	u.beginRowsStride(b, start, count, 1)
}

// beginRowsStride is beginRows for rows that span several output lines.
func (u *uiState) beginRowsStride(b *strings.Builder, start, count, stride int) {
	if u == nil {
		return
	}
	if stride < 1 {
		stride = 1
	}
	u.rowLine = strings.Count(b.String(), "\n")
	u.rowStart = start
	u.rowCount = count
	u.rowStride = stride
	u.gridCols = 0
}

// beginGrid marks the next output line as the first line of a grid.
func (u *uiState) beginGrid(b *strings.Builder, start, count, cols, colWidth int) {
	if u == nil || cols <= 0 || colWidth <= 0 {
		return
	}
	u.rowLine = strings.Count(b.String(), "\n")
	u.rowStart = start
	u.rowCount = count
	u.gridCols = cols
	u.colWidth = colWidth
}

// rowAt maps a click position to a data index in the current list.
func (u *uiState) rowAt(x, y int) (int, bool) {
	if u == nil || u.rowLine < 0 || y < u.rowLine || u.rowCount <= 0 {
		return 0, false
	}
	line := y - u.rowLine

	if u.gridCols > 0 {
		col := x / u.colWidth
		if col < 0 {
			col = 0
		}
		if col >= u.gridCols {
			col = u.gridCols - 1
		}
		idx := line*u.gridCols + col
		if idx < 0 || idx >= u.rowCount {
			return 0, false
		}
		return u.rowStart + idx, true
	}

	st := u.rowStride
	if st < 1 {
		st = 1
	}
	if line >= u.rowCount*st {
		return 0, false
	}
	idx := line / st
	if idx >= u.rowCount {
		return 0, false
	}
	return u.rowStart + idx, true
}

// hintKeyAt resolves a click on a clickable key hint into that key.
func (u *uiState) hintKeyAt(x, y int) (tea.KeyMsg, bool) {
	if u == nil {
		return tea.KeyMsg{}, false
	}
	for _, h := range u.hintHits {
		if h.y == y && x >= h.x0 && x < h.x1 {
			if msg, ok := keyMsgFor(h.trigger); ok {
				return msg, true
			}
		}
	}
	return tea.KeyMsg{}, false
}

// addKeyHit marks a screen region that replays a key when clicked.
func (u *uiState) addKeyHit(x0, x1, y int, trigger string) {
	if u == nil || x1 <= x0 {
		return
	}
	u.hintHits = append(u.hintHits, hintHit{x0: x0, x1: x1, y: y, trigger: trigger})
}

// shiftRowLine moves the clickable row block down by n output lines.
func (u *uiState) shiftRowLine(n int) {
	if u != nil && u.rowLine >= 0 {
		u.rowLine += n
	}
}

// ─── Footer ───────────────────────────────────────────────────────────────────

// writeStatusBar renders clickable key hints. Each hint is "Key: mô tả".
// Clicking a hint replays the key it advertises.
func writeStatusBar(b *strings.Builder, u *uiState, width int, hints []string) {
	if width <= 0 {
		width = 96
	}
	divWidth := width - 4
	if divWidth < 20 {
		divWidth = 20
	}

	b.WriteString("\n  ")
	b.WriteString(lipgloss.NewStyle().Foreground(theme.TextDim).Render(strings.Repeat("─", divWidth)))
	b.WriteString("\n  ")

	if u != nil {
		u.hintLine = strings.Count(b.String(), "\n")
	}

	x := 2
	for i, raw := range hints {
		if i > 0 {
			sep := theme.Muted.Render("  │  ")
			b.WriteString(sep)
			x += lipgloss.Width(sep)
		}

		key, desc, hasDesc := splitHint(raw)
		keyText := theme.HintKey.Render(key)
		seg := keyText
		if hasDesc {
			seg += theme.Hint.Render(": " + desc)
		}
		segWidth := lipgloss.Width(seg)
		if x+segWidth > width-2 {
			break
		}
		b.WriteString(seg)

		if trigger := hintTrigger(key); trigger != "" && u != nil {
			u.addKeyHit(x, x+segWidth, u.hintLine, trigger)
		}
		x += segWidth
	}
}

func splitHint(raw string) (key, desc string, hasDesc bool) {
	if idx := strings.Index(raw, ":"); idx >= 0 {
		return strings.TrimSpace(raw[:idx]), strings.TrimSpace(raw[idx+1:]), true
	}
	return strings.TrimSpace(raw), "", false
}

// hintTrigger maps the displayed key label to a tea key string.
func hintTrigger(label string) string {
	k := strings.ToLower(strings.TrimSpace(label))
	if k == "" {
		return ""
	}
	if strings.Contains(k, "ctrl+") {
		if idx := strings.Index(k, "ctrl+"); idx >= 0 && len(k) > idx+5 {
			return "ctrl+" + string(k[idx+5])
		}
	}
	switch {
	case strings.HasPrefix(k, "esc"):
		return "esc"
	case strings.HasPrefix(k, "enter"):
		return "enter"
	case strings.Contains(k, "pgup/dn"), strings.Contains(k, "↑↓"), strings.Contains(k, "↓"):
		return "down"
	case strings.Contains(k, "←→"):
		return "right"
	case strings.HasPrefix(k, "g/g"):
		return "g"
	case strings.HasPrefix(k, "n/p"):
		return "n"
	case strings.HasPrefix(k, "1/2"):
		return "1"
	}
	if len([]rune(k)) == 1 {
		return k
	}
	return ""
}

// keyMsgFor converts a key string into a synthetic KeyMsg.
func keyMsgFor(key string) (tea.KeyMsg, bool) {
	if len([]rune(key)) == 1 {
		return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}, true
	}
	switch key {
	case "esc":
		return tea.KeyMsg{Type: tea.KeyEsc}, true
	case "enter":
		return tea.KeyMsg{Type: tea.KeyEnter}, true
	case "up":
		return tea.KeyMsg{Type: tea.KeyUp}, true
	case "down":
		return tea.KeyMsg{Type: tea.KeyDown}, true
	case "left":
		return tea.KeyMsg{Type: tea.KeyLeft}, true
	case "right":
		return tea.KeyMsg{Type: tea.KeyRight}, true
	case "pgup":
		return tea.KeyMsg{Type: tea.KeyPgUp}, true
	case "pgdown":
		return tea.KeyMsg{Type: tea.KeyPgDown}, true
	case "backspace":
		return tea.KeyMsg{Type: tea.KeyBackspace}, true
	case "tab":
		return tea.KeyMsg{Type: tea.KeyTab}, true
	case "ctrl+o":
		return tea.KeyMsg{Type: tea.KeyCtrlO}, true
	case "ctrl+u":
		return tea.KeyMsg{Type: tea.KeyCtrlU}, true
	case "ctrl+d":
		return tea.KeyMsg{Type: tea.KeyCtrlD}, true
	case "ctrl+c":
		return tea.KeyMsg{Type: tea.KeyCtrlC}, true
	}
	return tea.KeyMsg{}, false
}

// ─── Layout helpers ───────────────────────────────────────────────────────────

// listRowWidth is the width of a full-width selected row.
func listRowWidth(termWidth int) int {
	w := termWidth - 8
	if w < 24 {
		w = 24
	}
	return w
}

// listCapacity computes how many rows fit between header/footer.
func listCapacity(windowHeight, reserved int) int {
	if windowHeight <= 0 {
		return 14
	}
	capacity := windowHeight - reserved
	if capacity < 5 {
		capacity = 5
	}
	if capacity > 24 {
		capacity = 24
	}
	return capacity
}

func clampIndex(index, delta, total int) int {
	if total <= 0 {
		return 0
	}
	next := index + delta
	if next < 0 {
		next = 0
	}
	if next >= total {
		next = total - 1
	}
	return next
}

// ─── Help overlay ─────────────────────────────────────────────────────────────

func renderHelp(width, height int, hints []string) string {
	if width <= 0 {
		width = 80
	}
	inner := width - 12
	if inner < 40 {
		inner = 40
	}
	if inner > 64 {
		inner = 64
	}

	var b strings.Builder
	b.WriteString("\n  ")
	b.WriteString(lipgloss.NewStyle().Foreground(theme.Primary).Bold(true).Render("⌨  PHÍM TẮT"))
	b.WriteString("\n\n")

	for _, raw := range hints {
		key, desc, ok := splitHint(raw)
		if !ok {
			key = raw
			desc = ""
		}
		b.WriteString("  ")
		b.WriteString(lipgloss.NewStyle().Foreground(theme.Gold).Bold(true).Width(14).Render(key))
		b.WriteString(theme.Muted.Render(desc))
		b.WriteString("\n")
	}
	b.WriteString("  ")
	b.WriteString(lipgloss.NewStyle().Foreground(theme.Gold).Bold(true).Width(14).Render("?"))
	b.WriteString(theme.Muted.Render("Đóng trợ giúp"))
	b.WriteString("\n")

	card := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Primary).
		Padding(1, 2).
		Width(inner).
		Render(strings.TrimRight(b.String(), "\n"))

	return "\n  " + card + "\n"
}
