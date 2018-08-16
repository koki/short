package converters

import (
	rbac "k8s.io/api/rbac/v1"

	"github.com/koki/short/types"
)

func Convert_Kube_Role_to_Koki(kube *rbac.Role) (*types.RoleWrapper, error) {
	koki := &types.Role{}

	koki.Name = kube.Name
	koki.Namespace = kube.Namespace
	koki.Version = kube.APIVersion
	koki.Cluster = kube.ClusterName
	koki.Labels = kube.Labels
	koki.Annotations = kube.Annotations

	koki.Rules = convertPolicyRules(kube.Rules)

	return &types.RoleWrapper{
		Role: *koki,
	}, nil
}

