package converters

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kr/pretty"

	"github.com/koki/short/converter/converters/affinity"
)

var nodeAffinity0 = &v1.NodeSelectorTerm{
	MatchExpressions: []v1.NodeSelectorRequirement{
		v1.NodeSelectorRequirement{
			Key:      "existentKey",
			Operator: "Exists",
		},
	},
}

var nodeAffinity1 = &v1.NodeSelectorTerm{
	MatchExpressions: []v1.NodeSelectorRequirement{
		v1.NodeSelectorRequirement{
			Key:      "nonexistentKey",
			Operator: "DoesNotExist",
		},
	},
}

var nodeAffinity2 = &v1.NodeSelectorTerm{
	MatchExpressions: []v1.NodeSelectorRequirement{
		v1.NodeSelectorRequirement{
			Key:      "key1",
			Operator: "In",
			Values:   []string{"value1", "value2"},
		},
	},
}

var nodeAffinity3 = &v1.NodeSelectorTerm{
	MatchExpressions: []v1.NodeSelectorRequirement{
		v1.NodeSelectorRequirement{
			Key:      "key2",
			Operator: "NotIn",
			Values:   []string{"value2"},
		},
		v1.NodeSelectorRequirement{
			Key:      "key3",
			Operator: "DoesNotExist",
		},
		v1.NodeSelectorRequirement{
			Key:      "key4",
			Operator: "Lt",
			Values:   []string{"10"},
		},
	},
}

var nodeAffinity4 = &v1.PreferredSchedulingTerm{
	Preference: v1.NodeSelectorTerm{
		MatchExpressions: []v1.NodeSelectorRequirement{
			v1.NodeSelectorRequirement{
				Key:      "key5",
				Operator: "In",
				Values:   []string{"value5", "value6"},
			},
		},
	},
}

var nodeAffinity5 = &v1.PreferredSchedulingTerm{
	Weight: 100,
	Preference: v1.NodeSelectorTerm{
		MatchExpressions: []v1.NodeSelectorRequirement{
			v1.NodeSelectorRequirement{
				Key:      "value6weighted",
				Operator: "Exists",
			},
		},
	},
}

var podAffinity0 = &v1.PodAffinityTerm{
	LabelSelector: &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			metav1.LabelSelectorRequirement{
				Key:      "existentKeyPa",
				Operator: "Exists",
			},
		},
	},
}

var podAffinity1 = &v1.PodAffinityTerm{
	LabelSelector: &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			metav1.LabelSelectorRequirement{
				Key:      "!nonexistentKeyPa",
				Operator: "DoesNotExist",
			},
		},
	},
	TopologyKey: "topoKey",
}

var podAffinity2 = &v1.PodAffinityTerm{
	LabelSelector: &metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			metav1.LabelSelectorRequirement{
				Key:      "key1pa",
				Operator: "In",
				Values:   []string{"value1", "value2"},
			},
		},
	},
	Namespaces: []string{"ns1", "ns2", "ns3"},
}

var podAffinity3 = &v1.PodAffinityTerm{
	LabelSelector: &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"key1pa": "value1",
		},
		MatchExpressions: []metav1.LabelSelectorRequirement{
			metav1.LabelSelectorRequirement{
				Key:      "key2pa",
				Operator: "NotIn",
				Values:   []string{"value2"},
			},
			metav1.LabelSelectorRequirement{
				Key:      "key3pa",
				Operator: "DoesNotExist",
			},
		},
	},
}

var podAffinity4 = &v1.WeightedPodAffinityTerm{
	Weight: 50,
	PodAffinityTerm: v1.PodAffinityTerm{
		LabelSelector: &metav1.LabelSelector{
			MatchExpressions: []metav1.LabelSelectorRequirement{
				metav1.LabelSelectorRequirement{
					Key:      "key5pa",
					Operator: "In",
					Values:   []string{"value5", "value6"},
				},
			},
		},
	},
}

func TestAffinityConversion(t *testing.T) {
	doTestHardNodeAffinity(nodeAffinity0, t)
	doTestHardNodeAffinity(nodeAffinity1, t)
	doTestHardNodeAffinity(nodeAffinity2, t)
	doTestHardNodeAffinity(nodeAffinity3, t)
	doTestSoftNodeAffinity(nodeAffinity4, t)
	doTestSoftNodeAffinity(nodeAffinity5, t)

	doTestHardPodAffinity(podAffinity0, t)
	doTestHardPodAffinity(podAffinity1, t)
	doTestHardPodAffinity(podAffinity2, t)
	doTestHardPodAffinity(podAffinity3, t)
	doTestSoftPodAffinity(podAffinity4, t)

	doTestHardPodAntiAffinity(podAffinity0, t)
	doTestHardPodAntiAffinity(podAffinity1, t)
	doTestHardPodAntiAffinity(podAffinity2, t)
	doTestHardPodAntiAffinity(podAffinity3, t)
	doTestSoftPodAntiAffinity(podAffinity4, t)
}

func doTestHardNodeAffinity(hardNodeAffinity *v1.NodeSelectorTerm, t *testing.T) {
	doTestNodeAffinity(&v1.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: &v1.NodeSelector{
			NodeSelectorTerms: []v1.NodeSelectorTerm{*hardNodeAffinity},
		},
	}, t)
}

func doTestSoftNodeAffinity(softNodeAffinity *v1.PreferredSchedulingTerm, t *testing.T) {
	doTestNodeAffinity(&v1.NodeAffinity{
		PreferredDuringSchedulingIgnoredDuringExecution: []v1.PreferredSchedulingTerm{
			*softNodeAffinity,
		},
	}, t)
}

func doTestNodeAffinity(nodeAffinity *v1.NodeAffinity, t *testing.T) {
	kubeAffinity := &v1.Affinity{
		NodeAffinity: nodeAffinity,
	}
	kokiAffinities, err := convertNodeAffinity(kubeAffinity.NodeAffinity)
	if err != nil {
		t.Error(err)
	}
	result, err := affinity.Convert_Koki_Affinity_to_Kube_v1_Affinity(kokiAffinities)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(kubeAffinity, result) {
		t.Error(pretty.Sprintf("failed round-trip:\n(%# v)\n(%# v)\n(%# v)", kubeAffinity, result, kokiAffinities))
	}
}

func doTestHardPodAffinity(podAffinity *v1.PodAffinityTerm, t *testing.T) {
	doTestPodAffinity(&v1.PodAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: []v1.PodAffinityTerm{
			*podAffinity,
		},
	}, t)
}

func doTestSoftPodAffinity(podAffinity *v1.WeightedPodAffinityTerm, t *testing.T) {
	doTestPodAffinity(&v1.PodAffinity{
		PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
			*podAffinity,
		},
	}, t)
}

func doTestPodAffinity(podAffinity *v1.PodAffinity, t *testing.T) {
	kubeAffinity := &v1.Affinity{
		PodAffinity: podAffinity,
	}
	kokiAffinities, err := convertPodAffinity(kubeAffinity.PodAffinity)
	if err != nil {
		t.Error(err)
	}
	result, err := affinity.Convert_Koki_Affinity_to_Kube_v1_Affinity(kokiAffinities)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(kubeAffinity, result) {
		t.Error(pretty.Sprintf("failed round-trip:\n(%# v)\n(%# v)\n(%# v)", kubeAffinity, result, kokiAffinities))
	}
}

func doTestHardPodAntiAffinity(podAffinity *v1.PodAffinityTerm, t *testing.T) {
	doTestPodAntiAffinity(&v1.PodAntiAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution: []v1.PodAffinityTerm{
			*podAffinity,
		},
	}, t)
}

func doTestSoftPodAntiAffinity(podAffinity *v1.WeightedPodAffinityTerm, t *testing.T) {
	doTestPodAntiAffinity(&v1.PodAntiAffinity{
		PreferredDuringSchedulingIgnoredDuringExecution: []v1.WeightedPodAffinityTerm{
			*podAffinity,
		},
	}, t)
}

func doTestPodAntiAffinity(podAntiAffinity *v1.PodAntiAffinity, t *testing.T) {
	kubeAffinity := &v1.Affinity{
		PodAntiAffinity: podAntiAffinity,
	}
	kokiAffinities, err := convertPodAntiAffinity(kubeAffinity.PodAntiAffinity)
	if err != nil {
		t.Error(err)
	}
	result, err := affinity.Convert_Koki_Affinity_to_Kube_v1_Affinity(kokiAffinities)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(kubeAffinity, result) {
		t.Error(pretty.Sprintf("failed round-trip:\n(%# v)\n(%# v)\n(%# v)", kubeAffinity, result, kokiAffinities))
	}
}
