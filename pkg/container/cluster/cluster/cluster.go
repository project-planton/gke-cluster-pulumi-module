package cluster

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network/subnet"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/enums/gkereleasechannel"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/pulumi/pulumicustomoutput"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	gkeReleaseChannel           = gkereleasechannel.GkeReleaseChannel_STABLE
	autoscalingProfileBalanced  = "BALANCED"
	autoscalingProfileOptimized = "OPTIMIZE_UTILIZATION"
)

type Input struct {
	KubeClusterId                string
	GcpZone                      string
	ClusterName                  string
	AddedContainerClusterProject *organizations.Project
	AddedNetworkResources        *network.AddedNetworkResources
	IsWorkloadLogsEnabled        bool
	ClusterConfig                *model.ClusterConfig
	ClusterAutoscalingConfig     *code2cloudv1deployk8cmodel.GkeClusterClusterAutoscalingConfigSpec
}

func Resources(ctx *pulumi.Context, input *Input) (*container.Cluster, error) {
	clusterName := input.KubeClusterId
	cc, err := container.NewCluster(ctx, clusterName, &container.ClusterArgs{
		Name:                  pulumi.String(clusterName),
		Project:               input.AddedContainerClusterProject.ProjectId,
		Location:              pulumi.String(input.GcpZone),
		Network:               input.AddedNetworkResources.AddedVpc.SelfLink,
		Subnetwork:            input.AddedNetworkResources.AddedSubnet.SelfLink,
		RemoveDefaultNodePool: pulumi.Bool(true),
		DeletionProtection:    pulumi.Bool(false),
		WorkloadIdentityConfig: container.ClusterWorkloadIdentityConfigPtrInput(
			&container.ClusterWorkloadIdentityConfigArgs{WorkloadPool: getWorkloadIdentityNamespace(input.AddedContainerClusterProject)}),
		InitialNodeCount: pulumi.Int(1),
		ReleaseChannel: container.ClusterReleaseChannelPtrInput(
			&container.ClusterReleaseChannelArgs{Channel: pulumi.String(gkeReleaseChannel.String())}),
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
			MasterIpv4CidrBlock:   pulumi.String(input.ClusterConfig.ApiServerIpCidr),
		}),
		IpAllocationPolicy: container.ClusterIpAllocationPolicyPtrInput(&container.ClusterIpAllocationPolicyArgs{
			ClusterSecondaryRangeName:  pulumi.String(subnet.GetPodsSecondaryRangeName(input.ClusterConfig.KubePodSecondaryRangeCidrSetNum)),     // required for shared vpc
			ServicesSecondaryRangeName: pulumi.String(subnet.GetServicesSecondaryRangeName(input.ClusterConfig.KubePodSecondaryRangeCidrSetNum)), // required for shared vpc
		}),
		MasterAuthorizedNetworksConfig: container.ClusterMasterAuthorizedNetworksConfigPtrInput(
			&container.ClusterMasterAuthorizedNetworksConfigArgs{
				CidrBlocks: container.ClusterMasterAuthorizedNetworksConfigCidrBlockArray{container.ClusterMasterAuthorizedNetworksConfigCidrBlockArgs{
					CidrBlock:   pulumi.String("0.0.0.0/0"),
					DisplayName: pulumi.String("all-for-testing"),
				}},
			}),
		ClusterAutoscaling: getClusterAutoScalingInput(input.ClusterAutoscalingConfig),
		//todo: disabling billing export temporarily
		//ResourceUsageExportConfig: container.ClusterResourceUsageExportConfigPtrInput(&container.ClusterResourceUsageExportConfigArgs{
		//	BigqueryDestination: container.ClusterResourceUsageExportConfigBigqueryDestinationArgs{
		//		DatasetId: pulumi.String(input.UsageMeteringDatasetId)},
		//	EnableNetworkEgressMetering:       pulumi.Bool(false),
		//	EnableResourceConsumptionMetering: pulumi.Bool(true),
		//}),
		LoggingConfig: container.ClusterLoggingConfigPtrInput(&container.ClusterLoggingConfigArgs{
			EnableComponents: getLoggingComponents(input.IsWorkloadLogsEnabled),
		}),
	}, pulumi.Parent(input.AddedContainerClusterProject), pulumi.DependsOn(input.AddedNetworkResources.AddedNetworkIamResources))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add container cluster")
	}
	ctx.Export(ClusterNameOutputName
	clusterName), cc.Name)
	ctx.Export(ClusterEndpointOutputName
	clusterName), cc.Endpoint)
	ctx.Export(ApiServerCidrBlockOutputName
	clusterName), cc.PrivateClusterConfig.MasterIpv4CidrBlock())
	ctx.Export(ClusterCaDataOutputName
	clusterName), cc.MasterAuth.ClusterCaCertificate().Elem())
	return cc, nil
}

func getWorkloadIdentityNamespace(addedGcpProject *organizations.Project) pulumi.StringOutput {
	return pulumi.Sprintf("%s.svc.id.goog", addedGcpProject.ProjectId)
}

func getLoggingComponents(isWorkloadLogsEnabled bool) pulumi.StringArray {
	comps := pulumi.StringArray{
		pulumi.String("SYSTEM_COMPONENTS"),
	}
	if isWorkloadLogsEnabled {
		comps = append(comps, pulumi.String("WORKLOADS"))
	}
	return comps
}

func getClusterAutoScalingInput(input *code2cloudv1deployk8cmodel.GkeClusterClusterAutoscalingConfigSpec) container.ClusterClusterAutoscalingPtrInput {
	if input == nil || !input.IsEnabled {
		return container.ClusterClusterAutoscalingPtrInput(&container.ClusterClusterAutoscalingArgs{
			Enabled: pulumi.Bool(false),
		})
	}
	return container.ClusterClusterAutoscalingPtrInput(&container.ClusterClusterAutoscalingArgs{
		Enabled:            pulumi.Bool(input.IsEnabled),
		AutoscalingProfile: pulumi.String(autoscalingProfileOptimized),
		ResourceLimits: container.ClusterClusterAutoscalingResourceLimitArray{
			container.ClusterClusterAutoscalingResourceLimitArgs{
				ResourceType: pulumi.String("cpu"),
				Minimum:      pulumi.Int(input.CpuMinCores),
				Maximum:      pulumi.Int(input.CpuMaxCores),
			},
			container.ClusterClusterAutoscalingResourceLimitArgs{
				ResourceType: pulumi.String("memory"),
				Minimum:      pulumi.Int(input.MemoryMinGb),
				Maximum:      pulumi.Int(input.MemoryMaxGb),
			},
		},
	})
}

func GetApiServerCidrBlockOutputName     clusterFullName string) string {
	return pulumicustomoutput.Name(clusterFullName, "api-server-ip-cidr")
}

func GetClusterNameOutputName     clusterFullName string) string {
	return pulumigoogleprovider.PulumiOutputName
	container.Cluster{}, clusterFullName, englishword.EnglishWord_name.String())
}

func GetClusterEndpointOutputName     clusterFullName string) string {
	return pulumigoogleprovider.PulumiOutputName
	container.Cluster{}, clusterFullName, englishword.EnglishWord_endpoint.String())
}

func GetClusterCaDataOutputName     clusterFullName string) string {
	return pulumigoogleprovider.PulumiOutputName
	container.Cluster{}, clusterFullName, "ca-data")
}
