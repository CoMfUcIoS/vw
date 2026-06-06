package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <query>",
	Short: "Print a password for a matching item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := bwClient(true)
		if err != nil {
			return err
		}
		item, err := c.FindOne(args[0])
		if err != nil {
			return err
		}
		pw, err := c.GetPassword(item.ID)
		if err != nil {
			return err
		}
		fmt.Println(pw)
		return nil
	},
}

func init() { rootCmd.AddCommand(getCmd) }
