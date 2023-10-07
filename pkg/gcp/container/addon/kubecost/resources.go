package kubecost

import (
	"buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	"github.com/pkg/errors"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
)

type Input struct {
	KubernetesProvider *pulumikubernetes.Provider
	KubeCostAddonInput *gcp.AddonsKubeCost
}

func Resources(ctx *pulumi.Context, input *Input) error {
	if input.KubeCostAddonInput == nil || !input.KubeCostAddonInput.Enabled {
		return nil
	}
	if err := addHelmRelease(ctx, input); err != nil {
		return errors.Wrap(err, "failed to add helm release")
	}
	return nil
}

func addHelmRelease(ctx *pulumi.Context, input *Input) error {
	helmVal := getHelmVal()
	helmChart := getHelmChart()
	var helmValInput map[string]interface{}
	helmValBytes, err := yaml.Marshal(helmVal)
	if err != nil {
		return errors.Wrap(err, "failed to marshal helm val to bytes")
	}
	if err := yaml.Unmarshal(helmValBytes, &helmValInput); err != nil {
		return errors.Wrap(err, "failed to unmarshal helm val")
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
		Atomic:          pulumi.Bool(false),
		CleanupOnFail:   pulumi.Bool(true),
		WaitForJobs:     pulumi.Bool(true),
		Timeout:         pulumi.Int(180), // 3 minutes
		Values:          pulumi.ToMap(helmValInput),
		RepositoryOpts: helm.RepositoryOptsArgs{
			Repo: pulumi.String(helmChart.Repo),
		},
	}, pulumi.Provider(input.KubernetesProvider),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s helm release", helmChart.ReleaseName)
	}
	return nil
}
