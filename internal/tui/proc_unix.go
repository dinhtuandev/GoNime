//go:build !windows

package tui

import "os/exec"

// setDetachedProcess is a no-op on non-Windows platforms.
// On Linux/macOS, the child process naturally inherits a separate process group
// and can create GUI windows without detaching from the parent console.
func setDetachedProcess(cmd *exec.Cmd) {
	// no-op on non-Windows
}
