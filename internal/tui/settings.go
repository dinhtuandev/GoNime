package tui

import (
	"fmt"
	"strings"

	"gonime/internal/config"
	"gonime/internal/player"
	"gonime/internal/render"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type settingsField int

const (
	settingsFieldPlayer settingsField = iota
	settingsFieldDownloadDir
	settingsFieldMangaDir
	settingsFieldMangaImageSize
	settingsFieldMangaURL
)

type SettingsModel struct {
	parent       MainModel
	cfg          *config.Config
	manager      *config.Manager
	index        int
	editingField settingsField
	isEditing    bool
	editBuffer   string
	err          error
	statusMsg    string

	// Trạng thái dependency
	hasMPV   bool
	hasVLC   bool
	hasChafa bool
}

var settingsLabels = []string{
	"🎬 Trình phát video",
	"📁 Thư mục tải xuống",
	"🖼  Thư mục ảnh manga",
	"📐 Kích thước ảnh manga (%)",
	"🌐 URL nguồn Manga",
}

func NewSettingsModel(parent MainModel) SettingsModel {
	mgr, _ := config.NewManager()
	cfg, _ := mgr.LoadConfig()
	m := SettingsModel{
		parent:  parent,
		cfg:     cfg,
		manager: mgr,
	}
	m.checkDependencies()
	return m
}

func (m *SettingsModel) checkDependencies() {
	// Check players using internal logic
	if p, err := player.NewPlayer("mpv"); err == nil && p != nil {
		m.hasMPV = true
	}
	if p, err := player.NewPlayer("vlc"); err == nil && p != nil {
		m.hasVLC = true
	} else if p, err := player.NewPlayer("vlc.exe"); err == nil && p != nil {
		m.hasVLC = true
	}

	// Check chafa
	m.hasChafa = render.IsAvailable()
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
	case settingsFieldMangaURL:
		return m.cfg.MangaURL
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

	case settingsFieldMangaURL:
		if val == "" {
			return fmt.Errorf("URL không được để trống")
		}
		m.cfg.MangaURL = val
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
					runes := []rune(m.editBuffer)
					m.editBuffer = string(runes[:len(runes)-1])
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
		Bold(true)

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9f75ff")).
		Bold(true)

	normalStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#f0f4ff"))

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#606880"))

	errStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff6b6b"))

	successStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4caf50"))

	statusLabelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#a0a8c8"))

	b.WriteString("\n")
	b.WriteString(headerStyle.Render("  ⚙️ CÀI ĐẶT"))
	b.WriteString("\n\n")

	if m.cfg == nil {
		b.WriteString("  " + errStyle.Render("❌ Không thể tải cấu hình."))
		return b.String()
	}

	// ─── PHẦN KIỂM TRA HỆ THỐNG ───
	b.WriteString(statusLabelStyle.Render("  🔍 KIỂM TRA HỆ THỐNG:"))
	b.WriteString("\n")

	checkItem := func(name string, has bool) string {
		if has {
			return successStyle.Render(fmt.Sprintf("    ✓ %-8s: Sẵn sàng", name))
		}
		return errStyle.Render(fmt.Sprintf("    ✗ %-8s: Chưa thấy", name))
	}

	b.WriteString(checkItem("mpv", m.hasMPV))
	b.WriteString("\n")
	b.WriteString(checkItem("VLC", m.hasVLC))
	b.WriteString("\n")
	b.WriteString(checkItem("chafa", m.hasChafa))
	if !m.hasChafa {
		b.WriteString(hintStyle.Render(" (để xem manga)"))
	}
	b.WriteString("\n\n")

	// ─── PHẦN DANH SÁCH CÀI ĐẶT ───
	values := []string{
		m.cfg.Player,
		m.cfg.DownloadDir,
		m.cfg.MangaImageDir,
		fmt.Sprintf("%d%%", m.cfg.MangaImageSize),
		m.cfg.MangaURL,
	}

	for i, label := range settingsLabels {
		isSelected := i == m.index

		if isSelected && m.isEditing {
			// Đang chỉnh sửa
			b.WriteString(drawInputBox(label, m.editBuffer+"▌", 56, "#9f75ff"))
			b.WriteString("\n")
		} else {
			// Hiển thị bình thường
			marker := "  "
			titleColor := normalStyle
			valColor := hintStyle

			if isSelected {
				marker = "▶ "
				titleColor = selectedStyle
				valColor = lipgloss.NewStyle().Foreground(lipgloss.Color("#33a5f5"))
			}

			b.WriteString(fmt.Sprintf("  %s%s\n", marker, titleColor.Render(label)))
			b.WriteString(fmt.Sprintf("    %s\n", valColor.Render(values[i])))
			b.WriteString("\n")
		}
	}

	if m.err != nil {
		b.WriteString("  " + errStyle.Render(fmt.Sprintf("❌ %v", m.err)))
		b.WriteString("\n")
	}

	if m.statusMsg != "" {
		b.WriteString("  " + successStyle.Render(m.statusMsg))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	if m.isEditing {
		b.WriteString(hintStyle.Render("  Enter: Lưu  │  Esc: Hủy  │  Backspace: Xóa"))
	} else {
		b.WriteString(hintStyle.Render("  Esc: Quay lại  │  ↑↓/jk: Điều hướng  │  Enter: Chỉnh sửa  │  r: Mặc định"))
	}

	return b.String()
}
