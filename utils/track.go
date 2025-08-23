package utils

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	TheLibrary Library
	TrackByID  map[string]*Track
)

type strconvErr struct{ msg string }

var trackRegex = regexp.MustCompile(`(?i)^(?:[^-]+\s-\s)?(\d{1,2})?\.?\s*(.+?)\s*\.[^.]+'?$`)

func (e *strconvErr) Error() string { return e.msg }

func strconvAtoi(s string) (int, error) {
	var n int
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			return 0, &strconvErr{msg: "invalid syntax"}
		}
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func ParseTrackInfo(fileName string) (int, string) {
	base := filepath.Base(fileName)
	m := trackRegex.FindStringSubmatch(base)
	if len(m) == 0 {
		return 0, strings.TrimSuffix(base, filepath.Ext(base))
	}

	var num int
	if m[1] != "" {
		for _, ch := range m[1] {
			if ch < '0' || ch > '9' {
				num = 0
				break
			}
		}

		if n, err := strconvAtoi(m[1]); err == nil {
			num = n
		}
	}
	return num, strings.TrimSpace(m[2])
}
