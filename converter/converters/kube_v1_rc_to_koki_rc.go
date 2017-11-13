package converters

import (
	"reflect"

	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
)

func Convert_Kube_v1_ReplicationController_to_Koki_ReplicationController(kubeRC *v1.ReplicationController) (*types.ReplicationControllerWrapper, error) {
	var err error
	kokiRC := &types.ReplicationController{}

	kokiRC.Name = kubeRC.Name
	kokiRC.Namespace = kubeRC.Namespace
	kokiRC.Version = kubeRC.APIVersion
	kokiRC.Cluster = kubeRC.ClusterName
	kokiRC.Labels = kubeRC.Labels
	kokiRC.Annotations = kubeRC.Annotations

	kokiRC.Replicas = kubeRC.Spec.Replicas
	kokiRC.MinReadySeconds = kubeRC.Spec.MinReadySeconds
	kokiRC.Template, err = convertTemplate(kubeRC.Spec.Template, kubeRC.Spec.Selector, kubeRC.Name)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(kubeRC.Status, v1.ReplicationControllerStatus{}) {
		kokiRC.Status = &kubeRC.Status
	}

	return &types.ReplicationControllerWrapper{
		ReplicationController: *kokiRC,
	}, nil
}

func convertTemplate(kubeTemplate *v1.PodTemplateSpec, selector map[string]string, parentName string) (*types.Pod, error) {
	if kubeTemplate == nil {
		return nil, nil
	}

	kubePod := &v1.Pod{
		Spec: kubeTemplate.Spec,
	}

	kubePod.Name = kubeTemplate.Name
	kubePod.Namespace = kubeTemplate.Namespace
	kubePod.Labels = kubeTemplate.Labels
	kubePod.Annotations = kubeTemplate.Annotations

	if len(kubePod.Labels) == 0 {
		// If the template doesn't already specify a selector, try filling it from elsewhere.
		if len(selector) > 0 {
			kubePod.Labels = selector
		} else {
			kubePod.Labels = map[string]string{
				"koki.io/selector.name": parentName,
			}
		}
	}

	kokiPod, err := Convert_Kube_v1_Pod_to_Koki_Pod(kubePod)
	if err != nil {
		return nil, err
	}

	return &kokiPod.Pod, nil
}
