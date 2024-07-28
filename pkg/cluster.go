package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/localz"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/outputs"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/vars"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/enums/gkereleasechannel"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/projects"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func cluster(ctx *pulumi.Context, locals *localz.Locals,
	createdFolder *organizations.Folder) (*container.Cluster, error) {

	//create random suffix for container-cluster-project-id
	clusterProjectRandomString, err := random.NewRandomString(ctx,
		"cluster-project-id-suffix",
		&random.RandomStringArgs{
			Special: pulumi.Bool(false),
			Lower:   pulumi.Bool(true),
			Upper:   pulumi.Bool(false),
			Numeric: pulumi.Bool(true),
			Length:  pulumi.Int(2), //increasing this can result in violation of project name length <30
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create random suffix for cluster-project-id")
	}

	//build container-cluster-project-id using the random-id suffix
	clusterProjectId := clusterProjectRandomString.Result.ApplyT(func(suffix string) string {
		//project id is created by prefixing character "c" to the random string.
		//this is to easily distinguish between network project and cluster project in shared-vpc setup.
		return fmt.Sprintf("%s-c%s", locals.GkeCluster.Metadata.Id, suffix)
	}).(pulumi.StringOutput)

	//create container-cluster project
	createdClusterProject, err := organizations.NewProject(ctx,
		"cluster-project",
		&organizations.ProjectArgs{
			BillingAccount:    pulumi.String(locals.GkeCluster.Spec.BillingAccountId),
			Name:              clusterProjectId,
			AutoCreateNetwork: pulumi.Bool(false),
			Labels:            pulumi.ToStringMap(locals.GcpLabels),
			ProjectId:         clusterProjectId,
			FolderId:          createdFolder.FolderId,
		}, pulumi.Parent(createdFolder))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cluster project")
	}

	var createdNetworkProject *organizations.Project

	if !locals.GkeCluster.Spec.IsCreateSharedVpc {
		//when the cluster does not need to have shared-vpc, both cluster and network are created
		//in the same gcp project.
		createdNetworkProject = createdClusterProject
	} else {
		//create random suffix for network-cluster-project-id
		networkProjectRandomString, err := random.NewRandomString(ctx,
			"network-project-id-suffix",
			&random.RandomStringArgs{
				Special: pulumi.Bool(false),
				Lower:   pulumi.Bool(true),
				Upper:   pulumi.Bool(false),
				Numeric: pulumi.Bool(true),
				Length:  pulumi.Int(2), //increasing this can result in violation of project name length <30
			})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create random suffix for network-project-id")
		}

		//build network-project-id suffix using its random-id suffix
		networkProjectId := networkProjectRandomString.Result.ApplyT(func(suffix string) string {
			//project id is created by prefixing character "c" to the random string.
			//this is to easily distinguish between network project and cluster project in shared-vpc setup.
			return fmt.Sprintf("%s-n%s", locals.GkeCluster.Metadata.Id, suffix)
		}).(pulumi.StringOutput)

		//create network project
		createdNetworkProject, err = organizations.NewProject(ctx,
			"network-project",
			&organizations.ProjectArgs{
				BillingAccount:    pulumi.String(locals.GkeCluster.Spec.BillingAccountId),
				Name:              networkProjectId,
				AutoCreateNetwork: pulumi.Bool(false),
				Labels:            pulumi.ToStringMap(locals.GcpLabels),
				ProjectId:         networkProjectId,
				FolderId:          createdFolder.FolderId,
			}, pulumi.Parent(createdFolder))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create network project")
		}
	}

	//keep track of all the apis enabled to add as dependencies
	createdGoogleApiResources := make([]pulumi.Resource, 0)

	//enable apis for container cluster project
	for _, api := range vars.ContainerClusterProjectApis {
		addedProjectService, err := projects.NewService(ctx,
			fmt.Sprintf("container-cluster-%s", api),
			&projects.ServiceArgs{
				Project:                  createdClusterProject.ProjectId,
				DisableDependentServices: pulumi.BoolPtr(true),
				Service:                  pulumi.String(api),
			}, pulumi.Parent(createdClusterProject))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to enable %s api for container cluster project", api)
		}
		createdGoogleApiResources = append(createdGoogleApiResources, addedProjectService)
	}

	//enable apis for network project
	for _, api := range vars.NetworkProjectApis {
		addedProjectService, err := projects.NewService(ctx,
			fmt.Sprintf("container-cluster-%s", api),
			&projects.ServiceArgs{
				Project:                  createdNetworkProject.ProjectId,
				DisableDependentServices: pulumi.BoolPtr(true),
				Service:                  pulumi.String(api),
			}, pulumi.Parent(createdClusterProject))
		if err != nil {
			return nil, errors.Wrapf(err, "failed to enable %s api for network project", api)
		}
		createdGoogleApiResources = append(createdGoogleApiResources, addedProjectService)
	}

	//create vpc network
	createdNetwork, err := compute.NewNetwork(ctx,
		"vpc",
		&compute.NetworkArgs{
			Project:               createdNetworkProject.ProjectId,
			AutoCreateSubnetworks: pulumi.BoolPtr(false),
		}, pulumi.Parent(createdNetworkProject))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create network")
	}

	//export network self-link
	ctx.Export(outputs.NetworkSelfLink, createdNetwork.SelfLink)

	//create subnetwork
	createdSubNetwork, err := compute.NewSubnetwork(ctx, "sub-network", &compute.SubnetworkArgs{
		Name:                  pulumi.String(locals.GkeCluster.Metadata.Id),
		Project:               createdNetworkProject.ProjectId,
		Network:               createdNetwork.ID(),
		Region:                pulumi.String(locals.GkeCluster.Spec.Region),
		IpCidrRange:           pulumi.String(vars.SubNetworkCidr),
		PrivateIpGoogleAccess: pulumi.BoolPtr(true),
		//these two ranges will be referred in the cluster input
		SecondaryIpRanges: &compute.SubnetworkSecondaryIpRangeArray{
			&compute.SubnetworkSecondaryIpRangeArgs{
				RangeName:   pulumi.String(locals.KubernetesPodSecondaryIpRangeName),
				IpCidrRange: pulumi.String(vars.KubernetesPodSecondaryIpRange),
			},
			&compute.SubnetworkSecondaryIpRangeArgs{
				RangeName:   pulumi.String(locals.KubernetesServiceSecondaryIpRangeName),
				IpCidrRange: pulumi.String(vars.KubernetesServiceSecondaryIpRange),
			},
		},
	}, pulumi.Parent(createdNetwork))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create subnetwork")
	}

	//export subnetwork self-link
	ctx.Export(outputs.SubNetworkSelfLink, createdSubNetwork.SelfLink)

	//create firewall
	createdFirewall, err := compute.NewFirewall(ctx, "firewall", &compute.FirewallArgs{
		Name:    pulumi.Sprintf("%s-gke-webhook", locals.GkeCluster.Metadata.Id),
		Project: createdNetworkProject.ProjectId,
		Network: createdNetwork.Name,
		SourceRanges: pulumi.StringArray{
			pulumi.String(vars.ApiServerIpCidr),
		},
		Allows: compute.FirewallAllowArray{
			&compute.FirewallAllowArgs{
				Protocol: pulumi.String("tcp"),
				Ports: pulumi.StringArray{
					pulumi.String(vars.ApiServerWebhookPort),
					pulumi.String(vars.IstioPilotWebhookPort),
				},
			},
		},
		TargetTags: pulumi.StringArray{
			pulumi.String(locals.NetworkTag),
		},
	}, pulumi.Parent(createdNetwork))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create firewall")
	}

	//export firewall self-link
	ctx.Export(outputs.GkeWebhooksFirewallSelfLink, createdFirewall.SelfLink)

	//create router
	createdRouter, err := compute.NewRouter(ctx,
		"router",
		&compute.RouterArgs{
			Name:    pulumi.String(locals.GkeCluster.Metadata.Id),
			Network: createdNetwork.SelfLink,
			Region:  pulumi.String(locals.GkeCluster.Spec.Region),
			Project: createdNetworkProject.ProjectId,
		}, pulumi.Parent(createdNetwork))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create router")
	}

	//export router self-link
	ctx.Export(outputs.RouterSelfLink, createdRouter.SelfLink)

	//create ip-address for router nat
	createdRouterNatIp, err := compute.NewAddress(ctx,
		"router-nat-ip",
		&compute.AddressArgs{
			Name:        pulumi.Sprintf("%s-router-nat", locals.GkeCluster.Metadata.Id),
			Project:     createdNetworkProject.ProjectId,
			Region:      createdRouter.Region,
			AddressType: pulumi.String("external"),
			Labels:      pulumi.ToStringMap(locals.GcpLabels),
		}, pulumi.Parent(createdRouter))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add new compute address")
	}

	//export router nat ip
	ctx.Export(outputs.NatIpAddress, createdRouterNatIp.Address)

	//create router nat
	createdRouterNat, err := compute.NewRouterNat(ctx,
		"nat-router",
		&compute.RouterNatArgs{
			Name:                          pulumi.String(locals.GkeCluster.Metadata.Id),
			Router:                        createdRouter.Name,
			Region:                        createdRouter.Region,
			Project:                       createdNetworkProject.ProjectId,
			NatIpAllocateOption:           pulumi.String("MANUAL_ONLY"),
			NatIps:                        pulumi.StringArray{createdRouterNatIp.SelfLink},
			SourceSubnetworkIpRangesToNat: pulumi.String("ALL_SUBNETWORKS_ALL_IP_RANGES"),
		}, pulumi.Parent(createdRouter))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create network router nat")
	}

	//export router nat name
	ctx.Export(outputs.RouterNatName, createdRouterNat.Name)

	createdSharedVpcIamResources := make([]pulumi.Resource, 0)

	if locals.GkeCluster.Spec.IsCreateSharedVpc {
		createdSharedVpcIamResources, err = sharedVpcIam(ctx,
			locals,
			createdClusterProject,
			createdNetworkProject,
			createdSubNetwork)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create shared vpc iam resources")
		}
	}

	clusterAutoscalingArgs := &container.ClusterClusterAutoscalingArgs{
		Enabled: pulumi.Bool(false),
	}

	//determine autoscaling input based on gke-cluster input spec
	if locals.GkeCluster.Spec.ClusterAutoscalingConfig != nil &&
		locals.GkeCluster.Spec.ClusterAutoscalingConfig.IsEnabled {
		clusterAutoscalingArgs = &container.ClusterClusterAutoscalingArgs{
			Enabled:            pulumi.Bool(true),
			AutoscalingProfile: pulumi.String("OPTIMIZE_UTILIZATION"),
			ResourceLimits: container.ClusterClusterAutoscalingResourceLimitArray{
				container.ClusterClusterAutoscalingResourceLimitArgs{
					ResourceType: pulumi.String("cpu"),
					Minimum:      pulumi.Int(locals.GkeCluster.Spec.ClusterAutoscalingConfig.CpuMinCores),
					Maximum:      pulumi.Int(locals.GkeCluster.Spec.ClusterAutoscalingConfig.CpuMaxCores),
				},
				container.ClusterClusterAutoscalingResourceLimitArgs{
					ResourceType: pulumi.String("memory"),
					Minimum:      pulumi.Int(locals.GkeCluster.Spec.ClusterAutoscalingConfig.MemoryMinGb),
					Maximum:      pulumi.Int(locals.GkeCluster.Spec.ClusterAutoscalingConfig.MemoryMaxGb),
				},
			},
		}
	}

	//create container cluster
	createdCluster, err := container.NewCluster(ctx,
		"cluster",
		&container.ClusterArgs{
			Name:                  pulumi.String(locals.GkeCluster.Metadata.Id),
			Project:               createdClusterProject.ProjectId,
			Location:              pulumi.String(locals.GkeCluster.Spec.Zone),
			Network:               createdNetwork.SelfLink,
			Subnetwork:            createdSubNetwork.SelfLink,
			RemoveDefaultNodePool: pulumi.Bool(true),
			DeletionProtection:    pulumi.Bool(false),
			WorkloadIdentityConfig: container.ClusterWorkloadIdentityConfigPtrInput(
				&container.ClusterWorkloadIdentityConfigArgs{
					WorkloadPool: pulumi.Sprintf("%s.svc.id.goog", createdClusterProject.ProjectId),
				}),
			//warning: cluster is not coming into ready state with value set to 0
			InitialNodeCount: pulumi.Int(1),
			ReleaseChannel: container.ClusterReleaseChannelPtrInput(
				&container.ClusterReleaseChannelArgs{
					Channel: pulumi.String(gkereleasechannel.GkeReleaseChannel_STABLE.String()),
				}),
			VerticalPodAutoscaling: container.ClusterVerticalPodAutoscalingPtrInput(
				&container.ClusterVerticalPodAutoscalingArgs{Enabled: pulumi.Bool(true)}),
			AddonsConfig: container.ClusterAddonsConfigPtrInput(&container.ClusterAddonsConfigArgs{
				HorizontalPodAutoscaling: container.ClusterAddonsConfigHorizontalPodAutoscalingPtrInput(
					&container.ClusterAddonsConfigHorizontalPodAutoscalingArgs{
						Disabled: pulumi.Bool(false)}),
				HttpLoadBalancing: container.ClusterAddonsConfigHttpLoadBalancingPtrInput(
					&container.ClusterAddonsConfigHttpLoadBalancingArgs{
						Disabled: pulumi.Bool(true)}),
				IstioConfig: container.ClusterAddonsConfigIstioConfigPtrInput(
					&container.ClusterAddonsConfigIstioConfigArgs{
						Disabled: pulumi.Bool(true)}),
				NetworkPolicyConfig: container.ClusterAddonsConfigNetworkPolicyConfigPtrInput(
					&container.ClusterAddonsConfigNetworkPolicyConfigArgs{
						Disabled: pulumi.Bool(true)}),
			}),
			PrivateClusterConfig: container.ClusterPrivateClusterConfigPtrInput(&container.ClusterPrivateClusterConfigArgs{
				EnablePrivateEndpoint: pulumi.Bool(false),
				EnablePrivateNodes:    pulumi.Bool(true),
				MasterIpv4CidrBlock:   pulumi.String(vars.ApiServerIpCidr),
			}),
			IpAllocationPolicy: container.ClusterIpAllocationPolicyPtrInput(
				// setting this is mandatory for shared vpc setup
				&container.ClusterIpAllocationPolicyArgs{
					ClusterSecondaryRangeName:  pulumi.String(locals.KubernetesPodSecondaryIpRangeName),
					ServicesSecondaryRangeName: pulumi.String(locals.KubernetesServiceSecondaryIpRangeName),
				}),
			MasterAuthorizedNetworksConfig: container.ClusterMasterAuthorizedNetworksConfigPtrInput(
				&container.ClusterMasterAuthorizedNetworksConfigArgs{
					CidrBlocks: container.ClusterMasterAuthorizedNetworksConfigCidrBlockArray{
						container.ClusterMasterAuthorizedNetworksConfigCidrBlockArgs{
							CidrBlock:   pulumi.String(vars.ClusterMasterAuthorizedNetworksCidrBlock),
							DisplayName: pulumi.String(vars.ClusterMasterAuthorizedNetworksCidrBlockDescription),
						}},
				}),
			ClusterAutoscaling: clusterAutoscalingArgs,
			//todo: disabling billing export temporarily
			//ResourceUsageExportConfig: container.ClusterResourceUsageExportConfigPtrInput(&container.ClusterResourceUsageExportConfigArgs{
			//	BigqueryDestination: container.ClusterResourceUsageExportConfigBigqueryDestinationArgs{
			//		DatasetId: pulumi.String(input.UsageMeteringDatasetId)},
			//	EnableNetworkEgressMetering:       pulumi.Bool(false),
			//	EnableResourceConsumptionMetering: pulumi.Bool(true),
			//}),
			LoggingConfig: container.ClusterLoggingConfigPtrInput(
				&container.ClusterLoggingConfigArgs{
					EnableComponents: pulumi.ToStringArray(locals.ContainerClusterLoggingComponentList),
				}),
		},
		pulumi.Parent(createdFolder),
		pulumi.DependsOn(createdSharedVpcIamResources))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add container cluster")
	}

	return createdCluster, nil
}
