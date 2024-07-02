package aws

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/aws/container/cluster"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/aws/network"
	c2cv1deployk8cstackawsmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/aws/model"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/aws/pulumiawsprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input            *c2cv1deployk8cstackawsmodel.KubeClusterAwsStackInput
	AwsLabels        map[string]string
	KubernetesLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	kubeCluster := s.Input.ResourceInput.KubeCluster

	awsClassicProvider, err := pulumiawsprovider.GetClassic(ctx,
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
