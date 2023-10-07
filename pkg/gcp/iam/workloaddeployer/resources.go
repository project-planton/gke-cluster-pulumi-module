package workloaddeployer

import (
	"buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/english/rpc/enums"
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	puluminameoutputgcp "github.com/plantoncloud-inc/pulumi-stack-runner-sdk/go/pulumi/name/provider/cloud/gcp/output"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	GsaName = "workload-deployer"
)

type Input struct {
	AddedContainerClusterProject *organizations.Project
}

type AddedWorkloadServiceAccountResources struct {
	AddedWorkloadDeployerGsa    *serviceaccount.Account
	AddedWorkloadDeployerGsaKey *serviceaccount.Key
}

func Resources(ctx *pulumi.Context, input *Input) (*AddedWorkloadServiceAccountResources, error) {
	addedGsa, err := addGsa(ctx, input.AddedContainerClusterProject)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add gsa")
	}
	addedGsaKey, err := addGsaKey(ctx, addedGsa)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add gsa key")
	}
	if err := addWorkloadDeployerRoleBinding(ctx, input, addedGsa); err != nil {
		return nil, errors.Wrap(err, "failed to add workload deployer custom role")
	}
	return &AddedWorkloadServiceAccountResources{
		AddedWorkloadDeployerGsa:    addedGsa,
		AddedWorkloadDeployerGsaKey: addedGsaKey,
	}, nil
}

func addGsa(ctx *pulumi.Context, addedGcpProject *organizations.Project) (*serviceaccount.Account, error) {
	gsa, err := serviceaccount.NewAccount(ctx, GsaName, &serviceaccount.AccountArgs{
		Project:     addedGcpProject.ProjectId,
		Description: pulumi.String("service account to deploy workloads"),
		AccountId:   pulumi.String(GsaName),
		DisplayName: pulumi.String(GsaName),
	}, pulumi.Parent(addedGcpProject))
	if err != nil {
		return nil, errors.Wrapf(err, "failed add new %s svc acct", GsaName)
	}
	ctx.Export(GetGsaEmailOutputName(), gsa.Email)
	return gsa, nil
}

func addGsaKey(ctx *pulumi.Context, addedGsa *serviceaccount.Account) (*serviceaccount.Key, error) {
	addedGsaKey, err := serviceaccount.NewKey(ctx, GsaName, &serviceaccount.KeyArgs{
		ServiceAccountId: addedGsa.Name,
	}, pulumi.Parent(addedGsa))
	if err != nil {
		return nil, errors.Wrapf(err, "failed add new %s svc acct", GsaName)
	}
	ctx.Export(GetGsaKeyOutputName(), addedGsaKey.PrivateKey)
	return addedGsaKey, nil
}

// addWorkloadDeployerRoleBinding grants the required roles to do the following
// * add and delete secrets in kube-cluster project
// * add and delete dns managed zones and records in kube-cluster project
// * manage resources inside container clusters in kube-cluster
// * add and delete storage buckets in kube-cluster project
func addWorkloadDeployerRoleBinding(ctx *pulumi.Context, input *Input, addedGsa *serviceaccount.Account) error {
	//container cluster project role bindings
	//dns-admin role required for workload-deployer service account is managed in dns package as cert-manager also requires the same role in sha
	_, err := projects.NewIAMBinding(ctx, fmt.Sprintf("%s-secrets-admin", GsaName), &projects.IAMBindingArgs{
		Members: pulumi.StringArray{pulumi.Sprintf("serviceAccount:%s", addedGsa.Email)},
		Project: input.AddedContainerClusterProject.ProjectId,
		Role:    pulumi.String(standard.Secretmanager_admin),
	}, pulumi.Parent(addedGsa))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role binding for %s service account in share project", standard.Secretmanager_admin, GsaName)
	}
	//required for creating workload identity service accounts
	_, err = projects.NewIAMBinding(ctx, fmt.Sprintf("%s-iam-service-account-admin", GsaName), &projects.IAMBindingArgs{
		Members: pulumi.StringArray{pulumi.Sprintf("serviceAccount:%s", addedGsa.Email)},
		Project: input.AddedContainerClusterProject.ProjectId,
		Role:    pulumi.String(standard.Iam_serviceAccountAdmin),
	}, pulumi.Parent(addedGsa))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role binding for %s service account in share project", standard.Iam_serviceAccountAdmin, GsaName)
	}
	_, err = projects.NewIAMBinding(ctx, fmt.Sprintf("%s-kube-cluster-admin", GsaName), &projects.IAMBindingArgs{
		Members: pulumi.StringArray{pulumi.Sprintf("serviceAccount:%s", addedGsa.Email)},
		Project: input.AddedContainerClusterProject.ProjectId,
		Role:    pulumi.String(standard.Container_clusterAdmin),
	}, pulumi.Parent(addedGsa))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role binding for %s service account in kube-cluster project",
			standard.Container_clusterAdmin, GsaName)
	}
	_, err = projects.NewIAMBinding(ctx, fmt.Sprintf("%s-container-admin", GsaName), &projects.IAMBindingArgs{
		Members: pulumi.StringArray{pulumi.Sprintf("serviceAccount:%s", addedGsa.Email)},
		Project: input.AddedContainerClusterProject.ProjectId,
		Role:    pulumi.String(standard.Container_admin),
	}, pulumi.Parent(addedGsa))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role binding for %s service account in kube-cluster project",
			standard.Container_admin, GsaName)
	}
	_, err = projects.NewIAMBinding(ctx, fmt.Sprintf("%s-storage-admin", GsaName), &projects.IAMBindingArgs{
		Members: pulumi.StringArray{pulumi.Sprintf("serviceAccount:%s", addedGsa.Email)},
		Project: input.AddedContainerClusterProject.ProjectId,
		Role:    pulumi.String(standard.Storage_admin),
	}, pulumi.Parent(addedGsa))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s role binding for %s service account in kube-cluster project",
			standard.Storage_admin, GsaName)
	}
	return nil
}

func GetGsaEmailOutputName() string {
	return puluminameoutputgcp.Name(serviceaccount.Account{}, GsaName)
}

func GetGsaKeyOutputName() string {
	return puluminameoutputgcp.Name(serviceaccount.Key{}, GsaName, enums.Word_key.String())
}
