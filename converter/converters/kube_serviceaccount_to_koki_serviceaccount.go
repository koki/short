package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
)

func Convert_Kube_ServiceAccount_to_Koki_ServiceAccount(kubeServiceAccount *v1.ServiceAccount) (*types.ServiceAccountWrapper, error) {
	kokiWrapper := &types.ServiceAccountWrapper{}
	kokiServiceAccount := &kokiWrapper.ServiceAccount

	kokiServiceAccount.Name = kubeServiceAccount.Name
	kokiServiceAccount.Namespace = kubeServiceAccount.Namespace
	kokiServiceAccount.Version = kubeServiceAccount.APIVersion
	kokiServiceAccount.Cluster = kubeServiceAccount.ClusterName
	kokiServiceAccount.Labels = kubeServiceAccount.Labels
	kokiServiceAccount.Annotations = kubeServiceAccount.Annotations

	kokiServiceAccount.AutomountServiceAccountToken = kubeServiceAccount.AutomountServiceAccountToken

	for i := range kubeServiceAccount.ImagePullSecrets {
		kokiServiceAccount.ImagePullSecrets = append(kokiServiceAccount.ImagePullSecrets, kubeServiceAccount.ImagePullSecrets[i].Name)
	}

	for i := range kubeServiceAccount.Secrets {
		kubeSecret := kubeServiceAccount.Secrets[i]
		kokiServiceAccount.Secrets = append(kokiServiceAccount.Secrets, types.ObjectReference{
			Kind:            kubeSecret.Kind,
			Namespace:       kubeSecret.Namespace,
			Name:            kubeSecret.Name,
			UID:             string(kubeSecret.UID),
			Version:         kubeSecret.APIVersion,
			ResourceVersion: kubeSecret.ResourceVersion,
			FieldPath:       kubeSecret.FieldPath,
		})
	}

	return kokiWrapper, nil
}
