package router

import (
	rpc "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/english/rpc/enums"
	"fmt"
	"github.com/pkg/errors"
	puluminameoutputgcp "github.com/plantoncloud-inc/pulumi-stack-runner-sdk/go/pulumi/name/provider/cloud/gcp/output"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
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
	ctx.Export(GetRouterSelfLinkOutputName(name), nr.SelfLink)
	return nr, nil
}

func GetRouterName(kubeClusterId string) string {
	return fmt.Sprintf("%s-%s-router", rpc.Word_kubernetes, kubeClusterId)
}

func GetRouterSelfLinkOutputName(routerName string) string {
	return puluminameoutputgcp.Name(compute.Router{}, routerName)
}
