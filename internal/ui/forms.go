package ui

import "github.com/charmbracelet/huh"

type SetupAnswers struct {
	ServerURL string
	Bootstrap bool
}

func SetupForm(defaultServer string) (*SetupAnswers, error) {
	a := &SetupAnswers{ServerURL: defaultServer, Bootstrap: true}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Bitwarden/Vaultwarden server URL").Description("Leave empty for Bitwarden cloud.").Value(&a.ServerURL),
			huh.NewConfirm().Title("Download managed bw if missing?").Value(&a.Bootstrap),
		),
	)
	return a, form.Run()
}

type NewItemAnswers struct {
	Name     string
	Username string
	Password string
	URL      string
}

func NewItemForm(defaultName string) (*NewItemAnswers, error) {
	a := &NewItemAnswers{Name: defaultName}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Value(&a.Name),
			huh.NewInput().Title("Username").Value(&a.Username),
			huh.NewInput().Title("Password").EchoMode(huh.EchoModePassword).Value(&a.Password),
			huh.NewInput().Title("URL").Value(&a.URL),
		),
	)
	return a, form.Run()
}
