package types

import (
	"strings"

	"github.com/koki/json"
	serrors "github.com/koki/structurederrors"
	"k8s.io/api/admissionregistration/v1beta1"
)

const (
	MutatingKind string = "MutatingWebhookConfiguration"
	ValidatingKind string = "ValidatingWebhookConfiguration"
)

type MutatingWebhookConfigWrapper struct {
	WebhookConfig WebhookConfig `json:"mutating_webhook"`
}

type ValidatingWebhookConfigWrapper struct {
	WebhookConfig WebhookConfig `json:"validating_webhook"`
}

type WebhookConfig struct {
	Version     string             `json:"version,omitempty"`
	Cluster     string             `json:"cluster,omitempty"`
	Name        string             `json:"name,omitempty"`
	Namespace   string             `json:"namespace,omitempty"`
	Labels      map[string]string  `json:"labels,omitempty"`
	Annotations map[string]string  `json:"annotations,omitempty"`
	Webhooks    map[string]Webhook `json:"webhooks,omitempty"`
}

type Webhook struct {
	Name          string                      `json:"name,omitempty"`
	Client        string                      `json:"client,omitempty"`
	CaBundle      []byte                      `json:"caBundle,omitempty"`
	Service       string                      `json:"service,omitempty"`
	FailurePolicy *v1beta1.FailurePolicyType  `json:"on_fail,omitempty"`
	Rules         []WebhookRuleWithOperations `json:"rules,omitempty"`
	Selector      *RSSelector                 `json:"selector,omitempty"`
}

type WebhookRuleWithOperations struct {
	Groups     []string `json:"groups,omitempty"`
	Versions   []string `json:"versions,omitempty"`
	Operations []string   `json:"operations,omitempty"`
	Resources  []string `json:"resources,omitempty"`
}

type LabelSelector struct {
	MatchLabels      map[string]string          `json:"matchLabels,omitempty"`
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty"`
}

type LabelSelectorRequirement struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values,omitempty"`
}

func (i *WebhookRuleWithOperations) UnmarshalJSON(data []byte) error {
	var ruleString string

	strErr1 := json.Unmarshal(data, &ruleString)
	if strErr1 == nil {
		parts := strings.SplitN(ruleString, "/", 4)
		if len(parts) < 3 {
			return serrors.InvalidValueForTypeErrorf(ruleString, i, "couldn't parse JSON: Invalid format")
		}

		if len(parts) == 3 {
			i.Groups = []string{""}
			i.Versions = []string{parts[0]}
			i.Resources = []string{parts[1]}
			i.Operations = []string{parts[2]}
			return nil
		}

		i.Groups = []string{parts[0]}
		i.Versions = []string{parts[1]}
		i.Resources = []string{parts[2]}
		i.Operations = []string{parts[3]}

		return nil
	}
	var ruleStruct map[string][]string
	strErr2 := json.Unmarshal(data, &ruleStruct)
	if strErr2 != nil {
		return strErr2
	}

	var ok bool
	if i.Groups, ok = ruleStruct["groups"]; !ok {
		return serrors.InvalidInstanceError("couldn't parse JSON: Groups cannot be empty")
	}
	if i.Versions, ok = ruleStruct["versions"]; !ok {
		return serrors.InvalidInstanceError("couldn't parse JSON: Versions cannot be empty")
	}
	if i.Resources, ok = ruleStruct["resources"]; !ok {
		return serrors.InvalidInstanceError("couldn't parse JSON: Resources cannot be empty")
	}

	if len(ruleStruct["operations"]) > 0 {
		i.Operations = ruleStruct["operations"]
	}
	return nil
}

func (i WebhookRuleWithOperations) MarshalJSON() ([]byte, error) {
	if len(i.Resources) == 0 || len(i.Versions) == 0 || len(i.Groups) == 0 || len(i.Operations) == 0 {
		return []byte{}, serrors.InvalidInstanceErrorf(i, "Invalid Webhook Format")
	}

	var rule interface{}
	var ruleType string
	ruleType = "struct"

	if len(i.Resources) == 1 && len(i.Versions) == 1 && len(i.Groups) == 1 && len(i.Operations) == 1 {
		ruleString := strings.Join(append(append(append(i.Groups, i.Versions...), i.Resources...), i.Operations...), "/")
		rule = ruleString
		ruleType = "string"
	}
	if ruleType == "struct" {
		rule = map[string]interface{}{
			"groups":     i.Groups,
			"versions":   i.Versions,
			"resources":  i.Resources,
			"operations": i.Operations,
		}
	}

	b, err := json.Marshal(rule)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, i, "marshalling webhook rule %s to JSON", ruleType)
	}
	return b, nil
}
