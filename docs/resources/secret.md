# Introduction

Secret holds secret data for pods to consume

| API group | Resource | Kube Skeleton                                   |
|:----------|:---------|:------------------------------------------------|
| core/v1  | Secret |  [skel](../skel/secret.kube.skel.yaml)         |

Here's an example Kubernetes Secret:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: mysecret
type: Opaque
data:
  username: YWRtaW4=
  password: MWYyZDFlMmU2N2Rm
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:------|:-----|:--------|:-----------------------|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this Secret is running |
|name | `string` | `metadata.name`| The name of the Secret | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this Secret will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the Secret, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the Secret | 
|data| `map[string][]byte` | `data`| Base64 encoded secret data |
|string_data| `map[string]string` | `stringData` | Non-Binary secret data in string form can be stored using this field|
|type | `string` | `secretType` | Types used to facilitate programmatic handling of secrets. See [Secret Types](#secret-types) | 

#### Secret Types

| Secret Type | Description |
|:------------|:------------|
| opaque | Default type. Arbitrary user defined data |
| kubernetes.io/service-account-token | Secret contains a token that identifies a service account to the API. See [Service Account Secrets](#service-account-secrets) |
| kubernetes.io/dockercfg | Secret contains a dockercfg file that follows the same format rules as `~/.dockercfg`. See [Docker Config Secrets](#docker-config-secrets) |
| kubernetes.io/dockerconfigjson | Secret contains a dockercfg file that follows the same format rules as ~/.docker/config.json. See [Docker Config JSON secrets](#docker-config-json-secrets)|
| kubernetes.io/basic-auth | Secret contains credentials for basic auth. See [Basic Auth Secrets](#basic-auth-secrets)|
| kubernetes.io/ssh-auth | Secret contains credentials for SSH auth. See [SSH Auth Secrets](#ssh-auth-secrets)|
| kubernetes.io/tls | Secret contains information about TLS server or client certificate. See [TLS Secrets](#tls-secrets)|

#### Service Account Secrets

If the secret type is set to `kubernetes.io/service-account-token`, then the secret should have the following required fields

| Field | Description |
|:------|:------------|
| Secret.Annotations["kubernetes.io/service-account.name"] | The name of the ServiceAccount the token identifies|
| Secret.Annotations["kubernetes.io/service-account.uid"]  | the UID of the ServiceAccount the token identifies |
| Secret.Data["token"]  | a token that identifies the service account to the API ||

#### Docker Config Secrets

If the secret type is set to `kubernetes.io/dockercfg`, then the secret should have the following required field

| Field | Description |
|:------|:------------|
| Secret.Data[".dockercfg"] | A serialized ~/.dockercfg file | 

#### Docker Config JSON Secrets

If the secret type is set to `kubernetes.io/dockerconfigjson`, then the secret should have the following required field

| Field | Description |
|:------|:------------|
| Secret.Data[".dockerconfigjson"] | A serialized ~/.docker/config.json file | 

#### Basic Auth Secrets

If the secret type is set to `kubernetes.io/basic-auth`, then the secret should have atleast one of the following fields

| Field | Description |
|:------|:------------|
| Secret.Data["username"] | Username used for authentication | 
| Secret.Data["password"] | Password or token needed for authentication |

#### SSH Auth Secrets

If the secret type is set to `kubernetes.io/ssh-auth`, then the secret should have the following required field

| Field | Description |
|:------|:------------|
| Secret.Data["ssh-privatekey"] | Private SSH key needed for authentication|

#### TLS Secrets

If the secret type is set to `kubernetes.io/tls`, then the secret should have the following required fields

| Field | Description |
|:------|:------------|
| Secret.Data["tls.key"] | TLS private key |
| Secret.Data["tls.crt"] | TLS certificate |

# Examples 

 - Secret example

```yaml
secret:
  data:
    password: MWYyZDFlMmU2N2Rm
    username: YWRtaW4=
  name: mysecret
  type: opaque
  version: v1
```

# Skeleton

| Short Type           | Skeleton                                       |
|:---------------------|:-----------------------------------------------|
| Secret          | [skel](../skel/secret.short.skel.yaml)     |

Here's a starter skeleton of a Short Secret.
```yaml
secret:
  data:
    password: MWYyZDFlMmU2N2Rm
    username: YWRtaW4=
  name: mysecret
  type: opaque
  version: v1
```
