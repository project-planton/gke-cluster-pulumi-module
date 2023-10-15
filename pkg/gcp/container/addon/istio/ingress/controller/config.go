package controller

import (
	wordpb "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/commons/english/rpc/enums"
)

var (
	SelectorLabels = map[string]string{
		wordpb.Word_app.String(): "istio-ingress",
		"istio":                  "ingress",
	}
)
