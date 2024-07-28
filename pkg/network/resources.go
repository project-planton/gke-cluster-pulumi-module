package network

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network/ip"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network/vpc"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/projects/project"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	KubeClusterId          string
	IsCreateSharedVpc      bool
	GcpRegion              string
	AddedProjectsResources *project.AddedProjectsResources
	Labels                 map[string]string
}

type AddedNetworkResources struct {
	AddedVpc                 *compute.Network
	AddedSubnet              *compute.Subnetwork
	AddedIpAddresses         *ip.AddedIngressIpAddresses
	AddedNetworkIamResources []pulumi.Resource
}

// Resources sets up org network by
// * creates a private vpc or a shared vpc
// * creates subnetwork
// * creates compute router
// * creates router nat
func Resources(ctx *pulumi.Context, input *Input) (*AddedNetworkResources, error) {
	addedVpcNetwork, err := vpc.Resources(ctx, &vpc.Input{
		KubeClusterId:          input.KubeClusterId,
		IsSharedVpcEnabled:     input.IsCreateSharedVpc,
		AddedProjectsResources: input.AddedProjectsResources,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add vpc resources")
	}
	addedIpAddresses, err := ip.Resources(ctx, &ip.Input{
		KubeClusterId:          input.KubeClusterId,
		GcpRegion:              input.GcpRegion,
		AddedProjectsResources: input.AddedProjectsResources,
		AddedSubnet:            addedSubnet,
		Labels:                 input.Labels,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add ip resources")
	}
	return nil, nil
}
