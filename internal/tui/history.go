package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gonime/internal/api"
	"gonime/internal/config"
	"gonime/internal/player"
)

type HistoryModel struct {
	parent     MainModel
	items      []config.HistoryItem
	index      int
	loading    bool
	err        error
	statusMsg  string
	manager    *config.Manager
	animeClient *api.AnimeClient
}

func NewHistoryModel(parent MainModel) HistoryModel {
	mgr, _ := config.NewManager()
	items, _ := mgr.LoadHistory()
	return HistoryModel{
		parent:      parent,
		items:       items,
		manager:     mgr,
		animeClient: api.NewAnimeClient(),
	}
}

func (m HistoryModel) Init() tea.Cmd {
	return nil
}

func (m HistoryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case playAnimeMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			if m.index < len(m.items) {
				m.statusMsg = fmt.Sprintf("▶ Đang phát: %s - %s", m.items[m.index].Name, m.items[m.index].Episode)
			}
		}
		return m, nil

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
			if len(m.items) > 0 && m.index < len(m.items) {
				item := m.items[m.index]
				m.loading = true
				m.err = nil
				m.statusMsg = ""
				return m, m.playHistoryItem(item)
			}
		case "d":
			if len(m.items) > 0 && m.index < len(m.items) {
				m.items = append(m.items[:m.index], m.items[m.index+1:]...)
				_ = m.manager.SaveHistory(m.items)
				if m.index >= len(m.items) && m.index > 0 {
					m.index--
				}
				m.statusMsg = "Đã xóa khỏi lịch sử"
			}
		case "c":
			if err := m.manager.ClearHistory(); err == nil {
				m.items = nil
				m.index = 0
				m.statusMsg = "Đã xóa toàn bộ lịch sử"
			}
		}
	}

	return m, nil
}

func (m HistoryModel) playHistoryItem(item config.HistoryItem) tea.Cmd {
	return func() tea.Msg {
		if item.URL == "" {
			return playAnimeMsg{err: fmt.Errorf("không có link phát")}
		}

		// Lấy config
		var playerName = player.GetDefaultPlayer()
		if m.manager != nil {
			cfg, err := m.manager.LoadConfig()
			if err == nil && cfg.Player != "" {
				playerName = cfg.Player
			}
		}

		var p *player.Player
		var err error
		if playerName != "" {
			p, err = player.NewPlayer(playerName)
		}
		if p == nil || err != nil {
			p, err = player.FindAnyPlayer()
			if err != nil {
				return playAnimeMsg{err: fmt.Errorf("%v", err)}
			}
		}

		title := fmt.Sprintf("GoNime: %s - %s", item.Name, item.Episode)
		err = p.Play(item.URL, title)
		return playAnimeMsg{err: err}
	}
}

func (m HistoryModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#33a5f5")).
		Bold(true).
		Margin(1, 0)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9f75ff")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f0f4ff"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a0a8c8"))

	errStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff6b6b"))

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4caf50"))

	b.WriteString(headerStyle.Render("📜 LỊCH SỬ XEM"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(hintStyle.Render("⏳ Đang khởi động trình phát..."))
	} else if m.err != nil {
		b.WriteString(errStyle.Render(fmt.Sprintf("❌ Lỗi: %v", m.err)))
		b.WriteString("\n")
	}

	if len(m.items) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
		b.WriteString(emptyStyle.Render("Chưa có lịch sử xem."))
	} else {
		typeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#33a5f5"))
		timeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))

		for i, item := range m.items {
			typeIcon := "🎬"
			if item.Type == "manga" {
				typeIcon = "📚"
			}

			line := fmt.Sprintf("%s %s - %s", typeIcon, item.Name, item.Episode)

			if i == m.index {
				b.WriteString(selectedStyle.Render("▶ " + line))
				if item.Timestamp != "" {
					b.WriteString("  ")
					b.WriteString(timeStyle.Render(item.Timestamp))
				}
			} else {
				b.WriteString(normalStyle.Render("  " + line))
				if item.Timestamp != "" {
					b.WriteString("  ")
					b.WriteString(timeStyle.Render(item.Timestamp))
				}
			}
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(typeStyle.Render(fmt.Sprintf("Tổng: %d mục", len(m.items))))
	}

	if m.statusMsg != "" {
		b.WriteString("\n")
		b.WriteString(successStyle.Render(m.statusMsg))
	}

	b.WriteString("\n\n")
	b.WriteString(hintStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Mở lại | d: Xóa mục | c: Xóa tất cả"))

	return b.String()
}
