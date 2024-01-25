package path

const InvalidPathCharsOS = InvalidPathCharsDarwin

func SanitizeOS(path string) string {
	return SanitizeDarwin(path)
}
