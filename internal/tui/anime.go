package tui

import (
	"fmt"
	"strings"
	"time"

	"gonime/internal/api"
	"gonime/internal/config"
	"gonime/internal/player"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Messages ---

type searchAnimeMsg struct {
	items  []api.AnimeItem
	cdnURL string
	err    error
}

type trendingAnimeMsg struct {
	items  []api.AnimeItem
	cdnURL string
	err    error
}

type loadEpisodesMsg struct {
	movie    *api.MovieInfo
	episodes []api.EpisodeData
	servers  []api.EpisodeServer
	err      error
}

type playAnimeMsg struct {
	err error
}

// --- State ---

type animeScreen int

const (
	animeScreenTrending animeScreen = iota
	animeScreenSearch
	animeScreenEpisodes
)

// --- Model ---

type AnimeModel struct {
	parent        MainModel
	screen        animeScreen
	searchQuery   string
	animeList     []api.AnimeItem
	trendingList  []api.AnimeItem
	cdnURL        string
	selectedAnime *api.AnimeItem
	movie         *api.MovieInfo
	servers       []api.EpisodeServer
	serverIdx     int
	episodeList   []api.EpisodeData
	currentEpIdx  int
	index         int
	loading       bool
	err           error
	client        *api.AnimeClient
	cfgManager    *config.Manager
	statusMsg     string
}

func NewAnimeModel(parent MainModel) AnimeModel {
	mgr, _ := config.NewManager()
	return AnimeModel{
		parent:     parent,
		screen:     animeScreenTrending,
		client:     api.NewAnimeClient(),
		cfgManager: mgr,
	}
}

func (m AnimeModel) Init() tea.Cmd {
	return m.fetchTrending()
}

func (m AnimeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case trendingAnimeMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.trendingList = msg.items
		m.cdnURL = msg.cdnURL
		m.index = 0
		m.err = nil
		return m, nil

	case searchAnimeMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.animeList = msg.items
		m.cdnURL = msg.cdnURL
		m.index = 0
		m.err = nil
		return m, nil

	case loadEpisodesMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.movie = msg.movie
		m.servers = msg.servers
		m.serverIdx = 0
		m.episodeList = currentServerEpisodes(msg.servers, 0)
		m.currentEpIdx = 0
		m.index = 0
		m.screen = animeScreenEpisodes
		m.err = nil
		return m, nil

	case playAnimeMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			if m.currentEpIdx < len(m.episodeList) {
				m.statusMsg = fmt.Sprintf("▶ Đang phát tập %s", m.episodeList[m.currentEpIdx].Name)
			}
		}
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// ─── TRENDING SCREEN ───
		if m.screen == animeScreenTrending {
			switch key {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				return m.parent, nil
			case "s":
				m.screen = animeScreenSearch
				m.index = 0
				m.err = nil
				m.statusMsg = ""
			case "up", "k":
				if m.index > 0 {
					m.index--
				}
			case "down", "j":
				if m.index < len(m.trendingList)-1 {
					m.index++
				}
			case "enter":
				if len(m.trendingList) > 0 && m.index < len(m.trendingList) {
					selected := m.trendingList[m.index]
					m.selectedAnime = &selected
					m.loading = true
					m.err = nil
					return m, m.loadEpisodes(selected.Slug)
				}
			}
			return m, nil
		}

		// ─── SEARCH SCREEN: tất cả ký tự vào query ───
		if m.screen == animeScreenSearch {
			switch key {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				// Về trending thay vì về menu
				m.screen = animeScreenTrending
				m.index = 0
				m.err = nil
				m.animeList = nil
				m.searchQuery = ""
				m.statusMsg = ""
			case "up":
				if m.index > 0 {
					m.index--
				}
			case "down":
				if m.index < len(m.animeList)-1 {
					m.index++
				}
			case "backspace":
				if len(m.searchQuery) > 0 {
					runes := []rune(m.searchQuery)
					m.searchQuery = string(runes[:len(runes)-1])
				}
			case "ctrl+u":
				m.searchQuery = ""
				m.animeList = nil
			case "enter":
				if len(m.animeList) > 0 && m.index < len(m.animeList) {
					selected := m.animeList[m.index]
					m.selectedAnime = &selected
					m.loading = true
					m.err = nil
					return m, m.loadEpisodes(selected.Slug)
				}
				if len([]rune(m.searchQuery)) >= 2 {
					m.loading = true
					m.err = nil
					m.animeList = nil
					return m, m.searchAnime(m.searchQuery)
				}
			default:
				runes := []rune(key)
				if len(runes) == 1 && !m.loading {
					m.searchQuery += key
					m.statusMsg = ""
				}
			}
			return m, nil
		}

		// ─── EPISODE SCREEN: hotkeys ───
		switch key {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.screen = animeScreenTrending
			m.selectedAnime = nil
			m.movie = nil
			m.servers = nil
			m.episodeList = nil
			m.index = 0
			m.err = nil
			m.statusMsg = ""
		case "up", "k":
			if m.index > 0 {
				m.index--
			}
		case "down", "j":
			if m.index < len(m.episodeList)-1 {
				m.index++
			}
		case "enter":
			if len(m.episodeList) > 0 && m.index < len(m.episodeList) {
				m.currentEpIdx = m.index
				m.statusMsg = ""
				m.err = nil
				return m, m.playEpisode(m.index)
			}
		case "n":
			if m.currentEpIdx < len(m.episodeList)-1 {
				m.currentEpIdx++
				m.index = m.currentEpIdx
				m.err = nil
				return m, m.playEpisode(m.currentEpIdx)
			}
			m.statusMsg = "⚠ Đây là tập cuối!"
		case "p":
			if m.currentEpIdx > 0 {
				m.currentEpIdx--
				m.index = m.currentEpIdx
				m.err = nil
				return m, m.playEpisode(m.currentEpIdx)
			}
			m.statusMsg = "⚠ Đây là tập đầu!"
		case "s":
			if len(m.servers) > 1 {
				m.serverIdx = (m.serverIdx + 1) % len(m.servers)
				m.episodeList = currentServerEpisodes(m.servers, m.serverIdx)
				m.index = 0
				m.currentEpIdx = 0
				m.statusMsg = fmt.Sprintf("⚡ Server: %s", m.servers[m.serverIdx].ServerName)
			} else {
				m.statusMsg = "Chỉ có 1 server"
			}
		case "f":
			if m.selectedAnime != nil && m.cfgManager != nil {
				item := config.FavoriteItem{
					Name: m.selectedAnime.Name,
					Slug: m.selectedAnime.Slug,
					Type: "anime",
				}
				if err := m.cfgManager.AddToFavorites(item); err == nil {
					m.statusMsg = fmt.Sprintf("⭐ Đã thêm '%s'!", m.selectedAnime.Name)
				}
			}
		}
	}

	return m, nil
}

func currentServerEpisodes(servers []api.EpisodeServer, idx int) []api.EpisodeData {
	if len(servers) == 0 || idx >= len(servers) {
		return nil
	}
	return servers[idx].ServerData
}

// --- Commands ---

func (m AnimeModel) fetchTrending() tea.Cmd {
	return func() tea.Msg {
		items, cdnURL, err := m.client.GetTrending(1)
		return trendingAnimeMsg{items: items, cdnURL: cdnURL, err: err}
	}
}

func (m AnimeModel) searchAnime(query string) tea.Cmd {
	return func() tea.Msg {
		items, cdnURL, err := m.client.Search(query)
		return searchAnimeMsg{items: items, cdnURL: cdnURL, err: err}
	}
}

func (m AnimeModel) loadEpisodes(slug string) tea.Cmd {
	return func() tea.Msg {
		movie, servers, err := m.client.GetAllServers(slug)
		if err != nil {
			return loadEpisodesMsg{err: err}
		}
		bestIdx := 0
		for i, srv := range servers {
			for _, ep := range srv.ServerData {
				if ep.LinkM3U8 != "" {
					bestIdx = i
					goto found
				}
			}
		}
	found:
		var eps []api.EpisodeData
		if len(servers) > bestIdx {
			eps = servers[bestIdx].ServerData
		}
		return loadEpisodesMsg{movie: movie, servers: servers, episodes: eps}
	}
}

func (m AnimeModel) playEpisode(epIdx int) tea.Cmd {
	return func() tea.Msg {
		if epIdx < 0 || epIdx >= len(m.episodeList) {
			return playAnimeMsg{err: fmt.Errorf("tập phim không hợp lệ")}
		}

		ep := m.episodeList[epIdx]
		link := ep.LinkM3U8
		if link == "" {
			link = ep.LinkEmbed
		}
		if link == "" {
			return playAnimeMsg{err: fmt.Errorf("không tìm thấy link phát cho tập này")}
		}

		// Lưu lịch sử trước
		if m.cfgManager != nil && m.selectedAnime != nil {
			item := config.HistoryItem{
				Name:      m.selectedAnime.Name,
				Episode:   "Tập " + ep.Name,
				URL:       link,
				Type:      "anime",
				Timestamp: time.Now().Format("2006-01-02 15:04"),
			}
			_ = m.cfgManager.AddToHistory(item)
		}

		// Tìm player
		var playerName = player.GetDefaultPlayer()
		if m.cfgManager != nil {
			if cfg, err := m.cfgManager.LoadConfig(); err == nil && cfg.Player != "" {
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
		}

		if p == nil || err != nil {
			// Fallback: mở bằng trình duyệt
			browserErr := player.PlayWithBrowser(link)
			if browserErr != nil {
				return playAnimeMsg{err: fmt.Errorf("không tìm thấy player và không mở được trình duyệt: %v", browserErr)}
			}
			return playAnimeMsg{err: nil}
		}

		title := "GoNime"
		if m.selectedAnime != nil {
			title = fmt.Sprintf("GoNime: %s - Tập %s", m.selectedAnime.Name, ep.Name)
		}

		err = p.Play(link, title)
		return playAnimeMsg{err: err}
	}
}

// --- View ---

func (m AnimeModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#33a5f5")).
		Bold(true)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9f75ff")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f0f4ff"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#606880"))

	errStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff6b6b"))

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4caf50"))

	badgeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f5a623")).
		Bold(true)

	switch m.screen {

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// TRENDING SCREEN
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case animeScreenTrending:
		b.WriteString("\n")
		b.WriteString(headerStyle.Render("  📺 ANIME"))
		b.WriteString("  ")
		b.WriteString(badgeStyle.Render("[🔥 Xu hướng]"))
		b.WriteString("\n\n")

		if m.err != nil {
			b.WriteString("  ")
			b.WriteString(errStyle.Render("❌ " + m.err.Error()))
			b.WriteString("\n")
		} else if m.loading || len(m.trendingList) == 0 {
			b.WriteString("  ")
			b.WriteString(hintStyle.Render("⏳ Đang tải xu hướng..."))
		} else {
			for i, anime := range m.trendingList {
				yearStr := ""
				if anime.Year > 0 {
					yearStr = fmt.Sprintf(" (%d)", anime.Year)
				}
				line := anime.Name + yearStr
				if i == m.index {
					b.WriteString(selectedStyle.Render("  ▶ " + line))
				} else {
					b.WriteString(normalStyle.Render("    " + line))
				}
				b.WriteString("\n")
			}
		}

		if m.statusMsg != "" {
			b.WriteString("\n  ")
			b.WriteString(successStyle.Render(m.statusMsg))
		}

		b.WriteString("\n\n")
		b.WriteString(hintStyle.Render("  s: 🔍 Tìm kiếm  │  Enter: Xem  │  ↑↓: Điều hướng  │  Esc: Về menu"))

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// SEARCH SCREEN
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case animeScreenSearch:
		b.WriteString("\n")
		b.WriteString(headerStyle.Render("  📺 ANIME"))
		b.WriteString("  ")
		b.WriteString(badgeStyle.Render("[🔍 Tìm kiếm]"))
		b.WriteString("\n\n")

		cursor := "▌"
		if m.loading {
			cursor = ""
		}

		b.WriteString(drawInputBox("Nhập tên anime:", m.searchQuery+cursor, 56, "#9f75ff"))
		b.WriteString("\n\n")

		if m.err != nil {
			b.WriteString("  ")
			b.WriteString(errStyle.Render("❌ " + m.err.Error()))
			b.WriteString("\n")
		} else if m.loading {
			b.WriteString("  ")
			b.WriteString(hintStyle.Render("⏳ Đang tìm kiếm..."))
		} else if len(m.animeList) == 0 {
			b.WriteString("  ")
			b.WriteString(hintStyle.Render("Gõ tên anime rồi nhấn Enter (≥2 ký tự)"))
		} else {
			b.WriteString("  ")
			b.WriteString(hintStyle.Render(fmt.Sprintf("Tìm thấy %d kết quả:\n", len(m.animeList))))
			b.WriteString("\n")
			for i, anime := range m.animeList {
				yearStr := ""
				if anime.Year > 0 {
					yearStr = fmt.Sprintf(" (%d)", anime.Year)
				}
				line := anime.Name + yearStr
				if i == m.index {
					b.WriteString(selectedStyle.Render("  ▶ " + line))
				} else {
					b.WriteString(normalStyle.Render("    " + line))
				}
				b.WriteString("\n")
			}
		}

		if m.statusMsg != "" {
			b.WriteString("\n  ")
			b.WriteString(successStyle.Render(m.statusMsg))
		}

		b.WriteString("\n\n")
		b.WriteString(hintStyle.Render("  Esc: Về xu hướng  │  ↑↓: Điều hướng  │  Enter: Tìm/Chọn  │  Ctrl+U: Xóa"))

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// EPISODES SCREEN
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case animeScreenEpisodes:
		b.WriteString("\n")

		title := "DANH SÁCH TẬP"
		if m.selectedAnime != nil {
			title = "📺 " + m.selectedAnime.Name
		}
		b.WriteString(headerStyle.Render(title))
		b.WriteString("\n")

		if m.movie != nil {
			infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#606880"))
			parts := []string{}
			if m.movie.Episode != "" {
				parts = append(parts, m.movie.Episode)
			}
			if m.movie.Quality != "" {
				parts = append(parts, m.movie.Quality)
			}
			if m.movie.Lang != "" {
				parts = append(parts, m.movie.Lang)
			}
			if len(parts) > 0 {
				b.WriteString(infoStyle.Render("  " + strings.Join(parts, " · ")))
			}
		}

		if len(m.servers) > 0 {
			serverStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f5a623")).Bold(true)
			currentServer := "?"
			if m.serverIdx < len(m.servers) {
				currentServer = m.servers[m.serverIdx].ServerName
			}
			serverLine := fmt.Sprintf("  ⚡ Server [%d/%d]: %s", m.serverIdx+1, len(m.servers), currentServer)
			b.WriteString("  ")
			b.WriteString(serverStyle.Render(serverLine))
			if len(m.servers) > 1 {
				b.WriteString(hintStyle.Render("  (s: đổi server)"))
			}
		}
		b.WriteString("\n\n")

		if m.loading {
			b.WriteString("  ")
			b.WriteString(hintStyle.Render("⏳ Đang tải danh sách tập..."))
		} else if m.err != nil {
			b.WriteString("  ")
			b.WriteString(errStyle.Render("❌ " + m.err.Error()))
		} else if len(m.episodeList) == 0 {
			b.WriteString("  ")
			b.WriteString(hintStyle.Render("Không có tập nào."))
		} else {
			total := len(m.episodeList)
			start := 0
			if m.index > 10 {
				start = m.index - 10
			}
			end := start + 20
			if end > total {
				end = total
			}

			if total > 30 {
				m.renderEpisodeGrid(&b, start, end, selectedStyle, normalStyle, successStyle)
			} else {
				m.renderEpisodeList(&b, start, end, selectedStyle, normalStyle, successStyle)
			}

			b.WriteString("\n")
			b.WriteString(hintStyle.Render(fmt.Sprintf("  Tổng: %d tập  │  Đang xem: Tập %s",
				total, func() string {
					if m.currentEpIdx < len(m.episodeList) {
						return m.episodeList[m.currentEpIdx].Name
					}
					return "-"
				}())))
		}

		if m.statusMsg != "" {
			b.WriteString("\n  ")
			if strings.HasPrefix(m.statusMsg, "⚠") {
				b.WriteString(errStyle.Render(m.statusMsg))
			} else {
				b.WriteString(successStyle.Render(m.statusMsg))
			}
		}

		b.WriteString("\n\n")
		b.WriteString(hintStyle.Render("  Esc: Quay lại  │  ↑↓: Điều hướng  │  Enter: Phát  │  n/p: Tập sau/trước  │  s: Đổi server  │  f: ⭐"))
	}

	return b.String()
}

func (m *AnimeModel) renderEpisodeList(b *strings.Builder, start, end int, sel, norm, ok lipgloss.Style) {
	for i := start; i < end; i++ {
		ep := m.episodeList[i]
		line := fmt.Sprintf("Tập %s", ep.Name)

		if i == m.currentEpIdx {
			line += "  ▶"
		}

		if i == m.index {
			b.WriteString(sel.Render("  ▶ " + line))
		} else {
			b.WriteString(norm.Render("    " + line))
		}
		b.WriteString("\n")
	}
}

func (m *AnimeModel) renderEpisodeGrid(b *strings.Builder, start, end int, sel, norm, ok lipgloss.Style) {
	cols := 5
	colWidth := 14

	for i := start; i < end; i++ {
		ep := m.episodeList[i]
		label := fmt.Sprintf("Tập %-4s", ep.Name)
		if len(label) > colWidth-2 {
			label = label[:colWidth-2]
		}

		if i == m.index && i == m.currentEpIdx {
			b.WriteString(sel.Render(fmt.Sprintf("▶%-*s", colWidth-1, label+"▶")))
		} else if i == m.index {
			b.WriteString(sel.Render(fmt.Sprintf("▶%-*s", colWidth-1, label)))
		} else if i == m.currentEpIdx {
			b.WriteString(ok.Render(fmt.Sprintf(" %-*s", colWidth-1, label+"▶")))
		} else {
			b.WriteString(norm.Render(fmt.Sprintf(" %-*s", colWidth-1, label)))
		}

		if (i-start+1)%cols == 0 {
			b.WriteString("\n")
		}
	}
	b.WriteString("\n")
}
