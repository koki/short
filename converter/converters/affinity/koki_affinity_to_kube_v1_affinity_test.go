package affinity

import (
	"testing"

	"github.com/koki/short/types"
)

var nodeAffinity0 = types.Affinity{
	NodeAffinity: "existentKey",
}

var nodeAffinity1 = types.Affinity{
	NodeAffinity: "!nonexistentKey",
}

var nodeAffinity2 = types.Affinity{
	NodeAffinity: "key1=value1,value2",
}

var nodeAffinity3 = types.Affinity{
	NodeAffinity: "key2!=value2&!key3&key4<10",
}

var nodeAffinity4 = types.Affinity{
	NodeAffinity: "key5=value5,value6:soft",
}

var nodeAffinity5 = types.Affinity{
	NodeAffinity: "value6weighted:soft:100",
}

var podAffinity0 = types.Affinity{
	PodAffinity: "existentKeyPa",
}

var podAffinity1 = types.Affinity{
	PodAffinity: "!nonexistentKeyPa",
	Topology:    "topoKey",
}

var podAffinity2 = types.Affinity{
	PodAffinity: "key1pa=value1,value2",
	Namespaces:  []string{"ns1", "ns2", "ns3"},
}

var podAffinity3 = types.Affinity{
	PodAffinity: "key1pa=value1&key2pa!=value2&!key3pa",
}

var podAffinity4 = types.Affinity{
	PodAffinity: "key5pa=value5,value6:soft:50",
}

var podAntiAffinity0 = types.Affinity{
	PodAntiAffinity: "existentKeyPaa",
}

var podAntiAffinity1 = types.Affinity{
	PodAntiAffinity: "!nonexistentKeyPaaSoft:soft",
}

var podAntiAffinity2 = types.Affinity{
	PodAntiAffinity: "!nonexistentKeyPaaSoftWeighted:soft:10",
}

// TestRevert make sure the reversion doesn't throw errors.
func TestRevert(t *testing.T) {
	testRevertAffinities(t)
	testRevertAffinities(t, nodeAffinity0)
	testRevertAffinities(t, nodeAffinity0, nodeAffinity1, nodeAffinity2, nodeAffinity3, nodeAffinity4, nodeAffinity5)
	testRevertAffinities(t, podAffinity0)
	testRevertAffinities(t, podAffinity0, podAffinity1, podAffinity2, podAffinity3, podAffinity4)
	testRevertAffinities(t, podAntiAffinity0)
	testRevertAffinities(t, podAntiAffinity0, podAntiAffinity1)
	testRevertAffinities(t, nodeAffinity0, podAffinity0, podAntiAffinity0, nodeAffinity1, nodeAffinity4, podAffinity1, podAntiAffinity1, podAntiAffinity2)
}

func testRevertAffinities(t *testing.T, as ...types.Affinity) {
	_, err := Convert_Koki_Affinity_to_Kube_v1_Affinity(as)
	if err != nil {
		t.Error(err)
	}
}
