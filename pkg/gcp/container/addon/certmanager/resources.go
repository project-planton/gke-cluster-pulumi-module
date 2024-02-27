package certmanager

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/workloadidentity"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/certmanager/clusterissuer"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/gcp/model"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	pulk8scv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	v12 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
)

const (
	Namespace = "cert-manager"
	Ksa       = "cert-manager"
)

type Input struct {
	Workspace             string
	KubernetesProvider    *pulumikubernetes.Provider
	CertManagerAddonInput *c2cv1deployk8cstackgcpmodel.AddonsCertManager
	AddedCertManagerGsa   *serviceaccount.Account
}

func Resources(ctx *pulumi.Context, input *Input) error {
	if input.CertManagerAddonInput == nil || !input.CertManagerAddonInput.Enabled {
		return nil
	}
	addedNamespace, err := addNamespace(ctx, input)
	if err != nil {
		return errors.Wrap(err, "failed to add namespace")
	}
	addedServiceAccount, err := addServiceAccount(ctx, input, addedNamespace)
	if err != nil {
		return errors.Wrap(err, "failed to add service account")
	}
	helmRelease, err := addHelmRelease(ctx, addedNamespace, addedServiceAccount)
	if err != nil {
		return errors.Wrap(err, "failed to add helm release")
	}
	if err := clusterissuer.Resources(ctx, &clusterissuer.Input{
		Workspace:              input.Workspace,
		CertManagerAddonInput:  input.CertManagerAddonInput,
		CertManagerHelmRelease: helmRelease,
	}); err != nil {
		return errors.Wrap(err, "failed to add cluster issuer")
	}
	return nil
}

func addHelmRelease(ctx *pulumi.Context, addedNamespace *pulk8scv1.Namespace, addedServiceAccount *pulk8scv1.ServiceAccount) (*helm.Release, error) {
	helmVal := getHelmVal()
	helmChart := getHelmChart()
	var helmValInput map[string]interface{}
	helmValBytes, err := yaml.Marshal(helmVal)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal helm val to bytes")
	}
	if err := yaml.Unmarshal(helmValBytes, &helmValInput); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal helm val")
	}
	r, err := helm.NewRelease(ctx, helmChart.ReleaseName,
		&helm.ReleaseArgs{
			Name:            pulumi.String(helmChart.Name),
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
		}, pulumi.Parent(addedNamespace),
		pulumi.DependsOn([]pulumi.Resource{addedServiceAccount}),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s helm release", helmChart.ReleaseName)
	}
	return r, nil
}

func addNamespace(ctx *pulumi.Context, input *Input) (*pulk8scv1.Namespace, error) {
	ns, err := pulk8scv1.NewNamespace(ctx, Namespace, &pulk8scv1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("AddedNamespace"),
		Metadata: v12.ObjectMetaPtrInput(&v12.ObjectMetaArgs{
			Name: pulumi.String(Namespace),
		}),
	}, pulumi.Provider(input.KubernetesProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s namespace", Namespace)
	}
	return ns, nil
}

// addServiceAccount adds service account to be used by the cert-manager.
// reason for not configuring the helm release to manager service is because of the way pulumi output values are retrieved, it is not easy to inject the derived values into helm values.
// so, instead, disable service account creation in helm release and create it separately and add the google workload identity annotation to the service account which requires the email id of the google service account added as part of IAM module.
func addServiceAccount(ctx *pulumi.Context, input *Input, addedNamespace *pulk8scv1.Namespace) (*pulk8scv1.ServiceAccount, error) {
	addedServiceAccount, err := pulk8scv1.NewServiceAccount(ctx, Ksa, &pulk8scv1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: v12.ObjectMetaPtrInput(&v12.ObjectMetaArgs{
			Name:        pulumi.String(Ksa),
			Namespace:   addedNamespace.Metadata.Name(),
			Annotations: pulumi.StringMap{workloadidentity.WorkloadIdentityKubeAnnotationKey: input.AddedCertManagerGsa.Email},
		}),
	}, pulumi.Parent(addedNamespace))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add service account")
	}
	return addedServiceAccount, nil
}
