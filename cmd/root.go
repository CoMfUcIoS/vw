package cmd

import (
	"fmt"

	"github.com/comfucios/vw/internal/config"
	"github.com/comfucios/vw/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "vw",
	Short:   "A friendly CLI wrapper for Bitwarden and Vaultwarden",
	Long:    "vw is an opinionated CLI wrapper around the official Bitwarden CLI (bw).",
	Version: fmt.Sprintf("%s (commit %s, built %s)", version.Version(), version.Commit(), version.Date()),
}

func Execute() error { return rootCmd.Execute() }

func init() {
	rootCmd.PersistentFlags().BoolP("help", "h", false, "show help")
}

func loadConfig() (*config.Config, error) {
	cfg, _, err := config.Load()
	return cfg, err
}
