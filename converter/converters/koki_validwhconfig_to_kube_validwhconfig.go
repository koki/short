package converters

import (
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	//admissionregv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"

	"github.com/koki/short/types"
)

func Convert_Koki_ValidatingWebhookConfiguration_to_Kube_ValidtingWebhookConfiguration(webHookConfig *types.ValidatingWebhookConfigWrapper) (*admissionregv1beta1.ValidatingWebhookConfiguration, error) {
	kubeWebHookConfig := &admissionregv1beta1.ValidatingWebhookConfiguration{}
	kokiWebHookConfig := &webHookConfig.ValidatingWebhookConfig

	kubeWebHookConfig.Name = kokiWebHookConfig.Name
	kubeWebHookConfig.Namespace = kokiWebHookConfig.Namespace
	if len(kokiWebHookConfig.Version) == 0 {
		kubeWebHookConfig.APIVersion = "admissionregistration/v1beta1"
	} else {
		kubeWebHookConfig.APIVersion = kokiWebHookConfig.Version
	}
	kubeWebHookConfig.Kind = "ValidatingWebhookConfiguration"
	kubeWebHookConfig.ClusterName = kokiWebHookConfig.Cluster
	kubeWebHookConfig.Labels = kokiWebHookConfig.Labels
	kubeWebHookConfig.Annotations = kokiWebHookConfig.Annotations

	kubeWebHookConfig.Webhooks = revertVWCs(kokiWebHookConfig.Rules)

	return kubeWebHookConfig, nil
}

func revertVWCs(kokiRules map[string][]types.ValidatingWebhookRuleWithOperations) []admissionregv1beta1.Webhook {
	var webHooks []admissionregv1beta1.Webhook

	for key := range kokiRules {
		webHook := admissionregv1beta1.Webhook{}

		webHook.Name = key
		kokiInitRules := kokiRules[key]
		webHook.Rules = revertVWC(kokiInitRules)
		webHooks = append(webHooks, webHook)
	}

	return webHooks
}

func revertVWC(kokiRules []types.ValidatingWebhookRuleWithOperations) []admissionregv1beta1.RuleWithOperations {
	var kubeRules []admissionregv1beta1.RuleWithOperations

	for i := range kokiRules {
		kokiRule := kokiRules[i]
		internalRule := admissionregv1beta1.Rule{
			APIGroups:  kokiRule.Groups,
			APIVersions:  kokiRule.Versions,
			Resources:  kokiRule.Resources,
		}
		kubeRule := admissionregv1beta1.RuleWithOperations {
			Operations:   kokiRule.Operations,
			Rule: internalRule,
		}
		kubeRules = append(kubeRules, kubeRule)
	}

	return kubeRules
}
