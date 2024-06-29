package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var green = lipgloss.Color("#03BF87")

func runForm(config *Config) (*Config, error) {
	var (
		padding = strings.Trim(fmt.Sprintf("%v", config.Padding), "[]")
	)

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

	config.Padding = parsePadding(padding)

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

func parsePadding(v string) []float64 {
	var values []float64
	for _, p := range strings.Fields(v) {
		pi, _ := strconv.ParseFloat(p, 64) // already validated
		values = append(values, pi)
	}
	return expandPadding(values, 1)
}

var parseMargin = parsePadding
