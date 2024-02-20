package system

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/system/base"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/system/istiod"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/system/namespace"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/gcp/model"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AddedResources struct {
	IstioBaseHelmRelease *helm.Release
	IstiodHelmRelease    *helm.Release
}

type Input struct {
	Workspace          string
	KubernetesProvider *pulumikubernetes.Provider
	IstioAddonInput    *c2cv1deployk8cstackgcpmodel.AddonsIstio
}

func Resources(ctx *pulumi.Context, input *Input) (*AddedResources, error) {
	istioSystemNamespace, err := namespace.Resources(ctx, &namespace.Input{
		KubernetesProvider: input.KubernetesProvider,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add namespace")
	}
	istioBaseHelmRelease, err := base.Resources(ctx, &base.Input{
		IstioAddonBaseInput: input.IstioAddonInput.Base,
		Namespace:           istioSystemNamespace,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add base resources")
	}
	istiodHelmRelease, err := istiod.Resources(ctx, &istiod.Input{
		Workspace:             input.Workspace,
		IstioAddonDaemonInput: input.IstioAddonInput.Daemon,
		Namespace:             istioSystemNamespace,
		IstioBaseHelmRelease:  istioBaseHelmRelease,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add daemon resources")
	}
	return &AddedResources{
		IstioBaseHelmRelease: istioBaseHelmRelease,
		IstiodHelmRelease:    istiodHelmRelease,
	}, nil
}
