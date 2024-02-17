package prometheus

import (
	"github.com/plantoncloud-inc/go-commons/kubernetes/helm"
	c2cv1deployk8cstackgcpmodel "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/code2cloud/deploy/kubecluster/stack/gcp/model"
)

const (
	Namespace       = "prometheus"
	KubeServiceName = "prometheus-server"
)

type HelmVal struct {
	AlertManager           *ValAlertManager           `yaml:"alertmanager"`
	KubeStateMetrics       *ValKubeStateMetrics       `yaml:"kube-state-metrics"`
	PrometheusNodeExporter *ValPrometheusNodeExporter `yaml:"prometheus-node-exporter"`
	PrometheusPushGateway  *ValPrometheusPushGateway  `yaml:"prometheus-pushgateway"`
	PrometheusServer       *ValPrometheusServer       `yaml:"server"`
}

type ValAlertManager struct {
	Enabled bool `yaml:"enabled"`
}

type ValKubeStateMetrics struct {
	Enabled bool `yaml:"enabled"`
}

type ValPrometheusNodeExporter struct {
	Enabled bool `yaml:"enabled"`
}

type ValPrometheusPushGateway struct {
	Enabled bool `yaml:"enabled"`
}

// ValPrometheusServer is prometheus-server configuration
type ValPrometheusServer struct {
	Enabled          bool                       `yaml:"enabled"`
	PersistentVolume *ValServerPersistentVolume `yaml:"persistentVolume"`
}

// ValServerPersistentVolume is persistent volume for prometheus-server configuration
type ValServerPersistentVolume struct {
	Enabled bool   `yaml:"enabled"`
	Size    string `yaml:"size"`
}

//https://github.com/prometheus-community/helm-charts/blob/main/charts/prometheus/values.yaml
/*
alertmanager:
	enabled: false
kube-state-metrics:
	enabled: false
prometheus-node-exporter:
	enabled: false
prometheus-pushgateway:
	enabled: false
server:
	persistentVolume:
	  enabled: true
	  size: 2Gi
*/
func getHelmVal(openCostInput *c2cv1deployk8cstackgcpmodel.AddonsOpenCost) *HelmVal {
	helmVal := &HelmVal{
		AlertManager:           &ValAlertManager{Enabled: false},
		KubeStateMetrics:       &ValKubeStateMetrics{Enabled: false},
		PrometheusNodeExporter: &ValPrometheusNodeExporter{Enabled: false},
		PrometheusPushGateway:  &ValPrometheusPushGateway{Enabled: false},
		PrometheusServer: &ValPrometheusServer{
			Enabled: true,
			PersistentVolume: &ValServerPersistentVolume{
				Enabled: true,
				Size:    "2Gi",
			},
		},
	}
	if openCostInput.PrometheusDataDiskSizeGb != "" {
		helmVal.PrometheusServer.PersistentVolume.Size = openCostInput.PrometheusDataDiskSizeGb
	}
	return helmVal
}

// helm install prometheus prometheus-community/prometheus -n prometheus --values val.yaml
func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "prometheus",
		Repo:        "https://prometheus-community.github.io/helm-charts",
		Name:        "prometheus",
		//https://artifacthub.io/packages/helm/prometheus-community/prometheus
		Version: "20.2.0",
	}
}
