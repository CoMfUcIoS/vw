package cmd

import (
	"fmt"

	"github.com/comfucios/vw/internal/bw"
	"github.com/comfucios/vw/internal/session"
)

func bwClient(requireSession bool) (*bw.Client, error) {
	cfg, err := loadConfig()
	if err != nil {
		return nil, err
	}
	sess := ""
	if requireSession {
		sess, err = session.Load(cfg.UseKeyring)
		if err != nil {
			return nil, fmt.Errorf("vault is locked; run 'vw unlock': %w", err)
		}
	}
	return bw.New(cfg.BWPath, sess)
}
