package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
)

func Convert_Koki_ReplicationController_to_Kube_v1_ReplicationController(rc *types.ReplicationControllerWrapper) (*v1.ReplicationController, error) {
	var err error
	kubeRC := &v1.ReplicationController{}
	kokiRC := rc.ReplicationController

	kubeRC.Name = kokiRC.Name
	kubeRC.Namespace = kokiRC.Namespace
	if len(kokiRC.Version) == 0 {
		kubeRC.APIVersion = "v1"
	} else {
		kubeRC.APIVersion = kokiRC.Version
	}
	kubeRC.Kind = "ReplicationController"
	kubeRC.ClusterName = kokiRC.Cluster
	kubeRC.Labels = kokiRC.Labels
	kubeRC.Annotations = kokiRC.Annotations

	kubeSpec := &kubeRC.Spec
	kubeSpec.Replicas = kokiRC.Replicas
	kubeSpec.MinReadySeconds = kokiRC.MinReadySeconds

	// We won't repopulate kubeSpec.Selector because it's
	// defaulted to the Template's labels.
	kokiPod := kokiRC.GetTemplate()
	// Make sure there's at least one Label in the Template.
	if len(kokiPod.Labels) == 0 {
		kokiPod.Labels = map[string]string{
			"koki.io/selector.name": kokiRC.Name,
		}
	}
	kubeSpec.Template, err = revertTemplate(kokiPod)
	if err != nil {
		return nil, err
	}

	if kokiRC.Status != nil {
		kubeRC.Status = *kokiRC.Status
	}

	return kubeRC, nil
}

func revertTemplate(kokiPod *types.Pod) (*v1.PodTemplateSpec, error) {
	if kokiPod == nil {
		return nil, nil
	}

	kubePod, err := Convert_Koki_Pod_to_Kube_v1_Pod(&types.PodWrapper{Pod: *kokiPod})
	if err != nil {
		return nil, err
	}
	kubeTemplate := &v1.PodTemplateSpec{
		Spec: kubePod.Spec,
	}

	kubeTemplate.Name = kubePod.Name
	kubeTemplate.Namespace = kubePod.Namespace
	kubeTemplate.Labels = kubePod.Labels
	kubeTemplate.Annotations = kubePod.Annotations

	return kubeTemplate, nil
}
