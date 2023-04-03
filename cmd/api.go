package main

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	helmCLI "helm.sh/helm/v3/pkg/cli"
)

func yamlGenerator(w http.ResponseWriter, r *http.Request) {
	var values map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		logrus.Errorf("cannot decode values: %s", err)
		http.Error(w, http.StatusText(400), 400)
		return
	}

	chartName := "onechart/onechart"
	client, settings := helmClient()
	chart, err := loadChart(chartName, client, settings)
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

func helmClient() (*action.Install, *helmCLI.EnvSettings) {
	actionConfig := new(action.Configuration)
	client := action.NewInstall(actionConfig)

	client.DryRun = true
	client.ReleaseName = "release-name"
	client.Namespace = "default"
	client.Replace = true
	client.ClientOnly = true
	client.APIVersions = []string{}
	client.IncludeCRDs = false

	var settings = helmCLI.New()
	return client, settings
}

func loadChart(chartName string, client *action.Install, settings *helmCLI.EnvSettings) (*chart.Chart, error) {
	name, chart, err := client.NameAndChart([]string{chartName})
	if err != nil {
		return nil, err
	}
	client.ReleaseName = name

	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return nil, err
	}

	return loader.Load(cp)
}
