package projects

import (
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/projects/folder"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/projects/project"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	gcpfolderrpc "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/cloudaccount/provider/gcp/resource/folder"
	gcpresourceprojectv1 "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/cloudaccount/provider/gcp/resource/project"
	kubernetesclustergcpstack "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
)

func Output(input *kubernetesclustergcpstack.KubeClusterGcpStackResourceInput, stackOutput map[string]interface{}) *kubernetesclustergcpstack.KubeClusterGcpStackProjectsOutputs {
	return &kubernetesclustergcpstack.KubeClusterGcpStackProjectsOutputs{
		GcpFolder: &gcpfolderrpc.GcpFolder{
			Id:          backend.GetVal(stackOutput, folder.GetKubeClusterFolderIdOutputName(input.KubeCluster.Metadata.Id)),
			DisplayName: backend.GetVal(stackOutput, folder.GetKubeClusterFolderDisplayNameOutputName(input.KubeCluster.Metadata.Id)),
			Parent:      backend.GetVal(stackOutput, folder.GetKubeClusterFolderParentOutputName(input.KubeCluster.Metadata.Id)),
		},
		ContainerClusterProject: &gcpresourceprojectv1.GcpProject{
			Id:     backend.GetVal(stackOutput, project.GetContainerClusterProjectIdOutputName(input.KubeCluster.Metadata.Id)),
			Number: backend.GetVal(stackOutput, project.GetContainerClusterProjectNumberOutputName(input.KubeCluster.Metadata.Id)),
		},
		VpcNetworkProject: &gcpresourceprojectv1.GcpProject{
			Id:     backend.GetVal(stackOutput, project.GetVpcNetworkProjectIdOutputName(input.KubeCluster.Metadata.Id)),
			Number: backend.GetVal(stackOutput, project.GetVpcNetworkProjectNumberOutputName(input.KubeCluster.Metadata.Id)),
		},
	}
}
