package types

type SecretWrapper struct {
	Secret Secret `json:"secret,omitempty"`
}

type Secret struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	StringData  map[string]string `json:"string_data,omitempty"`
	Data        map[string][]byte `json:"data,omitempty"`
	SecretType  SecretType        `json:"type,omitempty"`
}

type SecretType string

const (
	SecretTypeOpaque              SecretType = "opaque"
	SecretTypeServiceAccountToken SecretType = "kubernetes.io/service-account-token"
	SecretTypeDockercfg           SecretType = "kubernetes.io/dockercfg"
	SecretTypeDockerConfigJson    SecretType = "kubernetes.io/dockerconfigjson"
	SecretTypeBasicAuth           SecretType = "kubernetes.io/basic-auth"
	SecretTypeSSHAuth             SecretType = "kubernetes.io/ssh-auth"
	SecretTypeTLS                 SecretType = "kubernetes.io/tls"
)
