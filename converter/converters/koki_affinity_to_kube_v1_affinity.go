package converters

import (
	"fmt"
	"strconv"
	"strings"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/glog"
	"github.com/koki/short/types"
)

func Convert_Koki_Affinity_to_Kube_v1_Affinity(a []types.Affinity) (*v1.Affinity, error) {
	node, pod, antiPod, err := splitAffinities(a)
	if err != nil {
		return nil, err
	}

	n, err := nodeAffinity(node)
	if err != nil {
		return nil, err
	}

	p, err := podAffinity(pod)
	if err != nil {
		return nil, err
	}

	ap, err := podAntiAffinity(antiPod)
	if err != nil {
		return nil, err
	}

	if n == nil && p == nil && ap == nil {
		return nil, nil
	}

	return &v1.Affinity{
		NodeAffinity:    n,
		PodAffinity:     p,
		PodAntiAffinity: ap,
	}, nil
}

type PodAffinity struct {
	Affinity   string
	Topology   string
	Namespaces []string
}

func splitAffinities(as []types.Affinity) (node []string, pod, antiPod []PodAffinity, err error) {
	node = []string{}
	pod = []PodAffinity{}
	antiPod = []PodAffinity{}

	for _, a := range as {
		switch {
		case len(a.NodeAffinity) > 0:
			node = append(node, a.NodeAffinity)
		case len(a.PodAffinity) > 0:
			pod = append(pod, PodAffinity{
				Affinity:   a.PodAffinity,
				Topology:   a.Topology,
				Namespaces: a.Namespaces,
			})
		case len(a.PodAntiAffinity) > 0:
			antiPod = append(antiPod, PodAffinity{
				Affinity:   a.PodAntiAffinity,
				Topology:   a.Topology,
				Namespaces: a.Namespaces,
			})
		default:
			err = fmt.Errorf("unrecognized affinity %#v", a)
		}
	}

	return
}

func podAntiAffinity(as []PodAffinity) (*v1.PodAntiAffinity, error) {
	if len(as) == 0 {
		return nil, nil
	}

	hard, soft, err := podAffinities(as)
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

func podAffinity(as []PodAffinity) (*v1.PodAffinity, error) {
	if len(as) == 0 {
		return nil, nil
	}

	hard, soft, err := podAffinities(as)
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

func podAffinities(as []PodAffinity) (hard []v1.PodAffinityTerm, soft []v1.WeightedPodAffinityTerm, err error) {
	hard = []v1.PodAffinityTerm{}
	soft = []v1.WeightedPodAffinityTerm{}
	for _, a := range as {
		segs := strings.Split(a.Affinity, ":")
		l := len(segs)

		var term *v1.PodAffinityTerm
		if l < 1 {
			err = fmt.Errorf("unrecognized PodAffinity %s", a)
			return
		} else {
			term, err = podExprs(segs[0])
			if err != nil {
				return
			}

			term.TopologyKey = a.Topology
			term.Namespaces = a.Namespaces
		}

		if l < 2 {
			hard = append(hard, *term)
			continue
		} else {
			if segs[1] != "soft" {
				err = fmt.Errorf("unrecognized NodeAffinity term %s", a)
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

func podExprs(s string) (*v1.PodAffinityTerm, error) {
	labels := map[string]string{}
	reqs := []metav1.LabelSelectorRequirement{}
	segs := strings.Split(s, "&")
	for _, seg := range segs {
		expr, err := parseExpr(seg, []string{"!=", "="})
		if err != nil {
			return nil, err
		}

		if expr == nil {
			if seg[0] == '!' {
				reqs = append(reqs, metav1.LabelSelectorRequirement{
					Key:      seg[1:],
					Operator: metav1.LabelSelectorOpDoesNotExist,
				})
			} else {
				reqs = append(reqs, metav1.LabelSelectorRequirement{
					Key:      seg,
					Operator: metav1.LabelSelectorOpExists,
				})
			}

			continue
		}

		var op metav1.LabelSelectorOperator
		switch expr.Op {
		case "=":
			op = metav1.LabelSelectorOpIn
		case "!=":
			op = metav1.LabelSelectorOpNotIn
		default:
			glog.Fatal("unreachable")
		}
		reqs = append(reqs, metav1.LabelSelectorRequirement{
			Key:      expr.Key,
			Operator: op,
			Values:   expr.Values,
		})
	}

	if len(labels) == 0 {
		labels = nil
	}

	if len(reqs) == 0 {
		reqs = nil
	}

	return &v1.PodAffinityTerm{
		LabelSelector: &metav1.LabelSelector{
			MatchLabels:      labels,
			MatchExpressions: reqs,
		},
	}, nil
}

func nodeAffinity(as []string) (*v1.NodeAffinity, error) {
	if len(as) == 0 {
		return nil, nil
	}

	hard, soft, err := nodeAffinities(as)
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

func nodeAffinities(as []string) (hard []v1.NodeSelectorTerm, soft []v1.PreferredSchedulingTerm, err error) {
	hard = []v1.NodeSelectorTerm{}
	soft = []v1.PreferredSchedulingTerm{}
	for _, a := range as {
		segs := strings.Split(a, ":")
		l := len(segs)

		var term *v1.NodeSelectorTerm
		if l < 1 {
			err = fmt.Errorf("unrecognized NodeAffinity term %s", a)
			return
		} else {
			term, err = nodeExprs(segs[0])
			if err != nil {
				return
			}
		}

		if l < 2 {
			hard = append(hard, *term)
			continue
		} else {
			if segs[1] != "soft" {
				err = fmt.Errorf("unrecognized NodeAffinity term %s", a)
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

func nodeExprs(s string) (*v1.NodeSelectorTerm, error) {
	reqs := []v1.NodeSelectorRequirement{}
	segs := strings.Split(s, "&")
	for _, seg := range segs {
		expr, err := parseExpr(seg, []string{"!=", "=", ">", "<"})
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

type Expr struct {
	Key    string
	Op     string
	Values []string
}

func parseExpr(s string, ops []string) (*Expr, error) {
	for _, op := range ops {
		x, err := parseOp(s, op)
		if err != nil {
			return nil, err
		}

		if x != nil {
			return x, nil
		}
	}

	return nil, nil
}

func parseOp(s string, op string) (*Expr, error) {
	if strings.Contains(s, op) {
		segs := strings.Split(s, op)
		if len(segs) != 2 {
			return nil, fmt.Errorf(
				"Unrecognized expression (%s), op (%s)", s, op)
		}

		return &Expr{
			Key:    segs[0],
			Op:     op,
			Values: parseValues(segs[1]),
		}, nil
	}

	return nil, nil
}

func parseValues(s string) []string {
	return strings.Split(s, ",")
}
