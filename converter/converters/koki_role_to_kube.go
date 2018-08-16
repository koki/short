package converters

import (
	rbac "k8s.io/api/rbac/v1"

	"github.com/koki/short/types"
)

func Convert_Koki_Role_to_Kube(wrapper *types.RoleWrapper) (*rbac.Role, error) {
	kube := &rbac.Role{}
	koki := &wrapper.Role

	kube.Name = koki.Name
	kube.Namespace = koki.Namespace
	if len(koki.Version) == 0 {
		kube.APIVersion = "rbac.authorization.k8s.io/v1"
	} else {
		kube.APIVersion = koki.Version
	}
	kube.Kind = "Role"
	kube.ClusterName = koki.Cluster
	kube.Labels = koki.Labels
	kube.Annotations = koki.Annotations

	kube.Rules = revertPolicyRules(koki.Rules)

	return kube, nil
}

