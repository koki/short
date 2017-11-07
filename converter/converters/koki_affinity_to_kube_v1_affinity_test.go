package converters

import (
	"testing"

	"github.com/koki/short/types"
)

var na0 = types.Affinity{
	NodeAffinity: "existentKey",
}

var na1 = types.Affinity{
	NodeAffinity: "!nonexistentKey",
}

var na2 = types.Affinity{
	NodeAffinity: "key1=value1,value2",
}

var na3 = types.Affinity{
	NodeAffinity: "key2!=value2&!key3&key4<10",
}

var na4 = types.Affinity{
	NodeAffinity: "key5=value5,value6:soft",
}

var na5 = types.Affinity{
	NodeAffinity: "value6weighted:soft:100",
}

var pa0 = types.Affinity{
	PodAffinity: "existentKeyPa",
}

var pa1 = types.Affinity{
	PodAffinity: "!nonexistentKeyPa",
	Topology:    "topoKey",
}

var pa2 = types.Affinity{
	PodAffinity: "key1pa=value1,value2",
	Namespaces:  []string{"ns1", "ns2", "ns3"},
}

var pa3 = types.Affinity{
	PodAffinity: "key1pa=value1&key2pa!=value2&!key3pa",
}

var pa4 = types.Affinity{
	PodAffinity: "key5pa=value5,value6:soft:50",
}

var paa0 = types.Affinity{
	PodAntiAffinity: "existentKeyPaa",
}

var paa1 = types.Affinity{
	PodAntiAffinity: "!nonexistentKeyPaaSoft:soft",
}

var paa2 = types.Affinity{
	PodAntiAffinity: "!nonexistentKeyPaaSoftWeighted:soft:10",
}

// TestRevert make sure the reversion doesn't throw errors.
func TestRevert(t *testing.T) {
	testRevertAffinities(t)
	testRevertAffinities(t, na0)
	testRevertAffinities(t, na0, na1, na2, na3, na4)
	testRevertAffinities(t, pa0)
	testRevertAffinities(t, pa0, pa1, pa2, pa3, pa4)
	testRevertAffinities(t, paa0)
	testRevertAffinities(t, paa0, paa1)
	testRevertAffinities(t, na0, pa0, paa0, na1, na4, pa1, paa1, paa2)
}

func testRevertAffinities(t *testing.T, as ...types.Affinity) {
	_, err := Convert_Koki_Affinity_to_Kube_v1_Affinity(as)
	if err != nil {
		t.Error(err)
	}
}
