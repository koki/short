package converters

import (
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"

	"github.com/koki/short/types"
	"strings"
)

func Convert_Koki_MutatingWebhookConfiguration_to_Kube_MutatingWebhookConfiguration(mutatingWebhookConfig *types.MutatingWebhookConfigWrapper) (*admissionregv1beta1.MutatingWebhookConfiguration, error) {
	kubeMutatingWebhookConfig := &admissionregv1beta1.MutatingWebhookConfiguration{}
	kokiMutatingWebhookConfig := &mutatingWebhookConfig.MutatingWebhookConfig

	kubeMutatingWebhookConfig.Name = kokiMutatingWebhookConfig.Name
	kubeMutatingWebhookConfig.Namespace = kokiMutatingWebhookConfig.Namespace
	if len(kokiMutatingWebhookConfig.Version) == 0 {
		kubeMutatingWebhookConfig.APIVersion = "admissionregistration/v1beta1"
	} else {
		kubeMutatingWebhookConfig.APIVersion = kokiMutatingWebhookConfig.Version
	}
	kubeMutatingWebhookConfig.Kind = "MutatingWebhookConfiguration"
	kubeMutatingWebhookConfig.ClusterName = kokiMutatingWebhookConfig.Cluster
	kubeMutatingWebhookConfig.Labels = kokiMutatingWebhookConfig.Labels
	kubeMutatingWebhookConfig.Annotations = kokiMutatingWebhookConfig.Annotations

	kubeMutatingWebhookConfig.Webhooks = revertMWCs(kokiMutatingWebhookConfig.Webhooks)
	return kubeMutatingWebhookConfig, nil
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

func revertMWCRules(kokiRules []types.MutatingWebhookRuleWithOperations) []admissionregv1beta1.RuleWithOperations {
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
