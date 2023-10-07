package plantoncloudkubeagent

import (
	gcp "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	"github.com/plantoncloud-inc/go-commons/kubernetes/helm"
)

const (
	Namespace = "planton-cloud-kube-agent"
)

type HelmVal struct {
	Image                               string `yaml:"image"`
	Company                             string `yaml:"company"`
	KubeClusterId                       string `yaml:"kubeClusterId"`
	MachineAccountEmail                 string `yaml:"machineAccountEmail"`
	ClientSecret                        string `yaml:"clientSecret"`
	PlantonCloudServiceApiEndpoint      string `yaml:"plantonCloudServiceApiEndpoint"`
	OpenCostApiEndpoint                 string `yaml:"openCostApiEndpoint"`
	OpenCostPollingIntervalSeconds      int32  `yaml:"openCostPollingIntervalSeconds"`
	TokenExpirationBufferMinutes        int32  `yaml:"tokenExpirationBufferMinutes"`
	TokenExpirationCheckIntervalSeconds int32  `yaml:"tokenExpirationCheckIntervalSeconds"`
}

func getHelmVal(input *gcp.AddonsPlantonCloudKubeAgent) *HelmVal {
	return &HelmVal{
		Image:                               input.DockerImage,
		Company:                             input.CompanyId,
		KubeClusterId:                       input.KubeClusterId,
		MachineAccountEmail:                 input.MachineAccountEmail,
		ClientSecret:                        input.ClientSecret,
		PlantonCloudServiceApiEndpoint:      input.PlantonCloudServiceApiEndpoint,
		OpenCostApiEndpoint:                 input.OpenCostApiEndpoint,
		OpenCostPollingIntervalSeconds:      input.OpenCostPollingIntervalSeconds,
		TokenExpirationBufferMinutes:        input.TokenExpirationBufferMinutes,
		TokenExpirationCheckIntervalSeconds: input.TokenExpirationCheckIntervalSeconds,
	}
}

// https://github.com/plantoncloud/helm-charts/tree/main/planton-cloud-kube-agent
func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "planton-cloud-kube-agent",
		Repo:        "https://plantoncloud.github.io/helm-charts/planton-cloud-kube-agent/",
		Name:        "planton-cloud-kube-agent",
		Version:     "v0.0.7",
	}
}