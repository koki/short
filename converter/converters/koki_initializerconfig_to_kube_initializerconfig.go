package converters

import (
	admissionregv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"

	"github.com/koki/short/types"
)

func Convert_Koki_InitializerConfig_to_Kube_InitializerConfig(initConfig *types.InitializerConfigWrapper) (*admissionregv1alpha1.InitializerConfiguration, error) {
	kubeInitConfig := &admissionregv1alpha1.InitializerConfiguration{}
	kokiInitConfig := &initConfig.InitializerConfig

	kubeInitConfig.Name = kokiInitConfig.Name
	kubeInitConfig.Namespace = kokiInitConfig.Namespace
	if len(kokiInitConfig.Version) == 0 {
		kubeInitConfig.APIVersion = "admissionregistration/v1alpha1"
	} else {
		kubeInitConfig.APIVersion = kokiInitConfig.Version
	}
	kubeInitConfig.Kind = "InitializerConfiguration"
	kubeInitConfig.ClusterName = kokiInitConfig.Cluster
	kubeInitConfig.Labels = kokiInitConfig.Labels
	kubeInitConfig.Annotations = kokiInitConfig.Annotations

	kubeInitConfig.Initializers = revertKokiRules(kokiInitConfig.Rules)

	return kubeInitConfig, nil
}

func revertKokiRules(kokiRules map[string][]types.InitializerRule) []admissionregv1alpha1.Initializer {
	var kubeInitializers []admissionregv1alpha1.Initializer

	for key := range kokiRules {
		kubeInitializer := admissionregv1alpha1.Initializer{}

		kubeInitializer.Name = key
		kokiInitRules := kokiRules[key]
		kubeInitializer.Rules = revertKokiRule(kokiInitRules)

		kubeInitializers = append(kubeInitializers, kubeInitializer)
	}

	return kubeInitializers
}

func revertKokiRule(kokiRules []types.InitializerRule) []admissionregv1alpha1.Rule {
	var kubeRules []admissionregv1alpha1.Rule

	for i := range kokiRules {
		kokiRule := kokiRules[i]
		kubeRule := admissionregv1alpha1.Rule{
			APIGroups:   kokiRule.Groups,
			APIVersions: kokiRule.Versions,
			Resources:   kokiRule.Resources,
		}
		kubeRules = append(kubeRules, kubeRule)
	}

	return kubeRules
}
