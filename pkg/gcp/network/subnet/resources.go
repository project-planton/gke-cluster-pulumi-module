package subnet

import (
	"fmt"
	"github.com/pkg/errors"
	puluminameoutputgcp "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/rpc/enums"
	envgcpnetworkapi "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	rpc "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/commons/english/rpc/enums"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

//https://jodies.de/ipcalc?host=10.0.0.0&mask1=10&mask2=14

//10.0.0.0/14 // primary subnetwork cidr
//10.32.0.0/14 //primary subnetwork cidr(reserve)

//pod range
//10.4.0.0/14
//10.8.0.0/14
//10.12.0.0/14
//10.16.0.0/14
//10.20.0.0/14
//10.24.0.0/14
//10.28.0.0/14
//service range
//10.36.0.0/14
//10.40.0.0/14
//10.44.0.0/14
//10.48.0.0/14
//10.52.0.0/14
//10.56.0.0/14
//10.60.0.0/14

const (
	subNetworkCidr                     = "10.0.0.0/14"
	secondaryIpRangeNamePodsPrefix     = "gke-pods"
	secondaryIpRangeNameServicesPrefix = "gke-services"
)

var (
	podCidrSecondaryRangeMap = map[enums.GkeKubePodServiceSecondaryRangeCidrSetNum]*envgcpnetworkapi.KubePodServiceSecondaryRangeCidr{
		enums.GkeKubePodServiceSecondaryRangeCidrSetNum_ONE: {
			Pod:     "10.4.0.0/14",
			Service: "10.36.0.0/14",
		}, enums.GkeKubePodServiceSecondaryRangeCidrSetNum_TWO: {
			Pod:     "10.8.0.0/14",
			Service: "10.40.0.0/14",
		}, enums.GkeKubePodServiceSecondaryRangeCidrSetNum_THREE: {
			Pod:     "10.12.0.0/14",
			Service: "10.44.0.0/14",
		}, enums.GkeKubePodServiceSecondaryRangeCidrSetNum_FOUR: {
			Pod:     "10.16.0.0/14",
			Service: "10.48.0.0/14",
		}, enums.GkeKubePodServiceSecondaryRangeCidrSetNum_FIVE: {
			Pod:     "10.20.0.0/14",
			Service: "10.52.0.0/14",
		}, enums.GkeKubePodServiceSecondaryRangeCidrSetNum_SIX: {
			Pod:     "10.24.0.0/14",
			Service: "10.56.0.0/14",
		}, enums.GkeKubePodServiceSecondaryRangeCidrSetNum_SEVEN: {
			Pod:     "10.28.0.0/14",
			Service: "10.60.0.0/14",
		},
	}
)

type Input struct {
	KubeClusterId string
	GcpRegion     string
	ShareProject  *organizations.Project
	VpcNetwork    *compute.Network
}

func Resources(ctx *pulumi.Context, input *Input) (*compute.Subnetwork, error) {
	snw, err := addSubNetwork(ctx, input)
	if err != nil {
		return nil, err
	}
	return snw, nil
}

func addSubNetwork(ctx *pulumi.Context, input *Input) (*compute.Subnetwork, error) {
	name := GetSubNetworkName(input.KubeClusterId)
	sn, err := compute.NewSubnetwork(ctx, name, &compute.SubnetworkArgs{
		Name:                  pulumi.String(name),
		Project:               input.ShareProject.ProjectId,
		Network:               input.VpcNetwork.ID(),
		Region:                pulumi.String(input.GcpRegion),
		IpCidrRange:           pulumi.String(subNetworkCidr),
		PrivateIpGoogleAccess: pulumi.BoolPtr(true),
		SecondaryIpRanges: &compute.SubnetworkSecondaryIpRangeArray{
			getPodSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_ONE),
			getServiceSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_ONE),
			getPodSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_TWO),
			getServiceSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_TWO),
			getPodSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_THREE),
			getServiceSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_THREE),
			getPodSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_FOUR),
			getServiceSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_FOUR),
			getPodSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_FIVE),
			getServiceSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_FIVE),
			getPodSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_SIX),
			getServiceSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_SIX),
			getPodSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_SEVEN),
			getServiceSecondaryRanges(enums.GkeKubePodServiceSecondaryRangeCidrSetNum_SEVEN),
		},
	}, pulumi.Parent(input.VpcNetwork))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add subnetwork")
	}
	ctx.Export(GetSubNetworkSelfLinkOutputName(name), sn.SelfLink)
	return sn, nil
}

// todo: this is a suboptimal code as a workaround for in ability to create an pulumi input array with looping
func getPodSecondaryRanges(setNum enums.GkeKubePodServiceSecondaryRangeCidrSetNum) *compute.SubnetworkSecondaryIpRangeArgs {
	rangeSet := podCidrSecondaryRangeMap[setNum]
	return &compute.SubnetworkSecondaryIpRangeArgs{
		RangeName:   pulumi.String(GetPodsSecondaryRangeName(setNum)),
		IpCidrRange: pulumi.String(rangeSet.Pod),
	}
}

// todo: this is a suboptimal code as a workaround for in ability to create an pulumi input array with looping
func getServiceSecondaryRanges(setNum enums.GkeKubePodServiceSecondaryRangeCidrSetNum) *compute.SubnetworkSecondaryIpRangeArgs {
	rangeSet := podCidrSecondaryRangeMap[setNum]
	return &compute.SubnetworkSecondaryIpRangeArgs{
		RangeName:   pulumi.String(GetServicesSecondaryRangeName(setNum)),
		IpCidrRange: pulumi.String(rangeSet.Service),
	}
}

func GetSubNetworkName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s", rpc.Word_kubernetes, kubeClusterId)
}

func GetServicesSecondaryRangeName(setNum enums.GkeKubePodServiceSecondaryRangeCidrSetNum) string {
	return fmt.Sprintf("%s-%s", secondaryIpRangeNameServicesPrefix, strings.ToLower(setNum.String()))
}

func GetPodsSecondaryRangeName(setNum enums.GkeKubePodServiceSecondaryRangeCidrSetNum) string {
	return fmt.Sprintf("%s-%s", secondaryIpRangeNamePodsPrefix, strings.ToLower(setNum.String()))
}

func GetSubNetworkSelfLinkOutputName(subNetworkName string) string {
	return puluminameoutputgcp.Name(compute.Subnetwork{}, subNetworkName)
}
