package ingressnginx

import "github.com/plantoncloud-inc/go-commons/kubernetes/helm"

const (
	Namespace        = "ingress-nginx"
	SvcTypeClusterIp = "ClusterIP"
)

type HelmVal struct {
	Controller *Controller `yaml:"controller"`
}

type Controller struct {
	Service              *Service              `yaml:"service"`
	IngressClassResource *IngressClassResource `yaml:"ingressClassResource"`
}

type Service struct {
	Type string `yaml:"type"`
}

type IngressClassResource struct {
	Default bool `yaml:"default"`
}

func getHelmVal() *HelmVal {
	return &HelmVal{
		Controller: &Controller{
			Service:              &Service{Type: SvcTypeClusterIp},
			IngressClassResource: &IngressClassResource{Default: true},
		},
	}
}

func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "ingress-nginx",
		Repo:        "https://kubernetes.github.io/ingress-nginx",
		Name:        "ingress-nginx",
		//https://github.com/kubernetes/ingress-nginx/blob/main/charts/ingress-nginx/Chart.yaml#L26C9-L26C14
		Version: "4.7.1",
	}
}
