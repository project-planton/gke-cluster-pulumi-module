package firewall

import (
	"fmt"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/pulumi/pulumicustomoutput"

	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/container/cluster/nodepool/tag"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	gkeWebhookAllowFirewallNameSuffix  = "gke-webhooks"
	containerClusterApiServerCidrBlock = "172.16.0.0/24"
	kubeWebhookPort                    = "8443"
	istioPilotWebhookPort              = "15017"
)

type Input struct {
	KubeClusterId     string
	AddedShareProject *organizations.Project
	AddedVpcNetwork   *compute.Network
}

func Resources(ctx *pulumi.Context, input *Input) error {
	if err := addNetworkFirewall(ctx, input); err != nil {
		return errors.Wrap(err, "failed to add network firewall")
	}
	return nil
}

func addNetworkFirewall(ctx *pulumi.Context, input *Input) error {
	firewallName := GetGkeWebhooksFirewallName(input.KubeClusterId)
	fw, err := compute.NewFirewall(ctx, firewallName, &compute.FirewallArgs{
		Name:    pulumi.String(firewallName),
		Project: input.AddedShareProject.ProjectId,
		Network: input.AddedVpcNetwork.Name,
		SourceRanges: pulumi.StringArray{
			pulumi.String(containerClusterApiServerCidrBlock),
		},
		Allows: compute.FirewallAllowArray{
			&compute.FirewallAllowArgs{
				Protocol: pulumi.String("tcp"),
				Ports: pulumi.StringArray{
					pulumi.String(kubeWebhookPort),
					pulumi.String(istioPilotWebhookPort),
				},
			},
		},
		TargetTags: pulumi.StringArray{
			pulumi.String(tag.Get(input.KubeClusterId)),
		},
	}, pulumi.Parent(input.AddedVpcNetwork))
	if err != nil {
		return errors.Wrap(err, "failed to add firewall")
	}
	ctx.Export(GkeWebhooksFirewallSelfLinkOutputName
	firewallName), fw.SelfLink)
	ctx.Export(ContainerClusterApiServersCidrBlockOutputName), pulumi.String(containerClusterApiServerCidrBlock))
	return nil
}

func GetGkeWebhooksFirewallName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s-%s", englishword.EnglishWord_kubernetes, kubeClusterId, gkeWebhookAllowFirewallNameSuffix)
}

func GetContainerClusterApiServersCidrBlockOutputName     ) string {
	return pulumicustomoutput.Name("container-cluster-api-servers-cidr-block")
}

func GetGkeWebhooksFirewallSelfLinkOutputName     firewallName string) string {
	return pulumigoogleprovider.PulumiOutputName
	compute.Firewall{}, firewallName)
}
