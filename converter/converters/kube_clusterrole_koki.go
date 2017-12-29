package converters

import (
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/short/parser/expressions"
	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_ClusterRole_to_Koki(kube *rbac.ClusterRole) (*types.ClusterRoleWrapper, error) {
	var err error
	koki := &types.ClusterRole{}

	koki.Name = kube.Name
	koki.Namespace = kube.Namespace
	koki.Version = kube.APIVersion
	koki.Cluster = kube.ClusterName
	koki.Labels = kube.Labels
	koki.Annotations = kube.Annotations

	koki.Rules = convertPolicyRules(kube.Rules)
	koki.AggregationRule, err = convertClusterRoleAggregation(kube.AggregationRule)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "ClusterRole AggregationRule")
	}

	return &types.ClusterRoleWrapper{
		ClusterRole: *koki,
	}, nil
}

func convertPolicyRules(kubeRules []rbac.PolicyRule) []types.PolicyRule {
	kokiRules := make([]types.PolicyRule, len(kubeRules))
	for i, kubeRule := range kubeRules {
		kokiRules[i] = types.PolicyRule{
			Verbs:           kubeRule.Verbs,
			APIGroups:       kubeRule.APIGroups,
			Resources:       kubeRule.Resources,
			ResourceNames:   kubeRule.ResourceNames,
			NonResourceURLs: kubeRule.NonResourceURLs,
		}
	}

	return kokiRules
}

func convertClusterRoleAggregation(kubeAggregation *rbac.AggregationRule) ([]string, error) {
	if kubeAggregation == nil {
		return nil, nil
	}

	kokiSelectors, err := convertLabelSelectors(kubeAggregation.ClusterRoleSelectors)
	if err != nil {
		return nil, err
	}

	if len(kokiSelectors) == 0 {
		return nil, nil
	}

	return kokiSelectors, nil
}

func convertLabelSelectors(kubeSelectors []metav1.LabelSelector) ([]string, error) {
	kokiSelectors := make([]string, len(kubeSelectors))
	for i, kubeSelector := range kubeSelectors {
		kokiSelector, err := expressions.UnparseLabelSelector(&kubeSelector)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "[%d]", i)
		}

		if len(kokiSelector) > 0 {
			kokiSelectors[i] = kokiSelector
		}
	}

	return kokiSelectors, nil
}
