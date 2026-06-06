package cmd

import (
	"fmt"

	"github.com/comfucios/vw/internal/ui"
	"github.com/spf13/cobra"
)

var (
	newUser        string
	newURL         string
	newPassword    string
	newInteractive bool
)

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new login item",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := ""
		if len(args) > 0 {
			name = args[0]
		}
		username := newUser
		url := newURL
		password := newPassword
		if newInteractive || name == "" || password == "" {
			answers, err := ui.NewItemForm(name)
			if err != nil {
				return err
			}
			name = answers.Name
			username = answers.Username
			password = answers.Password
			url = answers.URL
		}
		if name == "" {
			return fmt.Errorf("name is required")
		}
		if password == "" {
			return fmt.Errorf("password is required")
		}
		c, err := bwClient(true)
		if err != nil {
			return err
		}
		if err := c.CreateLogin(name, username, password, url); err != nil {
			return err
		}
		fmt.Println("created", name)
		return nil
	},
}

func init() {
	newCmd.Flags().StringVar(&newUser, "user", "", "username")
	newCmd.Flags().StringVar(&newURL, "url", "", "URL")
	newCmd.Flags().StringVar(&newPassword, "password", "", "password")
	newCmd.Flags().BoolVarP(&newInteractive, "interactive", "i", false, "force interactive form")
	rootCmd.AddCommand(newCmd)
}
