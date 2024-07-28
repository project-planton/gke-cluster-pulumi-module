package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/localz"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input            *model.GkeClusterStackInput
	GcpLabels        map[string]string
	KubernetesLabels map[string]string
}

func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	locals := localz.Initialize(ctx, s.Input)

	//create gcp-provider using the gcp-credential from input
	gcpProvider, err := pulumigoogleprovider.Get(ctx, s.Input.GcpCredential)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	//create gcp folder
	createdFolder, err := s.folder(ctx, locals, gcpProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create folder")
	}

	//create cluster
	createdCluster, err := cluster(ctx, locals, createdFolder)
	if err != nil {
		return errors.Wrap(err, "failed to create container cluster")
	}

	//create node-pools
	createdNodePools, err := clusterNodePools(ctx, locals, createdCluster)
	if err != nil {
		return errors.Wrap(err, "failed to create cluster node-pools")
	}

	//create addons
	if err := clusterAddons(ctx, locals, gcpProvider,
		createdCluster, createdNodePools); err != nil {
		return errors.Wrap(err, "failed to create addons")
	}
	return nil
}
