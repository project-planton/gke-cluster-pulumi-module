package cluster

import (
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster/cluster"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp/model"
)

func Output(input *c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackResourceInput,
	stackOutput map[string]interface{}) *c2cv1deployk8cstackgcpmodel.ClusterOutputs {
	clusterName := input.KubeCluster.Metadata.Id
	return &c2cv1deployk8cstackgcpmodel.ClusterOutputs{
		ClusterEndpoint: backend.GetVal(stackOutput, cluster.GetClusterEndpointOutputName(clusterName)),
		ClusterCaData:   backend.GetVal(stackOutput, cluster.GetClusterCaDataOutputName(clusterName)),
	}
}
