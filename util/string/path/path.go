package path

import (
	"strings"
	"unicode"
)

const (
	InvalidPathCharsUNIX    = `/`
	InvalidPathCharsDarwin  = `/:`
	InvalidPathCharsWindows = `<>:"/\|?*`
)

func SanitizePath(path string, isInvalid func(rune) bool) string {
	var (
		sanitized strings.Builder
		prev      rune
	)

	const underscore = '_'

	for _, r := range path {
		var toWrite rune
		if isInvalid(r) {
			toWrite = underscore
		} else {
			toWrite = r
		}

		// replace two or more consecutive underscores with one underscore
		if (toWrite == underscore && prev != underscore) || toWrite != underscore {
			sanitized.WriteRune(toWrite)
		}

		prev = toWrite
	}

	return strings.TrimFunc(sanitized.String(), func(r rune) bool {
		return r == underscore || unicode.IsSpace(r)
	})
}

func SanitizeWhitespace(path string) string {
	return SanitizePath(path, func(r rune) bool {
		return unicode.IsSpace(r)
	})
}

func SanitizeUNIX(path string) string {
	return SanitizePath(path, func(r rune) bool {
		return strings.ContainsRune(InvalidPathCharsUNIX, r)
	})
}

func SanitizeDarwin(path string) string {
	return SanitizePath(path, func(r rune) bool {
		return strings.ContainsRune(InvalidPathCharsDarwin, r)
	})
}

func SanitizeWindows(path string) string {
	return SanitizePath(path, func(r rune) bool {
		return strings.ContainsRune(InvalidPathCharsWindows, r)
	})
}
