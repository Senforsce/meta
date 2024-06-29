package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/go-git/go-git/v5"
)

const defaultOutputFilename = "freeze.png"

type Project struct {
	RepoURL       string
	Directory     string
	DefaultBranch string
	State         string
	Name          string
}

// Config is the configuration options for a screenshot.
type Config struct {
	Input    string
	Projects []Project

	// Window
	Margin  []float64
	Padding []float64
	Window  bool
	Width   float64
	Height  float64

	// Settings
	Config      string
	Interactive bool

	Output         string
	Execute        string
	ExecuteTimeout time.Duration
}

// Shadow is the configuration options for a drop shadow.
type Shadow struct {
	Blur float64
	X    float64
	Y    float64
}

// Border is the configuration options for a window border.
type Border struct {
	Radius float64
	Width  float64
	Color  string
}

// Font is the configuration options for a font.
type Font struct {
	Family    string
	ily       string
	File      string
	Size      float64
	Ligatures bool
}

//go:embed configurations/*
var configs embed.FS

func expandPadding(p []float64, scale float64) []float64 {
	switch len(p) {
	case 1:
		return []float64{p[top] * scale, p[top] * scale, p[top] * scale, p[top] * scale}
	case 2:
		return []float64{p[top] * scale, p[right] * scale, p[top] * scale, p[right] * scale}
	case 4:
		return []float64{p[top] * scale, p[right] * scale, p[bottom] * scale, p[left] * scale}
	default:
		return []float64{0, 0, 0, 0}

	}
}

var expandMargin = expandPadding

type side int

const (
	top    side = 0
	right  side = 1
	bottom side = 2
	left   side = 3
)

var userConfigPath = filepath.Join(xdg.ConfigHome, "tndr", "meta.json")

func loadUserConfig() (*Config, error) {
	f, readFileErr := loadUserConfigFile()

	c := &Config{}

	if readFileErr != nil {
		return nil, readFileErr
	}
	defer f.Close()

	byteValue, _ := io.ReadAll(f)

	marshallErr := json.Unmarshal(byteValue, c)

	if marshallErr != nil {
		return nil, marshallErr
	}

	return c, nil
}

func loadUserConfigFile() (fs.File, error) {
	return os.Open(userConfigPath)
}

func addProject(directory string, url string, name string) error {

	c, configErr := loadUserConfig()

	if configErr != nil {
		return configErr
	}

	c.Projects = append(c.Projects, Project{RepoURL: url, Directory: directory, Name: name, DefaultBranch: "main", State: "added"})

	saveUserConfig(*c)

	return nil
}

func cloneProjects() error {

	c, configErr := loadUserConfig()

	if configErr != nil {
		return configErr
	}

	e, err := os.Getwd()
	if err != nil {
		return err
	}

	for _, p := range c.Projects {
		fmt.Printf("Will clone : %s \n", path.Dir(e+"/"+p.Directory))
		_, err := git.PlainClone(path.Dir(e+"/"+p.Directory), false, &git.CloneOptions{
			URL:      p.RepoURL,
			Progress: os.Stdout,
		})

		if err != nil {
			return err
		}

	}

	saveUserConfig(*c)

	return nil
}

func saveUserConfig(config Config) error {
	config.Input = ""
	config.Output = ""
	config.Interactive = false

	err := os.MkdirAll(filepath.Dir(userConfigPath), os.ModePerm)
	if err != nil {
		return err
	}
	f, err := os.Create(userConfigPath)
	if err != nil {
		return err
	}
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = f.Write(b)

	print(userConfigPath)

	return err
}
