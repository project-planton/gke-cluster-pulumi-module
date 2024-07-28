package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/addons"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/localz"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func clusterAddons(ctx *pulumi.Context, locals *localz.Locals,
	createdCluster *container.Cluster, gcpProvider *gcp.Provider,
	kubernetesProvider *pulumikubernetes.Provider) error {

	if locals.GkeCluster.Spec.KubernetesAddons.IsInstallIngressNginx {
		if err := addons.Istio(ctx); err != nil {
			return errors.Wrap(err, "failed to install ingress-nginx resources")
		}
	}

	if locals.GkeCluster.Spec.KubernetesAddons.IsInstallIstio {
		if err := addons.Istio(ctx); err != nil {
			return errors.Wrap(err, "failed to install istio resources")
		}
	}

	if locals.GkeCluster.Spec.KubernetesAddons.IsInstallCertManager {
		if err := addons.CertManager(ctx,
			locals,
			createdCluster,
			gcpProvider,
			kubernetesProvider); err != nil {
			return errors.Wrap(err, "failed to install cert-manager resources")
		}
	}

	if locals.GkeCluster.Spec.KubernetesAddons.IsInstallExternalSecrets {
		if err := addons.ExternalSecrets(ctx,
			locals,
			createdCluster,
			gcpProvider,
			kubernetesProvider); err != nil {
			return errors.Wrap(err, "failed to install external-secrets resources")
		}
	}

	if locals.GkeCluster.Spec.KubernetesAddons.IsInstallExternalDns {
		if err := addons.ExternalDns(ctx); err != nil {
			return errors.Wrap(err, "failed to install external-dns resources")
		}
	}

	if locals.GkeCluster.Spec.KubernetesAddons.IsInstallPostgresOperator {
		if err := addons.ZalandoOperator(ctx); err != nil {
			return errors.Wrap(err, "failed to install zalando postgres-operator resources")
		}
	}

	if locals.GkeCluster.Spec.KubernetesAddons.IsInstallSolrOperator {
		if err := addons.SolrOperator(ctx); err != nil {
			return errors.Wrap(err, "failed to install solr-operator resources")
		}
	}

	if locals.GkeCluster.Spec.KubernetesAddons.IsInstallKafkaOperator {
		if err := addons.StrimziOperator(ctx); err != nil {
			return errors.Wrap(err, "failed to install strimzi kafka-operator resources")
		}
	}

	return nil
}
