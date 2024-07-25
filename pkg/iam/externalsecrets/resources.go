package externalsecrets

import (
	"fmt"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/externalsecrets"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	AddedContainerClusterProject   *organizations.Project
	AddedContainerClusterResources *cluster.AddedContainerClusterResources
}

func Resources(ctx *pulumi.Context, input *Input) (*serviceaccount.Account, error) {
	gsa, err := addGsa(ctx, input.AddedContainerClusterProject)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add gsa")
	}
	if err := addIamRole(ctx, input.AddedContainerClusterProject, gsa); err != nil {
		return nil, errors.Wrap(err, "failed to add iam role")
	}
	if err := addWorkloadIdentityBinding(ctx, input, gsa); err != nil {
		return nil, errors.Wrap(err, "failed to add workload identity binding")
	}
	return gsa, nil
}

func addGsa(ctx *pulumi.Context, addedContainerClusterProject *organizations.Project) (*serviceaccount.Account, error) {
	gsa, err := serviceaccount.NewAccount(ctx, externalsecrets.Ksa, &serviceaccount.AccountArgs{
		Project:     addedContainerClusterProject.ProjectId,
		Description: pulumi.String("external-secrets service account to retrieve secrets from google secrets manager"),
		AccountId:   pulumi.String(externalsecrets.Ksa),
		DisplayName: pulumi.String(externalsecrets.Ksa),
	}, pulumi.Parent(addedContainerClusterProject))
	if err != nil {
		return nil, errors.Wrapf(err, "failed add new %s ksa", externalsecrets.Ksa)
	}
	ctx.Export(GsaEmailOutputName), gsa.Email)
	return gsa, nil
}

func addIamRole(ctx *pulumi.Context, addedContainerClusterProject *organizations.Project, gsa *serviceaccount.Account) error {
	_, err := projects.NewIAMMember(ctx, externalsecrets.Ksa, &projects.IAMMemberArgs{
		Member:  pulumi.Sprintf("serviceAccount:%s", gsa.Email),
		Project: addedContainerClusterProject.ProjectId,
		Role:    pulumi.String(standard.Secretmanager_secretAccessor),
	}, pulumi.Parent(addedContainerClusterProject))
	if err != nil {
		return errors.Wrap(err, "failed to add iam roles to google service account")
	}
	return nil
}

func addWorkloadIdentityBinding(ctx *pulumi.Context, input *Input, gsa *serviceaccount.Account) error {
	_, err := serviceaccount.NewIAMBinding(ctx, fmt.Sprintf("%s-workload-identity", externalsecrets.Ksa), &serviceaccount.IAMBindingArgs{
		ServiceAccountId: gsa.Name,
		Role:             pulumi.String(standard.Iam_workloadIdentityUser),
		Members: pulumi.StringArray(getMembers(
			input.AddedContainerClusterProject,
			externalsecrets.Namespace,
			externalsecrets.Ksa,
		)),
	}, pulumi.Parent(gsa), pulumi.DependsOn([]pulumi.Resource{input.AddedContainerClusterResources.Cluster}))
	if err != nil {
		return errors.Wrapf(err, "failed to add workload identity binding for external secrets ksa to %v gsa", gsa.Email)
	}
	return nil
}

func GetGsaEmailOutputName     ) string {
	return pulumigoogleprovider.PulumiOutputName
	serviceaccount.Account{}, externalsecrets.Ksa)
}

func getMembers(addedContainerClusterProject *organizations.Project, kubernetesNamespace, kubernetesServiceAccount string) []pulumi.StringInput {
	return []pulumi.StringInput{
		pulumi.Sprintf("serviceAccount:%s.svc.id.goog[%s/%s]", addedContainerClusterProject.ProjectId, kubernetesNamespace, kubernetesServiceAccount),
	}
}
