package cleanup

import (
	stackrpc "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/stack/rpc/enums"
)

type Input struct {
	SourceKubeconfigBase64 string
	ReqWorkspace           string
	StackOperationType     stackrpc.StackOperationType
}
