package tui

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gonime/internal/config"
	"gonime/internal/render"
	"gonime/internal/scraper"
)

// --- Messages ---

type searchMangaMsg struct {
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

// chafaRenderMsg mang kết quả render ảnh từ chafa
type chafaRenderMsg struct {
	rendered string
	imageIdx int
	err      error
}

// --- State ---

type mangaScreen int

const (
	mangaScreenSearch   mangaScreen = iota
	mangaScreenChapters
	mangaScreenReading
)

// --- Model ---

type MangaModel struct {
	parent          MainModel
	screen          mangaScreen
	searchQuery     string
	mangaList       []scraper.Manga
	selectedManga   *scraper.Manga
	chapterList     []scraper.Chapter
	currentChapIdx  int
	images          []string
	currentImageIdx int
	renderedImage   string // output ANSI từ chafa
	renderLoading   bool
	index           int
	loading         bool
	err             error
	statusMsg       string
	scraper         *scraper.MangaScraper
	cfgManager      *config.Manager
	cacheDir        string
	termWidth       int
	termHeight      int
}

func NewMangaModel(parent MainModel) MangaModel {
	mgr, _ := config.NewManager()

	// Xác định cache dir
	cacheDir := ""
	if mgr != nil {
		if cfg, err := mgr.LoadConfig(); err == nil {
			cacheDir = cfg.MangaImageDir
		}
	}
	if cacheDir == "" {
		cacheDir = "manga_cache"
	}

	return MangaModel{
		parent:     parent,
		screen:     mangaScreenSearch,
		scraper:    scraper.NewMangaScraper("https://truyenvn.shop"),
		cfgManager: mgr,
		cacheDir:   cacheDir,
		termWidth:  100,
		termHeight: 30,
	}
}

func (m MangaModel) Init() tea.Cmd {
	return nil
}

func (m MangaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.termWidth = msg.Width
		m.termHeight = msg.Height

	case searchMangaMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.mangaList = msg.results
		m.index = 0
		m.err = nil
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
		}
		return m, nil

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
		// Bắt đầu render ảnh đầu tiên ngay
		if len(m.images) > 0 {
			return m, m.renderImageCmd(0)
		}
		return m, nil

	case chafaRenderMsg:
		m.renderLoading = false
		if msg.err != nil {
			m.err = msg.err
			// Vẫn ở màn reading, hiển thị lỗi
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
		switch msg.String() {
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
				m.screen = mangaScreenSearch
				m.selectedManga = nil
				m.chapterList = nil
				m.index = 0
				m.err = nil
			default:
				return m.parent, nil
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
			default:
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
			default:
				if m.index < len(m.mangaList)-1 {
					m.index++
				}
			}

		case "left", "h":
			// Trang trước (alias)
			if m.screen == mangaScreenReading && m.currentImageIdx > 0 {
				m.currentImageIdx--
				m.renderedImage = ""
				m.renderLoading = true
				return m, m.renderImageCmd(m.currentImageIdx)
			}

		case "right", "l":
			// Trang tiếp (alias)
			if m.screen == mangaScreenReading && m.currentImageIdx < len(m.images)-1 {
				m.currentImageIdx++
				m.renderedImage = ""
				m.renderLoading = true
				return m, m.renderImageCmd(m.currentImageIdx)
			}

		case "backspace":
			if m.screen == mangaScreenSearch && len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}

		case "enter":
			switch m.screen {
			case mangaScreenSearch:
				if len(m.mangaList) > 0 && m.index < len(m.mangaList) {
					selected := m.mangaList[m.index]
					m.selectedManga = &selected
					m.loading = true
					m.err = nil
					return m, m.getMangaDetail(selected.Slug)
				}
				if len(m.searchQuery) >= 2 {
					m.loading = true
					m.err = nil
					m.mangaList = nil
					return m, m.searchManga(m.searchQuery)
				}

			case mangaScreenChapters:
				if len(m.chapterList) > 0 && m.index < len(m.chapterList) {
					m.currentChapIdx = m.index
					m.loading = true
					m.err = nil
					return m, m.loadChapterImages(m.index)
				}

			case mangaScreenReading:
				// Enter = trang tiếp
				if m.currentImageIdx < len(m.images)-1 {
					m.currentImageIdx++
					m.renderedImage = ""
					m.renderLoading = true
					return m, m.renderImageCmd(m.currentImageIdx)
				}
			}

		// 'n' → chapter tiếp theo (mới hơn, index nhỏ hơn trong list đảo ngược)
		case "n":
			if m.screen == mangaScreenReading || m.screen == mangaScreenChapters {
				if m.currentChapIdx > 0 {
					m.currentChapIdx--
					m.index = m.currentChapIdx
					m.loading = true
					return m, m.loadChapterImages(m.currentChapIdx)
				}
				m.statusMsg = "⚠ Đây là chapter mới nhất!"
			}

		// 'p' → chapter cũ hơn
		case "p":
			if m.screen == mangaScreenReading || m.screen == mangaScreenChapters {
				if m.currentChapIdx < len(m.chapterList)-1 {
					m.currentChapIdx++
					m.index = m.currentChapIdx
					m.loading = true
					return m, m.loadChapterImages(m.currentChapIdx)
				}
				m.statusMsg = "⚠ Đây là chapter đầu tiên!"
			}

		// 'r' → reload/re-render ảnh hiện tại
		case "r":
			if m.screen == mangaScreenReading && len(m.images) > 0 {
				m.renderedImage = ""
				m.renderLoading = true
				return m, m.renderImageCmd(m.currentImageIdx)
			}

		// 'w' → mở chapter trong trình duyệt web
		case "w":
			if m.screen == mangaScreenReading && m.currentChapIdx < len(m.chapterList) {
				chapURL := m.chapterList[m.currentChapIdx].URL
				cmd := exec.Command("cmd", "/c", "start", "", chapURL)
				_ = cmd.Start()
				m.statusMsg = "🌐 Đang mở trong trình duyệt..."
			}

		// 'o' → mở file ảnh downloaded bằng viewer
		case "o":
			if m.screen == mangaScreenReading && len(m.images) > 0 {
				return m, m.openInViewer(m.currentImageIdx)
			}

		// 'f' → thêm vào yêu thích
		case "f":
			if m.selectedManga != nil && m.cfgManager != nil {
				item := config.FavoriteItem{
					Name: m.selectedManga.Name,
					Slug: m.selectedManga.Slug,
					Type: "manga",
				}
				if err := m.cfgManager.AddToFavorites(item); err == nil {
					m.statusMsg = fmt.Sprintf("⭐ Đã thêm '%s' vào yêu thích!", m.selectedManga.Name)
				}
			}

		default:
			if m.screen == mangaScreenSearch {
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

		// Lưu lịch sử
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

// renderImageCmd render ảnh tại index bằng chafa, trả về ANSI string
func (m MangaModel) renderImageCmd(imageIdx int) tea.Cmd {
	return func() tea.Msg {
		if imageIdx < 0 || imageIdx >= len(m.images) {
			return chafaRenderMsg{imageIdx: imageIdx, err: fmt.Errorf("index ảnh không hợp lệ")}
		}

		imageURL := m.images[imageIdx]

		// Tính toán kích thước: trừ đi header, hint, status
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

		renderer, err := render.NewChafaRenderer(m.cacheDir, w, h)
		if err != nil {
			// chafa không có → fallback hiển thị URL
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

// openInViewer download ảnh và mở bằng Default Image Viewer của Windows
func (m MangaModel) openInViewer(imageIdx int) tea.Cmd {
	return func() tea.Msg {
		if imageIdx < 0 || imageIdx >= len(m.images) {
			return mangaPlayMsg{err: fmt.Errorf("không có ảnh để mở")}
		}

		imageURL := m.images[imageIdx]
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		path, err := m.scraper.DownloadImage(ctx, imageURL, m.cacheDir)
		if err != nil {
			return mangaPlayMsg{err: fmt.Errorf("không tải được ảnh: %v", err)}
		}

		cmd := exec.Command("cmd", "/c", "start", "", path)
		err = cmd.Start()
		return mangaPlayMsg{err: err}
	}
}

// --- View ---

func (m MangaModel) View() string {
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

	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9f75ff")).
		Bold(true)

	switch m.screen {

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// MÀN TÌM KIẾM
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case mangaScreenSearch:
		b.WriteString(headerStyle.Render("📚 TÌM KIẾM MANGA"))
		b.WriteString("\n\n")

		cursor := "▌"
		if m.loading {
			cursor = " ⏳"
		}
		b.WriteString("Nhập tên manga: ")
		b.WriteString(inputStyle.Render(m.searchQuery + cursor))
		b.WriteString("\n\n")

		if m.err != nil {
			b.WriteString(errStyle.Render(fmt.Sprintf("❌ Lỗi: %v", m.err)))
			b.WriteString("\n")
		} else if m.loading {
			b.WriteString(hintStyle.Render("⏳ Đang tìm kiếm..."))
		} else if len(m.mangaList) == 0 {
			b.WriteString(hintStyle.Render("Nhập tên manga (≥2 ký tự) → Enter để tìm kiếm"))
		} else {
			b.WriteString(hintStyle.Render(fmt.Sprintf("Tìm thấy %d kết quả:\n", len(m.mangaList))))
			for i, manga := range m.mangaList {
				if i == m.index {
					b.WriteString(selectedStyle.Render("▶ " + manga.Name))
				} else {
					b.WriteString(normalStyle.Render("  " + manga.Name))
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

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// MÀN DANH SÁCH CHAPTER
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case mangaScreenChapters:
		title := "DANH SÁCH CHAPTER"
		if m.selectedManga != nil {
			title = fmt.Sprintf("📚 %s", m.selectedManga.Name)
		}
		b.WriteString(headerStyle.Render(title))

		if m.selectedManga != nil {
			infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
			infoParts := []string{}
			if m.selectedManga.Author != "" {
				infoParts = append(infoParts, "Tác giả: "+m.selectedManga.Author)
			}
			if m.selectedManga.Status != "" {
				infoParts = append(infoParts, "Trạng thái: "+m.selectedManga.Status)
			}
			if len(infoParts) > 0 {
				b.WriteString(infoStyle.Render("  " + strings.Join(infoParts, " | ")))
			}
		}
		b.WriteString("\n\n")

		if m.loading {
			b.WriteString(hintStyle.Render("⏳ Đang tải danh sách chapter..."))
		} else if m.err != nil {
			b.WriteString(errStyle.Render(fmt.Sprintf("❌ Lỗi: %v", m.err)))
		} else if len(m.chapterList) == 0 {
			b.WriteString(hintStyle.Render("Không có chapter nào."))
		} else {
			start := 0
			if m.index > 7 {
				start = m.index - 7
			}
			end := start + 15
			if end > len(m.chapterList) {
				end = len(m.chapterList)
			}

			for i := start; i < end; i++ {
				ch := m.chapterList[i]
				line := ch.Title
				if ch.Date != "" {
					line += "  " + ch.Date
				}
				if i == m.currentChapIdx {
					line += " ◀ đang đọc"
				}
				if i == m.index {
					b.WriteString(selectedStyle.Render("▶ " + line))
				} else {
					b.WriteString(normalStyle.Render("  " + line))
				}
				b.WriteString("\n")
			}

			if len(m.chapterList) > 15 {
				b.WriteString(hintStyle.Render(fmt.Sprintf("\n(Tổng: %d chapters)\n", len(m.chapterList))))
			}
		}

		if m.statusMsg != "" {
			b.WriteString("\n")
			b.WriteString(successStyle.Render(m.statusMsg))
		}

		b.WriteString("\n")
		b.WriteString(hintStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Đọc | n: Tiếp | p: Trước | f: Yêu thích"))

	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━
	// MÀN ĐỌC MANGA (chính)
	// ━━━━━━━━━━━━━━━━━━━━━━━━━━━
	case mangaScreenReading:
		chapTitle := "ĐỌC CHAPTER"
		if m.currentChapIdx < len(m.chapterList) {
			chapTitle = fmt.Sprintf("📖 %s", m.chapterList[m.currentChapIdx].Title)
		}

		// Header + progress rút gọn
		total := len(m.images)
		current := m.currentImageIdx + 1
		header := fmt.Sprintf("%s  [%d/%d]", chapTitle, current, total)
		b.WriteString(headerStyle.Render(header))
		b.WriteString("\n")

		if m.loading {
			// Đang tải danh sách ảnh chapter
			b.WriteString(hintStyle.Render("⏳ Đang tải chapter..."))
		} else if m.err != nil {
			b.WriteString(errStyle.Render(fmt.Sprintf("❌ %v", m.err)))
			b.WriteString("\n")
			b.WriteString(hintStyle.Render("r: Thử lại | Esc: Quay lại"))
		} else if len(m.images) == 0 {
			b.WriteString(errStyle.Render("Không tìm thấy ảnh trong chapter này."))
			b.WriteString("\n")
			b.WriteString(hintStyle.Render("Esc: Quay lại | w: Mở trong trình duyệt"))
		} else {
			// ─── Khu vực hiển thị ảnh ───
			if m.renderLoading {
				b.WriteString(progressStyle.Render("  ⏳ Đang render ảnh..."))
				b.WriteString("\n")
			} else if m.renderedImage != "" {
				// Hiển thị ảnh đã render bởi chafa
				b.WriteString(m.renderedImage)
				b.WriteString("\n")
			} else {
				// chafa không có hoặc chưa render
				imgURL := m.images[m.currentImageIdx]
				if len(imgURL) > 80 {
					imgURL = imgURL[:77] + "..."
				}
				imgURLStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
				b.WriteString(imgURLStyle.Render(fmt.Sprintf("  🖼 Trang %d: %s", current, imgURL)))
				b.WriteString("\n")
			}

			// ─── Thanh phím tắt ───
			chafaStatus := ""
			if render.IsAvailable() {
				chafaStatus = " [chafa ✓]"
			}
			b.WriteString("\n")
			b.WriteString(hintStyle.Render(
				fmt.Sprintf("← →/↑↓: Trang | Enter: Trang sau | n/p: Chapter | r: Reload | o: Viewer | w: Browser | f: Yêu thích%s", chafaStatus),
			))
		}

		if m.statusMsg != "" {
			b.WriteString("\n")
			b.WriteString(successStyle.Render(m.statusMsg))
		}
	}

	return b.String()
}
