package cmd

import (
	"fmt"

	"github.com/comfucios/vw/internal/clip"
	"github.com/spf13/cobra"
)

var copyField string

var copyCmd = &cobra.Command{
	Use:   "copy <query>",
	Short: "Copy a password, username, URL, or TOTP code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		value, label, err := resolveField(args[0], copyField)
		if err != nil {
			return err
		}
		if err := clip.Copy(value); err != nil {
			return err
		}
		fmt.Println(label, "copied")
		return nil
	},
}

func init() {
	copyCmd.Flags().StringVar(&copyField, "field", "password", "field to copy: password, user, url, code")
	rootCmd.AddCommand(copyCmd)
}

func resolveField(query, field string) (string, string, error) {
	c, err := bwClient(true)
	if err != nil {
		return "", "", err
	}
	item, err := c.FindOne(query)
	if err != nil {
		return "", "", err
	}
	switch field {
	case "password", "pass":
		v, err := c.GetPassword(item.ID)
		return v, "password", err
	case "user", "username", "login":
		if item.Login == nil {
			return "", "", fmt.Errorf("item %q has no login", item.Name)
		}
		return item.Login.Username, "username", nil
	case "url", "uri":
		if item.Login == nil || len(item.Login.URIs) == 0 {
			return "", "", fmt.Errorf("item %q has no URL", item.Name)
		}
		return item.Login.URIs[0].URI, "url", nil
	case "code", "totp":
		v, err := c.GetTOTP(item.ID)
		return v, "totp", err
	default:
		return "", "", fmt.Errorf("unknown field %q", field)
	}
}
