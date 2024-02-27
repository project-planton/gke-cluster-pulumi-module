package istio

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/system"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/gcp/model"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	Workspace          string
	KubernetesProvider *pulumikubernetes.Provider
	IstioAddonInput    *c2cv1deployk8cstackgcpmodel.AddonsIstio
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
