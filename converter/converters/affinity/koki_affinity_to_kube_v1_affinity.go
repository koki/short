package affinity

import (
	"strconv"
	"strings"

	"k8s.io/api/core/v1"

	"github.com/golang/glog"
	"github.com/koki/short/parser/expressions"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Koki_Affinity_to_Kube_v1_Affinity(kokiAffinities []types.Affinity) (*v1.Affinity, error) {
	node, pod, antiPod, err := splitAffinities(kokiAffinities)
	if err != nil {
		return nil, err
	}

	kubeNode, err := revertNodeAffinity(node)
	if err != nil {
		return nil, err
	}

	kubePod, err := revertPodAffinity(pod)
	if err != nil {
		return nil, err
	}

	kubeAntiPod, err := revertPodAntiAffinity(antiPod)
	if err != nil {
		return nil, err
	}

	if kubeNode == nil && kubePod == nil && kubeAntiPod == nil {
		return nil, nil
	}

	return &v1.Affinity{
		NodeAffinity:    kubeNode,
		PodAffinity:     kubePod,
		PodAntiAffinity: kubeAntiPod,
	}, nil
}

// PodAffinity is the subset of koki Affinity fields for pod affinity.
type PodAffinity struct {
	Affinity   string
	Topology   string
	Namespaces []string
}

// Separate a generic list of Affinities into a list for each type of Affinity.
func splitAffinities(affinities []types.Affinity) (node []string, pod, antiPod []PodAffinity, err error) {
	node = []string{}
	pod = []PodAffinity{}
	antiPod = []PodAffinity{}

	for _, affinity := range affinities {
		switch {
		case len(affinity.NodeAffinity) > 0:
			node = append(node, affinity.NodeAffinity)
		case len(affinity.PodAffinity) > 0:
			pod = append(pod, PodAffinity{
				Affinity:   affinity.PodAffinity,
				Topology:   affinity.Topology,
				Namespaces: affinity.Namespaces,
			})
		case len(affinity.PodAntiAffinity) > 0:
			antiPod = append(antiPod, PodAffinity{
				Affinity:   affinity.PodAntiAffinity,
				Topology:   affinity.Topology,
				Namespaces: affinity.Namespaces,
			})
		default:
			err = util.InvalidInstanceError(affinity)
		}
	}

	return
}

func revertPodAntiAffinity(affinities []PodAffinity) (*v1.PodAntiAffinity, error) {
	if len(affinities) == 0 {
		return nil, nil
	}

	hard, soft, err := splitAndRevertPodAffinity(affinities)
	if err != nil {
		return nil, err
	}

	if len(hard) == 0 {
		hard = nil
	}

	if len(soft) == 0 {
		soft = nil
	}

	return &v1.PodAntiAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution:  hard,
		PreferredDuringSchedulingIgnoredDuringExecution: soft,
	}, nil
}

func revertPodAffinity(affinities []PodAffinity) (*v1.PodAffinity, error) {
	if len(affinities) == 0 {
		return nil, nil
	}

	hard, soft, err := splitAndRevertPodAffinity(affinities)
	if err != nil {
		return nil, err
	}

	if len(hard) == 0 {
		hard = nil
	}

	if len(soft) == 0 {
		soft = nil
	}

	return &v1.PodAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution:  hard,
		PreferredDuringSchedulingIgnoredDuringExecution: soft,
	}, nil
}

func splitAndRevertPodAffinity(affinities []PodAffinity) (hard []v1.PodAffinityTerm, soft []v1.WeightedPodAffinityTerm, err error) {
	hard = []v1.PodAffinityTerm{}
	soft = []v1.WeightedPodAffinityTerm{}
	for _, affinity := range affinities {
		segs := strings.Split(affinity.Affinity, ":")
		l := len(segs)

		var term *v1.PodAffinityTerm
		if l < 1 {
			err = util.InvalidInstanceError(affinity)
			return
		} else {
			term, err = parsePodExprs(segs[0])
			if err != nil {
				return
			}

			term.TopologyKey = affinity.Topology
			term.Namespaces = affinity.Namespaces
		}

		if l < 2 {
			hard = append(hard, *term)
			continue
		} else {
			if segs[1] != "soft" {
				err = util.InvalidInstanceError(affinity)
				return
			}
		}

		var weight int64
		if l < 3 {
			weight = 1
		} else {
			weight, err = strconv.ParseInt(segs[2], 10, 32)
			if err != nil {
				return
			}
		}

		soft = append(soft, v1.WeightedPodAffinityTerm{
			Weight:          int32(weight),
			PodAffinityTerm: *term,
		})
	}

	return
}

func parsePodExprs(s string) (*v1.PodAffinityTerm, error) {
	labelSelector, err := expressions.ParseLabelSelector(s)
	if err != nil {
		return nil, err
	}

	return &v1.PodAffinityTerm{
		LabelSelector: labelSelector,
	}, nil
}

func revertNodeAffinity(affinities []string) (*v1.NodeAffinity, error) {
	if len(affinities) == 0 {
		return nil, nil
	}

	hard, soft, err := splitAndRevertNodeAffinity(affinities)
	if err != nil {
		return nil, err
	}

	if len(hard) == 0 {
		hard = nil
	}

	if len(soft) == 0 {
		soft = nil
	}

	return &v1.NodeAffinity{
		RequiredDuringSchedulingIgnoredDuringExecution:  &v1.NodeSelector{hard},
		PreferredDuringSchedulingIgnoredDuringExecution: soft,
	}, nil
}

func splitAndRevertNodeAffinity(affinities []string) (hard []v1.NodeSelectorTerm, soft []v1.PreferredSchedulingTerm, err error) {
	hard = []v1.NodeSelectorTerm{}
	soft = []v1.PreferredSchedulingTerm{}
	for _, affinity := range affinities {
		segs := strings.Split(affinity, ":")
		l := len(segs)

		var term *v1.NodeSelectorTerm
		if l < 1 {
			err = util.InvalidInstanceError(affinity)
			return
		} else {
			term, err = parseNodeExprs(segs[0])
			if err != nil {
				return
			}
		}

		if l < 2 {
			hard = append(hard, *term)
			continue
		} else {
			if segs[1] != "soft" {
				err = util.InvalidInstanceError(affinity)
				return
			}
		}

		var weight int64
		if l < 3 {
			weight = 1
		} else {
			weight, err = strconv.ParseInt(segs[2], 10, 32)
			if err != nil {
				return
			}
		}

		soft = append(soft, v1.PreferredSchedulingTerm{
			Weight:     int32(weight),
			Preference: *term,
		})
	}

	return
}

func parseNodeExprs(s string) (*v1.NodeSelectorTerm, error) {
	reqs := []v1.NodeSelectorRequirement{}
	segs := strings.Split(s, "&")
	for _, seg := range segs {
		expr, err := expressions.ParseExpr(seg, []string{"!=", "=", ">", "<"})
		if err != nil {
			return nil, err
		}

		if expr == nil {
			if seg[0] == '!' {
				reqs = append(reqs, v1.NodeSelectorRequirement{
					Key:      seg[1:],
					Operator: v1.NodeSelectorOpDoesNotExist,
				})
			} else {
				reqs = append(reqs, v1.NodeSelectorRequirement{
					Key:      seg,
					Operator: v1.NodeSelectorOpExists,
				})
			}

			continue
		}

		var op v1.NodeSelectorOperator
		switch expr.Op {
		case "=":
			op = v1.NodeSelectorOpIn
		case "!=":
			op = v1.NodeSelectorOpNotIn
		case ">":
			op = v1.NodeSelectorOpGt
		case "<":
			op = v1.NodeSelectorOpLt
		default:
			glog.Fatal("unreachable")
		}
		reqs = append(reqs, v1.NodeSelectorRequirement{
			Key:      expr.Key,
			Operator: op,
			Values:   expr.Values,
		})
	}

	return &v1.NodeSelectorTerm{reqs}, nil
}
