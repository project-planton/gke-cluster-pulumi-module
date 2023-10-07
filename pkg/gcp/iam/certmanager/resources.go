package certmanager

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/addon/certmanager"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/cluster"
	puluminameoutputgcp "github.com/plantoncloud-inc/pulumi-stack-runner-sdk/go/pulumi/name/provider/cloud/gcp/output"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
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
	if err := addWorkloadIdentityBinding(ctx, input, gsa); err != nil {
		return nil, errors.Wrap(err, "failed to add workload identity binding")
	}
	return gsa, nil
}

func addGsa(ctx *pulumi.Context, addedContainerClusterProject *organizations.Project) (*serviceaccount.Account, error) {
	gsa, err := serviceaccount.NewAccount(ctx, certmanager.Ksa, &serviceaccount.AccountArgs{
		Project:     addedContainerClusterProject.ProjectId,
		Description: pulumi.String("cert-manager service account for solving dns challenges to issue certificates"),
		AccountId:   pulumi.String(certmanager.Ksa),
		DisplayName: pulumi.String(certmanager.Ksa),
	}, pulumi.Parent(addedContainerClusterProject))
	if err != nil {
		return nil, errors.Wrapf(err, "failed add new %s svc acct", certmanager.Ksa)
	}
	ctx.Export(GetGsaEmailOutputName(), gsa.Email)
	return gsa, nil
}

func addWorkloadIdentityBinding(ctx *pulumi.Context, input *Input, gsa *serviceaccount.Account) error {
	_, err := serviceaccount.NewIAMBinding(ctx, fmt.Sprintf("%s-workload-identity", certmanager.Ksa), &serviceaccount.IAMBindingArgs{
		ServiceAccountId: gsa.Name,
		Role:             pulumi.String(standard.Iam_workloadIdentityUser),
		Members:          pulumi.StringArray(getMembers(input.AddedContainerClusterProject, certmanager.Namespace, certmanager.Ksa)),
	}, pulumi.Parent(gsa),
		pulumi.DependsOn([]pulumi.Resource{input.AddedContainerClusterResources.Cluster}))
	if err != nil {
		return errors.Wrapf(err, "failed to add workload identity binding for external secrets ksa to %v gsa", gsa.Email)
	}
	return nil
}

func GetGsaEmailOutputName() string {
	return puluminameoutputgcp.Name(serviceaccount.Account{}, certmanager.Ksa)
}

func getMembers(addedProject *organizations.Project, kubernetesNamespace, kubernetesServiceAccount string) []pulumi.StringInput {
	return []pulumi.StringInput{
		pulumi.Sprintf("serviceAccount:%s.svc.id.goog[%s/%s]", addedProject.ProjectId, kubernetesNamespace, kubernetesServiceAccount),
	}
}
