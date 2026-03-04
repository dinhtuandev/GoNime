package tui

import (
	"fmt"
	"strings"

	"gonime/internal/api"
	"gonime/internal/config"
	"gonime/internal/ui/i18n"
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
	tab         listTab
	cfg         *config.Config
}

func NewFavoritesModel(parent MainModel) FavoritesModel {
	mgr, _ := config.NewManager()
	items, _ := mgr.LoadFavorites()
	cfg, _ := mgr.LoadConfig()
	m := FavoritesModel{
		parent:      parent,
		items:       items,
		manager:     mgr,
		animeClient: api.NewAnimeClient(),
		prevIndex:   -1,
		tab:         tabAnime,
		cfg:         cfg,
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
			if m.tab != tabAnime {
				m.tab = tabAnime
				m.index = 0
				m.posterArt = ""
				if vis := m.visibleIndices(); len(vis) > 0 {
					return m, fetchPosterForFav(m.items[vis[0]])
				}
			}
		case "2":
			// Tab Manga
			if m.tab != tabManga {
				m.tab = tabManga
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
				// realIdx := vis[m.index] // Not needed if we use m.index directly on m.items after filtering
				if m.manager != nil {
					if len(m.items) > 0 && m.index < len(m.items) {
						m.items = append(m.items[:m.index], m.items[m.index+1:]...)
						_ = m.manager.SaveFavorites(m.items)
						if m.index >= len(m.items) && m.index > 0 {
							m.index--
						}
						m.statusMsg = i18n.T(m.cfg, "msg_del_item")
					}
				}
			}
		case "c":
			if m.manager != nil {
				if err := m.manager.ClearFavorites(); err == nil {
					m.items = nil
					m.index = 0
					m.statusMsg = i18n.T(m.cfg, "msg_del_all")
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
		// Delegate to thumbnailCmd logic
		cmd := thumbnailCmd(imageURL)
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
	if m.index >= len(idxs) {
		// Caller sẽ clamp lại index nếu cần.
	}
	return idxs
}

func (m FavoritesModel) View() string {
	var b strings.Builder

	// Tabs: Anime / Manga
	var tabs []string
	animeTab := theme.BadgeActive.Render(i18n.T(m.cfg, "tab_anime"))
	mangaTab := theme.Badge.Render(i18n.T(m.cfg, "tab_manga"))
	if m.tab == tabManga {
		animeTab = theme.Badge.Render(i18n.T(m.cfg, "tab_anime"))
		mangaTab = theme.BadgeActive.Render(i18n.T(m.cfg, "tab_manga"))
	}
	tabs = append(tabs, animeTab, mangaTab)

	b.WriteString(drawHeader("⭐", i18n.T(m.cfg, "title_favorites"), nil, theme.Gold))
	b.WriteString(drawTabs(tabs))
	b.WriteString(drawHeaderLine(theme.Gold))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(drawLoadingLine(i18n.T(m.cfg, "loading_generic"), m.tick))
		b.WriteString("\n")
	} else if m.err != nil {
		b.WriteString(drawErrorLine(fmt.Sprintf(i18n.T(m.cfg, "err_prefix"), m.err)))
		b.WriteString("\n\n")
	}

	visible := m.visibleIndices()

	if len(visible) == 0 {
		b.WriteString("  ")
		if m.tab == tabAnime {
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "empty_fav_anime")))
		} else {
			b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "empty_fav_manga")))
		}
		b.WriteString("\n  ")
		b.WriteString(theme.Muted.Render(i18n.T(m.cfg, "empty_fav_hint_a"))) // Note: I'll use generic hint or split
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
			nameStr := truncateString(item.Name, 35) // Truncate to avoid overflow
			if row == m.index {
				marker = "▶ "
				nameStr = selNameStyle.Render(nameStr)
			} else {
				nameStr = normNameStyle.Render(nameStr)
			}
			leftB.WriteString(fmt.Sprintf("  %s %-10s  %s\n", typeTag, marker, nameStr))
		}
		leftB.WriteString("\n  ")
		leftB.WriteString(theme.HintKey.Render(fmt.Sprintf("%d", len(visible))))
		if m.tab == tabAnime {
			leftB.WriteString(theme.Muted.Render(fmt.Sprintf(" %s", i18n.T(m.cfg, "count_fav_anime"))))
		} else {
			leftB.WriteString(theme.Muted.Render(fmt.Sprintf(" %s", i18n.T(m.cfg, "count_fav_manga"))))
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
			info.WriteString("\n" + accentStyle.Render(truncateString(item.Name, 45)) + "\n")
			if item.Type == "manga" {
				info.WriteString(theme.Muted.Render("📚 Manga") + "\n")
			} else {
				info.WriteString(theme.Muted.Render("📺 Anime") + "\n")
			}
			posterStyle := lipgloss.NewStyle().Padding(0, 1)
			infoStyle := lipgloss.NewStyle().Padding(0, 1).Width(45)
			rightStr = lipgloss.JoinHorizontal(lipgloss.Top,
				posterStyle.Render(m.posterArt),
				infoStyle.Render(info.String()),
			)
		}

		leftStyle := lipgloss.NewStyle().Width(50)
		rightStyle := lipgloss.NewStyle().Width(45)
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
		i18n.T(m.cfg, "hint_enter_open"),
		i18n.T(m.cfg, "hint_switch_tab"),
		i18n.T(m.cfg, "hint_del_item"),
		i18n.T(m.cfg, "hint_del_all"),
	}))
	b.WriteString("\n")

	return b.String()
}
