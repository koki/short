package param

import (
	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func ApplyReplicationControllerParams(params map[string]interface{}, wrapper *types.ReplicationControllerWrapper) error {
	rc := &wrapper.ReplicationController

	if name, ok := params["name"]; ok {
		if name, ok := name.(string); ok {
			rc.Name = name
			if rc.Labels != nil {
				rc.Labels["name"] = name
			} else {
				rc.Labels = map[string]string{
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
			rc.SetTemplate(&kokiPod.Pod)
			// We'll either use the Pod's Labels or generate a Selector and Labels on conversion to kube obj.
			rc.Selector = nil
		} else {
			return util.TypeErrorf(kokiObj, "expected a pod")
		}
	}
	return nil
}
