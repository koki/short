package converters

import (
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_Ingress_to_Koki_Ingress(kubeIngress *v1beta1.Ingress) (*types.IngressWrapper, error) {
	var err error
	kokiWrapper := &types.IngressWrapper{}
	kokiIngress := &kokiWrapper.Ingress

	kokiIngress.Name = kubeIngress.Name
	kokiIngress.Namespace = kubeIngress.Namespace
	kokiIngress.Version = kubeIngress.APIVersion
	kokiIngress.Cluster = kubeIngress.ClusterName
	kokiIngress.Labels = kubeIngress.Labels
	kokiIngress.Annotations = kubeIngress.Annotations

	kubeSpec := kubeIngress.Spec

	kokiIngress.ServiceName, kokiIngress.ServicePort = convertIngressBackend(kubeSpec.Backend)
	kokiIngress.TLS = convertIngressTLS(kubeSpec.TLS)
	kokiIngress.Rules, err = convertIngressRules(kubeSpec.Rules)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "ingress rules")
	}

	kokiIngress.LoadBalancerIngress, err = convertLoadBalancerIngress(kubeIngress.Status.LoadBalancer.Ingress)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "ingress load balancer status")
	}

	return kokiWrapper, nil
}

func convertIngressBackend(kubeBackend *v1beta1.IngressBackend) (string, *intstr.IntOrString) {
	if kubeBackend == nil {
		return "", nil
	}

	return kubeBackend.ServiceName, &kubeBackend.ServicePort
}

func convertIngressTLS(kubeTLS []v1beta1.IngressTLS) []types.IngressTLS {
	if kubeTLS == nil {
		return nil
	}

	kokiTLS := make([]types.IngressTLS, len(kubeTLS))
	for i, kubeItem := range kubeTLS {
		kokiTLS[i] = types.IngressTLS{
			Hosts:      kubeItem.Hosts,
			SecretName: kubeItem.SecretName,
		}
	}

	return kokiTLS
}

func convertIngressRules(kubeRules []v1beta1.IngressRule) ([]types.IngressRule, error) {
	if kubeRules == nil {
		return nil, nil
	}

	kokiRules := make([]types.IngressRule, len(kubeRules))
	for i, kubeRule := range kubeRules {
		if kubeRule.HTTP == nil {
			return nil, serrors.InvalidInstanceErrorf(kubeRule, "HTTP is the only supported rule type, but this rule is missing its HTTP entry.")
		}
		kokiRules[i] = types.IngressRule{
			Host:  kubeRule.Host,
			Paths: convertHTTPIngressPaths(kubeRule.HTTP.Paths),
		}
	}

	return kokiRules, nil
}

func convertHTTPIngressPaths(kubePaths []v1beta1.HTTPIngressPath) []types.HTTPIngressPath {
	if kubePaths == nil {
		return nil
	}

	kokiPaths := make([]types.HTTPIngressPath, len(kubePaths))
	for i, kubePath := range kubePaths {
		kokiPaths[i] = types.HTTPIngressPath{
			Path:        kubePath.Path,
			ServiceName: kubePath.Backend.ServiceName,
			ServicePort: kubePath.Backend.ServicePort,
		}
	}

	return kokiPaths
}
