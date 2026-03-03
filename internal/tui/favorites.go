package tui

import (
	"fmt"
	"strings"

	"gonime/internal/api"
	"gonime/internal/config"
	"gonime/internal/ui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FavoritesModel struct {
	parent      MainModel
	items       []config.FavoriteItem
	index       int
	prevIndex   int
	loading     bool
	err         error
	statusMsg   string
	manager     *config.Manager
	animeClient *api.AnimeClient
	posterArt   string // chafa-rendered thumbnail for selected item
	tick        int
	tab         favoritesTab
}

type favoritesTab int

const (
	favTabAnime favoritesTab = iota
	favTabManga
)


func NewFavoritesModel(parent MainModel) FavoritesModel {
	mgr, _ := config.NewManager()
	items, _ := mgr.LoadFavorites()
	m := FavoritesModel{
		parent:      parent,
		items:       items,
		manager:     mgr,
		animeClient: api.NewAnimeClient(),
		prevIndex:   -1,
		tab:         favTabAnime,
	}
	return m
}

func (m FavoritesModel) Init() tea.Cmd {
	cmds := []tea.Cmd{tickCmd()}
	if len(m.items) > 0 {
		cmds = append(cmds, fetchPosterForFav(m.items[0]))
	}
	return tea.Batch(cmds...)
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
			if m.tab != favTabAnime {
				m.tab = favTabAnime
				m.index = 0
				m.posterArt = ""
				if vis := m.visibleIndices(); len(vis) > 0 {
					return m, fetchPosterForFav(m.items[vis[0]])
				}
			}
		case "2":
			// Tab Manga
			if m.tab != favTabManga {
				m.tab = favTabManga
				m.index = 0
				m.posterArt = ""
				if vis := m.visibleIndices(); len(vis) > 0 {
					return m, fetchPosterForFav(m.items[vis[0]])
				}
			}
		case "up", "k":
			if vis := m.visibleIndices(); len(vis) > 0 && m.index > 0 {
				m.index--
				m.posterArt = ""
				return m, fetchPosterForFav(m.items[vis[m.index]])
			}
		case "down", "j":
			if vis := m.visibleIndices(); len(vis) > 0 && m.index < len(vis)-1 {
				m.index++
				m.posterArt = ""
				return m, fetchPosterForFav(m.items[vis[m.index]])
			}
		case "enter":
			if vis := m.visibleIndices(); len(vis) > 0 && m.index < len(vis) {
				item := m.items[vis[m.index]]
				m.loading = true
				m.err = nil
				m.statusMsg = ""
				return m, m.openFavorite(item)
			}
		case "d":
			if vis := m.visibleIndices(); len(vis) > 0 && m.index < len(vis) {
				realIdx := vis[m.index]
				slug := m.items[realIdx].Slug
				if m.manager != nil {
					_ = m.manager.RemoveFromFavorites(slug)
				}
				m.items = append(m.items[:realIdx], m.items[realIdx+1:]...)
				if m.index > 0 {
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

// fetchPosterForFav fetches the thumbnail URL for a favorite item then renders it with chafa
func fetchPosterForFav(item config.FavoriteItem) tea.Cmd {
	return func() tea.Msg {
		var imageURL string
		switch item.Type {
		case "manga":
			s := newScraper()
			manga, err := s.GetDetails(item.Slug)
			if err == nil && manga != nil {
				imageURL = manga.Cover
			}
		default: // anime
			client := api.NewAnimeClient()
			movie, _, err := client.GetEpisodes(item.Slug)
			if err == nil && movie != nil {
				imageURL = movie.ThumbURL
				if imageURL == "" {
					imageURL = movie.PosterURL
				}
			}
		}
		if imageURL == "" {
			return posterRenderMsg{}
		}
		// Delegate to posterCmd's logic
		cmd := posterCmd(imageURL)
		return cmd()
	}
}

// scraper_instance tạo scraper cho manga (helper không cần import)
func scraper_instance() scraper_iface {
	return newScraper()
}

// visibleIndices trả về index thực trong m.items theo tab hiện tại.
func (m FavoritesModel) visibleIndices() []int {
	var idxs []int
	for i, it := range m.items {
		switch m.tab {
		case favTabAnime:
			if it.Type != "manga" {
				idxs = append(idxs, i)
			}
		case favTabManga:
			if it.Type == "manga" {
				idxs = append(idxs, i)
			}
		}
	}
	if m.index >= len(idxs) {
		// Caller sẽ clamp lại index nếu cần.
	}
	return idxs
}

func (m FavoritesModel) View() string {
	var b strings.Builder

	// Tabs: Anime / Manga
	var tabs []string
	animeTab := theme.BadgeActive.Render("Anime / Phim")
	mangaTab := theme.Badge.Render("Manga")
	if m.tab == favTabManga {
		animeTab = theme.Badge.Render("Anime / Phim")
		mangaTab = theme.BadgeActive.Render("Manga")
	}
	tabs = append(tabs, animeTab, mangaTab)

	b.WriteString(drawHeader("⭐", "YÊU THÍCH", tabs, theme.Gold))
	b.WriteString(drawHeaderLine(theme.Gold))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(drawLoadingLine("Đang tải...", m.tick))
		b.WriteString("\n")
	} else if m.err != nil {
		b.WriteString(drawErrorLine(fmt.Sprintf("Lỗi: %v", m.err)))
		b.WriteString("\n\n")
	}

	visible := m.visibleIndices()

	if len(visible) == 0 {
		b.WriteString("  ")
		if m.tab == favTabAnime {
			b.WriteString(theme.Muted.Render("Chưa có anime/phim yêu thích nào."))
		} else {
			b.WriteString(theme.Muted.Render("Chưa có manga yêu thích nào."))
		}
		b.WriteString("\n  ")
		b.WriteString(theme.Muted.Render("Nhấn "))
		b.WriteString(theme.HintKey.Render("f"))
		if m.tab == favTabAnime {
			b.WriteString(theme.Muted.Render(" khi xem anime/phim để thêm vào đây."))
		} else {
			b.WriteString(theme.Muted.Render(" khi đọc manga để thêm vào đây."))
		}
	} else {
		selNameStyle := lipgloss.NewStyle().Foreground(theme.Gold).Bold(true)
		normNameStyle := lipgloss.NewStyle().Foreground(theme.TextSecondary)
		animeTagStyle := lipgloss.NewStyle().Foreground(theme.AnimeAccent).Bold(true)
		mangaTagStyle := lipgloss.NewStyle().Foreground(theme.MangaAccent).Bold(true)

		// Build left list (filtered by tab)
		var leftB strings.Builder
		leftB.WriteString("\n")
		for row, idx := range visible {
			item := m.items[idx]
			var typeTag string
			if item.Type == "manga" {
				typeTag = mangaTagStyle.Render("[Manga]")
			} else {
				typeTag = animeTagStyle.Render("[Anime]")
			}
			marker := "  "
			var nameStr string
			if row == m.index {
				marker = "▶ "
				nameStr = selNameStyle.Render(item.Name)
			} else {
				nameStr = normNameStyle.Render(item.Name)
			}
			leftB.WriteString(fmt.Sprintf("  %s %-10s  %s\n", typeTag, marker, nameStr))
		}
		leftB.WriteString("\n  ")
		leftB.WriteString(theme.HintKey.Render(fmt.Sprintf("%d", len(visible))))
		if m.tab == favTabAnime {
			leftB.WriteString(theme.Muted.Render(" anime/phim yêu thích"))
		} else {
			leftB.WriteString(theme.Muted.Render(" manga yêu thích"))
		}

		// Build right poster panel
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
			if item.Type == "manga" {
				info.WriteString(theme.Muted.Render("📚 Manga") + "\n")
			} else {
				info.WriteString(theme.Muted.Render("📺 Anime") + "\n")
			}
			posterStyle := lipgloss.NewStyle().Padding(0, 1)
			infoStyle := lipgloss.NewStyle().Padding(0, 1).Width(40)
			rightStr = lipgloss.JoinHorizontal(lipgloss.Top,
				posterStyle.Render(m.posterArt),
				infoStyle.Render(info.String()),
			)
		}

		leftStyle := lipgloss.NewStyle().Width(50)
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
		"Enter: Mở",
		"1/2: Anime/Manga",
		"d: Xóa mục",
		"c: Xóa tất cả",
	}))
	b.WriteString("\n")

	return b.String()
}
