package converters

import (
	admissionregv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"strings"
	"github.com/koki/short/types"
	//"bytes"
)

func Convert_Kube_MutatingWebhookConfiguration_to_Koki_MutatingWebhookConfiguration(kubeMutatingConfig *admissionregv1beta1.MutatingWebhookConfiguration) (*types.MutatingWebhookConfigWrapper, error) {
	var err error
	kokiWrapper := &types.MutatingWebhookConfigWrapper{}
	kokiMutatingConfig := &kokiWrapper.MutatingWebhookConfig

	kokiMutatingConfig.Name = kubeMutatingConfig.Name
	kokiMutatingConfig.Namespace = kubeMutatingConfig.Namespace
	kokiMutatingConfig.Version = kubeMutatingConfig.APIVersion
	kokiMutatingConfig.Cluster = kubeMutatingConfig.ClusterName
	kokiMutatingConfig.Labels = kubeMutatingConfig.Labels
	kokiMutatingConfig.Annotations = kubeMutatingConfig.Annotations

	webhooks, err := convertMutatingWebhooks(kubeMutatingConfig.Webhooks)
	if err != nil {
		return nil, err
	}
	kokiMutatingConfig.Webhooks = webhooks
	return kokiWrapper, nil
}

func convertMutatingWebhooks(webhooks []admissionregv1beta1.Webhook) (map[string]types.Webhook, error) {
	kokiWebhooks := map[string]types.Webhook{}

	for i := range webhooks {
		webhook := webhooks[i]
		name, kokiWebhook, err := convertMutatingWebhook(webhook)
		if err != nil {
			return nil, err
		}

		kokiWebhooks[name] = kokiWebhook
	}
	return kokiWebhooks, nil
}

func convertMutatingWebhook(webhook admissionregv1beta1.Webhook) (name string, kokiWebhook types.Webhook, err error) {
	var rules []types.MutatingWebhookRuleWithOperations = nil
	name = webhook.Name
	if len(webhook.Rules) != 0 {
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

	//construct namespaceselector
	namespaceselector := webhook.NamespaceSelector
	matchLabels := namespaceselector.MatchLabels
	matchExpressions := namespaceselector.MatchExpressions
	var kokiNamespaceSelector map[string]string
	kokiNamespaceSelector = make(map[string]string)
	for key := range matchLabels {
		kokiNamespaceSelector[key] = matchLabels[key]
	}

	for _, matchExpression := range matchExpressions {
		kokiNamespaceSelector[matchExpression.Key] = strings.Join(matchExpression.Values[:], ",")
	}

 	kokiWebhook = types.Webhook {
		Name: name,
		Client: *kokiWebhookConfig.URL,
		//CaBundle: string(kokiWebhookConfig.CABundle[:]),
		//CaBundle: bytes.NewBuffer(kokiWebhookConfig.CABundle).String(),
		CaBundle: kokiWebhookConfig.CABundle,
		Service: kokiService,
		FailurePolicy: webhook.FailurePolicy,
		NSSelector: kokiNamespaceSelector,
		Rules: rules,
	}
	return name, kokiWebhook, err
}


