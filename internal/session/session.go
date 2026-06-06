package session

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/99designs/keyring"
	"github.com/comfucios/vw/internal/paths"
)

const (
	service = "vw"
	user    = "bw_session"
)

func Save(token string, useKeyring bool) error {
	if strings.TrimSpace(token) == "" {
		return errors.New("empty session token")
	}
	if useKeyring {
		kr, err := openKeyring()
		if err == nil {
			return kr.Set(keyring.Item{Key: user, Data: []byte(token)})
		}
	}
	return saveFile(token)
}

func Load(useKeyring bool) (string, error) {
	if useKeyring {
		kr, err := openKeyring()
		if err == nil {
			item, err := kr.Get(user)
			if err == nil {
				return string(item.Data), nil
			}
		}
	}
	return loadFile()
}

func Delete(useKeyring bool) error {
	if useKeyring {
		kr, err := openKeyring()
		if err == nil {
			_ = kr.Remove(user)
		}
	}
	_ = os.Remove(filePath())
	return nil
}

func openKeyring() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{ServiceName: service})
}

func filePath() string { return filepath.Join(paths.StateDir(), "session") }

func saveFile(token string) error {
	if err := os.MkdirAll(paths.StateDir(), 0o700); err != nil {
		return err
	}
	return os.WriteFile(filePath(), []byte(token), 0o600)
}

func loadFile() (string, error) {
	b, err := os.ReadFile(filePath())
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}
