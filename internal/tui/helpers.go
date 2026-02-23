package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// drawInputBox vẽ input box thủ công bằng Unicode box-drawing characters
// Tránh dùng lipgloss Border vì gây lệch khung.
//
// Output:
//
//	Nhập tên anime:
//	╭──────────────────────────────────────────────────────╮
//	│ bleach▌                                              │
//	╰──────────────────────────────────────────────────────╯
func drawInputBox(label, value string, width int, borderColor string) string {
	var b strings.Builder

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(borderColor))
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f0f4ff"))

	// Đảm bảo width hợp lý
	if width < 20 {
		width = 20
	}

	innerWidth := width - 4 // trừ "│ " bên trái + " │" bên phải

	// Pad hoặc trim value cho vừa inner width
	displayVal := value
	valRunes := []rune(displayVal)
	if len(valRunes) > innerWidth {
		displayVal = string(valRunes[len(valRunes)-innerWidth:])
	}
	padding := innerWidth - len([]rune(displayVal))
	if padding < 0 {
		padding = 0
	}

	// Tạo horizontal bar
	bar := strings.Repeat("─", width-2)

	// Label
	b.WriteString("  ")
	b.WriteString(labelStyle.Render(label))
	b.WriteString("\n")

	// Top border: ╭───────╮
	b.WriteString("  ")
	b.WriteString(borderStyle.Render("╭" + bar + "╮"))
	b.WriteString("\n")

	// Content:    │ text   │
	b.WriteString("  ")
	b.WriteString(borderStyle.Render("│"))
	b.WriteString(" ")
	b.WriteString(textStyle.Render(displayVal))
	b.WriteString(strings.Repeat(" ", padding))
	b.WriteString(" ")
	b.WriteString(borderStyle.Render("│"))
	b.WriteString("\n")

	// Bottom:     ╰───────╯
	b.WriteString("  ")
	b.WriteString(borderStyle.Render("╰" + bar + "╯"))

	return b.String()
}
