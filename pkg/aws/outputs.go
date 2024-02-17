package aws

import (
	"context"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/stack/job/enums/operationtype"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/org"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	c2cv1deployk8cstackawsmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/aws/model"
)

func Outputs(ctx context.Context, input *c2cv1deployk8cstackawsmodel.KubeClusterAwsStackInput) (*c2cv1deployk8cstackawsmodel.KubeClusterAwsStackOutputs, error) {
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

func Get(stackOutput map[string]interface{}, input *c2cv1deployk8cstackawsmodel.KubeClusterAwsStackInput) *c2cv1deployk8cstackawsmodel.KubeClusterAwsStackOutputs {
	if input.StackJob.Spec.OperationType != operationtype.StackJobOperationType_apply || stackOutput == nil {
		return &c2cv1deployk8cstackawsmodel.KubeClusterAwsStackOutputs{}
	}

	return &c2cv1deployk8cstackawsmodel.KubeClusterAwsStackOutputs{
		ClusterVpcId:    "coming-soon",
		ClusterEndpoint: "coming-soon",
		ClusterCaData:   "coming-soon",
	}
}
