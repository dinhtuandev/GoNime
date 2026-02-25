package theme

import "github.com/charmbracelet/lipgloss"

// ─── Backgrounds ──────────────────────────────────────────────────────────────
var (
	BackgroundDark  = lipgloss.Color("#0d0d12") // near-black base
	BackgroundMed   = lipgloss.Color("#12121a") // card / panel bg
	BackgroundLight = lipgloss.Color("#1a1a26") // hover / selected bg
	Surface         = lipgloss.Color("#1e1e2e") // input field bg
	SurfaceHigh     = lipgloss.Color("#252538") // elevated surface
)

// ─── Neon Accents ─────────────────────────────────────────────────────────────
var (
	// Primary — electric cyan (anime accent, main highlights)
	Primary = lipgloss.Color("#00d4ff")

	// Secondary — neon purple (manga accent, badges)
	Secondary = lipgloss.Color("#bf5fcf")

	// Tertiary — hot pink (movie accent, error-adjacent)
	Tertiary = lipgloss.Color("#ff2d78")

	// Gold — amber/golden, used for servers, warnings, keys in hints
	Gold = lipgloss.Color("#ffb347")

	// Mint — softer green, for success indicators
	Mint = lipgloss.Color("#39d353")
)

// ─── Gradient Pair ────────────────────────────────────────────────────────────
var (
	GradientStart = lipgloss.Color("#00d4ff") // cyan
	GradientEnd   = lipgloss.Color("#bf5fcf") // purple
)

// ─── Text ─────────────────────────────────────────────────────────────────────
var (
	TextPrimary   = lipgloss.Color("#e8e8f0") // bright near-white
	TextSecondary = lipgloss.Color("#8888aa") // medium muted
	TextMuted     = lipgloss.Color("#4a4a6a") // very muted, hints
	TextDim       = lipgloss.Color("#2e2e4a") // dimmed decorators
)

// ─── Semantic ─────────────────────────────────────────────────────────────────
var (
	Error   = lipgloss.Color("#ff4757")
	Warning = lipgloss.Color("#ffb347")
	Success = lipgloss.Color("#39d353")
)

// ─── Module Accents ───────────────────────────────────────────────────────────
var (
	AnimeAccent = Primary   // cyan
	MangaAccent = Secondary // purple
	MovieAccent = Tertiary  // hot pink
)
