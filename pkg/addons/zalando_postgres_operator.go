package addons

import (
	"github.com/pkg/errors"
	kuberneteslabels "github.com/plantoncloud-inc/go-commons/kubernetes/labels"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/localz"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/vars"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ZalandoPostgresOperator(ctx *pulumi.Context, locals *localz.Locals,
	kubernetesProvider *pulumikubernetes.Provider) error {
	//create namespace resource
	createdNamespace, err := corev1.NewNamespace(ctx,
		vars.ExternalSecrets.Namespace,
		&corev1.NamespaceArgs{
			Metadata: metav1.ObjectMetaPtrInput(
				&metav1.ObjectMetaArgs{
					Name:   pulumi.String(vars.ZalandoPostgresOperator.Namespace),
					Labels: pulumi.ToStringMap(locals.KubernetesLabels),
				}),
		},
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create namespace")
	}

	//create helm-release
	_, err = helm.NewRelease(ctx, "zalando-postgres-operator",
		&helm.ReleaseArgs{
			Name:            pulumi.String(vars.ZalandoPostgresOperator.HelmChartName),
			Namespace:       pulumi.String(vars.ZalandoPostgresOperator.Namespace),
			Chart:           pulumi.String(vars.ZalandoPostgresOperator.HelmChartName),
			Version:         pulumi.String(vars.ZalandoPostgresOperator.HelmChartVersion),
			CreateNamespace: pulumi.Bool(false),
			Atomic:          pulumi.Bool(false),
			CleanupOnFail:   pulumi.Bool(true),
			WaitForJobs:     pulumi.Bool(true),
			Timeout:         pulumi.Int(180),
			Values: pulumi.Map{
				"configKubernetes": pulumi.Map{
					"inherited_labels": pulumi.ToStringArray(
						[]string{
							kuberneteslabels.Resource,
							kuberneteslabels.Organization,
							kuberneteslabels.Environment,
							kuberneteslabels.ResourceKind,
							kuberneteslabels.ResourceId,
						},
					),
				},
			},
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String(vars.ZalandoPostgresOperator.HelmChartRepo),
			},
		}, pulumi.Parent(createdNamespace),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return errors.Wrap(err, "failed to create helm release")
	}
	return nil
}
