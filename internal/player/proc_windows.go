//go:build windows

package player

import (
	"os/exec"
	"syscall"
)

// setDetachedProcess configures the command to spawn as a detached process on Windows.
// Without DETACHED_PROCESS, mpv/vlc installed via Windows Store or winget (MSIX packages)
// may fail to create their GUI window when launched from a TUI application.
func setDetachedProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// CREATE_NEW_PROCESS_GROUP (0x00000200): isolates the child so Ctrl+C in TUI doesn't kill the player
		// DETACHED_PROCESS (0x00000008): detaches from parent console, lets mpv create its own GUI window
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP | 0x00000008,
	}
}
