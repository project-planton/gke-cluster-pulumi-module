package controller

import (
	"github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/commons/english/enums/englishword"
)

var (
	SelectorLabels = map[string]string{
		englishword.EnglishWord_app.String(): "istio-ingress",
		"istio":                              "ingress",
	}
)
