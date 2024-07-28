package base

import "github.com/plantoncloud-inc/go-commons/kubernetes/helm"

type HelmVal struct {
}

func getHelmVal() *HelmVal {
	cv := &HelmVal{}
	return cv
}

// https://artifacthub.io/packages/helm/istio-official/base
func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "istio-base",
		Repo:        "https://istio-release.storage.googleapis.com/charts",
		Name:        "base",
		Version:     "1.15.0-beta.1",
	}
}
