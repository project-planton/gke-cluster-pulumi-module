package addons

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/localz"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/outputs"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/vars"
	certmanagerv1 "github.com/plantoncloud/kubernetes-crd-pulumi-types/pkg/certmanager/certmanager/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

func CertManager(ctx *pulumi.Context, locals *localz.Locals,
	createdCluster *container.Cluster,
	gcpProvider *gcp.Provider,
	kubernetesProvider *pulumikubernetes.Provider) error {

	//create google service account required to create workload identity binding
	createdGoogleServiceAccount, err := serviceaccount.NewAccount(ctx,
		vars.CertManager.KsaName,
		&serviceaccount.AccountArgs{
			Project:     createdCluster.Project,
			Description: pulumi.String("cert-manager service account for solving dns challenges to issue certificates"),
			AccountId:   pulumi.String(vars.CertManager.KsaName),
			DisplayName: pulumi.String(vars.CertManager.KsaName),
		}, pulumi.Parent(createdCluster), pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to create cert-manager google service account")
	}

	//export cert-manager gsa email
	ctx.Export(outputs.CertManagerGsaEmail, createdGoogleServiceAccount.Email)

	//create workload-identity binding
	_, err = serviceaccount.NewIAMBinding(ctx,
		fmt.Sprintf("%s-workload-identity", vars.CertManager.KsaName),
		&serviceaccount.IAMBindingArgs{
			ServiceAccountId: createdGoogleServiceAccount.Name,
			Role:             pulumi.String(standard.Iam_workloadIdentityUser),
			Members: pulumi.StringArray{
				pulumi.Sprintf("serviceAccount:%s.svc.id.goog[%s/%s]",
					createdCluster.Project,
					vars.CertManager.Namespace,
					vars.CertManager.KsaName),
			},
		},
		pulumi.Parent(createdGoogleServiceAccount),
		pulumi.DependsOn([]pulumi.Resource{createdCluster}))
	if err != nil {
		return errors.Wrap(err, "failed to create workload-identity binding for cert-manager")
	}

	//create namespace resource
	createdNamespace, err := corev1.NewNamespace(ctx,
		vars.CertManager.Namespace,
		&corev1.NamespaceArgs{
			Metadata: metav1.ObjectMetaPtrInput(
				&metav1.ObjectMetaArgs{
					Name:   pulumi.String(vars.CertManager.Namespace),
					Labels: pulumi.ToStringMap(locals.KubernetesLabels),
				}),
		},
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create cert-manager namespace")
	}

	//create kubernetes service account to be used by the cert-manager.
	//it is not straight forward to add the gsa email as one of the helm values.
	// so, instead, disable service account creation in helm release and create it separately add
	// the Google workload identity annotation to the service account which requires the email id
	// of the Google service account added as part of IAM module.
	createdKubernetesServiceAccount, err := corev1.NewServiceAccount(ctx,
		vars.CertManager.KsaName,
		&corev1.ServiceAccountArgs{
			Metadata: metav1.ObjectMetaPtrInput(
				&metav1.ObjectMetaArgs{
					Name:      pulumi.String(vars.CertManager.KsaName),
					Namespace: createdNamespace.Metadata.Name(),
					Annotations: pulumi.StringMap{
						vars.WorkloadIdentityKubeAnnotationKey: createdGoogleServiceAccount.Email,
					},
				}),
		}, pulumi.Parent(createdNamespace))
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes service account")
	}

	//created helm-release
	_, err = helm.NewRelease(ctx, "cert-manager",
		&helm.ReleaseArgs{
			Name:            pulumi.String(vars.CertManager.HelmChartName),
			Namespace:       createdNamespace.Metadata.Name(),
			Chart:           pulumi.String(vars.CertManager.HelmChartName),
			Version:         pulumi.String(vars.CertManager.HelmChartVersion),
			CreateNamespace: pulumi.Bool(false),
			Atomic:          pulumi.Bool(false),
			CleanupOnFail:   pulumi.Bool(true),
			WaitForJobs:     pulumi.Bool(true),
			Timeout:         pulumi.Int(180),
			Values: pulumi.Map{
				"installCRDs": pulumi.Sprintf("%t", true),
				"extraArgs": pulumi.String(strings.Join([]string{
					"--dns01-recursive-nameservers=\"1.1.1.1:53\"",
					"--dns01-recursive-nameservers-only=true",
				}, ",")),
				"serviceAccount": pulumi.StringMap{
					"create": pulumi.Sprintf("%t", false),
					"name":   pulumi.String(vars.CertManager.KsaName),
				},
				"startupapicheck": pulumi.StringMap{
					"enabled": pulumi.Sprintf("%t", true),
					"timeout": pulumi.String("5m"),
				},
				"featureGates": pulumi.String("AdditionalCertificateOutputFormats=true"),
				"webhook": pulumi.StringMap{
					"extraArgs": pulumi.String(strings.Join([]string{
						"--feature-gates=AdditionalCertificateOutputFormats=true",
					}, ",")),
				},
			},
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String(vars.CertManager.HelmChartRepo),
			},
		}, pulumi.Parent(createdNamespace),
		pulumi.DependsOn([]pulumi.Resource{createdKubernetesServiceAccount}),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return errors.Wrap(err, "failed to create cert-manager helm release")
	}

	//create self-signed issuer
	_, err = certmanagerv1.NewClusterIssuer(ctx, "self-signed",
		&certmanagerv1.ClusterIssuerArgs{
			Metadata: metav1.ObjectMetaArgs{
				Name:   pulumi.String(vars.CertManager.SelfSignedIssuerName),
				Labels: pulumi.ToStringMap(locals.KubernetesLabels),
			},
			Spec: certmanagerv1.ClusterIssuerSpecArgs{
				SelfSigned: certmanagerv1.ClusterIssuerSpecSelfSignedArgs{},
			},
		})
	if err != nil {
		return errors.Wrap(err, "failed to create self-signed cluster-issuer")
	}
	return nil
}
