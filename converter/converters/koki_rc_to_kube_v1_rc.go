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
	kubeRC.APIVersion = kokiRC.Version
	kubeRC.Kind = "ReplicationController"
	kubeRC.ClusterName = kokiRC.Cluster
	kubeRC.Labels = kokiRC.Labels
	kubeRC.Annotations = kokiRC.Annotations

	kubeRC.Spec.Replicas = kokiRC.Replicas
	kubeRC.Spec.MinReadySeconds = kokiRC.MinReadySeconds
	kubeRC.Spec.Selector = kokiRC.PodLabels

	kubeRC.Spec.Template, err = revertTemplate(kokiRC.Template)
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
