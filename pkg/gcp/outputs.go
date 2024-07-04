package gcp

import (
	"context"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/projects"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/gcp/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/stack/output/backend"
)

func Outputs(ctx context.Context, input *c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackInput) (*c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackOutputs, error) {
	stackOutput, err := backend.StackOutput(input.StackJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stack output")
	}
	return OutputMapTransformer(stackOutput, input), nil
}

func OutputMapTransformer(stackOutput map[string]interface{}, input *c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackInput) *c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackOutputs{}
	}

	projectsOutputs := projects.Output(input.ResourceInput, stackOutput)
	networkOutputs := network.Output(input.ResourceInput, stackOutput)
	iamOutputs := iam.Output(input.ResourceInput, stackOutput)
	containerClusterOutputs := cluster.Output(input.ResourceInput, stackOutput)

	return &c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackOutputs{
		Projects:  projectsOutputs,
		Network:   networkOutputs,
		Iam:       iamOutputs,
		Container: containerClusterOutputs,
	}
}
