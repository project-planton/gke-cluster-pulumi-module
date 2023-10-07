package clustersecretstore

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	Name = "gcp-backend"
	Kind = "ClusterSecretStore"
	//clusterSecretStoreManifestTemplate requires following inputs in the same order
	// 1. cluster secret store kind
	// 3. name of the cluster secret store
	// 3. gcp project id
	clusterSecretStoreManifestTemplate = `
apiVersion: external-secrets.io/v1beta1
kind: %s
metadata:
  name: %s
spec:
  provider:
    gcpsm:
      projectID: %s
`
)

type Input struct {
	AddedContainerClusterProject *organizations.Project
	ExternalSecretsHelmRelease   *helm.Release
}

func Resources(ctx *pulumi.Context, input *Input) error {
	if err := addClusterSecretStore(ctx, input); err != nil {
		return errors.Wrap(err, "failed to add cluster secret store")
	}
	return nil
}

func addClusterSecretStore(ctx *pulumi.Context, input *Input) error {
	input.AddedContainerClusterProject.ProjectId.ApplyT(func(shareProjectId string) error {
		_, err := pulumik8syaml.NewConfigGroup(ctx, "addon-external-secrets-cluster-secret-store",
			&pulumik8syaml.ConfigGroupArgs{
				YAML: []string{
					fmt.Sprintf(`
apiVersion: external-secrets.io/v1beta1
kind: %s
metadata:
  name: %s
spec:
  provider:
    gcpsm:
      projectID: %s
        `, Kind, Name, shareProjectId),
				},
			}, pulumi.Parent(input.ExternalSecretsHelmRelease), pulumi.DependsOn([]pulumi.Resource{input.ExternalSecretsHelmRelease}))
		if err != nil {
			return errors.Wrap(err, "failed to add cluster secret store manifest")
		}
		return nil
	})
	return nil
}
