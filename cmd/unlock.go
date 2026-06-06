package cmd

import (
	"fmt"

	"github.com/comfucios/vw/internal/session"
	"github.com/spf13/cobra"
)

var unlockCmd = &cobra.Command{
	Use:   "unlock",
	Short: "Unlock the vault and store the session for vw",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		c, err := bwClient(false)
		if err != nil {
			return err
		}
		token, err := c.UnlockRaw()
		if err != nil {
			return err
		}
		if err := session.Save(token, cfg.UseKeyring); err != nil {
			return err
		}
		fmt.Println("vault unlocked")
		return nil
	},
}

func init() { rootCmd.AddCommand(unlockCmd) }
