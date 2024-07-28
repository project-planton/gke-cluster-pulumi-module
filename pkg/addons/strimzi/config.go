package strimzi

import "github.com/plantoncloud-inc/go-commons/kubernetes/helm"

const (
	Namespace = "strimzi"
)

type HelmVal struct {
	WatchAnyNamespace bool `yaml:"watchAnyNamespace"`
}

func getHelmVal() *HelmVal {
	return &HelmVal{
		WatchAnyNamespace: true,
	}
}

func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "strimzi",
		//https://artifacthub.io/packages/helm/strimzi/strimzi-kafka-operator
		Repo:    "https://strimzi.io/charts/",
		Name:    "strimzi-kafka-operator",
		Version: "0.35.1",
	}
}
