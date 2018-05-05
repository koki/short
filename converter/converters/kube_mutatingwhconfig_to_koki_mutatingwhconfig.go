package converters

import (
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"

	"github.com/koki/short/types"
)

func Convert_Kube_MutatingWebhookConfiguration_to_Koki_MutatingWebhookConfiguration(kubeInitConfig *admissionregv1beta1.MutatingWebhookConfiguration) (*types.MutatingWebhookConfigWrapper, error) {
	var err error
	kokiWrapper := &types.MutatingWebhookConfigWrapper{}
	kokiInitConfig := &kokiWrapper.MutatingWebhookConfig

	kokiInitConfig.Name = kubeInitConfig.Name
	kokiInitConfig.Namespace = kubeInitConfig.Namespace
	kokiInitConfig.Version = kubeInitConfig.APIVersion
	kokiInitConfig.Cluster = kubeInitConfig.ClusterName
	kokiInitConfig.Labels = kubeInitConfig.Labels
	kokiInitConfig.Annotations = kubeInitConfig.Annotations

	rules, err := convertMutatingWebhooks(kubeInitConfig.Webhooks)
	if err != nil {
		return nil, err
	}
	kokiInitConfig.Rules = rules

	return kokiWrapper, nil
}

func convertMutatingWebhooks(webhooks []admissionregv1beta1.Webhook) (map[string][]types.MutatingWebhookRuleWithOperations, error) {
	kokiRules := map[string][]types.MutatingWebhookRuleWithOperations{}

	for i := range webhooks {
		webhook := webhooks[i]
		name, rules, err := convertMutatingWebhook(webhook)
		if err != nil {
			return nil, err
		}

		if len(rules) > 0 {
			kokiRules[name] = rules
		}
	}
	return kokiRules, nil
}

func convertMutatingWebhook(webhook admissionregv1beta1.Webhook) (name string, rules []types.MutatingWebhookRuleWithOperations, err error) {
	name = webhook.Name

	if len(webhook.Rules) == 0 {
		return name, rules, err
	}

	for i := range webhook.Rules {
		rule := webhook.Rules[i]
		kokiRule := types.MutatingWebhookRuleWithOperations{
			Groups:    rule.APIGroups,
			Versions:  rule.APIVersions,
			Operations: rule.Operations,
			Resources: rule.Resources,
		}
		rules = append(rules, kokiRule)
	}
	return name, rules, err
}
