package controller

import (
	wordpb "buf.build/gen/go/plantoncloud/planton-cloud-apis/protocolbuffers/go/cloud/planton/apis/v1/commons/english/rpc/enums"
)

var (
	SelectorLabels = map[string]string{
		wordpb.Word_app.String(): "istio-ingress",
		"istio":                  "ingress",
	}
)
