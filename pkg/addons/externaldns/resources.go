package externaldns

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/workloadidentity"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	pulk8scv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	v12 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	Namespace = "external-dns"
	Ksa       = "external-dns"
)

type Input struct {
	Workspace             string
	KubernetesProvider    *pulumikubernetes.Provider
	ExternalDnsAddonInput *model.AddonsExternalDns
	AddedExternalDnsGsa   *serviceaccount.Account
}

func Resources(ctx *pulumi.Context, input *Input) error {
	if input.ExternalDnsAddonInput == nil || !input.ExternalDnsAddonInput.Enabled {
		return nil
	}

	addedNamespace, err := addNamespace(ctx, input)
	if err != nil {
		return errors.Wrap(err, "failed to add namespace")
	}
	_, err = addServiceAccount(ctx, input, addedNamespace)
	if err != nil {
		return errors.Wrap(err, "failed to add service account")
	}
	return nil
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

// addServiceAccount adds service account to be used by the external-dns.
// reason for not configuring the helm release to manager service is because of the way pulumi output values are retrieved, it is not easy to inject the derived values into helm values.
// so, instead, disable service account creation in helm release and create it separately and add the google workload identity annotation to the service account which requires the email id of the google service account added as part of IAM module.
func addServiceAccount(ctx *pulumi.Context, input *Input, addedNamespace *pulk8scv1.Namespace) (*pulk8scv1.ServiceAccount, error) {
	addedServiceAccount, err := pulk8scv1.NewServiceAccount(ctx, Ksa, &pulk8scv1.ServiceAccountArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("ServiceAccount"),
		Metadata: v12.ObjectMetaPtrInput(&v12.ObjectMetaArgs{
			Name:        pulumi.String(Ksa),
			Namespace:   addedNamespace.Metadata.Name(),
			Annotations: pulumi.StringMap{workloadidentity.WorkloadIdentityKubeAnnotationKey: input.AddedExternalDnsGsa.Email},
		}),
	}, pulumi.Parent(addedNamespace))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add service account")
	}
	return addedServiceAccount, nil
}
