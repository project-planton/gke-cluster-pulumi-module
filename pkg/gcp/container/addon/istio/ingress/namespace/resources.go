package namespace

import (
	"github.com/pkg/errors"
	pulumikubernetes "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	pk8smv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const (
	Name = "istio-ingress"
)

var Labels = map[string]string{"istio-injection": "enabled"}

type Input struct {
	KubernetesProvider *pulumikubernetes.Provider
}

func Resources(ctx *pulumi.Context, input *Input) (*v1.Namespace, error) {
	ns, err := v1.NewNamespace(ctx, Name, &v1.NamespaceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("AddedNamespace"),
		Metadata: pk8smv1.ObjectMetaArgs{
			Name:   pulumi.String(Name),
			Labels: pulumi.ToStringMap(Labels),
		},
	}, pulumi.Provider(input.KubernetesProvider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to add namespace")
	}
	return ns, nil
}
