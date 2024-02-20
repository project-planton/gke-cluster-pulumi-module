package vpc

import (
	"github.com/pkg/errors"
	code2cloudv1deployk8cmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/model"
	awsclassic "github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/zalando/postgres-operator/pkg/util/k8sutil"
)

const DefaultCidrBlock = "10.0.0.0/16"

type Input struct {
	AwsProvider *awsclassic.Provider
	KubeCluster *code2cloudv1deployk8cmodel.KubeCluster
	Labels      map[string]string
}

func Resources(ctx *pulumi.Context, input *Input) (*ec2.Vpc, error) {
	addedVpc, err := ec2.NewVpc(ctx, input.KubeCluster.Metadata.Id, &ec2.VpcArgs{
		CidrBlock:          k8sutil.StringToPointer(DefaultCidrBlock),
		EnableDnsHostnames: pulumi.Bool(true),
		EnableDnsSupport:   pulumi.Bool(true),
		Tags:               pulumi.ToStringMap(input.Labels),
	}, pulumi.Provider(input.AwsProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create addedVpc")
	}
	return addedVpc, nil
}
