package converters

import (
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"strings"
	"github.com/koki/short/types"
)

const (
	create string = "CREATE"
	update string = "UPDATE"
	delete string = "DELETE"
	connect string = "CONNECT"
)

func Convert_Kube_WebhookConfiguration_to_Koki_WebhookConfiguration(webhookConfig interface{}, kind string) (interface{}, error) {
	var err error
	switch kind {
	case "MutatingWebhookConfiguration":
		kokiWrapper := &types.MutatingWebhookConfigWrapper{}
		kokiConfig := &kokiWrapper.WebhookConfig
		kubeConfig := webhookConfig.(*admissionregv1beta1.MutatingWebhookConfiguration)
		kokiConfig.Name = kubeConfig.Name
		kokiConfig.Namespace = kubeConfig.Namespace
		kokiConfig.Version = kubeConfig.APIVersion
		kokiConfig.Cluster = kubeConfig.ClusterName
		kokiConfig.Labels = kubeConfig.Labels
		kokiConfig.Annotations = kubeConfig.Annotations

		webhooks, err := convertWebhooks(kubeConfig.Webhooks)
		if err != nil {
			return nil, err
		}
		kokiConfig.Webhooks = webhooks
		return kokiWrapper, nil
	case "ValidatingWebhookConfiguration":
		kokiWrapper := &types.ValidatingWebhookConfigWrapper{}
		kokiConfig := &kokiWrapper.WebhookConfig
		kubeConfig := webhookConfig.(*admissionregv1beta1.ValidatingWebhookConfiguration)
		kokiConfig.Name = kubeConfig.Name
		kokiConfig.Namespace = kubeConfig.Namespace
		kokiConfig.Version = kubeConfig.APIVersion
		kokiConfig.Cluster = kubeConfig.ClusterName
		kokiConfig.Labels = kubeConfig.Labels
		kokiConfig.Annotations = kubeConfig.Annotations

		webhooks, err := convertWebhooks(kubeConfig.Webhooks)
		if err != nil {
			return nil, err
		}
		kokiConfig.Webhooks = webhooks
		return kokiWrapper, nil
	default:
		return nil, err
	}
}

func convertWebhooks(webhooks []admissionregv1beta1.Webhook) (map[string]types.Webhook, error) {
	kokiWebhooks := map[string]types.Webhook{}

	for i := range webhooks {
		webhook := webhooks[i]
		name, kokiWebhook, err := convertWebhook(webhook)
		if err != nil {
			return nil, err
		}

		kokiWebhooks[name] = kokiWebhook
	}
	return kokiWebhooks, nil
}

func convertWebhook(webhook admissionregv1beta1.Webhook) (name string, kokiWebhook types.Webhook, err error) {
	var rules []types.WebhookRuleWithOperations = nil
	name = webhook.Name
	if len(webhook.Rules) != 0 {
		for i := range webhook.Rules {
			rule := webhook.Rules[i]

			kokiOperationsArr := []string{}
			var createOp bool = false
			var updateOp bool = false
			var deleteOp bool = false
			var connectOp bool = false
			for _, operation := range(rule.Operations) {
				kokiOperationsArr = append(kokiOperationsArr, string(operation))
				if string(operation) == create {
					createOp = true
				}
				if string(operation) == update {
					updateOp = true
				}
				if string(operation) == delete {
					deleteOp = true
				}
				if string(operation) == connect {
					connectOp = true
				}
			}

			var kokiOperations string = ""
			if createOp && updateOp && deleteOp && connectOp {
				kokiOperations = "*"
			} else {
				kokiOperations = strings.Join(kokiOperationsArr,"|")
			}

			kokiRule := types.WebhookRuleWithOperations{
				Groups:    rule.APIGroups,
				Versions:  rule.APIVersions,
				Operations: kokiOperations,
				Resources: rule.Resources,
			}
			rules = append(rules, kokiRule)
		}
	}
	var kokiWebhookConfig = webhook.ClientConfig

	//construct service
	serviceReference := *kokiWebhookConfig.Service
	service := serviceReference.Name
	serviceNS := serviceReference.Namespace
	servicePath := *serviceReference.Path
	s1 := []string{serviceNS, service}
	s1Str := strings.Join(s1, "/")
	s2 := []string{s1Str, servicePath}
	kokiService := strings.Join(s2, ":")

	//get selector
	selector, _, err := convertRSLabelSelector(webhook.NamespaceSelector, nil)

	//align the structure using above variables
 	kokiWebhook = types.Webhook {
		Name: name,
		Client: *kokiWebhookConfig.URL,
		CaBundle: kokiWebhookConfig.CABundle,
		Service: kokiService,
		FailurePolicy: webhook.FailurePolicy,
		Selector: selector,
		Rules: rules,
	}
	return name, kokiWebhook, err
}


