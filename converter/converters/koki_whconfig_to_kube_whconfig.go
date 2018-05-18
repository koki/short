package converters

import (
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/koki/short/types"
	"strings"
)

const apiVersion string = "admissionregistration/v1beta1"

func Convert_Koki_WebhookConfiguration_to_Kube_WebhookConfiguration(webhookConfig interface{}, kind string) (interface{}, error) {
	switch kind {
	case types.MutatingKind:
		kokiWebhookConfig := webhookConfig.(*types.MutatingWebhookConfigWrapper).WebhookConfig
		kubeWebhookConfig := &admissionregv1beta1.MutatingWebhookConfiguration{}
		convertMetaKokiToKube(&kubeWebhookConfig.TypeMeta, &kubeWebhookConfig.ObjectMeta, kokiWebhookConfig, kind)
		kubeWebhookConfig.Webhooks = revertMWCs(kokiWebhookConfig.Webhooks)
		return kubeWebhookConfig, nil
	case types.ValidatingKind:
		kokiWebhookConfig := webhookConfig.(*types.ValidatingWebhookConfigWrapper).WebhookConfig
		kubeWebhookConfig := &admissionregv1beta1.ValidatingWebhookConfiguration{}
		convertMetaKokiToKube(&kubeWebhookConfig.TypeMeta, &kubeWebhookConfig.ObjectMeta, kokiWebhookConfig, kind)
		kubeWebhookConfig.Webhooks = revertMWCs(kokiWebhookConfig.Webhooks)
		return kubeWebhookConfig, nil
	default:
		return nil, nil
	}
}

func convertMetaKokiToKube(typeMeta *metav1.TypeMeta, objectMeta *metav1.ObjectMeta, kokiWebhookConfig types.WebhookConfig, kind string) () {
	objectMeta.Name = kokiWebhookConfig.Name
	objectMeta.Namespace = kokiWebhookConfig.Namespace
	if len(kokiWebhookConfig.Version) == 0 {
		typeMeta.APIVersion = apiVersion
	} else {
		typeMeta.APIVersion = kokiWebhookConfig.Version
	}
	typeMeta.Kind = kind
	objectMeta.ClusterName = kokiWebhookConfig.Cluster
	objectMeta.Labels = kokiWebhookConfig.Labels
	objectMeta.Annotations = kokiWebhookConfig.Annotations
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

	//set servicereference
	kokiServiceStringSplit := strings.Split(kokiWebhook.Service, ":")
	kokiNamespaceStringSplit := strings.Split(kokiServiceStringSplit[0], "/")
	var kubeServiceReference admissionregv1beta1.ServiceReference
	if len(kokiNamespaceStringSplit) > 0 {
		kubeServiceReference.Namespace = kokiNamespaceStringSplit[0]
	}
	if len(kokiNamespaceStringSplit) > 1 {
		kubeServiceReference.Name = kokiNamespaceStringSplit[1]
	}
	if len(kokiServiceStringSplit) > 1 {
		kubeServiceReference.Path = &kokiServiceStringSplit[1]
	}

	//set client config
	var kubeWebhookClientConfig admissionregv1beta1.WebhookClientConfig
	if len(kokiWebhook.CaBundle) != 0 {
		kubeWebhookClientConfig.CABundle = kokiWebhook.CaBundle
	}
	if len(kokiWebhook.Client) != 0 {
		kubeWebhookClientConfig.URL = &kokiWebhook.Client
	}
	if len(kubeServiceReference.Name) != 0 || len(kubeServiceReference.Namespace) != 0 || kubeServiceReference.Path != nil {
		kubeWebhookClientConfig.Service = &kubeServiceReference
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

		kokiOperations := strings.Split(kokiRule.Operations[0], "|")
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
