// Package sanitize removes terminal control sequences from untrusted text.
package sanitize

import (
	"regexp"
	"strings"
)

var escapeSeq = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]|\x1b\][^\x07\x1b]*(?:\x07|\x1b\\)|\x1b[@-Z\\-_]`)

// Text strips ANSI/OSC escape sequences and control characters from remote
// text before it reaches the terminal. Newlines and tabs are kept.
func Text(s string) string {
	if s == "" {
		return s
	}
	s = escapeSeq.ReplaceAllString(s, "")
	return strings.Map(func(r rune) rune {
		switch {
		case r == '\n' || r == '\t':
			return r
		case r == 0x7f:
			return -1
		case r < 0x20:
			return -1
		case r >= 0x80 && r <= 0x9f:
			return -1
		default:
			return r
		}
	}, s)
}
