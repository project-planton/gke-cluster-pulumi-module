package aws

import (
	"context"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

	"github.com/pkg/errors"
	c2cv1deployk8cstackawsmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/aws/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/stack/output/backend"
)

func Outputs(ctx context.Context, input *c2cv1deployk8cstackawsmodel.KubeClusterAwsStackInput) (*c2cv1deployk8cstackawsmodel.KubeClusterAwsStackOutputs, error) {
	stackOutput, err := backend.StackOutput(input.StackJob)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get stack output")
	}
	return OutputMapTransformer(stackOutput, input), nil
}

func OutputMapTransformer(stackOutput map[string]interface{}, input *c2cv1deployk8cstackawsmodel.KubeClusterAwsStackInput) *c2cv1deployk8cstackawsmodel.KubeClusterAwsStackOutputs {
	if input.StackJob.Spec.OperationType != stackjoboperationtype.StackJobOperationType_apply || stackOutput == nil {
		return &c2cv1deployk8cstackawsmodel.KubeClusterAwsStackOutputs{}
	}

	return &c2cv1deployk8cstackawsmodel.KubeClusterAwsStackOutputs{
		ClusterVpcId:    "coming-soon",
		ClusterEndpoint: "coming-soon",
		ClusterCaData:   "coming-soon",
	}
}
