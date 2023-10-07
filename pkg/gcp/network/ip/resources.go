package ip

import (
	wordpb "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/english/rpc/enums"
	"buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/network/ip/rpc/enums"
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/projects/project"
	puluminamegcpoutput "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

type Input struct {
	KubeClusterId          string
	GcpRegion              string
	AddedProjectsResources *project.AddedProjectsResources
	AddedSubnet            *compute.Subnetwork
	Labels                 map[string]string
}

type AddedComputeIpAddresses struct {
	External *compute.Address
	Internal *compute.Address
}

// Resources adds one set of external and one internal ip addresses reservations.
// these ip address are attached to the load balancers created by services on container cluster to configure ingress
func Resources(ctx *pulumi.Context, input *Input) (*AddedComputeIpAddresses, error) {
	addedIpAddresses, err := addIpAddresses(ctx, input,
		input.AddedProjectsResources.KubeClusterProjects.ContainerClusterProject,
		input.AddedProjectsResources.AddedProjectApis.ContainerClusterProject)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add ingress ip addresses")
	}
	exportOutputs(ctx, input.KubeClusterId, addedIpAddresses)
	return addedIpAddresses, nil
}

func exportOutputs(ctx *pulumi.Context, kubeClusterId string, addedIpAddresses *AddedComputeIpAddresses) {
	ctx.Export(getIngressIpOutputName(enums.IpAddressVisibility_IP_ADDRESS_VISIBILITY_EXTERNAL, kubeClusterId),
		addedIpAddresses.External.Address)
	ctx.Export(getIngressIpOutputName(enums.IpAddressVisibility_IP_ADDRESS_VISIBILITY_INTERNAL, kubeClusterId),
		addedIpAddresses.Internal.Address)
}

func addIpAddresses(ctx *pulumi.Context, input *Input, addedGcpProject *organizations.Project,
	dependencies []pulumi.Resource) (*AddedComputeIpAddresses, error) {
	externalIpAddress, err := addExternalIp(ctx, input, addedGcpProject, dependencies)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add external ingress ip")
	}
	internalIpAddress, err := addInternalIp(ctx, input, addedGcpProject, dependencies)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add internal ingress ip")
	}
	return &AddedComputeIpAddresses{
		External: externalIpAddress,
		Internal: internalIpAddress,
	}, nil
}

func addExternalIp(ctx *pulumi.Context, input *Input, addedGcpProject *organizations.Project,
	dependencies []pulumi.Resource) (*compute.Address, error) {
	ingIpName := getIngressExternalIpName(input.KubeClusterId)
	ca, err := compute.NewAddress(ctx, ingIpName, &compute.AddressArgs{
		Name:        pulumi.String(ingIpName),
		Project:     addedGcpProject.ProjectId,
		Region:      pulumi.String(input.GcpRegion),
		AddressType: pulumi.String(strings.ToUpper(wordpb.Word_external.String())),
		Labels:      pulumi.ToStringMap(input.Labels),
	}, pulumi.Parent(addedGcpProject), pulumi.DependsOn(dependencies))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add new compute address")
	}
	return ca, nil
}

func addInternalIp(ctx *pulumi.Context, input *Input, addedGcpProject *organizations.Project,
	dependencies []pulumi.Resource) (*compute.Address, error) {
	ingIpName := getIngressInternalIpName(input.KubeClusterId)
	ca, err := compute.NewAddress(ctx, ingIpName, &compute.AddressArgs{
		Name:        pulumi.String(ingIpName),
		Project:     addedGcpProject.ProjectId,
		Region:      pulumi.String(input.GcpRegion),
		AddressType: pulumi.String(strings.ToUpper(wordpb.Word_internal.String())),
		Subnetwork:  input.AddedSubnet.SelfLink,
		Labels:      pulumi.ToStringMap(input.Labels),
	}, pulumi.Parent(input.AddedSubnet), pulumi.DependsOn(dependencies))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add new compute address")
	}
	return ca, nil
}

func getIngressExternalIpName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s-ingress-ip", kubeClusterId, wordpb.Word_external.String())
}

func getIngressInternalIpName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s-ingress-ip", kubeClusterId, wordpb.Word_internal.String())
}

func getIngressIpOutputName(visibility enums.IpAddressVisibility, kubeClusterId string) string {
	switch visibility {
	case enums.IpAddressVisibility_IP_ADDRESS_VISIBILITY_EXTERNAL:
		return GetIngressExternalIpOutputName(kubeClusterId)
	case enums.IpAddressVisibility_IP_ADDRESS_VISIBILITY_INTERNAL:
		return GetIngressInternalIpOutputName(kubeClusterId)
	}
	return ""
}

func GetIngressExternalIpOutputName(kubeClusterId string) string {
	return puluminamegcpoutput.Name(compute.Address{},
		fmt.Sprintf("%s-%s-ingress-ip", kubeClusterId, wordpb.Word_external.String()))
}

func GetIngressInternalIpOutputName(kubeClusterId string) string {
	return puluminamegcpoutput.Name(compute.Address{},
		fmt.Sprintf("%s-%s-ingress-ip", kubeClusterId, wordpb.Word_internal.String()))
}
