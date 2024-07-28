package addons

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/localz"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/gcp/pulumigkekubernetesprovider"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Addons(ctx *pulumi.Context,
	locals *localz.Locals,
	gcpProvider *gcp.Provider,
	createdContainerCluster *container.Cluster,
	createdNodePools []*pulumi.Resource) error {

	if locals.GkeCluster.Spec.KubernetesAddons == nil {
		return nil
	}

	pulumigkekubernetesprovider.GetWithCreatedGkeClusterAndCreatedGsaKey(ctx, createdContainerCluster)
	if locals.GkeCluster.Spec.KubernetesAddons.IsInstallCertManager {
		if err := certManager(ctx, locals, kubernetesProvider, createdContainerCluster, createdNodePools); err != nil {
			return errors.Wrap(err, "failed to install cert-manager")
		}
	}

	return nil
}
