package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gonime/internal/api"
	"gonime/internal/config"
)

type FavoritesModel struct {
	parent      MainModel
	items       []config.FavoriteItem
	index       int
	loading     bool
	err         error
	statusMsg   string
	manager     *config.Manager
	animeClient *api.AnimeClient
}

func NewFavoritesModel(parent MainModel) FavoritesModel {
	mgr, _ := config.NewManager()
	items, _ := mgr.LoadFavorites()
	return FavoritesModel{
		parent:      parent,
		items:       items,
		manager:     mgr,
		animeClient: api.NewAnimeClient(),
	}
}

func (m FavoritesModel) Init() tea.Cmd {
	return nil
}

func (m FavoritesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case loadEpisodesMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		// Chuyển sang anime model để xem
		if m.index < len(m.items) {
			item := m.items[m.index]
			animeModel := NewAnimeModel(m.parent)
			// Tạo AnimeItem giả để có thể xem
			animeModel.selectedAnime = &api.AnimeItem{
				Name: item.Name,
				Slug: item.Slug,
			}
			animeModel.movie = msg.movie
			animeModel.episodeList = msg.episodes
			animeModel.screen = animeScreenEpisodes
			return animeModel, nil
		}
		return m, nil

	case mangaDetailMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		if msg.manga != nil && m.index < len(m.items) {
			item := m.items[m.index]
			mangaModel := NewMangaModel(m.parent)
			mangaModel.selectedManga = msg.manga
			mangaModel.chapterList = msg.manga.Chapters
			mangaModel.screen = mangaScreenChapters
			// Đặt selectedManga đúng
			mangaModel.selectedManga.Name = item.Name
			return mangaModel, nil
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
				return m, m.openFavorite(item)
			}
		case "d":
			if len(m.items) > 0 && m.index < len(m.items) {
				slug := m.items[m.index].Slug
				if m.manager != nil {
					_ = m.manager.RemoveFromFavorites(slug)
				}
				m.items = append(m.items[:m.index], m.items[m.index+1:]...)
				if m.index >= len(m.items) && m.index > 0 {
					m.index--
				}
				m.statusMsg = "Đã xóa khỏi danh sách yêu thích"
			}
		case "c":
			if m.manager != nil {
				if err := m.manager.ClearFavorites(); err == nil {
					m.items = nil
					m.index = 0
					m.statusMsg = "Đã xóa toàn bộ danh sách yêu thích"
				}
			}
		}
	}

	return m, nil
}

func (m FavoritesModel) openFavorite(item config.FavoriteItem) tea.Cmd {
	return func() tea.Msg {
		switch item.Type {
		case "manga":
			s := scraper_instance()
			manga, err := s.GetDetails(item.Slug)
			return mangaDetailMsg{manga: manga, err: err}
		default: // anime
			client := api.NewAnimeClient()
			movie, episodes, err := client.GetEpisodes(item.Slug)
			return loadEpisodesMsg{movie: movie, episodes: episodes, err: err}
		}
	}
}

// scraper_instance tạo scraper cho manga (helper không cần import)
func scraper_instance() scraper_iface {
	return newScraper()
}

func (m FavoritesModel) View() string {
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

	b.WriteString(headerStyle.Render("⭐ YÊU THÍCH"))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(hintStyle.Render("⏳ Đang tải..."))
		b.WriteString("\n")
	} else if m.err != nil {
		b.WriteString(errStyle.Render(fmt.Sprintf("❌ Lỗi: %v", m.err)))
		b.WriteString("\n\n")
	}

	if len(m.items) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
		b.WriteString(emptyStyle.Render("Chưa có mục yêu thích nào."))
		b.WriteString("\n")
		b.WriteString(hintStyle.Render("Thêm yêu thích bằng phím 'f' khi xem anime/manga."))
	} else {
		typeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#33a5f5"))

		for i, item := range m.items {
			typeIcon := "🎬"
			typeLabel := "Anime"
			if item.Type == "manga" {
				typeIcon = "📚"
				typeLabel = "Manga"
			}

			line := fmt.Sprintf("%s [%s] %s", typeIcon, typeLabel, item.Name)

			if i == m.index {
				b.WriteString(selectedStyle.Render("▶ " + line))
			} else {
				b.WriteString(normalStyle.Render("  " + line))
			}
			b.WriteString("\n")
		}

		b.WriteString("\n")
		b.WriteString(typeStyle.Render(fmt.Sprintf("Tổng: %d mục yêu thích", len(m.items))))
	}

	if m.statusMsg != "" {
		b.WriteString("\n")
		b.WriteString(successStyle.Render(m.statusMsg))
	}

	b.WriteString("\n\n")
	b.WriteString(hintStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Mở | d: Xóa | c: Xóa tất cả"))

	return b.String()
}
