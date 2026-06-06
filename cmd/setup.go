package cmd

import (
	"fmt"

	"github.com/comfucios/vw/internal/bootstrap"
	"github.com/comfucios/vw/internal/bw"
	"github.com/comfucios/vw/internal/config"
	"github.com/comfucios/vw/internal/ui"
	"github.com/spf13/cobra"
)

var (
	setupServer      string
	setupInteractive bool
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure vw for Bitwarden or Vaultwarden",
	RunE: func(cmd *cobra.Command, args []string) error {
		server := setupServer
		bootstrapMissing := true
		if setupInteractive || server == "" {
			answers, err := ui.SetupForm(server)
			if err != nil {
				return err
			}
			server = answers.ServerURL
			bootstrapMissing = answers.Bootstrap
		}
		if err := config.SaveValue("server_url", server); err != nil {
			return err
		}
		if bootstrapMissing {
			if _, err := bw.Resolve(""); err != nil {
				if _, err := bootstrap.InstallBW(bootstrap.Options{}); err != nil {
					return err
				}
			}
		}
		if server != "" {
			c, err := bw.New("", "")
			if err != nil {
				return err
			}
			if err := c.ConfigServer(server); err != nil {
				return err
			}
			fmt.Println("server configured:", server)
		} else {
			fmt.Println("using Bitwarden cloud defaults")
		}
		return nil
	},
}

func init() {
	setupCmd.Flags().StringVar(&setupServer, "server", "", "Vaultwarden or self-hosted Bitwarden server URL")
	setupCmd.Flags().BoolVarP(&setupInteractive, "interactive", "i", false, "force interactive setup")
	rootCmd.AddCommand(setupCmd)
}
