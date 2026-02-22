package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gonime/internal/scraper"
)

type searchMangaMsg struct {
	results []scraper.Manga
	err     error
}

type mangaDetailMsg struct {
	manga *scraper.Manga
	err   error
}

type MangaModel struct {
	parent        MainModel
	searchQuery   string
	mangaList     []scraper.Manga
	selectedManga *scraper.Manga
	chapterList   []scraper.Chapter
	index         int
	loading       bool
	err           error
	scraper       *scraper.MangaScraper
}

func NewMangaModel(parent MainModel) MangaModel {
	return MangaModel{
		parent:  parent,
		scraper: scraper.NewMangaScraper("https://truyenvn.shop"),
	}
}

func (m MangaModel) Init() tea.Cmd {
	return nil
}

func (m MangaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case searchMangaMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.mangaList = msg.results
		m.index = 0
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
			m.index = 0
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.selectedManga != nil {
				m.selectedManga = nil
				m.chapterList = nil
				m.index = 0
				return m, nil
			}
			return m.parent, nil
		case "up", "k":
			if m.selectedManga != nil {
				if m.index > 0 {
					m.index--
				}
			} else {
				if m.index > 0 {
					m.index--
				}
			}
		case "down", "j":
			if m.selectedManga != nil {
				if m.index < len(m.chapterList)-1 {
					m.index++
				}
			} else {
				if m.index < len(m.mangaList)-1 {
					m.index++
				}
			}
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}
		case "enter":
			if m.selectedManga != nil {
				if len(m.chapterList) > 0 && m.index < len(m.chapterList) {
					return m, nil
				}
				return m, nil
			}
			if len(m.mangaList) > 0 && m.index < len(m.mangaList) {
				m.loading = true
				m.err = nil
				slug := m.mangaList[m.index].Slug
				return m, m.getMangaDetail(slug)
			}
			if len(m.searchQuery) >= 2 {
				m.loading = true
				m.err = nil
				return m, m.searchManga(m.searchQuery)
			}
		default:
			if len(msg.String()) == 1 && !m.loading {
				m.searchQuery += msg.String()
			}
		}
	}

	return m, nil
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

func (m MangaModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#33a5f5")).
		Bold(true).
		Margin(1, 0)

	if m.selectedManga != nil {
		b.WriteString(headerStyle.Render(fmt.Sprintf("📚 %s - DANH SÁCH CHAPTER", m.selectedManga.Name)))
		b.WriteString("\n\n")

		if len(m.chapterList) == 0 {
			emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
			b.WriteString(emptyStyle.Render("Không có chapter nào."))
		} else {
			listStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f0f4ff"))
			selectedListStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9f75ff")).
				Bold(true)

			start := 0
			if m.index > 10 {
				start = m.index - 5
			}
			end := start + 15
			if end > len(m.chapterList) {
				end = len(m.chapterList)
			}

			for i := start; i < end; i++ {
				ch := m.chapterList[i]
				line := fmt.Sprintf("%s - %s", ch.Number, ch.Title)
				if i == m.index {
					b.WriteString(selectedListStyle.Render("▶ " + line))
				} else {
					b.WriteString(listStyle.Render("  " + line))
				}
				b.WriteString("\n")
			}

			if len(m.chapterList) > 15 {
				infoStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
				b.WriteString(infoStyle.Render(fmt.Sprintf("\n... và %d chapter khác\n", len(m.chapterList)-end+start)))
			}
		}

		b.WriteString("\n")
		footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
		b.WriteString(footerStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Đọc"))
		return b.String()
	}

	b.WriteString(headerStyle.Render("📚 TÌM KIẾM MANGA"))
	b.WriteString("\n\n")

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#9f75ff")).
		Padding(0, 1).
		Width(50)

	b.WriteString("Nhập tên manga: ")
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
	} else if len(m.mangaList) == 0 {
		hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
		b.WriteString(hintStyle.Render("Nhập tên manga (ít nhất 2 ký tự) và nhấn Enter để tìm kiếm..."))
	} else {
		listStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f0f4ff"))
		selectedListStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9f75ff")).
			Bold(true)

		for i, manga := range m.mangaList {
			if i == m.index {
				b.WriteString(selectedListStyle.Render(fmt.Sprintf("▶ %s", manga.Name)))
			} else {
				b.WriteString(listStyle.Render(fmt.Sprintf("  %s", manga.Name)))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
	b.WriteString(footerStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Chọn"))

	return b.String()
}
