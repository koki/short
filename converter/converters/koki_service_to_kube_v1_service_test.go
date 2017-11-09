package converters

import (
	"net"
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/koki/short/types"
	"github.com/kr/pretty"
)

var s0 = &types.ServiceWrapper{
	Service: types.Service{
		Name:         "example",
		Version:      "v1",
		ExternalName: "externalDnsName",
	}}

var s1 = &types.ServiceWrapper{
	Service: types.Service{
		Name:    "example",
		Version: "v1",
		PodLabels: map[string]string{
			"labelKey": "labelValue",
		},
		ExternalIPs: []types.IPAddr{types.IPAddr("1.1.1.1")},
		Port: &types.ServicePort{
			Expose:  80,
			PodPort: intstr.FromInt(8080),
		},
		ClusterIP:        types.ClusterIPAddr(types.IPAddr("1.1.1.10")),
		ClientIPAffinity: types.ClientIPAffinityDefault(),
	}}

var s2 = &types.ServiceWrapper{
	Service: types.Service{
		Name:    "example",
		Version: "v1",
		PodLabels: map[string]string{
			"labelKey": "labelValue",
		},
		ExternalIPs: []types.IPAddr{types.IPAddr("1.1.1.1")},
		Ports: map[string]types.ServicePort{
			"http": types.ServicePort{
				Expose:   80,
				PodPort:  intstr.FromInt(8080),
				NodePort: 999,
			},
		},
		ClusterIP:        types.ClusterIPAddr(types.IPAddr("1.1.1.10")),
		ClientIPAffinity: types.ClientIPAffinityDefault(),
	}}

var s3 = &types.ServiceWrapper{
	Service: types.Service{
		Name:    "example",
		Version: "v1",
		PodLabels: map[string]string{
			"labelKey": "labelValue",
		},
		ExternalIPs: []types.IPAddr{types.IPAddr("1.1.1.1")},
		Ports: map[string]types.ServicePort{
			"http": types.ServicePort{
				Expose:   80,
				PodPort:  intstr.FromInt(8080),
				NodePort: 999,
				Protocol: types.ProtocolTCP,
			},
		},
		ClusterIP:             types.ClusterIPAddr(types.IPAddr("1.1.1.10")),
		ExternalTrafficPolicy: types.ExternalTrafficPolicyLocal,
		ClientIPAffinity:      types.ClientIPAffinitySeconds(30),
		LoadBalancer: &types.LoadBalancer{
			IP: types.IPAddr("100.1.1.1"),
			Allowed: []types.CIDR{
				"0.0.0.0/0",
			},
			Ingress: []types.Ingress{
				types.Ingress{Hostname: "ingressHostname"},
				types.Ingress{IP: net.ParseIP("1.2.3.4")},
			},
		},
	}}

func TestRevertService(t *testing.T) {
	var kubeService *v1.Service
	kubeService = tryService(s0, t)
	if kubeService.Spec.Type != v1.ServiceTypeExternalName {
		t.Errorf("unexpected type %s", kubeService.Spec.Type)
	}

	kubeService = tryService(s1, t)
	if kubeService.Spec.Type != v1.ServiceTypeClusterIP {
		t.Errorf("unexpected type %s", kubeService.Spec.Type)
	}

	kubeService = tryService(s2, t)
	if kubeService.Spec.Type != v1.ServiceTypeNodePort {
		t.Errorf("unexpected type %s", kubeService.Spec.Type)
	}

	kubeService = tryService(s3, t)
	if kubeService.Spec.Type != v1.ServiceTypeLoadBalancer {
		t.Errorf("unexpected type %s", kubeService.Spec.Type)
	}
}

func tryService(kokiService *types.ServiceWrapper, t *testing.T) *v1.Service {
	kubeService, err := Convert_Koki_Service_To_Kube_v1_Service(kokiService)
	if err != nil {
		t.Error(pretty.Sprintf("failed converting (%# v) with error (%s)", kokiService, err.Error()))
	}

	roundTripped, err := Convert_Kube_v1_Service_to_Koki_Service(kubeService)
	if err != nil {
		t.Error(pretty.Sprintf("failed reverting (%# v) with error (%s)", kubeService, err.Error()))
	}

	if !reflect.DeepEqual(kokiService, roundTripped) {
		t.Error(pretty.Sprintf(
			"failed round-trip:\n(%# v)\n(%# v)\n(%# v)",
			kokiService,
			roundTripped,
			kubeService,
		))
	}

	return kubeService
}
