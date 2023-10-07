package envoyfilter

import (
	"github.com/pkg/errors"
	"github.com/plantoncloud-inc/go-commons/util/file"
	v1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/helm/v3"
	pulumik8syaml "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/yaml"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"os"
	"path/filepath"
)

const ManifestYaml = `
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: gateway-grpc-web-filter
  namespace: istio-ingress
spec:
  workloadSelector:
    labels:
      app: istio-ingress
      istio: ingress
  configPatches:
    - applyTo: HTTP_FILTER
      match:
        # importantly we're patching to the GATEWAY envoy, not sidecar
        context: GATEWAY
        listener:
          filterChain:
            filter:
              name: "envoy.filters.network.http_connection_manager"
              subFilter:
                # apply the patch before the cors filter, just like the one in
                # grpc-web example
                name: "envoy.filters.http.cors"
      patch:
        operation: INSERT_BEFORE
        value:
          name: envoy.filters.http.grpc_web
`

type Input struct {
	WorkspaceDir                      string
	AddedIstioIngressNamespace        *v1.Namespace
	AddedIngressControllerHelmRelease *helm.Release
}

// Resources installs envoy-filter required to support for grpc-web ingress requests
func Resources(ctx *pulumi.Context, input *Input) error {
	manifestPath := filepath.Join(input.WorkspaceDir, "envoy-filter.yaml")
	if !file.IsDirExists(filepath.Dir(manifestPath)) {
		if err := os.MkdirAll(filepath.Dir(manifestPath), os.ModePerm); err != nil {
			return errors.Wrapf(err, "failed to ensure %s dir", filepath.Dir(manifestPath))
		}
	}
	if err := os.WriteFile(manifestPath, []byte(ManifestYaml), os.ModePerm); err != nil {
		return errors.Wrapf(err, "failed to write %s file", manifestPath)
	}
	_, err := pulumik8syaml.NewConfigFile(ctx, "envoy-filter",
		&pulumik8syaml.ConfigFileArgs{File: manifestPath},
		pulumi.Parent(input.AddedIstioIngressNamespace),
		pulumi.DependsOn([]pulumi.Resource{input.AddedIngressControllerHelmRelease}))
	if err != nil {
		return errors.Wrap(err, "failed to add ingress-gateway manifest")
	}
	return nil
}
