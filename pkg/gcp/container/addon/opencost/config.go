package opencost

import (
	"github.com/plantoncloud-inc/go-commons/kubernetes/helm"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/prometheus"
)

const (
	Namespace = "opencost"
)

type HelmVal struct {
	OpenCost *ValOpenCost `yaml:"opencost"`
}

type ValOpenCost struct {
	Prometheus *ValPrometheus `yaml:"prometheus"`
}

type ValPrometheus struct {
	Internal *ValPrometheusInternal `yaml:"internal"`
}

type ValPrometheusInternal struct {
	Enabled       bool   `yaml:"enabled"`
	ServiceName   string `yaml:"serviceName"`
	NamespaceName string `yaml:"namespaceName"`
	Port          int32  `yaml:"port"`
}

func getHelmVal() *HelmVal {
	return &HelmVal{
		OpenCost: &ValOpenCost{
			//https://github.com/opencost/opencost-helm-chart/blob/main/charts/opencost/values.yaml#L137
			Prometheus: &ValPrometheus{
				Internal: &ValPrometheusInternal{
					Enabled:       true,
					ServiceName:   prometheus.KubeServiceName,
					NamespaceName: prometheus.Namespace,
					Port:          80,
				},
			},
		},
	}
}

// https://github.com/opencost/opencost-helm-chart
func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "opencost",
		Repo:        "https://opencost.github.io/opencost-helm-chart",
		Name:        "opencost",
		Version:     "1.11.0",
	}
}
