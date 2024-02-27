package network

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/aws/network/vpc"
	code2cloudv1deployk8cmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/model"
	awsclassic "github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	AwsProvider *awsclassic.Provider
	KubeCluster *code2cloudv1deployk8cmodel.KubeCluster
	Labels      map[string]string
}

type AddedResources struct {
	AddedVpc *ec2.Vpc
}

// Resources sets up network by
// * optionally creates a vpc depending on the presence of value for vpc-id in the input
func Resources(ctx *pulumi.Context, input *Input) (*AddedResources, error) {
	addedVpc, err := vpc.Resources(ctx, &vpc.Input{
		AwsProvider: input.AwsProvider,
		KubeCluster: input.KubeCluster,
		Labels:      input.Labels,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add vpc resources")
	}
	return &AddedResources{AddedVpc: addedVpc}, nil
}
