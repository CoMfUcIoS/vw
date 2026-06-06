package cmd

import (
	"fmt"
	"os"

	"github.com/comfucios/vw/internal/bw"
	"github.com/comfucios/vw/internal/paths"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check vw configuration and dependencies",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		fmt.Println("config:", paths.ConfigFile())
		fmt.Println("managed bw:", paths.ManagedBWPath())
		p, err := bw.Resolve(cfg.BWPath)
		if err != nil {
			fmt.Println("bw: not found")
			fmt.Fprintln(os.Stderr, err)
			return nil
		}
		fmt.Println("bw:", p)
		c, err := bw.New(cfg.BWPath, "")
		if err != nil {
			return err
		}
		out, err := c.Run("--version")
		if err != nil {
			return err
		}
		fmt.Println("bw version:", out)
		return nil
	},
}

func init() { rootCmd.AddCommand(doctorCmd) }
