package converters

import (
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_Ingress_to_Kube_Ingress(ingress *types.IngressWrapper) (*v1beta1.Ingress, error) {
	var err error
	kubeIngress := &v1beta1.Ingress{}
	kokiIngress := &ingress.Ingress

	kubeIngress.Name = kokiIngress.Name
	kubeIngress.Namespace = kokiIngress.Namespace
	if len(kokiIngress.Version) == 0 {
		kubeIngress.APIVersion = "extensions/v1beta1"
	} else {
		kubeIngress.APIVersion = kokiIngress.Version
	}
	kubeIngress.Kind = "Ingress"
	kubeIngress.ClusterName = kokiIngress.Cluster
	kubeIngress.Labels = kokiIngress.Labels
	kubeIngress.Annotations = kokiIngress.Annotations

	kubeIngress.Spec.Backend, err = revertIngressBackend(kokiIngress.ServiceName, kokiIngress.ServicePort)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "ingress backend/backend_port")
	}
	kubeIngress.Spec.TLS = revertIngressTLS(kokiIngress.TLS)
	kubeIngress.Spec.Rules = revertIngressRules(kokiIngress.Rules)
	kubeIngress.Status.LoadBalancer.Ingress = revertLoadBalancerIngress(kokiIngress.LoadBalancerIngress)

	return kubeIngress, nil
}

func revertIngressBackend(serviceName string, servicePort *intstr.IntOrString) (*v1beta1.IngressBackend, error) {
	if len(serviceName) == 0 && servicePort == nil {
		return nil, nil
	}

	if servicePort == nil {
		return nil, serrors.InvalidValueErrorf(servicePort, "if service name is specified, service port is also required")
	}
	if len(serviceName) == 0 {
		return nil, serrors.InvalidValueErrorf(serviceName, "if service port is specified, service name is also required")
	}

	return &v1beta1.IngressBackend{
		ServiceName: serviceName,
		ServicePort: *servicePort,
	}, nil
}

func revertIngressTLS(kokiTLS []types.IngressTLS) []v1beta1.IngressTLS {
	if kokiTLS == nil {
		return nil
	}

	kubeTLS := make([]v1beta1.IngressTLS, len(kokiTLS))
	for i, kokiItem := range kokiTLS {
		kubeTLS[i] = v1beta1.IngressTLS{
			Hosts:      kokiItem.Hosts,
			SecretName: kokiItem.SecretName,
		}
	}

	return kubeTLS
}

func revertIngressRules(kokiRules []types.IngressRule) []v1beta1.IngressRule {
	if kokiRules == nil {
		return nil
	}

	kubeRules := make([]v1beta1.IngressRule, len(kokiRules))
	for i, kokiRule := range kokiRules {
		kubeRules[i] = v1beta1.IngressRule{
			Host: kokiRule.Host,
			IngressRuleValue: v1beta1.IngressRuleValue{
				HTTP: &v1beta1.HTTPIngressRuleValue{
					Paths: revertHTTPIngressPaths(kokiRule.Paths),
				},
			},
		}
	}

	return kubeRules
}

func revertHTTPIngressPaths(kokiPaths []types.HTTPIngressPath) []v1beta1.HTTPIngressPath {
	if kokiPaths == nil {
		return nil
	}

	kubePaths := make([]v1beta1.HTTPIngressPath, len(kokiPaths))
	for i, kokiPath := range kokiPaths {
		kubePaths[i] = v1beta1.HTTPIngressPath{
			Path: kokiPath.Path,
			Backend: v1beta1.IngressBackend{
				ServiceName: kokiPath.ServiceName,
				ServicePort: kokiPath.ServicePort,
			},
		}
	}

	return kubePaths
}
