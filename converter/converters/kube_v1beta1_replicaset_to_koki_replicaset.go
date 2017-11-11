package converters

import (
	"fmt"
	"reflect"
	"strings"

	exts "k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/short/types"
)

func Convert_Kube_v1beta2_ReplicaSet_to_Koki_ReplicaSet(kubeRS *exts.ReplicaSet) (*types.ReplicaSetWrapper, error) {
	var err error
	kokiRS := &types.ReplicaSet{}

	kokiRS.Name = kubeRS.Name
	kokiRS.Namespace = kubeRS.Namespace
	kokiRS.Version = kubeRS.APIVersion
	kokiRS.Cluster = kubeRS.ClusterName
	kokiRS.Labels = kubeRS.Labels
	kokiRS.Annotations = kubeRS.Annotations

	kokiRS.Replicas = kubeRS.Spec.Replicas
	kokiRS.MinReadySeconds = kubeRS.Spec.MinReadySeconds
	kokiRS.PodSelector, err = convertLabelSelector(kubeRS.Spec.Selector)
	if err != nil {
		return nil, err
	}

	kokiRS.Template, err = convertTemplate(&kubeRS.Spec.Template)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(kubeRS.Status, exts.ReplicaSetStatus{}) {
		kokiRS.Status = &kubeRS.Status
	}

	return &types.ReplicaSetWrapper{
		ReplicaSet: *kokiRS,
	}, nil
}

func convertLabelSelector(kubeSelector *metav1.LabelSelector) (string, error) {
	if kubeSelector == nil {
		return "", nil
	}

	affinityString := ""
	// parse through match labels first
	for k, v := range kubeSelector.MatchLabels {
		kokiExpr := fmt.Sprintf("%s=%s", k, v)
		if len(affinityString) > 0 {
			affinityString = fmt.Sprintf("%s&%s", affinityString, kokiExpr)
		} else {
			affinityString = kokiExpr
		}
	}

	// parse through match expressions now
	for i := range kubeSelector.MatchExpressions {
		expr := kubeSelector.MatchExpressions[i]
		value := strings.Join(expr.Values, ",")
		op, err := convertOperatorLabelSelector(expr.Operator)
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

		if len(affinityString) > 0 {
			affinityString = fmt.Sprintf("%s&%s", affinityString, kokiExpr)
		} else {
			affinityString = kokiExpr
		}
	}

	return affinityString, nil
}
