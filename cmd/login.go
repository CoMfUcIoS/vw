package cmd

import "github.com/spf13/cobra"

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in with bw",
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := bwClient(false)
		if err != nil {
			return err
		}
		return c.Login()
	},
}

func init() { rootCmd.AddCommand(loginCmd) }
