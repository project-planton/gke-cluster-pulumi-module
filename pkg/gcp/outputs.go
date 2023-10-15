package gcp

import (
	"context"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/iam"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/network"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/projects"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/stack/rpc/enums"
)

func Outputs(ctx context.Context, input *gcp.KubeClusterGcpStackInput) (*gcp.KubeClusterGcpStackOutputs, error) {
	pulumiOrgName, err := org.GetOrgName()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pulumi org name")
	}
	stackOutput, err := backend.StackOutput(pulumiOrgName, input.StackJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stack output")
	}
	return Get(stackOutput, input), nil
}

func Get(stackOutput map[string]interface{}, input *gcp.KubeClusterGcpStackInput) *gcp.KubeClusterGcpStackOutputs {
	if input.StackJob.OperationType != enums.StackOperationType_apply || stackOutput == nil {
		return &gcp.KubeClusterGcpStackOutputs{}
	}

	projectsOutputs := projects.Output(input.ResourceInput, stackOutput)
	networkOutputs := network.Output(input.ResourceInput, stackOutput)
	iamOutputs := iam.Output(input.ResourceInput, stackOutput)
	containerClusterOutputs := cluster.Output(input.ResourceInput, stackOutput)

	return &gcp.KubeClusterGcpStackOutputs{
		Projects:  projectsOutputs,
		Network:   networkOutputs,
		Iam:       iamOutputs,
		Container: containerClusterOutputs,
	}
}
