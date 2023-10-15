package iam

import (
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/iam/certmanager"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/iam/externalsecrets"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/iam/workloaddeployer"
	"github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/stack/output/backend"
	kubernetesclustergcpstack "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
)

func Output(input *kubernetesclustergcpstack.KubeClusterGcpStackResourceInput, stackOutput map[string]interface{}) *kubernetesclustergcpstack.KubeClusterGcpStackIamOutputs {
	return &kubernetesclustergcpstack.KubeClusterGcpStackIamOutputs{
		CertManagerGsaEmail:          backend.GetVal(stackOutput, certmanager.GetGsaEmailOutputName()),
		ExternalSecretsGsaEmail:      backend.GetVal(stackOutput, externalsecrets.GetGsaEmailOutputName()),
		WorkloadDeployerGsaEmail:     backend.GetVal(stackOutput, workloaddeployer.GetGsaEmailOutputName()),
		WorkloadDeployerGsaKeyBase64: backend.GetVal(stackOutput, workloaddeployer.GetGsaKeyOutputName()),
	}
}
