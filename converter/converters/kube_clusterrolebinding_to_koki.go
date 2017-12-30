package converters

import (
	rbac "k8s.io/api/rbac/v1"

	"github.com/koki/short/types"
)

func Convert_Kube_ClusterRoleBinding_to_Koki(kube *rbac.ClusterRoleBinding) (*types.ClusterRoleBindingWrapper, error) {
	koki := &types.ClusterRoleBinding{}

	koki.Name = kube.Name
	koki.Namespace = kube.Namespace
	koki.Version = kube.APIVersion
	koki.Cluster = kube.ClusterName
	koki.Labels = kube.Labels
	koki.Annotations = kube.Annotations

	koki.Subjects = convertSubjects(kube.Subjects)
	koki.RoleRef = convertRoleRef(kube.RoleRef)

	return &types.ClusterRoleBindingWrapper{
		ClusterRoleBinding: *koki,
	}, nil
}

func convertSubjects(kubeSubjects []rbac.Subject) []types.Subject {
	kokiSubjects := make([]types.Subject, len(kubeSubjects))
	for i, kubeSubject := range kubeSubjects {
		kokiSubjects[i] = types.Subject{
			Kind:      kubeSubject.Kind,
			APIGroup:  kubeSubject.APIGroup,
			Name:      types.Name(kubeSubject.Name),
			Namespace: kubeSubject.Namespace,
		}
	}

	return kokiSubjects
}

func convertRoleRef(kubeRef rbac.RoleRef) types.RoleRef {
	return types.RoleRef{
		APIGroup: kubeRef.APIGroup,
		Kind:     kubeRef.Kind,
		Name:     types.Name(kubeRef.Name),
	}
}
