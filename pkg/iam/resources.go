package iam

import (
	"github.com/pkg/errors"
	addoncertmanager "github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/certmanager"
	addonexternaldns "github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/externaldns"
	addonexternalsecrets "github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/externalsecrets"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam/certmanager"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam/dns"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam/externaldns"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam/externalsecrets"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam/workloaddeployer"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	AddedContainerClusterProject *organizations.Project
	AddedContainerClusters       *cluster.AddedContainerClusterResources
}

type AddedIamResources struct {
	CertManagerGsa         *serviceaccount.Account
	ExternalDnsGsa         *serviceaccount.Account
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
	addedExternalDnsGsa, err := externaldns.Resources(ctx, &externaldns.Input{
		AddedContainerClusterProject:   input.AddedContainerClusterProject,
		AddedContainerClusterResources: input.AddedContainerClusters,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add %s gsa", addonexternaldns.Ksa)
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
		ExternalDnsGsa:         addedExternalDnsGsa,
		ExternalSecretsGsa:     addedExternalSecretsGsa,
		WorkloadDeployerGsa:    addedWorkloadDeployerResources.AddedWorkloadDeployerGsa,
		WorkloadDeployerGsaKey: addedWorkloadDeployerResources.AddedWorkloadDeployerGsaKey,
	}, nil
}
