package parser

import (
	"encoding/json"
	"fmt"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func ParseKokiNativeObject(obj interface{}) (interface{}, error) {
	if _, ok := obj.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("Error casting input object into map[string]interface{}")
	}

	objMap := obj.(map[string]interface{})

	if len(objMap) != 1 {
		return nil, util.TypeValueErrorf(objMap, "Invalid koki syntax")
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
		case "persistentVolume":
			pv := &types.PersistentVolumeWrapper{}
			err := json.Unmarshal(bytes, pv)
			return pv, err
		case "pod":
			pod := &types.PodWrapper{}
			err := json.Unmarshal(bytes, pod)
			return pod, err
		case "replicaSet":
			replicaSet := &types.ReplicaSetWrapper{}
			err := json.Unmarshal(bytes, replicaSet)
			return replicaSet, err
		case "replicationController":
			replicationController := &types.ReplicationControllerWrapper{}
			err := json.Unmarshal(bytes, replicationController)
			return replicationController, err
		case "service":
			service := &types.ServiceWrapper{}
			err := json.Unmarshal(bytes, service)
			return service, err
		}

		return nil, util.TypeValueErrorf(objMap, "Unexpected value %s", k)
	}

	return nil, nil
}
