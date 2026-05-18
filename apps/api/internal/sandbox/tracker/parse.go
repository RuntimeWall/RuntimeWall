package tracker

import (
	"regexp"
	"strings"
)

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)

// ParseInput extracts completed commands from terminal keystrokes.
func ParseInput(buf *[]byte, data []byte) []string {
	var submitted []string

	for _, c := range data {
		switch c {
		case '\r', '\n':
			line := cleanCommand(string(*buf))
			*buf = (*buf)[:0]
			if line != "" {
				submitted = append(submitted, line)
			}
		case 127, 8: // backspace / delete
			if len(*buf) > 0 {
				*buf = (*buf)[:len(*buf)-1]
			}
		case 3, 21: // ctrl+c, ctrl+u — clear line
			*buf = (*buf)[:0]
		case 4: // ctrl+d — ignore
		default:
			if c >= 32 && c < 127 {
				*buf = append(*buf, c)
			}
		}
	}

	return submitted
}

func cleanCommand(line string) string {
	line = ansiEscape.ReplaceAllString(line, "")
	line = strings.TrimSpace(line)
	// Drop bare prompts accidentally captured.
	if line == "" || strings.HasSuffix(line, "$") || strings.HasSuffix(line, "#") {
		return ""
	}
	return line
}
