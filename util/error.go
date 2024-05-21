package util

import (
	"fmt"
	"os"
)

// TODO: also use the logger to possibly write to file
//
// Printf wrapper that also exits, uses stderrr.
func Errorf(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
