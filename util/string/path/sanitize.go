//go:build !windows && !darwin

package path

const InvalidPathCharsOS = InvalidPathCharsUNIX

func SanitizeOS(path string) string {
	return SanitizeUNIX(path)
}
