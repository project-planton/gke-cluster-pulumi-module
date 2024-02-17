package cluster

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster/cluster"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster/nodepool"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/network"
	code2cloudv1deployk8cmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/model"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp/model"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	KubeClusterId                string
	GcpZone                      string
	AddedContainerClusterProject *organizations.Project
	ContainerClusterInput        *c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackContainerClusterInput
	AddedNetworkResources        *network.AddedNetworkResources
	Labels                       map[string]string
	IsWorkloadLogsEnabled        bool
	NodePools                    []*code2cloudv1deployk8cmodel.KubeClusterNodePoolGcp
	ClusterAutoscalingConfig     *code2cloudv1deployk8cmodel.KubeClusterGcpClusterAutoscalingConfigSpec
}

type AddedContainerClusterResources struct {
	Cluster   *container.Cluster
	NodePools []pulumi.Resource
}

func Resources(ctx *pulumi.Context, input *Input) (*AddedContainerClusterResources, error) {
	addedCluster, err := cluster.Resources(ctx, &cluster.Input{
		KubeClusterId:                input.KubeClusterId,
		ClusterName:                  input.ContainerClusterInput.ClusterName,
		GcpZone:                      input.GcpZone,
		AddedNetworkResources:        input.AddedNetworkResources,
		AddedContainerClusterProject: input.AddedContainerClusterProject,
		ClusterConfig:                input.ContainerClusterInput.ContainerClusterConfig,
		IsWorkloadLogsEnabled:        input.IsWorkloadLogsEnabled,
		ClusterAutoscalingConfig:     input.ClusterAutoscalingConfig,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add kube-cluster resources")
	}
	addedNodePools, err := nodepool.Resources(ctx, &nodepool.Input{
		KubeClusterId: input.KubeClusterId,
		GcpZone:       input.GcpZone,
		Cluster:       addedCluster,
		Labels:        input.Labels,
		NodePools:     input.NodePools,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add cluster node-pool resources")
	}
	return &AddedContainerClusterResources{
		Cluster:   addedCluster,
		NodePools: addedNodePools,
	}, nil
}
