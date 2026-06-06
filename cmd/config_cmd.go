package cmd

import (
	"fmt"

	"github.com/comfucios/vw/internal/config"
	"github.com/comfucios/vw/internal/paths"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Read and update vw configuration",
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Print config file path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(paths.ConfigFile())
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Print configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		fmt.Println("server_url:", cfg.ServerURL)
		fmt.Println("bw_path:", cfg.BWPath)
		fmt.Println("clipboard_clear_seconds:", cfg.ClipboardClearSeconds)
		fmt.Println("use_keyring:", cfg.UseKeyring)
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.SaveValue(args[0], args[1]); err != nil {
			return err
		}
		fmt.Println("set", args[0])
		return nil
	},
}

func init() {
	configCmd.AddCommand(configPathCmd, configGetCmd, configSetCmd)
	rootCmd.AddCommand(configCmd)
}
