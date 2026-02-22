package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gonime/internal/config"
)

type SettingsModel struct {
	parent  MainModel
	cfg     *config.Config
	manager *config.Manager
	index   int
	items   []string
}

func NewSettingsModel(parent MainModel) SettingsModel {
	mgr, _ := config.NewManager()
	cfg, _ := mgr.LoadConfig()
	return SettingsModel{
		parent:  parent,
		cfg:     cfg,
		manager: mgr,
		items: []string{
			fmt.Sprintf("Trình phát video: %s", cfg.Player),
			fmt.Sprintf("Thư mục tải: %s", cfg.DownloadDir),
			fmt.Sprintf("Thư mục manga: %s", cfg.MangaImageDir),
			fmt.Sprintf("Kích thước ảnh manga: %d%%", cfg.MangaImageSize),
		},
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
			// TODO: Implement setting edit
		}
	}

	return m, nil
}

func (m SettingsModel) View() string {
	var b strings.Builder

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#33a5f5")).
		Bold(true).
		Margin(1, 0)

	b.WriteString(headerStyle.Render("⚙️ CÀI ĐẶT"))
	b.WriteString("\n\n")

	listStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f0f4ff"))
	selectedListStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9f75ff")).
		Bold(true)

	for i, item := range m.items {
		if i == m.index {
			b.WriteString(selectedListStyle.Render(fmt.Sprintf("▶ %s", item)))
		} else {
			b.WriteString(listStyle.Render(fmt.Sprintf("  %s", item)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	footerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#a0a8c8"))
	b.WriteString(footerStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Chỉnh sửa"))

	return b.String()
}
