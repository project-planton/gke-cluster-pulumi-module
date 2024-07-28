package reflector

import "github.com/plantoncloud-inc/go-commons/kubernetes/helm"

const (
	Namespace = "reflector"
)

type HelmVal struct{}

func getHelmVal() *HelmVal {
	return &HelmVal{}
}

func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "reflector",
		Repo:        "https://emberstack.github.io/helm-charts",
		Name:        "reflector",
		Version:     "6.1.47",
	}
}
