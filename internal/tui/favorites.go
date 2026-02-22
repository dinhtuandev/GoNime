package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gonime/internal/config"
)

type FavoritesModel struct {
	parent  MainModel
	items   []config.FavoriteItem
	index   int
	manager *config.Manager
}

func NewFavoritesModel(parent MainModel) FavoritesModel {
	mgr, _ := config.NewManager()
	items, _ := mgr.LoadFavorites()
	return FavoritesModel{
		parent:  parent,
		items:   items,
		manager: mgr,
	}
}

func (m FavoritesModel) Init() tea.Cmd {
	return nil
}

func (m FavoritesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				// TODO: Open selected favorite
			}
		case "d":
			if len(m.items) > 0 {
				m.items = append(m.items[:m.index], m.items[m.index+1:]...)
				m.manager.SaveFavorites(m.items)
				if m.index >= len(m.items) && m.index > 0 {
					m.index--
				}
			}
		}
	}

	return m, nil
}

func (m FavoritesModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#33a5f5")).
		Bold(true).
		Margin(1, 0)

	b.WriteString(headerStyle.Render("⭐ YÊU THÍCH"))
	b.WriteString("\n\n")

	if len(m.items) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
		b.WriteString(emptyStyle.Render("Chưa có mục yêu thích nào."))
	} else {
		listStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f0f4ff"))
		selectedListStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9f75ff")).
			Bold(true)

		for i, item := range m.items {
			if i == m.index {
				b.WriteString(selectedListStyle.Render(fmt.Sprintf("▶ %s", item.Name)))
			} else {
				b.WriteString(listStyle.Render(fmt.Sprintf("  %s", item.Name)))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
	b.WriteString(footerStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Mở | d: Xóa"))

	return b.String()
}
