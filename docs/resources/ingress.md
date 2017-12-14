# Introduction

Ingress is a collection of rules that allow inbound connections to reach endpoints defined by a backend.

| API group | Resource | Kube Skeleton                                   |
|:----------|:---------|:------------------------------------------------|
| extensions/v1beta1   | Ingress  |  [skel](../skel/ingress.extensions.v1beta1.kube.skel.yaml)         |

Here's an example Kubernetes Ingress spec:
```yaml
ingress:
  name: echomap
  rules:
  - host: foo.bar.com
    paths:
    - path: /foo
      port: 80
      service: echoheaders-x
  - host: bar.baz.com
    paths:
    - path: /bar
      port: 80
      service: echoheaders-y
    - path: /foo
      port: 80
      service: echoheaders-x
  version: extensions/v1beta1
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:------|:-----|:--------|:-----------------------|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.` `clusterName` | The name of the cluster on which this Ingress is running |
|name | `string` | `metadata.name`| The name of the Ingress | 
|namespace | `string` | `metadata.` `namespace` | The K8s namespace this Ingress will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the Ingress, including identifying information | 
|annotations| `string` | `metadata.` `annotations`| Non-identifying information about the Ingress | 
|backend | `string` | `backend.serviceName` | Name of the referenced service |
|backend_port | `int` or `string` | `backend.servicePort` | Port of the referenced service |
|tls | `[]IngressTLS` | `spec.tls` | TLS configuration for this ingress. Currently only port 443 is supported. See [Ingress TLS](#ingress-tls) |
|rules | `[]IngressRule` | `spec.rules` | List of host rules used to configure ingress. If unspecified, all traffic is sent to default backend. See [Ingress Rule](#ingress-rule) | 

#### Ingress TLS

| Field | Type | K8s counterpart(s) | Description |
|:------|:------|:------------------|:------------|
|hosts | `[]string` | `hosts` | List of hosts included in the TLS cert |
|secret | `string` | `secretName` | Name of the secret used to terminate SSL |

#### Ingress Rule

| Field | Type | K8s counterpart(s) | Description |
|:------|:-----|:-------------------|:------------|
|host | `string` | `host` | FQDN of the host that should match incoming traffic and be processed by this rule |
|paths | `IngressPath` | `ingressRuleValue` | Paths and their respective backends. See [Ingress Path](#ingress-path) |

#### Ingress Path

| Field | Type | Description |
|:------|:-----|:------------|
|path | `string` | Path of the specified HTTP URL |
|service | `string` | Name of the service which should receive this traffic |
|port | `int` or `string` | Port of the Service at which it should receive traffic |

# Examples 

 - An ingress with multiple rules and tls config

```yaml
ingress:
  name: cafe-ingress
  rules:
  - host: cafe.example.com
    paths:
    - path: /tea
      port: 80
      service: tea-svc
    - path: /coffee
      port: 80
      service: coffee-svc
  tls:
  - hosts:
    - cafe.example.com
    secret: cafe-secret
  version: extensions/v1beta1
```

 - Ingress without TLS config

```yaml
ingress:
  annotations:
    ingress.kubernetes.io/auth-url: https://httpbin.org/basic-auth/user/passwd
  name: external-auth
  rules:
  - host: external-auth-01.sample.com
    paths:
    - path: /
      port: 80
      service: echoheaders
  version: extensions/v1beta1
```

# Skeleton

| Short Type           | Skeleton                                       |
|:---------------------|:-----------------------------------------------|
| Ingress                  | [skel](../skel/ingress.short.skel.yaml)            |

Here's a starter skeleton of a Short Ingress.
```yaml
ingress:
  name: cafe-ingress
  rules:
  - host: cafe.example.com
    paths:
    - path: /tea
      port: 80
      service: tea-svc
    - path: /coffee
      port: 80
      service: coffee-svc
  tls:
  - hosts:
    - cafe.example.com
    secret: cafe-secret
  version: extensions/v1beta1
```
