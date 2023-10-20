package cluster

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/aws/container/cluster/cluster"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/aws/network"
	kubeclusterv1 "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster"
	awsclassic "github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	KubeClusterId         string
	KubeCluster           *kubeclusterv1.KubeCluster
	Labels                map[string]string
	AddedNetworkResources *network.AddedResources
	AwsProvider           *awsclassic.Provider
}

func Resources(ctx *pulumi.Context, input *Input) error {
	err := cluster.Resources(ctx, &cluster.Input{
		AwsProvider:           input.AwsProvider,
		KubeCluster:           input.KubeCluster,
		Labels:                input.Labels,
		AddedNetworkResources: input.AddedNetworkResources,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add cluster resources")
	}
	return nil
}
