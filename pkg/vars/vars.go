package vars

var (
	NetworkProjectApis = []string{
		"compute.googleapis.com",
		"container.googleapis.com",
		"dns.googleapis.com",
	}

	ContainerClusterProjectApis = []string{
		"compute.googleapis.com",
		"container.googleapis.com",
		"secretmanager.googleapis.com",
		"dns.googleapis.com",
	}

	WorkloadIdentityKubeAnnotationKey = "iam.gke.io/gcp-service-account"

	// SubNetworkCidr 10.0.0.0/14
	// this subnet will be divided into two equal halves for pod-secondary-ip-range and service-secondary-ip-range
	//https://jodies.de/ipcalc?host=10.0.0.0&mask1=14&mask2=15
	SubNetworkCidr = "10.0.0.0/14"

	// KubernetesPodSecondaryIpRange https://cloud.google.com/kubernetes-engine/docs/concepts/alias-ips#cluster_sizing_secondary_range_pods
	KubernetesPodSecondaryIpRange = "10.0.0.0/15"
	// KubernetesServiceSecondaryIpRange https://cloud.google.com/kubernetes-engine/docs/concepts/alias-ips#cluster_sizing_secondary_range_svcs
	KubernetesServiceSecondaryIpRange = "10.2.0.0/15"

	ApiServerIpCidr                                     = "172.16.0.0/24"
	ClusterMasterAuthorizedNetworksCidrBlock            = "0.0.0.0/0"
	ClusterMasterAuthorizedNetworksCidrBlockDescription = "kubectl-from-anywhere"
	ApiServerWebhookPort                                = "8443"
	IstioPilotWebhookPort                               = "15017"

	// WorkloadDeployServiceAccountName name of the google service account to
	//be used for deploying workloads to the gke cluster.
	WorkloadDeployServiceAccountName = "workload-deployer"

	CertManager = struct {
		HelmChartName        string
		HelmChartRepo        string
		HelmChartVersion     string
		KsaName              string
		Namespace            string
		SelfSignedIssuerName string
	}{
		"cert-manager",
		"https://charts.jetstack.io",
		"v1.15.1",
		"cert-manager",
		"cert-manager",
		"self-signed",
	}

	ExternalSecrets = struct {
		HelmChartName                           string
		HelmChartRepo                           string
		HelmChartVersion                        string
		KsaName                                 string
		Namespace                               string
		SecretsPollingIntervalSeconds           int
		GcpSecretsManagerClusterSecretStoreName string
	}{
		"external-secrets",
		"https://charts.external-secrets.io",
		"v0.9.20",
		"external-secrets",
		"external-secrets",
		//caution: polling interval frequency may have effect on provider costs on some platforms
		10,
		"gcp-secrets-manager",
	}

	IngressNginx = struct {
		HelmChartName    string
		HelmChartRepo    string
		HelmChartVersion string
		Namespace        string
	}{
		"ingress-nginx",
		"https://kubernetes.github.io/ingress-nginx",
		//https://github.com/kubernetes/ingress-nginx/blob/main/charts/ingress-nginx/Chart.yaml#L26C9-L26C14
		"4.11.1",
		"ingress-nginx",
	}
)
