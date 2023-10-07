package iam

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/projects/project"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	networkAdminRoleName = "network.admin"
)

type Input struct {
	GcpRegion                string
	AddedSubnet              *compute.Subnetwork
	AddedKubeClusterProjects *project.AddedKubeClusterProjects
}

// Resources implements role creation as explained in the below blog post
// //https://www.linkedin.com/pulse/fixing-gkes-load-balancing-permissions-when-using-shared-dmitri-lerko/
func Resources(ctx *pulumi.Context, input *Input) ([]pulumi.Resource, error) {
	addedResources := make([]pulumi.Resource, 0)
	_, err := addNetworkAdminRole(ctx, input.AddedSubnet, input.AddedKubeClusterProjects.VpcNetworkProject)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add network admin role")
	}

	addedSubnetIamResources, err := addSubnetIam(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add iam policy for subnet work")
	}
	addedResources = append(addedResources, addedSubnetIamResources...)

	addedSubnetServiceAgentUserRoleBindings, err := addHostSvcAgentUserRoleBinding(ctx, input.AddedSubnet, input.AddedKubeClusterProjects)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add host svc agent roles")
	}
	addedResources = append(addedResources, addedSubnetServiceAgentUserRoleBindings...)

	addedNetworkAdminPolicyBinding, err := addNetworkAdminIamPolicyBinding(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add iam policy binding for service project container bot service accounts")
	}
	return append(addedResources, addedNetworkAdminPolicyBinding), nil
}

func addNetworkAdminRole(ctx *pulumi.Context, snw *compute.Subnetwork, addedShareProject *organizations.Project) (*projects.IAMCustomRole, error) {
	role, err := projects.NewIAMCustomRole(
		ctx,
		networkAdminRoleName,
		&projects.IAMCustomRoleArgs{
			Description: pulumi.String("This role allows to administer network and security of the host project. Intended for use by GKE automation on service projects."),
			Project:     addedShareProject.ProjectId,
			Permissions: pulumi.StringArray{
				pulumi.String("compute.firewalls.create"),
				pulumi.String("compute.firewalls.delete"),
				pulumi.String("compute.firewalls.get"),
				pulumi.String("compute.firewalls.list"),
				pulumi.String("compute.firewalls.update"),
				pulumi.String("compute.networks.updatePolicy"),
			},
			RoleId: pulumi.String(networkAdminRoleName),
			Title:  pulumi.String("Host Project Network and Security Admin"),
		}, pulumi.Parent(snw))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s custom iam role", networkAdminRoleName)
	}
	return role, nil
}

// addNetworkAdminIamPolicyBinding binds network admin role to container engine robot service accounts
// auto created for each service project.
func addNetworkAdminIamPolicyBinding(ctx *pulumi.Context, input *Input) (*projects.IAMBinding, error) {
	addedIamBinding, err := projects.NewIAMBinding(
		ctx,
		networkAdminRoleName,
		&projects.IAMBindingArgs{
			Members: getNetworkAdminIamBindingMembers(input.AddedKubeClusterProjects),
			Project: input.AddedKubeClusterProjects.VpcNetworkProject.ProjectId,
			Role: pulumi.Sprintf(
				"projects/%s/roles/%s",
				input.AddedKubeClusterProjects.VpcNetworkProject.ProjectId,
				networkAdminRoleName,
			),
		}, pulumi.Parent(input.AddedSubnet))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create role binding for %s role", networkAdminRoleName)
	}
	return addedIamBinding, nil
}

// getNetworkAdminIamBindingMembers returns the members for binding network admin role
// for container engine robot service account
func getNetworkAdminIamBindingMembers(addedKubeClusterProjects *project.AddedKubeClusterProjects) pulumi.StringArray {
	resp := make([]pulumi.StringInput, 0)
	resp = append(resp, pulumi.Sprintf(
		"serviceAccount:service-%s@container-engine-robot.iam.gserviceaccount.com",
		addedKubeClusterProjects.ContainerClusterProject.Number,
	))
	return resp
}

//   - serviceAccount:SERVICE_PROJECT_1_NUM@cloudservices.gserviceaccount.com
//   - serviceAccount:service-SERVICE_PROJECT_1_NUM@container-engine-robot.iam.gserviceaccount.com
//
// https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-shared-vpc#enabling_and_granting_roles
func addSubnetIam(ctx *pulumi.Context, input *Input) ([]pulumi.Resource, error) {
	addedIamMemberSubnetCloudServices, err := compute.NewSubnetworkIAMMember(
		ctx,
		"subnetwork-iam-policy-cloudservices",
		&compute.SubnetworkIAMMemberArgs{
			Member: pulumi.Sprintf(
				"serviceAccount:%s@cloudservices.gserviceaccount.com",
				input.AddedKubeClusterProjects.ContainerClusterProject.Number,
			),
			Project:    input.AddedKubeClusterProjects.VpcNetworkProject.ProjectId,
			Region:     pulumi.String(input.GcpRegion),
			Role:       pulumi.String(standard.Compute_networkUser),
			Subnetwork: input.AddedSubnet.SelfLink,
		}, pulumi.Parent(input.AddedSubnet))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add gke service accounts as iam members for subnetwork")
	}
	addedIamMemberSubnetContainerEngine, err := compute.NewSubnetworkIAMMember(
		ctx,
		"subnetwork-iam-policy-container-engine-robot",
		&compute.SubnetworkIAMMemberArgs{
			Member: pulumi.Sprintf(
				"serviceAccount:service-%s@container-engine-robot.iam.gserviceaccount.com",
				input.AddedKubeClusterProjects.ContainerClusterProject.Number,
			),
			Project:    input.AddedKubeClusterProjects.VpcNetworkProject.ProjectId,
			Region:     pulumi.String(input.GcpRegion),
			Role:       pulumi.String(standard.Compute_networkUser),
			Subnetwork: input.AddedSubnet.SelfLink,
		}, pulumi.Parent(input.AddedSubnet))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add gke service accounts as iam members for subnetwork")
	}
	return append([]pulumi.Resource{
		addedIamMemberSubnetCloudServices,
		addedIamMemberSubnetContainerEngine,
	}), nil
}

func addHostSvcAgentUserRoleBinding(ctx *pulumi.Context, addedSubnet *compute.Subnetwork, addedKubeClusterProjects *project.AddedKubeClusterProjects) ([]pulumi.Resource, error) {
	addedIamMemberContainerEngineServiceAgent, err := projects.NewIAMMember(ctx,
		"host-service-agent-role",
		&projects.IAMMemberArgs{
			Member: pulumi.Sprintf(
				"serviceAccount:service-%s@container-engine-robot.iam.gserviceaccount.com",
				addedKubeClusterProjects.ContainerClusterProject.Number,
			),
			Project: addedKubeClusterProjects.VpcNetworkProject.ProjectId,
			Role:    pulumi.String(standard.Container_hostServiceAgentUser),
		}, pulumi.Parent(addedSubnet))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add network host service agent role")
	}
	return []pulumi.Resource{
		addedIamMemberContainerEngineServiceAgent,
	}, nil
}
