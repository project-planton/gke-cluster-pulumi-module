package istio

import (
	"buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/addon/istio/ingress"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/addon/istio/system"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
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
		ReqWorkspace:              input.Workspace,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add istio-ingress resources")
	}
	return &AddedResources{
		AddedIngressControllerHelmRelease: ingressAddedResources.AddedIngressControllerHelmRelease,
	}, nil
}
