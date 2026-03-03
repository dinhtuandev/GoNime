package tui

import (
	"fmt"
	"strings"

	"gonime/internal/api"
	"gonime/internal/config"
	"gonime/internal/player"
	"gonime/internal/ui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HistoryModel struct {
	parent      MainModel
	items       []config.HistoryItem
	index       int
	loading     bool
	err         error
	statusMsg   string
	manager     *config.Manager
	animeClient *api.AnimeClient
	posterArt   string // chafa-rendered thumbnail
	tick        int
}

func NewHistoryModel(parent MainModel) HistoryModel {
	mgr, _ := config.NewManager()
	items, _ := mgr.LoadHistory()
	m := HistoryModel{
		parent:      parent,
		items:       items,
		manager:     mgr,
		animeClient: api.NewAnimeClient(),
	}
	return m
}

func (m HistoryModel) Init() tea.Cmd {
	cmds := []tea.Cmd{tickCmd()}
	if len(m.items) > 0 && m.items[0].ThumbURL != "" {
		cmds = append(cmds, posterCmd(m.items[0].ThumbURL))
	}
	return tea.Batch(cmds...)
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

	case posterRenderMsg:
		if msg.err == nil {
			m.posterArt = msg.rendered
		}
		return m, nil

	case tickMsg:
		if m.loading {
			m.tick++
		}
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m.parent, nil
		case "up", "k":
			if m.index > 0 {
				m.index--
				m.posterArt = ""
				if m.index < len(m.items) && m.items[m.index].ThumbURL != "" {
					return m, posterCmd(m.items[m.index].ThumbURL)
				}
			}
		case "down", "j":
			if m.index < len(m.items)-1 {
				m.index++
				m.posterArt = ""
				if m.index < len(m.items) && m.items[m.index].ThumbURL != "" {
					return m, posterCmd(m.items[m.index].ThumbURL)
				}
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

	b.WriteString(drawHeader("📜", "LỊCH SỬ XEM", nil, theme.Primary))
	b.WriteString(drawHeaderLine(theme.Primary))
	b.WriteString("\n")

	if m.loading {
		b.WriteString(drawLoadingLine("Đang khởi động trình phát...", m.tick))
	} else if m.err != nil {
		b.WriteString(drawErrorLine(fmt.Sprintf("Lỗi: %v", m.err)))
		b.WriteString("\n")
	}

	if len(m.items) == 0 {
		b.WriteString("  ")
		b.WriteString(theme.Muted.Render("Chưa có lịch sử xem."))
		b.WriteString("\n  ")
		b.WriteString(theme.Muted.Render("Xem anime/manga rồi quay lại đây nhé!"))
	} else {
		selStyle := lipgloss.NewStyle().
			Foreground(theme.Primary).
			Background(theme.BackgroundLight).
			Bold(true).
			Padding(0, 1)
		normStyle := lipgloss.NewStyle().
			Foreground(theme.TextSecondary).
			Padding(0, 1)
		timeStyle := lipgloss.NewStyle().Foreground(theme.TextMuted)
		animeTypeStyle := lipgloss.NewStyle().Foreground(theme.AnimeAccent)
		mangaTypeStyle := lipgloss.NewStyle().Foreground(theme.MangaAccent)

		// Build LEFT list
		var leftB strings.Builder
		leftB.WriteString("\n")
		for i, item := range m.items {
			var typeTag string
			if item.Type == "manga" {
				typeTag = mangaTypeStyle.Render("📚")
			} else {
				typeTag = animeTypeStyle.Render("🎦")
			}
			epStr := ""
			if item.Episode != "" {
				epStr = "  " + theme.Muted.Render("·  "+item.Episode)
			}
			timestamp := ""
			if item.Timestamp != "" {
				timestamp = "  " + timeStyle.Render(item.Timestamp)
			}
			if i == m.index {
				leftB.WriteString("  " + typeTag + "  " + selStyle.Render("▶  "+item.Name) + epStr + timestamp + "\n")
			} else {
				leftB.WriteString("  " + typeTag + "  " + normStyle.Render("   "+item.Name) + epStr + timestamp + "\n")
			}
		}
		leftB.WriteString("\n  ")
		leftB.WriteString(theme.HintKey.Render(fmt.Sprintf("%d", len(m.items))))
		leftB.WriteString(theme.Muted.Render(" mục trong lịch sử"))

		// Build RIGHT poster panel
		var rightStr string
		if m.posterArt != "" && m.index < len(m.items) {
			item := m.items[m.index]
			accentColor := theme.AnimeAccent
			if item.Type == "manga" {
				accentColor = theme.MangaAccent
			}
			accentStyle := lipgloss.NewStyle().Foreground(accentColor).Bold(true)
			var info strings.Builder
			info.WriteString("\n" + accentStyle.Render(item.Name) + "\n")
			if item.Episode != "" {
				info.WriteString(theme.Muted.Render(item.Episode) + "\n")
			}
			if item.Timestamp != "" {
				info.WriteString(theme.Muted.Render("📅 "+item.Timestamp) + "\n")
			}
			posterStyle := lipgloss.NewStyle().Padding(0, 1)
			infoStyle := lipgloss.NewStyle().Padding(0, 1).Width(40)
			rightStr = lipgloss.JoinHorizontal(lipgloss.Top,
				posterStyle.Render(m.posterArt),
				infoStyle.Render(info.String()),
			)
		}

		leftStyle := lipgloss.NewStyle().Width(52)
		rightStyle := lipgloss.NewStyle().Width(40)
		if rightStr == "" {
			rightStr = rightStyle.Render("")
		} else {
			rightStr = rightStyle.Render(rightStr)
		}
		b.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, leftStyle.Render(leftB.String()), rightStr))
	}

	if m.statusMsg != "" {
		b.WriteString(drawSuccessLine(m.statusMsg))
	}

	b.WriteString(drawStatusBar([]string{
		"Esc: Quay lại",
		"↑↓: Di chuyển",
		"Enter: Mở lại",
		"d: Xóa mục",
		"c: Xóa tất cả",
	}))
	b.WriteString("\n")

	return b.String()
}
