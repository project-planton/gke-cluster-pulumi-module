package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	CertManagerGsaEmail                 = "cert-manager-gsa-email"
	ClusterCaData                       = "cluster-ca-data"
	ClusterEndpoint                     = "cluster-endpoint"
	ClusterName                         = "cluster-name"
	ContainerClusterApiServersCidrBlock = "container-cluster-api-servers-cidr-block"
	ContainerClusterProjectId           = "container-cluster-project-id"
	ContainerClusterProjectNumber       = "container-cluster-project-number"
	ExternalSecretsGsaEmail             = "external-secrets-gsa-email"
	FolderDisplayName                   = "folder-name"
	FolderId                            = "folder-id"
	FolderParent                        = "folder-parent"
	GkeWebhooksFirewallSelfLink         = "gke-webhooks-firewall-self-link"
	NatIpAddress                        = "nat-ip-address"
	NetworkSelfLink                     = "network-self-link"
	RouterNatName                       = "router-nat-name"
	RouterSelfLink                      = "router-self-link"
	SubNetworkSelfLink                  = "sub-network-self-link"
	VpcNetworkProjectId                 = "vpc-network-project-id"
	VpcNetworkProjectNumber             = "vpc-network-project-number"
	WorkloadDeployerGsaEmail            = "workload-deployer-gsa-email"
	WorkloadDeployerGsaKey              = "workload-deployer-gsa-key"
)

func PulumiOutputToStackOutputsConverter(pulumiOutputs auto.OutputMap,
	input *model.GkeClusterStackInput) *model.GkeClusterStackOutputs {
	return &model.GkeClusterStackOutputs{
		Folder:                       nil,
		ContainerClusterProject:      nil,
		VpcNetworkProject:            nil,
		ClusterEndpoint:              "",
		ClusterCaData:                "",
		ExternalNatIp:                "",
		InternalIngressIp:            "",
		ExternalIngressIp:            "",
		CertManagerGsaEmail:          "",
		ExternalSecretsGsaEmail:      "",
		WorkloadDeployerGsaEmail:     "",
		WorkloadDeployerGsaKeyBase64: "",
		ExternalDnsGsaEmail:          "",
	}
}
