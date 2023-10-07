package gateway

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/kubernetes/manifest"
	"github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/addon/istio/ingress/controller"
	ingressnamespace "github.com/plantoncloud-inc/kube-cluster-pulumi-stack/pkg/gcp/container/addon/istio/ingress/namespace"
	"github.com/plantoncloud-inc/stack-runner-service/internal/domain/code2cloud/deploy/kafka/kubernetes/listener"
	productstoragedatabasepostgresclusterstackimplcluster "github.com/plantoncloud-inc/stack-runner-service/internal/domain/code2cloud/deploy/postgres/kubernetes/cluster"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	networkingv1beta1 "istio.io/api/networking/v1beta1"
	"istio.io/client-go/pkg/apis/networking/v1beta1"
	k8smetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path/filepath"
)

const (
	Namespace           = ingressnamespace.Name
	KafkaGatewayName    = "kafka"
	PostgresGatewayName = "postgres"
)

type Input struct {
	Workspace                              string
	AddedIstioIngressControllerHelmRelease *helm.Release
}

// Resources adds a kafka gateway
func Resources(ctx *pulumi.Context, input *Input) error {
	if err := addKafkaGateway(ctx, input); err != nil {
		return errors.Wrap(err, "failed to add kafka  gateway")
	}
	if err := addPostgresGateway(ctx, input); err != nil {
		return errors.Wrap(err, "failed to add kafka  gateway")
	}
	return nil
}

func addKafkaGateway(ctx *pulumi.Context, input *Input) error {
	gatewayObject := buildTlsPassThroughGatewayObject(
		KafkaGatewayName,
		Namespace,
		listener.ExternalPublicListenerPortNumber,
	)
	resourceName := fmt.Sprintf("gateway-%s", KafkaGatewayName)
	manifestPath := filepath.Join(input.Workspace, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, gatewayObject); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}
	_, err := pulumik8syaml.NewConfigFile(ctx, resourceName, &pulumik8syaml.ConfigFileArgs{File: manifestPath},
		pulumi.Parent(input.AddedIstioIngressControllerHelmRelease),
		pulumi.DependsOn([]pulumi.Resource{input.AddedIstioIngressControllerHelmRelease}))
	if err != nil {
		return errors.Wrap(err, "failed to add ingress-gateway manifest")
	}
	return nil
}

func addPostgresGateway(ctx *pulumi.Context, input *Input) error {
	gatewayObject := buildTlsPassThroughGatewayObject(
		PostgresGatewayName,
		Namespace,
		productstoragedatabasepostgresclusterstackimplcluster.PostgresContainerPort,
	)
	resourceName := fmt.Sprintf("gateway-%s", PostgresGatewayName)
	manifestPath := filepath.Join(input.Workspace, fmt.Sprintf("%s.yaml", resourceName))
	if err := manifest.Create(manifestPath, gatewayObject); err != nil {
		return errors.Wrapf(err, "failed to create %s manifest file", manifestPath)
	}
	_, err := pulumik8syaml.NewConfigFile(ctx, resourceName, &pulumik8syaml.ConfigFileArgs{File: manifestPath},
		pulumi.Parent(input.AddedIstioIngressControllerHelmRelease),
		pulumi.DependsOn([]pulumi.Resource{input.AddedIstioIngressControllerHelmRelease}))
	if err != nil {
		return errors.Wrap(err, "failed to add ingress-gateway manifest")
	}
	return nil
}

/*
apiVersion: networking.istio.io/v1beta1
kind: Gateway
metadata:

	name: kafka
	namespace: istio-ingress

spec:

	selector:
	  app: istio-ingress
	  istio: ingress
	servers:
	  - port:
	      number: 9092 or 5432
	      name: tcp-kafka or tcp-postgres
	      protocol: TLS
	    hosts:
	      - "*"
	    tls:
	      mode: PASSTHROUGH
*/
func buildTlsPassThroughGatewayObject(name, namespace string, portNumber uint32) *v1beta1.Gateway {
	return &v1beta1.Gateway{
		TypeMeta: k8smetav1.TypeMeta{
			APIVersion: "networking.istio.io/v1beta1",
			Kind:       "Gateway",
		},
		ObjectMeta: k8smetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: networkingv1beta1.Gateway{
			Selector: controller.SelectorLabels,
			Servers: []*networkingv1beta1.Server{
				{
					Port: &networkingv1beta1.Port{
						Number:   portNumber,
						Protocol: "TLS",
						Name:     fmt.Sprintf("tcp-%s", name),
					},
					Hosts: []string{"*"},
					Tls: &networkingv1beta1.ServerTLSSettings{
						Mode: networkingv1beta1.ServerTLSSettings_PASSTHROUGH,
					},
				},
			},
		},
	}
}
