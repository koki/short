package param

import (
	"github.com/koki/short/imports"
	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func ApplyParams(params map[string]interface{}, module *imports.Module) error {
	if len(params) == 0 {
		return nil
	}

	var err error
	switch obj := module.TypedResult.(type) {
	case *types.PodWrapper:
		err = ApplyPodParams(params, obj)
	case *types.ServiceWrapper:
		err = ApplyServiceParams(params, obj)
	case *types.ReplicaSetWrapper:
		err = ApplyReplicaSetParams(params, obj)
	case *types.ReplicationControllerWrapper:
		err = ApplyReplicationControllerParams(params, obj)
	case *types.DeploymentWrapper:
		err = ApplyDeploymentParams(params, obj)
	default:
		err = util.TypeErrorf(obj, "unsupported type for parameterization")
	}

	if err != nil {
		return nil
	}

	// Update Raw to reflect new TypedResult.
	module.Raw, err = parser.UnparseKokiNativeObject(module.TypedResult)

	return err
}
