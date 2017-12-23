package converters

import (
	"fmt"

	settingsv1alpha1 "k8s.io/api/settings/v1alpha1"

	"github.com/koki/short/types"
)

func Convert_Kube_PodPreset_to_Koki_PodPreset(kubePodPreset *settingsv1alpha1.PodPreset) (*types.PodPresetWrapper, error) {
	var err error
	kokiWrapper := &types.PodPresetWrapper{}
	kokiPodPreset := &kokiWrapper.PodPreset

	kokiPodPreset.Name = kubePodPreset.Name
	kokiPodPreset.Namespace = kubePodPreset.Namespace
	kokiPodPreset.Version = kubePodPreset.APIVersion
	kokiPodPreset.Cluster = kubePodPreset.ClusterName
	kokiPodPreset.Labels = kubePodPreset.Labels
	kokiPodPreset.Annotations = kubePodPreset.Annotations

	err = convertPodPresetSpec(kubePodPreset.Spec, kokiPodPreset)
	if err != nil {
		return nil, err
	}

	return kokiWrapper, nil
}

func convertPodPresetSpec(kubeSpec settingsv1alpha1.PodPresetSpec, kokiPodPreset *types.PodPreset) error {
	if kokiPodPreset == nil {
		return fmt.Errorf("Writing to uninitialized pod preset pointer")
	}

	kokiPodPreset.Env = convertEnvVars(kubeSpec.Env, kubeSpec.EnvFrom)

	volumeMounts, err := convertVolumeMounts(kubeSpec.VolumeMounts)
	if err != nil {
		return err
	}
	kokiPodPreset.VolumeMounts = volumeMounts

	volumes, err := convertVolumes(kubeSpec.Volumes)
	if err != nil {
		return err
	}
	kokiPodPreset.Volumes = volumes

	kubeSelector := kubeSpec.Selector
	// Fill out the Selector
	selector, _, err := convertRSLabelSelector(&kubeSelector, nil)
	if err != nil {
		return err
	}

	if selector != nil && (selector.Labels != nil || selector.Shorthand != "") {
		kokiPodPreset.Selector = selector
	}

	return nil
}
