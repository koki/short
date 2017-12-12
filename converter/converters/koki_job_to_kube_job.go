package converters

import (
	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
	batchv1 "k8s.io/api/batch/v1"
)

func Convert_Koki_Job_to_Kube_Job(job *types.JobWrapper) (interface{}, error) {
	kubeJob := &batchv1.Job{}
	kokiJob := job.Job

	kubeJob.Name = kokiJob.Name
	kubeJob.Namespace = kokiJob.Namespace
	kubeJob.APIVersion = kokiJob.Version
	kubeJob.Kind = "Job"
	kubeJob.ClusterName = kokiJob.Cluster
	kubeJob.Labels = kokiJob.Labels
	kubeJob.Annotations = kokiJob.Annotations

	spec, err := revertJobSpec(kokiJob.Name, kokiJob.JobTemplate)
	if err != nil {
		return nil, err
	}
	kubeJob.Spec = spec

	status, err := revertJobStatus(kokiJob)
	if err != nil {
		return nil, err
	}
	kubeJob.Status = status

	return kubeJob, nil
}

func revertJobSpec(name string, kokiJob types.JobTemplate) (batchv1.JobSpec, error) {
	kubeSpec := batchv1.JobSpec{}

	// Setting the Selector and Template is identical to ReplicaSet
	// Get the right Selector and Template Labels.
	var templateLabelsOverride map[string]string
	var kokiTemplateLabels map[string]string
	if kokiJob.TemplateMetadata != nil {
		kokiTemplateLabels = kokiJob.TemplateMetadata.Labels
	}
	var err error
	kubeSpec.Selector, templateLabelsOverride, err = revertRSSelector(name, kokiJob.Selector, kokiTemplateLabels)
	if err != nil {
		return batchv1.JobSpec{}, err
	}
	// Set the right Labels before we fill in the Pod template with this metadata.
	kokiJob.TemplateMetadata = applyTemplateLabelsOverride(templateLabelsOverride, kokiJob.TemplateMetadata)

	// Fill in the rest of the Pod template.
	kubeTemplate, err := revertTemplate(kokiJob.TemplateMetadata, kokiJob.PodTemplate)
	if err != nil {
		return batchv1.JobSpec{}, serrors.ContextualizeErrorf(err, "pod template")
	}
	if kubeTemplate == nil {
		return batchv1.JobSpec{}, serrors.InvalidInstanceErrorf(kokiJob, "missing pod template")
	}

	kubeSpec.Template = *kubeTemplate
	kubeSpec.Parallelism = kokiJob.Parallelism
	kubeSpec.Completions = kokiJob.Completions
	kubeSpec.BackoffLimit = kokiJob.MaxRetries
	kubeSpec.ManualSelector = kokiJob.ManualSelector
	kubeSpec.ActiveDeadlineSeconds = kokiJob.ActiveDeadlineSeconds

	return kubeSpec, nil
}

func revertJobStatus(kokiJob types.Job) (batchv1.JobStatus, error) {
	kubeJobStatus := batchv1.JobStatus{}
	kokiStatus := kokiJob.JobStatus

	conditions, err := revertJobConditions(kokiStatus.Conditions)
	if err != nil {
		return batchv1.JobStatus{}, err
	}

	kubeJobStatus.Conditions = conditions

	kubeJobStatus.StartTime = kokiStatus.StartTime
	kubeJobStatus.CompletionTime = kokiStatus.EndTime
	if kokiStatus.Running != nil {
		kubeJobStatus.Active = *kokiStatus.Running
	}
	if kokiStatus.Successful != nil {
		kubeJobStatus.Succeeded = *kokiStatus.Successful
	}
	if kokiStatus.Failed != nil {
		kubeJobStatus.Failed = *kokiStatus.Failed
	}

	return kubeJobStatus, nil
}

func revertJobConditions(kokiConditions []types.JobCondition) ([]batchv1.JobCondition, error) {
	if len(kokiConditions) == 0 {
		return nil, nil
	}

	kubeConditions := make([]batchv1.JobCondition, len(kokiConditions))
	for i, condition := range kokiConditions {
		status, err := revertConditionStatus(condition.Status)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "job conditions[%d]", i)
		}
		conditionType, err := revertJobConditionType(condition.Type)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "job conditions[%d]", i)
		}
		kubeConditions[i] = batchv1.JobCondition{
			Type:               conditionType,
			Status:             status,
			LastProbeTime:      condition.LastProbeTime,
			LastTransitionTime: condition.LastTransitionTime,
			Reason:             condition.Reason,
			Message:            condition.Message,
		}
	}
	return kubeConditions, nil
}

func revertJobConditionType(kokiType types.JobConditionType) (batchv1.JobConditionType, error) {
	switch kokiType {
	case types.JobComplete:
		return batchv1.JobComplete, nil
	case types.JobFailed:
		return batchv1.JobFailed, nil
	default:
		return batchv1.JobFailed, serrors.InvalidValueErrorf(kokiType, "unrecognized job condition type")
	}
}
