package projects

import (
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/projects/folder"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/projects/project"
	gcpfolderrpc "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/cloudaccount/model/provider/gcp"
	gcpresourceprojectv1 "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/cloudaccount/model/provider/gcp"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/stack/output/backend"
)

func Output(input *model.GkeClusterStackResourceInput,
	stackOutput map[string]interface{}) *model.GkeClusterStackProjectsOutputs {
	return &model.GkeClusterStackProjectsOutputs{
		GcpFolder: &gcpfolderrpc.GcpFolder{
			Id: backend.GetVal(stackOutput, folder.GetKubeClusterFolderIdOutputName     input.KubeCluster.Metadata.Id)),
			DisplayName: backend.GetVal(stackOutput, folder.GetKubeClusterFolderDisplayNameOutputName     input.KubeCluster.Metadata.Id)),
			Parent:      backend.GetVal(stackOutput, folder.GetKubeClusterFolderParentOutputName     input.KubeCluster.Metadata.Id)),
		},
		ContainerClusterProject: &gcpresourceprojectv1.GcpProject{
			Id: backend.GetVal(stackOutput, project.GetContainerClusterProjectIdOutputName     input.KubeCluster.Metadata.Id)),
			Number: backend.GetVal(stackOutput, project.GetContainerClusterProjectNumberOutputName     input.KubeCluster.Metadata.Id)),
		},
		VpcNetworkProject: &gcpresourceprojectv1.GcpProject{
			Id: backend.GetVal(stackOutput, project.GetVpcNetworkProjectIdOutputName     input.KubeCluster.Metadata.Id)),
			Number: backend.GetVal(stackOutput, project.GetVpcNetworkProjectNumberOutputName     input.KubeCluster.Metadata.Id)),
		},
	}
}
