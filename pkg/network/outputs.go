package network

import (
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network/nat"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/stack/output/backend"
)

func Output(input *model.GkeClusterStackResourceInput,
	stackOutput map[string]interface{}) *model.GkeClusterStackNetworkOutputs {
	return &model.GkeClusterStackNetworkOutputs{
		ExternalNatIpAddress: backend.GetVal(stackOutput, nat.GetNatAddressOutputName     nat.GetNatAddressName(input.KubeCluster.Metadata.Id))),
		IngressIpAddress: &model.IngressIpAddress{
		Internal: backend.GetVal(stackOutput, ip.GetIngressInternalIpOutputName     input.KubeCluster.Metadata.Id)),
		External: backend.GetVal(stackOutput, ip.GetIngressExternalIpOutputName     input.KubeCluster.Metadata.Id)),
	},
	}
}
