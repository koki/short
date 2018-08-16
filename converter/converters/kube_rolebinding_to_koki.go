package converters

import (
	rbac "k8s.io/api/rbac/v1"

	"github.com/koki/short/types"
)

func Convert_Kube_RoleBinding_to_Koki(kube *rbac.RoleBinding) (*types.RoleBindingWrapper, error) {
	koki := &types.RoleBinding{}

	koki.Name = kube.Name
	koki.Namespace = kube.Namespace
	koki.Version = kube.APIVersion
	koki.Cluster = kube.ClusterName
	koki.Labels = kube.Labels
	koki.Annotations = kube.Annotations

	koki.Subjects = convertSubjects(kube.Subjects)
	koki.RoleRef = convertRoleRef(kube.RoleRef)

	return &types.RoleBindingWrapper{
		RoleBinding: *koki,
	}, nil
}

