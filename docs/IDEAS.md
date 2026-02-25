# GoNime Ideas & Future Features

Dựa trên dự án Sudachi và các ý tưởng phát triển thêm cho GoNime.

## 🎬 Movie/Anime Watching Features

### 1. Tích hợp API Phim (OPhim, NguonC, PhimAPI)

**Mô tả:** Thêm chức năng xem phim anime sử dụng các API có sẵn từ Sudachi

**API Sources:**
- **OPhim** (ophim.cc) - API chính
- **NguonC** (nguonc.com) - API dự phòng  
- **PhimAPI** (phimapi.com) - API bổ sung

**Tính năng:**
- [ ] Tìm kiếm phim anime
- [ ] Xem phim trực tuyến với MPV/VLC
- [ ] Danh sách tập phim
- [ ] Chọn chất lượng video
- [ ] Tải phim với yt-dlp

**Implementation:**
```
internal/
  api/
    ophim.go      # OPhim API client
    nguonc.go     # NguonC API client
    phimapi.go    # PhimAPI client
```

### 2. Download Manager

**Mô tả:** Quản lý tải phim/truyện với hỗ trợ multi-thread

**Tính năng:**
- [ ] Tải phim anime (yt-dlp)
- [ ] Tải truyện tranh (images)
- [ ] Queue management
- [ ] Tốc độ tải (aria2c integration)
- [ ] Lịch sử tải

### 3. Anime Streaming Server

**Mô tả:** Local streaming server để stream anime không cần player ngoài

**Tính năng:**
- [ ] HTTP server for streaming
- [ ] Transcode video (nếu cần)
- [ ] Subtitle support
- [ ] Quality selection

## 📚 Manga/Comic Features

### 4. Multiple Source Support

**Mô tả:** Thêm nhiều nguồn truyện ngoài truyenvn

**Nguồn tiềm năng:**
- [ ] Nettruyen
- [ ] TruyenQQ
- [ ] MangaDex
- [ ] COMICWWW



## 🎮 UI/UX Features

### 6. TUI Enhancements

**Tính năng:**
- [ ] Fuzzy search với fzf-style
- [ ] Preview images trong menu
- [ ] Color schemes/Themes
- [ ] Keybindings tùy chỉnh
- [ ] History/Favorites management

### 7. Media Server Integration

**Tính năng:**
- [ ] Jellyfin/Plex sync
- [ ] Streaming qua local network
- [ ] Cast to Chromecast (nếu có thể)

## 🔧 Technical Features

### 8. Caching & Performance

**Tính năng:**
- [ ] Image cache với LRU
- [ ] Prefetch chapters
- [ ] Offline mode
- [ ] Database cho metadata

### 9. Plugin System

**Mô tả:** Hệ thống plugin để mở rộng

**Tính năng:**
- [ ] Source plugins
- [ ] Player plugins
- [ ] UI themes plugins

### 10. Cross-Platform Features

**Tính năng:**
- [ ] Windows notifications
- [ ] System tray
- [ ] Global hotkeys
- [ ] Auto-update

## 📋 Implementation Priority

### Phase 1 (Cao ưu tiên)
1. ✅ Manga Viewer với Wails - DONE
2. ⏳ Movie API Integration (OPhim/NguonC)
3. ⏳ Download Manager
4. ⏳ Chapter navigation trong viewer

### Phase 2 (Trung ưu tiên)
5. Multiple manga sources
6. Reader enhancements
7. TUI improvements
8. Caching system

### Phase 3 (Thấp ưu tiên)
9. Plugin system
10. Media server integration
11. Cross-platform features

## 🔗 Related Projects

- [Sudachi](https://github.com/KabosuNeko/sudachi) - CLI anime player tham khảo
- [Mpv](https://mpv.io/) - Media player
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - Video downloader
- [chafa](https://github.com/hpjansson/chafa) - Image to ASCII

## 💡 Ideas Notes

- Sử dụng cùng API như Sudachi để tương thích
- Giữ kiến trúc CLI orchestrator đã có
- Wails viewer có thể mở rộng cho movie player
- Có thể tách riêng module cho từng source

---

*Document created: 2026*
*Last updated: 2026*
