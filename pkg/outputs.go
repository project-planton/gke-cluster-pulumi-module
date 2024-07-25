package pkg

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/projects"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
)

func PulumiOutputToStackOutputsConverter(stackOutput map[string]interface{}, input *model.GkeClusterStackInput) *model.GkeClusterStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &model.GkeClusterStackOutputs{}
	}

	projectsOutputs := projects.Output(input.ResourceInput, stackOutput)
	networkOutputs := network.Output(input.ResourceInput, stackOutput)
	iamOutputs := iam.Output(input.ResourceInput, stackOutput)
	containerClusterOutputs := cluster.Output(input.ResourceInput, stackOutput)

	return &model.GkeClusterStackOutputs{
		Projects:  projectsOutputs,
		Network:   networkOutputs,
		Iam:       iamOutputs,
		Container: containerClusterOutputs,
	}
}
