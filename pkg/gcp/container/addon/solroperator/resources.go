package solroperator

import (
	"github.com/pkg/errors"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/gcp/model"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	pulumiyaml "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"gopkg.in/yaml.v3"
)

type Input struct {
	KubernetesProvider     *pulumikubernetes.Provider
	SolrOperatorAddonInput *c2cv1deployk8cstackgcpmodel.AddonsSolrOperator
}

func Resources(ctx *pulumi.Context, input *Input) error {
	if input.SolrOperatorAddonInput == nil || !input.SolrOperatorAddonInput.Enabled {
		return nil
	}
	addedCrds, err := addCrdResources(ctx, input)
	if err != nil {
		return errors.Wrap(err, "failed to add solr-operator crds")
	}
	if err := addHelmRelease(ctx, input, addedCrds); err != nil {
		return errors.Wrap(err, "failed to add helm release")
	}
	return nil
}

func addCrdResources(ctx *pulumi.Context, input *Input) (*pulumiyaml.ConfigFile, error) {
	addedConfigFile, err := pulumiyaml.NewConfigFile(ctx, "solr-operator-crds",
		&pulumiyaml.ConfigFileArgs{
			File: getCrdsDownloadUrl(),
		}, pulumi.Provider(input.KubernetesProvider),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add solr-operator crds manifest")
	}
	return addedConfigFile, nil
}

func addHelmRelease(ctx *pulumi.Context, input *Input, addedCrds *pulumiyaml.ConfigFile) error {
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
		pulumi.DependsOn([]pulumi.Resource{addedCrds}),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s helm release", helmChart.ReleaseName)
	}
	return nil
}
