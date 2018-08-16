package converters

import (
	rbac "k8s.io/api/rbac/v1"

	"github.com/koki/short/types"
)

func Convert_Koki_RoleBinding_to_Kube(wrapper *types.RoleBindingWrapper) (*rbac.RoleBinding, error) {
	kube := &rbac.RoleBinding{}
	koki := wrapper.RoleBinding

	kube.Name = koki.Name
	kube.Namespace = koki.Namespace
	if len(koki.Version) == 0 {
		kube.APIVersion = "rbac.authorization.k8s.io/v1"
	} else {
		kube.APIVersion = koki.Version
	}
	kube.Kind = "RoleBinding"
	kube.ClusterName = koki.Cluster
	kube.Labels = koki.Labels
	kube.Annotations = koki.Annotations

	kube.Subjects = revertSubjects(koki.Subjects)
	kube.RoleRef = revertRoleRef(koki.RoleRef)

	return kube, nil
}

