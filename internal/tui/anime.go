package tui

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gonime/internal/api"
	"gonime/internal/config"
	"gonime/internal/fuzzy"
	"gonime/internal/player"
	"gonime/internal/render"
	"gonime/internal/ui/i18n"
	"gonime/internal/ui/theme"

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

// downloadAnimeMsg báo kết quả tải xuống một tập phim.
type downloadAnimeMsg struct {
	err error
}

type posterRenderMsg struct {
	rendered string
	err      error
}

type tickMsg struct{}

func tickCmd() tea.Cmd {
	return tea.Tick(200*time.Millisecond, func(time.Time) tea.Msg {
		return tickMsg{}
	})
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
	animeList     []api.AnimeItem // kết quả thô từ API
	fuzzyResults  []api.AnimeItem // kết quả sau fuzzy filter + rank
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
	tick          int // spinner animation frame counter
	err           error
	cfgManager    *config.Manager
	cfg           *config.Config
	statusMsg     string
	source        string
	sources       []string
	posterArt     string // chafa-rendered poster for episodes screen
	lastPlayTime  time.Time
}

func NewAnimeModel(parent MainModel) AnimeModel {
	mgr, _ := config.NewManager()
	cfg, _ := mgr.LoadConfig()
	sources := api.GetMovieSourceNames()
	return AnimeModel{
		parent:     parent,
		screen:     animeScreenTrending,
		cfgManager: mgr,
		cfg:        cfg,
		source:     "phimapi",
		sources:    sources,
		loading:    true,
	}
}

func (m AnimeModel) Init() tea.Cmd {
	return tea.Batch(m.fetchTrending(), tickCmd())
}

func (m *AnimeModel) getClient() api.MovieSourceWithPage {
	return api.GetMovieSourceWithPage(m.source)
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
		// Áp dụng fuzzy filter ngay sau khi nhận kết quả từ API
		m.fuzzyResults = applyAnimeFuzzy(m.animeList, m.searchQuery)
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
		// Trigger poster render if available
		var posterURL string
		if msg.movie != nil {
			posterURL = msg.movie.ThumbURL
			if posterURL == "" {
				posterURL = msg.movie.PosterURL
			}
		}
		if posterURL != "" {
			m.posterArt = ""
			return m, posterCmd(posterURL)
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

	case playAnimeMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			if m.currentEpIdx < len(m.episodeList) {
				m.statusMsg = fmt.Sprintf("%s %s", i18n.T(m.cfg, "anime_playing_ep"), m.episodeList[m.currentEpIdx].Name)
			}
		}
		return m, nil

	case downloadAnimeMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.statusMsg = ""
		} else if m.currentEpIdx < len(m.episodeList) {
			m.err = nil
			m.statusMsg = fmt.Sprintf("%s %s", i18n.T(m.cfg, "anime_dl_sent"), episodeLabel(m.episodeList[m.currentEpIdx].Name, m.cfg))
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
			case "m":
				// Cycle through movie sources
				currentIdx := 0
				for i, s := range m.sources {
					if s == m.source {
						currentIdx = i
						break
					}
				}
				m.source = m.sources[(currentIdx+1)%len(m.sources)]
				m.statusMsg = fmt.Sprintf("⚡ %s: %s", i18n.T(m.cfg, "status_source_switched"), m.source)
				m.loading = true
				m.trendingList = nil
				return m, m.fetchTrending()
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
				m.fuzzyResults = nil
				m.searchQuery = ""
				m.statusMsg = ""
			case "m":
				// Cycle through movie sources
				currentIdx := 0
				for i, s := range m.sources {
					if s == m.source {
						currentIdx = i
						break
					}
				}
				m.source = m.sources[(currentIdx+1)%len(m.sources)]
				m.statusMsg = fmt.Sprintf("⚡ %s: %s", i18n.T(m.cfg, "status_source_switched"), m.source)
			case "up":
				if m.index > 0 {
					m.index--
				}
			case "down":
				// Điều hướng theo fuzzy results (hoặc raw list nếu chưa có)
				maxIdx := len(m.animeList) - 1
				if m.searchQuery != "" {
					maxIdx = len(m.fuzzyResults) - 1
				}
				if m.index < maxIdx {
					m.index++
				}
			case "pgup":
				if m.index > 0 {
					m.index -= 10
					if m.index < 0 {
						m.index = 0
					}
				}
			case "pgdown":
				maxIdx := len(m.animeList) - 1
				if m.searchQuery != "" {
					maxIdx = len(m.fuzzyResults) - 1
				}
				if maxIdx >= 0 {
					m.index += 10
					if m.index > maxIdx {
						m.index = maxIdx
					}
				}
			case "backspace":
				if len(m.searchQuery) > 0 {
					runes := []rune(m.searchQuery)
					m.searchQuery = string(runes[:len(runes)-1])
					// Cập nhật fuzzy filter real-time khi xóa
					m.fuzzyResults = applyAnimeFuzzy(m.animeList, m.searchQuery)
					m.index = 0
				}
			case "ctrl+u":
				m.searchQuery = ""
				m.animeList = nil
				m.fuzzyResults = nil
				m.index = 0
			case "enter":
				// Ưu tiên fuzzy results nếu có
				displayList := m.animeList
				if m.searchQuery != "" {
					displayList = m.fuzzyResults
				}
				if len(displayList) > 0 && m.index < len(displayList) {
					selected := displayList[m.index]
					m.selectedAnime = &selected
					m.loading = true
					m.err = nil
					return m, m.loadEpisodes(selected.Slug)
				}
				if len([]rune(m.searchQuery)) >= 2 {
					m.loading = true
					m.err = nil
					m.animeList = nil
					m.fuzzyResults = nil
					m.index = 0
					return m, m.searchAnime(m.searchQuery)
				}
			default:
				runes := []rune(key)
				if len(runes) == 1 && !m.loading {
					m.searchQuery += key
					m.statusMsg = ""
					// Cập nhật fuzzy filter real-time khi gõ thêm ký tự
					m.fuzzyResults = applyAnimeFuzzy(m.animeList, m.searchQuery)
					m.index = 0
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
			if len(m.trendingList) == 0 {
				m.loading = true
				return m, m.fetchTrending()
			}
		case "up", "k":
			if m.index > 0 {
				m.index--
			}
		case "down", "j":
			if m.index < len(m.episodeList)-1 {
				m.index++
			}
		case "pgup", "ctrl+u":
			m.index -= 10
			if m.index < 0 {
				m.index = 0
			}
		case "pgdown", "ctrl+d":
			m.index += 10
			if m.index >= len(m.episodeList) {
				m.index = len(m.episodeList) - 1
			}
		case "g":
			m.index = 0
		case "G":
			if len(m.episodeList) > 0 {
				m.index = len(m.episodeList) - 1
			}
		case "enter":
			if time.Since(m.lastPlayTime) < 2*time.Second {
				return m, nil
			}
			if len(m.episodeList) > 0 && m.index < len(m.episodeList) {
				m.lastPlayTime = time.Now()
				m.currentEpIdx = m.index
				m.statusMsg = ""
				m.err = nil
				return m, m.playEpisode(m.index)
			}
		case "n":
			if time.Since(m.lastPlayTime) < 2*time.Second {
				return m, nil
			}
			if m.currentEpIdx < len(m.episodeList)-1 {
				m.lastPlayTime = time.Now()
				m.currentEpIdx++
				m.index = m.currentEpIdx
				m.err = nil
				return m, m.playEpisode(m.currentEpIdx)
			}
			m.statusMsg = "⚠ " + i18n.T(m.cfg, "anime_msg_last_ep")
		case "p":
			if time.Since(m.lastPlayTime) < 2*time.Second {
				return m, nil
			}
			if m.currentEpIdx > 0 {
				m.lastPlayTime = time.Now()
				m.currentEpIdx--
				m.index = m.currentEpIdx
				m.err = nil
				return m, m.playEpisode(m.currentEpIdx)
			}
			m.statusMsg = "⚠ " + i18n.T(m.cfg, "anime_msg_first_ep")
		case "d":
			// Tải tập hiện tại bằng yt-dlp
			if len(m.episodeList) > 0 && m.currentEpIdx < len(m.episodeList) {
				m.loading = true
				m.err = nil
				return m, m.downloadEpisode(m.currentEpIdx)
			}
		case "s":
			if len(m.servers) > 1 {
				m.serverIdx = (m.serverIdx + 1) % len(m.servers)
				m.episodeList = currentServerEpisodes(m.servers, m.serverIdx)
				m.index = 0
				m.currentEpIdx = 0
				m.statusMsg = fmt.Sprintf("⚡ Server: %s", m.servers[m.serverIdx].ServerName)
			} else {
				m.statusMsg = i18n.T(m.cfg, "anime_msg_one_server")
			}
		case "f":
			if m.selectedAnime != nil && m.cfgManager != nil {
				item := config.FavoriteItem{
					Name: m.selectedAnime.Name,
					Slug: m.selectedAnime.Slug,
					Type: "anime",
				}
				if err := m.cfgManager.AddToFavorites(item); err == nil {
					m.statusMsg = fmt.Sprintf("⭐ %s '%s'!", i18n.T(m.cfg, "status_fav_added"), m.selectedAnime.Name)
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
	client := m.getClient()
	return func() tea.Msg {
		items, cdnURL, err := client.GetTrending(1)
		return trendingAnimeMsg{items: items, cdnURL: cdnURL, err: err}
	}
}

func (m AnimeModel) searchAnime(query string) tea.Cmd {
	client := m.getClient()
	return func() tea.Msg {
		items, cdnURL, err := client.Search(query)
		return searchAnimeMsg{items: items, cdnURL: cdnURL, err: err}
	}
}

func (m AnimeModel) loadEpisodes(slug string) tea.Cmd {
	client := m.getClient()
	return func() tea.Msg {
		movie, episodes, err := client.GetEpisodes(slug)
		if err != nil {
			return loadEpisodesMsg{err: err}
		}
		var servers []api.EpisodeServer
		if len(episodes) > 0 {
			servers = []api.EpisodeServer{
				{ServerName: "Default", ServerData: episodes},
			}
		}
		return loadEpisodesMsg{movie: movie, servers: servers, episodes: episodes}
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
			thumbURL := ""
			if m.movie != nil {
				thumbURL = m.movie.ThumbURL
				if thumbURL == "" {
					thumbURL = m.movie.PosterURL
				}
			}
			item := config.HistoryItem{
				Name:      m.selectedAnime.Name,
				Episode:   episodeLabel(ep.Name, m.cfg),
				URL:       link,
				ThumbURL:  thumbURL,
				Type:      "anime",
				Timestamp: time.Now().Format("2006-01-02 15:04"),
			}
			_ = m.cfgManager.AddToHistory(item)
		}

		// Tìm player & config
		var cfg *config.Config
		var playerName = player.GetDefaultPlayer()
		if m.cfgManager != nil {
			if c, err := m.cfgManager.LoadConfig(); err == nil {
				cfg = c
				if c.Player != "" {
					playerName = c.Player
				}
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
			title = fmt.Sprintf("GoNime: %s - %s", m.selectedAnime.Name, episodeLabel(ep.Name, m.cfg))
		}

		// Chuẩn bị phụ đề & tốc độ
		var subPath string
		var speed float64 = 1.0
		if cfg != nil {
			if cfg.SubtitlesDir != "" && m.selectedAnime != nil {
				subPath = findLocalSubtitle(cfg.SubtitlesDir, m.selectedAnime.Slug, ep.Name)
			}
			if cfg.PlaybackSpeed > 0 {
				speed = cfg.PlaybackSpeed
			}
		}

		// Gọi player với args tuỳ biến
		if err := playWithOptions(p, link, title, subPath, speed); err != nil {
			return playAnimeMsg{err: err}
		}
		return playAnimeMsg{err: nil}
	}
}

// downloadEpisode tải tập hiện tại bằng yt-dlp vào DownloadDir trong config.
func (m AnimeModel) downloadEpisode(epIdx int) tea.Cmd {
	return func() tea.Msg {
		if epIdx < 0 || epIdx >= len(m.episodeList) {
			return downloadAnimeMsg{err: fmt.Errorf("tập phim không hợp lệ")}
		}

		ep := m.episodeList[epIdx]
		link := ep.LinkM3U8
		if link == "" {
			link = ep.LinkEmbed
		}
		if link == "" {
			return downloadAnimeMsg{err: fmt.Errorf("không tìm thấy link để tải tập này")}
		}

		if m.cfgManager == nil {
			return downloadAnimeMsg{err: fmt.Errorf("không tìm thấy cấu hình tải xuống")}
		}

		cfg, err := m.cfgManager.LoadConfig()
		if err != nil {
			return downloadAnimeMsg{err: fmt.Errorf("không tải được config: %v", err)}
		}

		downloader, err := player.NewDownloader()
		if err != nil {
			return downloadAnimeMsg{err: fmt.Errorf("yt-dlp chưa được cài đặt hoặc không tìm thấy trong PATH: %v", err)}
		}

		outputDir := cfg.DownloadDir
		if outputDir == "" {
			return downloadAnimeMsg{err: fmt.Errorf("DownloadDir không được cấu hình")}
		}

		// Gọi Download, để yt-dlp tự đặt tên file theo tiêu đề.
		go func() {
			_ = os.MkdirAll(outputDir, 0o755)
			_ = downloader.Download(link, outputDir, nil)
		}()

		return downloadAnimeMsg{err: nil}
	}
}

// posterSizedCmd downloads a poster image URL and renders it with chafa with a specific size
func posterSizedCmd(imageURL, size string) tea.Cmd {
	return func() tea.Msg {
		if imageURL == "" {
			return posterRenderMsg{}
		}
		// Download image to temp file
		resp, err := http.Get(imageURL) //nolint:noctx
		if err != nil {
			return posterRenderMsg{err: err}
		}
		defer resp.Body.Close()

		tmpFile, err := os.CreateTemp("", "gonime-poster-*.jpg")
		if err != nil {
			return posterRenderMsg{err: err}
		}
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		if _, err := io.Copy(tmpFile, resp.Body); err != nil {
			tmpFile.Close()
			return posterRenderMsg{err: err}
		}
		tmpFile.Close()

		// Cleanup extension for chafa
		newPath := filepath.Join(filepath.Dir(tmpPath), filepath.Base(tmpPath)+".jpg")
		_ = os.Rename(tmpPath, newPath)

		w, h := 28, 16
		_, _ = fmt.Sscanf(size, "%dx%d", &w, &h)
		mode := "auto"
		if mgr, errCfg := config.NewManager(); errCfg == nil {
			if cfg, errLoad := mgr.LoadConfig(); errLoad == nil && cfg != nil {
				mode = cfg.ImageMode
			}
		}
		renderer, err := render.NewImageRenderer(os.TempDir(), w, h, mode)
		if err != nil {
			return posterRenderMsg{rendered: ""}
		}
		out, err := renderer.RenderFile(newPath)
		_ = os.Remove(newPath)
		if err != nil {
			return posterRenderMsg{err: err}
		}
		return posterRenderMsg{rendered: out}
	}
}

// posterCmd downloads a poster image URL and renders it with chafa for detail views
func posterCmd(imageURL string) tea.Cmd {
	return posterSizedCmd(imageURL, "28x16")
}

// thumbnailCmd downloads a poster image URL and renders it with chafa for list views
func thumbnailCmd(imageURL string) tea.Cmd {
	return posterSizedCmd(imageURL, "14x10")
}

// --- View ---

func (m AnimeModel) View() string {
	var b strings.Builder

	switch m.screen {

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// TRENDING SCREEN
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case animeScreenTrending:
		badges := []string{"🔥 " + i18n.T(m.cfg, "anime_tab_trending"), m.source}
		b.WriteString(drawHeader("📺", i18n.T(m.cfg, "anime_title"), badges, theme.AnimeAccent))
		b.WriteString(drawHeaderLine(theme.AnimeAccent))
		b.WriteString("\n\n")

		if m.err != nil {
			b.WriteString(drawErrorLine(m.err.Error()))
		} else if m.loading {
			b.WriteString(drawLoadingLine(i18n.T(m.cfg, "anime_loading_trending"), m.tick))
		} else if len(m.trendingList) == 0 {
			b.WriteString("  ")
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "anime_empty_trending") + " "))
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "anime_empty_trending_hint_1") + " "))
			b.WriteString(theme.HintKey.Render("m"))
			b.WriteString(theme.Muted.Render(" " + i18n.T(m.cfg, "anime_empty_trending_hint_2")))
			b.WriteString("\n")
		} else {
			selStyle := lipgloss.NewStyle().
				Foreground(theme.Primary).
				Background(theme.BackgroundLight).
				Bold(true).
				Padding(0, 1)
			normStyle := lipgloss.NewStyle().
				Foreground(theme.TextSecondary).
				Padding(0, 1)

			for i, anime := range m.trendingList {
				yearStr := ""
				if anime.Year > 0 {
					yearStr = lipgloss.NewStyle().Foreground(theme.TextMuted).Render(
						fmt.Sprintf(" (%d)", anime.Year))
				}
				rank := lipgloss.NewStyle().Foreground(theme.TextDim).Render(
					fmt.Sprintf("%3d.", i+1))
				if i == m.index {
					b.WriteString("  " + rank + " " + selStyle.Render("▶  "+anime.Name) + yearStr)
				} else {
					b.WriteString("  " + rank + " " + normStyle.Render("   "+anime.Name) + yearStr)
				}
				b.WriteString("\n")
			}
		}

		if m.statusMsg != "" {
			b.WriteString(drawSuccessLine(m.statusMsg))
		}

		b.WriteString(drawStatusBar([]string{
			`m: ` + i18n.T(m.cfg, "hint_switch_source"),
			`s: ` + i18n.T(m.cfg, "hint_search"),
			i18n.T(m.cfg, "hint_enter_open"),
			i18n.T(m.cfg, "hint_nav"),
			i18n.T(m.cfg, "hint_esc_back"),
		}))
		b.WriteString("\n")

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// SEARCH SCREEN
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case animeScreenSearch:
		badges := []string{"🔍 " + i18n.T(m.cfg, "anime_tab_search"), m.source}
		b.WriteString(drawHeader("📺", i18n.T(m.cfg, "anime_title"), badges, theme.AnimeAccent))
		b.WriteString(drawHeaderLine(theme.AnimeAccent))
		b.WriteString("\n\n")

		cursor := "▌"
		if m.loading {
			cursor = ""
		}
		b.WriteString(drawInputBox("  "+i18n.T(m.cfg, "anime_search_prompt"), m.searchQuery+cursor, 56,
			string(theme.Secondary)))
		b.WriteString("\n\n")

		if m.err != nil {
			b.WriteString(drawErrorLine(m.err.Error()))
		} else if m.loading {
			b.WriteString(drawLoadingLine(i18n.T(m.cfg, "anime_loading_search"), m.tick))
		} else if len(m.animeList) == 0 {
			b.WriteString("\n  ")
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "anime_search_hint_1") + " "))
			b.WriteString(theme.HintKey.Render("Enter"))
			b.WriteString(theme.Muted.Render(" " + i18n.T(m.cfg, "anime_search_hint_2")))
		} else {
			displayList := m.animeList
			if m.searchQuery != "" {
				displayList = m.fuzzyResults // filłered; may be nil = no match
			}

			// Result counter (plain inline text — bordered badges break layout)
			b.WriteString("\n  ")
			b.WriteString(theme.Muted.Render(fmt.Sprintf("%d %s", len(m.animeList), i18n.T(m.cfg, "anime_search_results"))))
			if m.searchQuery != "" && len(m.fuzzyResults) < len(m.animeList) {
				b.WriteString(theme.Muted.Render("  ·  "))
				b.WriteString(theme.HintKey.Render(fmt.Sprintf("~%d %s", len(m.fuzzyResults), i18n.T(m.cfg, "anime_search_fuzzy"))))
			}
			b.WriteString("\n\n")

			selStyle := lipgloss.NewStyle().
				Foreground(theme.Primary).
				Background(theme.BackgroundLight).
				Bold(true).
				Padding(0, 1)
			normStyle := lipgloss.NewStyle().
				Foreground(theme.TextSecondary).
				Padding(0, 1)

			if len(displayList) == 0 && m.searchQuery != "" {
				b.WriteString("\n  ")
				b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "anime_search_no_match") + " '"))
				b.WriteString(theme.HintKey.Render(m.searchQuery))
				b.WriteString(theme.Muted.Render("'"))
				b.WriteString("\n")
			} else {
				maxRows := 14
				if m.parent.height > 0 {
					maxRows = m.parent.height - 18
					if maxRows < 6 {
						maxRows = 6
					}
					if maxRows > 24 {
						maxRows = 24
					}
				}
				start, end := windowedRange(len(displayList), m.index, maxRows)
				for i, anime := range displayList {
					if i < start || i >= end {
						continue
					}
					yearStr := ""
					if anime.Year > 0 {
						yearStr = lipgloss.NewStyle().Foreground(theme.TextMuted).Render(
							fmt.Sprintf(" (%d)", anime.Year))
					}
					related := ""
					if anime.OriginName != "" && !strings.EqualFold(anime.OriginName, anime.Name) {
						related = lipgloss.NewStyle().Foreground(theme.TextMuted).Render("  ·  " + truncateString(anime.OriginName, 28))
					}
					if i == m.index {
						b.WriteString("  " + selStyle.Render("▶  "+anime.Name) + yearStr + related)
					} else {
						b.WriteString("  " + normStyle.Render("   "+anime.Name) + yearStr + related)
					}
					b.WriteString("\n")
				}
				if len(displayList) > (end - start) {
					b.WriteString("\n  ")
					b.WriteString(theme.Muted.Render(fmt.Sprintf("Showing %d-%d / %d", start+1, end, len(displayList))))
				}
			}
		}

		if m.statusMsg != "" {
			b.WriteString(drawSuccessLine(m.statusMsg))
		}

		b.WriteString(drawStatusBar([]string{
			"Esc: " + i18n.T(m.cfg, "anime_tab_trending"), // Actually wait, better to use standard hint
			i18n.T(m.cfg, "hint_nav"),
			i18n.T(m.cfg, "hint_pgup_dn"),
			i18n.T(m.cfg, "hint_enter_select_2"),
			`m: ` + i18n.T(m.cfg, "hint_switch_source"),
			`Ctrl+U: ` + i18n.T(m.cfg, "hint_clear"),
		}))
		b.WriteString("\n")

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// EPISODES SCREEN
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case animeScreenEpisodes:
		animeName := i18n.T(m.cfg, "anime_ep_list_title")
		if m.selectedAnime != nil {
			animeName = m.selectedAnime.Name
		}
		b.WriteString(drawHeader("📺", animeName, nil, theme.AnimeAccent))

		// Movie meta info (quality, lang, episodes)
		if m.movie != nil {
			parts := []string{}
			if m.movie.Episode != "" {
				parts = append(parts, "📽  "+m.movie.Episode)
			}
			if m.movie.Quality != "" {
				parts = append(parts, "🎬 "+m.movie.Quality)
			}
			if m.movie.Lang != "" {
				parts = append(parts, "🗣 "+m.movie.Lang)
			}
			b.WriteString(drawMetaLine(parts))
		}

		// Server info — plain text, no bordered badge
		if len(m.servers) > 0 {
			currentServer := "?"
			if m.serverIdx < len(m.servers) {
				currentServer = m.servers[m.serverIdx].ServerName
			}
			serverStyle := lipgloss.NewStyle().Foreground(theme.Gold).Bold(true)
			b.WriteString("\n  ")
			b.WriteString(serverStyle.Render(
				fmt.Sprintf("⚡ Server %d/%d: %s", m.serverIdx+1, len(m.servers), currentServer)))
			if len(m.servers) > 1 {
				b.WriteString("  ")
				b.WriteString(theme.Muted.Render("(s: đổi server)"))
			}
		}

		b.WriteString(drawHeaderLine(theme.AnimeAccent))

		// -- Detail panel: poster + movie info --
		if m.movie != nil && m.posterArt != "" {
			accentStyle := lipgloss.NewStyle().Foreground(theme.AnimeAccent).Bold(true)
			mutedStyle := theme.Muted
			var info strings.Builder
			info.WriteString(accentStyle.Render(m.movie.Name) + "\n")
			if m.movie.OriginName != "" {
				info.WriteString(mutedStyle.Render(m.movie.OriginName) + "\n")
			}
			var metaParts []string
			if m.movie.Year > 0 {
				metaParts = append(metaParts, fmt.Sprintf("%d", m.movie.Year))
			}
			if m.movie.Time != "" {
				metaParts = append(metaParts, m.movie.Time)
			}
			if m.movie.Quality != "" {
				metaParts = append(metaParts, m.movie.Quality)
			}
			if m.movie.Lang != "" {
				metaParts = append(metaParts, m.movie.Lang)
			}
			if len(metaParts) > 0 {
				info.WriteString(mutedStyle.Render(strings.Join(metaParts, "  ")) + "\n")
			}
			if m.movie.Description != "" {
				desc := m.movie.Description
				if len([]rune(desc)) > 200 {
					desc = string([]rune(desc)[:200]) + ""
				}
				info.WriteString("\n" + mutedStyle.Render(wordWrap(desc, 50)))
			}
			posterStyle := lipgloss.NewStyle().Padding(0, 1)
			infoStyle := lipgloss.NewStyle().Padding(0, 1).Width(56)
			panel := lipgloss.JoinHorizontal(lipgloss.Top,
				posterStyle.Render(m.posterArt),
				infoStyle.Render(info.String()),
			)
			b.WriteString("\n")
			b.WriteString(panel)
			b.WriteString("\n")
		} else {
			b.WriteString("\n\n")
		}

		if m.loading {
			b.WriteString(drawLoadingLine(i18n.T(m.cfg, "anime_loading_eps"), m.tick))
		} else if m.err != nil {
			b.WriteString(drawErrorLine(m.err.Error()))
		} else if len(m.episodeList) == 0 {
			b.WriteString("\n  ")
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "anime_empty_eps")))
		} else {
			total := len(m.episodeList)

			// Progress bar
			b.WriteString("\n  ")
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "anime_ep_progress") + "  "))
			b.WriteString(theme.ProgressBar(m.currentEpIdx, total, 32))
			b.WriteString("  ")
			b.WriteString(theme.Muted.Render(fmt.Sprintf("%d/%d", m.currentEpIdx+1, total)))
			b.WriteString("\n\n")

			start := 0
			if m.index > 10 {
				start = m.index - 10
			}
			end := start + 20
			if end > total {
				end = total
			}

			selStyle := theme.EpisodeActive.Copy()
			normStyle := theme.Episode.Copy()
			okStyle := theme.EpisodePlaying.Copy()

			if total > 30 {
				m.renderEpisodeGrid(&b, start, end, selStyle, normStyle, okStyle)
			} else {
				m.renderEpisodeList(&b, start, end, selStyle, normStyle, okStyle)
			}

			b.WriteString("\n")
			b.WriteString(theme.Muted.Render(fmt.Sprintf(
				"  %s: %d %s  ·  %s: %s",
				i18n.T(m.cfg, "anime_ep_total"), total, i18n.T(m.cfg, "anime_ep_unit"), i18n.T(m.cfg, "anime_ep_watching"), func() string {
					if m.currentEpIdx < len(m.episodeList) {
						return episodeLabel(m.episodeList[m.currentEpIdx].Name, m.cfg)
					}
					return "-"
				}())))
		}

		if m.statusMsg != "" {
			b.WriteString("\n")
			if strings.HasPrefix(m.statusMsg, "⚠") {
				b.WriteString(drawErrorLine(m.statusMsg))
			} else {
				b.WriteString(drawSuccessLine(m.statusMsg))
			}
		}

		b.WriteString(drawStatusBar([]string{
			i18n.T(m.cfg, "hint_esc_back"),
			i18n.T(m.cfg, "hint_nav"),
			i18n.T(m.cfg, "hint_pgup_dn"),
			i18n.T(m.cfg, "hint_start_end"),
			i18n.T(m.cfg, "hint_enter_play"),
			i18n.T(m.cfg, "hint_next_prev"),
			`s: Server`,
			`f: ⭐`,
		}))
		b.WriteString("\n")
	}

	return b.String()

}

func (m *AnimeModel) renderEpisodeList(b *strings.Builder, start, end int, sel, norm, ok lipgloss.Style) {
	for i := start; i < end; i++ {
		ep := m.episodeList[i]
		line := episodeLabel(ep.Name, m.cfg)

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
		label := episodeLabel(ep.Name, m.cfg)
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

// episodeLabel chuẩn hóa nhãn tập, tránh lặp "Tập Tập ..."
// Nếu tên tập đã bắt đầu bằng "Tập ", "Tap " hoặc "Episode " (không phân biệt hoa/thường),
// trả về nguyên văn, ngược lại thêm prefix "Tập ".
func episodeLabel(name string, cfg *config.Config) string {
	n := strings.TrimSpace(name)
	if n == "" {
		return ""
	}
	lower := strings.ToLower(n)
	if strings.HasPrefix(lower, "tập ") ||
		strings.HasPrefix(lower, "tap ") ||
		strings.HasPrefix(lower, "episode ") {
		return n
	}
	return i18n.T(cfg, "anime_ep_prefix") + " " + n
}

// applyAnimeFuzzy lọc và sắp xếp danh sách anime theo query dùng thuật toán fuzzy.
// Trả về slice mới đã được filter + rank theo điểm số giảm dần.
// Nếu query rỗng, trả về toàn bộ danh sách gốc.
func applyAnimeFuzzy(items []api.AnimeItem, query string) []api.AnimeItem {
	if query == "" || len(items) == 0 {
		return items
	}

	// Fuzzy theo cả tên chính và tên gốc/related title nếu có.
	names := make([]string, len(items))
	for i, item := range items {
		candidate := item.Name
		if item.OriginName != "" && !strings.EqualFold(item.OriginName, item.Name) {
			candidate += " " + item.OriginName
		}
		names[i] = candidate
	}

	ranked := fuzzy.FilterAndRank(names, query)

	result := make([]api.AnimeItem, 0, len(ranked))
	for _, r := range ranked {
		result = append(result, items[r.Index])
	}
	return result
}
