package istiod

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
)

type Input struct {
	Workspace             string
	IstioAddonDaemonInput *model.AddonsIstioDaemon
	Namespace             *v1.Namespace
	IstioBaseHelmRelease  *helm.Release
}

func Resources(ctx *pulumi.Context, input *Input) (*helm.Release, error) {
	release, err := addHelmRelease(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add helm release")
	}
	return release, nil
}

func addHelmRelease(ctx *pulumi.Context, input *Input) (*helm.Release, error) {
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
	release, err := helm.NewRelease(ctx, helmChart.ReleaseName, &helm.ReleaseArgs{
		Name:            pulumi.String(helmChart.ReleaseName),
		Namespace:       input.Namespace.Metadata.Name(),
		Chart:           pulumi.String(helmChart.Name),
		Version:         pulumi.String(helmChart.Version),
		CreateNamespace: pulumi.Bool(true),
		Atomic:          pulumi.Bool(false),
		CleanupOnFail:   pulumi.Bool(true),
		WaitForJobs:     pulumi.Bool(true),
		Timeout:         pulumi.Int(300), // 5 minutes
		Values:          pulumi.ToMap(helmValInput),
		RepositoryOpts: helm.RepositoryOptsArgs{
			Repo: pulumi.String(helmChart.Repo),
		},
	}, pulumi.Parent(input.Namespace),
		pulumi.DependsOn([]pulumi.Resource{input.IstioBaseHelmRelease}), pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s helm release", helmChart.ReleaseName)
	}
	return release, nil
}