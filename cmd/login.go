package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/comfucios/vw/internal/bw"
	"github.com/spf13/cobra"
)

var (
	loginEmail    string
	loginPassword string
	loginMethod   int
	loginCode     string
	loginOTP      bool
	loginAPIKey   bool
	loginSSO      bool
	loginRaw      bool
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Bitwarden or Vaultwarden",
	Long: `Log in to Bitwarden or Vaultwarden using the managed Bitwarden CLI.

By default, vw delegates to 'bw login' interactively. This lets the
Bitwarden CLI handle email, master password, and supported two-step login
prompts directly.

Examples:
  vw login
  vw login --raw
  vw login --otp
  vw login --method 0 --code 123456
  vw login --email me@example.com --password 'master-password'
  vw login --email me@example.com --password 'master-password' --method 0 --code 123456
  vw login --apikey
  vw login --sso

Two-step method values are Bitwarden CLI enum values. Common values are:
  0  authenticator app
  1  email
  3  YubiKey

If your account uses FIDO2, Duo, or an unsupported two-step method, use:
  vw login --apikey`,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := bwClient(false)
		if err != nil {
			return err
		}

		if loginAPIKey && loginSSO {
			return fmt.Errorf("--apikey and --sso cannot be used together")
		}

		if loginAPIKey {
			return c.LoginAPIKey()
		}

		if loginSSO {
			return c.LoginSSO()
		}

		// Raw mode delegates completely to bw. This is the safest fallback
		// because bw owns any current/future login prompts.
		if loginRaw {
			return c.LoginInteractive()
		}

		email := strings.TrimSpace(loginEmail)
		password := loginPassword
		code := strings.TrimSpace(loginCode)

		// No flags: delegate directly to bw login. This lets bw prompt for
		// email, master password, and OTP using its native flow.
		if email == "" && password == "" && code == "" && !loginOTP && !cmd.Flags().Changed("method") {
			return c.LoginInteractive()
		}

		if email == "" || password == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Email address").
						Value(&email),
					huh.NewInput().
						Title("Master password").
						EchoMode(huh.EchoModePassword).
						Value(&password),
				),
			)

			if err := form.Run(); err != nil {
				return err
			}
		}

		if loginOTP && code == "" {
			form := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Two-step code").
						Description("Enter your authenticator/email/YubiKey two-step login code.").
						Value(&code),
				),
			)

			if err := form.Run(); err != nil {
				return err
			}
		}

		email = strings.TrimSpace(email)
		code = strings.TrimSpace(code)

		if email == "" {
			return fmt.Errorf("email address is required")
		}
		if password == "" {
			return fmt.Errorf("master password is required")
		}

		opts := bw.LoginOptions{
			Email:    email,
			Password: password,
		}

		if cmd.Flags().Changed("method") {
			opts.Method = &loginMethod
		}

		// Convenience: vw login --otp prompts for a code and defaults to
		// authenticator-app method 0 unless --method was explicitly provided.
		if loginOTP && opts.Method == nil {
			method := 0
			opts.Method = &method
		}

		if code != "" {
			opts.Code = code
		}

		return c.LoginWithPassword(opts)
	},
}

func init() {
	loginCmd.Flags().StringVar(&loginEmail, "email", "", "email address")
	loginCmd.Flags().StringVar(&loginPassword, "password", "", "master password")
	loginCmd.Flags().IntVar(&loginMethod, "method", 0, "Bitwarden two-step login provider method")
	loginCmd.Flags().StringVar(&loginCode, "code", "", "two-step login code")
	loginCmd.Flags().BoolVar(&loginOTP, "otp", false, "prompt for a two-step login code")
	loginCmd.Flags().BoolVar(&loginAPIKey, "apikey", false, "log in using Bitwarden personal API key")
	loginCmd.Flags().BoolVar(&loginSSO, "sso", false, "log in using SSO")
	loginCmd.Flags().BoolVar(&loginRaw, "raw", false, "delegate directly to 'bw login' interactive flow")

	rootCmd.AddCommand(loginCmd)
}
