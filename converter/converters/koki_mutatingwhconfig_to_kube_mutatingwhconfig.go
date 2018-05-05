package converters

import (
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"

	"github.com/koki/short/types"
)

func Convert_Koki_MutatingWebhookConfiguration_to_Kube_MutatingWebhookConfiguration(initConfig *types.MutatingWebhookConfigWrapper) (*admissionregv1beta1.MutatingWebhookConfiguration, error) {
	kubeInitConfig := &admissionregv1beta1.MutatingWebhookConfiguration{}
	kokiInitConfig := &initConfig.MutatingWebhookConfig

	kubeInitConfig.Name = kokiInitConfig.Name
	kubeInitConfig.Namespace = kokiInitConfig.Namespace
	if len(kokiInitConfig.Version) == 0 {
		kubeInitConfig.APIVersion = "admissionregistration/v1beta1"
	} else {
		kubeInitConfig.APIVersion = kokiInitConfig.Version
	}
	kubeInitConfig.Kind = "MutatingWebhookConfiguration"
	kubeInitConfig.ClusterName = kokiInitConfig.Cluster
	kubeInitConfig.Labels = kokiInitConfig.Labels
	kubeInitConfig.Annotations = kokiInitConfig.Annotations

	kubeInitConfig.Webhooks = revertMWCs(kokiInitConfig.Rules)
	return kubeInitConfig, nil
}

func revertMWCs(kokiRules map[string][]types.MutatingWebhookRuleWithOperations) []admissionregv1beta1.Webhook {
	var kubeInitializers []admissionregv1beta1.Webhook

	for key := range kokiRules {
		kubeInitializer := admissionregv1beta1.Webhook{}

		kubeInitializer.Name = key
		kokiInitRules := kokiRules[key]
		kubeInitializer.Rules = revertMWC(kokiInitRules)

		kubeInitializers = append(kubeInitializers, kubeInitializer)
	}
	return kubeInitializers
}

func revertMWC(kokiRules []types.MutatingWebhookRuleWithOperations) []admissionregv1beta1.RuleWithOperations {
	var kubeRules []admissionregv1beta1.RuleWithOperations

	for i := range kokiRules {
		kokiRule := kokiRules[i]
		internalRule := admissionregv1beta1.Rule{
			APIGroups:  kokiRule.Groups,
			APIVersions:  kokiRule.Versions,
			Resources:  kokiRule.Resources,
		}
		kubeRule := admissionregv1beta1.RuleWithOperations {
			Operations:  kokiRule.Operations,
			Rule: internalRule,
		}
		kubeRules = append(kubeRules, kubeRule)
	}

	return kubeRules
}
