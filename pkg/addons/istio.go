package addons

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud/gke-cluster-pulumi-module/pkg/localz"
	"github.com/plantoncloud/gke-cluster-pulumi-module/pkg/outputs"
	"github.com/plantoncloud/gke-cluster-pulumi-module/pkg/vars"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Istio installs the Istio service mesh in the Kubernetes cluster using Helm. It creates the necessary namespaces,
// installs the Helm charts for Istio base, Istiod, and gateway components, and sets up load balancers for ingress.
//
// Parameters:
// - ctx: The Pulumi context used for defining cloud resources.
// - locals: A struct containing local configuration and metadata.
// - createdCluster: The GKE cluster where Istio will be installed.
// - gcpProvider: The GCP provider for Pulumi.
// - kubernetesProvider: The Kubernetes provider for Pulumi.
//
// Returns:
// - error: An error object if there is any issue during the installation.
//
// The function performs the following steps:
// 1. Creates the `istio-system` namespace and labels it with metadata from locals.
// 2. Deploys the Istio base Helm chart into the `istio-system` namespace.
// 3. Deploys the Istiod Helm chart into the `istio-system` namespace with specific mesh configuration.
// 4. Creates the Istio gateway namespace and labels it with metadata from locals.
// 5. Deploys the Istio gateway Helm chart into the gateway namespace, configuring service ports for HTTP, HTTPS, and other protocols.
// 6. Creates a compute IP address for the internal load balancer and exports its address.
// 7. Creates a Kubernetes service for the internal load balancer using the created IP address and service port configurations.
// 8. Creates a compute IP address for the external load balancer and exports its address.
// 9. Creates a Kubernetes service for the external load balancer using the created IP address and service port configurations.
// 10. Handles errors and returns any errors encountered during the namespace creation, Helm release deployment, or service setup.
func Istio(ctx *pulumi.Context, locals *localz.Locals,
	createdCluster *container.Cluster, gcpProvider *gcp.Provider,
	kubernetesProvider *pulumikubernetes.Provider) error {
	//create istio-system namespace resource
	createdIstioSystemNamespace, err := corev1.NewNamespace(ctx,
		vars.Istio.SystemNamespace,
		&corev1.NamespaceArgs{
			Metadata: metav1.ObjectMetaPtrInput(
				&metav1.ObjectMetaArgs{
					Name:   pulumi.String(vars.Istio.SystemNamespace),
					Labels: pulumi.ToStringMap(locals.KubernetesLabels),
				}),
		},
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create istio-system namespace")
	}

	//create istio-base helm-release
	_, err = helm.NewRelease(ctx, "istio-base",
		&helm.ReleaseArgs{
			Name:            pulumi.String(vars.Istio.BaseHelmChartName),
			Namespace:       createdIstioSystemNamespace.Metadata.Name(),
			Chart:           pulumi.String(vars.Istio.BaseHelmChartName),
			Version:         pulumi.String(vars.Istio.HelmChartsVersion),
			CreateNamespace: pulumi.Bool(false),
			Atomic:          pulumi.Bool(false),
			CleanupOnFail:   pulumi.Bool(true),
			WaitForJobs:     pulumi.Bool(true),
			Timeout:         pulumi.Int(180),
			Values:          pulumi.Map{},
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String(vars.Istio.HelmChartsRepo),
			},
		}, pulumi.Parent(createdIstioSystemNamespace),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return errors.Wrap(err, "failed to create istio-base helm release")
	}

	//create istiod helm-release
	_, err = helm.NewRelease(ctx, "istiod",
		&helm.ReleaseArgs{
			Name:            pulumi.String(vars.Istio.IstiodHelmChartName),
			Namespace:       createdIstioSystemNamespace.Metadata.Name(),
			Chart:           pulumi.String(vars.Istio.IstiodHelmChartName),
			Version:         pulumi.String(vars.Istio.HelmChartsVersion),
			CreateNamespace: pulumi.Bool(false),
			Atomic:          pulumi.Bool(false),
			CleanupOnFail:   pulumi.Bool(true),
			WaitForJobs:     pulumi.Bool(true),
			Timeout:         pulumi.Int(180),
			Values: pulumi.Map{
				"meshConfig": pulumi.StringMap{
					"ingressClass":          pulumi.String("istio"),
					"ingressControllerMode": pulumi.String("STRICT"),
					"ingressService":        pulumi.String("ingress-external"),
					"ingressSelector":       pulumi.String("ingress"),
				},
			},
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String(vars.Istio.HelmChartsRepo),
			},
		}, pulumi.Parent(createdIstioSystemNamespace),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return errors.Wrap(err, "failed to create istiod helm release")
	}

	//create istio-gateway namespace resource
	createdIstioGatewayNamespace, err := corev1.NewNamespace(ctx,
		vars.Istio.GatewayNamespace,
		&corev1.NamespaceArgs{
			Metadata: metav1.ObjectMetaPtrInput(
				&metav1.ObjectMetaArgs{
					Name:   pulumi.String(vars.Istio.GatewayNamespace),
					Labels: pulumi.ToStringMap(locals.KubernetesLabels),
				}),
		},
		pulumi.Provider(kubernetesProvider))
	if err != nil {
		return errors.Wrapf(err, "failed to create istio-system namespace")
	}

	//create istio-gateway helm-release
	_, err = helm.NewRelease(ctx, "istio-gateway",
		&helm.ReleaseArgs{
			Name:            pulumi.String(vars.Istio.GatewayHelmChartName),
			Namespace:       createdIstioGatewayNamespace.Metadata.Name(),
			Chart:           pulumi.String(vars.Istio.GatewayHelmChartName),
			Version:         pulumi.String(vars.Istio.HelmChartsVersion),
			CreateNamespace: pulumi.Bool(false),
			Atomic:          pulumi.Bool(false),
			CleanupOnFail:   pulumi.Bool(true),
			WaitForJobs:     pulumi.Bool(true),
			Timeout:         pulumi.Int(180),
			Values: pulumi.Map{
				"service": pulumi.Map{
					"type":           pulumi.String("ClusterIP"),
					"loadBalancerIP": pulumi.String(""), // No LoadBalancer IP specified in the example
					"ports": pulumi.StringMapArray{
						pulumi.StringMap{
							"name":       pulumi.String("status-port"),
							"protocol":   pulumi.String("TCP"),
							"port":       pulumi.Sprintf("%d", 15021),
							"targetPort": pulumi.Sprintf("%d", 15021),
						},
						pulumi.StringMap{
							"name":       pulumi.String("http2"),
							"protocol":   pulumi.String("TCP"),
							"port":       pulumi.Sprintf("%d", 80),
							"targetPort": pulumi.Sprintf("%d", 80),
						},
						pulumi.StringMap{
							"name":       pulumi.String("https"),
							"protocol":   pulumi.String("TCP"),
							"port":       pulumi.Sprintf("%d", 443),
							"targetPort": pulumi.Sprintf("%d", 443),
						},
						pulumi.StringMap{
							"name":       pulumi.String("debug"),
							"protocol":   pulumi.String("TCP"),
							"port":       pulumi.Sprintf("%d", 5005),
							"targetPort": pulumi.Sprintf("%d", 5005),
						},
					},
				},
			},
			RepositoryOpts: helm.RepositoryOptsArgs{
				Repo: pulumi.String(vars.Istio.HelmChartsRepo),
			},
		}, pulumi.Parent(createdIstioGatewayNamespace),
		pulumi.IgnoreChanges([]string{"status", "description", "resourceNames"}))
	if err != nil {
		return errors.Wrap(err, "failed to create istio-gateway helm release")
	}

	//define array of ports to be configured for both internal and external ingress services
	loadBalancerServicePortArray := corev1.ServicePortArray{
		&corev1.ServicePortArgs{
			Name:       pulumi.String("status-port"),
			Protocol:   pulumi.String("TCP"),
			Port:       pulumi.Int(vars.Istio.IstiodStatusPort),
			TargetPort: pulumi.Int(vars.Istio.IstiodStatusPort),
		},
		&corev1.ServicePortArgs{
			Name:       pulumi.String("http2"),
			Protocol:   pulumi.String("TCP"),
			Port:       pulumi.Int(vars.Istio.HttpPort),
			TargetPort: pulumi.Int(vars.Istio.HttpPort),
		},
		&corev1.ServicePortArgs{
			Name:       pulumi.String("https"),
			Protocol:   pulumi.String("TCP"),
			Port:       pulumi.Int(vars.Istio.HttpsPort),
			TargetPort: pulumi.Int(vars.Istio.HttpsPort),
		},
		&corev1.ServicePortArgs{
			Name:       pulumi.String("kafka-private"),
			Protocol:   pulumi.String("TCP"),
			Port:       pulumi.Int(vars.Istio.KafkaConfig.ExternalPrivateListenerPortNumber),
			TargetPort: pulumi.Int(vars.Istio.KafkaConfig.ExternalPrivateListenerPortNumber),
		},
		&corev1.ServicePortArgs{
			Name:       pulumi.String("kafka-public"),
			Protocol:   pulumi.String("TCP"),
			Port:       pulumi.Int(vars.Istio.KafkaConfig.ExternalPublicListenerPortNumber),
			TargetPort: pulumi.Int(vars.Istio.KafkaConfig.ExternalPublicListenerPortNumber),
		},
		&corev1.ServicePortArgs{
			Name:       pulumi.String("postgres"),
			Protocol:   pulumi.String("TCP"),
			Port:       pulumi.Int(vars.Istio.PostgresPort),
			TargetPort: pulumi.Int(vars.Istio.PostgresPort),
		},
		&corev1.ServicePortArgs{
			Name:       pulumi.String("redis"),
			Protocol:   pulumi.String("TCP"),
			Port:       pulumi.Int(vars.Istio.RedisPort),
			TargetPort: pulumi.Int(vars.Istio.RedisPort),
		},
	}

	//create compute ip address for internal load-balancer
	createdIngressInternalLoadBalancerIp, err := compute.NewAddress(ctx,
		vars.Istio.IngressInternalLoadBalancerServiceName,
		&compute.AddressArgs{
			Name:        pulumi.Sprintf("%s-ingress-internal", locals.GkeCluster.Metadata.Id),
			Project:     createdCluster.Project,
			Region:      pulumi.String(locals.GkeCluster.Spec.Region),
			AddressType: pulumi.String("INTERNAL"),
			Labels:      pulumi.ToStringMap(locals.GcpLabels),
		}, pulumi.Parent(createdCluster))
	if err != nil {
		return errors.Wrap(err, "failed to create ip address for ingress-internal load-balancer")
	}

	//export ingress-internal ip
	ctx.Export(outputs.IngressInternalIp, createdIngressInternalLoadBalancerIp.Address)

	//create load-balancer service for internal load-balancer
	_, err = corev1.NewService(ctx,
		vars.Istio.IngressInternalLoadBalancerServiceName,
		&corev1.ServiceArgs{
			Metadata: metav1.ObjectMetaArgs{
				Name:        pulumi.String(vars.Istio.IngressInternalLoadBalancerServiceName),
				Namespace:   createdIstioGatewayNamespace.Metadata.Name(),
				Annotations: pulumi.ToStringMap(vars.Istio.IngressInternalServiceAnnotations),
				Labels:      pulumi.ToStringMap(vars.Istio.SelectorLabels),
			},
			Spec: &corev1.ServiceSpecArgs{
				Type:           pulumi.String("LoadBalancer"),
				Selector:       pulumi.ToStringMap(vars.Istio.SelectorLabels),
				LoadBalancerIP: createdIngressInternalLoadBalancerIp.Address,
				Ports:          loadBalancerServicePortArray,
			},
		}, pulumi.Parent(createdIstioGatewayNamespace))
	if err != nil {
		return errors.Wrapf(err, "failed to create ingress-external kubernetes service")
	}

	//create compute ip address for external load-balancer
	createdIngressExternalLoadBalancerIp, err := compute.NewAddress(ctx,
		vars.Istio.IngressExternalLoadBalancerServiceName,
		&compute.AddressArgs{
			Name:        pulumi.Sprintf("%s-ingress-external", locals.GkeCluster.Metadata.Id),
			Project:     createdCluster.Project,
			Region:      pulumi.String(locals.GkeCluster.Spec.Region),
			AddressType: pulumi.String("EXTERNAL"),
			Labels:      pulumi.ToStringMap(locals.GcpLabels),
		}, pulumi.Parent(createdCluster))
	if err != nil {
		return errors.Wrap(err, "failed to create ip address for ingress-internal load-balancer")
	}

	//export ingress-external ip
	ctx.Export(outputs.IngressExternalIp, createdIngressExternalLoadBalancerIp.Address)

	//create load-balancer service for external load-balancer
	_, err = corev1.NewService(ctx,
		vars.Istio.IngressExternalLoadBalancerServiceName,
		&corev1.ServiceArgs{
			Metadata: metav1.ObjectMetaArgs{
				Name:        pulumi.String(vars.Istio.IngressExternalLoadBalancerServiceName),
				Namespace:   createdIstioGatewayNamespace.Metadata.Name(),
				Annotations: pulumi.ToStringMap(vars.Istio.IngressExternalServiceAnnotations),
				Labels:      pulumi.ToStringMap(vars.Istio.SelectorLabels),
			},
			Spec: &corev1.ServiceSpecArgs{
				Type:           pulumi.String("LoadBalancer"),
				Selector:       pulumi.ToStringMap(vars.Istio.SelectorLabels),
				LoadBalancerIP: createdIngressExternalLoadBalancerIp.Address,
				Ports:          loadBalancerServicePortArray,
			},
		}, pulumi.Parent(createdIstioGatewayNamespace))
	if err != nil {
		return errors.Wrapf(err, "failed to create ingress-external kubernetes service")
	}

	return nil
}
