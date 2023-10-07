package postgresoperator

import (
	"github.com/plantoncloud-inc/go-commons/kubernetes/helm"
	kuberneteslabels "github.com/plantoncloud-inc/go-commons/kubernetes/labels"
)

// PlantonCloudStandardLabels are applied to all the resources created by Planton Cloud.
// these labels should be added to postgres-operator configuration to ensure that the labels are added to pods managed by postgres-operator.
var PlantonCloudStandardLabels = []string{
	kuberneteslabels.Resource,
	kuberneteslabels.Company,
	kuberneteslabels.Product,
	kuberneteslabels.Environment,
	kuberneteslabels.ResourceType,
	kuberneteslabels.ResourceId,
}

type HelmVal struct {
	//https://github.com/zalando/postgres-operator/blob/af084a5a650527c43f0c0fc579551a741e77f5c8/charts/postgres-operator/values.yaml#L96
	ConfigKubernetes *ValConfigKubernetes `yaml:"configKubernetes"`
}

type ValConfigKubernetes struct {
	InheritLabels []string `yaml:"inherited_labels"`
}

func getHelmVal() *HelmVal {
	return &HelmVal{
		ConfigKubernetes: &ValConfigKubernetes{
			InheritLabels: PlantonCloudStandardLabels,
		},
	}
}

func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "postgres-operator",
		Repo:        "https://opensource.zalando.com/postgres-operator/charts/postgres-operator",
		Name:        "postgres-operator",
		//https://github.com/zalando/postgres-operator/releases/tag/v1.10.0
		Version: "1.10.0",
	}
}
