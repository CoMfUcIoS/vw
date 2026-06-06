package config

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/comfucios/vw/internal/paths"
	"github.com/spf13/viper"
)

type Config struct {
	ServerURL             string
	BWPath                string
	ClipboardClearSeconds int
	UseKeyring            bool
}

func NewViper() *viper.Viper {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(paths.ConfigDir())
	v.SetEnvPrefix("VW")
	v.AutomaticEnv()

	v.SetDefault("server_url", "")
	v.SetDefault("bw_path", "")
	v.SetDefault("clipboard_clear_seconds", 45)
	v.SetDefault("use_keyring", true)
	return v
}

func Load() (*Config, *viper.Viper, error) {
	v := NewViper()
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, nil, err
		}
	}
	return &Config{
		ServerURL:             v.GetString("server_url"),
		BWPath:                v.GetString("bw_path"),
		ClipboardClearSeconds: v.GetInt("clipboard_clear_seconds"),
		UseKeyring:            v.GetBool("use_keyring"),
	}, v, nil
}

func SaveValue(key, value string) error {
	v := NewViper()
	_ = v.ReadInConfig()
	v.Set(key, value)
	if err := os.MkdirAll(paths.ConfigDir(), 0o700); err != nil {
		return err
	}
	if v.ConfigFileUsed() == "" {
		return v.WriteConfigAs(paths.ConfigFile())
	}
	return v.WriteConfig()
}

func EnsureConfigDir() error {
	return os.MkdirAll(filepath.Dir(paths.ConfigFile()), 0o700)
}
