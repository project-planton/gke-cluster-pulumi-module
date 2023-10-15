package network

import (
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/network/ip"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/network/nat"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	kubernetesclustergcpstack "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
)

func Output(input *kubernetesclustergcpstack.KubeClusterGcpStackResourceInput, stackOutput map[string]interface{}) *kubernetesclustergcpstack.KubeClusterGcpStackNetworkOutputs {
	return &kubernetesclustergcpstack.KubeClusterGcpStackNetworkOutputs{
		ExternalNatIpAddress: backend.GetVal(stackOutput, nat.GetNatAddressOutputName(nat.GetNatAddressName(input.KubeCluster.Metadata.Id))),
		IngressIpAddress: &kubernetesclustergcpstack.IngressIpAddress{
			Internal: backend.GetVal(stackOutput, ip.GetIngressInternalIpOutputName(input.KubeCluster.Metadata.Id)),
			External: backend.GetVal(stackOutput, ip.GetIngressExternalIpOutputName(input.KubeCluster.Metadata.Id)),
		},
	}
}
