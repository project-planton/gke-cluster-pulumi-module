package kubecost

import (
	"github.com/plantoncloud-inc/go-commons/kubernetes/helm"
)

const (
	Namespace = "kubecost"
)

type HelmVal struct {
}

func getHelmVal() *HelmVal {
	return &HelmVal{}
}

// https://github.com/kubecost/cost-analyzer-helm-chart
// https://docs.kubecost.com/install-and-configure/install#alternative-installation-methods
func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "kubecost",
		Repo:        "https://kubecost.github.io/cost-analyzer/",
		Name:        "cost-analyzer",
		Version:     "1.107.0",
	}
}
