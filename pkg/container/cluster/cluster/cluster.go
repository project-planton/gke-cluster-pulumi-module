package cluster

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/network/subnet"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/enums/gkereleasechannel"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/pulumi/pulumicustomoutput"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	gkeReleaseChannel           = gkereleasechannel.GkeReleaseChannel_STABLE
	autoscalingProfileBalanced  = "BALANCED"
	autoscalingProfileOptimized = "OPTIMIZE_UTILIZATION"
)

type Input struct {
	KubeClusterId                string
	GcpZone                      string
	ClusterName                  string
	AddedContainerClusterProject *organizations.Project
	AddedNetworkResources        *network.AddedNetworkResources
	IsWorkloadLogsEnabled        bool
	ClusterConfig                *model.ClusterConfig
	ClusterAutoscalingConfig     *code2cloudv1deployk8cmodel.GkeClusterClusterAutoscalingConfigSpec
}

func Resources(ctx *pulumi.Context, input *Input) (*container.Cluster, error) {
	clusterName := input.KubeClusterId

	ctx.Export(ClusterNameOutputName
	clusterName), cc.Name)
	ctx.Export(ClusterEndpointOutputName
	clusterName), cc.Endpoint)
	ctx.Export(ApiServerCidrBlockOutputName
	clusterName), cc.PrivateClusterConfig.MasterIpv4CidrBlock())
	ctx.Export(ClusterCaDataOutputName
	clusterName), cc.MasterAuth.ClusterCaCertificate().Elem())
	return cc, nil
}

func getWorkloadIdentityNamespace(addedGcpProject *organizations.Project) pulumi.StringOutput {
	return pulumi.Sprintf("%s.svc.id.goog", addedGcpProject.ProjectId)
}

func getLoggingComponents(isWorkloadLogsEnabled bool) pulumi.StringArray {
	comps := pulumi.StringArray{
		pulumi.String("SYSTEM_COMPONENTS"),
	}
	if isWorkloadLogsEnabled {
		comps = append(comps, pulumi.String("WORKLOADS"))
	}
	return comps
}

func GetApiServerCidrBlockOutputName     clusterFullName string) string {
	return pulumicustomoutput.Name(clusterFullName, "api-server-ip-cidr")
}

func GetClusterNameOutputName     clusterFullName string) string {
	return pulumigoogleprovider.PulumiOutputName
	container.Cluster{}, clusterFullName, englishword.EnglishWord_name.String())
}
