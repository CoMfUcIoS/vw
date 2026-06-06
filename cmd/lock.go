package cmd

import (
	"fmt"

	"github.com/comfucios/vw/internal/session"
	"github.com/spf13/cobra"
)

var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock the vault and clear the stored session",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		if c, err := bwClient(true); err == nil {
			_ = c.Lock()
		}
		if err := session.Delete(cfg.UseKeyring); err != nil {
			return err
		}
		fmt.Println("vault locked")
		return nil
	},
}

func init() { rootCmd.AddCommand(lockCmd) }
