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
		case "pod":
			pod := &types.PodWrapper{}
			err := json.Unmarshal(bytes, pod)
			return pod, err
		case "service":
			service := &types.ServiceWrapper{}
			err := json.Unmarshal(bytes, service)
			return service, err
		}

		if k == "replicaSet" {
			replicaSet := &types.ReplicaSetWrapper{}
			err := json.Unmarshal(bytes, replicaSet)
			return replicaSet, err
		}

		if k == "replicationController" {
			replicationController := &types.ReplicationControllerWrapper{}
			err := json.Unmarshal(bytes, replicationController)
			return replicationController, err
		}

		return nil, util.TypeValueErrorf(objMap, "Unexpected value %s", k)
	}

	return nil, nil
}
