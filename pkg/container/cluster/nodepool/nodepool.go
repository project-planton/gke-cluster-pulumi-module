package nodepool

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster/nodepool/tag"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	autoRepair  = true
	autoUpgrade = true
)

type Input struct {
	KubeClusterId string
	GcpZone       string
	Cluster       *container.Cluster
	NodePools     []*code2cloudv1deployk8cmodel.KubeClusterNodePoolGcp
	Labels        map[string]string
}

func Resources(ctx *pulumi.Context, input *Input) ([]pulumi.Resource, error) {
	addedNodePools := make([]pulumi.Resource, 0)
	for _, np := range input.NodePools {
		addedNodePool, err := addNodePool(ctx, input, np)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to add %s node-pool", np.Name)
		}
		addedNodePools = append(addedNodePools, addedNodePool)
	}
	return addedNodePools, nil
}

func addNodePool(ctx *pulumi.Context, input *Input, clusterNodePoolInput *code2cloudv1deployk8cmodel.KubeClusterNodePoolGcp) (*container.NodePool, error) {
	addedNodePool, err := container.NewNodePool(ctx, clusterNodePoolInput.Name, &container.NodePoolArgs{
		Location:  pulumi.String(input.GcpZone),
		Project:   input.Cluster.Project,
		Cluster:   input.Cluster.Name,
		NodeCount: pulumi.Int(clusterNodePoolInput.MinNodeCount),
		Autoscaling: container.NodePoolAutoscalingPtrInput(&container.NodePoolAutoscalingArgs{
			MinNodeCount: pulumi.Int(clusterNodePoolInput.MinNodeCount),
			MaxNodeCount: pulumi.Int(clusterNodePoolInput.MaxNodeCount),
		}),
		Management: container.NodePoolManagementPtrInput(&container.NodePoolManagementArgs{
			AutoRepair:  pulumi.Bool(autoRepair),
			AutoUpgrade: pulumi.Bool(autoUpgrade),
		}),
		NodeConfig: &container.NodePoolNodeConfigArgs{
			Labels:      pulumi.ToStringMap(input.Labels),
			MachineType: pulumi.String(clusterNodePoolInput.MachineType),
			Metadata:    pulumi.StringMap{"disable-legacy-endpoints": pulumi.String("true")},
			OauthScopes: getOauthScopes(),
			Preemptible: pulumi.Bool(clusterNodePoolInput.IsSpotEnabled),
			Tags:        pulumi.StringArray{pulumi.String(tag.Get(input.KubeClusterId))},
			WorkloadMetadataConfig: container.NodePoolNodeConfigWorkloadMetadataConfigPtrInput(
				&container.NodePoolNodeConfigWorkloadMetadataConfigArgs{
					Mode: pulumi.String("GKE_METADATA")}),
		},
		UpgradeSettings: container.NodePoolUpgradeSettingsPtrInput(&container.NodePoolUpgradeSettingsArgs{
			MaxSurge:       pulumi.Int(2),
			MaxUnavailable: pulumi.Int(1),
		}),
	},
		pulumi.Parent(input.Cluster),
		pulumi.IgnoreChanges([]string{"nodeCount"}),
		pulumi.DeleteBeforeReplace(true),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add node pool")
	}
	ctx.Export(NodePoolNameOutputName
	clusterNodePoolInput.Name), addedNodePool.Name)
	ctx.Export(NodePoolMachineTypeOutputName
	clusterNodePoolInput.Name), addedNodePool.NodeConfig.MachineType())
	ctx.Export(NodePoolIsSpotInstancesOutputName
	clusterNodePoolInput.Name), addedNodePool.NodeConfig.Preemptible())
	return addedNodePool, nil
}
func getOauthScopes() pulumi.StringArrayInput {
	scopes := pulumi.StringArray{
		pulumi.String("https://www.googleapis.com/auth/monitoring"),
		pulumi.String("https://www.googleapis.com/auth/monitoring.write"),
		pulumi.String("https://www.googleapis.com/auth/devstorage.read_only"),
		pulumi.String("https://www.googleapis.com/auth/logging.write"),
	}
	return scopes
}

func GetNodePoolNameOutputName     nodePoolName string) string {
	return pulumigoogleprovider.PulumiOutputName
	container.NodePool{}, nodePoolName,
		englishword.EnglishWord_name.String())
}

func GetNodePoolMachineTypeOutputName     nodePoolName string) string {
	return pulumigoogleprovider.PulumiOutputName
	container.NodePool{}, nodePoolName,
		englishword.EnglishWord_machine.String(), englishword.EnglishWord_type.String())
}

func GetNodePoolIsSpotInstancesOutputName     nodePoolName string) string {
	return pulumigoogleprovider.PulumiOutputName
	container.NodePool{}, nodePoolName,
		englishword.EnglishWord_spot.String(), englishword.EnglishWord_instances.String())
}
