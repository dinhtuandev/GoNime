package tui

import (
	"fmt"
	"strings"

	"gonime/internal/ui/theme"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Screen int

const (
	ScreenMain Screen = iota
	ScreenAnime
	ScreenAnimeSearch
	ScreenAnimeList
	ScreenAnimePlayer
	ScreenManga
	ScreenMangaSearch
	ScreenMangaList
	ScreenMangaReader
	ScreenHistory
	ScreenFavorites
	ScreenSettings
)

// ─── Menu item ────────────────────────────────────────────────────────────────

type MenuItem struct {
	Icon  string
	Title string
	Desc  string
}

// ─── Model ────────────────────────────────────────────────────────────────────

type MainModel struct {
	screen    Screen
	menuIndex int
	menuItems []MenuItem
	width     int
	height    int
	tick      int
}

func NewMainModel() MainModel {
	return MainModel{
		screen: ScreenMain,
		menuItems: []MenuItem{
			{"📺", "Xem Anime", "Khám phá & xem anime vietsub"},
			{"📚", "Đọc Manga", "Đọc manga từ nhiều nguồn"},
			{"📜", "Lịch sử", "Xem lại các anime/manga đã xem"},
			{"⭐", "Yêu thích", "Danh sách anime/manga yêu thích"},
			{"⚙️ ", "Cài đặt", "Cấu hình ứng dụng"},
			{"🚪", "Thoát", "Thoát khỏi ứng dụng"},
		},
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			if m.menuIndex > 0 {
				m.menuIndex--
			}

		case "down", "j":
			if m.menuIndex < len(m.menuItems)-1 {
				m.menuIndex++
			}

		case "enter", " ":
			return m.handleMenuSelection()

		case "esc":
			if m.screen != ScreenMain {
				m.screen = ScreenMain
			}
		}
	}

	return m, nil
}

func (m MainModel) handleMenuSelection() (tea.Model, tea.Cmd) {
	switch m.menuIndex {
	case 0:
		m.screen = ScreenAnime
		return NewAnimeModel(m), nil
	case 1:
		m.screen = ScreenManga
		return NewMangaModel(m), nil
	case 2:
		m.screen = ScreenHistory
		return NewHistoryModel(m), nil
	case 3:
		m.screen = ScreenFavorites
		return NewFavoritesModel(m), nil
	case 4:
		m.screen = ScreenSettings
		return NewSettingsModel(m), nil
	case 5:
		return m, tea.Quit
	}
	return m, nil
}

// ─── View ─────────────────────────────────────────────────────────────────────

func (m MainModel) View() string {
	var b strings.Builder

	// ── Logo ──────────────────────────────────────────────────────────────────
	logo := buildLogo()
	b.WriteString(logo)
	b.WriteString("\n")

	// ── Subtitle ──────────────────────────────────────────────────────────────
	subtitleStyle := lipgloss.NewStyle().
		Foreground(theme.TextSecondary).
		Padding(0, 4)
	b.WriteString(subtitleStyle.Render("🎌  Anime & Manga CLI Player  ·  Vietnamese Subtitles  🎌"))
	b.WriteString("\n")

	// ── Divider ───────────────────────────────────────────────────────────────
	b.WriteString("\n  ")
	b.WriteString(lipgloss.NewStyle().Foreground(theme.TextDim).Render(strings.Repeat("─", 56)))
	b.WriteString("\n\n")

	// ── Menu Items ────────────────────────────────────────────────────────────
	for i, item := range m.menuItems {
		if i == m.menuIndex {
			b.WriteString(renderSelectedMenuItem(i+1, item))
		} else {
			b.WriteString(renderNormalMenuItem(i+1, item))
		}
		b.WriteString("\n")
	}

	// ── Bottom divider ────────────────────────────────────────────────────────
	b.WriteString("\n  ")
	b.WriteString(lipgloss.NewStyle().Foreground(theme.TextDim).Render(strings.Repeat("─", 56)))

	// ── Status bar ────────────────────────────────────────────────────────────
	b.WriteString(drawStatusBar([]string{
		"↑↓/jk: Di chuyển",
		"Enter: Chọn",
		"q: Thoát",
	}))
	b.WriteString("\n")

	return b.String()
}

// renderSelectedMenuItem renders the currently highlighted menu item.
func renderSelectedMenuItem(num int, item MenuItem) string {
	// Number badge
	numStyle := lipgloss.NewStyle().
		Foreground(theme.BackgroundDark).
		Background(theme.Primary).
		Bold(true).
		Padding(0, 1)

	// Selected row box
	rowStyle := lipgloss.NewStyle().
		Foreground(theme.Primary).
		Background(theme.BackgroundLight).
		Bold(true).
		Padding(0, 1)

	descStyle := lipgloss.NewStyle().
		Foreground(theme.TextSecondary).
		Background(theme.BackgroundLight)

	numBadge := numStyle.Render(fmt.Sprintf("%d", num))
	row := rowStyle.Render(fmt.Sprintf("  ▶  %s  %s", item.Icon, item.Title))
	desc := descStyle.Render("  " + item.Desc)

	return "  " + numBadge + row + desc
}

// renderNormalMenuItem renders an unselected menu item.
func renderNormalMenuItem(num int, item MenuItem) string {
	numStyle := lipgloss.NewStyle().
		Foreground(theme.TextMuted).
		Padding(0, 1)

	titleStyle := lipgloss.NewStyle().
		Foreground(theme.TextSecondary).
		Padding(0, 1)

	descStyle := lipgloss.NewStyle().
		Foreground(theme.TextMuted)

	numBadge := numStyle.Render(fmt.Sprintf("%d", num))
	title := titleStyle.Render(fmt.Sprintf("     %s  %s", item.Icon, item.Title))
	desc := descStyle.Render("  " + item.Desc)

	return "  " + numBadge + title + desc
}

// buildLogo constructs the multi-line GONIME ASCII logo with alternating colors.
func buildLogo() string {
	lines := []string{
		` ██████╗  ██████╗  ███╗  ██╗ ██╗ ███╗ ███╗ ███████╗`,
		`██╔════╝ ██╔═══██╗ ████╗ ██║ ██║ ████████║ ██╔════╝`,
		`██║  ██╗ ██║   ██║ ██╔██╗██║ ██║ ██╔██╔██║ █████╗  `,
		`██║  ╚██╗██║   ██║ ██║╚████║ ██║ ██║╚═╝██║ ██╔══╝  `,
		`╚██████╔╝╚██████╔╝ ██║ ╚███║ ██║ ██║   ██║ ███████╗`,
		` ╚═════╝  ╚═════╝  ╚═╝  ╚══╝ ╚═╝ ╚═╝   ╚═╝ ╚══════╝`,
	}

	colors := []lipgloss.Color{
		theme.Primary,
		theme.Primary,
		lipgloss.Color("#40e0ff"),
		lipgloss.Color("#80eeff"),
		theme.Secondary,
		theme.Secondary,
	}

	var b strings.Builder
	b.WriteString("\n")
	for i, line := range lines {
		c := colors[i%len(colors)]
		s := lipgloss.NewStyle().Foreground(c).Bold(true).Padding(0, 2)
		b.WriteString(s.Render(line))
		b.WriteString("\n")
	}
	return b.String()
}
