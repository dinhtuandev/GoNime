package theme

import "github.com/charmbracelet/lipgloss"

// ─── Core Structural Styles ───────────────────────────────────────────────────

var (
	// AppTitle: used for the big logo text
	AppTitle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	// Header: screen header (module name)
	Header = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)

	// ScreenTitle: large section heading
	ScreenTitle = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Bold(true)

	// Subtitle: secondary info line under title
	Subtitle = lipgloss.NewStyle().
			Foreground(TextSecondary)

	// Muted: dim helper text
	Muted = lipgloss.NewStyle().
		Foreground(TextMuted)
)

// ─── List Items ───────────────────────────────────────────────────────────────

var (
	// Selected: highlighted list row — cyan text, dark bg
	Selected = lipgloss.NewStyle().
			Foreground(Primary).
			Background(BackgroundLight).
			Bold(true)

	// Normal: unselected list row
	Normal = lipgloss.NewStyle().
		Foreground(TextSecondary)
)

// ─── Badges ───────────────────────────────────────────────────────────────────

var (
	// Badge: inactive tab/source badge
	Badge = lipgloss.NewStyle().
		Foreground(TextMuted).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(TextDim).
		Padding(0, 1)

	// BadgeActive: active tab with bright border
	BadgeActive = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(0, 1)

	// BadgePurple: manga-tagged badge
	BadgePurple = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Secondary).
			Padding(0, 1)

	// BadgeGold: server/warning badge
	BadgeGold = lipgloss.NewStyle().
			Foreground(Gold).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Gold).
			Padding(0, 1)
)

// ─── Input ────────────────────────────────────────────────────────────────────

var (
	// InputLabel: label above the input box
	InputLabel = lipgloss.NewStyle().
			Foreground(TextSecondary)

	// InputText: text inside input
	InputText = lipgloss.NewStyle().
			Foreground(TextPrimary)

	// InputCursor: cursor character style
	InputCursor = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)
)

// ─── Status / Feedback ────────────────────────────────────────────────────────

var (
	ErrorText   = lipgloss.NewStyle().Foreground(Error).Bold(true)
	SuccessText = lipgloss.NewStyle().Foreground(Success).Bold(true)
	WarningText = lipgloss.NewStyle().Foreground(Warning).Bold(true)

	// Hint: keyboard shortcut hints in status bar
	Hint = lipgloss.NewStyle().
		Foreground(TextMuted)

	// HintKey: the key part (e.g. "Enter") in bright color
	HintKey = lipgloss.NewStyle().
		Foreground(Gold).
		Bold(true)
)

// ─── Episode / Chapter Cells ──────────────────────────────────────────────────

var (
	Episode = lipgloss.NewStyle().
		Foreground(TextSecondary).
		Background(SurfaceHigh).
		Padding(0, 1)

	EpisodeActive = lipgloss.NewStyle().
			Foreground(BackgroundDark).
			Background(Primary).
			Bold(true).
			Padding(0, 1)

	EpisodePlaying = lipgloss.NewStyle().
			Foreground(Primary).
			Background(BackgroundLight).
			Bold(true).
			Padding(0, 1)

	Chapter = lipgloss.NewStyle().
		Foreground(TextSecondary)

	ChapterActive = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	// Loading: spinner / loading indicator
	Loading = lipgloss.NewStyle().
		Foreground(Primary).
		Bold(true)
)

// ─── Cards ────────────────────────────────────────────────────────────────────

var (
	Card = lipgloss.NewStyle().
		Background(BackgroundMed).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(BackgroundLight).
		Padding(1, 2)

	CardSelected = lipgloss.NewStyle().
			Background(BackgroundLight).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Primary).
			Padding(1, 2)
)
