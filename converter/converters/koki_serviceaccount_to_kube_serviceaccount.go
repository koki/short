package converters

import (
	"k8s.io/api/core/v1"
	apimachinery "k8s.io/apimachinery/pkg/types"

	"github.com/koki/short/types"
)

func Convert_Koki_ServiceAccount_to_Kube_ServiceAccount(kokiWrapper *types.ServiceAccountWrapper) (*v1.ServiceAccount, error) {
	kubeServiceAccount := &v1.ServiceAccount{}
	kokiServiceAccount := kokiWrapper.ServiceAccount

	kubeServiceAccount.Name = kokiServiceAccount.Name
	kubeServiceAccount.Namespace = kokiServiceAccount.Namespace
	if len(kokiServiceAccount.Version) == 0 {
		kubeServiceAccount.APIVersion = "v1"
	} else {
		kubeServiceAccount.APIVersion = kokiServiceAccount.Version
	}
	kubeServiceAccount.Kind = "ServiceAccount"
	kubeServiceAccount.ClusterName = kokiServiceAccount.Cluster
	kubeServiceAccount.Labels = kokiServiceAccount.Labels
	kubeServiceAccount.Annotations = kokiServiceAccount.Annotations

	kubeServiceAccount.AutomountServiceAccountToken = kokiServiceAccount.AutomountServiceAccountToken

	for i := range kokiServiceAccount.ImagePullSecrets {
		kubeServiceAccount.ImagePullSecrets = append(kubeServiceAccount.ImagePullSecrets, v1.LocalObjectReference{
			Name: kokiServiceAccount.ImagePullSecrets[i],
		})
	}

	for i := range kokiServiceAccount.Secrets {
		kokiSecret := kokiServiceAccount.Secrets[i]
		kubeServiceAccount.Secrets = append(kubeServiceAccount.Secrets, v1.ObjectReference{
			Kind:            kokiSecret.Kind,
			Namespace:       kokiSecret.Namespace,
			Name:            kokiSecret.Name,
			UID:             apimachinery.UID(kokiSecret.UID),
			APIVersion:      kokiSecret.Version,
			ResourceVersion: kokiSecret.ResourceVersion,
			FieldPath:       kokiSecret.FieldPath,
		})
	}

	return kubeServiceAccount, nil
}
