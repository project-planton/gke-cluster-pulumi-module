package cluster

import (
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster/cluster"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/stack/output/backend"
)

func Output(input *model.GkeClusterStackResourceInput,
	stackOutput map[string]interface{}) *model.ClusterOutputs {
	clusterName := input.KubeCluster.Metadata.Id
	return &model.ClusterOutputs{
		ClusterEndpoint: backend.GetVal(stackOutput, cluster.GetClusterEndpointOutputName     clusterName)),
		ClusterCaData:   backend.GetVal(stackOutput, cluster.GetClusterCaDataOutputName     clusterName)),
	}
}
