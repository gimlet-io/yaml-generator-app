package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
)

const DEFAULT_CHARTS = "name=onechart,repo=https://chart.onechart.dev,version=0.53.0;name=static-site,repo=https://chart.onechart.dev,version=0.53.0;name=cron-job,repo=https://chart.onechart.dev,version=0.53.0"

// Environ returns the settings from the environment.
func Environ() (*Config, error) {
	cfg := Config{}
	err := envconfig.Process("", &cfg)
	defaults(&cfg)

	return &cfg, err
}

func defaults(c *Config) {
	if c.Charts == nil {
		c.Charts.Decode(DEFAULT_CHARTS)
	}
}

// String returns the configuration in string format.
func (c *Config) String() string {
	out, _ := yaml.Marshal(c)
	return string(out)
}

type Config struct {
	Charts Charts
}

type Charts []Chart

type Chart struct {
	Name    string
	Repo    string
	Version string
}

func (c *Charts) Decode(value string) error {
	charts := []Chart{}
	splittedCharts := strings.Split(value, ";")

	for _, chartsString := range splittedCharts {
		parsedChart, err := parseChartString(chartsString)
		if err != nil {
			return fmt.Errorf("invalid chart format: %s", err)
		}

		if parsedChart != nil {
			charts = append(charts, *parsedChart)
		}
	}
	*c = charts
	return nil
}

func parseChartString(chartsString string) (*Chart, error) {
	if chartsString == "" {
		return nil, nil
	}

	parsedValues, err := parse(chartsString)
	if err != nil {
		return nil, err
	}

	chart := &Chart{
		Name:    parsedValues.Get("name"),
		Repo:    parsedValues.Get("repo"),
		Version: parsedValues.Get("version"),
	}

	return chart, nil
}

func parse(query string) (url.Values, error) {
	values := make(url.Values)
	err := populateValues(values, query)
	return values, err
}

func populateValues(values url.Values, query string) error {
	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, ",")
		if strings.Contains(key, ";") {
			return fmt.Errorf("invalid semicolon separator in query")
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")
		key, err := url.QueryUnescape(key)
		if err != nil {
			return err
		}
		value, err = url.QueryUnescape(value)
		if err != nil {
			return err
		}
		values[key] = append(values[key], value)
	}
	return nil
}
