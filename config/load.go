package config

import (
	"fmt"

	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/util/afs"
	"github.com/spf13/viper"
)

func Load(path string) error {
	if path == "" {
		return fmt.Errorf("empty config directory provided")
	}

	viper.SetConfigName(meta.AppName)
	viper.SetConfigType("toml")
	viper.SetFs(afs.Afero.Fs)
	viper.AddConfigPath(path)
	viper.KeyDelimiter(".")
	viper.SetTypeByDefaultValue(false)

	for _, field := range Fields {
		marshalled, err := field.Marshal(field.Default)
		if err != nil {
			return err
		}
		viper.SetDefault(field.Key, marshalled)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	for _, field := range Fields {
		unmarshalled, err := field.Unmarshal(viper.Get(field.Key))
		if err != nil {
			return err
		}

		if err := field.Validate(unmarshalled); err != nil {
			return fmt.Errorf("%s: %s", field.Key, err)
		}

		if err := field.SetValue(unmarshalled); err != nil {
			return err
		}
	}

	return nil
}
