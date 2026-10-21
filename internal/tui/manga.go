package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gonime/internal/config"
	"gonime/internal/fuzzy"
	"gonime/internal/player"
	"gonime/internal/render"
	"gonime/internal/scraper"
	"gonime/internal/ui/i18n"
	"gonime/internal/ui/theme"
	"gonime/internal/viewer"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Messages ---

type searchMangaMsg struct {
	results []scraper.Manga
	err     error
}

type trendingMangaMsg struct {
	results []scraper.Manga
	err     error
}

type mangaDetailMsg struct {
	manga *scraper.Manga
	err   error
}

type chapterImagesMsg struct {
	images      []string
	chapterName string
	chapterURL  string
	err         error
}

type mangaPlayMsg struct {
	err error
}

type chafaRenderMsg struct {
	rendered string
	imageIdx int
	err      error
}

// --- State ---

type mangaScreen int

const (
	mangaScreenTrending mangaScreen = iota
	mangaScreenSearch
	mangaScreenChapters
	mangaScreenReading
)

// --- Model ---

type MangaModel struct {
	parent          MainModel
	screen          mangaScreen
	searchQuery     string
	mangaList       []scraper.Manga // kết quả thô từ API
	fuzzyResults    []scraper.Manga // kết quả sau fuzzy filter + rank
	trendingList    []scraper.Manga
	selectedManga   *scraper.Manga
	chapterList     []scraper.Chapter
	currentChapIdx  int
	images          []string
	currentImageIdx int
	renderedImage   string
	renderLoading   bool
	index           int
	loading         bool
	tick            int // spinner animation frame counter
	err             error
	statusMsg       string
	posterArt       string // chafa-rendered poster for chapters screen
	scraper         scraper.MangaSource
	cfgManager      *config.Manager
	cfg             *config.Config
	cacheDir        string
	termWidth       int
	termHeight      int
	viewer          viewer.ViewerClient
	source          string
	sources         []string
	lastPlayTime    time.Time
}

func NewMangaModel(parent MainModel) MangaModel {
	mgr, _ := config.NewManager()

	cacheDir := ""
	var cfg *config.Config
	if mgr != nil {
		if c, err := mgr.LoadConfig(); err == nil {
			cacheDir = c.MangaImageDir
			cfg = c
		}
	}
	if cacheDir == "" {
		cacheDir = "manga_cache"
	}

	sources := scraper.GetMangaSourceNames()

	return MangaModel{
		parent:     parent,
		screen:     mangaScreenTrending,
		scraper:    scraper.GetMangaSource("truyenvn"),
		cfgManager: mgr,
		cfg:        cfg,
		cacheDir:   cacheDir,
		termWidth:  100,
		termHeight: 40,
		source:     "truyenvn",
		sources:    sources,
		viewer:     viewer.NewWailsViewerReal(),
		loading:    true,
	}
}

func (m MangaModel) Init() tea.Cmd {
	return tea.Batch(m.fetchTrending(), tickCmd())
}

func (m MangaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

	case trendingMangaMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.trendingList = msg.results
		m.index = 0
		m.err = nil
		return m, nil

	case searchMangaMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.mangaList = msg.results
		m.index = 0
		m.err = nil
		// Áp dụng fuzzy filter ngay sau khi nhận kết quả từ API
		m.fuzzyResults = applyMangaFuzzy(m.mangaList, m.searchQuery)
		return m, nil

	case mangaDetailMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		if msg.manga != nil {
			m.selectedManga = msg.manga
			m.chapterList = msg.manga.Chapters
			m.currentChapIdx = 0
			m.index = 0
			m.screen = mangaScreenChapters
			m.err = nil
			m.syncViewerChapters(msg.manga)
			// Trigger poster render if cover URL available
			if msg.manga.Cover != "" {
				m.posterArt = ""
				return m, posterCmd(msg.manga.Cover)
			}
		}
		return m, nil

	case posterRenderMsg:
		if msg.err == nil {
			m.posterArt = msg.rendered
		}
		return m, nil

	case tickMsg:
		if m.loading || m.renderLoading {
			m.tick++
		}
		return m, tickCmd()

	case chapterImagesMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.images = msg.images
		m.currentImageIdx = 0
		m.renderedImage = ""
		m.screen = mangaScreenReading
		m.err = nil
		if len(m.images) > 0 {
			return m, m.renderImageCmd(0)
		}
		return m, nil

	case chafaRenderMsg:
		m.renderLoading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		if msg.imageIdx == m.currentImageIdx {
			m.renderedImage = msg.rendered
			m.err = nil
		}
		return m, nil

	case mangaPlayMsg:
		if msg.err != nil {
			m.err = msg.err
		}
		return m, nil

	case tea.KeyMsg:
		key := msg.String()

		// ─── TRENDING SCREEN ───
		if m.screen == mangaScreenTrending {
			switch key {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				return m.parent, nil
			case "s":
				m.screen = mangaScreenSearch
				m.index = 0
				m.err = nil
				m.statusMsg = ""
			case "ctrl+o":
				currentIdx := 0
				for i, s := range m.sources {
					if s == m.source {
						currentIdx = i
						break
					}
				}
				m.source = m.sources[(currentIdx+1)%len(m.sources)]
				m.scraper = scraper.GetMangaSource(m.source)
				m.syncViewerScraper()
				m.statusMsg = fmt.Sprintf("⚡ %s: %s", i18n.T(m.cfg, "status_source_switched"), m.scraper.Name())
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
					m.selectedManga = &selected
					m.loading = true
					m.err = nil
					return m, m.getMangaDetail(selected.Slug)
				}
			}
			return m, nil
		}

		// ─── SEARCH SCREEN: tất cả ký tự vào query ───
		if m.screen == mangaScreenSearch {
			switch key {
			case "ctrl+c":
				return m, tea.Quit
			case "esc":
				m.screen = mangaScreenTrending
				m.index = 0
				m.err = nil
				m.mangaList = nil
				m.fuzzyResults = nil
				m.searchQuery = ""
				m.statusMsg = ""
			case "ctrl+o":
				// Đổi nguồn bằng Ctrl+o — không xung đột với gõ text
				currentIdx := 0
				for i, s := range m.sources {
					if s == m.source {
						currentIdx = i
						break
					}
				}
				m.source = m.sources[(currentIdx+1)%len(m.sources)]
				m.scraper = scraper.GetMangaSource(m.source)
				m.syncViewerScraper()
				m.statusMsg = fmt.Sprintf("⚡ %s: %s", i18n.T(m.cfg, "status_source_switched"), m.scraper.Name())
				m.loading = true
				m.mangaList = nil
				m.fuzzyResults = nil
				m.index = 0
				return m, m.searchManga(m.searchQuery)
			case "up":
				if m.index > 0 {
					m.index--
				}
			case "down":
				// Điều hướng theo fuzzy results (hoặc raw list nếu chưa có)
				maxIdx := len(m.mangaList) - 1
				if m.searchQuery != "" {
					maxIdx = len(m.fuzzyResults) - 1
				}
				if m.index < maxIdx {
					m.index++
				}
			case "backspace":
				if len(m.searchQuery) > 0 {
					runes := []rune(m.searchQuery)
					m.searchQuery = string(runes[:len(runes)-1])
					// Cập nhật fuzzy filter real-time khi xóa
					m.fuzzyResults = applyMangaFuzzy(m.mangaList, m.searchQuery)
					m.index = 0
				}
			case "ctrl+u":
				m.searchQuery = ""
				m.mangaList = nil
				m.fuzzyResults = nil
				m.index = 0
			case "enter":
				// Ưu tiên fuzzy results nếu có
				displayList := m.mangaList
				if m.searchQuery != "" {
					displayList = m.fuzzyResults
				}
				if len(displayList) > 0 && m.index < len(displayList) {
					selected := displayList[m.index]
					m.selectedManga = &selected
					m.loading = true
					m.err = nil
					return m, m.getMangaDetail(selected.Slug)
				}
				if len([]rune(m.searchQuery)) >= 2 {
					m.loading = true
					m.err = nil
					m.mangaList = nil
					m.fuzzyResults = nil
					m.index = 0
					return m, m.searchManga(m.searchQuery)
				}
			default:
				runes := []rune(key)
				if len(runes) == 1 && !m.loading {
					m.searchQuery += key
					m.statusMsg = ""
					// Cập nhật fuzzy filter real-time khi gõ thêm ký tự
					m.fuzzyResults = applyMangaFuzzy(m.mangaList, m.searchQuery)
					m.index = 0
				}
			}
			return m, nil
		}

		// ─── CHAPTER & READING SCREENS: hotkeys ───
		switch key {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			switch m.screen {
			case mangaScreenReading:
				m.screen = mangaScreenChapters
				m.images = nil
				m.renderedImage = ""
				m.err = nil
			case mangaScreenChapters:
				m.screen = mangaScreenTrending
				m.selectedManga = nil
				m.chapterList = nil
				m.index = 0
				m.err = nil
				if len(m.trendingList) == 0 {
					m.loading = true
					return m, m.fetchTrending()
				}
			}

		case "up", "k":
			switch m.screen {
			case mangaScreenReading:
				if m.currentImageIdx > 0 {
					m.currentImageIdx--
					m.renderedImage = ""
					m.renderLoading = true
					return m, m.renderImageCmd(m.currentImageIdx)
				}
			case mangaScreenChapters:
				if m.index > 0 {
					m.index--
				}
			}

		case "down", "j":
			switch m.screen {
			case mangaScreenReading:
				if m.currentImageIdx < len(m.images)-1 {
					m.currentImageIdx++
					m.renderedImage = ""
					m.renderLoading = true
					return m, m.renderImageCmd(m.currentImageIdx)
				}
			case mangaScreenChapters:
				if m.index < len(m.chapterList)-1 {
					m.index++
				}
			}

		case "pgup", "ctrl+u":
			if m.screen == mangaScreenChapters {
				m.index -= 10
				if m.index < 0 {
					m.index = 0
				}
			}

		case "pgdown", "ctrl+d":
			if m.screen == mangaScreenChapters {
				m.index += 10
				if m.index >= len(m.chapterList) {
					m.index = len(m.chapterList) - 1
				}
			}

		case "g":
			if m.screen == mangaScreenChapters {
				m.index = 0
			}

		case "G":
			if m.screen == mangaScreenChapters && len(m.chapterList) > 0 {
				m.index = len(m.chapterList) - 1
			}

		case "left", "h":
			if m.screen == mangaScreenReading && m.currentImageIdx > 0 {
				m.currentImageIdx--
				m.renderedImage = ""
				m.renderLoading = true
				return m, m.renderImageCmd(m.currentImageIdx)
			}

		case "right", "l":
			if m.screen == mangaScreenReading && m.currentImageIdx < len(m.images)-1 {
				m.currentImageIdx++
				m.renderedImage = ""
				m.renderLoading = true
				return m, m.renderImageCmd(m.currentImageIdx)
			}

		case "enter":
			if time.Since(m.lastPlayTime) < 2*time.Second {
				return m, nil
			}
			switch m.screen {
			case mangaScreenChapters:
				if len(m.chapterList) > 0 && m.index < len(m.chapterList) {
					m.lastPlayTime = time.Now()
					m.currentChapIdx = m.index
					m.loading = true
					m.err = nil
					return m, m.loadChapterImages(m.index)
				}
			case mangaScreenReading:
				if m.currentImageIdx < len(m.images)-1 {
					m.lastPlayTime = time.Now()
					m.currentImageIdx++
					m.renderedImage = ""
					m.renderLoading = true
					return m, m.renderImageCmd(m.currentImageIdx)
				}
			}

		case "n":
			if time.Since(m.lastPlayTime) < 2*time.Second {
				return m, nil
			}
			if m.screen == mangaScreenReading || m.screen == mangaScreenChapters {
				if m.currentChapIdx > 0 {
					m.lastPlayTime = time.Now()
					m.currentChapIdx--
					m.index = m.currentChapIdx
					m.loading = true
					return m, m.loadChapterImages(m.currentChapIdx)
				}
				m.statusMsg = "⚠ " + i18n.T(m.cfg, "manga_msg_last_chap")
			}

		case "p":
			if time.Since(m.lastPlayTime) < 2*time.Second {
				return m, nil
			}
			if m.screen == mangaScreenReading || m.screen == mangaScreenChapters {
				if m.currentChapIdx < len(m.chapterList)-1 {
					m.lastPlayTime = time.Now()
					m.currentChapIdx++
					m.index = m.currentChapIdx
					m.loading = true
					return m, m.loadChapterImages(m.currentChapIdx)
				}
				m.statusMsg = "⚠ " + i18n.T(m.cfg, "manga_msg_first_chap")
			}

		case "r":
			if time.Since(m.lastPlayTime) < 2*time.Second {
				return m, nil
			}
			if m.screen == mangaScreenReading && len(m.images) > 0 {
				m.lastPlayTime = time.Now()
				m.renderedImage = ""
				m.renderLoading = true
				return m, m.renderImageCmd(m.currentImageIdx)
			}

		case "w":
			if time.Since(m.lastPlayTime) < 2*time.Second {
				return m, nil
			}
			if m.screen == mangaScreenReading && m.currentChapIdx < len(m.chapterList) {
				m.lastPlayTime = time.Now()
				chapURL := m.chapterList[m.currentChapIdx].URL
				_ = player.PlayWithBrowser(chapURL)
				m.statusMsg = "🌐 " + i18n.T(m.cfg, "manga_msg_browser_open")
			}

		case "o":
			if time.Since(m.lastPlayTime) < 2*time.Second {
				return m, nil
			}
			if m.screen == mangaScreenReading && len(m.images) > 0 {
				m.lastPlayTime = time.Now()
				return m, m.openInViewer(m.currentImageIdx)
			}

		case "f":
			if m.selectedManga != nil && m.cfgManager != nil {
				item := config.FavoriteItem{
					Name: m.selectedManga.Name,
					Slug: m.selectedManga.Slug,
					Type: "manga",
				}
				if err := m.cfgManager.AddToFavorites(item); err == nil {
					m.statusMsg = fmt.Sprintf("⭐ %s '%s'!", i18n.T(m.cfg, "status_fav_added"), m.selectedManga.Name)
				}
			}
		}

	case tea.MouseMsg:
		if msg.Action != tea.MouseActionPress {
			return m, nil
		}
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			return m.handleMouseWheel(-1)
		case tea.MouseButtonWheelDown:
			return m.handleMouseWheel(1)
		}
	}

	return m, nil
}

func (m MangaModel) handleMouseWheel(delta int) (tea.Model, tea.Cmd) {
	switch m.screen {
	case mangaScreenReading:
		nextIdx := m.currentImageIdx + delta
		if nextIdx < 0 || nextIdx >= len(m.images) {
			return m, nil
		}
		m.currentImageIdx = nextIdx
		m.renderedImage = ""
		m.renderLoading = true
		return m, m.renderImageCmd(m.currentImageIdx)
	case mangaScreenChapters:
		nextIdx := m.index + delta
		if nextIdx < 0 {
			nextIdx = 0
		}
		if nextIdx >= len(m.chapterList) {
			nextIdx = len(m.chapterList) - 1
		}
		if nextIdx >= 0 {
			m.index = nextIdx
		}
		return m, nil
	case mangaScreenTrending:
		nextIdx := m.index + delta
		if nextIdx < 0 {
			nextIdx = 0
		}
		if nextIdx >= len(m.trendingList) {
			nextIdx = len(m.trendingList) - 1
		}
		if nextIdx >= 0 {
			m.index = nextIdx
		}
		return m, nil
	case mangaScreenSearch:
		displayLen := len(m.mangaList)
		if m.searchQuery != "" {
			displayLen = len(m.fuzzyResults)
		}
		if displayLen == 0 {
			return m, nil
		}
		nextIdx := m.index + delta
		if nextIdx < 0 {
			nextIdx = 0
		}
		if nextIdx >= displayLen {
			nextIdx = displayLen - 1
		}
		m.index = nextIdx
		return m, nil
	default:
		return m, nil
	}
}

// --- Commands ---

func (m MangaModel) fetchTrending() tea.Cmd {
	return func() tea.Msg {
		results, err := m.scraper.GetTrending()
		return trendingMangaMsg{results: results, err: err}
	}
}

func (m MangaModel) searchManga(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := m.scraper.Search(query)
		return searchMangaMsg{results: results, err: err}
	}
}

func (m MangaModel) getMangaDetail(slug string) tea.Cmd {
	return func() tea.Msg {
		manga, err := m.scraper.GetDetails(slug)
		return mangaDetailMsg{manga: manga, err: err}
	}
}

func (m MangaModel) loadChapterImages(chapIdx int) tea.Cmd {
	return func() tea.Msg {
		if chapIdx < 0 || chapIdx >= len(m.chapterList) {
			return chapterImagesMsg{err: fmt.Errorf("chapter không hợp lệ")}
		}

		chap := m.chapterList[chapIdx]

		if m.cfgManager != nil && m.selectedManga != nil {
			item := config.HistoryItem{
				Name:      m.selectedManga.Name,
				Episode:   chap.Title,
				URL:       chap.URL,
				Type:      "manga",
				Timestamp: time.Now().Format("2006-01-02 15:04"),
			}
			_ = m.cfgManager.AddToHistory(item)
		}

		images, err := m.scraper.GetChapterImages(chap.URL)
		return chapterImagesMsg{
			images:      images,
			chapterName: chap.Title,
			chapterURL:  chap.URL,
			err:         err,
		}
	}
}

func (m MangaModel) renderImageCmd(imageIdx int) tea.Cmd {
	return func() tea.Msg {
		if imageIdx < 0 || imageIdx >= len(m.images) {
			return chafaRenderMsg{imageIdx: imageIdx, err: fmt.Errorf("index ảnh không hợp lệ")}
		}

		imageURL := m.images[imageIdx]

		w := m.termWidth - 2
		h := m.termHeight - 8
		if w <= 0 {
			w = 80
		}
		if h <= 10 {
			h = 20
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		mode := "auto"
		if m.cfg != nil {
			mode = m.cfg.ImageMode
		}
		renderer, err := render.NewImageRenderer(m.cacheDir, w, h, mode)
		if err != nil {
			return chafaRenderMsg{
				imageIdx: imageIdx,
				rendered: fmt.Sprintf("  [Không có chafa - URL ảnh]\n  %s", imageURL),
				err:      nil,
			}
		}

		rendered, err := renderer.RenderURL(ctx, imageURL)
		return chafaRenderMsg{
			imageIdx: imageIdx,
			rendered: rendered,
			err:      err,
		}
	}
}

// syncViewerChapters pushes the chapter list and active scraper to the web viewer.
func (m MangaModel) syncViewerChapters(manga *scraper.Manga) {
	wv, ok := m.viewer.(*viewer.WailsViewerReal)
	if !ok || manga == nil {
		return
	}
	chapters := make([]viewer.ChapterInfo, 0, len(manga.Chapters))
	for _, ch := range manga.Chapters {
		chapters = append(chapters, viewer.ChapterInfo{Title: ch.Title, URL: ch.URL})
	}
	wv.SetChapters(chapters)
	wv.SetScraper(m.scraper)
}

// syncViewerScraper keeps the web viewer on the same source as the TUI.
func (m MangaModel) syncViewerScraper() {
	if wv, ok := m.viewer.(*viewer.WailsViewerReal); ok {
		wv.SetScraper(m.scraper)
	}
}

func (m MangaModel) openInViewer(imageIdx int) tea.Cmd {
	return func() tea.Msg {
		if imageIdx < 0 || imageIdx >= len(m.images) {
			return mangaPlayMsg{err: fmt.Errorf("không có ảnh để mở")}
		}
		if m.viewer == nil {
			return mangaPlayMsg{err: fmt.Errorf("viewer chưa sẵn sàng")}
		}
		return mangaPlayMsg{err: m.viewer.Open(m.images, imageIdx)}
	}
}

// --- View ---

func (m MangaModel) View() string {
	var b strings.Builder

	switch m.screen {

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// TRENDING
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case mangaScreenTrending:
		badges := []string{"🔥 " + i18n.T(m.cfg, "manga_tab_trending"), m.source}
		b.WriteString(drawHeader("📚", i18n.T(m.cfg, "manga_title"), badges, theme.MangaAccent))
		b.WriteString(drawHeaderLine(theme.MangaAccent))
		b.WriteString("\n\n")

		if m.err != nil {
			b.WriteString(drawErrorLine(m.err.Error()))
		} else if m.loading {
			b.WriteString(drawLoadingLine(i18n.T(m.cfg, "manga_loading_trending"), m.tick))
		} else if len(m.trendingList) == 0 {
			b.WriteString("  ")
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "manga_empty_trending") + " "))
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "manga_empty_trending_hint_1") + " "))
			b.WriteString(theme.HintKey.Render("m"))
			b.WriteString(theme.Muted.Render(" " + i18n.T(m.cfg, "manga_empty_trending_hint_2")))
			b.WriteString("\n")
		} else {
			selStyle := lipgloss.NewStyle().
				Foreground(theme.Secondary).
				Background(theme.BackgroundLight).
				Bold(true).
				Padding(0, 1)
			normStyle := lipgloss.NewStyle().
				Foreground(theme.TextSecondary).
				Padding(0, 1)
			for i, manga := range m.trendingList {
				rank := lipgloss.NewStyle().Foreground(theme.TextDim).Render(
					fmt.Sprintf("%3d.", i+1))
				if i == m.index {
					b.WriteString("  " + rank + " " + selStyle.Render("▶  "+manga.Name))
				} else {
					b.WriteString("  " + rank + " " + normStyle.Render("   "+manga.Name))
				}
				b.WriteString("\n")
			}
		}

		if m.statusMsg != "" {
			b.WriteString(drawSuccessLine(m.statusMsg))
		}

		b.WriteString(drawStatusBar([]string{
			`Ctrl+O: ` + i18n.T(m.cfg, "hint_switch_source"),
			`s: ` + i18n.T(m.cfg, "hint_search"),
			i18n.T(m.cfg, "hint_enter_open"),
			i18n.T(m.cfg, "hint_nav"),
			i18n.T(m.cfg, "hint_esc_back"),
		}))
		b.WriteString("\n")

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// SEARCH
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case mangaScreenSearch:
		badges := []string{"🔍 " + i18n.T(m.cfg, "manga_tab_search"), m.source}
		b.WriteString(drawHeader("📚", i18n.T(m.cfg, "manga_title"), badges, theme.MangaAccent))
		b.WriteString(drawHeaderLine(theme.MangaAccent))
		b.WriteString("\n\n")

		cursor := "▌"
		if m.loading {
			cursor = ""
		}
		b.WriteString(drawInputBox("  "+i18n.T(m.cfg, "manga_search_prompt"), m.searchQuery+cursor, 56,
			string(theme.Secondary)))
		b.WriteString("\n\n")

		if m.err != nil {
			b.WriteString(drawErrorLine(fmt.Sprintf("%v", m.err)))
		} else if m.loading {
			b.WriteString(drawLoadingLine(i18n.T(m.cfg, "manga_loading_search"), m.tick))
		} else if len(m.mangaList) == 0 {
			b.WriteString("\n  ")
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "manga_search_hint_1") + " "))
			b.WriteString(theme.HintKey.Render("Enter"))
			b.WriteString(theme.Muted.Render(" " + i18n.T(m.cfg, "manga_search_hint_2")))
		} else {
			displayList := m.mangaList
			if m.searchQuery != "" {
				displayList = m.fuzzyResults
			}

			// Result counter (plain inline text — bordered badges break layout)
			b.WriteString("\n  ")
			b.WriteString(theme.Muted.Render(fmt.Sprintf("%d %s", len(m.mangaList), i18n.T(m.cfg, "manga_search_results"))))
			if m.searchQuery != "" && len(m.fuzzyResults) < len(m.mangaList) {
				b.WriteString(theme.Muted.Render("  ·  "))
				b.WriteString(lipgloss.NewStyle().Foreground(theme.MangaAccent).Bold(true).Render(
					fmt.Sprintf("~%d %s", len(m.fuzzyResults), i18n.T(m.cfg, "manga_search_fuzzy"))))
			}
			b.WriteString("\n\n")

			selStyle := lipgloss.NewStyle().
				Foreground(theme.Secondary).
				Background(theme.BackgroundLight).
				Bold(true).
				Padding(0, 1)
			normStyle := lipgloss.NewStyle().
				Foreground(theme.TextSecondary).
				Padding(0, 1)

			if len(displayList) == 0 && m.searchQuery != "" {
				b.WriteString("\n  ")
				b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "manga_search_no_match") + " '"))
				b.WriteString(theme.HintKey.Render(m.searchQuery))
				b.WriteString(theme.Muted.Render("'"))
				b.WriteString("\n")
			} else {
				for i, manga := range displayList {
					if i == m.index {
						b.WriteString("  " + selStyle.Render("▶  "+manga.Name))
					} else {
						b.WriteString("  " + normStyle.Render("   "+manga.Name))
					}
					b.WriteString("\n")
				}
			}
		}

		if m.statusMsg != "" {
			b.WriteString(drawSuccessLine(m.statusMsg))
		}

		b.WriteString(drawStatusBar([]string{
			"Esc: " + i18n.T(m.cfg, "manga_tab_trending"),
			i18n.T(m.cfg, "hint_nav"),
			i18n.T(m.cfg, "hint_enter_select_2"),
			`Ctrl+o: ` + i18n.T(m.cfg, "hint_switch_source"),
			`Ctrl+U: ` + i18n.T(m.cfg, "hint_clear"),
		}))
		b.WriteString("\n")

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// CHAPTERS
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case mangaScreenChapters:
		mangaName := i18n.T(m.cfg, "manga_chap_list_title")
		if m.selectedManga != nil {
			mangaName = m.selectedManga.Name
		}
		b.WriteString(drawHeader("📚", mangaName, nil, theme.MangaAccent))

		// Author / Status meta
		if m.selectedManga != nil {
			parts := []string{}
			if m.selectedManga.Author != "" {
				parts = append(parts, "✍  "+m.selectedManga.Author)
			}
			if m.selectedManga.Status != "" {
				parts = append(parts, "📌 "+m.selectedManga.Status)
			}
			b.WriteString(drawMetaLine(parts))
		}

		b.WriteString(drawHeaderLine(theme.MangaAccent))

		// -- Detail panel: poster + manga info --
		if m.selectedManga != nil && m.posterArt != "" {
			accentStyle := lipgloss.NewStyle().Foreground(theme.MangaAccent).Bold(true)
			mutedStyle := theme.Muted
			var info strings.Builder
			info.WriteString(accentStyle.Render(m.selectedManga.Name) + "\n")
			var metaParts []string
			if m.selectedManga.Author != "" {
				metaParts = append(metaParts, " "+m.selectedManga.Author)
			}
			if m.selectedManga.Status != "" {
				metaParts = append(metaParts, " "+m.selectedManga.Status)
			}
			if len(metaParts) > 0 {
				info.WriteString(mutedStyle.Render(strings.Join(metaParts, "  ")) + "\n")
			}
			if len(m.selectedManga.Genres) > 0 {
				genres := strings.Join(m.selectedManga.Genres, ", ")
				if len(genres) > 50 {
					genres = genres[:50] + ""
				}
				info.WriteString(mutedStyle.Render(" "+genres) + "\n")
			}
			if m.selectedManga.Description != "" {
				desc := m.selectedManga.Description
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
			b.WriteString(drawLoadingLine(i18n.T(m.cfg, "manga_loading_chaps"), m.tick))
		} else if m.err != nil {
			b.WriteString(drawErrorLine(fmt.Sprintf("%v", m.err)))
		} else if len(m.chapterList) == 0 {
			b.WriteString("\n  ")
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "manga_empty_chaps")))
		} else {
			total := len(m.chapterList)
			start := 0
			if m.index > 7 {
				start = m.index - 7
			}
			end := start + 15
			if end > total {
				end = total
			}

			selStyle := lipgloss.NewStyle().
				Foreground(theme.Secondary).
				Background(theme.BackgroundLight).
				Bold(true).
				Padding(0, 1)
			normStyle := theme.Chapter.Copy()
			readingStyle := theme.ChapterActive.Copy()

			for i := start; i < end; i++ {
				ch := m.chapterList[i]
				dateStr := ""
				if ch.Date != "" {
					dateStr = lipgloss.NewStyle().Foreground(theme.TextMuted).Render("  " + ch.Date)
				}
				readingTag := ""
				if i == m.currentChapIdx {
					readingTag = "  " + lipgloss.NewStyle().Foreground(theme.Gold).Bold(true).Render("◀ "+i18n.T(m.cfg, "manga_tag_reading"))
				}

				if i == m.index {
					b.WriteString("  " + selStyle.Render("▶  "+ch.Title) + dateStr + readingTag)
				} else if i == m.currentChapIdx {
					b.WriteString("  " + readingStyle.Render("   "+ch.Title) + dateStr + readingTag)
				} else {
					b.WriteString("  " + normStyle.Render("   "+ch.Title) + dateStr)
				}
				b.WriteString("\n")
			}

			if total > 15 {
				b.WriteString("\n  ")
				b.WriteString(theme.Muted.Render(fmt.Sprintf("%s: %d", i18n.T(m.cfg, "manga_chap_total"), total)))
			}
		}

		if m.statusMsg != "" {
			b.WriteString(drawSuccessLine(m.statusMsg))
		}

		b.WriteString(drawStatusBar([]string{
			i18n.T(m.cfg, "hint_esc_back"),
			i18n.T(m.cfg, "hint_nav"),
			i18n.T(m.cfg, "hint_pgup_dn"),
			i18n.T(m.cfg, "hint_start_end"),
			i18n.T(m.cfg, "hint_enter_read"),
			i18n.T(m.cfg, "hint_next_prev"),
			"f: ⭐",
		}))
		b.WriteString("\n")

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// READING
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case mangaScreenReading:
		chapTitle := i18n.T(m.cfg, "manga_reading_title")
		if m.currentChapIdx < len(m.chapterList) {
			chapTitle = m.chapterList[m.currentChapIdx].Title
		}

		total := len(m.images)
		current := m.currentImageIdx + 1

		pageInfo := fmt.Sprintf("[%d/%d]", current, total)
		b.WriteString(drawHeader("📖", chapTitle, []string{pageInfo}, theme.MangaAccent))

		// Page progress bar
		if total > 0 {
			b.WriteString("\n  ")
			b.WriteString(theme.ProgressBar(m.currentImageIdx, total, 40))
			b.WriteString("  ")
			b.WriteString(theme.Muted.Render(fmt.Sprintf("%s %d/%d", i18n.T(m.cfg, "manga_page_prefix"), current, total)))
		}

		b.WriteString(drawHeaderLine(theme.MangaAccent))
		b.WriteString("\n\n")

		if m.loading {
			b.WriteString(drawLoadingLine(i18n.T(m.cfg, "manga_loading_pages"), m.tick))
		} else if m.err != nil {
			b.WriteString(drawErrorLine(fmt.Sprintf("%v", m.err)))
			b.WriteString("\n")
			b.WriteString(drawStatusBar([]string{"r: " + i18n.T(m.cfg, "hint_retry"), i18n.T(m.cfg, "hint_esc_back")}))
		} else if len(m.images) == 0 {
			b.WriteString(drawErrorLine(i18n.T(m.cfg, "manga_empty_pages")))
			b.WriteString(drawStatusBar([]string{i18n.T(m.cfg, "hint_esc_back"), "w: " + i18n.T(m.cfg, "hint_browser")}))
		} else {
			if m.renderLoading {
				b.WriteString(drawLoadingLine(i18n.T(m.cfg, "manga_rendering"), m.tick))
			} else if m.renderedImage != "" {
				b.WriteString(m.renderedImage)
				b.WriteString("\n")
			} else {
				imgURL := m.images[m.currentImageIdx]
				if len(imgURL) > 80 {
					imgURL = imgURL[:77] + "..."
				}
				b.WriteString("  ")
				b.WriteString(theme.Subtitle.Render(fmt.Sprintf("🖼  %s %d: %s", i18n.T(m.cfg, "manga_page_prefix"), current, imgURL)))
				b.WriteString("\n")
			}

			chafaStatus := ""
			if render.IsAvailable() {
				chafaStatus = "  " + lipgloss.NewStyle().Foreground(theme.Mint).Render("[chafa ✓]")
			}
			b.WriteString(chafaStatus)

			b.WriteString(drawStatusBar([]string{
				i18n.T(m.cfg, "hint_nav_page"),
				i18n.T(m.cfg, "hint_wheel_page"),
				i18n.T(m.cfg, "hint_next_page"),
				i18n.T(m.cfg, "hint_np_chapter"),
				"r: " + i18n.T(m.cfg, "hint_retry"),
				i18n.T(m.cfg, "hint_viewer"),
				i18n.T(m.cfg, "hint_browser"),
				"f: ⭐",
			}))
		}

		if m.statusMsg != "" {
			b.WriteString(drawSuccessLine(m.statusMsg))
		}
		b.WriteString("\n")
	}

	return b.String()
}

// applyMangaFuzzy lọc và sắp xếp danh sách manga theo query dùng thuật toán fuzzy.
// Trả về slice mới đã được filter + rank theo điểm số giảm dần.
// Nếu query rỗng, trả về toàn bộ danh sách gốc.
func applyMangaFuzzy(items []scraper.Manga, query string) []scraper.Manga {
	if query == "" || len(items) == 0 {
		return items
	}

	// Chuyển tên thành slice string để dùng fuzzy.FilterAndRank
	names := make([]string, len(items))
	for i, item := range items {
		names[i] = item.Name
	}

	ranked := fuzzy.FilterAndRank(names, query)

	result := make([]scraper.Manga, 0, len(ranked))
	for _, r := range ranked {
		result = append(result, items[r.Index])
	}
	return result
}
