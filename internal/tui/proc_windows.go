//go:build windows

package tui

import (
	"os/exec"
	"syscall"
)

// setDetachedProcess configures the command to spawn as a detached process on Windows.
// This is critical for Windows Store/MSIX apps (like mpv installed via winget) because:
// 1. CREATE_NEW_PROCESS_GROUP: prevents Ctrl+C in the TUI from killing the player
// 2. DETACHED_PROCESS: detaches from the parent console so the player can create its own window
// Without these flags, mpv launched from a TUI app may fail to start its GUI window.
func setDetachedProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008, // 0x00000008 = DETACHED_PROCESS
	}
}
