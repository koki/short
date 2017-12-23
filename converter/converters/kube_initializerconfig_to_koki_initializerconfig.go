package converters

import (
	admissionregv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"

	"github.com/koki/short/types"
)

func Convert_Kube_InitializerConfig_to_Koki_InitializerConfig(kubeInitConfig *admissionregv1alpha1.InitializerConfiguration) (*types.InitializerConfigWrapper, error) {
	var err error
	kokiWrapper := &types.InitializerConfigWrapper{}
	kokiInitConfig := &kokiWrapper.InitializerConfig

	kokiInitConfig.Name = kubeInitConfig.Name
	kokiInitConfig.Namespace = kubeInitConfig.Namespace
	kokiInitConfig.Version = kubeInitConfig.APIVersion
	kokiInitConfig.Cluster = kubeInitConfig.ClusterName
	kokiInitConfig.Labels = kubeInitConfig.Labels
	kokiInitConfig.Annotations = kubeInitConfig.Annotations

	rules, err := convertInitializers(kubeInitConfig.Initializers)
	if err != nil {
		return nil, err
	}
	kokiInitConfig.Rules = rules

	return kokiWrapper, nil
}

func convertInitializers(initializers []admissionregv1alpha1.Initializer) (map[string][]types.InitializerRule, error) {
	kokiRules := map[string][]types.InitializerRule{}

	for i := range initializers {
		initializer := initializers[i]
		name, rules, err := convertInitializer(initializer)
		if err != nil {
			return nil, err
		}

		if len(rules) > 0 {
			kokiRules[name] = rules
		}
	}
	return kokiRules, nil
}

func convertInitializer(initializer admissionregv1alpha1.Initializer) (name string, rules []types.InitializerRule, err error) {
	name = initializer.Name

	if len(initializer.Rules) == 0 {
		return name, rules, err
	}

	for i := range initializer.Rules {
		rule := initializer.Rules[i]
		kokiRule := types.InitializerRule{
			Groups:    rule.APIGroups,
			Versions:  rule.APIVersions,
			Resources: rule.Resources,
		}
		rules = append(rules, kokiRule)
	}
	return name, rules, err
}
