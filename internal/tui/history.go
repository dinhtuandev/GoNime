package tui

import (
	"fmt"
	"strings"

	"gonime/internal/api"
	"gonime/internal/config"
	"gonime/internal/player"
	"gonime/internal/ui/i18n"
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
	tab         listTab
	cfg         *config.Config
}

// visibleIndices trả về index thực trong m.items theo tab hiện tại.
func (m HistoryModel) visibleIndices() []int {
	var idxs []int
	for i, it := range m.items {
		switch m.tab {
		case tabAnime:
			if it.Type != "manga" {
				idxs = append(idxs, i)
			}
		case tabManga:
			if it.Type == "manga" {
				idxs = append(idxs, i)
			}
		}
	}
	return idxs
}

func NewHistoryModel(parent MainModel) HistoryModel {
	mgr, _ := config.NewManager()
	items, _ := mgr.LoadHistory()
	cfg, _ := mgr.LoadConfig()
	m := HistoryModel{
		parent:      parent,
		items:       items,
		manager:     mgr,
		animeClient: api.NewAnimeClient(),
		tab:         tabAnime,
		cfg:         cfg,
	}
	return m
}

func (m HistoryModel) Init() tea.Cmd {
	cmds := []tea.Cmd{tickCmd()}
	if vis := m.visibleIndices(); len(vis) > 0 {
		if url := m.items[vis[0]].ThumbURL; url != "" {
			cmds = append(cmds, thumbnailCmd(url))
		}
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
			if vis := m.visibleIndices(); m.index < len(vis) {
				realIdx := vis[m.index]
				m.statusMsg = fmt.Sprintf("▶ Đang phát: %s - %s", m.items[realIdx].Name, m.items[realIdx].Episode)
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
		case "1":
			// Tab Anime
			if m.tab != tabAnime {
				m.tab = tabAnime
				m.index = 0
				m.posterArt = ""
				if vis := m.visibleIndices(); len(vis) > 0 {
					if url := m.items[vis[0]].ThumbURL; url != "" {
						return m, thumbnailCmd(url)
					}
				}
			}
		case "2":
			// Tab Manga
			if m.tab != tabManga {
				m.tab = tabManga
				m.index = 0
				m.posterArt = ""
				if vis := m.visibleIndices(); len(vis) > 0 {
					if url := m.items[vis[0]].ThumbURL; url != "" {
						return m, thumbnailCmd(url)
					}
				}
			}
		case "up", "k":
			if m.index > 0 {
				m.index--
				m.posterArt = ""
				if vis := m.visibleIndices(); m.index < len(vis) {
					if url := m.items[vis[m.index]].ThumbURL; url != "" {
						return m, thumbnailCmd(url)
					}
				}
			}
		case "down", "j":
			if vis := m.visibleIndices(); m.index < len(vis)-1 {
				m.index++
				m.posterArt = ""
				if url := m.items[vis[m.index]].ThumbURL; url != "" {
					return m, thumbnailCmd(url)
				}
			}
		case "enter":
			if vis := m.visibleIndices(); len(vis) > 0 && m.index < len(vis) {
				item := m.items[vis[m.index]]
				m.loading = true
				m.err = nil
				m.statusMsg = ""
				return m, m.playHistoryItem(item)
			}
		case "d":
			if vis := m.visibleIndices(); len(vis) > 0 && m.index < len(vis) {
				realIdx := vis[m.index]
				m.items = append(m.items[:realIdx], m.items[realIdx+1:]...)
				_ = m.manager.SaveHistory(m.items)
				if m.index >= len(m.visibleIndices()) && m.index > 0 {
					m.index--
				}
				m.statusMsg = i18n.T(m.cfg, "msg_del_item")
			}
		case "c":
			if err := m.manager.ClearHistory(); err == nil {
				m.items = nil
				m.index = 0
				m.statusMsg = i18n.T(m.cfg, "msg_del_all")
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

	var tabs []string
	animeTab := theme.BadgeActive.Render(i18n.T(m.cfg, "tab_anime"))
	mangaTab := theme.Badge.Render(i18n.T(m.cfg, "tab_manga"))
	if m.tab == tabManga {
		animeTab = theme.Badge.Render(i18n.T(m.cfg, "tab_anime"))
		mangaTab = theme.BadgeActive.Render(i18n.T(m.cfg, "tab_manga"))
	}
	tabs = append(tabs, animeTab, mangaTab)

	b.WriteString(drawHeader("📜", i18n.T(m.cfg, "title_history"), nil, theme.Primary))
	b.WriteString(drawTabs(tabs))
	b.WriteString(drawHeaderLine(theme.Primary))
	b.WriteString("\n")

	if m.loading {
		b.WriteString(drawLoadingLine("Đang khởi động trình phát...", m.tick))
	} else if m.err != nil {
		b.WriteString(drawErrorLine(fmt.Sprintf("Lỗi: %v", m.err)))
		b.WriteString("\n")
	}

	visible := m.visibleIndices()

	if len(visible) == 0 {
		b.WriteString("  ")
		if m.tab == tabAnime {
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "empty_hist_anime")))
		} else {
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "empty_hist_manga")))
		}
		b.WriteString("\n  ")
		b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "empty_hist_hint")))
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
		for row, idx := range visible {
			item := m.items[idx]
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

			dispName := truncateString(item.Name, 30)

			if row == m.index {
				leftB.WriteString("  " + typeTag + "  " + selStyle.Render("▶  "+dispName) + epStr + timestamp + "\n")
			} else {
				leftB.WriteString("  " + typeTag + "  " + normStyle.Render("   "+dispName) + epStr + timestamp + "\n")
			}
		}
		leftB.WriteString("\n  ")
		leftB.WriteString(theme.HintKey.Render(fmt.Sprintf("%d", len(visible))))
		if m.tab == tabAnime {
			leftB.WriteString(theme.Muted.Render(fmt.Sprintf(" %s", i18n.T(m.cfg, "count_hist_anime"))))
		} else {
			leftB.WriteString(theme.Muted.Render(fmt.Sprintf(" %s", i18n.T(m.cfg, "count_hist_manga"))))
		}

		// Build RIGHT poster panel
		var rightStr string
		if m.posterArt != "" && m.index < len(visible) {
			item := m.items[visible[m.index]]
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
		i18n.T(m.cfg, "hint_esc_back"),
		i18n.T(m.cfg, "hint_nav"),
		i18n.T(m.cfg, "hint_enter_reopen"),
		i18n.T(m.cfg, "hint_del_item"),
		i18n.T(m.cfg, "hint_del_all"),
	}))
	b.WriteString("\n")

	return b.String()
}
