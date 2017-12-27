package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_Binding_to_Kube_Binding(kokiWrapper *types.BindingWrapper) (*v1.Binding, error) {
	var err error
	kubeBinding := &v1.Binding{}
	kokiBinding := kokiWrapper.Binding

	kubeBinding.Name = kokiBinding.Name
	kubeBinding.Namespace = kokiBinding.Namespace
	if len(kokiBinding.Version) == 0 {
		kubeBinding.APIVersion = "v1"
	} else {
		kubeBinding.APIVersion = kokiBinding.Version
	}
	kubeBinding.Kind = "Binding"
	kubeBinding.ClusterName = kokiBinding.Cluster
	kubeBinding.Labels = kokiBinding.Labels
	kubeBinding.Annotations = kokiBinding.Annotations

	target, err := revertTarget(&kokiBinding.Target)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "binding target")
	}
	kubeBinding.Target = *target

	return kubeBinding, nil
}
