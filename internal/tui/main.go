package tui

import (
	"fmt"
	"strings"

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

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#33a5f5")).
			Bold(true).
			Padding(1, 2).
			Render

	menuStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f0f4ff")).
			Padding(0, 4).
			Render

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9f75ff")).
			Bold(true).
			Padding(0, 4).
			Render

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#33a5f5")).
			Padding(1, 2).
			Margin(1, 2)
)

type MenuItem struct {
	Title string
	Desc  string
}

type MainModel struct {
	screen    Screen
	menuIndex int
	menuItems []MenuItem
	width     int
	height    int
}

func NewMainModel() MainModel {
	return MainModel{
		screen: ScreenMain,
		menuItems: []MenuItem{
			{"Xem Anime", "Tìm kiếm và xem anime vietsub"},
			{"Đọc Manga", "Đọc manga từ truyenvn.shop"},
			{"Lịch sử", "Xem lại các anime/manga đã xem"},
			{"Yêu thích", "Danh sách anime/manga yêu thích"},
			{"Cài đặt", "Thay đổi cấu hình ứng dụng"},
			{"Thoát", "Thoát khỏi ứng dụng"},
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

func (m MainModel) View() string {
	var b strings.Builder

	title := `
   ██████╗  ██████╗ ██████╗ ███╗   ██╗
  ██╔════╝ ██╔═══██╗██╔══██╗████╗  ██║
  ██║      ██║   ██║██║  ██║██╔██╗ ██║
  ██║      ██║   ██║██║  ██║██║╚██╗██║
  ╚██████╗ ╚██████╔╝██████╔╝██║ ╚████║
   ╚═════╝  ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝`

	b.WriteString(titleStyle(title))
	b.WriteString("\n\n")
	b.WriteString(menuStyle("🎬 Anime & Manga CLI Player"))
	b.WriteString("\n\n")

	for i, item := range m.menuItems {
		if i == m.menuIndex {
			b.WriteString(selectedStyle(fmt.Sprintf("▶ %s - %s", item.Title, item.Desc)))
		} else {
			b.WriteString(menuStyle(fmt.Sprintf("  %s - %s", item.Title, item.Desc)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(menuStyle("↑/↓: Điều hướng | Enter: Chọn | q: Thoát"))

	return b.String()
}
