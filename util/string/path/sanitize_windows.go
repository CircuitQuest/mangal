package path

const InvalidPathCharsOS = InvalidPathCharsWindows

func SanitizeOS(path string) string {
	return SanitizeWindows(path)
}
