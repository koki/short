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

	kubeSpec := &kubeRC.Spec

	kokiRC.Replicas = kubeSpec.Replicas
	kokiRC.MinReadySeconds = kubeSpec.MinReadySeconds

	kokiPod, err := convertRSTemplate(kubeSpec.Template)
	if err != nil {
		return nil, err
	}
	kokiRC.SetTemplate(kokiPod)

	if !reflect.DeepEqual(kubeRC.Status, v1.ReplicationControllerStatus{}) {
		kokiRC.Status = &kubeRC.Status
	}

	return &types.ReplicationControllerWrapper{
		ReplicationController: *kokiRC,
	}, nil
}
