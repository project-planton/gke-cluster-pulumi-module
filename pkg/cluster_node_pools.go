package pkg

import (
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/localz"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func clusterNodePools(ctx *pulumi.Context,
	locals *localz.Locals,
	createdCluster *container.Cluster) ([]pulumi.Resource, error) {
	return nil, nil
}
