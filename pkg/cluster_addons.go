package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/addons"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/localz"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func clusterAddons(ctx *pulumi.Context, locals *localz.Locals,
	gcpProvider *gcp.Provider,
	kubernetesProvider *pulumikubernetes.Provider) error {
	if locals.GkeCluster.Spec.KubernetesAddons.IsInstallCertManager {
		if err := addons.CertManager(ctx,
			locals,
			kubernetesProvider); err != nil {
			return errors.Wrap(err, "failed to install cert-manager")
		}
	}

	return nil
}
