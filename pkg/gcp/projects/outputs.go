package projects

import (
	kubernetesclustergcpstack "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	gcpfolderrpc "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/cloud/gcp/resource/folder/rpc"
	"buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/cloud/gcp/resource/project/rpc"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/projects/folder"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/projects/project"
	"github.com/plantoncloud-inc/pulumi-stack-runner-sdk/go/pulumi/stack/output/backend"
)

func Output(input *kubernetesclustergcpstack.KubeClusterGcpStackResourceInput, stackOutput map[string]interface{}) *kubernetesclustergcpstack.KubeClusterGcpStackProjectsOutputs {
	return &kubernetesclustergcpstack.KubeClusterGcpStackProjectsOutputs{
		GcpFolder: &gcpfolderrpc.GcpFolder{
			Id:          backend.GetVal(stackOutput, folder.GetKubeClusterFolderIdOutputName(input.KubeCluster.Metadata.Id)),
			DisplayName: backend.GetVal(stackOutput, folder.GetKubeClusterFolderDisplayNameOutputName(input.KubeCluster.Metadata.Id)),
			Parent:      backend.GetVal(stackOutput, folder.GetKubeClusterFolderParentOutputName(input.KubeCluster.Metadata.Id)),
		},
		ContainerClusterProject: &rpc.GcpProject{
			Id:     backend.GetVal(stackOutput, project.GetContainerClusterProjectIdOutputName(input.KubeCluster.Metadata.Id)),
			Number: backend.GetVal(stackOutput, project.GetContainerClusterProjectNumberOutputName(input.KubeCluster.Metadata.Id)),
		},
		VpcNetworkProject: &rpc.GcpProject{
			Id:     backend.GetVal(stackOutput, project.GetVpcNetworkProjectIdOutputName(input.KubeCluster.Metadata.Id)),
			Number: backend.GetVal(stackOutput, project.GetVpcNetworkProjectNumberOutputName(input.KubeCluster.Metadata.Id)),
		},
	}
}
