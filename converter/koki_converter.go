package converter

import (
	"github.com/koki/short/converter/converters"
	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"

	admissionregv1alpha1 "k8s.io/api/admissionregistration/v1alpha1"
	apps "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"
	autoscaling "k8s.io/api/autoscaling/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	batchv2alpha1 "k8s.io/api/batch/v2alpha1"
	"k8s.io/api/core/v1"
	exts "k8s.io/api/extensions/v1beta1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	schedulingv1alpha1 "k8s.io/api/scheduling/v1alpha1"
	settingsv1alpha1 "k8s.io/api/settings/v1alpha1"
	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiregistrationv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"

	"k8s.io/apimachinery/pkg/runtime"
)

func DetectAndConvertFromKokiObj(kokiObj interface{}) (interface{}, error) {
	switch kokiObj := kokiObj.(type) {
	case *types.APIServiceWrapper:
		return converters.Convert_Koki_APIService_to_Kube_APIService(kokiObj)
	case *types.BindingWrapper:
		return converters.Convert_Koki_Binding_to_Kube_Binding(kokiObj)
	case *types.ConfigMapWrapper:
		return converters.Convert_Koki_ConfigMap_to_Kube_v1_ConfigMap(kokiObj)
	case *types.ControllerRevisionWrapper:
		return converters.Convert_Koki_ControllerRevision_to_Kube(kokiObj)
	case *types.CronJobWrapper:
		return converters.Convert_Koki_CronJob_to_Kube_CronJob(kokiObj)
	case *types.CRDWrapper:
		return converters.Convert_Koki_CRD_to_Kube(kokiObj)
	case *types.DaemonSetWrapper:
		return converters.Convert_Koki_DaemonSet_to_Kube_DaemonSet(kokiObj)
	case *types.DeploymentWrapper:
		return converters.Convert_Koki_Deployment_to_Kube_Deployment(kokiObj)
	case *types.EndpointsWrapper:
		return converters.Convert_Koki_Endpoints_to_Kube_v1_Endpoints(kokiObj)
	case *types.EventWrapper:
		return converters.Convert_Koki_Event_to_Kube(kokiObj)
	case *types.HorizontalPodAutoscalerWrapper:
		return converters.Convert_Koki_HPA_to_Kube(kokiObj)
	case *types.IngressWrapper:
		return converters.Convert_Koki_Ingress_to_Kube_Ingress(kokiObj)
	case *types.InitializerConfigWrapper:
		return converters.Convert_Koki_InitializerConfig_to_Kube_InitializerConfig(kokiObj)
	case *types.JobWrapper:
		return converters.Convert_Koki_Job_to_Kube_Job(kokiObj)
	case *types.LimitRangeWrapper:
		return converters.Convert_Koki_LimitRange_to_Kube(kokiObj)
	case *types.NamespaceWrapper:
		return converters.Convert_Koki_Namespace_to_Kube_Namespace(kokiObj)
	case *types.PersistentVolumeClaimWrapper:
		return converters.Convert_Koki_PVC_to_Kube_PVC(kokiObj)
	case *types.PersistentVolumeWrapper:
		return converters.Convert_Koki_PersistentVolume_to_Kube_v1_PersistentVolume(kokiObj)
	case *types.PodDisruptionBudgetWrapper:
		return converters.Convert_Koki_PodDisruptionBudget_to_Kube_PodDisruptionBudget(kokiObj)
	case *types.PodPresetWrapper:
		return converters.Convert_Koki_PodPreset_to_Kube_PodPreset(kokiObj)
	case *types.PodSecurityPolicyWrapper:
		return converters.Convert_Koki_PodSecurityPolicy_to_Kube_PodSecurityPolicy(kokiObj)
	case *types.PodWrapper:
		return converters.Convert_Koki_Pod_to_Kube_v1_Pod(kokiObj)
	case *types.PriorityClassWrapper:
		return converters.Convert_Koki_PriorityClass_to_Kube_PriorityClass(kokiObj)
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
		return nil, serrors.TypeErrorf(kokiObj, "can't convert from unsupported koki type")
	}
}

func DetectAndConvertFromKubeObj(kubeObj runtime.Object) (interface{}, error) {
	switch kubeObj := kubeObj.(type) {
	case *apiregistrationv1beta1.APIService:
		return converters.Convert_Kube_APIService_to_Koki_APIService(kubeObj)
	case *v1.Binding:
		return converters.Convert_Kube_Binding_to_Koki_Binding(kubeObj)
	case *v1.ConfigMap:
		return converters.Convert_Kube_v1_ConfigMap_to_Koki_ConfigMap(kubeObj)
	case *apps.ControllerRevision, *appsv1beta1.ControllerRevision, *appsv1beta2.ControllerRevision:
		return converters.Convert_Kube_ControllerRevision_to_Koki(kubeObj)
	case *batchv1beta1.CronJob, *batchv2alpha1.CronJob:
		return converters.Convert_Kube_CronJob_to_Koki_CronJob(kubeObj)
	case *apiext.CustomResourceDefinition:
		return converters.Convert_Kube_CRD_to_Koki(kubeObj)
	case *appsv1beta2.DaemonSet, *exts.DaemonSet:
		return converters.Convert_Kube_DaemonSet_to_Koki_DaemonSet(kubeObj)
	case *appsv1beta1.Deployment, *appsv1beta2.Deployment, *exts.Deployment:
		return converters.Convert_Kube_Deployment_to_Koki_Deployment(kubeObj)
	case *v1.Endpoints:
		return converters.Convert_Kube_v1_Endpoints_to_Koki_Endpoints(kubeObj)
	case *v1.Event:
		return converters.Convert_Kube_Event_to_Koki(kubeObj)
	case *autoscaling.HorizontalPodAutoscaler:
		return converters.Convert_Kube_HPA_to_Koki(kubeObj)
	case *exts.Ingress:
		return converters.Convert_Kube_Ingress_to_Koki_Ingress(kubeObj)
	case *admissionregv1alpha1.InitializerConfiguration:
		return converters.Convert_Kube_InitializerConfig_to_Koki_InitializerConfig(kubeObj)
	case *batchv1.Job:
		return converters.Convert_Kube_Job_to_Koki_Job(kubeObj)
	case *v1.LimitRange:
		return converters.Convert_Kube_LimitRange_to_Koki(kubeObj)
	case *v1.Namespace:
		return converters.Convert_Kube_Namespace_to_Koki_Namespace(kubeObj)
	case *v1.PersistentVolume:
		return converters.Convert_Kube_v1_PersistentVolume_to_Koki_PersistentVolume(kubeObj)
	case *v1.PersistentVolumeClaim:
		return converters.Convert_Kube_PVC_to_Koki_PVC(kubeObj)
	case *schedulingv1alpha1.PriorityClass:
		return converters.Convert_Kube_PriorityClass_to_Koki_PriorityClass(kubeObj)
	case *v1.Pod:
		return converters.Convert_Kube_v1_Pod_to_Koki_Pod(kubeObj)
	case *policyv1beta1.PodDisruptionBudget:
		return converters.Convert_Kube_PodDisruptionBudget_to_Koki_PodDisruptionBudget(kubeObj)
	case *settingsv1alpha1.PodPreset:
		return converters.Convert_Kube_PodPreset_to_Koki_PodPreset(kubeObj)
	case *exts.PodSecurityPolicy:
		return converters.Convert_Kube_PodSecurityPolicy_to_Koki_PodSecurityPolicy(kubeObj)
	case *v1.ReplicationController:
		return converters.Convert_Kube_v1_ReplicationController_to_Koki_ReplicationController(kubeObj)
	case *appsv1beta2.ReplicaSet, *exts.ReplicaSet:
		return converters.Convert_Kube_ReplicaSet_to_Koki_ReplicaSet(kubeObj)
	case *v1.Secret:
		return converters.Convert_Kube_v1_Secret_to_Koki_Secret(kubeObj)
	case *v1.Service:
		return converters.Convert_Kube_v1_Service_to_Koki_Service(kubeObj)
	case *appsv1beta1.StatefulSet, *appsv1beta2.StatefulSet:
		return converters.Convert_Kube_StatefulSet_to_Koki_StatefulSet(kubeObj)
	case *storagev1.StorageClass, *storagev1beta1.StorageClass:
		return converters.Convert_Kube_StorageClass_to_Koki_StorageClass(kubeObj)

	default:
		return nil, serrors.TypeErrorf(kubeObj, "can't convert from unsupported kube type")
	}
}
