package ingress

import (
	"buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/addon/istio/ingress/controller"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/addon/istio/ingress/envoyfilter"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/addon/istio/ingress/namespace"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/addon/istio/system"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	ReqWorkspace              string
	KubernetesProvider        *pulumikubernetes.Provider
	IstioAddonIngressInput    *gcp.AddonsIstioIngress
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
		ReqWorkspace:                      input.ReqWorkspace,
		AddedIstioIngressNamespace:        addedNamespace,
		AddedIngressControllerHelmRelease: addedControllerResources.AddedIngressControllerHelmRelease,
	}); err != nil {
		return nil, errors.Wrap(err, "failed to add envoy-filter")
	}
	return &AddedResources{
		AddedIngressControllerHelmRelease: addedControllerResources.AddedIngressControllerHelmRelease,
	}, nil
}
