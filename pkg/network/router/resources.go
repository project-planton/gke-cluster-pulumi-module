package router

import (
	"fmt"
	"github.com/plantoncloud/pulumi-blueprint-golang-commons/pkg/google/pulumigoogleprovider"

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
	VpcNetwork             *compute.Network
}

func Resources(ctx *pulumi.Context, input *Input) (*compute.Router, error) {
	name := GetRouterName(input.KubeClusterId)
	nr, err := compute.NewRouter(ctx, name, &compute.RouterArgs{
		Name:    pulumi.String(name),
		Network: input.VpcNetwork.SelfLink,
		Region:  pulumi.String(input.GcpRegion),
		Project: input.AddedVpcNetworkProject.ProjectId,
	}, pulumi.Parent(input.VpcNetwork))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add compute router")
	}
	ctx.Export(RouterSelfLinkOutputName
	name), nr.SelfLink)
	return nr, nil
}

func GetRouterName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s-router", englishword.EnglishWord_kubernetes, kubeClusterId)
}
