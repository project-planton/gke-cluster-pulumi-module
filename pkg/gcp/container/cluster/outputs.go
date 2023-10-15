package cluster

import (
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster/cluster"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	kubernetesclustergcpstack "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
)

func Output(input *kubernetesclustergcpstack.KubeClusterGcpStackResourceInput,
	stackOutput map[string]interface{}) *kubernetesclustergcpstack.ClusterOutputs {
	clusterName := input.KubeCluster.Metadata.Id
	return &kubernetesclustergcpstack.ClusterOutputs{
		ClusterEndpoint: backend.GetVal(stackOutput, cluster.GetClusterEndpointOutputName(clusterName)),
		ClusterCaData:   backend.GetVal(stackOutput, cluster.GetClusterCaDataOutputName(clusterName)),
	}
}
