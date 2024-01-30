package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// TODO: change this to operate on the Fields Entries instead of directly on viper?
// not sure if needed as these operations are run once and then the program exits,
// so they don't need to keep living in the Config

func Set(key string, value any) error {
	viper.Set(key, value)
	return nil
}

func Get(key string) any {
	return viper.Get(key)
}

func Exists(key string) bool {
	return viper.IsSet(key)
}

func Keys() []string {
	return viper.AllKeys()
}

func Write() error {
	switch viper.WriteConfig().(type) {
	case nil:
		return nil
	case viper.ConfigFileNotFoundError:
		return viper.SafeWriteConfig()
	default:
		return fmt.Errorf("unexpected error")
	}
}
