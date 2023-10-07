package network

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/network/firewall"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/network/iam"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/network/ip"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/network/nat"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/network/router"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/network/subnet"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/network/vpc"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/projects/project"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
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
	AddedIpAddresses         *ip.AddedComputeIpAddresses
	AddedNetworkIamResources []pulumi.Resource
}

// Resources sets up org network by
// * creates a private vpc or a shared vpc
// * creates subnetwork
// * creates compute router
// * creates router nat
// * create iam resources to allow the Google container engine service account in kube-cluster project to update firewall rules in shared project
func Resources(ctx *pulumi.Context, input *Input) (*AddedNetworkResources, error) {
	addedVpcNetwork, err := vpc.Resources(ctx, &vpc.Input{
		KubeClusterId:          input.KubeClusterId,
		IsSharedVpcEnabled:     input.IsCreateSharedVpc,
		AddedProjectsResources: input.AddedProjectsResources,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add vpc resources")
	}
	if err := firewall.Resources(ctx, &firewall.Input{
		KubeClusterId:     input.KubeClusterId,
		AddedShareProject: input.AddedProjectsResources.KubeClusterProjects.VpcNetworkProject,
		AddedVpcNetwork:   addedVpcNetwork,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to add network firewall")
	}
	addedSubnet, err := subnet.Resources(ctx, &subnet.Input{
		KubeClusterId: input.KubeClusterId,
		GcpRegion:     input.GcpRegion,
		ShareProject:  input.AddedProjectsResources.KubeClusterProjects.VpcNetworkProject,
		VpcNetwork:    addedVpcNetwork,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add subnetwork")
	}

	addedNetworkRouter, err := router.Resources(ctx, &router.Input{
		KubeClusterId:          input.KubeClusterId,
		GcpRegion:              input.GcpRegion,
		AddedVpcNetworkProject: input.AddedProjectsResources.KubeClusterProjects.VpcNetworkProject,
		VpcNetwork:             addedVpcNetwork,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add compute router")
	}
	if err := nat.Resources(ctx, &nat.Input{
		KubeClusterId:          input.KubeClusterId,
		GcpRegion:              input.GcpRegion,
		AddedVpcNetworkProject: input.AddedProjectsResources.KubeClusterProjects.VpcNetworkProject,
		AddedNetworkRouter:     addedNetworkRouter,
		Labels:                 input.Labels,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to add compute router nat")
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
	//skip creating iam resources as they are only needed in case of shared-vpc
	if !input.IsCreateSharedVpc {
		return &AddedNetworkResources{
			AddedVpc:                 addedVpcNetwork,
			AddedSubnet:              addedSubnet,
			AddedIpAddresses:         addedIpAddresses,
			AddedNetworkIamResources: []pulumi.Resource{},
		}, nil
	}
	addedNetworkIamResources, err := iam.Resources(ctx, &iam.Input{
		GcpRegion:                input.GcpRegion,
		AddedSubnet:              addedSubnet,
		AddedKubeClusterProjects: input.AddedProjectsResources.KubeClusterProjects,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add iam resources")
	}
	return &AddedNetworkResources{
		AddedVpc:                 addedVpcNetwork,
		AddedSubnet:              addedSubnet,
		AddedIpAddresses:         addedIpAddresses,
		AddedNetworkIamResources: addedNetworkIamResources,
	}, nil
}
