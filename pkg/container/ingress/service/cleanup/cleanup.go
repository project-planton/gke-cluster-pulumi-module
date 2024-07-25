package cleanup

import "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/iac/v1/stackjob/enums/stackjoboperationtype"

type Input struct {
	SourceKubeconfigBase64 string
	WorkspaceDir           string
	StackJobOperationType  stackjoboperationtype.StackJobOperationType
}
