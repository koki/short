package param

import (
	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func ApplyDeploymentParams(params map[string]interface{}, wrapper *types.DeploymentWrapper) error {
	deployment := &wrapper.Deployment
	if name, ok := params["name"]; ok {
		if name, ok := name.(string); ok {
			deployment.Name = name
			if deployment.Labels != nil {
				deployment.Labels["name"] = name
			} else {
				deployment.Labels = map[string]string{
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
			deployment.Template = kokiPod.Pod
		} else {
			return util.PrettyTypeError(kokiObj, "expected a pod")
		}
	}

	return nil
}
