package network

import (
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network/ip"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network/nat"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/gcp/model"
	"github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
)

func Output(input *c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackResourceInput,
	stackOutput map[string]interface{}) *c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackNetworkOutputs {
	return &c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackNetworkOutputs{
		ExternalNatIpAddress: backend.GetVal(stackOutput, nat.GetNatAddressOutputName(nat.GetNatAddressName(input.KubeCluster.Metadata.Id))),
		IngressIpAddress: &c2cv1deployk8cstackgcpmodel.IngressIpAddress{
			Internal: backend.GetVal(stackOutput, ip.GetIngressInternalIpOutputName(input.KubeCluster.Metadata.Id)),
			External: backend.GetVal(stackOutput, ip.GetIngressExternalIpOutputName(input.KubeCluster.Metadata.Id)),
		},
	}
}
