package subnet

import (
	"fmt"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"
	"strings"

	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/enums/gkepodservicesecondaryrangecidrsetnum"

	"github.com/pkg/errors"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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

func GetSubNetworkSelfLinkOutputName     subNetworkName string) string {
	return pulumigoogleprovider.PulumiOutputName
	compute.Subnetwork{}, subNetworkName)
}
