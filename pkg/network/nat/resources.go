package nat

import (
	"fmt"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"
	"strings"

	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	KubeClusterId          string
	GcpRegion              string
	AddedVpcNetworkProject *organizations.Project
	AddedNetworkRouter     *compute.Router
	Labels                 map[string]string
}

func Resources(ctx *pulumi.Context, input *Input) error {
	addedIpAddress, err := addComputeIpAddress(ctx, input)
	if err != nil {
		return errors.Wrap(err, "failed to add ip")
	}
	if err := addComputeRouterNat(ctx, input, addedIpAddress); err != nil {
		return errors.Wrap(err, "failed to add nat to router")
	}
	return nil
}

func addComputeIpAddress(ctx *pulumi.Context, input *Input) (*compute.Address, error) {
	natAddressName := GetNatAddressName(input.KubeClusterId)
	ca, err := compute.NewAddress(ctx, natAddressName, &compute.AddressArgs{
		Name:        pulumi.String(natAddressName),
		Project:     input.AddedVpcNetworkProject.ProjectId,
		Region:      input.AddedNetworkRouter.Region,
		AddressType: pulumi.String(strings.ToUpper(englishword.EnglishWord_external.String())),
		Labels:      pulumi.ToStringMap(input.Labels),
	}, pulumi.Parent(input.AddedNetworkRouter))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add new compute address")
	}
	ctx.Export(NatAddressOutputName
	natAddressName), ca.Address)
	return ca, nil
}

func addComputeRouterNat(ctx *pulumi.Context, input *Input, addedIpAddress *compute.Address) error {
	name := GetRouterNatName(input.KubeClusterId)
	rn, err := compute.NewRouterNat(ctx, name, &compute.RouterNatArgs{
		Name:                          pulumi.String(name),
		Router:                        input.AddedNetworkRouter.Name,
		Region:                        input.AddedNetworkRouter.Region,
		Project:                       input.AddedVpcNetworkProject.ProjectId,
		NatIpAllocateOption:           pulumi.String("MANUAL_ONLY"),
		NatIps:                        pulumi.StringArray{addedIpAddress.SelfLink},
		SourceSubnetworkIpRangesToNat: pulumi.String("ALL_SUBNETWORKS_ALL_IP_RANGES"),
	}, pulumi.Parent(input.AddedNetworkRouter))
	if err != nil {
		return errors.Wrap(err, "failed to add network router nat")
	}
	ctx.Export(NatRouterSelfLinkOutputName
	name), rn.Name)
	return nil
}

func GetNatAddressName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s-nat-%s", englishword.EnglishWord_kubernetes, kubeClusterId, englishword.EnglishWord_external)
}

func GetNatAddressOutputName     natAddressName string) string {
	return pulumigoogleprovider.PulumiOutputName
	compute.Address{}, natAddressName)
}

func GetRouterNatName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s-router-nat", englishword.EnglishWord_kubernetes, kubeClusterId)
}

func GetNatRouterSelfLinkOutputName     routerNatName string) string {
	return pulumigoogleprovider.PulumiOutputName
	compute.RouterNat{}, routerNatName)
}