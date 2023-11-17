package cleanup

import "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/stack/job/enums/operationtype"

type Input struct {
	SourceKubeconfigBase64 string
	WorkspaceDir           string
	StackJobOperationType  operationtype.StackJobOperationType
}
