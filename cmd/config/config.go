package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

// Environ returns the settings from the environment.
func Environ() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	defaults(&cfg)

	return &cfg, err
}

func defaults(c *Config) {
	if c.Chart.Name == "" {
		c.Chart.Name = "onechart"
	}
	if c.Chart.Repo == "" {
		c.Chart.Repo = "https://chart.onechart.dev"
	}
	if c.Chart.Version == "" {
		c.Chart.Version = "0.46.0"
	}
}

// String returns the configuration in string format.
func (c *Config) String() string {
	out, _ := yaml.Marshal(c)
	return string(out)
}

type Config struct {
	Chart Chart
}

type Chart struct {
	Name    string `envconfig:"CHART_NAME"`
	Repo    string `envconfig:"CHART_REPO"`
	Version string `envconfig:"CHART_VERSION"`
}

type Multiline string

func (m *Multiline) Decode(value string) error {
	value = strings.ReplaceAll(value, "\\n", "\n")
	*m = Multiline(value)
	return nil
}

func (m *Multiline) String() string {
	return string(*m)
}
