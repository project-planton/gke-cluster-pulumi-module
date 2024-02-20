package traefik

import (
	"github.com/pkg/errors"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/gcp/model"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
)

type Input struct {
	KubernetesProvider *pulumikubernetes.Provider
	TraefikAddonInput  *c2cv1deployk8cstackgcpmodel.AddonsTraefik
}

func Resources(ctx *pulumi.Context, input *Input) error {
	if input.TraefikAddonInput == nil || !input.TraefikAddonInput.Enabled {
		return nil
	}
	if err := addHelmRelease(ctx, input); err != nil {
		return errors.Wrap(err, "failed to add traefik helm release")
	}
	return nil
}

// addHelmRelease adds traefik helm release to the stack
func addHelmRelease(ctx *pulumi.Context, input *Input) error {
	helmChart := getHelmChart()
	helmVal := getHelmVal()
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
	}, pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}), pulumi.Provider(input.KubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s helm release", helmChart.ReleaseName)
	}
	return nil
}
