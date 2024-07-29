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
		Namespace            string
		HelmChartName        string
		HelmChartRepo        string
		HelmChartVersion     string
		KsaName              string
		SelfSignedIssuerName string
	}{
		Namespace:            "cert-manager",
		HelmChartName:        "cert-manager",
		HelmChartRepo:        "https://charts.jetstack.io",
		HelmChartVersion:     "v1.15.1",
		KsaName:              "cert-manager",
		SelfSignedIssuerName: "self-signed",
	}

	ExternalSecrets = struct {
		Namespace                               string
		HelmChartName                           string
		HelmChartRepo                           string
		HelmChartVersion                        string
		KsaName                                 string
		SecretsPollingIntervalSeconds           int
		GcpSecretsManagerClusterSecretStoreName string
	}{
		Namespace:        "external-secrets",
		HelmChartName:    "external-secrets",
		HelmChartRepo:    "https://charts.external-secrets.io",
		HelmChartVersion: "v0.9.20",
		KsaName:          "external-secrets",
		//caution: polling interval frequency may have effect on provider costs on some platforms
		SecretsPollingIntervalSeconds:           10,
		GcpSecretsManagerClusterSecretStoreName: "gcp-secrets-manager",
	}

	IngressNginx = struct {
		Namespace        string
		HelmChartName    string
		HelmChartRepo    string
		HelmChartVersion string
	}{
		Namespace:     "ingress-nginx",
		HelmChartName: "ingress-nginx",
		HelmChartRepo: "https://kubernetes.github.io/ingress-nginx",
		//https://github.com/kubernetes/ingress-nginx/blob/main/charts/ingress-nginx/Chart.yaml#L26C9-L26C14
		HelmChartVersion: "4.11.1",
	}

	ZalandoPostgresOperator = struct {
		Namespace        string
		HelmChartName    string
		HelmChartRepo    string
		HelmChartVersion string
	}{
		Namespace:     "postgres-operator",
		HelmChartName: "postgres-operator",
		HelmChartRepo: "https://opensource.zalando.com/postgres-operator/charts/postgres-operator",
		//https://github.com/zalando/postgres-operator/blob/v1.12.2/charts/postgres-operator/Chart.yaml#L3
		HelmChartVersion: "1.12.2",
	}

	SolrOperator = struct {
		Namespace              string
		CrdManifestDownloadUrl string
		HelmChartName          string
		HelmChartRepo          string
		HelmChartVersion       string
	}{
		Namespace: "solr-operator",
		//version in the url should match the helm-chart version and should be prefixed with 'v'
		CrdManifestDownloadUrl: "https://solr.apache.org/operator/downloads/crds/v0.8.1/all-with-dependencies.yaml",
		HelmChartName:          "solr-operator",
		HelmChartRepo:          "https://solr.apache.org/charts",
		//https://github.com/apache/solr-operator/blob/v0.8.1/helm/solr-operator/Chart.yaml#L18
		HelmChartVersion: "0.8.1",
	}

	StrimziKafkaOperator = struct {
		Namespace        string
		HelmChartName    string
		HelmChartRepo    string
		HelmChartVersion string
	}{
		Namespace:     "strimzi-operator",
		HelmChartName: "strimzi-kafka-operator",
		//https://artifacthub.io/packages/helm/strimzi/strimzi-kafka-operator
		HelmChartRepo: "https://strimzi.io/charts/",
		//check artifact-hub for the latest version
		HelmChartVersion: "0.42.0",
	}

	Istio = struct {
		SystemNamespace                        string
		GatewayNamespace                       string
		HelmChartsRepo                         string
		HelmChartsVersion                      string
		BaseHelmChartName                      string
		IstiodHelmChartName                    string
		GatewayHelmChartName                   string
		IngressInternalLoadBalancerServiceName string
		IngressExternalLoadBalancerServiceName string
		IngressExternalServiceAnnotations      map[string]string
		IngressInternalServiceAnnotations      map[string]string
		KafkaConfig                            struct {
			GatewayName                       string
			ExternalPublicListenerName        string
			ExternalPublicListenerPortNumber  int //this port is intended to be used by clients output the private network and outside the container cluster
			ExternalPrivateListenerName       string
			ExternalPrivateListenerPortNumber int //this port is intended to be used by clients inside the private network but outside the container cluster
			InternalListenerPortNumber        int //this port is intended to be used by clients inside the container cluster
		}
		SelectorLabels   map[string]string
		HttpPort         int
		HttpsPort        int
		IstiodStatusPort int
		PostgresPort     int
		RedisPort        int
	}{
		SystemNamespace:  "istio-system",
		GatewayNamespace: "istio-ingress",
		HelmChartsRepo:   "https://istio-release.storage.googleapis.com/charts",
		//all three charts are versioned separately but consistently. so we use the same version for all three charts.
		HelmChartsVersion: "1.22.3",
		//https://artifacthub.io/packages/helm/istio-official/base
		BaseHelmChartName: "base",
		//https://artifacthub.io/packages/helm/istio-official/istiod
		IstiodHelmChartName: "istiod",
		//https://artifacthub.io/packages/helm/istio-official/gateway
		GatewayHelmChartName:                   "gateway",
		IngressExternalLoadBalancerServiceName: "ingress-external",
		IngressInternalLoadBalancerServiceName: "ingress-internal",
		IngressExternalServiceAnnotations: map[string]string{
			"cloud.google.com/load-balancer-type": "external",
		},
		IngressInternalServiceAnnotations: map[string]string{
			"cloud.google.com/load-balancer-type": "internal",
		},
		KafkaConfig: struct {
			GatewayName                       string
			ExternalPublicListenerName        string
			ExternalPublicListenerPortNumber  int
			ExternalPrivateListenerName       string
			ExternalPrivateListenerPortNumber int
			InternalListenerPortNumber        int
		}{
			GatewayName:                       "kafka",
			ExternalPublicListenerName:        "extpub",
			ExternalPublicListenerPortNumber:  9092,
			ExternalPrivateListenerName:       "extpvt",
			ExternalPrivateListenerPortNumber: 9093,
			InternalListenerPortNumber:        9094,
		},
		SelectorLabels: map[string]string{
			"app":   "istio-ingress",
			"istio": "ingress",
		},
		HttpPort:         80,
		HttpsPort:        443,
		IstiodStatusPort: 15021,
		PostgresPort:     5432,
		RedisPort:        6789,
	}
)
