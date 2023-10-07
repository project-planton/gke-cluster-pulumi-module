package iam

import (
	kubernetesclustergcpstack "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/iam/certmanager"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/iam/externalsecrets"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/iam/workloaddeployer"
	"github.com/plantoncloud-inc/pulumi-stack-runner-sdk/go/pulumi/stack/output/backend"
)

func Output(input *kubernetesclustergcpstack.KubeClusterGcpStackResourceInput, stackOutput map[string]interface{}) *kubernetesclustergcpstack.KubeClusterGcpStackIamOutputs {
	return &kubernetesclustergcpstack.KubeClusterGcpStackIamOutputs{
		CertManagerGsaEmail:          backend.GetVal(stackOutput, certmanager.GetGsaEmailOutputName()),
		ExternalSecretsGsaEmail:      backend.GetVal(stackOutput, externalsecrets.GetGsaEmailOutputName()),
		WorkloadDeployerGsaEmail:     backend.GetVal(stackOutput, workloaddeployer.GetGsaEmailOutputName()),
		WorkloadDeployerGsaKeyBase64: backend.GetVal(stackOutput, workloaddeployer.GetGsaKeyOutputName()),
	}
}
