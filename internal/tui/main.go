package tui

import (
	"fmt"
	"strings"

	"gonime/internal/config"
	"gonime/internal/ui/i18n"
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
	cfg       *config.Config
	ui        *uiState
	help      bool
}

func NewMainModel() MainModel {
	mgr, _ := config.NewManager()
	cfg, _ := mgr.LoadConfig()

	return MainModel{
		screen:    ScreenMain,
		cfg:       cfg,
		menuItems: getMenu(cfg),
		ui:        newUIState(),
	}
}

func getMenu(cfg *config.Config) []MenuItem {
	return []MenuItem{
		{"📺", i18n.T(cfg, "menu.watch_anime"), i18n.T(cfg, "menu.watch_anime_desc")},
		{"📚", i18n.T(cfg, "menu.read_manga"), i18n.T(cfg, "menu.read_manga_desc")},
		{"📜", i18n.T(cfg, "menu.history"), i18n.T(cfg, "menu.history_desc")},
		{"⭐", i18n.T(cfg, "menu.favorites"), i18n.T(cfg, "menu.favorites_desc")},
		{"⚙️ ", i18n.T(cfg, "menu.settings"), i18n.T(cfg, "menu.settings_desc")},
		{"🚪", i18n.T(cfg, "menu.exit"), i18n.T(cfg, "menu.exit_desc")},
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

		case "?":
			m.help = !m.help

		case "esc":
			m.help = false

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
		}

	case tea.MouseMsg:
		return m.handleMouse(msg)
	}

	return m, nil
}

func (m MainModel) handleMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if msg.Action != tea.MouseActionPress {
		return m, nil
	}
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		if m.menuIndex > 0 {
			m.menuIndex--
		}
	case tea.MouseButtonWheelDown:
		if m.menuIndex < len(m.menuItems)-1 {
			m.menuIndex++
		}
	case tea.MouseButtonLeft:
		if key, ok := m.ui.hintKeyAt(msg.X, msg.Y); ok {
			return m.Update(key)
		}
		if row, ok := m.ui.rowAt(msg.X, msg.Y); ok && row >= 0 && row < len(m.menuItems) {
			if row == m.menuIndex {
				return m.handleMenuSelection()
			}
			m.menuIndex = row
		}
	}
	return m, nil
}

func (m MainModel) handleMenuSelection() (tea.Model, tea.Cmd) {
	switch m.menuIndex {
	case 0:
		m.screen = ScreenAnime
		am := NewAnimeModel(m)
		return am, am.Init()
	case 1:
		m.screen = ScreenManga
		mm := NewMangaModel(m)
		return mm, mm.Init()
	case 2:
		m.screen = ScreenHistory
		hm := NewHistoryModel(m)
		return hm, hm.Init()
	case 3:
		m.screen = ScreenFavorites
		fm := NewFavoritesModel(m)
		return fm, fm.Init()
	case 4:
		m.screen = ScreenSettings
		sm := NewSettingsModel(m)
		return sm, sm.Init()
	case 5:
		return m, tea.Quit
	}
	return m, nil
}

// ─── View ─────────────────────────────────────────────────────────────────────

func (m MainModel) hints() []string {
	return []string{
		i18n.T(m.cfg, "hint_nav"),
		i18n.T(m.cfg, "hint_enter_select_2"),
		"?",
		i18n.T(m.cfg, "hint_quit"),
	}
}

func (m MainModel) View() string {
	m.ui.reset()
	if m.help {
		return renderHelp(m.width, m.height, m.hints())
	}

	var b strings.Builder

	// ── Logo ──────────────────────────────────────────────────────────────────
	logo := buildLogo()
	b.WriteString(logo)
	b.WriteString("\n")

	// ── Subtitle ──────────────────────────────────────────────────────────────
	subtitleStyle := lipgloss.NewStyle().
		Foreground(theme.TextSecondary).
		Padding(0, 4)
	b.WriteString(subtitleStyle.Render(i18n.T(m.cfg, "app.subtitle")))
	b.WriteString("\n")

	// ── Divider ───────────────────────────────────────────────────────────────
	b.WriteString("\n  ")
	b.WriteString(lipgloss.NewStyle().Foreground(theme.TextDim).Render(strings.Repeat("─", 56)))
	b.WriteString("\n\n")

	// ── Menu Items ────────────────────────────────────────────────────────────
	m.ui.beginRows(&b, 0, len(m.menuItems))
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
	writeStatusBar(&b, m.ui, m.width, m.hints())
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
		` ██████╗   ██████╗  ███╗  ██╗ ██╗ ███╗ ███╗ ███████╗`,
		`██╔════╝  ██╔═══██╗ ████╗ ██║ ██║ ████████║ ██╔════╝`,
		`██║  ██╗  ██║   ██║ ██╔██╗██║ ██║ ██╔██╔██║ █████╗  `,
		`██║  ╚██╗ ██║   ██║ ██║╚████║ ██║ ██║╚═╝██║ ██╔══╝  `,
		`╚██████╔╝ ╚██████╔╝ ██║ ╚███║ ██║ ██║   ██║ ███████╗`,
		` ╚═════╝   ╚═════╝  ╚═╝  ╚══╝ ╚═╝ ╚═╝   ╚═╝ ╚══════╝`,
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
