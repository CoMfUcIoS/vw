package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync the Bitwarden/Vaultwarden vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := bwClient(true)
		if err != nil {
			return err
		}
		if err := c.Sync(); err != nil {
			return err
		}
		fmt.Println("synced")
		return nil
	},
}

func init() { rootCmd.AddCommand(syncCmd) }
