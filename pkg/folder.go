package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/opname"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/connect/v1/gcpcredential/model"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (s *ResourceStack) folder(ctx *pulumi.Context, gcpProvider *gcp.Provider, gcpCredential *model.GcpCredential) (*organizations.Folder, error) {
	//create variable with descriptive name for the api-resource in the input
	gkeCluster := s.Input.ApiResource

	//create a random suffix to be used for naming the folder
	//random suffix is to ensure uniqueness of the folder name on google cloud
	randomString, err := random.NewRandomString(ctx, "folder-suffix", &random.RandomStringArgs{
		Special: pulumi.Bool(false),
		Lower:   pulumi.Bool(true),
		Upper:   pulumi.Bool(false),
		Numeric: pulumi.Bool(true),
		Length:  pulumi.Int(2), //increasing this can result in violation of folder display name length <30
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create random suffix for folder-name")
	}

	//create google cloud folder with the organization from gcp-creadential
	createdFolder, err := organizations.NewFolder(ctx, "folder",
		&organizations.FolderArgs{
			DisplayName: pulumi.Sprintf("%s-%s", gkeCluster.Metadata.Id, randomString.Result),
			Parent:      pulumi.Sprintf("organizations/%s", gcpCredential.Spec.GcpOrganizationId),
		}, pulumi.Provider(gcpProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add cloud account folder")
	}

	//export important attributes of the created folder
	ctx.Export(opname.FolderId, createdFolder.FolderId)
	ctx.Export(opname.FolderDisplayName, createdFolder.DisplayName)
	ctx.Export(opname.FolderParent, createdFolder.Parent)

	return createdFolder, nil
}
