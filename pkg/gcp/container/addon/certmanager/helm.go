package certmanager

import (
	helmcommons "github.com/plantoncloud-inc/go-commons/kubernetes/helm"
)

// helmVal https://artifacthub.io/packages/helm/cert-manager/cert-manager?modal=values
type helmVal struct {
	InstallCrds     bool             `yaml:"installCRDs"`
	ExtraArgs       []string         `yaml:"extraArgs"`
	ServiceAccount  *serviceAccount  `yaml:"serviceAccount"`
	StartupApiCheck *StartupApiCheck `yaml:"startupapicheck"`
	//comma separated list of feature gates that should be enabled
	FeatureGates string          `yaml:"featureGates"`
	Webhook      *helmValWebhook `yaml:"webhook"`
}

type helmValWebhook struct {
	ExtraArgs []string `yaml:"extraArgs"`
}

type serviceAccount struct {
	Create bool   `yaml:"create"`
	Name   string `yaml:"name"`
}

type StartupApiCheck struct {
	Enabled bool   `yaml:"enabled"`
	Timeout string `yaml:"timeout"`
}

// getHelmChart https://cert-manager.io/docs/installation/helm/
func getHelmChart() *helmcommons.Chart {
	return &helmcommons.Chart{
		ReleaseName: "cert-manager",
		Repo:        "https://charts.jetstack.io",
		Name:        "cert-manager",
		//https://github.com/cert-manager/cert-manager/releases/tag/v1.12.2
		Version: "v1.12.2",
	}
}

func getHelmVal() *helmVal {
	return &helmVal{
		InstallCrds: true,
		ExtraArgs: []string{
			"--dns01-recursive-nameservers=\"1.1.1.1:53\"",
			"--dns01-recursive-nameservers-only=true",
		},
		ServiceAccount: &serviceAccount{
			Create: false,
			Name:   Ksa,
		},
		StartupApiCheck: &StartupApiCheck{
			Enabled: true,
			//https://github.com/cert-manager/cert-manager/issues/4646#issuecomment-1023741490
			Timeout: "5m",
		},
		//feature gate required to create a single ca.pem file required for mounting the self-signed cert on postgres stunnel sidecar container.
		FeatureGates: "AdditionalCertificateOutputFormats=true",
		Webhook:      &helmValWebhook{ExtraArgs: []string{"--feature-gates=AdditionalCertificateOutputFormats=true"}},
	}
}
