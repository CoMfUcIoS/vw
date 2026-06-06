package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [query]",
	Short: "List matching vault items",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) > 0 {
			query = args[0]
		}
		c, err := bwClient(true)
		if err != nil {
			return err
		}
		items, err := c.ListItems(query)
		if err != nil {
			return err
		}
		for _, item := range items {
			user := ""
			if item.Login != nil {
				user = item.Login.Username
			}
			fmt.Printf("%s\t%s\t%s\n", item.ID, item.Name, user)
		}
		return nil
	},
}

func init() { rootCmd.AddCommand(listCmd) }
