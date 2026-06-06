package paths

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigFile(t *testing.T) {
	require.True(t, strings.HasSuffix(ConfigFile(), "vw/config.yaml"))
}
