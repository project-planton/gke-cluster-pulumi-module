package apis

import "github.com/plantoncloud-inc/go-commons/cloud/gcp/apis"

var (
	VpcNetworkProjectApis = []apis.Api{
		apis.Compute,
		apis.Dns,
		apis.Container,
	}

	ContainerClusterProjectApis = []apis.Api{
		apis.Compute,
		apis.Container,
		apis.SecretManager,
	}
)
