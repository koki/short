package converters

import (
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"

	"github.com/koki/short/types"
	"strings"
)

func Convert_Koki_WebhookConfiguration_to_Kube_WebhookConfiguration(webhookConfig interface{}, kind string) (interface{}, error) {

	switch kind {
	case "MutatingWebhookConfiguration":
		kokiWebhookConfig := webhookConfig.(*types.MutatingWebhookConfigWrapper).WebhookConfig
		kubeWebhookConfig := &admissionregv1beta1.MutatingWebhookConfiguration{}
		kubeWebhookConfig.Name = kokiWebhookConfig.Name
		kubeWebhookConfig.Namespace = kokiWebhookConfig.Namespace
		if len(kokiWebhookConfig.Version) == 0 {
			kubeWebhookConfig.APIVersion = "admissionregistration/v1beta1"
		} else {
			kubeWebhookConfig.APIVersion = kokiWebhookConfig.Version
		}
		kubeWebhookConfig.Kind = kind
		kubeWebhookConfig.ClusterName = kokiWebhookConfig.Cluster
		kubeWebhookConfig.Labels = kokiWebhookConfig.Labels
		kubeWebhookConfig.Annotations = kokiWebhookConfig.Annotations
		kubeWebhookConfig.Webhooks = revertMWCs(kokiWebhookConfig.Webhooks)
		return kubeWebhookConfig, nil
	case "ValidatingWebhookConfiguration":
		kokiWebhookConfig := webhookConfig.(*types.ValidatingWebhookConfigWrapper).WebhookConfig
		kubeWebhookConfig := &admissionregv1beta1.ValidatingWebhookConfiguration{}
		kubeWebhookConfig.Name = kokiWebhookConfig.Name
		kubeWebhookConfig.Namespace = kokiWebhookConfig.Namespace
		if len(kokiWebhookConfig.Version) == 0 {
			kubeWebhookConfig.APIVersion = "admissionregistration/v1beta1"
		} else {
			kubeWebhookConfig.APIVersion = kokiWebhookConfig.Version
		}
		kubeWebhookConfig.Kind = kind
		kubeWebhookConfig.ClusterName = kokiWebhookConfig.Cluster
		kubeWebhookConfig.Labels = kokiWebhookConfig.Labels
		kubeWebhookConfig.Annotations = kokiWebhookConfig.Annotations
		kubeWebhookConfig.Webhooks = revertMWCs(kokiWebhookConfig.Webhooks)
		return kubeWebhookConfig, nil
	default:
		return nil, nil
	}



}

func revertMWCs(kokiWebhooks map[string]types.Webhook) []admissionregv1beta1.Webhook {
	var kubeWebhooks []admissionregv1beta1.Webhook

	for name := range kokiWebhooks {
		kubeWebhook := revertMWC(kokiWebhooks[name])
		kubeWebhooks = append(kubeWebhooks, kubeWebhook)
	}
	return kubeWebhooks
}

func revertMWC(kokiWebhook types.Webhook) admissionregv1beta1.Webhook {

	kokiServiceStringSplit := strings.Split(kokiWebhook.Service, ":")
	kokiNamespaceStringSplit := strings.Split(kokiServiceStringSplit[0], "/")

	kubeServiceReference := admissionregv1beta1.ServiceReference {
		Namespace: kokiNamespaceStringSplit[0],
		Name: kokiNamespaceStringSplit[1],
		Path: &kokiServiceStringSplit[1],
	}

	kubeWebhookClientConfig := admissionregv1beta1.WebhookClientConfig {
		CABundle: kokiWebhook.CaBundle,
		URL: &kokiWebhook.Client,
		Service: &kubeServiceReference,
	}

	kubeWebhookRules := revertMWCRules(kokiWebhook.Rules)
	kubeSelector, _, _ := revertRSSelector(kokiWebhook.Name, kokiWebhook.Selector, nil)
	kubeWebhook := admissionregv1beta1.Webhook {
		Name: kokiWebhook.Name,
		ClientConfig: kubeWebhookClientConfig,
		FailurePolicy: kokiWebhook.FailurePolicy,
		Rules: kubeWebhookRules,
		NamespaceSelector: kubeSelector,
	}

	return kubeWebhook
}

func revertMWCRules(kokiRules []types.WebhookRuleWithOperations) []admissionregv1beta1.RuleWithOperations {
	var kubeRules []admissionregv1beta1.RuleWithOperations

	for i := range kokiRules {
		kokiRule := kokiRules[i]
		internalRule := admissionregv1beta1.Rule {
			APIGroups:  kokiRule.Groups,
			APIVersions:  kokiRule.Versions,
			Resources:  kokiRule.Resources,
		}
		kokiOperations := strings.Split(kokiRule.Operations, "|")
		kubeOperations := []admissionregv1beta1.OperationType{}
		for _, operation := range(kokiOperations) {
			switch operation {
			case "*":
				kubeOperations = append(kubeOperations, admissionregv1beta1.OperationAll)
			case "CREATE":
				kubeOperations = append(kubeOperations, admissionregv1beta1.Create)
			case "UPDATE":
				kubeOperations = append(kubeOperations, admissionregv1beta1.Update)
			case "DELETE":
				kubeOperations = append(kubeOperations, admissionregv1beta1.Delete)
			case "CONNECT":
				kubeOperations = append(kubeOperations, admissionregv1beta1.Connect)
			}
		}

		kubeRule := admissionregv1beta1.RuleWithOperations {
			Operations:  kubeOperations,
			Rule: internalRule,
		}
		kubeRules = append(kubeRules, kubeRule)
	}

	return kubeRules
}
