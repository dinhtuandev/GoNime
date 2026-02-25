package theme

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// spinnerFrames: braille spinner for an animated loading feel
var spinnerFrames = []string{"⣾", "⣽", "⣻", "⢿", "⡿", "⣟", "⣯", "⣷"}

// LoadingSpinner returns the spinner frame for the given tick count.
// Call with an incrementing tick (e.g. from a tea.Tick command) and take modulo internally.
func LoadingSpinner(tick int) string {
	frame := spinnerFrames[tick%len(spinnerFrames)]
	return Loading.Render(frame)
}

// ProgressBar renders a colorful progress bar.
//   - current: current value (0-based index, pass currentEp)
//   - total:   total items
//   - width:   number of character columns for the bar
func ProgressBar(current, total, width int) string {
	if total == 0 || width <= 0 {
		return Muted.Render(strings.Repeat("░", width))
	}

	filled := int(float64(current) / float64(total) * float64(width))
	if filled > width {
		filled = width
	}
	empty := width - filled

	filledStyle := lipgloss.NewStyle().Foreground(Primary)
	emptyStyle := lipgloss.NewStyle().Foreground(TextDim)

	return filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))
}

// HorizontalRule renders a full-width neon rule in the given color.
func HorizontalRule(width int, color lipgloss.Color) string {
	if width <= 0 {
		width = 60
	}
	ruleStyle := lipgloss.NewStyle().Foreground(color)
	return ruleStyle.Render(strings.Repeat("─", width))
}

// SectionDivider renders an inline "── TITLE ──" style divider.
func SectionDivider(title string, width int, color lipgloss.Color) string {
	if width <= 0 {
		width = 60
	}
	side := 3
	if width > len(title)+8 {
		side = (width - len(title) - 2) / 2
	}
	s := lipgloss.NewStyle().Foreground(color)
	t := lipgloss.NewStyle().Foreground(color).Bold(true)
	line := s.Render(strings.Repeat("─", side)) + " " +
		t.Render(title) + " " +
		s.Render(strings.Repeat("─", side))
	return line
}

// GlowTitle renders text with a vivid neon border using ThickBorder.
func GlowTitle(text string, color lipgloss.Color) string {
	return lipgloss.NewStyle().
		Foreground(color).
		Bold(true).
		Border(lipgloss.ThickBorder()).
		BorderForeground(color).
		Padding(0, 2).
		Render(text)
}

// GradientLabel renders text alternating between two colors character by character.
// Since lipgloss doesn't support per-char gradients, we alternate on words.
func GradientLabel(text string, c1, c2 lipgloss.Color) string {
	words := strings.Fields(text)
	var sb strings.Builder
	for i, w := range words {
		if i > 0 {
			sb.WriteString(" ")
		}
		if i%2 == 0 {
			sb.WriteString(lipgloss.NewStyle().Foreground(c1).Bold(true).Render(w))
		} else {
			sb.WriteString(lipgloss.NewStyle().Foreground(c2).Bold(true).Render(w))
		}
	}
	return sb.String()
}
