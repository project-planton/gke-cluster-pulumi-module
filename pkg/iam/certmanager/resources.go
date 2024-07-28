package certmanager

import (
	"fmt"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"

	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/certmanager"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
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
	if err := addWorkloadIdentityBinding(ctx, input, gsa); err != nil {
		return nil, errors.Wrap(err, "failed to add workload identity binding")
	}
	return gsa, nil
}

func addGsa(ctx *pulumi.Context, addedContainerClusterProject *organizations.Project) (*serviceaccount.Account, error) {

}

func addWorkloadIdentityBinding(ctx *pulumi.Context, input *Input, gsa *serviceaccount.Account) error {

}

func GetGsaEmailOutputName     ) string {
	return pulumigoogleprovider.PulumiOutputName
	serviceaccount.Account{}, certmanager.Ksa)
}

func getMembers(addedProject *organizations.Project, kubernetesNamespace, kubernetesServiceAccount string) []pulumi.StringInput {
	return []pulumi.StringInput{
		pulumi.Sprintf("serviceAccount:%s.svc.id.goog[%s/%s]", addedProject.ProjectId, kubernetesNamespace, kubernetesServiceAccount),
	}
}
