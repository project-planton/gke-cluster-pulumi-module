package istiod

import "github.com/plantoncloud-inc/go-commons/kubernetes/helm"

type helmValues struct {
	HelmValuesMeshConfig *helmValuesMeshConfig `yaml:"meshConfig"`
}

// helmValuesMeshConfig https://github.com/istio/istio/issues/37329#issuecomment-1042353593
type helmValuesMeshConfig struct {
	IngressClass          string `yaml:"ingressClass"`
	IngressControllerMode string `yaml:"ingressControllerMode"`
	IngressService        string `yaml:"ingressService"`
	IngressSelector       string `yaml:"ingressSelector"`
}

// getHelmVal https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#MeshConfig
func getHelmVal() *helmValues {
	return &helmValues{
		HelmValuesMeshConfig: &helmValuesMeshConfig{
			IngressClass:          "istio",
			IngressControllerMode: "STRICT",
			IngressService:        "ingress-external",
			IngressSelector:       "ingress",
		},
	}
}

func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "istiod",
		Repo:        "https://istio-release.storage.googleapis.com/charts",
		Name:        "istiod",
		Version:     "1.15.0-beta.1",
	}
}
