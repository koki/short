package converters

import (
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_CronJob_to_Koki_CronJob(kubeCronJob runtime.Object) (*types.CronJobWrapper, error) {
	groupVersionKind := kubeCronJob.GetObjectKind().GroupVersionKind()
	groupVersionString := groupVersionKind.GroupVersion().String()
	groupVersionKind.Version = "v1beta1"
	groupVersionKind.Group = "batch"
	kubeCronJob.GetObjectKind().SetGroupVersionKind(groupVersionKind)

	// Serialize as batch/v1beta1
	b, err := yaml.Marshal(kubeCronJob)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, kubeCronJob, "couldn't serialize kube CronJob after setting apiVersion to batch/v1beta1")
	}

	// Deserialize the "generic" kube CronJob
	genericCronJob, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, string(b), "couldn't deserialize 'generic' kube CronJob")
	}

	if genericCronJob, ok := genericCronJob.(*batchv1beta1.CronJob); ok {
		kokiWrapper, err := Convert_Kube_batch_v1beta1_CronJob_to_Koki_CronJob(genericCronJob)
		if err != nil {
			return nil, err
		}

		kokiCronJob := &kokiWrapper.CronJob

		kokiCronJob.Version = groupVersionString
		return kokiWrapper, nil
	}

	return nil, serrors.InvalidInstanceErrorf(genericCronJob, "didn't deserialize 'generic' kube CronJob as batch/v1beta1.CronJob")
}

func Convert_Kube_batch_v1beta1_CronJob_to_Koki_CronJob(kubeCronJob *batchv1beta1.CronJob) (*types.CronJobWrapper, error) {
	kokiObj := &types.CronJobWrapper{}
	kokiCronJob := types.CronJob{}

	kokiCronJob.Version = kubeCronJob.APIVersion
	kokiCronJob.Name = kubeCronJob.Name
	kokiCronJob.Namespace = kubeCronJob.Namespace
	kokiCronJob.Cluster = kubeCronJob.ClusterName
	kokiCronJob.Labels = kubeCronJob.Labels
	kokiCronJob.Annotations = kubeCronJob.Annotations

	kubeSpec := &kubeCronJob.Spec

	jobSpec, err := convertJobSpec(kubeSpec.JobTemplate.Spec)
	if err != nil {
		return nil, err
	}
	kokiCronJob.JobTemplate = *jobSpec
	templateMetadata := convertPodObjectMeta(kubeSpec.JobTemplate.ObjectMeta)
	if templateMetadata != nil {
		kokiCronJob.TemplateMetadata = templateMetadata
	}

	kokiCronJob.Schedule = kubeSpec.Schedule
	kokiCronJob.Suspend = kubeSpec.Suspend
	concurrencyPolicy, err := convertConcurrencyPolicy(kubeSpec.ConcurrencyPolicy)
	if err != nil {
		return nil, err
	}
	kokiCronJob.ConcurrencyPolicy = concurrencyPolicy

	kokiCronJob.CronJobStatus.Active = kubeCronJob.Status.Active
	kokiCronJob.CronJobStatus.LastScheduled = kubeCronJob.Status.LastScheduleTime

	kokiObj.CronJob = kokiCronJob

	return kokiObj, nil
}

func convertConcurrencyPolicy(concurrencyPolicy batchv1beta1.ConcurrencyPolicy) (types.ConcurrencyPolicy, error) {
	if concurrencyPolicy == "" {
		return "", nil
	}
	if concurrencyPolicy == batchv1beta1.AllowConcurrent {
		return types.AllowConcurrent, nil
	} else if concurrencyPolicy == batchv1beta1.ForbidConcurrent {
		return types.ForbidConcurrent, nil
	} else if concurrencyPolicy == batchv1beta1.ReplaceConcurrent {
		return types.ReplaceConcurrent, nil
	}
	return "", serrors.InvalidValueErrorf(concurrencyPolicy, "unrecognized Concurreny Policy")
}
