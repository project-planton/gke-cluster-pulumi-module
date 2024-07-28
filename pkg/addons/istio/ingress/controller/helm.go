package controller

import (
	"github.com/plantoncloud-inc/go-commons/kubernetes/helm"
)

const (
	IstioStatusPort = 15021
	HttpPort        = 80
	HttpsPort       = 443
	DebugPort       = 5005
)

// HelmVal https://github.com/istio/istio/blob/master/manifests/charts/gateways/istio-ingress/values.yaml
type HelmVal struct {
	Service *IstioIngressHelmValIngressService
}

type IstioIngressHelmValIngressService struct {
	Type           string                                   `yaml:"type"`
	LoadBalancerIp string                                   `yaml:"loadBalancerIP"`
	Ports          []*IstioIngressHelmValIngressServicePort `yaml:"ports"`
}

type IstioIngressHelmValIngressServicePort struct {
	Name       string `yaml:"name"`
	Protocol   string `yaml:"protocol"`
	Port       int32  `yaml:"port"`
	TargetPort int32  `yaml:"targetPort"`
	NodePort   int32  `yaml:"nodePort"`
}

func getHelmVal() *HelmVal {
	cv := &HelmVal{
		Service: &IstioIngressHelmValIngressService{
			Type: "ClusterIP",
			Ports: []*IstioIngressHelmValIngressServicePort{
				{
					Name:       "status-port",
					Protocol:   "TCP",
					Port:       IstioStatusPort,
					TargetPort: IstioStatusPort,
				},
				{
					Name:       "http2",
					Protocol:   "TCP",
					Port:       HttpPort,
					TargetPort: HttpPort,
				},
				{
					Name:       "https",
					Protocol:   "TCP",
					Port:       HttpsPort,
					TargetPort: HttpsPort,
				},
				{
					Name:       "debug",
					Protocol:   "TCP",
					Port:       DebugPort,
					TargetPort: DebugPort,
				},
			},
		},
	}
	return cv
}

func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "istio-ingress",
		Repo:        "https://istio-release.storage.googleapis.com/charts",
		Name:        "gateway",
		//https://github.com/istio/istio/releases
		Version: "1.18.0",
	}
}
