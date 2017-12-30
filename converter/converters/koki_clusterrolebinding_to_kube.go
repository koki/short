package converters

import (
	rbac "k8s.io/api/rbac/v1"

	"github.com/koki/short/types"
)

func Convert_Koki_ClusterRoleBinding_to_Kube(wrapper *types.ClusterRoleBindingWrapper) (*rbac.ClusterRoleBinding, error) {
	kube := &rbac.ClusterRoleBinding{}
	koki := wrapper.ClusterRoleBinding

	kube.Name = koki.Name
	kube.Namespace = koki.Namespace
	if len(koki.Version) == 0 {
		kube.APIVersion = "rbac.authorization.k8s.io/v1"
	} else {
		kube.APIVersion = koki.Version
	}
	kube.Kind = "ClusterRoleBinding"
	kube.ClusterName = koki.Cluster
	kube.Labels = koki.Labels
	kube.Annotations = koki.Annotations

	kube.Subjects = revertSubjects(koki.Subjects)
	kube.RoleRef = revertRoleRef(koki.RoleRef)

	return kube, nil
}

func revertSubjects(kokiSubjects []types.Subject) []rbac.Subject {
	kubeSubjects := make([]rbac.Subject, len(kokiSubjects))
	for i, kokiSubject := range kokiSubjects {
		kubeSubjects[i] = rbac.Subject{
			Kind:      kokiSubject.Kind,
			APIGroup:  kokiSubject.APIGroup,
			Name:      string(kokiSubject.Name),
			Namespace: kokiSubject.Namespace,
		}
	}

	return kubeSubjects
}

func revertRoleRef(kokiRef types.RoleRef) rbac.RoleRef {
	return rbac.RoleRef{
		APIGroup: kokiRef.APIGroup,
		Kind:     kokiRef.Kind,
		Name:     string(kokiRef.Name),
	}
}
