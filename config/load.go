package config

import (
	"fmt"
	"path"

	"github.com/luevano/mangal/meta"
	"github.com/luevano/mangal/util/afs"
	"github.com/spf13/viper"
)

// Load sets the default values from registered values,
// then reads the config file and sets the keys found in viper,
// validating that it is a valid config value.
func Load(configDir string) error {
	configFile := path.Join(configDir, fmt.Sprintf("%s.toml", meta.AppName))
	// Without Reset it wouldn't be possible to specify a
	// different config file after one has been loaded.
	//
	// Uncomment if different config files might be loaded in the same session.
	// viper.Reset()
	viper.SetConfigFile(configFile)
	viper.KeyDelimiter(".")
	viper.SetFs(afs.Afero.Fs)
	viper.SetTypeByDefaultValue(false)

	// Set default registered values
	for _, field := range Fields {
		marshalled, err := field.Marshal(field.Default)
		if err != nil {
			return err
		}
		viper.SetDefault(field.Key, marshalled)
	}

	// Read config file and set the discovered keys
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// Validate the just set keys (and defaults when the key is not found),
	// and set back the value into the Fields map.
	// SetValue also re-sets the viper key for some reason?
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
