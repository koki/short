package types

import (
	"strings"

	"github.com/koki/json"
	serrors "github.com/koki/structurederrors"
	"k8s.io/api/admissionregistration/v1beta1"
)

type MutatingWebhookConfigWrapper struct {
	MutatingWebhookConfig MutatingWebhookConfig `json:"mutatingwh_config"`
}

type MutatingWebhookConfig struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Rules map[string][]MutatingWebhookRuleWithOperations `json:"rules,omitempty"`
}

type MutatingWebhookRule struct {
	Groups    []string
	Versions  []string
	Resources []string
}

type MutatingWebhookRuleWithOperations struct {
	Groups    []string
	Versions  []string
	Operations  []v1beta1.OperationType
	Resources []string
}

func (i *MutatingWebhookRuleWithOperations) UnmarshalJSON(data []byte) error {
	var ruleString string
	strErr1 := json.Unmarshal(data, &ruleString)
	if strErr1 == nil {
		parts := strings.SplitN(ruleString, "/", 3)
		if len(parts) < 2 {
			return serrors.InvalidValueForTypeErrorf(ruleString, i, "couldn't parse JSON: Invalid format")
		}

		if len(parts) == 2 {
			i.Groups = []string{""}
			i.Versions = []string{parts[0]}
			i.Resources = []string{parts[1]}
			return nil
		}

		i.Groups = []string{parts[0]}
		i.Versions = []string{parts[1]}
		i.Resources = []string{parts[2]}

		return nil
	}
	var ruleStruct MutatingWebhookRule
	strErr2 := json.Unmarshal(data, &ruleStruct)
	if strErr2 != nil {
		return strErr2
	}
	i.Groups = ruleStruct.Groups
	i.Versions = ruleStruct.Versions
	i.Resources = ruleStruct.Resources

	return nil
}

func (i MutatingWebhookRuleWithOperations) MarshalJSON() ([]byte, error) {
	if len(i.Resources) == 0 || len(i.Versions) == 0 {
		return []byte{}, nil
	}

	var rule interface{}
	var ruleType string

	rule = i
	ruleType = "struct"

	if len(i.Resources) == 1 && len(i.Versions) == 1 && len(i.Resources) == 1 {
		ruleString := strings.Join(append(append(i.Groups, i.Versions...), i.Resources...), "/")

		rule = ruleString
		ruleType = "string"
	}
	b, err := json.Marshal(rule)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, i, "marshalling mutating webhook rule %s to JSON", ruleType)
	}
	return b, nil
}
