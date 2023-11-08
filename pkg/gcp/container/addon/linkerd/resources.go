package linkerd

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	pk8scv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	pk8smv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
)

type Input struct {
	KubernetesProvider *pulumikubernetes.Provider
	LinkerdAddonInput  *gcp.AddonsLinkerd
}

func Resources(ctx *pulumi.Context, input *Input) error {
	if input.LinkerdAddonInput == nil || !input.LinkerdAddonInput.Enabled {
		return nil
	}
	ns, err := addNamespace(ctx, input)
	if err != nil {
		return errors.Wrap(err, "failed to add namespace")
	}
	if err := addHelmRelease(ctx, input, ns); err != nil {
		return errors.Wrap(err, "failed to add helm release")
	}
	return nil
}

func addNamespace(ctx *pulumi.Context, input *Input) (*pk8scv1.Namespace, error) {
	ns, err := pk8scv1.NewNamespace(ctx, Namespace, &pk8scv1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("AddedNamespace"),
		Metadata: pk8smv1.ObjectMetaArgs{
			Labels: pulumi.ToStringMap(NamespaceLabels),
			Name:   pulumi.String(Namespace),
		},
	}, pulumi.Provider(input.KubernetesProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add namespace")
	}
	return ns, nil
}

func addHelmRelease(ctx *pulumi.Context, input *Input, ns *pk8scv1.Namespace) error {
	helmVal := GetHelmVal()
	helmChart := GetHelmChart()
	var helmValInput map[string]interface{}
	helmValBytes, err := yaml.Marshal(helmVal)
	if err != nil {
		return errors.Wrap(err, "failed to marshal helm val to bytes")
	}
	if err := yaml.Unmarshal(helmValBytes, &helmValInput); err != nil {
		return errors.Wrap(err, "failed to unmarshal helm val")
	}
	_, err = helm.NewRelease(ctx, helmChart.ReleaseName, &helm.ReleaseArgs{
		Name:            pulumi.String(helmChart.ReleaseName),
		Namespace:       pulumi.String(Namespace),
		Chart:           pulumi.String(helmChart.Name),
		Version:         pulumi.String(helmChart.Version),
		CreateNamespace: pulumi.Bool(true),
		Atomic:          pulumi.Bool(true),
		CleanupOnFail:   pulumi.Bool(true),
		WaitForJobs:     pulumi.Bool(true),
		Timeout:         pulumi.Int(180), // 3 minutes
		Values:          pulumi.ToMap(helmValInput),
		RepositoryOpts: helm.RepositoryOptsArgs{
			Repo: pulumi.String(helmChart.Repo),
		},
	}, pulumi.Parent(ns),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}), pulumi.DependsOn([]pulumi.Resource{ns}))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s helm release", helmChart.ReleaseName)
	}
	return nil
}
