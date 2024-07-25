package traefik

import "github.com/plantoncloud-inc/go-commons/kubernetes/helm"

// HelmVal https://github.com/traefik/traefik-helm-chart/blob/master/traefik/values.yaml
type HelmVal struct {
	Rbac                *Rbac                             `yaml:"rbac"`
	Providers           *Providers                        `yaml:"providers"`
	Logs                *Logs                             `yaml:"logs"`
	Ports               map[string]map[string]interface{} `yaml:"ports"`
	AdditionalArguments []string                          `yaml:"additionalArguments"`
	IngressRoute        *IngressRoute                     `yaml:"ingressRoute"`
	Service             *Service                          `yaml:"service"`
}

type Rbac struct {
	Enabled    bool `yaml:"enabled"`
	Namespaced bool `yaml:"namespaced"`
}

type KubernetesCRD struct {
	Enabled             bool `yaml:"enabled"`
	AllowCrossNamespace bool `yaml:"allowCrossNamespace"`
}
type KubernetesIngress struct {
	Enabled                   bool `yaml:"enabled"`
	AllowExternalNameServices bool `yaml:"allowExternalNameServices"`
}
type Providers struct {
	KubernetesCRD     *KubernetesCRD     `yaml:"kubernetesCRD"`
	KubernetesIngress *KubernetesIngress `yaml:"kubernetesIngress"`
}
type General struct {
	Level string `yaml:"level"`
}
type Names struct {
	ClientUsername string `yaml:"ClientUsername"`
}
type Fields struct {
	DefaultMode string `yaml:"defaultMode"`
	Names       *Names `yaml:"names"`
}
type Access struct {
	Enabled bool    `yaml:"enabled"`
	Fields  *Fields `yaml:"fields"`
}
type Logs struct {
	General *General `yaml:"general"`
	Access  *Access  `yaml:"access"`
}
type Metrics struct {
	Port        int  `yaml:"port"`
	Expose      bool `yaml:"expose"`
	ExposedPort int  `yaml:"exposedPort"`
}
type Dashboard struct {
	Enabled bool `yaml:"enabled"`
}
type IngressRoute struct {
	Dashboard *Dashboard `yaml:"dashboard"`
}
type Annotations struct {
	PrometheusIoPort   string `yaml:"prometheus.io/port"`
	PrometheusIoScrape string `yaml:"prometheus.io/scrape"`
}
type Spec struct {
	ExternalTrafficPolicy string `yaml:"externalTrafficPolicy"`
	LoadBalancerIp        string `yaml:"loadBalancerIP"`
}
type Service struct {
	Enabled     bool         `yaml:"enabled"`
	Type        string       `yaml:"type"`
	Annotations *Annotations `yaml:"annotations"`
	Spec        *Spec        `yaml:"spec"`
}

func getHelmVal() *HelmVal {
	cv := &HelmVal{
		Rbac: &Rbac{
			Enabled:    true,
			Namespaced: false,
		},
		Providers: &Providers{
			KubernetesCRD: &KubernetesCRD{
				Enabled:             true,
				AllowCrossNamespace: true,
			},
			KubernetesIngress: &KubernetesIngress{
				Enabled:                   true,
				AllowExternalNameServices: true,
			},
		},
		Logs: &Logs{
			General: &General{Level: "info"},
			Access: &Access{
				Enabled: true,
				Fields: &Fields{
					DefaultMode: "keep",
					Names:       &Names{ClientUsername: "drop"},
				},
			},
		},
		Ports: getDefaultPorts(),
		AdditionalArguments: []string{
			"--metrics.prometheus=true",
			"--metrics.prometheus.buckets=0.100000, 0.300000, 1.200000, 5.000000",
			"--metrics.prometheus.addEntryPointsLabels=true",
			"--metrics.prometheus.addServicesLabels=true",
			"--metrics.prometheus.entryPoint=metrics",
			"--providers.kubernetescrd.allowcrossnamespace=true",
			"--providers.kubernetescrd.allowexternalnameservices=true",
			"--providers.kubernetesingress.allowexternalnameservices=true",
		},
		IngressRoute: &IngressRoute{
			Dashboard: &Dashboard{Enabled: false},
		},
		Service: &Service{
			Enabled: true,
			Type:    "ClusterIP",
			Annotations: &Annotations{
				PrometheusIoPort:   "8080",
				PrometheusIoScrape: "true",
			},
			Spec: &Spec{ExternalTrafficPolicy: "Local"},
		},
	}
	return cv
}

func getHelmChart() *helm.Chart {
	return &helm.Chart{
		ReleaseName: "traefik",
		Repo:        "https://helm.traefik.io/traefik",
		Name:        "traefik",
		Version:     "10.24.0",
	}
}

func getDefaultPorts() map[string]map[string]interface{} {
	ports := make(map[string]map[string]interface{}, 0)
	ports["web"] = map[string]interface{}{"redirectTo": "websecure"}
	ports["metrics"] = map[string]interface{}{"port": 8080, "expose": true, "exposedPort": 8080}
	return ports
}

func AppendKafkaPort(ports map[string]map[string]interface{}) map[string]map[string]interface{} {
	ports[KafkaPortName] = map[string]interface{}{"port": KafkaPort}
	return ports
}
