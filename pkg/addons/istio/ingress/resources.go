package ingress

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/controller"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/envoyfilter"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/namespace"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/system"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	WorkspaceDir              string
	KubernetesProvider        *pulumikubernetes.Provider
	IstioAddonIngressInput    *model.AddonsIstioIngress
	IstioSystemAddedResources *system.AddedResources
}

type AddedResources struct {
	AddedIngressControllerHelmRelease *helm.Release
}

func Resources(ctx *pulumi.Context, input *Input) (*AddedResources, error) {
	addedNamespace, err := namespace.Resources(ctx, &namespace.Input{
		KubernetesProvider: input.KubernetesProvider,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add namespace resources")
	}
	addedControllerResources, err := controller.Resources(ctx, &controller.Input{
		IstioAddonIngressInput:    input.IstioAddonIngressInput,
		Namespace:                 addedNamespace,
		IstioSystemAddedResources: input.IstioSystemAddedResources,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add controller resources")
	}
	if err := envoyfilter.Resources(ctx, &envoyfilter.Input{
		WorkspaceDir:                      input.WorkspaceDir,
		AddedIstioIngressNamespace:        addedNamespace,
		AddedIngressControllerHelmRelease: addedControllerResources.AddedIngressControllerHelmRelease,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to add envoy-filter")
	}
	return &AddedResources{
		AddedIngressControllerHelmRelease: addedControllerResources.AddedIngressControllerHelmRelease,
	}, nil
}
