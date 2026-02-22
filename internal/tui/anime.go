package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gonime/internal/api"
	"gonime/internal/config"
	"gonime/internal/player"
)

// --- Messages ---

type searchAnimeMsg struct {
	items  []api.AnimeItem
	cdnURL string
	err    error
}

type loadEpisodesMsg struct {
	movie    *api.MovieInfo
	episodes []api.EpisodeData
	err      error
}

type playAnimeMsg struct {
	err error
}

// --- State ---

type animeScreen int

const (
	animeScreenSearch animeScreen = iota
	animeScreenEpisodes
	animeScreenPlaying
)

// --- Model ---

type AnimeModel struct {
	parent        MainModel
	screen        animeScreen
	searchQuery   string
	animeList     []api.AnimeItem
	cdnURL        string
	selectedAnime *api.AnimeItem
	movie         *api.MovieInfo
	episodeList   []api.EpisodeData
	currentEpIdx  int
	index         int
	loading       bool
	err           error
	client        *api.AnimeClient
	cfgManager    *config.Manager
	playerInst    *player.Player
	statusMsg     string
}

func NewAnimeModel(parent MainModel) AnimeModel {
	mgr, _ := config.NewManager()
	return AnimeModel{
		parent:     parent,
		screen:     animeScreenSearch,
		client:     api.NewAnimeClient(),
		cfgManager: mgr,
	}
}

func (m AnimeModel) Init() tea.Cmd {
	return nil
}

func (m AnimeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
		m.episodeList = msg.episodes
		m.currentEpIdx = 0
		m.index = 0
		m.screen = animeScreenEpisodes
		m.err = nil
		return m, nil

	case playAnimeMsg:
		if msg.err != nil {
			m.err = msg.err
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			switch m.screen {
			case animeScreenEpisodes:
				m.screen = animeScreenSearch
				m.selectedAnime = nil
				m.movie = nil
				m.episodeList = nil
				m.index = 0
				m.err = nil
				return m, nil
			default:
				return m.parent, nil
			}

		case "up", "k":
			if m.index > 0 {
				m.index--
			}

		case "down", "j":
			limit := len(m.animeList) - 1
			if m.screen == animeScreenEpisodes {
				limit = len(m.episodeList) - 1
			}
			if m.index < limit {
				m.index++
			}

		case "backspace":
			if m.screen == animeScreenSearch && len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}

		case "enter":
			switch m.screen {
			case animeScreenSearch:
				// Nếu đang hiện danh sách anime → load episodes
				if len(m.animeList) > 0 && m.index < len(m.animeList) {
					selected := m.animeList[m.index]
					m.selectedAnime = &selected
					m.loading = true
					m.err = nil
					return m, m.loadEpisodes(selected.Slug)
				}
				// Nếu chưa có danh sách → search
				if len(m.searchQuery) >= 2 {
					m.loading = true
					m.err = nil
					m.animeList = nil
					return m, m.searchAnime(m.searchQuery)
				}

			case animeScreenEpisodes:
				if len(m.episodeList) > 0 && m.index < len(m.episodeList) {
					m.currentEpIdx = m.index
					return m, m.playEpisode(m.index)
				}
			}

		// Phím 'n' để xem tập tiếp theo (khi đang ở màn episodes)
		case "n":
			if m.screen == animeScreenEpisodes && m.currentEpIdx < len(m.episodeList)-1 {
				m.currentEpIdx++
				m.index = m.currentEpIdx
				return m, m.playEpisode(m.currentEpIdx)
			}

		// Phím 'p' để xem tập trước
		case "p":
			if m.screen == animeScreenEpisodes && m.currentEpIdx > 0 {
				m.currentEpIdx--
				m.index = m.currentEpIdx
				return m, m.playEpisode(m.currentEpIdx)
			}

		// Phím 'f' để thêm vào yêu thích
		case "f":
			if m.selectedAnime != nil && m.cfgManager != nil {
				item := config.FavoriteItem{
					Name: m.selectedAnime.Name,
					Slug: m.selectedAnime.Slug,
					Type: "anime",
				}
				if err := m.cfgManager.AddToFavorites(item); err == nil {
					m.statusMsg = fmt.Sprintf("⭐ Đã thêm '%s' vào yêu thích!", m.selectedAnime.Name)
				}
			}

		// Phím 'd' để xóa search query
		case "d":
			if m.screen == animeScreenSearch {
				m.searchQuery = ""
				m.animeList = nil
			}

		default:
			if m.screen == animeScreenSearch {
				// Gõ ký tự vào ô search
				if len(msg.String()) == 1 && !m.loading {
					m.searchQuery += msg.String()
					m.statusMsg = ""
				}
			}
		}
	}

	return m, nil
}

// --- Commands ---

func (m AnimeModel) searchAnime(query string) tea.Cmd {
	return func() tea.Msg {
		items, cdnURL, err := m.client.Search(query)
		return searchAnimeMsg{items: items, cdnURL: cdnURL, err: err}
	}
}

func (m AnimeModel) loadEpisodes(slug string) tea.Cmd {
	return func() tea.Msg {
		movie, episodes, err := m.client.GetEpisodes(slug)
		return loadEpisodesMsg{movie: movie, episodes: episodes, err: err}
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
			return playAnimeMsg{err: fmt.Errorf("không tìm thấy link phát")}
		}

		// Lấy config để biết dùng player nào
		var playerName = player.GetDefaultPlayer()
		if m.cfgManager != nil {
			cfg, err := m.cfgManager.LoadConfig()
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

		title := "GoNime"
		if m.selectedAnime != nil {
			title = fmt.Sprintf("GoNime: %s - Tập %s", m.selectedAnime.Name, ep.Name)
		}

		// Lưu lịch sử
		if m.cfgManager != nil && m.selectedAnime != nil {
			item := config.HistoryItem{
				Name:      m.selectedAnime.Name,
				Episode:   ep.Name,
				URL:       link,
				Type:      "anime",
				Timestamp: time.Now().Format("2006-01-02 15:04"),
			}
			_ = m.cfgManager.AddToHistory(item)
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

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#9f75ff")).
		Padding(0, 1).
		Width(50)

	switch m.screen {

	// ━━━ MÀNG TÌM KIẾM ━━━
	case animeScreenSearch:
		b.WriteString(headerStyle.Render("📺 TÌM KIẾM ANIME"))
		b.WriteString("\n\n")

		cursor := "▌"
		if m.loading {
			cursor = " ⏳"
		}
		b.WriteString("Nhập tên anime: ")
		b.WriteString(inputStyle.Render(m.searchQuery + cursor))
		b.WriteString("\n\n")

		if m.err != nil {
			b.WriteString(errStyle.Render(fmt.Sprintf("❌ Lỗi: %v", m.err)))
			b.WriteString("\n")
		} else if m.loading {
			b.WriteString(hintStyle.Render("⏳ Đang tìm kiếm..."))
		} else if len(m.animeList) == 0 {
			b.WriteString(hintStyle.Render("Nhập tên anime (≥2 ký tự) → Enter để tìm kiếm"))
		} else {
			b.WriteString(hintStyle.Render(fmt.Sprintf("Tìm thấy %d kết quả:\n", len(m.animeList))))
			for i, anime := range m.animeList {
				yearStr := ""
				if anime.Year > 0 {
					yearStr = fmt.Sprintf(" (%d)", anime.Year)
				}
				line := fmt.Sprintf("%s%s", anime.Name, yearStr)
				if i == m.index {
					b.WriteString(selectedStyle.Render("▶ " + line))
				} else {
					b.WriteString(normalStyle.Render("  " + line))
				}
				b.WriteString("\n")
			}
		}

		if m.statusMsg != "" {
			b.WriteString("\n")
			b.WriteString(successStyle.Render(m.statusMsg))
		}

		b.WriteString("\n")
		b.WriteString(hintStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Tìm/Chọn"))

	// ━━━ MÀNG DANH SÁCH TẬP ━━━
	case animeScreenEpisodes:
		title := "DANH SÁCH TẬP"
		if m.selectedAnime != nil {
			title = fmt.Sprintf("📺 %s", m.selectedAnime.Name)
		}
		b.WriteString(headerStyle.Render(title))

		if m.movie != nil && m.movie.Episode != "" {
			infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
			b.WriteString(infoStyle.Render(fmt.Sprintf("  Trạng thái: %s | Chất lượng: %s | Ngôn ngữ: %s",
				m.movie.Episode, m.movie.Quality, m.movie.Lang)))
		}
		b.WriteString("\n\n")

		if m.loading {
			b.WriteString(hintStyle.Render("⏳ Đang tải danh sách tập..."))
		} else if m.err != nil {
			b.WriteString(errStyle.Render(fmt.Sprintf("❌ Lỗi: %v", m.err)))
		} else if len(m.episodeList) == 0 {
			b.WriteString(hintStyle.Render("Không có tập nào."))
		} else {
			// Hiển thị tối đa 20 tập, cuộn theo index
			start := 0
			if m.index > 10 {
				start = m.index - 10
			}
			end := start + 20
			if end > len(m.episodeList) {
				end = len(m.episodeList)
			}

			for i := start; i < end; i++ {
				ep := m.episodeList[i]
				line := fmt.Sprintf("Tập %s", ep.Name)
				if i == m.currentEpIdx {
					line += " ▶ Đang phát"
				}
				if i == m.index {
					b.WriteString(selectedStyle.Render("▶ " + line))
				} else {
					b.WriteString(normalStyle.Render("  " + line))
				}
				b.WriteString("\n")
			}

			if len(m.episodeList) > 20 {
				b.WriteString(hintStyle.Render(fmt.Sprintf("\n(Tổng: %d tập)\n", len(m.episodeList))))
			}
		}

		if m.statusMsg != "" {
			b.WriteString("\n")
			b.WriteString(successStyle.Render(m.statusMsg))
		}

		b.WriteString("\n")
		b.WriteString(hintStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Phát | n: Tập sau | p: Tập trước | f: Yêu thích"))
	}

	return b.String()
}
