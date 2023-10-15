package istio

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/system"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	Workspace          string
	KubernetesProvider *pulumikubernetes.Provider
	IstioAddonInput    *gcp.AddonsIstio
}

type AddedResources struct {
	AddedIngressControllerHelmRelease *helm.Release
}

func Resources(ctx *pulumi.Context, input *Input) (*AddedResources, error) {
	if input.IstioAddonInput == nil || !input.IstioAddonInput.Enabled {
		return nil, nil
	}
	istioSystemAddedResources, err := system.Resources(ctx, &system.Input{
		Workspace:          input.Workspace,
		KubernetesProvider: input.KubernetesProvider,
		IstioAddonInput:    input.IstioAddonInput,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add istio-system resources")
	}

	ingressAddedResources, err := ingress.Resources(ctx, &ingress.Input{
		KubernetesProvider:        input.KubernetesProvider,
		IstioAddonIngressInput:    input.IstioAddonInput.Ingress,
		IstioSystemAddedResources: istioSystemAddedResources,
		WorkspaceDir:              input.Workspace,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add istio-ingress resources")
	}
	return &AddedResources{
		AddedIngressControllerHelmRelease: ingressAddedResources.AddedIngressControllerHelmRelease,
	}, nil
}
