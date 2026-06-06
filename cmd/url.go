package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:   "url <query>",
	Short: "Print the URL for a matching item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		v, _, err := resolveField(args[0], "url")
		if err != nil {
			return err
		}
		fmt.Println(v)
		return nil
	},
}

func init() { rootCmd.AddCommand(urlCmd) }
