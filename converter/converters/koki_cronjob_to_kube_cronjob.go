package converters

import (
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_CronJob_to_Kube_CronJob(cronJob *types.CronJobWrapper) (interface{}, error) {
	// Perform version-agnostic conversion into batch/v1beta1 CronJob.
	kubeCronJob, err := Convert_Koki_CronJob_to_Kube_batch_v1beta1_CronJob(cronJob)
	if err != nil {
		return nil, err
	}

	// Serialize the "generic" kube CronJob.
	b, err := yaml.Marshal(kubeCronJob)
	if err != nil {
		return nil, serrors.InvalidValueContextErrorf(err, kubeCronJob, "couldn't serialize 'generic' kube CronJob")
	}

	// Deserialize a versioned kube CronJob using its apiVersion.
	versionedCronJob, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, err
	}

	switch versionedCronJob := versionedCronJob.(type) {
	case *batchv1beta1.CronJob:
		// Perform batch/v1beta1-specific initialization here.
	case *batchv2alpha1.CronJob:
		// Perform batch/v2alpha-specific initialization here.
	default:
		return nil, serrors.TypeErrorf(versionedCronJob, "deserialized the manifest, but not as a supported kube CronJob")
	}

	return versionedCronJob, nil
}

func Convert_Koki_CronJob_to_Kube_batch_v1beta1_CronJob(kokiObj *types.CronJobWrapper) (batchv1beta1.CronJob, error) {
	kubeCronJob := batchv1beta1.CronJob{}
	kokiCronJob := kokiObj.CronJob

	kubeCronJob.Name = kokiCronJob.Name
	kubeCronJob.Namespace = kokiCronJob.Namespace
	kubeCronJob.APIVersion = kokiCronJob.Version
	kubeCronJob.Kind = "CronJob"
	kubeCronJob.ClusterName = kokiCronJob.Cluster
	kubeCronJob.Labels = kokiCronJob.Labels
	kubeCronJob.Annotations = kokiCronJob.Annotations

	spec, err := revertCronJobSpec(kokiCronJob)
	if err != nil {
		return kubeCronJob, err
	}
	kubeCronJob.Spec = spec

	kubeCronJob.Status.Active = kokiCronJob.Active
	kubeCronJob.Status.LastScheduleTime = kokiCronJob.LastScheduled

	return kubeCronJob, nil
}

func revertCronJobSpec(kokiCronJob types.CronJob) (batchv1beta1.CronJobSpec, error) {
	kubeSpec := batchv1beta1.CronJobSpec{}

	var objectMeta metav1.ObjectMeta

	if kokiCronJob.TemplateMetadata != nil {
		objectMeta = revertPodObjectMeta(*kokiCronJob.TemplateMetadata)
	}

	jobSpec, err := revertJobSpec(kokiCronJob.Name, kokiCronJob.JobTemplate)
	if err != nil {
		return kubeSpec, err
	}

	kubeSpec.JobTemplate = batchv1beta1.JobTemplateSpec{
		ObjectMeta: objectMeta,
		Spec:       jobSpec,
	}

	kubeSpec.Schedule = kokiCronJob.Schedule
	kubeSpec.StartingDeadlineSeconds = kokiCronJob.StartingDeadlineSeconds
	kubeSpec.Suspend = kokiCronJob.Suspend
	concurrencyPolicy, err := revertConcurrencyPolicy(kokiCronJob.ConcurrencyPolicy)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.ConcurrencyPolicy = concurrencyPolicy

	kubeSpec.SuccessfulJobsHistoryLimit = kokiCronJob.MaxSuccessHistory
	kubeSpec.FailedJobsHistoryLimit = kokiCronJob.MaxFailureHistory

	return kubeSpec, nil
}

func revertConcurrencyPolicy(concurrencyPolicy types.ConcurrencyPolicy) (batchv1beta1.ConcurrencyPolicy, error) {
	if concurrencyPolicy == "" {
		return "", nil
	}
	switch concurrencyPolicy {
	case types.AllowConcurrent:
		return batchv1beta1.AllowConcurrent, nil
	case types.ForbidConcurrent:
		return batchv1beta1.ForbidConcurrent, nil
	case types.ReplaceConcurrent:
		return batchv1beta1.ReplaceConcurrent, nil
	default:
		return "", serrors.InvalidValueErrorf(concurrencyPolicy, "unrecognized Concurreny Policy")
	}
}
