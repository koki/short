package converter

import (
	"github.com/koki/short/converter/converters"
	"github.com/koki/short/types"
	"github.com/koki/short/util"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	"k8s.io/api/core/v1"
	exts "k8s.io/api/extensions/v1beta1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

func DetectAndConvertFromKokiObj(kokiObj interface{}) (interface{}, error) {
	switch kokiObj := kokiObj.(type) {
	case *types.ConfigMapWrapper:
		return converters.Convert_Koki_ConfigMap_to_Kube_v1_ConfigMap(kokiObj)
	case *types.CronJobWrapper:
		return converters.Convert_Koki_CronJob_to_Kube_CronJob(kokiObj)
	case *types.DaemonSetWrapper:
		return converters.Convert_Koki_DaemonSet_to_Kube_DaemonSet(kokiObj)
	case *types.DeploymentWrapper:
		return converters.Convert_Koki_Deployment_to_Kube_Deployment(kokiObj)
	case *types.EndpointsWrapper:
		return converters.Convert_Koki_Endpoints_to_Kube_v1_Endpoints(kokiObj)
	case *types.IngressWrapper:
		return converters.Convert_Koki_Ingress_to_Kube_Ingress(kokiObj)
	case *types.JobWrapper:
		return converters.Convert_Koki_Job_to_Kube_Job(kokiObj)
	case *types.PersistentVolumeClaimWrapper:
		return converters.Convert_Koki_PVC_to_Kube_PVC(kokiObj)
	case *types.PersistentVolumeWrapper:
		return converters.Convert_Koki_PersistentVolume_to_Kube_v1_PersistentVolume(kokiObj)
	case *types.PodWrapper:
		return converters.Convert_Koki_Pod_to_Kube_v1_Pod(kokiObj)
	case *types.ReplicationControllerWrapper:
		return converters.Convert_Koki_ReplicationController_to_Kube_v1_ReplicationController(kokiObj)
	case *types.ReplicaSetWrapper:
		return converters.Convert_Koki_ReplicaSet_to_Kube_ReplicaSet(kokiObj)
	case *types.SecretWrapper:
		return converters.Convert_Koki_Secret_to_Kube_v1_Secret(kokiObj)
	case *types.ServiceWrapper:
		return converters.Convert_Koki_Service_To_Kube_v1_Service(kokiObj)
	case *types.StatefulSetWrapper:
		return converters.Convert_Koki_StatefulSet_to_Kube_StatefulSet(kokiObj)
	case *types.StorageClassWrapper:
		return converters.Convert_Koki_StorageClass_to_Kube_StorageClass(kokiObj)
	case *types.VolumeWrapper:
		return &kokiObj.Volume, nil
	default:
		return nil, util.TypeErrorf(kokiObj, "can't convert from unsupported koki type")
	}
}

func DetectAndConvertFromKubeObj(kubeObj runtime.Object) (interface{}, error) {
	switch kubeObj := kubeObj.(type) {
	case *v1.ConfigMap:
		return converters.Convert_Kube_v1_ConfigMap_to_Koki_ConfigMap(kubeObj)
	case *batchv1beta1.CronJob:
		return converters.Convert_Kube_CronJob_to_Koki_CronJob(kubeObj)
	case *batchv2alpha1.CronJob:
		return converters.Convert_Kube_CronJob_to_Koki_CronJob(kubeObj)
	case *appsv1beta2.DaemonSet:
		return converters.Convert_Kube_DaemonSet_to_Koki_DaemonSet(kubeObj)
	case *exts.DaemonSet:
		return converters.Convert_Kube_DaemonSet_to_Koki_DaemonSet(kubeObj)
	case *appsv1beta1.Deployment:
		return converters.Convert_Kube_Deployment_to_Koki_Deployment(kubeObj)
	case *appsv1beta2.Deployment:
		return converters.Convert_Kube_Deployment_to_Koki_Deployment(kubeObj)
	case *exts.Deployment:
		return converters.Convert_Kube_Deployment_to_Koki_Deployment(kubeObj)
	case *v1.Endpoints:
		return converters.Convert_Kube_v1_Endpoints_to_Koki_Endpoints(kubeObj)
	case *exts.Ingress:
		return converters.Convert_Kube_Ingress_to_Koki_Ingress(kubeObj)
	case *batchv1.Job:
		return converters.Convert_Kube_Job_to_Koki_Job(kubeObj)
	case *v1.PersistentVolume:
		return converters.Convert_Kube_v1_PersistentVolume_to_Koki_PersistentVolume(kubeObj)
	case *v1.PersistentVolumeClaim:
		return converters.Convert_Kube_PVC_to_Koki_PVC(kubeObj)
	case *v1.Pod:
		return converters.Convert_Kube_v1_Pod_to_Koki_Pod(kubeObj)
	case *v1.ReplicationController:
		return converters.Convert_Kube_v1_ReplicationController_to_Koki_ReplicationController(kubeObj)
	case *appsv1beta2.ReplicaSet:
		return converters.Convert_Kube_ReplicaSet_to_Koki_ReplicaSet(kubeObj)
	case *exts.ReplicaSet:
		return converters.Convert_Kube_ReplicaSet_to_Koki_ReplicaSet(kubeObj)
	case *v1.Secret:
		return converters.Convert_Kube_v1_Secret_to_Koki_Secret(kubeObj)
	case *v1.Service:
		return converters.Convert_Kube_v1_Service_to_Koki_Service(kubeObj)
	case *appsv1beta1.StatefulSet:
		return converters.Convert_Kube_StatefulSet_to_Koki_StatefulSet(kubeObj)
	case *appsv1beta2.StatefulSet:
		return converters.Convert_Kube_StatefulSet_to_Koki_StatefulSet(kubeObj)
	case *storagev1.StorageClass:
		return converters.Convert_Kube_StorageClass_to_Koki_StorageClass(kubeObj)
	case *storagev1beta1.StorageClass:
		return converters.Convert_Kube_StorageClass_to_Koki_StorageClass(kubeObj)

	default:
		return nil, util.TypeErrorf(kubeObj, "can't convert from unsupported kube type")
	}
}
