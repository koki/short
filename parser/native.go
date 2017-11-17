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
		return nil, err
	}

	for k := range objMap {
		switch k {
		case "deployment":
			deployment := &types.DeploymentWrapper{}
			err := json.Unmarshal(bytes, deployment)
			return deployment, err
		case "persistent_volume":
			pv := &types.PersistentVolumeWrapper{}
			err := json.Unmarshal(bytes, pv)
			return pv, err
		case "pod":
			pod := &types.PodWrapper{}
			err := json.Unmarshal(bytes, pod)
			return pod, err
		case "replica_set":
			replicaSet := &types.ReplicaSetWrapper{}
			err := json.Unmarshal(bytes, replicaSet)
			return replicaSet, err
		case "replication_controller":
			replicationController := &types.ReplicationControllerWrapper{}
			err := json.Unmarshal(bytes, replicationController)
			return replicationController, err
		case "service":
			service := &types.ServiceWrapper{}
			err := json.Unmarshal(bytes, service)
			return service, err
		case "volume":
			volume := &types.VolumeWrapper{}
			err := json.Unmarshal(bytes, volume)
			return volume, err
		}

		return nil, util.TypeErrorf(objMap, "Unexpected key (%s)", k)
	}

	return nil, nil
}

func UnparseKokiNativeObject(kokiObj interface{}) (map[string]interface{}, error) {
	// Marshal the koki object back into yaml.
	bytes, err := yaml.Marshal(kokiObj)
	if err != nil {
		return nil, err
	}

	obj := map[string]interface{}{}
	err = yaml.Unmarshal(bytes, &obj)
	if err != nil {
		return nil, err
	}

	return obj, err
}
