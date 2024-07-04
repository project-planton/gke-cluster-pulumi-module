package iam

import (
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam/certmanager"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam/externaldns"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam/externalsecrets"
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam/workloaddeployer"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/kubecluster/stack/gcp/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/stack/output/backend"
)

func Output(input *c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackResourceInput,
	stackOutput map[string]interface{}) *c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackIamOutputs {
	return &c2cv1deployk8cstackgcpmodel.KubeClusterGcpStackIamOutputs{
		CertManagerGsaEmail:          backend.GetVal(stackOutput, certmanager.GetGsaEmailOutputName()),
		ExternalSecretsGsaEmail:      backend.GetVal(stackOutput, externalsecrets.GetGsaEmailOutputName()),
		ExternalDnsGsaEmail:          backend.GetVal(stackOutput, externaldns.GetGsaEmailOutputName()),
		WorkloadDeployerGsaEmail:     backend.GetVal(stackOutput, workloaddeployer.GetGsaEmailOutputName()),
		WorkloadDeployerGsaKeyBase64: backend.GetVal(stackOutput, workloaddeployer.GetGsaKeyOutputName()),
	}
}
