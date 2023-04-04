package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	helmCLI "helm.sh/helm/v3/pkg/cli"
)

type Chart struct {
	Repository string `yaml:"repository" json:"repository"`
	Name       string `yaml:"name" json:"name"`
	Version    string `yaml:"version" json:"version"`
}

func yamlGenerator(w http.ResponseWriter, r *http.Request) {
	var values map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		logrus.Errorf("cannot decode values: %s", err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	onechart := Chart{
		Repository: "https://chart.onechart.dev",
		Name:       "onechart",
		Version:    "0.41.0",
	}
	client, settings := helmClient(&onechart)
	chart, err := loadChart(&onechart, client, settings)
	if err != nil {
		logrus.Errorf("couldn't load chart", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rel, err := client.Run(chart, values)
	if err != nil {
		logrus.Errorf("couldn't template values", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	manifestString, err := yaml.Marshal(rel.Manifest)
	if err != nil {
		logrus.Errorf("couldn't marshal manifest", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(200)
	w.Write(manifestString)
}

func helmClient(chart *Chart) (*action.Install, *helmCLI.EnvSettings) {
	actionConfig := new(action.Configuration)
	client := action.NewInstall(actionConfig)

	client.DryRun = true
	client.ReleaseName = "release-name"
	client.Replace = true
	client.ClientOnly = true
	client.APIVersions = []string{}
	client.IncludeCRDs = false
	client.ChartPathOptions.RepoURL = chart.Repository
	client.ChartPathOptions.Version = chart.Version
	client.Namespace = "default"

	var settings = helmCLI.New()
	return client, settings
}

func loadChart(chart *Chart, client *action.Install, settings *helmCLI.EnvSettings) (*chart.Chart, error) {
	cp, err := client.ChartPathOptions.LocateChart(chart.Name, settings)
	if err != nil {
		return nil, fmt.Errorf("cannot locate chart %s", err)
	}

	return loader.Load(cp)
}
