package nat

import (
	"fmt"
	"github.com/pkg/errors"
	puluminameoutputgcp "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	wordpb "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/commons/english/enums"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
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
		AddressType: pulumi.String(strings.ToUpper(wordpb.Word_external.String())),
		Labels:      pulumi.ToStringMap(input.Labels),
	}, pulumi.Parent(input.AddedNetworkRouter))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add new compute address")
	}
	ctx.Export(GetNatAddressOutputName(natAddressName), ca.Address)
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
	ctx.Export(GetNatRouterSelfLinkOutputName(name), rn.Name)
	return nil
}

func GetNatAddressName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s-nat-%s", wordpb.Word_kubernetes, kubeClusterId, wordpb.Word_external)
}

func GetNatAddressOutputName(natAddressName string) string {
	return puluminameoutputgcp.Name(compute.Address{}, natAddressName)
}

func GetRouterNatName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s-router-nat", wordpb.Word_kubernetes, kubeClusterId)
}

func GetNatRouterSelfLinkOutputName(routerNatName string) string {
	return puluminameoutputgcp.Name(compute.RouterNat{}, routerNatName)
}
