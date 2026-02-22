package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gonime/internal/config"
)

type HistoryModel struct {
	parent  MainModel
	items   []config.HistoryItem
	index   int
	manager *config.Manager
}

func NewHistoryModel(parent MainModel) HistoryModel {
	mgr, _ := config.NewManager()
	items, _ := mgr.LoadHistory()
	return HistoryModel{
		parent:  parent,
		items:   items,
		manager: mgr,
	}
}

func (m HistoryModel) Init() tea.Cmd {
	return nil
}

func (m HistoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m.parent, nil
		case "up", "k":
			if m.index > 0 {
				m.index--
			}
		case "down", "j":
			if m.index < len(m.items)-1 {
				m.index++
			}
		case "enter":
			if len(m.items) > 0 {
				// TODO: Play selected history item
			}
		case "d":
			if len(m.items) > 0 {
				m.items = append(m.items[:m.index], m.items[m.index+1:]...)
				m.manager.SaveHistory(m.items)
				if m.index >= len(m.items) && m.index > 0 {
					m.index--
				}
			}
		}
	}

	return m, nil
}

func (m HistoryModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#33a5f5")).
		Bold(true).
		Margin(1, 0)

	b.WriteString(headerStyle.Render("📜 LỊCH SỬ XEM"))
	b.WriteString("\n\n")

	if len(m.items) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
		b.WriteString(emptyStyle.Render("Chưa có lịch sử xem."))
	} else {
		listStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f0f4ff"))
		selectedListStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9f75ff")).
			Bold(true)

		for i, item := range m.items {
			line := fmt.Sprintf("%s - %s (%s)", item.Name, item.Episode, item.Timestamp)
			if i == m.index {
				b.WriteString(selectedListStyle.Render("▶ " + line))
			} else {
				b.WriteString(listStyle.Render("  " + line))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
	b.WriteString(footerStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Xem lại | d: Xóa"))

	return b.String()
}
