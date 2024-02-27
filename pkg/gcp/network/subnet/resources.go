package subnet

import (
	"fmt"
	"strings"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/enums/gkepodservicesecondaryrangecidrsetnum"

	"github.com/pkg/errors"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/gcp/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	puluminameoutputgcp "github.com/plantoncloud/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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
	podCidrSecondaryRangeMap = map[gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum]*c2cv1deployk8cstackgcpmodel.KubePodServiceSecondaryRangeCidr{
		gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_one: {
			Pod:     "10.4.0.0/14",
			Service: "10.36.0.0/14",
		}, gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_two: {
			Pod:     "10.8.0.0/14",
			Service: "10.40.0.0/14",
		}, gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_three: {
			Pod:     "10.12.0.0/14",
			Service: "10.44.0.0/14",
		}, gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_four: {
			Pod:     "10.16.0.0/14",
			Service: "10.48.0.0/14",
		}, gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_five: {
			Pod:     "10.20.0.0/14",
			Service: "10.52.0.0/14",
		}, gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_six: {
			Pod:     "10.24.0.0/14",
			Service: "10.56.0.0/14",
		}, gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_seven: {
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
			getPodSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_one),
			getServiceSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_one),
			getPodSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_two),
			getServiceSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_two),
			getPodSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_three),
			getServiceSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_three),
			getPodSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_four),
			getServiceSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_four),
			getPodSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_five),
			getServiceSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_five),
			getPodSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_six),
			getServiceSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_six),
			getPodSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_seven),
			getServiceSecondaryRanges(gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum_seven),
		},
	}, pulumi.Parent(input.VpcNetwork))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add subnetwork")
	}
	ctx.Export(GetSubNetworkSelfLinkOutputName(name), sn.SelfLink)
	return sn, nil
}

// todo: this is a suboptimal code as a workaround for in ability to create an pulumi input array with looping
func getPodSecondaryRanges(setNum gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum) *compute.SubnetworkSecondaryIpRangeArgs {
	rangeSet := podCidrSecondaryRangeMap[setNum]
	return &compute.SubnetworkSecondaryIpRangeArgs{
		RangeName:   pulumi.String(GetPodsSecondaryRangeName(setNum)),
		IpCidrRange: pulumi.String(rangeSet.Pod),
	}
}

// todo: this is a suboptimal code as a workaround for in ability to create an pulumi input array with looping
func getServiceSecondaryRanges(setNum gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum) *compute.SubnetworkSecondaryIpRangeArgs {
	rangeSet := podCidrSecondaryRangeMap[setNum]
	return &compute.SubnetworkSecondaryIpRangeArgs{
		RangeName:   pulumi.String(GetServicesSecondaryRangeName(setNum)),
		IpCidrRange: pulumi.String(rangeSet.Service),
	}
}

func GetSubNetworkName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s", englishword.EnglishWord_kubernetes, kubeClusterId)
}

func GetServicesSecondaryRangeName(setNum gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum) string {
	return fmt.Sprintf("%s-%s", secondaryIpRangeNameServicesPrefix, strings.ToLower(setNum.String()))
}

func GetPodsSecondaryRangeName(setNum gkepodservicesecondaryrangecidrsetnum.GkePodServiceSecondaryRangeCidrSetNum) string {
	return fmt.Sprintf("%s-%s", secondaryIpRangeNamePodsPrefix, strings.ToLower(setNum.String()))
}

func GetSubNetworkSelfLinkOutputName(subNetworkName string) string {
	return puluminameoutputgcp.Name(compute.Subnetwork{}, subNetworkName)
}
