package converters

import (
	rbac "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/short/parser/expressions"
	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_ClusterRole_to_Kube(wrapper *types.ClusterRoleWrapper) (*rbac.ClusterRole, error) {
	var err error
	kube := &rbac.ClusterRole{}
	koki := &wrapper.ClusterRole

	kube.Name = koki.Name
	kube.Namespace = koki.Namespace
	if len(koki.Version) == 0 {
		kube.APIVersion = "rbac.authorization.k8s.io/v1"
	} else {
		kube.APIVersion = koki.Version
	}
	kube.Kind = "ClusterRole"
	kube.ClusterName = koki.Cluster
	kube.Labels = koki.Labels
	kube.Annotations = koki.Annotations

	kube.Rules = revertPolicyRules(koki.Rules)
	kube.AggregationRule, err = revertClusterRoleAggregation(koki.AggregationRule)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "cluster_role aggregation")
	}

	return kube, nil
}

func revertPolicyRules(kokiRules []types.PolicyRule) []rbac.PolicyRule {
	kubeRules := make([]rbac.PolicyRule, len(kokiRules))
	for i, kokiRule := range kokiRules {
		kubeRules[i] = rbac.PolicyRule{
			Verbs:           kokiRule.Verbs,
			APIGroups:       kokiRule.APIGroups,
			Resources:       kokiRule.Resources,
			ResourceNames:   kokiRule.ResourceNames,
			NonResourceURLs: kokiRule.NonResourceURLs,
		}
	}

	return kubeRules
}

func revertClusterRoleAggregation(kokiAggregation []string) (*rbac.AggregationRule, error) {
	kubeSelectors, err := revertLabelSelectors(kokiAggregation)
	if err != nil {
		return nil, err
	}

	if len(kokiAggregation) == 0 {
		return nil, nil
	}

	return &rbac.AggregationRule{
		ClusterRoleSelectors: kubeSelectors,
	}, nil
}

func revertLabelSelectors(kokiSelectors []string) ([]metav1.LabelSelector, error) {
	kubeSelectors := make([]metav1.LabelSelector, len(kokiSelectors))
	for i, kokiSelector := range kokiSelectors {
		kubeSelector, err := expressions.ParseLabelSelector(kokiSelector)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "[%d]", i)
		}

		if kubeSelector != nil {
			kubeSelectors[i] = *kubeSelector
		}
	}

	return kubeSelectors, nil
}
