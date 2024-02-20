package folder

import (
	"fmt"

	"github.com/pkg/errors"
	puluminamegcpoutput "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/provider/cloud/gcp/output"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
	pulumigcp "github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/organizations"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Input struct {
	GcpProvider          *pulumigcp.Provider
	CloudAccountFolderId string
	KubeClusterId        string
}

// Resources adds a folder to organize all projects created by planton cloud for the kube-cluster
// the parent for this folder is the cloud account folder created by planton cloud.
func Resources(ctx *pulumi.Context, input *Input) (*organizations.Folder, error) {
	randomString, err := random.NewRandomString(ctx, fmt.Sprintf("%s-folder-suffix", input.KubeClusterId), &random.RandomStringArgs{
		Special: pulumi.Bool(false),
		Lower:   pulumi.Bool(true),
		Upper:   pulumi.Bool(false),
		Number:  pulumi.Bool(true),
		Length:  pulumi.Int(2), //increasing this can result in violation of folder display name length <30
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create random suffix for cloud account folder")
	}
	kubeClusterFolder, err := organizations.NewFolder(ctx, input.KubeClusterId,
		&organizations.FolderArgs{
			DisplayName: pulumi.Sprintf("%s-%s", input.KubeClusterId, randomString.Result),
			Parent:      pulumi.Sprintf("folders/%s", input.CloudAccountFolderId),
		}, pulumi.Provider(input.GcpProvider))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to add cloud account folder")
	}
	ctx.Export(GetKubeClusterFolderIdOutputName(input.KubeClusterId), kubeClusterFolder.FolderId)
	ctx.Export(GetKubeClusterFolderDisplayNameOutputName(input.KubeClusterId), kubeClusterFolder.DisplayName)
	ctx.Export(GetKubeClusterFolderParentOutputName(input.KubeClusterId), kubeClusterFolder.Parent)
	return kubeClusterFolder, nil
}

func GetKubeClusterFolderIdOutputName(kubeClusterId string) string {
	return puluminamegcpoutput.Name(&organizations.Folder{}, kubeClusterId, englishword.EnglishWord_id.String())
}

func GetKubeClusterFolderDisplayNameOutputName(kubeClusterId string) string {
	return puluminamegcpoutput.Name(&organizations.Folder{}, kubeClusterId, englishword.EnglishWord_name.String())
}

func GetKubeClusterFolderParentOutputName(kubeClusterId string) string {
	return puluminamegcpoutput.Name(&organizations.Folder{}, kubeClusterId, englishword.EnglishWord_parent.String())
}
