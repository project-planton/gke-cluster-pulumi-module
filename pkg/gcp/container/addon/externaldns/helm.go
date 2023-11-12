package externaldns

import (
	helmcommons "github.com/plantoncloud-inc/go-commons/kubernetes/helm"
)

// helmVal https://github.com/kubernetes-sigs/external-dns/blob/438d06f3c45cf66d08945ae18d17e29c540d5c96/charts/external-dns/values.yaml
type helmVal struct {
	ServiceAccount *serviceAccount `yaml:"serviceAccount"`
}

type serviceAccount struct {
	Create bool   `yaml:"create"`
	Name   string `yaml:"name"`
}

// getHelmChart https://github.com/kubernetes-sigs/external-dns/tree/438d06f3c45cf66d08945ae18d17e29c540d5c96/charts/external-dns
func getHelmChart() *helmcommons.Chart {
	return &helmcommons.Chart{
		ReleaseName: "external-dns",
		Repo:        "https://kubernetes-sigs.github.io/external-dns/",
		Name:        "external-dns",
		//https://github.com/kubernetes-sigs/external-dns/blob/438d06f3c45cf66d08945ae18d17e29c540d5c96/charts/external-dns/Chart.yaml#L5
		Version: "1.13.1",
	}
}

func getHelmVal() *helmVal {
	return &helmVal{
		ServiceAccount: &serviceAccount{
			Create: false,
			Name:   Ksa,
		},
	}
}
