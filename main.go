package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/charmbracelet/lipgloss"
)

var (
	// Version contains the application version number. It's set via ldflags
	// when building.
	Version = ""

	// CommitSHA contains the SHA of the commit that this application was built
	// against. It's set via ldflags when building.
	CommitSHA = ""
)

func init() {
	driver, ok := os.LookupEnv("O8ROOT")

	if !ok {
		fmt.Println("O8ROOT is not present, please make sure O8ROOT is exported before running the command")
	} else {
		fmt.Printf("OntologyRootFolder: %s\n", driver)
	}

	ns, ok := os.LookupEnv("O8_META_NAMESPACE")

	if !ok {
		fmt.Println("O8NAMESPACE is not present, please make sure O8NAMESPACE is exported before running the command")
	} else {
		fmt.Printf("OntologyRootNamespace: %s\n", ns)
	}

}

func main() {
	const shaLen = 8

	var (
		err error
	)

	switch os.Args[1] {

	case "version":
		if Version == "" {
			if info, ok := debug.ReadBuildInfo(); ok && info.Main.Sum != "" {
				Version = info.Main.Version
			} else {
				Version = "unknown (built from source)"
			}
		}
		version := fmt.Sprintf("senforsce Meta version %s", Version)
		if len(CommitSHA) >= shaLen {
			version += " (" + CommitSHA[:shaLen] + ")"
		}

		fmt.Println(version)
		os.Exit(0)

	case "add":
		if os.Args[2] != "" || os.Args[3] != "" {
			overrideName := ""
			if len(os.Args) > 4 {
				overrideName = os.Args[4]
			}
			addProject(os.Args[2], os.Args[3], overrideName)
		}

		os.Exit(0)

	case "clone":
		cloneProjects()
		os.Exit(0)

	case "update":
		print("update", "NOT IMPLEMENTED")

		os.Exit(0)
	}

	//istty := isatty.IsTerminal(os.Stdout.Fd())

	// reading from file.
	//if istty {
	// config.Output = strings.TrimSuffix(filepath.Base(config.Input), filepath.Ext(config.Input)) + ".svg"
	// err = doc.WriteToFile(config.Output)
	// printFilenameOutput(config.Output)
	//} else {
	// _, err = doc.WriteTo(os.Stdout)
	//}
	if err != nil {
		printErrorFatal("Unable to write output", err)
	}

}

func print(line string, headerText string) {
	head := lipgloss.NewStyle().Foreground(lipgloss.Color("#F1F1F1")).Background(lipgloss.Color("#6C50FF")).Bold(true).Padding(0, 1).MarginRight(1).SetString(headerText)

	fmt.Println(lipgloss.JoinHorizontal(lipgloss.Center, head.String(), line))
}
