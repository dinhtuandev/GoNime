package tui

import (
	"strings"
	"unicode/utf8"

	"gonime/internal/ui/theme"

	"github.com/charmbracelet/lipgloss"
)

// ─── Shared Types ─────────────────────────────────────────────────────────────

type listTab int

const (
	tabAnime listTab = iota
	tabManga
)

// truncateString cuts a string to maxLen and appends "..." if it was longer.
func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen && maxLen > 3 {
		return string(runes[:maxLen-3]) + "..."
	}
	return s
}

// ─── drawInputBox ─────────────────────────────────────────────────────────────
// Renders a styled input box with neon border:
//
//	Label:
//	╭──────────────────────────────────────╮
//	│ value▌                               │
//	╰──────────────────────────────────────╯
func drawInputBox(label, value string, width int, borderColor string) string {
	var b strings.Builder

	if width < 20 {
		width = 20
	}

	labelStyle := theme.InputLabel.Copy()
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	textStyle := theme.InputText.Copy()

	innerWidth := width - 4 // subtract │ + space on each side

	valRunes := []rune(value)
	displayVal := value
	if len(valRunes) > innerWidth {
		displayVal = string(valRunes[len(valRunes)-innerWidth:])
	}
	padding := innerWidth - utf8.RuneCountInString(displayVal)
	if padding < 0 {
		padding = 0
	}

	bar := strings.Repeat("─", width-2)

	b.WriteString("  ")
	b.WriteString(labelStyle.Render(label))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(borderStyle.Render("╭" + bar + "╮"))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(borderStyle.Render("│"))
	b.WriteString(" ")
	b.WriteString(textStyle.Render(displayVal))
	b.WriteString(strings.Repeat(" ", padding))
	b.WriteString(" ")
	b.WriteString(borderStyle.Render("│"))
	b.WriteString("\n")

	b.WriteString("  ")
	b.WriteString(borderStyle.Render("╰" + bar + "╯"))

	return b.String()
}

// ─── drawHeader ───────────────────────────────────────────────────────────────
// Renders the top header bar over two lines:
//
//	Line 1: "  📺  ANIME / MOVIE"
//	Line 2: "  🔥 Xu hướng  ·  phimapi"
//
// accent controls the color of the module name text.
func drawHeader(moduleIcon, moduleName string, badges []string, accent lipgloss.Color) string {
	var b strings.Builder

	modStyle := lipgloss.NewStyle().Foreground(accent).Bold(true)
	b.WriteString("\n")
	b.WriteString("  ")
	b.WriteString(modStyle.Render(moduleIcon + "  " + moduleName))

	if len(badges) > 0 {
		b.WriteString("\n  ")
		for i, badge := range badges {
			if i > 0 {
				b.WriteString(theme.Muted.Render("  ·  "))
			}
			b.WriteString(theme.Muted.Render(badge))
		}
	}

	return b.String()
}

// ─── drawTabs ─────────────────────────────────────────────────────────────────
// Renders tabs side by side with gap, properly indenting all lines of
// multi-line badge output (rounded borders span 3 lines).
func drawTabs(tabs []string) string {
	if len(tabs) == 0 {
		return "\n"
	}
	rendered := make([]string, 0, len(tabs))
	space := lipgloss.NewStyle().MarginRight(1)
	for i, tab := range tabs {
		if i < len(tabs)-1 {
			rendered = append(rendered, space.Render(tab))
		} else {
			rendered = append(rendered, tab)
		}
	}
	joined := lipgloss.JoinHorizontal(lipgloss.Top, rendered...)
	// Indent every line so the border box is aligned, not just the first line.
	var b strings.Builder
	b.WriteString("\n")
	for i, line := range strings.Split(joined, "\n") {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString("  ")
		b.WriteString(line)
	}
	return b.String()
}

// ─── drawHeaderLine ───────────────────────────────────────────────────────────
// Renders a full-width ─── divider below the header.
func drawHeaderLine(accent lipgloss.Color) string {
	return "\n  " + lipgloss.NewStyle().Foreground(accent).Render(strings.Repeat("─", 64))
}

// ─── drawStatusBar ────────────────────────────────────────────────────────────
// Renders the bottom hint bar with pipe-separated hints.
// Each hint is "KEY: description" — KEY gets gold color, description gets muted.
//
// Example: drawStatusBar([]string{"Enter: Chọn", "↑↓: Di chuyển", "Esc: Menu"})
func drawStatusBar(hints []string) string {
	var b strings.Builder

	b.WriteString("\n  ")
	b.WriteString(lipgloss.NewStyle().Foreground(theme.TextDim).Render(strings.Repeat("─", 64)))
	b.WriteString("\n  ")

	for i, hint := range hints {
		if i > 0 {
			b.WriteString(theme.Muted.Render("  │  "))
		}
		// Split on first colon to colorize key separately
		if idx := strings.Index(hint, ":"); idx >= 0 {
			key := hint[:idx]
			desc := hint[idx:]
			b.WriteString(theme.HintKey.Render(key))
			b.WriteString(theme.Hint.Render(desc))
		} else {
			b.WriteString(theme.Hint.Render(hint))
		}
	}

	return b.String()
}

// ─── drawSectionTitle ─────────────────────────────────────────────────────────
// Renders a small "── SECTION ──" inline divider
func drawSectionTitle(title string, accent lipgloss.Color) string {
	s := lipgloss.NewStyle().Foreground(accent)
	t := lipgloss.NewStyle().Foreground(accent).Bold(true)
	return "  " + s.Render("── ") + t.Render(title) + s.Render(" ──")
}

// ─── drawLoadingLine ──────────────────────────────────────────────────────────
// Renders a loading line with spinner glyph.
func drawLoadingLine(msg string, tick int) string {
	spinner := theme.LoadingSpinner(tick)
	return "\n  " + spinner + "  " + theme.Loading.Render(msg)
}

// ─── drawErrorLine ────────────────────────────────────────────────────────────
func drawErrorLine(msg string) string {
	return "\n  " + theme.ErrorText.Render("✗  "+msg)
}

// ─── drawSuccessLine ──────────────────────────────────────────────────────────
func drawSuccessLine(msg string) string {
	return "\n  " + theme.SuccessText.Render("✓  "+msg)
}

// ─── drawWarningLine ──────────────────────────────────────────────────────────
func drawWarningLine(msg string) string {
	return "\n  " + theme.WarningText.Render("⚠  "+msg)
}

// ─── drawMetaLine ─────────────────────────────────────────────────────────────
// Renders a row of "• item" metadata pieces separated by spaces.
func drawMetaLine(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	dotStyle := lipgloss.NewStyle().Foreground(theme.TextDim)
	textStyle := theme.Subtitle.Copy()
	var b strings.Builder
	b.WriteString("\n  ")
	for i, p := range parts {
		if i > 0 {
			b.WriteString(dotStyle.Render("  ·  "))
		}
		b.WriteString(textStyle.Render(p))
	}
	return b.String()
}

// wordWrap wraps text at word boundaries so each line fits within maxWidth runes.
func wordWrap(text string, maxWidth int) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	var lines []string
	current := ""
	for _, w := range words {
		if current == "" {
			current = w
		} else if len([]rune(current))+1+len([]rune(w)) <= maxWidth {
			current += " " + w
		} else {
			lines = append(lines, current)
			current = w
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return strings.Join(lines, "\n")
}

func contentWidth(windowWidth int) int {
	if windowWidth <= 0 {
		return 96
	}
	if windowWidth < 68 {
		return 64
	}
	return windowWidth - 4
}

func splitPaneWidths(windowWidth int, minLeft int, preferredRight int, minRight int) (int, int) {
	total := contentWidth(windowWidth)
	if total <= minLeft+2 {
		return total, 0
	}
	right := preferredRight
	if total < minLeft+preferredRight+2 {
		right = total - minLeft - 2
	}
	if right < minRight {
		return total, 0
	}
	left := total - right - 2
	if left < minLeft {
		return total, 0
	}
	return left, right
}

func windowedRange(total int, index int, maxRows int) (int, int) {
	if total <= 0 {
		return 0, 0
	}
	if maxRows <= 0 || maxRows >= total {
		return 0, total
	}
	if index < 0 {
		index = 0
	}
	if index >= total {
		index = total - 1
	}
	start := index - maxRows/2
	if start < 0 {
		start = 0
	}
	end := start + maxRows
	if end > total {
		end = total
		start = end - maxRows
		if start < 0 {
			start = 0
		}
	}
	return start, end
}
