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
		case "ingress":
			ingress := &types.IngressWrapper{}
			err := json.Unmarshal(bytes, ingress)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, ingress)
			}
			return ingress, nil
		case "job":
			job := &types.JobWrapper{}
			err := json.Unmarshal(bytes, job)
			if err != nil {
				return nil, serrors.InvalidValueForTypeContextError(err, objMap, job)
			}
			return job, nil
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
