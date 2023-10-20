package cleanup

import (
	stackrpc "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/stack/enums"
)

type Input struct {
	SourceKubeconfigBase64 string
	WorkspaceDir           string
	StackOperationType     stackrpc.StackOperationType
}
