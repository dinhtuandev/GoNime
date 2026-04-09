# Fix: VLC/Player Issues and Loop Errors

## Problems Fixed

### 1. VLC m3u8 Loop Error
**Issue**: VLC keeps looping when playing .m3u8 streams from some APIs

**Root Cause**: VLC handles m3u8 streams poorly without proper caching and network settings

**Fix Applied** (`internal/tui/anime_player_opts.go`):
- Added `--network-caching=5000` (increased from 3000)
- Added `--file-caching=3000` 
- Added `--live-caching=3000`
- Added `--deinterlace=auto` for better stream handling
- Separated User-Agent to be always sent (not conditional)

### 2. Player Switching Not Working
**Issue**: When switching from mpv to VLC (or vice versa), player doesn't start

**Root Cause**: Variable `err` was shared between `NewPlayer()` and `FindAnyPlayer()` calls, causing incorrect fallback logic

**Fix Applied** (`internal/tui/anime.go` line 605-621):
- Renamed `err` to `playerErr` to isolate error state

### 3. Stream URL Detection Too Strict
**Issue**: Many m3u8 URLs from APIs don't have `.m3u8` extension in URL

**Fix Applied** (`internal/tui/anime.go` `isDirectStreamURL`):
- Added pattern detection: `/hls/`, `/stream/`, `/dash/`, `cdn-`, etc.

### 4. No Fallback When Player Fails
**Issue**: If player fails to start, user just gets error with no alternative

**Fix Applied** (`internal/tui/anime.go` line 640-650):
- Added automatic fallback to browser when player fails

### 5. Added Logging for Debugging
**Fix Applied** (`internal/tui/anime_player_opts.go`):
- Added `log.Printf` to track player start attempts
- Logs player name, args, and success/failure

## Files Changed

1. `internal/tui/anime_player_opts.go` - VLC args, logging
2. `internal/tui/anime.go` - player error handling, stream detection, fallback logic