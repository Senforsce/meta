package main

import (
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var green = lipgloss.Color("#03BF87")

func runForm(config *Config) (*Config, error) {

	theme := huh.ThemeCharm()
	theme.FieldSeparator = lipgloss.NewStyle()
	theme.Blurred.BlurredButton = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).PaddingRight(1)
	theme.Blurred.FocusedButton = lipgloss.NewStyle().Foreground(lipgloss.Color("7")).PaddingRight(1)
	theme.Focused.BlurredButton = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).PaddingRight(1)
	theme.Focused.FocusedButton = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).PaddingRight(1)
	theme.Focused.Base.BorderForeground(green)

	f := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().Title("\nCapture file"),

			huh.NewFilePicker().
				Title("").
				Picking(true).
				Height(10).
				Value(&config.Input),

			huh.NewNote().Description("Choose a code file to screenshot."),
		).WithHide(config.Input != "" || config.Execute != ""),
	).WithTheme(theme).WithWidth(40)

	err := f.Run()

	if config.Output == "" {
		config.Output = ".meta.ttl"
	}

	return config, err
}

// var colorRegex = regexp.MustCompile("^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$")

// func validateColor(s string) error {
// 	if len(s) <= 0 {
// 		return nil
// 	}

// 	if !colorRegex.MatchString(s) {
// 		return errors.New("must be valid color")
// 	}
// 	return nil
// }
