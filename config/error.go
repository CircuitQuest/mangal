package config

import "fmt"

// TODO: refactor all error msgs to a better error handling (see viper for example)
//
// just to get a prefixed error msg
func errorf(format string, a ...any) error {
	prefix := fmt.Sprintf("config error: %s", format)
	return fmt.Errorf(prefix, a...)
}
