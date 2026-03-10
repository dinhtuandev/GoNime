//go:build !windows

package player

import "os/exec"

// setDetachedProcess is a no-op on non-Windows platforms.
func setDetachedProcess(cmd *exec.Cmd) {
	// no-op
}
