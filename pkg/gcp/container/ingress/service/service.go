package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/controller"
	ingressnamespace "github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/addon/istio/ingress/namespace"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/ingress/gateway/kafka"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/ingress/gateway/postgres"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/container/ingress/gateway/redis"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-blueprint/pkg/gcp/network/ip"
	puluminameoutputcustom "github.com/plantoncloud-inc/pulumi-stack-runner-go-sdk/pkg/name/output/custom"
	wordpb "github.com/plantoncloud/planton-cloud-apis/zzgo/cloud/planton/apis/v1/commons/english/enums"
	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/compute"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	Namespace                       = ingressnamespace.Name
	ExternalLoadBalancerServiceName = "ingress-external"
	InternalLoadBalancerServiceName = "ingress-internal"
)

var (
	ExternalServiceAnnotations = map[string]string{
		"cloud.google.com/load-balancer-type": "external",
	}
	InternalServiceAnnotations = map[string]string{
		"cloud.google.com/load-balancer-type": "internal",
	}
)

type Input struct {
	AddedComputeIpAddress                  *ip.AddedIngressIpAddresses
	AddedIstioIngressControllerHelmRelease *helm.Release
}

// Resources adds two services of type LoadBalancer
// one load balancer should be internal and another one is external
func Resources(ctx *pulumi.Context, input *Input) error {
	if err := exportIpAddressOutputs(ctx, input.AddedComputeIpAddress); err != nil {
		return errors.Wrap(err, "failed to export ip address outputs")
	}
	if err := addService(ctx, input, InternalLoadBalancerServiceName, input.AddedComputeIpAddress.Internal, InternalServiceAnnotations); err != nil {
		return errors.Wrap(err, "failed to add internal service")
	}
	if err := addService(ctx, input, ExternalLoadBalancerServiceName, input.AddedComputeIpAddress.External, ExternalServiceAnnotations); err != nil {
		return errors.Wrap(err, "failed to add external service")
	}
	return nil
}

/*
apiVersion: v1
kind: Service
metadata:

	labels:
	  app: istio-ingress
	  istio: ingress
	name: ingress-internal
	namespace: istio-ingress

spec:

	externalTrafficPolicy: Cluster
	loadBalancerIP: 34.131.67.167
	ports:
	- name: status-port
	  port: 15021
	  protocol: TCP
	  targetPort: 15021
	- name: http2
	  port: 80
	  protocol: TCP
	  targetPort: 80
	- name: https
	  port: 443
	  protocol: TCP
	  targetPort: 443
	selector:
	  app: istio-ingress
	  istio: ingress
	type: LoadBalancer
*/
func addService(ctx *pulumi.Context, input *Input, serviceName string, addedIpAddress *compute.Address,
	annotations map[string]string) error {
	_, err := corev1.NewService(ctx, serviceName, &corev1.ServiceArgs{
		Metadata: metav1.ObjectMetaArgs{
			Name:        pulumi.String(serviceName),
			Namespace:   pulumi.String(Namespace),
			Annotations: pulumi.ToStringMap(annotations),
			Labels:      pulumi.ToStringMap(controller.SelectorLabels),
		},
		Spec: &corev1.ServiceSpecArgs{
			Type:           pulumi.String("LoadBalancer"),
			Selector:       pulumi.ToStringMap(controller.SelectorLabels),
			LoadBalancerIP: addedIpAddress.Address,
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Name:       pulumi.String("status-port"),
					Protocol:   pulumi.String("TCP"),
					Port:       pulumi.Int(controller.IstioStatusPort),
					TargetPort: pulumi.Int(controller.IstioStatusPort),
				},
				&corev1.ServicePortArgs{
					Name:       pulumi.String("http2"),
					Protocol:   pulumi.String("TCP"),
					Port:       pulumi.Int(controller.HttpPort),
					TargetPort: pulumi.Int(controller.HttpPort),
				},
				&corev1.ServicePortArgs{
					Name:       pulumi.String("https"),
					Protocol:   pulumi.String("TCP"),
					Port:       pulumi.Int(controller.HttpsPort),
					TargetPort: pulumi.Int(controller.HttpsPort),
				},
				&corev1.ServicePortArgs{
					Name:       pulumi.String("kafka-private"),
					Protocol:   pulumi.String("TCP"),
					Port:       pulumi.Int(kafka.ExternalPrivateListenerPortNumber),
					TargetPort: pulumi.Int(kafka.ExternalPrivateListenerPortNumber),
				},
				&corev1.ServicePortArgs{
					Name:       pulumi.String("kafka-public"),
					Protocol:   pulumi.String("TCP"),
					Port:       pulumi.Int(kafka.ExternalPublicListenerPortNumber),
					TargetPort: pulumi.Int(kafka.ExternalPublicListenerPortNumber),
				},
				&corev1.ServicePortArgs{
					Name:       pulumi.String("postgres"),
					Protocol:   pulumi.String("TCP"),
					Port:       pulumi.Int(postgres.ContainerPort),
					TargetPort: pulumi.Int(postgres.ContainerPort),
				},
				&corev1.ServicePortArgs{
					Name:       pulumi.String("debug"),
					Protocol:   pulumi.String("TCP"),
					Port:       pulumi.Int(controller.DebugPort),
					TargetPort: pulumi.Int(controller.DebugPort),
				},
				&corev1.ServicePortArgs{
					Name:       pulumi.String("redis"),
					Protocol:   pulumi.String("TCP"),
					Port:       pulumi.Int(redis.Port),
					TargetPort: pulumi.Int(redis.Port),
				},
			},
		},
	}, pulumi.Parent(input.AddedIstioIngressControllerHelmRelease),
		pulumi.DependsOn([]pulumi.Resource{input.AddedIstioIngressControllerHelmRelease}))
	if err != nil {
		return errors.Wrapf(err, "failed to add service")
	}
	return nil
}

func exportIpAddressOutputs(ctx *pulumi.Context, ingressIpAddress *ip.AddedIngressIpAddresses) error {
	ctx.Export(GetInternalIpOutputName(), ingressIpAddress.Internal.Address)
	ctx.Export(GetExternalIpOutputName(), ingressIpAddress.External.Address)
	return nil
}

func GetInternalIpOutputName() string {
	return puluminameoutputcustom.Name(fmt.Sprintf("%s-%s-%s", wordpb.Word_ingress.String(), wordpb.Word_internal.String(), wordpb.Word_ip.String()))
}

func GetExternalIpOutputName() string {
	return puluminameoutputcustom.Name(fmt.Sprintf("%s-%s-%s", wordpb.Word_ingress.String(), wordpb.Word_external.String(), wordpb.Word_ip.String()))
}
