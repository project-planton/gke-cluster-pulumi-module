package ingress

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/ingress/gateway"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/ingress/service"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/network/ip"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	WorkspaceDir           string
	AddedIpAddresses       *ip.AddedIngressIpAddresses
	AddedContainerClusters *cluster.AddedContainerClusterResources
	AddedAddonResources    *addon.AddedResources
	IstioAddonInput        *gcp.AddonsIstio
}

func Resources(ctx *pulumi.Context, input *Input) error {
	if input.IstioAddonInput == nil || !input.IstioAddonInput.Enabled {
		return nil
	}

	//Resources adds kubernetes service resources that create internal and external load balancers
	//adds a kafka gateway that opens up port 9092 and does a tls passthrough
	if err := service.Resources(ctx, &service.Input{
		AddedComputeIpAddress:                  input.AddedIpAddresses,
		AddedIstioIngressControllerHelmRelease: input.AddedAddonResources.IstioAddedResources.AddedIngressControllerHelmRelease,
	}); err != nil {
		return errors.Wrapf(err, "failed to add service resources in container cluster")
	}
	if err := gateway.Resources(ctx, &gateway.Input{
		Workspace:                              input.WorkspaceDir,
		AddedIstioIngressControllerHelmRelease: input.AddedAddonResources.IstioAddedResources.AddedIngressControllerHelmRelease,
	}); err != nil {
		return errors.Wrapf(err, "failed to add gateway resources")
	}
	return nil
}
