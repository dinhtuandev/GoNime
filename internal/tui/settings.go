package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gonime/internal/config"
)

type settingsField int

const (
	settingsFieldPlayer settingsField = iota
	settingsFieldDownloadDir
	settingsFieldMangaDir
	settingsFieldMangaImageSize
)

type SettingsModel struct {
	parent      MainModel
	cfg         *config.Config
	manager     *config.Manager
	index       int
	editingField settingsField
	isEditing   bool
	editBuffer  string
	err         error
	statusMsg   string
}

var settingsLabels = []string{
	"🎬 Trình phát video",
	"📁 Thư mục tải xuống",
	"🖼  Thư mục ảnh manga",
	"📐 Kích thước ảnh manga (%)",
}

func NewSettingsModel(parent MainModel) SettingsModel {
	mgr, _ := config.NewManager()
	cfg, _ := mgr.LoadConfig()
	return SettingsModel{
		parent:  parent,
		cfg:     cfg,
		manager: mgr,
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) getCurrentValue() string {
	if m.cfg == nil {
		return ""
	}
	switch settingsField(m.index) {
	case settingsFieldPlayer:
		return m.cfg.Player
	case settingsFieldDownloadDir:
		return m.cfg.DownloadDir
	case settingsFieldMangaDir:
		return m.cfg.MangaImageDir
	case settingsFieldMangaImageSize:
		return fmt.Sprintf("%d", m.cfg.MangaImageSize)
	}
	return ""
}

func (m *SettingsModel) applyEdit() error {
	if m.cfg == nil {
		return fmt.Errorf("cấu hình không khả dụng")
	}

	val := strings.TrimSpace(m.editBuffer)
	switch settingsField(m.index) {
	case settingsFieldPlayer:
		if val == "" {
			return fmt.Errorf("tên trình phát không được để trống")
		}
		m.cfg.Player = val

	case settingsFieldDownloadDir:
		if val == "" {
			return fmt.Errorf("thư mục tải không được để trống")
		}
		m.cfg.DownloadDir = val

	case settingsFieldMangaDir:
		if val == "" {
			return fmt.Errorf("thư mục manga không được để trống")
		}
		m.cfg.MangaImageDir = val

	case settingsFieldMangaImageSize:
		var size int
		if _, err := fmt.Sscanf(val, "%d", &size); err != nil {
			return fmt.Errorf("vui lòng nhập số từ 50 đến 200")
		}
		if size < 50 || size > 200 {
			return fmt.Errorf("kích thước phải từ 50 đến 200")
		}
		m.cfg.MangaImageSize = size
	}

	// Lưu config
	if err := m.manager.SaveConfig(m.cfg); err != nil {
		return fmt.Errorf("không lưu được cấu hình: %v", err)
	}
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		// --- Khi đang chỉnh sửa ---
		if m.isEditing {
			switch msg.String() {
			case "enter":
				// Lưu giá trị
				if err := m.applyEdit(); err != nil {
					m.err = err
				} else {
					m.statusMsg = fmt.Sprintf("✅ Đã lưu: %s", settingsLabels[m.index])
					m.err = nil
				}
				m.isEditing = false
				m.editBuffer = ""

			case "esc":
				// Hủy chỉnh sửa
				m.isEditing = false
				m.editBuffer = ""
				m.err = nil
				m.statusMsg = ""

			case "backspace":
				if len(m.editBuffer) > 0 {
					m.editBuffer = m.editBuffer[:len(m.editBuffer)-1]
				}

			case "ctrl+c":
				return m, tea.Quit

			default:
				if len(msg.String()) == 1 {
					m.editBuffer += msg.String()
				}
			}
			return m, nil
		}

		// --- Khi không chỉnh sửa ---
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
			if m.index < len(settingsLabels)-1 {
				m.index++
			}
		case "enter", " ":
			// Bắt đầu chỉnh sửa trường được chọn
			m.isEditing = true
			m.editingField = settingsField(m.index)
			m.editBuffer = m.getCurrentValue()
			m.err = nil
			m.statusMsg = ""
		case "r":
			// Reset về mặc định
			if m.cfg != nil && m.manager != nil {
				defaultCfg := config.GetDefaultConfigForManager(m.manager)
				if err := m.manager.SaveConfig(defaultCfg); err == nil {
					m.cfg = defaultCfg
					m.statusMsg = "✅ Đã đặt lại về cấu hình mặc định"
				}
			}
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

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9f75ff")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f0f4ff"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#33a5f5"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a0a8c8"))

	errStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff6b6b"))

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4caf50"))

	editingInputStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#9f75ff")).
		Padding(0, 1).
		Width(50)

	b.WriteString(headerStyle.Render("⚙️ CÀI ĐẶT"))
	b.WriteString("\n\n")

	if m.cfg == nil {
		b.WriteString(errStyle.Render("❌ Không thể tải cấu hình."))
		return b.String()
	}

	values := []string{
		m.cfg.Player,
		m.cfg.DownloadDir,
		m.cfg.MangaImageDir,
		fmt.Sprintf("%d%%", m.cfg.MangaImageSize),
	}

	for i, label := range settingsLabels {
		isSelected := i == m.index

		if isSelected && m.isEditing {
			// Đang chỉnh sửa trường này
			b.WriteString(selectedStyle.Render("▶ " + label))
			b.WriteString("\n")
			b.WriteString("  ")
			b.WriteString(editingInputStyle.Render(m.editBuffer + "▌"))
			b.WriteString("\n")
		} else if isSelected {
			b.WriteString(selectedStyle.Render(fmt.Sprintf("▶ %s", label)))
			b.WriteString("\n")
			b.WriteString("  ")
			b.WriteString(valueStyle.Render(values[i]))
			b.WriteString("\n")
		} else {
			b.WriteString(normalStyle.Render(fmt.Sprintf("  %s", label)))
			b.WriteString("\n")
			b.WriteString("  ")
			b.WriteString(hintStyle.Render(values[i]))
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}

	if m.err != nil {
		b.WriteString(errStyle.Render(fmt.Sprintf("❌ %v", m.err)))
		b.WriteString("\n")
	}

	if m.statusMsg != "" {
		b.WriteString(successStyle.Render(m.statusMsg))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	if m.isEditing {
		b.WriteString(hintStyle.Render("Enter: Lưu | Esc: Hủy | Backspace: Xóa ký tự"))
	} else {
		b.WriteString(hintStyle.Render("Esc: Quay lại | ↑/↓: Điều hướng | Enter: Chỉnh sửa | r: Đặt lại mặc định"))
	}

	return b.String()
}
