package dns

import (
	"buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/english/rpc/enums"
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/cloud/gcp/iam/roles/standard"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	AddedWorkloadDeployerGsa     *serviceaccount.Account
	AddedCertManagerGsa          *serviceaccount.Account
	AddedContainerClusterProject *organizations.Project
}

// Resources grants workload-deployer and cert-manager service accounts dns-admin role in kube-cluster project.
func Resources(ctx *pulumi.Context, input *Input) error {
	_, err := projects.NewIAMBinding(ctx, fmt.Sprintf("%s-dns-admin", enums.Word_share.String()),
		&projects.IAMBindingArgs{
			Members: pulumi.StringArray{
				pulumi.Sprintf("serviceAccount:%s", input.AddedWorkloadDeployerGsa.Email),
				pulumi.Sprintf("serviceAccount:%s", input.AddedCertManagerGsa.Email),
			},
			Project: input.AddedContainerClusterProject.ProjectId,
			Role:    pulumi.String(standard.Dns_admin),
		}, pulumi.Parent(input.AddedContainerClusterProject))
	if err != nil {
		return errors.Wrapf(err, "failed to add %s project role binding in share project", standard.Dns_admin)
	}
	return nil
}
