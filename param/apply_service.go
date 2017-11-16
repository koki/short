package param

import (
	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func ApplyServiceParams(params map[string]interface{}, wrapper *types.ServiceWrapper) error {
	service := &wrapper.Service

	if name, ok := params["name"]; ok {
		if name, ok := name.(string); ok {
			service.Name = name
			if service.Labels != nil {
				service.Labels["name"] = name
			} else {
				service.Labels = map[string]string{
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
			service.Selector = kokiPod.Pod.Labels
		} else {
			return util.TypeErrorf(kokiObj, "expected a pod")
		}
	}

	return nil
}
