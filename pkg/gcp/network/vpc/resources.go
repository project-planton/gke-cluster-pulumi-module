package vpc

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/projects/project"
	puluminameoutputgcp "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	wordpb "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/commons/english/rpc/enums"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	KubeClusterId          string
	IsSharedVpcEnabled     bool
	AddedProjectsResources *project.AddedProjectsResources
}

func Resources(ctx *pulumi.Context, input *Input) (*compute.Network, error) {
	nw, err := addNetwork(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add network")
	}
	//skip creating shared vpc resources when shared vpc is disabled.
	if !input.IsSharedVpcEnabled {
		return nw, nil
	}
	hostProject, err := addSharedVpcHostProject(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add shared vpc host project")
	}
	if err := addSharedVpcServiceProjects(ctx, input, hostProject); err != nil {
		return nil, errors.Wrap(err, "failed to add shared vpc service projects")
	}
	return nw, nil
}

// addSharedVpcServiceProjects adds kube-cluster project as a service project to the vpc-network project
func addSharedVpcServiceProjects(ctx *pulumi.Context, input *Input, hostProject *compute.SharedVPCHostProject) error {
	_, err := compute.NewSharedVPCServiceProject(ctx, fmt.Sprintf("%s-%s", wordpb.Word_kubernetes, input.KubeClusterId),
		&compute.SharedVPCServiceProjectArgs{
			HostProject:    hostProject.Project,
			ServiceProject: input.AddedProjectsResources.KubeClusterProjects.ContainerClusterProject.ProjectId,
		}, pulumi.Parent(hostProject), pulumi.DependsOn([]pulumi.Resource{hostProject}),
		pulumi.DependsOn(input.AddedProjectsResources.AddedProjectApis.VpcNetworkProject),
	)
	if err != nil {
		return errors.Wrap(err, "failed to add kube-cluster project as service project")
	}
	return nil
}

func addSharedVpcHostProject(ctx *pulumi.Context, input *Input) (*compute.SharedVPCHostProject, error) {
	hostProject, err := compute.NewSharedVPCHostProject(ctx,
		fmt.Sprintf("%s-%s-host", wordpb.Word_kubernetes, input.KubeClusterId),
		&compute.SharedVPCHostProjectArgs{
			Project: input.AddedProjectsResources.KubeClusterProjects.VpcNetworkProject.ProjectId,
		}, pulumi.Parent(input.AddedProjectsResources.KubeClusterProjects.VpcNetworkProject),
		pulumi.DependsOn(input.AddedProjectsResources.AddedProjectApis.VpcNetworkProject))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add host project")
	}
	return hostProject, nil
}

func addNetwork(ctx *pulumi.Context, input *Input) (*compute.Network, error) {
	name := GetNetworkName(input.KubeClusterId)
	nw, err := compute.NewNetwork(ctx, name, &compute.NetworkArgs{
		Project:               input.AddedProjectsResources.KubeClusterProjects.VpcNetworkProject.ProjectId,
		AutoCreateSubnetworks: pulumi.BoolPtr(false),
	}, pulumi.Parent(input.AddedProjectsResources.KubeClusterProjects.VpcNetworkProject),
		pulumi.DependsOn(input.AddedProjectsResources.AddedProjectApis.VpcNetworkProject))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create network")
	}
	ctx.Export(GetNetworkSelfLinkOutputName(name), nw.SelfLink)
	return nw, nil
}

func GetNetworkName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s", wordpb.Word_kubernetes, kubeClusterId)
}

func GetNetworkSelfLinkOutputName(networkName string) string {
	return puluminameoutputgcp.Name(compute.Network{}, networkName)
}
