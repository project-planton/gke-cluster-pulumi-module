package addons

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/addons/certmanager"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/localz"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/vars"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func certManager(ctx *pulumi.Context, locals *localz.Locals,
	kubernetesProvider *pulumikubernetes.Provider) error {

	createdGsa, err := serviceaccount.NewAccount(ctx,
		certmanager.Ksa,
		&serviceaccount.AccountArgs{
		Project:     addedContainerClusterProject.ProjectId,
		Description: pulumi.String("cert-manager service account for solving dns challenges to issue certificates"),
		AccountId:   pulumi.String(certmanager.Ksa),
		DisplayName: pulumi.String(certmanager.Ksa),
	}, pulumi.Parent(addedContainerClusterProject))
	if err != nil {
		return nil, errors.Wrapf(err, "failed add new %s svc acct", certmanager.Ksa)
	}
	ctx.Export(GsaEmailOutputName), createdGsa.Email)

	//create namespace resource
	createdNamespace, err := corev1.NewNamespace(ctx,
		vars.CertManagerNamespace,
		&corev1.NamespaceArgs{
			Metadata: metav1.ObjectMetaPtrInput(
				&metav1.ObjectMetaArgs{
					Name:   pulumi.String(vars.CertManagerNamespace),
					Labels: pulumi.ToStringMap(locals.KubernetesLabels),
				}),
		},
		pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "5s", Update: "5s", Delete: "5s"}),
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create %s namespace", locals.Namespace)
	}

	createdCertManagerHelmRelease, err := helm.NewRelease(ctx, "cert-manager",
		&helm.ReleaseArgs{
			Name:            pulumi.String("cert-manager"),
			Namespace:       pulumi.String(Namespace),
			Chart:           pulumi.String(helmChart.Name),
			Version:         pulumi.String(helmChart.Version),
			CreateNamespace: pulumi.Bool(false),
			Atomic:          pulumi.Bool(false),
			CleanupOnFail:   pulumi.Bool(true),
			WaitForJobs:     pulumi.Bool(true),
			Timeout:         pulumi.Int(180), // 3 minutes
			Values:          pulumi.ToMap(helmValInput),
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String(helmChart.Repo),
			},
		}, pulumi.Parent(createdNamespace),
		pulumi.DependsOn([]pulumi.Resource{addedServiceAccount}),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return
		errors.Wrapf(err, "failed to add %s helm release", helmChart.ReleaseName)
	}
	return nil
}
