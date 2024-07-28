// Package localz instead of locals to avoid naming collision w/ "locals" for the instance name created for the struct.
package localz

import (
	"fmt"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	gcpcredential "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/gcpcredential/model"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	GcpCredential                         *gcpcredential.GcpCredential
	GkeCluster                            *model.GkeCluster
	KubernetesPodSecondaryIpRangeName     string
	KubernetesServiceSecondaryIpRangeName string
	KubernetesLabels                      map[string]string
	GcpLabels                             map[string]string
	ContainerClusterLoggingComponentList  []string
	NetworkTag                            string
}

func Initialize(ctx *pulumi.Context, stackInput *model.GkeClusterStackInput) *Locals {
	gkeCluster := stackInput.ApiResource

	locals := &Locals{}

	locals.GcpCredential = stackInput.GcpCredential
	locals.GkeCluster = stackInput.ApiResource

	locals.KubernetesPodSecondaryIpRangeName = fmt.Sprintf("%s-pods", gkeCluster.Metadata.Id)
	locals.KubernetesServiceSecondaryIpRangeName = fmt.Sprintf("%s-services", gkeCluster.Metadata.Id)
	locals.NetworkTag = gkeCluster.Metadata.Id

	locals.ContainerClusterLoggingComponentList = []string{"SYSTEM_COMPONENTS"}

	if gkeCluster.Spec.IsWorkloadLogsEnabled {
		locals.ContainerClusterLoggingComponentList = append(locals.ContainerClusterLoggingComponentList,
			"WORKLOADS")
	}

	return locals
}
