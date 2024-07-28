package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/localz"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func clusterNodePools(ctx *pulumi.Context,
	locals *localz.Locals,
	createdCluster *container.Cluster) ([]pulumi.Resource, error) {
	createdNodePoolResources := make([]pulumi.Resource, 0)

	for _, nodePoolSpec := range locals.GkeCluster.Spec.NodePools {
		createdNodePool, err := container.NewNodePool(ctx, nodePoolSpec.Name, &container.NodePoolArgs{
			Location:  pulumi.String(locals.GkeCluster.Spec.Zone),
			Project:   createdCluster.Project,
			Cluster:   createdCluster.Name,
			NodeCount: pulumi.Int(nodePoolSpec.MinNodeCount),
			Autoscaling: container.NodePoolAutoscalingPtrInput(&container.NodePoolAutoscalingArgs{
				MinNodeCount: pulumi.Int(nodePoolSpec.MinNodeCount),
				MaxNodeCount: pulumi.Int(nodePoolSpec.MaxNodeCount),
			}),
			Management: container.NodePoolManagementPtrInput(&container.NodePoolManagementArgs{
				AutoRepair:  pulumi.Bool(true),
				AutoUpgrade: pulumi.Bool(true),
			}),
			NodeConfig: &container.NodePoolNodeConfigArgs{
				Labels:      pulumi.ToStringMap(locals.GcpLabels),
				MachineType: pulumi.String(nodePoolSpec.MachineType),
				Metadata:    pulumi.StringMap{"disable-legacy-endpoints": pulumi.String("true")},
				OauthScopes: pulumi.StringArray{
					pulumi.String("https://www.googleapis.com/auth/monitoring"),
					pulumi.String("https://www.googleapis.com/auth/monitoring.write"),
					pulumi.String("https://www.googleapis.com/auth/devstorage.read_only"),
					pulumi.String("https://www.googleapis.com/auth/logging.write"),
				},
				Preemptible: pulumi.Bool(nodePoolSpec.IsSpotEnabled),
				Tags: pulumi.StringArray{
					pulumi.String(locals.NetworkTag),
				},
				WorkloadMetadataConfig: container.NodePoolNodeConfigWorkloadMetadataConfigPtrInput(
					&container.NodePoolNodeConfigWorkloadMetadataConfigArgs{
						Mode: pulumi.String("GKE_METADATA")}),
			},
			UpgradeSettings: container.NodePoolUpgradeSettingsPtrInput(&container.NodePoolUpgradeSettingsArgs{
				MaxSurge:       pulumi.Int(2),
				MaxUnavailable: pulumi.Int(1),
			}),
		},
			pulumi.Parent(createdCluster),
			pulumi.IgnoreChanges([]string{"nodeCount"}),
			pulumi.DeleteBeforeReplace(true),
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create node-pool")
		}

		createdNodePoolResources = append(createdNodePoolResources, createdNodePool)
	}

	return createdNodePoolResources, nil
}
