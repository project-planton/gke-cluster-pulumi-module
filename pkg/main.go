package pkg

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/gke-cluster-pulumi-module/pkg/localz"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/gcp/pulumigkekubernetesprovider"
	"github.com/plantoncloud/pulumi-module-golang-commons/pkg/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ResourceStack struct {
	Input            *gkecluster.GkeClusterStackInput
	GcpLabels        map[string]string
	KubernetesLabels map[string]string
}

// Resources function is the pulumi program that deploys GKE cluster along with chosen addons.
//
// Parameters:
// - ctx: The Pulumi context used for defining cloud resources.
//
// Returns:
// - error: An error object if there is any issue during the resource creation.
//
// The function performs the following steps:
// 1. Initializes local variables and configuration from the input.
// 2. Sets up the GCP provider using the provided GCP credentials.
// 3. Creates a GCP folder for organizing the projects.
// 4. Creates the GKE cluster within the specified folder.
// 5. Creates the node pools for the GKE cluster.
// 6. Creates a service account and key for deploying workloads to the cluster.
// 7. If Kubernetes addons are specified, creates a Kubernetes provider for the cluster.
// 8. Installs the specified Kubernetes addons using the created providers.
func (s *ResourceStack) Resources(ctx *pulumi.Context) error {
	locals := localz.Initialize(ctx, s.Input)

	//create gcp-provider using the gcp-credential from input
	gcpProvider, err := pulumigoogleprovider.Get(ctx, s.Input.GcpCredential)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	//create gcp folder
	createdFolder, err := folder(ctx, locals, gcpProvider)
	if err != nil {
		return errors.Wrap(err, "failed to create folder")
	}

	//create cluster
	createdCluster, err := cluster(ctx, locals, createdFolder)
	if err != nil {
		return errors.Wrap(err, "failed to create container cluster")
	}

	//create node-pools
	createdNodePools, err := clusterNodePools(ctx, locals, createdCluster)
	if err != nil {
		return errors.Wrap(err, "failed to create cluster node-pools")
	}

	//create workload-deployer google service account resources
	createdWorkloadDeployerServiceAccountKey, err := workloadDeployer(ctx, createdCluster)
	if err != nil {
		return errors.Wrap(err, "failed to create workload-deployer resources")
	}

	//if kubernetes-addons is nil, nothing more to do
	if locals.GkeCluster.Spec.KubernetesAddons == nil {
		return nil
	}

	//create kubernetes provider for the created cluster
	kubernetesProvider, err := pulumigkekubernetesprovider.GetWithCreatedGkeClusterAndCreatedGsaKey(
		ctx,
		createdWorkloadDeployerServiceAccountKey,
		createdCluster,
		createdNodePools,
		"gke-cluster")
	if err != nil {
		return errors.Wrap(err, "failed to create kubernetes provider")
	}

	//create addons
	if err := clusterAddons(ctx, locals, createdCluster, gcpProvider, kubernetesProvider); err != nil {
		return errors.Wrap(err, "failed to create addons")
	}
	return nil
}
