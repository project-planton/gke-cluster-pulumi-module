// Package localz instead of locals to avoid naming collision w/ "locals" for the instance name created for the struct.
package localz

import (
	gcpcredentialv1 "buf.build/gen/go/plantoncloud/project-planton/protocolbuffers/go/project/planton/apis/credential/gcpcredential/v1"
	gkeclusterv1 "buf.build/gen/go/plantoncloud/project-planton/protocolbuffers/go/project/planton/apis/provider/gcp/gkecluster/v1"
	"fmt"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/gcp/gcplabelkeys"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/kubernetes/kuberneteslabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strconv"
)

type Locals struct {
	GcpCredentialSpec                     *gcpcredentialv1.GcpCredentialSpec
	GkeCluster                            *gkeclusterv1.GkeCluster
	KubernetesPodSecondaryIpRangeName     string
	KubernetesServiceSecondaryIpRangeName string
	KubernetesLabels                      map[string]string
	GcpLabels                             map[string]string
	ContainerClusterLoggingComponentList  []string
	NetworkTag                            string
}

func Initialize(ctx *pulumi.Context, stackInput *gkeclusterv1.GkeClusterStackInput) *Locals {
	gkeCluster := stackInput.Target

	locals := &Locals{}

	locals.GcpCredentialSpec = stackInput.GcpCredential
	locals.GkeCluster = stackInput.Target

	locals.GcpLabels = map[string]string{
		gcplabelkeys.Resource:     strconv.FormatBool(true),
		gcplabelkeys.Organization: locals.GkeCluster.Spec.EnvironmentInfo.OrgId,
		gcplabelkeys.ResourceKind: "gke_cluster",
		gcplabelkeys.ResourceId:   locals.GkeCluster.Metadata.Id,
	}

	locals.KubernetesLabels = map[string]string{
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.Organization: locals.GkeCluster.Spec.EnvironmentInfo.OrgId,
		kuberneteslabelkeys.ResourceKind: "gke_cluster",
		kuberneteslabelkeys.ResourceId:   locals.GkeCluster.Metadata.Id,
	}

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
