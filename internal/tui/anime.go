package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gonime/internal/api"
)

type searchAnimeMsg struct {
	results []api.Anime
	err     error
}

type AnimeModel struct {
	parent        MainModel
	searchQuery   string
	animeList     []api.Anime
	selectedAnime *api.Anime
	episodeList   []api.Episode
	index         int
	loading       bool
	err           error
	client        *api.AnimeClient
}

func NewAnimeModel(parent MainModel) AnimeModel {
	return AnimeModel{
		parent: parent,
		client: api.NewAnimeClient("https://phimapi.com"),
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
		m.animeList = msg.results
		m.index = 0
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.selectedAnime != nil {
				m.selectedAnime = nil
				m.episodeList = nil
				m.index = 0
				return m, nil
			}
			return m.parent, nil
		case "up", "k":
			if m.selectedAnime != nil {
				if m.index > 0 {
					m.index--
				}
			} else {
				if m.index > 0 {
					m.index--
				}
			}
		case "down", "j":
			if m.selectedAnime != nil {
				if m.index < len(m.episodeList)-1 {
					m.index++
				}
			} else {
				if m.index < len(m.animeList)-1 {
					m.index++
				}
			}
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}
		case "enter":
			if m.selectedAnime != nil {
				if len(m.episodeList) > 0 && m.index < len(m.episodeList) {
					return m, nil
				}
				return m, nil
			}
			if len(m.animeList) > 0 && m.index < len(m.animeList) {
				m.selectedAnime = &m.animeList[m.index]
				m.episodeList = m.selectedAnime.Episodes
				m.index = 0
				return m, nil
			}
			if len(m.searchQuery) >= 2 {
				m.loading = true
				m.err = nil
				return m, m.searchAnime(m.searchQuery)
			}
		default:
			if len(msg.String()) == 1 && !m.loading {
				m.searchQuery += msg.String()
			}
		}
	}

	return m, nil
}

func (m AnimeModel) searchAnime(query string) tea.Cmd {
	return func() tea.Msg {
		results, err := m.client.Search(query)
		return searchAnimeMsg{results: results, err: err}
	}
}

func (m AnimeModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#33a5f5")).
		Bold(true).
		Margin(1, 0)

	if m.selectedAnime != nil {
		b.WriteString(headerStyle.Render(fmt.Sprintf("📺 %s - DANH SÁCH TẬP", m.selectedAnime.Name)))
		b.WriteString("\n\n")

		if len(m.episodeList) == 0 {
			emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
			b.WriteString(emptyStyle.Render("Không có tập nào."))
		} else {
			listStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f0f4ff"))
			selectedListStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9f75ff")).
				Bold(true)

			for i, ep := range m.episodeList {
				line := fmt.Sprintf("Tập %s - %s", ep.Number, ep.Title)
				if i == m.index {
					b.WriteString(selectedListStyle.Render("▶ " + line))
				} else {
					b.WriteString(listStyle.Render("  " + line))
				}
				b.WriteString("\n")
			}
		}

		b.WriteString("\n")
		footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
		b.WriteString(footerStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Xem"))
		return b.String()
	}

	b.WriteString(headerStyle.Render("📺 TÌM KIẾM ANIME"))
	b.WriteString("\n\n")

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#9f75ff")).
		Padding(0, 1).
		Width(50)

	b.WriteString("Nhập tên anime: ")
	cursor := "▌"
	if m.loading {
		cursor = " ⏳"
	}
	b.WriteString(inputStyle.Render(m.searchQuery + cursor))
	b.WriteString("\n\n")

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6b6b"))
		b.WriteString(errStyle.Render(fmt.Sprintf("Lỗi: %v", m.err)))
	} else if m.loading {
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
		b.WriteString(hintStyle.Render("Đang tìm kiếm..."))
	} else if len(m.animeList) == 0 {
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
		b.WriteString(hintStyle.Render("Nhập tên anime (ít nhất 2 ký tự) và nhấn Enter để tìm kiếm..."))
	} else {
		listStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f0f4ff"))
		selectedListStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9f75ff")).
			Bold(true)

		for i, anime := range m.animeList {
			if i == m.index {
				b.WriteString(selectedListStyle.Render(fmt.Sprintf("▶ %s", anime.Name)))
			} else {
				b.WriteString(listStyle.Render(fmt.Sprintf("  %s", anime.Name)))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
	b.WriteString(footerStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Chọn"))

	return b.String()
}
