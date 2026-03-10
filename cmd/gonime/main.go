package main

import (
	"fmt"
	"os"
	"strings"

	"gonime/internal/api"
	"gonime/internal/config"
	"gonime/internal/player"
	"gonime/internal/tui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

// Build info — injected by ldflags during release builds:
//
//	go build -ldflags "-X main.version=v1.2.3 -X main.commit=abc1234 -X main.buildDate=2026-03-08T14:00:00Z"
var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gonime",
		Short: "GoNime - Trình xem Anime và Manga trên Terminal",
		Long: `GoNime là một ứng dụng CLI native cho phép xem anime, đọc manga trực tiếp trên terminal.
Lấy cảm hứng từ anisub, mang đến trải nghiệm nhanh, nhẹ và mạnh mẽ trên Windows.`,
		Run: runTUI,
	}

	rootCmd.AddCommand(
		NewAnimeCommand(),
		NewMangaCommand(),
		NewSettingsCommand(),
		NewVersionCommand(),
		NewHistoryCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// runTUI khởi động giao diện TUI chính
func runTUI(cmd *cobra.Command, args []string) {
	p := tea.NewProgram(
		tui.NewMainModel(),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi: %v\n", err)
		os.Exit(1)
	}
}

// NewAnimeCommand tìm kiếm anime và phát tập đầu tiên
func NewAnimeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "anime [tên anime]",
		Short: "Tìm kiếm và xem anime",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")
			fmt.Printf("🔍 Đang tìm kiếm: %s\n", query)

			client := api.NewAnimeClient()
			items, _, err := client.Search(query)
			if err != nil {
				return fmt.Errorf("lỗi tìm kiếm: %v", err)
			}
			if len(items) == 0 {
				fmt.Println("Không tìm thấy kết quả nào.")
				return nil
			}

			fmt.Printf("\nTìm thấy %d kết quả:\n", len(items))
			for i, item := range items {
				year := ""
				if item.Year > 0 {
					year = fmt.Sprintf(" (%d)", item.Year)
				}
				fmt.Printf("  [%d] %s%s\n", i+1, item.Name, year)
			}

			// Mở TUI với kết quả search sẵn
			fmt.Println("\nMở giao diện TUI để chọn...")
			p := tea.NewProgram(
				tui.NewMainModel(),
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)
			_, err = p.Run()
			return err
		},
	}
}

// NewMangaCommand tìm kiếm manga
func NewMangaCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "manga [tên manga]",
		Short: "Tìm kiếm và đọc manga",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := strings.Join(args, " ")
			fmt.Printf("🔍 Đang tìm kiếm manga: %s\n", query)
			fmt.Println("Mở giao diện TUI để chọn...")

			// Mở TUI - người dùng sẽ search trong TUI
			p := tea.NewProgram(
				tui.NewMainModel(),
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)
			_, err := p.Run()
			return err
		},
	}
}

// NewSettingsCommand mở màn Settings trực tiếp
func NewSettingsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "Quản lý cấu hình ứng dụng",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := config.NewManager()
			if err != nil {
				return fmt.Errorf("không tải được config manager: %v", err)
			}
			cfg, _ := mgr.LoadConfig()

			fmt.Println("⚙️  Cấu hình hiện tại:")
			fmt.Printf("  Trình phát video : %s\n", cfg.Player)
			fmt.Printf("  Thư mục tải      : %s\n", cfg.DownloadDir)
			fmt.Printf("  Thư mục manga    : %s\n", cfg.MangaImageDir)
			fmt.Printf("  Kích thước ảnh   : %d%%\n", cfg.MangaImageSize)
			fmt.Println("\nMở giao diện TUI để chỉnh sửa...")

			p := tea.NewProgram(
				tui.NewMainModel(),
				tea.WithAltScreen(),
				tea.WithMouseCellMotion(),
			)
			_, err = p.Run()
			return err
		},
	}
}

// NewVersionCommand hiển thị version
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Hiển thị thông tin phiên bản",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("GoNime %s\n", version)
			fmt.Printf("  Commit    : %s\n", commit)
			fmt.Printf("  Built     : %s\n", buildDate)
			fmt.Println("  Stack     : Go · bubbletea · lipgloss · cobra")
		},
	}
}

// NewHistoryCommand xem / phát lại lịch sử xem
func NewHistoryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "history",
		Short: "Xem và quản lý lịch sử xem",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Liệt kê lịch sử xem",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := config.NewManager()
			if err != nil {
				return err
			}
			items, err := mgr.LoadHistory()
			if err != nil {
				return err
			}
			if len(items) == 0 {
				fmt.Println("Lịch sử trống.")
				return nil
			}
			fmt.Printf("📜 Lịch sử xem (%d mục):\n", len(items))
			for i, item := range items {
				icon := "🎬"
				if item.Type == "manga" {
					icon = "📚"
				}
				fmt.Printf("  [%d] %s %s - %s  (%s)\n", i+1, icon, item.Name, item.Episode, item.Timestamp)
			}
			return nil
		},
	}

	clearCmd := &cobra.Command{
		Use:   "clear",
		Short: "Xóa toàn bộ lịch sử",
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := config.NewManager()
			if err != nil {
				return err
			}
			if err := mgr.ClearHistory(); err != nil {
				return err
			}
			fmt.Println("✅ Đã xóa toàn bộ lịch sử.")
			return nil
		},
	}

	playCmd := &cobra.Command{
		Use:   "play [số thứ tự]",
		Short: "Phát lại mục trong lịch sử (mặc định: mục mới nhất)",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			mgr, err := config.NewManager()
			if err != nil {
				return err
			}
			items, err := mgr.LoadHistory()
			if err != nil {
				return err
			}
			if len(items) == 0 {
				fmt.Println("Lịch sử trống.")
				return nil
			}

			idx := 0
			if len(args) == 1 {
				if _, err := fmt.Sscanf(args[0], "%d", &idx); err != nil || idx < 1 || idx > len(items) {
					return fmt.Errorf("số thứ tự không hợp lệ (1-%d)", len(items))
				}
				idx-- // Convert to 0-based
			}

			item := items[idx]
			if item.URL == "" {
				return fmt.Errorf("không có link phát cho mục này")
			}

			cfg, _ := mgr.LoadConfig()
			playerName := cfg.Player

			p, err := player.NewPlayer(playerName)
			if err != nil {
				p, err = player.NewPlayer("vlc")
				if err != nil {
					return fmt.Errorf("không tìm thấy trình phát (mpv/vlc)")
				}
			}

			title := fmt.Sprintf("GoNime: %s - %s", item.Name, item.Episode)
			fmt.Printf("▶ Đang phát: %s - %s\n", item.Name, item.Episode)
			return p.Play(item.URL, title)
		},
	}

	cmd.AddCommand(listCmd, clearCmd, playCmd)
	return cmd
}
