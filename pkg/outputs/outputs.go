package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/automationapi/autoapistackoutput"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	CertManagerGsaEmail           = "cert-manager-gsa-email"
	ExternalDnsGsaEmail           = "external-dns-gsa-email"
	ClusterCaData                 = "cluster-ca-data"
	ClusterEndpoint               = "cluster-endpoint"
	ContainerClusterProjectId     = "container-cluster-project-id"
	ContainerClusterProjectNumber = "container-cluster-project-number"
	ExternalSecretsGsaEmail       = "external-secrets-gsa-email"
	IngressExternalIp             = "ingress-external-ip"
	IngressInternalIp             = "ingress-internal-ip"
	FolderDisplayName             = "folder-name"
	FolderId                      = "folder-id"
	FolderParent                  = "folder-parent"
	GkeWebhooksFirewallSelfLink   = "gke-webhooks-firewall-self-link"
	NatIpAddress                  = "nat-ip-address"
	NetworkSelfLink               = "network-self-link"
	RouterNatName                 = "router-nat-name"
	RouterSelfLink                = "router-self-link"
	SubNetworkSelfLink            = "sub-network-self-link"
	VpcNetworkProjectId           = "vpc-network-project-id"
	VpcNetworkProjectNumber       = "vpc-network-project-number"
	WorkloadDeployerGsaEmail      = "workload-deployer-gsa-email"
	WorkloadDeployerGsaKey        = "workload-deployer-gsa-key"
)

func PulumiOutputsToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *gkecluster.GkeClusterStackInput) *gkecluster.GkeClusterStackOutputs {
	return &gkecluster.GkeClusterStackOutputs{
		FolderId: autoapistackoutput.GetVal(pulumiOutputs, FolderId),
		ContainerClusterProject: &gkecluster.GcpProject{
			Id:     autoapistackoutput.GetVal(pulumiOutputs, ContainerClusterProjectId),
			Number: autoapistackoutput.GetVal(pulumiOutputs, ContainerClusterProjectNumber),
		},
		VpcNetworkProject: &gkecluster.GcpProject{
			Id:     autoapistackoutput.GetVal(pulumiOutputs, VpcNetworkProjectId),
			Number: autoapistackoutput.GetVal(pulumiOutputs, VpcNetworkProjectNumber),
		},
		ClusterEndpoint:              autoapistackoutput.GetVal(pulumiOutputs, ClusterEndpoint),
		ClusterCaData:                autoapistackoutput.GetVal(pulumiOutputs, ClusterCaData),
		ExternalNatIp:                autoapistackoutput.GetVal(pulumiOutputs, NatIpAddress),
		InternalIngressIp:            autoapistackoutput.GetVal(pulumiOutputs, IngressInternalIp),
		ExternalIngressIp:            autoapistackoutput.GetVal(pulumiOutputs, IngressExternalIp),
		CertManagerGsaEmail:          autoapistackoutput.GetVal(pulumiOutputs, CertManagerGsaEmail),
		ExternalSecretsGsaEmail:      autoapistackoutput.GetVal(pulumiOutputs, ExternalSecretsGsaEmail),
		WorkloadDeployerGsaEmail:     autoapistackoutput.GetVal(pulumiOutputs, WorkloadDeployerGsaEmail),
		WorkloadDeployerGsaKeyBase64: autoapistackoutput.GetVal(pulumiOutputs, WorkloadDeployerGsaKey),
		ExternalDnsGsaEmail:          autoapistackoutput.GetVal(pulumiOutputs, ExternalDnsGsaEmail),
	}
}
