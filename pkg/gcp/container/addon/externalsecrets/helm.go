package externalsecrets

import (
	helmcommons "github.com/plantoncloud-inc/go-commons/kubernetes/helm"
)

// https://github.com/external-secrets/external-secrets
// https://external-secrets.io/v0.5.6/guides-getting-started/
func getHelmChart() *helmcommons.Chart {
	return &helmcommons.Chart{
		ReleaseName: "external-secrets",
		Repo:        "https://charts.external-secrets.io",
		Name:        "external-secrets",
		Version:     "0.6.1",
	}
}

// HelmVal https://github.com/external-secrets/external-secrets/blob/main/deploy/charts/external-secrets/values.yaml
type HelmVal struct {
	CustomResourceManagerDisabled bool            `yaml:"customResourceManagerDisabled"`
	Crds                          *Crds           `yaml:"crds"`
	Env                           *Env            `yaml:"env"`
	Rbac                          *Rbac           `yaml:"rbac"`
	ServiceAccount                *ServiceAccount `yaml:"serviceAccount"`
	ReplicaCount                  int             `yaml:"replicaCount"`
}

type Crds struct {
	Create bool `yaml:"create"`
}

type Env struct {
	PollerIntervalMilliseconds int    `yaml:"POLLER_INTERVAL_MILLISECONDS"`
	LogLevel                   string `yaml:"LOG_LEVEL"`
	LogMessageKey              string `yaml:"LOG_MESSAGE_KEY"`
	MetricsPort                int    `yaml:"METRICS_PORT"`
}

type Rbac struct {
	Create bool `yaml:"create"`
}

type ServiceAccount struct {
	Create      bool                   `yaml:"create"`
	Annotations map[string]interface{} `yaml:"annotations"`
	Name        string                 `yaml:"name"`
}

func getHelmVal() *HelmVal {
	return &HelmVal{
		CustomResourceManagerDisabled: false,
		Crds: &Crds{
			Create: true,
		},
		Env: &Env{
			PollerIntervalMilliseconds: SecretsPollingIntervalSeconds * 1000,
			LogLevel:                   "info",
			LogMessageKey:              "msg",
			MetricsPort:                3001,
		},
		Rbac: &Rbac{
			Create: true,
		},
		ServiceAccount: &ServiceAccount{
			Create: false,
			Name:   Ksa,
		},
		ReplicaCount: 1,
	}
}
