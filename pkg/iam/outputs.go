package iam

import (
	"github.com/plantoncloud/kube-cluster-pulumi-blueprint/pkg/gcp/iam/certmanager"
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/code2cloud/v1/gcp/gkecluster/model"
	"github.com/plantoncloud/stack-job-runner-golang-sdk/pkg/stack/output/backend"
)

func Output(input *model.GkeClusterStackResourceInput,
	stackOutput map[string]interface{}) *model.GkeClusterStackIamOutputs {
	return &model.GkeClusterStackIamOutputs{
		CertManagerGsaEmail: backend.GetVal(stackOutput, certmanager.GetGsaEmailOutputName)),
		ExternalSecretsGsaEmail:      backend.GetVal(stackOutput, externalsecrets.GetGsaEmailOutputName     )),
		ExternalDnsGsaEmail:          backend.GetVal(stackOutput, externaldns.GetGsaEmailOutputName     )),
		WorkloadDeployerGsaEmail:     backend.GetVal(stackOutput, workloaddeployer.GetGsaEmailOutputName     )),
		WorkloadDeployerGsaKeyBase64: backend.GetVal(stackOutput, workloaddeployer.GetGsaKeyOutputName     )),
	}
}
