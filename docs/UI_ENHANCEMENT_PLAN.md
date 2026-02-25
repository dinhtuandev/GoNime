# GoNime UI/UX Enhancement Plan

## 🎨 Design Vision

**Modern Dark Theme** như Spotify/Discord với:
- Gradient backgrounds
- Neon accent colors
- Glassmorphism effects
- Smooth animations

**Side-by-Side Layout** như Netflix:
- Poster/banner bên trái
- Details/info bên phải
- Full-width search bar
- Better content preview

---

## 📐 Layout Specification

### Trending Screen (Manga/Anime)

```
┌────────────────────────────────────────────────────────────────────┐
│  🎌 GoNime                                         [📺 Anime] [m] │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  ┌──────────────────┐  ┌────────────────────────────────────────┐ │
│  │                  │  │ 🔥 Xu hướng                            │ │
│  │   [POSTER]       │  │                                        │ │
│  │   240x320        │  │  ▶ Naruto Shippuden         (2007)     │ │
│  │                  │  │    One Piece                 (1999)    │ │
│  │                  │  │    Dragon Ball Super         (2015)     │ │
│  │                  │  │    Attack on Titan          (2013)     │ │
│  │                  │  │    My Hero Academia         (2016)     │ │
│  │                  │  │    Demon Slayer             (2019)     │ │
│  └──────────────────┘  └────────────────────────────────────────┘ │
│                                                                    │
│  [m] Đổi nguồn  [s] Tìm kiếm  [Enter] Xem  [↑↓] Di chuyển     │
└────────────────────────────────────────────────────────────────────┘
```

### Search Screen

```
┌────────────────────────────────────────────────────────────────────┐
│  🎌 GoNime                                         [📺 Anime] [m] │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  ┌──────────────────────────────────────────────────────────────┐  │
│  │  🔍 Tìm kiếm anime... ______________________________________ │  │
│  └──────────────────────────────────────────────────────────────┘  │
│                                                                    │
│  Kết quả: 14                                                     │
│                                                                    │
│  ┌──────────────────┐  ┌────────────────────────────────────────┐ │
│  │                  │  │  ▶ Naruto: Cuộc Chiến Ở Tuyết Quốc  │ │
│  │   [POSTER]       │  │     Naruto Shippuden                 │ │
│  │   180x240        │  │     Naruto                           │ │
│  │                  │  │                                        │ │
│  │  ⭐ 8.5         │  │  Thể loại: Action, Adventure         │ │
│  │  📺 220 tập    │  │  Năm: 2002                          │ │
│  └──────────────────┘  └────────────────────────────────────────┘ │
│                                                                    │
│  [Enter] Chọn  [↑↓] Di chuyển  [Esc] Quay lại                   │
└────────────────────────────────────────────────────────────────────┘
```

### Episodes Screen

```
┌────────────────────────────────────────────────────────────────────┐
│  🎌 GoNime                                         [📺 Anime] [m] │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  ┌──────────────────┐  ┌────────────────────────────────────────┐ │
│  │                  │  │ 📺 Naruto Shippuden                   │ │
│  │   [POSTER]       │  │                                        │ │
│  │   240x320        │  │ 📅 2007  │  ⏱ 24min  │  HD         │ │
│  │                  │  │                                   │ │
│  │                  │  │ ┌──────────────────────────────────┐ │ │
│  │                  │  │ │ Tập 1-50  │  Tập 51-100  │ ... │ │ │
│  │                  │  │ └──────────────────────────────────┘ │ │
│  │                  │  │                                        │ │
│  └──────────────────┘  │  ▶ Tập 1      Tập 2      Tập 3     │ │
│                        │    Tập 4      Tập 5      Tập 6     │ │
│                        │    Tập 7      Tập 8      Tập 9     │ │
│                        └────────────────────────────────────────┘ │
│                                                                    │
│  [Enter] Phát  [n/p] Tập tiếp/trước  [s] Server  [f] Yêu thích │
└────────────────────────────────────────────────────────────────────┘
```

---

## 🎨 Color Palette

### Primary Colors

| Name | Hex | Usage |
|------|-----|-------|
| Background Dark | `#0f0f0f` | Main background |
| Background Medium | `#181818` | Card background |
| Background Light | `#282828` | Hover states |
| Surface | `#1e1e1e` | Input fields |

### Accent Colors (Neon)

| Name | Hex | Usage |
|------|-----|-------|
| Primary | `#1db954` | Spotify green - buttons, highlights |
| Secondary | `#9b59b6` | Purple - manga accent |
| Tertiary | `#3498db` | Blue - anime accent |
| Gradient Start | `#e91e63` | Pink - gradients |
| Gradient End | `#9c27b0` | Purple - gradients |

### Text Colors

| Name | Hex | Usage |
|------|-----|-------|
| Text Primary | `#ffffff` | Main text |
| Text Secondary | `#b3b3b3` | Subtitles |
| Text Muted | `#727272` | Hints |

---

## 🔧 Implementation Plan

### Phase 1: Core UI Components

1. **Create theme package**
   - Define color palette constants
   - Create style builders
   - Add gradient support

2. **Create layout components**
   - SidePanel (poster + details)
   - ContentList (browse items)
   - SearchBar (input with styling)
   - StatusBar (footer hints)

### Phase 2: Screen Updates

1. **Update Trending Screen**
   - Add poster preview on left
   - Show item details on selection
   - Add source indicator badge

2. **Update Search Screen**
   - Full-width search bar with glow effect
   - Side-by-side results layout
   - Preview on selection

3. **Update Episodes Screen**
   - Large poster on left
   - Episode grid on right
   - Server selector

### Phase 3: Visual Enhancements

1. **Add animations**
   - Selection highlight transitions
   - Loading spinners
   - Screen transitions

2. **Add effects**
   - Gradient borders on selection
   - Glow effects on active items
   - Subtle shadows

---

## 📁 File Structure

```
internal/
├── ui/
│   ├── theme/
│   │   ├── colors.go        # Color palette
│   │   ├── styles.go        # Style builders
│   │   └── effects.go       # Gradients, glows
│   ├── components/
│   │   ├── panel.go        # SidePanel component
│   │   ├── list.go         # ContentList component
│   │   ├── search.go       # SearchBar component
│   │   └── status.go       # StatusBar component
│   └── layouts/
│       ├── trending.go      # Trending layout
│       ├── search.go       # Search layout
│       └── detail.go       # Detail layout
```

---

## ✅ Implementation Checklist

- [ ] Create theme package with colors
- [ ] Create style builders with gradients
- [ ] Update anime trending screen layout
- [ ] Update anime search screen layout
- [ ] Update anime episodes screen layout
- [ ] Update manga trending screen layout
- [ ] Update manga search screen layout
- [ ] Update manga chapters screen layout
- [ ] Add selection animations
- [ ] Add loading indicators
- [ ] Test all screens

---

## 🚀 Quick Start

After this plan is implemented:

```bash
# Run the app
go run ./cmd/gonime

# Navigate with arrow keys
# Press Enter to select
# Press m to switch sources
# Press s to search
```

---

*Last updated: 2026*
