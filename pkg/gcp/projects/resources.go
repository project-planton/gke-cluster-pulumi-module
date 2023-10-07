package projects

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/projects/folder"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/projects/project"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	CloudAccountFolderId string
	IsCreateSharedVpc    bool
	KubeClusterId        string
	GcpRegion            string
	GcpZone              string
	BillingAccountId     string
	GcpProvider          *pulumigcp.Provider
	Labels               map[string]string
}

// Resources sets up kube-cluster projects
// * creates kube-cluster and vpc-network projects
// * enabled required apis on kube-cluster and network projects
func Resources(ctx *pulumi.Context, input *Input) (*project.AddedProjectsResources, error) {
	addedFolder, err := folder.Resources(ctx, &folder.Input{
		GcpProvider:          input.GcpProvider,
		CloudAccountFolderId: input.CloudAccountFolderId,
		KubeClusterId:        input.KubeClusterId,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add kube-cluster folder")
	}
	addedProjectsResources, err := project.Resources(ctx, &project.Input{
		AddedKubeClusterFolder: addedFolder,
		IsSharedVpcEnabled:     input.IsCreateSharedVpc,
		KubeClusterId:          input.KubeClusterId,
		GcpRegion:              input.GcpRegion,
		GcpZone:                input.GcpZone,
		BillingAccountId:       input.BillingAccountId,
		GcpProvider:            input.GcpProvider,
		Labels:                 input.Labels,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to add projects resources")
	}
	return addedProjectsResources, nil
}
