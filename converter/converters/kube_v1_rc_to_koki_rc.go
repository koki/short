package converters

import (
	"reflect"

	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
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

	if kubeSpec.Template != nil {
		meta, template, err := convertTemplate(*kubeSpec.Template)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "pod template")
		}
		kokiRC.TemplateMetadata = meta
		kokiRC.PodTemplate = template
	}

	kokiRC.ReplicationControllerStatus, err = convertReplicationControllerStatus(kubeRC.Status)
	if err != nil {
		return nil, err
	}

	return &types.ReplicationControllerWrapper{
		ReplicationController: *kokiRC,
	}, nil
}

func convertReplicationControllerStatus(kubeStatus v1.ReplicationControllerStatus) (types.ReplicationControllerStatus, error) {
	conditions, err := convertReplicationControllerConditions(kubeStatus.Conditions)
	if err != nil {
		return types.ReplicationControllerStatus{}, err
	}

	replicasStatus := &types.ReplicationControllerReplicasStatus{
		Total:        kubeStatus.Replicas,
		FullyLabeled: kubeStatus.FullyLabeledReplicas,
		Ready:        kubeStatus.ReadyReplicas,
		Available:    kubeStatus.AvailableReplicas,
	}
	if reflect.DeepEqual(replicasStatus, &types.ReplicationControllerReplicasStatus{}) {
		replicasStatus = nil
	}

	return types.ReplicationControllerStatus{
		ObservedGeneration: kubeStatus.ObservedGeneration,
		Replicas:           replicasStatus,
		Conditions:         conditions,
	}, nil
}

func convertReplicationControllerConditions(kubeConditions []v1.ReplicationControllerCondition) ([]types.ReplicationControllerCondition, error) {
	if len(kubeConditions) == 0 {
		return nil, nil
	}

	kokiConditions := make([]types.ReplicationControllerCondition, len(kubeConditions))
	for i, condition := range kubeConditions {
		status, err := convertConditionStatus(condition.Status)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "replication-controller conditions[%d]", i)
		}
		conditionType, err := convertReplicationControllerConditionType(condition.Type)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "replication-controller conditions[%d]", i)
		}
		kokiConditions[i] = types.ReplicationControllerCondition{
			Type:               conditionType,
			Status:             status,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}

	return kokiConditions, nil
}

func convertReplicationControllerConditionType(kubeType v1.ReplicationControllerConditionType) (types.ReplicationControllerConditionType, error) {
	switch kubeType {
	case v1.ReplicationControllerReplicaFailure:
		return types.ReplicationControllerReplicaFailure, nil
	default:
		return types.ReplicationControllerReplicaFailure, serrors.InvalidValueErrorf(kubeType, "unrecognized replication-controller condition type")
	}
}

func convertTemplate(kubeTemplate v1.PodTemplateSpec) (*types.PodTemplateMeta, types.PodTemplate, error) {
	meta := convertPodObjectMeta(kubeTemplate.ObjectMeta)

	spec, err := convertPodSpec(kubeTemplate.Spec)
	if err != nil {
		return nil, types.PodTemplate{}, err
	}

	return meta, *spec, nil
}
