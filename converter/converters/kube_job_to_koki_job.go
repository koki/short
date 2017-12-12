package converters

import (
	batchv1 "k8s.io/api/batch/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_Job_to_Koki_Job(kubeJob *batchv1.Job) (*types.JobWrapper, error) {
	kokiObj := &types.JobWrapper{}
	kokiJob := types.Job{}

	kokiJob.Version = kubeJob.APIVersion
	kokiJob.Name = kubeJob.Name
	kokiJob.Namespace = kubeJob.Namespace
	kokiJob.Cluster = kubeJob.ClusterName
	kokiJob.Labels = kubeJob.Labels
	kokiJob.Annotations = kubeJob.Annotations

	jobSpec, err := convertJobSpec(kubeJob.Spec)
	if err != nil {
		return nil, err
	}
	kokiJob.JobTemplate = *jobSpec

	status, err := convertJobStatus(kubeJob.Status)
	if err != nil {
		return nil, err
	}
	kokiJob.JobStatus = status

	kokiObj.Job = kokiJob

	return kokiObj, nil
}

func convertJobSpec(kubeSpec batchv1.JobSpec) (*types.JobTemplate, error) {
	kokiJob := &types.JobTemplate{}

	selector, templateLabelsOverride, err := convertRSLabelSelector(kubeSpec.Selector, kubeSpec.Template.Labels)
	if err != nil {
		return nil, err
	}

	if selector != nil && (selector.Labels != nil || selector.Shorthand != "") {
		kokiJob.Selector = selector
	}

	// Build a Pod from the kube Template. Use it to set the koki Template.
	meta, template, err := convertTemplate(kubeSpec.Template)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod template")
	}
	kokiJob.TemplateMetadata = applyTemplateLabelsOverride(templateLabelsOverride, meta)
	kokiJob.PodTemplate = template

	kokiJob.Parallelism = kubeSpec.Parallelism
	kokiJob.Completions = kubeSpec.Completions
	kokiJob.MaxRetries = kubeSpec.BackoffLimit
	kokiJob.ActiveDeadlineSeconds = kubeSpec.ActiveDeadlineSeconds
	kokiJob.ManualSelector = kubeSpec.ManualSelector

	return kokiJob, nil
}

func convertJobStatus(kubeStatus batchv1.JobStatus) (types.JobStatus, error) {
	conditions, err := convertJobConditions(kubeStatus.Conditions)
	if err != nil {
		return types.JobStatus{}, err
	}

	var running *int32
	if kubeStatus.Active != 0 {
		running = &kubeStatus.Active
	}

	var successful *int32
	if kubeStatus.Succeeded != 0 {
		successful = &kubeStatus.Succeeded
	}

	var failed *int32
	if kubeStatus.Failed != 0 {
		failed = &kubeStatus.Failed
	}
	return types.JobStatus{
		Conditions: conditions,
		StartTime:  kubeStatus.StartTime,
		EndTime:    kubeStatus.CompletionTime,
		Running:    running,
		Successful: successful,
		Failed:     failed,
	}, nil
}

func convertJobConditions(kubeConditions []batchv1.JobCondition) ([]types.JobCondition, error) {
	if len(kubeConditions) == 0 {
		return nil, nil
	}

	kokiConditions := make([]types.JobCondition, len(kubeConditions))
	for i, condition := range kubeConditions {
		status, err := convertConditionStatus(condition.Status)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "job conditions[%d]", i)
		}
		conditionType, err := convertJobConditionType(condition.Type)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "job conditions[%d]", i)
		}
		kokiConditions[i] = types.JobCondition{
			Type:               conditionType,
			Status:             status,
			LastProbeTime:      condition.LastProbeTime,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}

	return kokiConditions, nil
}

func convertJobConditionType(kubeType batchv1.JobConditionType) (types.JobConditionType, error) {
	switch kubeType {
	case batchv1.JobComplete:
		return types.JobComplete, nil
	case batchv1.JobFailed:
		return types.JobFailed, nil
	default:
		return types.JobFailed, serrors.InvalidValueErrorf(kubeType, "unrecognized job condition type")
	}
}
