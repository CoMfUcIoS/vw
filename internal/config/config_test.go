package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewViperDefaults(t *testing.T) {
	v := NewViper()
	require.Equal(t, 45, v.GetInt("clipboard_clear_seconds"))
	require.True(t, v.GetBool("use_keyring"))
}
