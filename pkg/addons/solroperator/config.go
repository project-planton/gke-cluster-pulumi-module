package solroperator

import (
	"fmt"
	"github.com/plantoncloud-inc/go-commons/kubernetes/helm"
)

const (
	Namespace = "solr-operator"
	Version   = "0.7.0"
)

type HelmVal struct {
}

func getHelmVal() *HelmVal {
	return &HelmVal{}
}

// getHelmChart https://github.com/apache/solr-operator/tree/main/helm/solr-operator#installing-the-chart
func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "solr-operator",
		Repo:        "https://solr.apache.org/charts",
		Name:        "solr-operator",
		Version:     Version,
	}
}

func getCrdsDownloadUrl() string {
	return fmt.Sprintf("https://solr.apache.org/operator/downloads/crds/v%s/all-with-dependencies.yaml", Version)
}
