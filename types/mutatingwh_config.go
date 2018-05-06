package types

import (
	"strings"

	"github.com/koki/json"
	serrors "github.com/koki/structurederrors"
	"k8s.io/api/admissionregistration/v1beta1"
	//"k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MutatingWebhookConfigWrapper struct {
	MutatingWebhookConfig MutatingWebhookConfig `json:"mutating_webhook"`
}

type MutatingWebhookConfig struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Webhooks    map[string]Webhook `json:"webhooks,omitempty"`
}

type Webhook struct {
	Name           string            `json:"name,omitempty"`
	Client         string            `json:"client,omitempty"`
	//CaBundle       string            `json:"caBundle,omitempty"`
	CaBundle       []byte            `json:"caBundle,omitempty"`
	Service        string            `json:"service,omitempty"`
	FailurePolicy  *v1beta1.FailurePolicyType            `json:"on_fail,omitempty"`
	Rules          []MutatingWebhookRuleWithOperations `json:"rules,omitempty"`
	NSSelector     map[string]string            `json:"ns_selector,omitempty"`
}

type MutatingWebhookRuleWithOperations struct {
	Groups    []string		`json:"groups,omitempty"`
	Versions  []string		`json:"versions,omitempty"`
	Operations  []v1beta1.OperationType  `json:"operations,omitempty"`
	Resources []string		`json:"resources,omitempty"`
}

type LabelSelector struct {
	MatchLabels map[string]string `json:"matchLabels,omitempty" protobuf:"bytes,1,rep,name=matchLabels"`
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty" protobuf:"bytes,2,rep,name=matchExpressions"`
}

type LabelSelectorRequirement struct {
	Key string `json:"key" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,1,opt,name=key"`
	Operator string `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=LabelSelectorOperator"`
	Values []string `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
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
	//var ruleStruct MutatingWebhookRuleWithOperations
	var ruleStruct map[string][]string
	strErr2 := json.Unmarshal(data, &ruleStruct)
	if strErr2 != nil {
		return strErr2
	}

	var ok bool
	if i.Groups, ok = ruleStruct["groups"]; !ok{
		return serrors.InvalidInstanceError("couldn't parse JSON: Groups cannot be empty")
	}
	if i.Versions, ok = ruleStruct["versions"]; !ok{
		return serrors.InvalidInstanceError("couldn't parse JSON: Versions cannot be empty")
	}
	if i.Resources, ok = ruleStruct["resources"]; !ok{
		return serrors.InvalidInstanceError("couldn't parse JSON: Resources cannot be empty")
	}
	for _, operation := range(ruleStruct["operations"]) {
		switch operation {
		case "*":
			i.Operations = append(i.Operations, v1beta1.OperationAll)
		case "CREATE":
			i.Operations = append(i.Operations, v1beta1.Create)
		case "UPDATE":
			i.Operations = append(i.Operations, v1beta1.Update)
		case "DELETE":
			i.Operations = append(i.Operations, v1beta1.Delete)
		case "CONNECT":
			i.Operations = append(i.Operations, v1beta1.Connect)
		default:
			return serrors.InvalidInstanceError("couldn't parse JSON: Unknown Operation type")
		}
	}

	return nil
}

func (i MutatingWebhookRuleWithOperations) MarshalJSON() ([]byte, error) {
	if len(i.Resources) == 0 || len(i.Versions) == 0 || len(i.Groups) == 0 || len(i.Operations) == 0 {
		return []byte{}, serrors.InvalidInstanceErrorf(i, "Invalid Mutating Webhook Format")
	}

	var rule interface{}
	var ruleType string
	ruleType = "struct"

	if len(i.Resources) == 1 && len(i.Versions) == 1 && len(i.Groups) == 1 && len(i.Operations) == 1 {
		ruleString := strings.Join(append(append(i.Groups, i.Versions...), i.Resources...), "/")

		for _, operation := range(i.Operations) {
			ruleString = ruleString + "/" + string(operation)
		}
		rule = ruleString
		ruleType = "string"
	}
	if ruleType == "struct" {
		rule = map[string]interface{} {
			"groups" : i.Groups,
			"versions" : i.Versions,
			"resources" : i.Resources,
			"operations" : i.Operations,
		}
	}

	b, err := json.Marshal(rule)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, i, "marshalling mutating webhook rule %s to JSON", ruleType)
	}
	return b, nil
}
