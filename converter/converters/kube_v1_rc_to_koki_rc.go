package converters

import (
	"reflect"

	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Kube_v1_ReplicationController_to_Koki_ReplicationController(kubeRC *v1.ReplicationController) (*types.ReplicationControllerWrapper, error) {
	kokiRC := &types.ReplicationController{}

	kokiRC.Name = kubeRC.Name
	kokiRC.Namespace = kubeRC.Namespace
	kokiRC.Version = kubeRC.APIVersion
	kokiRC.Cluster = kubeRC.ClusterName
	kokiRC.Labels = kubeRC.Labels
	kokiRC.Annotations = kubeRC.Annotations

	kubeSpec := &kubeRC.Spec

	kokiRC.Replicas = kubeSpec.Replicas
	kokiRC.MinReadySeconds = kubeSpec.MinReadySeconds

	if kubeSpec.Template != nil {
		meta, template, err := convertTemplate(*kubeSpec.Template)
		if err != nil {
			return nil, util.ContextualizeErrorf(err, "pod template")
		}
		kokiRC.TemplateMetadata = meta
		kokiRC.PodTemplate = template
	}

	if !reflect.DeepEqual(kubeRC.Status, v1.ReplicationControllerStatus{}) {
		kokiRC.Status = &kubeRC.Status
	}

	return &types.ReplicationControllerWrapper{
		ReplicationController: *kokiRC,
	}, nil
}

func convertTemplate(kubeTemplate v1.PodTemplateSpec) (*types.PodTemplateMeta, types.PodTemplate, error) {
	meta := convertPodObjectMeta(kubeTemplate.ObjectMeta)

	spec, err := convertPodSpec(kubeTemplate.Spec)
	if err != nil {
		return nil, types.PodTemplate{}, err
	}

	return &meta, *spec, nil
}
