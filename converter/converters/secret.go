package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_v1_Secret_to_Koki_Secret(kubeSecret *v1.Secret) (*types.SecretWrapper, error) {
	kokiSecretWrapper := &types.SecretWrapper{}
	kokiSecret := types.Secret{}

	kokiSecret.Name = kubeSecret.Name
	kokiSecret.Namespace = kubeSecret.Namespace
	kokiSecret.Version = kubeSecret.APIVersion
	kokiSecret.Cluster = kubeSecret.ClusterName
	kokiSecret.Labels = kubeSecret.Labels
	kokiSecret.Annotations = kubeSecret.Annotations
	kokiSecret.Data = kubeSecret.Data
	kokiSecret.StringData = kubeSecret.StringData

	t, err := convertSecretType(kubeSecret.Type)
	if err != nil {
		return nil, err
	}
	kokiSecret.SecretType = t

	kokiSecretWrapper.Secret = kokiSecret
	return kokiSecretWrapper, nil
}

func convertSecretType(secret v1.SecretType) (types.SecretType, error) {
	if secret == "" {
		return "", nil
	}
	switch secret {
	case v1.SecretTypeOpaque:
		return types.SecretTypeOpaque, nil
	case v1.SecretTypeServiceAccountToken:
		return types.SecretTypeServiceAccountToken, nil
	case v1.SecretTypeDockercfg:
		return types.SecretTypeDockercfg, nil
	case v1.SecretTypeDockerConfigJson:
		return types.SecretTypeDockerConfigJson, nil
	case v1.SecretTypeBasicAuth:
		return types.SecretTypeBasicAuth, nil
	case v1.SecretTypeSSHAuth:
		return types.SecretTypeSSHAuth, nil
	case v1.SecretTypeTLS:
		return types.SecretTypeTLS, nil
	default:
		return "", serrors.InvalidValueErrorf(secret, "unrecognized Secret type")
	}
}

func Convert_Koki_Secret_to_Kube_v1_Secret(kokiSecretWrapper *types.SecretWrapper) (*v1.Secret, error) {
	kubeSecret := &v1.Secret{}
	kokiSecret := kokiSecretWrapper.Secret

	kubeSecret.Name = kokiSecret.Name
	kubeSecret.Namespace = kokiSecret.Namespace
	kubeSecret.APIVersion = kokiSecret.Version
	kubeSecret.ClusterName = kokiSecret.Cluster
	kubeSecret.Kind = "Secret"
	kubeSecret.Labels = kokiSecret.Labels
	kubeSecret.Annotations = kokiSecret.Annotations
	kubeSecret.Data = kokiSecret.Data
	kubeSecret.StringData = kubeSecret.StringData

	t, err := revertSecretType(kokiSecret.SecretType)
	if err != nil {
		return nil, err
	}
	kubeSecret.Type = t

	return kubeSecret, nil
}

func revertSecretType(secret types.SecretType) (v1.SecretType, error) {
	if secret == "" {
		return "", nil
	}

	switch secret {
	case types.SecretTypeOpaque:
		return v1.SecretTypeOpaque, nil
	case types.SecretTypeServiceAccountToken:
		return v1.SecretTypeServiceAccountToken, nil
	case types.SecretTypeDockercfg:
		return v1.SecretTypeDockercfg, nil
	case types.SecretTypeDockerConfigJson:
		return v1.SecretTypeDockerConfigJson, nil
	case types.SecretTypeBasicAuth:
		return v1.SecretTypeBasicAuth, nil
	case types.SecretTypeSSHAuth:
		return v1.SecretTypeSSHAuth, nil
	case types.SecretTypeTLS:
		return v1.SecretTypeTLS, nil
	default:
		return "", serrors.InvalidValueErrorf(secret, "unrecognized Secret type")
	}
}
