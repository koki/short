package parser

import (
	"github.com/koki/json"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func ParseKokiNativeObject(obj interface{}) (interface{}, error) {
	if _, ok := obj.(map[string]interface{}); !ok {
		return nil, serrors.TypeErrorf(obj, "can only parse map[string]interface{} as koki obj")
	}

	objMap := obj.(map[string]interface{})

	if len(objMap) != 1 {
		return nil, serrors.InvalidValueErrorf(objMap, "Invalid koki syntax")
	}

	bytes, err := json.Marshal(objMap)
	if err != nil {
		return nil, serrors.InvalidValueContextErrorf(err, objMap, "error converting to JSON before re-parsing as as koki obj")
	}

	for k := range objMap {
		switch k {
		case "api_service":
			apiService := &types.APIService{}
			err := json.Unmarshal(bytes, apiService)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, apiService)
			}
			return apiService, nil
		case "binding":
			binding := &types.BindingWrapper{}
			err := json.Unmarshal(bytes, binding)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, binding)
			}
			return binding, nil
		case "csr":
			csr := &types.CertificateSigningRequestWrapper{}
			err := json.Unmarshal(bytes, csr)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, csr)
			}
			return csr, nil
		case "cluster_role":
			result := &types.ClusterRoleWrapper{}
			err := json.Unmarshal(bytes, result)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, result)
			}
			return result, nil
		case "cluster_role_binding":
			result := &types.ClusterRoleBindingWrapper{}
			err := json.Unmarshal(bytes, result)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, result)
			}
			return result, nil
		case "config_map":
			configMap := &types.ConfigMapWrapper{}
			err := json.Unmarshal(bytes, configMap)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, configMap)
			}
			return configMap, nil
		case "controller_revision":
			rev := &types.ControllerRevisionWrapper{}
			err := json.Unmarshal(bytes, rev)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, rev)
			}
			return rev, nil
		case "crd":
			result := &types.CRDWrapper{}
			err := json.Unmarshal(bytes, result)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, result)
			}
			return result, nil
		case "cron_job":
			cronJob := &types.CronJobWrapper{}
			err := json.Unmarshal(bytes, cronJob)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, cronJob)
			}
			return cronJob, nil
		case "daemon_set":
			daemonSet := &types.DaemonSetWrapper{}
			err := json.Unmarshal(bytes, daemonSet)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, daemonSet)
			}
			return daemonSet, nil
		case "deployment":
			deployment := &types.DeploymentWrapper{}
			err := json.Unmarshal(bytes, deployment)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, deployment)
			}
			return deployment, nil
		case "endpoints":
			endpoints := &types.EndpointsWrapper{}
			err := json.Unmarshal(bytes, endpoints)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, endpoints)
			}
			return endpoints, nil
		case "event":
			result := &types.EventWrapper{}
			err := json.Unmarshal(bytes, result)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, result)
			}
			return result, nil
		case "hpa":
			result := &types.HorizontalPodAutoscalerWrapper{}
			err := json.Unmarshal(bytes, result)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, result)
			}
			return result, nil
		case "ingress":
			ingress := &types.IngressWrapper{}
			err := json.Unmarshal(bytes, ingress)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, ingress)
			}
			return ingress, nil
		case "initializer_config":
			initConfig := &types.InitializerConfigWrapper{}
			err := json.Unmarshal(bytes, initConfig)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, initConfig)
			}
			return initConfig, nil
		case "job":
			job := &types.JobWrapper{}
			err := json.Unmarshal(bytes, job)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, job)
			}
			return job, nil
		case "limit_range":
			result := &types.LimitRangeWrapper{}
			err := json.Unmarshal(bytes, result)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, result)
			}
			return result, nil
		case "namespace":
			namespace := &types.NamespaceWrapper{}
			err := json.Unmarshal(bytes, namespace)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, namespace)
			}
			return namespace, nil
		case "pdb":
			pdb := &types.PodDisruptionBudgetWrapper{}
			err := json.Unmarshal(bytes, pdb)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, pdb)
			}
			return pdb, nil
		case "persistent_volume":
			pv := &types.PersistentVolumeWrapper{}
			err := json.Unmarshal(bytes, pv)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, pv)
			}
			return pv, nil
		case "pod":
			pod := &types.PodWrapper{}
			err := json.Unmarshal(bytes, pod)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, pod)
			}
			return pod, nil
		case "pod_preset":
			podPreset := &types.PodPresetWrapper{}
			err = json.Unmarshal(bytes, podPreset)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, podPreset)
			}
			return podPreset, nil
		case "pod_security_policy":
			podSecurityPolicy := &types.PodSecurityPolicyWrapper{}
			err = json.Unmarshal(bytes, podSecurityPolicy)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, podSecurityPolicy)
			}
			return podSecurityPolicy, nil
		case "pod_template":
			podTemplate := &types.PodTemplateWrapper{}
			err = json.Unmarshal(bytes, podTemplate)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, podTemplate)
			}
			return podTemplate, nil
		case "priority_class":
			priorityClass := &types.PriorityClassWrapper{}
			err := json.Unmarshal(bytes, priorityClass)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, priorityClass)
			}
			return priorityClass, nil
		case "pvc":
			pvc := &types.PersistentVolumeClaimWrapper{}
			err := json.Unmarshal(bytes, pvc)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, pvc)
			}
			return pvc, nil
		case "replica_set":
			replicaSet := &types.ReplicaSetWrapper{}
			err := json.Unmarshal(bytes, replicaSet)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, replicaSet)
			}
			return replicaSet, nil
		case "replication_controller":
			replicationController := &types.ReplicationControllerWrapper{}
			err := json.Unmarshal(bytes, replicationController)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, replicationController)
			}
			return replicationController, nil
		case "secret":
			secret := &types.SecretWrapper{}
			err := json.Unmarshal(bytes, secret)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, secret)
			}
			return secret, nil
		case "service":
			service := &types.ServiceWrapper{}
			err := json.Unmarshal(bytes, service)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, service)
			}
			return service, nil
		case "service_account":
			serviceAccount := &types.ServiceAccountWrapper{}
			err := json.Unmarshal(bytes, serviceAccount)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, serviceAccount)
			}
			return serviceAccount, nil
		case "stateful_set":
			statefulSet := &types.StatefulSetWrapper{}
			err := json.Unmarshal(bytes, statefulSet)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, statefulSet)
			}
			return statefulSet, nil
		case "storage_class":
			storageClass := &types.StorageClassWrapper{}
			err := json.Unmarshal(bytes, storageClass)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, storageClass)
			}
			return storageClass, nil
		case "volume":
			volume := &types.VolumeWrapper{}
			err := json.Unmarshal(bytes, volume)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, volume)
			}
			return volume, nil
		}
		return nil, serrors.TypeErrorf(objMap, "Unexpected key (%s)", k)
	}
	return nil, nil
}

func UnparseKokiNativeObject(kokiObj interface{}) (map[string]interface{}, error) {
	// Marshal the koki object back into yaml.
	bytes, err := yaml.Marshal(kokiObj)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, kokiObj, "converting to yaml")
	}

	obj := map[string]interface{}{}
	err = yaml.Unmarshal(bytes, &obj)
	if err != nil {
		return nil, err
	}

	return obj, serrors.InvalidInstanceContextErrorf(err, kokiObj, "converting to dictionary")
}
