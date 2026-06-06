package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var codeCmd = &cobra.Command{
	Use:   "code <query>",
	Short: "Print a TOTP code for a matching item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		v, _, err := resolveField(args[0], "code")
		if err != nil {
			return err
		}
		fmt.Println(v)
		return nil
	},
}

func init() { rootCmd.AddCommand(codeCmd) }
