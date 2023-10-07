package gcp

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/addon"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/cluster"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/ingress"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/iam"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/network"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/projects"
	pulumigcpprovider "github.com/plantoncloud-inc/pulumi-stack-runner-sdk/go/pulumi/automation/provider/google"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// resources sets up kube-cluster projects
// * creates kube-cluster and vpc-network projects
// * enabled required apis on kube-cluster and vpc-network projects
func (s *Stack) resources(ctx *pulumi.Context) error {
	gcpProvider, err := pulumigcpprovider.Get(ctx, s.Input.CredentialsInput.Google)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}
	addedProjectsResources, err := projects.Resources(ctx, &projects.Input{
		CloudAccountFolderId: s.Input.ResourceInput.CloudAccount.Status.Gcp.CloudAccountFolder.Id,
		KubeClusterId:        s.Input.ResourceInput.KubeCluster.Metadata.Id,
		GcpRegion:            s.Input.ResourceInput.KubeCluster.Spec.Gcp.Region,
		GcpZone:              s.Input.ResourceInput.KubeCluster.Spec.Gcp.Zone,
		GcpProvider:          gcpProvider,
		BillingAccountId:     s.Input.ResourceInput.KubeCluster.Spec.Gcp.BillingAccountId,
		IsCreateSharedVpc:    s.Input.ResourceInput.KubeCluster.Spec.Gcp.IsCreateSharedVpc,
		Labels:               s.GcpLabels,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add projects resources")
	}

	addedNetworkResources, err := network.Resources(ctx, &network.Input{
		KubeClusterId:          s.Input.ResourceInput.KubeCluster.Metadata.Id,
		GcpRegion:              s.Input.ResourceInput.KubeCluster.Spec.Gcp.Region,
		IsCreateSharedVpc:      s.Input.ResourceInput.KubeCluster.Spec.Gcp.IsCreateSharedVpc,
		AddedProjectsResources: addedProjectsResources,
		Labels:                 s.GcpLabels,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add network resources")
	}

	addedContainerClusters, err := cluster.Resources(ctx, &cluster.Input{
		KubeClusterId:                s.Input.ResourceInput.KubeCluster.Metadata.Id,
		GcpZone:                      s.Input.ResourceInput.KubeCluster.Spec.Gcp.Zone,
		AddedContainerClusterProject: addedProjectsResources.KubeClusterProjects.ContainerClusterProject,
		ContainerClusterInput:        s.Input.ResourceInput.Container.Cluster,
		AddedNetworkResources:        addedNetworkResources,
		Labels:                       s.GcpLabels,
		IsWorkloadLogsEnabled:        s.Input.ResourceInput.KubeCluster.Spec.Gcp.IsWorkloadLogsEnabled,
		NodePools:                    s.Input.ResourceInput.KubeCluster.Spec.Gcp.NodePools,
		ClusterAutoscalingConfig:     s.Input.ResourceInput.KubeCluster.Spec.Gcp.ClusterAutoscalingConfig,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add network resources")
	}

	addedIamResources, err := iam.Resources(ctx, &iam.Input{
		AddedContainerClusterProject: addedProjectsResources.KubeClusterProjects.ContainerClusterProject,
		AddedContainerClusters:       addedContainerClusters,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add iam resources")
	}

	addedContainerClusterAddonResources, err := addon.Resources(ctx, &addon.Input{
		KubeClusterAddons:              s.Input.ResourceInput.KubeCluster.Spec.KubernetesAddons,
		ContainerAddonInput:            s.Input.ResourceInput.Container.Addon,
		ReqWorkspace:                   s.ReqWorkspace,
		AddedContainerClusterProject:   addedProjectsResources.KubeClusterProjects.ContainerClusterProject,
		AddedContainerClusterResources: addedContainerClusters,
		AddedIamResources:              addedIamResources,
	})
	if err != nil {
		return errors.Wrap(err, "failed to add container addon resources")
	}

	if err := ingress.Resources(ctx, &ingress.Input{
		ReqWorkspace:           s.ReqWorkspace,
		AddedIpAddresses:       addedNetworkResources.AddedIpAddresses,
		AddedContainerClusters: addedContainerClusters,
		AddedAddonResources:    addedContainerClusterAddonResources,
	}); err != nil {
		return errors.Wrap(err, "failed to add container ingress resources")
	}
	return nil
}
