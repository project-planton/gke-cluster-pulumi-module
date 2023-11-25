package aws

import (
	"context"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/stack/job/enums/operationtype"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/aws"
)

func Outputs(ctx context.Context, input *aws.KubeClusterAwsStackInput) (*aws.KubeClusterAwsStackOutputs, error) {
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

func Get(stackOutput map[string]interface{}, input *aws.KubeClusterAwsStackInput) *aws.KubeClusterAwsStackOutputs {
	if input.StackJob.Spec.OperationType != operationtype.StackJobOperationType_apply || stackOutput == nil {
		return &aws.KubeClusterAwsStackOutputs{}
	}

	return &aws.KubeClusterAwsStackOutputs{
		ClusterVpcId:    "coming-soon",
		ClusterEndpoint: "coming-soon",
		ClusterCaData:   "coming-soon",
	}
}
