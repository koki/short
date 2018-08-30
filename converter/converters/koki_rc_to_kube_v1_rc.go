package converters

import (
	"reflect"

	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
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

	kubeSpec.Selector = kokiRC.Selector
	kubeSpec.Template, err = revertTemplate(kokiRC.TemplateMetadata, kokiRC.PodTemplate)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod template")
	}

	// Make sure there's at least one Label in the Template and the Selector.
	if kubeSpec.Template != nil {
		if len(kubeSpec.Template.Labels) == 0 {
			if len(kubeSpec.Selector) > 0 {
				kubeSpec.Template.Labels = kubeSpec.Selector
			}
		}
	}

	kubeRC.Status, err = revertReplicationControllerStatus(kokiRC.ReplicationControllerStatus)
	if err != nil {
		return nil, err
	}

	return kubeRC, nil
}

func revertReplicationControllerStatus(kokiStatus types.ReplicationControllerStatus) (v1.ReplicationControllerStatus, error) {
	conditions, err := revertReplicationControllerConditions(kokiStatus.Conditions)
	if err != nil {
		return v1.ReplicationControllerStatus{}, err
	}
	replicasStatus := types.ReplicationControllerReplicasStatus{}
	if kokiStatus.Replicas != nil {
		replicasStatus = *kokiStatus.Replicas
	}
	return v1.ReplicationControllerStatus{
		ObservedGeneration:   kokiStatus.ObservedGeneration,
		Replicas:             replicasStatus.Total,
		FullyLabeledReplicas: replicasStatus.FullyLabeled,
		ReadyReplicas:        replicasStatus.Ready,
		AvailableReplicas:    replicasStatus.Available,
		Conditions:           conditions,
	}, nil
}

func revertReplicationControllerConditions(kokiConditions []types.ReplicationControllerCondition) ([]v1.ReplicationControllerCondition, error) {
	if len(kokiConditions) == 0 {
		return nil, nil
	}

	kubeConditions := make([]v1.ReplicationControllerCondition, len(kokiConditions))
	for i, condition := range kokiConditions {
		status, err := revertConditionStatus(condition.Status)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "replica-set conditions[%d]", i)
		}
		conditionType, err := revertReplicationControllerConditionType(condition.Type)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "replica-set conditions[%d]", i)
		}
		kubeConditions[i] = v1.ReplicationControllerCondition{
			Type:               conditionType,
			Status:             status,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}

	return kubeConditions, nil
}

func revertReplicationControllerConditionType(kokiType types.ReplicationControllerConditionType) (v1.ReplicationControllerConditionType, error) {
	switch kokiType {
	case types.ReplicationControllerReplicaFailure:
		return v1.ReplicationControllerReplicaFailure, nil
	default:
		return v1.ReplicationControllerReplicaFailure, serrors.InvalidValueErrorf(kokiType, "unrecognized replica-set condition type")
	}
}

func revertTemplate(kokiMeta *types.PodTemplateMeta, kokiSpec types.PodTemplate) (*v1.PodTemplateSpec, error) {
	var hasMeta = kokiMeta != nil
	var hasSpec = !reflect.DeepEqual(kokiSpec, types.PodTemplate{})
	if !hasMeta && !hasSpec {
		return nil, nil
	}

	template := v1.PodTemplateSpec{}

	if hasMeta {
		template.ObjectMeta = revertPodObjectMeta(*kokiMeta)
	}

	if hasSpec {
		spec, err := revertPodSpec(kokiSpec)
		if err != nil {
			return nil, err
		}
		template.Spec = *spec
	}

	return &template, nil
}
