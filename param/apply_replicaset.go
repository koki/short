package param

import (
	//"github.com/koki/short/parser"
	"github.com/koki/short/types"
	//"github.com/koki/short/util"
)

func ApplyReplicaSetParams(params map[string]interface{}, wrapper *types.ReplicaSetWrapper) error {
	/*
		replicaSet := &wrapper.ReplicaSet
		if name, ok := params["name"]; ok {
			if name, ok := name.(string); ok {
				replicaSet.Name = name
				if replicaSet.Labels != nil {
					replicaSet.Labels["name"] = name
				} else {
					replicaSet.Labels = map[string]string{
						"name": name,
					}
				}
			} else {
				return util.PrettyTypeError(params, "expected string for 'name'")
			}
		}

		if pod, ok := params["pod"]; ok {
			kokiObj, err := parser.ParseKokiNativeObject(pod)
			if err != nil {
				return err
			}
			if kokiPod, ok := kokiObj.(*types.PodWrapper); ok {
				replicaSet.Template = &kokiPod.Pod
				// Empty selector just uses template's labels.
				replicaSet.PodSelector = ""
			} else {
				return util.PrettyTypeError(kokiObj, "expected a pod")
			}
		}
	*/

	return nil
}
