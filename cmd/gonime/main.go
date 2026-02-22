package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"gonime/internal/tui"
)

var version = "1.0.0"

func main() {
	var rootCmd = &cobra.Command{
		Use:   "gonime",
		Short: "GoNime - Trình xem Anime và Manga trên Terminal",
		Long: `GoNime là một ứng dụng CLI native cho phép xem anime, đọc manga trực tiếp trên terminal.
Lấy cảm hứng từ anisub, mang đến trải nghiệm nhanh, nhẹ và mạnh mẽ.`,
		Run: runTUI,
	}

	rootCmd.AddCommand(
		NewAnimeCommand(),
		NewMangaCommand(),
		NewSettingsCommand(),
		NewVersionCommand(),
	)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runTUI(cmd *cobra.Command, args []string) {
	p := tea.NewProgram(tui.NewMainModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Lỗi: %v\n", err)
		os.Exit(1)
	}
}

func NewAnimeCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "anime [tên anime]",
		Short: "Tìm kiếm và xem anime",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			query := args[0]
			cmd.Printf("Đang tìm kiếm: %s\n", query)
			// TODO: Implement anime search
		},
	}
}

func NewMangaCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "manga [tên manga]",
		Short: "Tìm kiếm và đọc manga",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			query := args[0]
			cmd.Printf("Đang tìm kiếm: %s\n", query)
			// TODO: Implement manga search
		},
	}
}

func NewSettingsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "Cài đặt ứng dụng",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Mở cài đặt...")
			// TODO: Open settings TUI
		},
	}
}

func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Hiển thị thông tin phiên bản",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("GoNime v%s\n", version)
		},
	}
}
