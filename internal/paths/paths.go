package paths

import (
	"os"
	"path/filepath"
)

const AppName = "vw"

func ConfigDir() string {
	if v := os.Getenv("XDG_CONFIG_HOME"); v != "" {
		return filepath.Join(v, AppName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", AppName)
}

func StateDir() string {
	if v := os.Getenv("XDG_STATE_HOME"); v != "" {
		return filepath.Join(v, AppName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "state", AppName)
}

func DataDir() string {
	if v := os.Getenv("XDG_DATA_HOME"); v != "" {
		return filepath.Join(v, AppName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", AppName)
}

func CacheDir() string {
	if v := os.Getenv("XDG_CACHE_HOME"); v != "" {
		return filepath.Join(v, AppName)
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".cache", AppName)
}

func ConfigFile() string    { return filepath.Join(ConfigDir(), "config.yaml") }
func ManagedBWPath() string { return filepath.Join(DataDir(), "bin", "bw") }
