package parser

import (
	"encoding/json"

	"github.com/koki/short/types"
	"github.com/koki/short/util"

	"github.com/ghodss/yaml"
)

func ParseKokiNativeObject(obj interface{}) (interface{}, error) {
	if _, ok := obj.(map[string]interface{}); !ok {
		return nil, util.TypeErrorf(obj, "can only parse map[string]interface{} as koki obj")
	}

	objMap := obj.(map[string]interface{})

	if len(objMap) != 1 {
		return nil, util.InvalidValueErrorf(objMap, "Invalid koki syntax")
	}

	bytes, err := json.Marshal(objMap)
	if err != nil {
		return nil, util.InvalidValueErrorf(objMap, "error converting to JSON before re-parsing as as koki obj: %s", err.Error())
	}

	for k := range objMap {
		switch k {
		case "deployment":
			deployment := &types.DeploymentWrapper{}
			err := json.Unmarshal(bytes, deployment)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, deployment, err.Error())
			}
			return deployment, nil
		case "persistent_volume":
			pv := &types.PersistentVolumeWrapper{}
			err := json.Unmarshal(bytes, pv)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, pv, err.Error())
			}
			return pv, nil
		case "pod":
			pod := &types.PodWrapper{}
			err := json.Unmarshal(bytes, pod)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, pod, err.Error())
			}
			return pod, nil
		case "replica_set":
			replicaSet := &types.ReplicaSetWrapper{}
			err := json.Unmarshal(bytes, replicaSet)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, replicaSet, err.Error())
			}
			return replicaSet, nil
		case "replication_controller":
			replicationController := &types.ReplicationControllerWrapper{}
			err := json.Unmarshal(bytes, replicationController)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, replicationController, err.Error())
			}
			return replicationController, nil
		case "service":
			service := &types.ServiceWrapper{}
			err := json.Unmarshal(bytes, service)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, service, err.Error())
			}
			return service, nil
		case "volume":
			volume := &types.VolumeWrapper{}
			err := json.Unmarshal(bytes, volume)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, volume, err.Error())
			}
			return volume, nil
		case "job":
			job := &types.JobWrapper{}
			err := json.Unmarshal(bytes, job)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, job, err.Error())
			}
			return job, nil
		case "daemon_set":
			daemonSet := &types.DaemonSetWrapper{}
			err := json.Unmarshal(bytes, daemonSet)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, daemonSet, err.Error())
			}
			return daemonSet, nil
		case "cron_job":
			cronJob := &types.CronJobWrapper{}
			err := json.Unmarshal(bytes, cronJob)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, cronJob, err.Error())
			}
			return cronJob, nil
		case "pvc":
			pvc := &types.PersistentVolumeClaimWrapper{}
			err := json.Unmarshal(bytes, pvc)
			if err != nil {
				return nil, util.InvalidValueForTypeErrorf(objMap, pvc, err.Error())
			}
			return pvc, nil
		}
		return nil, util.TypeErrorf(objMap, "Unexpected key (%s)", k)
	}

	return nil, nil
}

func UnparseKokiNativeObject(kokiObj interface{}) (map[string]interface{}, error) {
	// Marshal the koki object back into yaml.
	bytes, err := yaml.Marshal(kokiObj)
	if err != nil {
		return nil, util.InvalidInstanceErrorf(kokiObj, "couldn't convert to yaml: %s", err.Error())
	}

	obj := map[string]interface{}{}
	err = yaml.Unmarshal(bytes, &obj)
	if err != nil {
		return nil, err
	}

	return obj, util.InvalidInstanceErrorf(kokiObj, "couldn't convert to dictionary: %s", err.Error())
}
