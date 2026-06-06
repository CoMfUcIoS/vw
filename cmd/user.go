package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user <query>",
	Short: "Print the username for a matching item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		v, _, err := resolveField(args[0], "user")
		if err != nil {
			return err
		}
		fmt.Println(v)
		return nil
	},
}

func init() { rootCmd.AddCommand(userCmd) }
