package project

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/projects/apis"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/project"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	KubeClusterId          string
	GcpRegion              string
	GcpZone                string
	GcpProvider            *pulumigcp.Provider
	AddedKubeClusterFolder *organizations.Folder
	Labels                 map[string]string
	IsSharedVpcEnabled     bool
	BillingAccountId       string
}

type AddedProjectsResources struct {
	KubeClusterProjects *AddedKubeClusterProjects
	AddedProjectApis    *AddedProjectsApis
}

type AddedProjectsApis struct {
	VpcNetworkProject       []pulumi.Resource
	ContainerClusterProject []pulumi.Resource
}

type AddedKubeClusterProjects struct {
	VpcNetworkProject       *organizations.Project
	ContainerClusterProject *organizations.Project
}

func Resources(ctx *pulumi.Context, input *Input) (*AddedProjectsResources, error) {
	//container cluster
	containerClusterProjectName := getContainerClusterProjectName(input.KubeClusterId)
	containerClusterProjectId, err := getContainerClusterProjectId(ctx, input.KubeClusterId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get kube-cluster project id")
	}
	containerClusterProject, err := addProject(ctx, input.GcpProvider,
		input.BillingAccountId,
		containerClusterProjectName,
		containerClusterProjectId,
		input.AddedKubeClusterFolder,
		input.Labels,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s project", containerClusterProjectName)
	}
	ctx.Export(ContainerClusterProjectIdOutputName, containerClusterProject.ProjectId)
	ctx.Export(ContainerClusterProjectNumberOutputName, containerClusterProject.Number)
	addedContainerClusterProjectApis, err := project.AddApi(ctx, containerClusterProjectName, containerClusterProject,
		apis.ContainerClusterProjectApis)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add services for %s project", containerClusterProjectName)
	}

	//skip creating shared-vpc project if disabled
	if !input.IsSharedVpcEnabled {
		//if shared-vpc is not enabled, vpc is created in the kube-cluster project
		ctx.Export(VpcNetworkProjectIdOutputName, containerClusterProject.ProjectId)
		ctx.Export(VpcNetworkProjectNumberOutputName, containerClusterProject.Number)
		return &AddedProjectsResources{
			KubeClusterProjects: &AddedKubeClusterProjects{
				VpcNetworkProject:       containerClusterProject,
				ContainerClusterProject: containerClusterProject,
			},
			AddedProjectApis: &AddedProjectsApis{
				VpcNetworkProject:       addedContainerClusterProjectApis,
				ContainerClusterProject: addedContainerClusterProjectApis,
			},
		}, nil
	}

	//vpc network project
	vpcNetworkProjectName := getVpcNetworkProjectName(input.KubeClusterId)
	vpcNetworkProjectId, err := getVpcNetworkProjectId(ctx, input.KubeClusterId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get vpc-network project id")
	}
	vpcNetworkProject, err := addProject(ctx, input.GcpProvider,
		input.BillingAccountId,
		vpcNetworkProjectName,
		vpcNetworkProjectId,
		input.AddedKubeClusterFolder,
		input.Labels,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s project", vpcNetworkProjectName)
	}
	ctx.Export(VpcNetworkProjectIdOutputName, vpcNetworkProject.ProjectId)
	ctx.Export(VpcNetworkProjectNumberOutputName, vpcNetworkProject.Number)
	addedVpcNetworkProjectApis, err := project.AddApi(ctx, vpcNetworkProjectName, vpcNetworkProject, apis.VpcNetworkProjectApis)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add services for %s project", vpcNetworkProjectName)
	}

	return &AddedProjectsResources{
		KubeClusterProjects: &AddedKubeClusterProjects{
			VpcNetworkProject:       vpcNetworkProject,
			ContainerClusterProject: containerClusterProject,
		},
		AddedProjectApis: &AddedProjectsApis{
			VpcNetworkProject:       addedVpcNetworkProjectApis,
			ContainerClusterProject: addedContainerClusterProjectApis,
		},
	}, nil
}

func addProject(ctx *pulumi.Context, gcpProvider *pulumigcp.Provider, gcpBillingAccountId, gcpProjectName string,
	gcpProjectId pulumi.StringOutput, addedKubeClusterFolder *organizations.Folder, labels map[string]string) (*organizations.Project, error) {

}

func getVpcNetworkProjectName(kubeClusterId string) string {
	return fmt.Sprintf("%s-nw", kubeClusterId)
}

func getContainerClusterProjectName(kubeClusterId string) string {
	return fmt.Sprintf("%s-co", kubeClusterId)
}

func getVpcNetworkProjectId(ctx *pulumi.Context, kubeClusterId string) (pulumi.StringOutput, error) {
	randomString, err := random.NewRandomString(ctx, fmt.Sprintf("%s-vpc-network-project-suffix", kubeClusterId),
		&random.RandomStringArgs{
			Special: pulumi.Bool(false),
			Lower:   pulumi.Bool(true),
			Upper:   pulumi.Bool(false),
			Number:  pulumi.Bool(true),
			Length:  pulumi.Int(2), //increasing this can result in violation of project name length <30
		})
	if err != nil {
		return pulumi.StringOutput{}, errors.Wrap(err, "failed to create random suffix for vpc-network project")
	}
	gcpProjectId := randomString.Result.ApplyT(func(suffix string) string {
		//project id is created by prefixing character "n" to the random string
		return fmt.Sprintf("%s-n%s", kubeClusterId, suffix)
	}).(pulumi.StringOutput)
	return gcpProjectId, nil
}

func getContainerClusterProjectId(ctx *pulumi.Context, kubeClusterId string) (pulumi.StringOutput, error) {
	randomString, err := random.NewRandomString(ctx, fmt.Sprintf("%s-container-cluster-project-suffix", kubeClusterId),
		&random.RandomStringArgs{
			Special: pulumi.Bool(false),
			Lower:   pulumi.Bool(true),
			Upper:   pulumi.Bool(false),
			Number:  pulumi.Bool(true),
			Length:  pulumi.Int(2), //increasing this can result in violation of project name length <30
		})
	if err != nil {
		return pulumi.StringOutput{}, errors.Wrap(err, "failed to create random suffix for kube-cluster project")
	}
	gcpProjectId := randomString.Result.ApplyT(func(suffix string) string {
		//project id is created by prefixing character "c" to the random string
		return fmt.Sprintf("%s-c%s", kubeClusterId, suffix)
	}).(pulumi.StringOutput)
	return gcpProjectId, nil
}
