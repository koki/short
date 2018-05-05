package converters

import (
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"

	"github.com/koki/short/types"
)

func Convert_Kube_ValidatingWebhookConfiguration_to_Koki_ValidatingWebhookConfiguration(kubeInitConfig *admissionregv1beta1.ValidatingWebhookConfiguration) (*types.ValidatingWebhookConfigWrapper, error) {
	var err error
	kokiWrapper := &types.ValidatingWebhookConfigWrapper{}
	kokiInitConfig := &kokiWrapper.ValidatingWebhookConfig

	kokiInitConfig.Name = kubeInitConfig.Name
	kokiInitConfig.Namespace = kubeInitConfig.Namespace
	kokiInitConfig.Version = kubeInitConfig.APIVersion
	kokiInitConfig.Cluster = kubeInitConfig.ClusterName
	kokiInitConfig.Labels = kubeInitConfig.Labels
	kokiInitConfig.Annotations = kubeInitConfig.Annotations

	rules, err := convertValidatingWebhooks(kubeInitConfig.Webhooks)
	if err != nil {
		return nil, err
	}
	kokiInitConfig.Rules = rules

	return kokiWrapper, nil
}

func convertValidatingWebhooks(webhooks []admissionregv1beta1.Webhook) (map[string][]types.ValidatingWebhookRuleWithOperations, error) {
	kokiRules := map[string][]types.ValidatingWebhookRuleWithOperations{}

	for i := range webhooks {
		webhook := webhooks[i]
		name, rules, err := convertValidatingWebhook(webhook)
		if err != nil {
			return nil, err
		}

		if len(rules) > 0 {
			kokiRules[name] = rules
		}
	}
	return kokiRules, nil
}

func convertValidatingWebhook(webhook admissionregv1beta1.Webhook) (name string, rules []types.ValidatingWebhookRuleWithOperations, err error) {
	name = webhook.Name

	if len(webhook.Rules) == 0 {
		return name, rules, err
	}

	for i := range webhook.Rules {
		rule := webhook.Rules[i]
		kokiRule := types.ValidatingWebhookRuleWithOperations {
			Groups:    rule.APIGroups,
			Versions:  rule.APIVersions,
			Operations: rule.Operations,
			Resources: rule.Resources,
		}
		rules = append(rules, kokiRule)
	}
	return name, rules, err
}
