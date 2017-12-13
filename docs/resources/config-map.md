# Introduction

ConfigMap holds configuration data for pods to consume

| API group | Resource | Kube Skeleton                                   |
|:---------:|:--------:|:-----------------------------------------------:|
| core/v1  | ConfigMap |  [skel](../skel/config-map.kube.skel.yaml)         |

Here's an example Kubernetes ConfigMap:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: special-config
  namespace: default
data:
  special.how: very
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this ConfigMap is running |
|name | `string` | `metadata.name`| The name of the ConfigMap | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this ConfigMap will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the ConfigMap, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the ConfigMap | 
|data| `map[string]string` | `data`| Configuration Data |

# Examples 

 - ConfigMap example

```yaml
config_map:
  data:
    log_level: INFO
  name: env-config
  namespace: default
  version: v1
```

# Skeleton

| Short Type           | Skeleton                                       |
|:--------------------:|:----------------------------------------------:|
| ConfigMap           | [skel](../skel/config-map.short.skel.yaml)     |

Here's a starter skeleton of a Short ConfigMap.
```yaml
config_map:
  data:
    log_level: INFO
  name: env-config
  namespace: default
  version: v1
```
