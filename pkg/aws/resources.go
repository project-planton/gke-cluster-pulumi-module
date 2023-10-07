package aws

import (
	"buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/aws"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/aws/container/cluster"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/aws/network"
	pulumiawsnativeprovider "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/automation/provider/aws"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input            *aws.KubeClusterAwsStackInput
	AwsLabels        map[string]string
	KubernetesLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	kubeCluster := s.Input.ResourceInput.KubeCluster

	awsClassicProvider, err := pulumiawsnativeprovider.GetClassic(ctx,
		s.Input.CredentialsInput.Aws, kubeCluster.Spec.Aws.Region)
	if err != nil {
		return errors.Wrap(err, "failed to setup aws provider")
	}

	addedNetworkResources, err := network.Resources(ctx, &network.Input{
		AwsProvider: awsClassicProvider,
		KubeCluster: kubeCluster,
		Labels:      s.AwsLabels,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add network resources")
	}

	err = cluster.Resources(ctx, &cluster.Input{
		AwsProvider:           awsClassicProvider,
		KubeCluster:           kubeCluster,
		Labels:                s.AwsLabels,
		AddedNetworkResources: addedNetworkResources,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add cluster resources")
	}
	return nil
}
