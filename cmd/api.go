package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gimlet-io/yaml-generator-app/cmd/config"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	helmCLI "helm.sh/helm/v3/pkg/cli"
)

func yamlGenerator(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	config := ctx.Value("config").(*config.Config)
	var values map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&values)
	if err != nil {
		logrus.Errorf("cannot decode values: %s", err)
		http.Error(w, http.StatusText(400), 400)
		return
	}
	logrus.Infof("%s", values)

	client, settings := helmClient(&config.Chart)
	chart, err := loadChart(&config.Chart, client, settings)
	if err != nil {
		logrus.Errorf("couldn't load chart: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	rel, err := client.Run(chart, values)
	if err != nil {
		logrus.Errorf("couldn't template values: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(rel.Manifest))
}

func helmClient(chart *config.Chart) (*action.Install, *helmCLI.EnvSettings) {
	actionConfig := new(action.Configuration)
	client := action.NewInstall(actionConfig)

	client.DryRun = true
	client.ReleaseName = "my-release"
	client.Replace = true
	client.ClientOnly = true
	client.APIVersions = []string{}
	client.IncludeCRDs = false
	client.ChartPathOptions.RepoURL = chart.Repo
	client.ChartPathOptions.Version = chart.Version
	client.Namespace = "default"

	var settings = helmCLI.New()
	return client, settings
}

func loadChart(chart *config.Chart, client *action.Install, settings *helmCLI.EnvSettings) (*chart.Chart, error) {
	cp, err := client.ChartPathOptions.LocateChart(chart.Name, settings)
	if err != nil {
		return nil, fmt.Errorf("cannot locate chart %s", err)
	}

	return loader.Load(cp)
}
