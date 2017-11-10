package converters

import (
	"fmt"
	"strings"

	apps "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/glog"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Koki_ReplicaSet_to_Kube_v1_ReplicaSet(rc *types.ReplicaSetWrapper) (*apps.ReplicaSet, error) {
	var err error
	kubeRC := &apps.ReplicaSet{}
	kokiRC := rc.ReplicaSet

	kubeRC.Name = kokiRC.Name
	kubeRC.Namespace = kokiRC.Namespace
	kubeRC.APIVersion = kokiRC.Version
	kubeRC.Kind = "ReplicaSet"
	kubeRC.ClusterName = kokiRC.Cluster
	kubeRC.Labels = kokiRC.Labels
	kubeRC.Annotations = kokiRC.Annotations

	kubeRC.Spec.Replicas = kokiRC.Replicas
	kubeRC.Spec.MinReadySeconds = kokiRC.MinReadySeconds
	kubeRC.Spec.Selector, err = parseLabelSelector(kokiRC.PodSelector)
	if err != nil {
		return nil, err
	}

	kubeTemplate, err := revertTemplate(kokiRC.Template)
	if err != nil {
		return nil, err
	}

	if kubeTemplate == nil {
		return nil, util.TypeValueErrorf(kokiRC, "missing pod template")
	}

	kubeRC.Spec.Template = *kubeTemplate

	if kokiRC.Status != nil {
		kubeRC.Status = *kokiRC.Status
	}

	return kubeRC, nil
}

func parseLabelSelector(s string) (*metav1.LabelSelector, error) {
	if len(s) == 0 {
		return nil, nil
	}

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

		if len(expr.Values) == 1 {
			labels[expr.Key] = expr.Values[0]
		} else {
			reqs = append(reqs, metav1.LabelSelectorRequirement{
				Key:      expr.Key,
				Operator: op,
				Values:   expr.Values,
			})
		}
	}

	if len(labels) == 0 {
		labels = nil
	}

	if len(reqs) == 0 {
		reqs = nil
	}

	return &metav1.LabelSelector{
		MatchLabels:      labels,
		MatchExpressions: reqs,
	}, nil
}

// Expr is the generic AST format of a koki NodeSelectorRequirement or LabelSelectorRequirement
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
