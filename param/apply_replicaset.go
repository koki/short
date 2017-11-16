package param

import (
	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func ApplyReplicaSetParams(params map[string]interface{}, wrapper *types.ReplicaSetWrapper) error {
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
			return util.InvalidValueErrorf(params, "expected string for param 'name'")
		}
	}

	if pod, ok := params["pod"]; ok {
		kokiObj, err := parser.ParseKokiNativeObject(pod)
		if err != nil {
			return err
		}
		if kokiPod, ok := kokiObj.(*types.PodWrapper); ok {
			// If the Pod has labels, pull them up into the Selector.
			// Otherwise, automatically set up Selector and Labels on conversion to kube obj.
			labels := kokiPod.Pod.Labels
			kokiPod.Pod.Labels = nil
			replicaSet.SetTemplate(&kokiPod.Pod)
			if len(kokiPod.Pod.Labels) > 0 {
				replicaSet.Selector = &types.RSSelector{
					Labels: labels,
				}
			} else {
				replicaSet.Selector = nil
			}
		} else {
			return util.TypeErrorf(kokiObj, "expected a pod")
		}
	}

	return nil
}
