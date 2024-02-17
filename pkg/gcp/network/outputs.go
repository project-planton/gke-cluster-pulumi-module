package network

import (
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/network/ip"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/network/nat"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp/model"
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
