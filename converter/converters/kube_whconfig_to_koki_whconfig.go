package converters

import (
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"strings"
	"github.com/koki/short/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Convert_Kube_WebhookConfiguration_to_Koki_WebhookConfiguration(webhookConfig interface{}, kind string) (interface{}, error) {
	var err error
	switch kind {
	case types.MutatingKind:
		kokiWrapper := &types.MutatingWebhookConfigWrapper{}
		kokiConfig := &kokiWrapper.WebhookConfig
		kubeConfig := webhookConfig.(*admissionregv1beta1.MutatingWebhookConfiguration)
		convertMeta(kokiConfig, kubeConfig.TypeMeta, kubeConfig.ObjectMeta)
		webhooks, err := convertWebhooks(kubeConfig.Webhooks)
		if err != nil {
			return nil, err
		}
		kokiConfig.Webhooks = webhooks
		return kokiWrapper, nil
	case types.ValidatingKind:
		kokiWrapper := &types.ValidatingWebhookConfigWrapper{}
		kokiConfig := &kokiWrapper.WebhookConfig
		kubeConfig := webhookConfig.(*admissionregv1beta1.ValidatingWebhookConfiguration)
		convertMeta(kokiConfig, kubeConfig.TypeMeta, kubeConfig.ObjectMeta)
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

func convertMeta(kokiConfig *types.WebhookConfig, typeMeta metav1.TypeMeta, objectMeta metav1.ObjectMeta) (){
	kokiConfig.Name = objectMeta.Name
	kokiConfig.Namespace = objectMeta.Namespace
	kokiConfig.Version = typeMeta.APIVersion
	kokiConfig.Cluster = objectMeta.ClusterName
	kokiConfig.Labels = objectMeta.Labels
	kokiConfig.Annotations = objectMeta.Annotations
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
				if operation == admissionregv1beta1.Create {
					createOp = true
				}
				if operation == admissionregv1beta1.Update {
					updateOp = true
				}
				if operation == admissionregv1beta1.Delete {
					deleteOp = true
				}
				if operation == admissionregv1beta1.Connect {
					connectOp = true
				}
			}

			var kokiOperations []string
			if createOp && updateOp && deleteOp && connectOp {
				kokiOperations = append(kokiOperations, "*")
			} else {
				kokiOperations = append(kokiOperations, strings.Join(kokiOperationsArr,"|"))
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
	var kokiService = ""
	var kokiWebhookConfig = webhook.ClientConfig
	//construct service
	if kokiWebhookConfig.Service != nil {
		serviceReference := *kokiWebhookConfig.Service
		service := serviceReference.Name
		serviceNS := serviceReference.Namespace
		servicePath := *serviceReference.Path
		s1 := []string{serviceNS, service}
		s1Str := strings.Join(s1, "/")
		s2 := []string{s1Str, servicePath}
		kokiService = strings.Join(s2, ":")
	}
	//get selector
	selector, _, err := convertRSLabelSelector(webhook.NamespaceSelector, nil)

	var kokiURL = ""
	if kokiWebhookConfig.URL != nil {
		kokiURL = *kokiWebhookConfig.URL
	}
	//align the structure using above variables
 	kokiWebhook = types.Webhook {
		Name: name,
		Client: kokiURL,
		CaBundle: kokiWebhookConfig.CABundle,
		Service: kokiService,
		FailurePolicy: webhook.FailurePolicy,
		Selector: selector,
		Rules: rules,
	}
	return name, kokiWebhook, err
}


