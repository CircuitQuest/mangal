package log

import (
	"fmt"
	"strings"
	"time"
)

// Aggregate is the buildup of all the libmangal.Logger logs.
var Aggregate = strings.Builder{}

// Log is a convenience function to add log messages to Aggregate in a custom format.
func Log(s string) {
	Aggregate.WriteString(fmt.Sprintf("[%s] %s\n", time.Now().Format(time.TimeOnly), s))
}
