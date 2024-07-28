package outputs

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

const (
	ClusterCaData                 = "cluster-ca-data"
	ClusterEndpoint               = "cluster-endpoint"
	ClusterName                   = "cluster-name"
	ContainerClusterProjectId     = "container-cluster-project-id"
	ContainerClusterProjectNumber = "container-cluster-project-number"
	FolderDisplayName             = "folder-name"
	FolderId                      = "folder-id"
	FolderParent                  = "folder-parent"
	NetworkSelfLink               = "network-self-link"
	SubNetworkSelfLink            = "sub-network-self-link"
	VpcNetworkProjectId           = "vpc-network-project-id"
	VpcNetworkProjectNumber       = "vpc-network-project-number"
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
