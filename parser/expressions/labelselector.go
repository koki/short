package expressions

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/golang/glog"

	"github.com/koki/short/util"
)

func ParseLabelSelector(s string) (*metav1.LabelSelector, error) {
	if len(s) == 0 {
		return nil, nil
	}

	labels := map[string]string{}
	reqs := []metav1.LabelSelectorRequirement{}
	segs := strings.Split(s, "&")
	for _, seg := range segs {
		expr, err := ParseExpr(seg, []string{"!=", "="})
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

func UnparseLabelSelector(kubeSelector *metav1.LabelSelector) (string, error) {
	if kubeSelector == nil {
		return "", nil
	}

	selectorString := ""
	// parse through match labels first
	for k, v := range kubeSelector.MatchLabels {
		kokiExpr := fmt.Sprintf("%s=%s", k, v)
		if len(selectorString) > 0 {
			selectorString = fmt.Sprintf("%s&%s", selectorString, kokiExpr)
		} else {
			selectorString = kokiExpr
		}
	}

	// parse through match expressions now
	for i := range kubeSelector.MatchExpressions {
		expr := kubeSelector.MatchExpressions[i]
		value := strings.Join(expr.Values, ",")
		op, err := ConvertOperatorLabelSelector(expr.Operator)
		if err != nil {
			return "", err
		}
		kokiExpr := fmt.Sprintf("%s%s%s", expr.Key, op, value)
		if expr.Operator == metav1.LabelSelectorOpExists {
			kokiExpr = fmt.Sprintf("%s", expr.Key)
		}
		if expr.Operator == metav1.LabelSelectorOpDoesNotExist {
			kokiExpr = fmt.Sprintf("!%s", expr.Key)
		}

		if len(selectorString) > 0 {
			selectorString = fmt.Sprintf("%s&%s", selectorString, kokiExpr)
		} else {
			selectorString = kokiExpr
		}
	}

	return selectorString, nil
}

func ConvertOperatorLabelSelector(op metav1.LabelSelectorOperator) (string, error) {
	if op == "" {
		return "", nil
	}
	if op == metav1.LabelSelectorOpIn {
		return "=", nil
	}
	if op == metav1.LabelSelectorOpNotIn {
		return "!=", nil
	}
	if op == metav1.LabelSelectorOpExists {
		return "", nil
	}
	if op == metav1.LabelSelectorOpDoesNotExist {
		return "", nil
	}
	return "", util.InvalidInstanceError(op)
}
