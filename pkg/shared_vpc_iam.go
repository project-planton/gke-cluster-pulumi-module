package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// sharedVpcIam sets up IAM permissions as explained in
// https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-shared-vpc#managing_firewall_resources
func sharedVpcIam(ctx *pulumi.Context,
	createdNetworkProject *organizations.Project,
	createdSubNetwork *compute.Subnetwork) ([]pulumi.Resource, error) {
	_, err := projects.NewIAMCustomRole(
		ctx,
		"network-admin-role",
		&projects.IAMCustomRoleArgs{
			Description: pulumi.String("This role allows to administer network and security of the host project. " +
				"Intended for use by GKE automation on service projects."),
			Project: createdNetworkProject.ProjectId,
			Permissions: pulumi.StringArray{
				pulumi.String("compute.firewalls.create"),
				pulumi.String("compute.firewalls.delete"),
				pulumi.String("compute.firewalls.get"),
				pulumi.String("compute.firewalls.list"),
				pulumi.String("compute.firewalls.update"),
				pulumi.String("compute.networks.updatePolicy"),
			},
			RoleId: pulumi.String("network.admin"),
			Title:  pulumi.String("Host Project Network and Security Admin"),
		}, pulumi.Parent(createdSubNetwork))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create custom-iam role for network-admin-role on host project")
	}

	//   - serviceAccount:SERVICE_PROJECT_1_NUM@cloudservices.gserviceaccount.com
	//   - serviceAccount:service-SERVICE_PROJECT_1_NUM@container-engine-robot.iam.gserviceaccount.com
	//
	// https://cloud.google.com/kubernetes-engine/docs/how-to/cluster-shared-vpc#enabling_and_granting_roles
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

	return nil, nil
}
