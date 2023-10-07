package cluster

import (
	kubernetesclustergcpstack "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/cluster/cluster"
	"github.com/plantoncloud-inc/pulumi-stack-runner-sdk/go/pulumi/stack/output/backend"
)

func Output(input *kubernetesclustergcpstack.KubeClusterGcpStackResourceInput,
	stackOutput map[string]interface{}) *kubernetesclustergcpstack.ClusterOutputs {
	clusterName := input.KubeCluster.Metadata.Id
	return &kubernetesclustergcpstack.ClusterOutputs{
		ClusterEndpoint: backend.GetVal(stackOutput, cluster.GetClusterEndpointOutputName(clusterName)),
		ClusterCaData:   backend.GetVal(stackOutput, cluster.GetClusterCaDataOutputName(clusterName)),
	}
}
