package converters

import (
	settingsv1alpha1 "k8s.io/api/settings/v1alpha1"

	"github.com/koki/short/types"
)

func Convert_Koki_PodPreset_to_Kube_PodPreset(podPreset *types.PodPresetWrapper) (*settingsv1alpha1.PodPreset, error) {
	kubePodPreset := &settingsv1alpha1.PodPreset{}
	kokiPodPreset := &podPreset.PodPreset

	kubePodPreset.Name = kokiPodPreset.Name
	kubePodPreset.Namespace = kokiPodPreset.Namespace
	if len(kokiPodPreset.Version) == 0 {
		kubePodPreset.APIVersion = "settings.k8s.io/v1alpha1"
	} else {
		kubePodPreset.APIVersion = kokiPodPreset.Version
	}
	kubePodPreset.Kind = "PodPreset"
	kubePodPreset.ClusterName = kokiPodPreset.Cluster
	kubePodPreset.Labels = kokiPodPreset.Labels
	kubePodPreset.Annotations = kokiPodPreset.Annotations

	spec, err := revertPodPreset(kokiPodPreset)
	if err != nil {
		return nil, err
	}
	kubePodPreset.Spec = spec

	return kubePodPreset, nil
}

func revertPodPreset(kokiPodPreset *types.PodPreset) (settingsv1alpha1.PodPresetSpec, error) {
	var kubeSpec settingsv1alpha1.PodPresetSpec

	volumes, err := revertVolumes(kokiPodPreset.Volumes)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.Volumes = volumes

	kubeSpec.VolumeMounts = revertVolumeMounts(kokiPodPreset.VolumeMounts)

	envs, envFroms, err := revertEnv(kokiPodPreset.Env)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.Env = envs
	kubeSpec.EnvFrom = envFroms

	selector, _, err := revertRSSelector(kokiPodPreset.Name, kokiPodPreset.Selector, nil)
	if err != nil {
		return kubeSpec, err
	}
	if selector != nil {
		kubeSpec.Selector = *selector
	}

	return kubeSpec, nil
}
