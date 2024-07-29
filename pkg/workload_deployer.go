package pkg

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/outputs"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/vars"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func workloadDeployer(ctx *pulumi.Context, createdCluster *container.Cluster) (*serviceaccount.Key, error) {
	//create workload deployer service account
	createdWorkloadDeployerServiceAccount, err := serviceaccount.NewAccount(ctx,
		vars.WorkloadDeployServiceAccountName,
		&serviceaccount.AccountArgs{
			Project:     createdCluster.Project,
			Description: pulumi.String("service account to deploy workloads"),
			AccountId:   pulumi.String(vars.WorkloadDeployServiceAccountName),
			DisplayName: pulumi.String(vars.WorkloadDeployServiceAccountName),
		}, pulumi.Parent(createdCluster))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create workload deployer service account")
	}

	//export email of the created workload deployer service account
	ctx.Export(outputs.WorkloadDeployerGsaEmail, createdWorkloadDeployerServiceAccount.Email)

	//create key for workload deployer service account.
	createdWorkloadDeployerServiceAccountKey, err := serviceaccount.NewKey(ctx,
		vars.WorkloadDeployServiceAccountName,
		&serviceaccount.KeyArgs{
			ServiceAccountId: createdWorkloadDeployerServiceAccount.Name,
		}, pulumi.Parent(createdWorkloadDeployerServiceAccount))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create key for workload-deployer service account")
	}

	//export workload deployer google service account key
	ctx.Export(outputs.WorkloadDeployerGsaKey, createdWorkloadDeployerServiceAccountKey.PrivateKey)

	// create iam-binding for workload-deployer to manage the container cluster itself
	_, err = projects.NewIAMBinding(ctx,
		fmt.Sprintf("%s-container-admin", vars.WorkloadDeployServiceAccountName),
		&projects.IAMBindingArgs{
			Members: pulumi.StringArray{pulumi.Sprintf("serviceAccount:%s", createdWorkloadDeployerServiceAccount.Email)},
			Project: createdCluster.Project,
			Role:    pulumi.String("roles/container.admin"),
		}, pulumi.Parent(createdWorkloadDeployerServiceAccount))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create container-admin iam binding for workload deployer")
	}

	// create iam-binding for workload-deployer to manage resources inside container clusters
	_, err = projects.NewIAMBinding(ctx,
		fmt.Sprintf("%s-kube-cluster-admin", vars.WorkloadDeployServiceAccountName),
		&projects.IAMBindingArgs{
			Members: pulumi.StringArray{pulumi.Sprintf("serviceAccount:%s", createdWorkloadDeployerServiceAccount.Email)},
			Project: createdCluster.Project,
			Role:    pulumi.String("roles/container.clusterAdmin"),
		}, pulumi.Parent(createdWorkloadDeployerServiceAccount))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create cluster-admin iam binding for workload deployer")
	}

	return createdWorkloadDeployerServiceAccountKey, nil
}
