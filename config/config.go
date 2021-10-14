package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type TextFormat string

const (
	Backlog  TextFormat = "backlog"
	Markdown TextFormat = "markdown"
)

// Config is for tfnotify config structure
type Config struct {
	CI         string     `yaml:"ci"`
	Terraform  Terraform  `yaml:"terraform"`
	APIKey     string     `yaml:"api_key"`
	BaseURL    string     `yaml:"base_url"`
	Project    string     `yaml:"project"`
	Repository string     `yaml:"repo"`
	PRNumber   int        `yaml:"pr"`
	Format     TextFormat `yaml:"format"`
}

// Terraform represents terraform configurations
type Terraform struct {
	Plan  Plan  `yaml:"plan"`
	Apply Apply `yaml:"apply"`
}

// Plan is a terraform plan config
type Plan struct {
	Title    string `yaml:"title"`
	Message  string `yaml:"message"`
	Template string `yaml:"template"`
}

// Apply is a terraform apply config
type Apply struct {
	Title    string `yaml:"title"`
	Message  string `yaml:"message"`
	Template string `yaml:"template"`
}

// LoadFile binds the config file to Config structure
func (cfg *Config) LoadFile(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("%s: no config file", path)
	}
	raw, _ := ioutil.ReadFile(path)
	return yaml.Unmarshal(raw, cfg)
}

// Validation validates config file
func (cfg *Config) Validation() error {
	switch strings.ToLower(cfg.CI) {
	case "":
		// ok pattern eg. exec in local
	case "circleci", "circle-ci":
		// ok pattern
	case "gitlabci", "gitlab-ci":
		// ok pattern
	case "travis", "travisci", "travis-ci":
		// ok pattern
	case "codebuild":
		// ok pattern
	case "teamcity":
		// ok pattern
	case "drone":
		// ok pattern
	case "jenkins":
		// ok pattern
	case "github-actions":
		// ok pattern
	case "cloud-build", "cloudbuild":
		// ok pattern
	default:
		return fmt.Errorf("%s: not supported yet", cfg.CI)
	}

	if cfg.BaseURL == "" {
		return errors.New("space is missing")
	}

	if cfg.Project == "" {
		return errors.New("project is missing")
	}

	if cfg.Repository == "" {
		return errors.New("repository is missing")
	}

	if cfg.PRNumber <= 0 {
		return errors.New("pull request number is needed")
	}

	return nil
}

// Find returns config path
func (cfg *Config) Find(file string) (string, error) {
	var files []string
	if file == "" {
		files = []string{
			"tf2b.yaml",
			"tf2b.yml",
			".tf2b.yaml",
			".tf2b.yml",
		}
	} else {
		files = []string{file}
	}
	for _, file := range files {
		_, err := os.Stat(file)
		if err == nil {
			return file, nil
		}
	}
	return "", nil
}
