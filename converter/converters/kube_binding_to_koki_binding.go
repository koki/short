package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_Binding_to_Koki_Binding(kubeBinding *v1.Binding) (*types.BindingWrapper, error) {
	kokiWrapper := &types.BindingWrapper{}
	kokiBinding := &kokiWrapper.Binding

	kokiBinding.Name = kubeBinding.Name
	kokiBinding.Namespace = kubeBinding.Namespace
	kokiBinding.Version = kubeBinding.APIVersion
	kokiBinding.Cluster = kubeBinding.ClusterName
	kokiBinding.Labels = kubeBinding.Labels
	kokiBinding.Annotations = kubeBinding.Annotations

	target, err := convertTarget(&kubeBinding.Target)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "Binding Target")
	}
	kokiBinding.Target = *target

	return kokiWrapper, nil
}
