package cmd

import (
	"fmt"

	"github.com/comfucios/vw/internal/bootstrap"
	"github.com/spf13/cobra"
)

var (
	bootstrapVersion string
	bootstrapForce   bool
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap-bw",
	Short: "Download and install a managed Bitwarden CLI binary",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := bootstrap.InstallBW(bootstrap.Options{Version: bootstrapVersion, Force: bootstrapForce})
		if err != nil {
			return err
		}
		fmt.Println("bw installed at", path)
		return nil
	},
}

func init() {
	bootstrapCmd.Flags().StringVar(&bootstrapVersion, "version", "", "specific bw version, for example 2026.5.0")
	bootstrapCmd.Flags().BoolVar(&bootstrapForce, "force", false, "redownload even if managed bw already exists")
	rootCmd.AddCommand(bootstrapCmd)
}
