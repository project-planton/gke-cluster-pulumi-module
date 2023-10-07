package iam

import (
	"github.com/pkg/errors"
	addoncertmanager "github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/certmanager"
	addonexternalsecrets "github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/externalsecrets"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/iam/certmanager"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/iam/dns"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/iam/externalsecrets"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/iam/workloaddeployer"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	AddedContainerClusterProject *organizations.Project
	AddedContainerClusters       *cluster.AddedContainerClusterResources
}

type AddedIamResources struct {
	CertManagerGsa         *serviceaccount.Account
	ExternalSecretsGsa     *serviceaccount.Account
	WorkloadDeployerGsa    *serviceaccount.Account
	WorkloadDeployerGsaKey *serviceaccount.Key
}

func Resources(ctx *pulumi.Context, input *Input) (*AddedIamResources, error) {
	addedCertManagerGsa, err := certmanager.Resources(ctx, &certmanager.Input{
		AddedContainerClusterProject:   input.AddedContainerClusterProject,
		AddedContainerClusterResources: input.AddedContainerClusters,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s gsa", addoncertmanager.Ksa)
	}
	addedExternalSecretsGsa, err := externalsecrets.Resources(ctx, &externalsecrets.Input{
		AddedContainerClusterProject:   input.AddedContainerClusterProject,
		AddedContainerClusterResources: input.AddedContainerClusters,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s gsa", addonexternalsecrets.Ksa)
	}
	addedWorkloadDeployerResources, err := workloaddeployer.Resources(ctx, &workloaddeployer.Input{
		AddedContainerClusterProject: input.AddedContainerClusterProject,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s gsa", workloaddeployer.GsaName)
	}
	err = dns.Resources(ctx, &dns.Input{
		AddedWorkloadDeployerGsa:     addedWorkloadDeployerResources.AddedWorkloadDeployerGsa,
		AddedCertManagerGsa:          addedCertManagerGsa,
		AddedContainerClusterProject: input.AddedContainerClusterProject,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add dns iam roles")
	}
	return &AddedIamResources{
		CertManagerGsa:         addedCertManagerGsa,
		ExternalSecretsGsa:     addedExternalSecretsGsa,
		WorkloadDeployerGsa:    addedWorkloadDeployerResources.AddedWorkloadDeployerGsa,
		WorkloadDeployerGsaKey: addedWorkloadDeployerResources.AddedWorkloadDeployerGsaKey,
	}, nil
}