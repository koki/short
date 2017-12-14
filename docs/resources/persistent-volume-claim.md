# Introduction

 PersistentVolumeClaim is a request for storage by a user

| API group | Resource | Kube Skeleton                                   |
|:----------|:---------|:------------------------------------------------|
| core/v1  | PersistentVolumeClaim |  [skel](../skel/persistent-volume-claim.kube.skel.yaml)         |

Here's an example Kubernetes PersistentVolumeClaim:
```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: myclaim
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 8Gi
  storageClassName: slow
  selector:
    matchLabels:
      release: "stable"
    matchExpressions:
      - {key: environment, operator: In, values: [dev]}
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:------|:-----|:--------|:-----------------------|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this PersistentVolumeClaim is running |
|name | `string` | `metadata.name`| The name of the PersistentVolumeClaim | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this PersistentVolumeClaim will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the PersistentVolumeClaim, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the PersistentVolumeClaim | 
|storage_class| `string` | `spec.storageClassName`| The number of storageclass required by the claim |
|volume | `string` | `spec.volumeName` | Binding reference to persistent volume claim holding this reference |
|access_modes | `[]string` | `spec.accessModes` | Desired access mode the volume should have. See [Access Modes](#access-modes) | 
|storage | `string` | `spec.resources.requests.limit` | Amount of storage the volume should have (eg. 4Gi)|
|selector | `map[string]string` or `string` | `selector` | An expression (string) or a set of key, value pairs (map) that is used to select a set of pods to manage using the PersistentVolumeClaim controller. See [Selector Overview](#selector-overview) |

#### Access Modes 

| Access Mode | Description |
|:----------------------|:------------|
| rw_once | Can be mounted read/write mode to exactly 1 host |
| ro_many | Can be mounted read only mode to many hosts |
| rw_many | Can be mounted read/write mode to many hosts |

#### Selector Overview

Selector can be a map value or a string value. If it is a string value, then it can be an expression of type

 - `Key=Value`

Valid Operators are 

| Operator | Syntax | Description         |
|:------|:-----|:-----------------------|
| Eq| `=` | Key should be equal to value |
| Exists| N/A | Key should exist | 
| NotExists| N/A | Key should not exist |
| In| `=` | Key should be one of the comma separated values |
| NotIn| `!=` | Key should not be one of the comma separated values |

Here are valid examples of all the expression operators
```yaml
selector: key=value # key should be equal to value
selector: key # key should exist
selector: !key # key should not exist
selector: key=value1,value2 # key's value can be any of value1 or value2
selector: key!=value1,value2 # key's value cannot be any of value1 or value2
selector: key&key!=value # composite expression
```

**Note that multiple expressions can be combined using the `&` symbol**

If the selector is a map, then the values in the map are expected to match directly with the labels of a pod. 

#### Template Metadata

| Field | Type | K8s counterpart(s) | Description         |
|:------|:-----|:--------|:-----------------------|
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this Pod is running |
|name | `string` | `metadata.name`| The name of the Pod | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this Pod will be a member of | 
|labels | `string` | `metadata.labels`| Metadata that could be identifying information about the Pod | 
|annotations| `string` | `metadata.annotations`| Non identifying information about the Pod| 

# Examples 

 - PersistentVolumeClaim requesting 20Gi of Storage with `rw_once` access mode

```yaml
pvc:
  access_modes:
  - rw_once
  labels:
    app: wordpress
  name: mysql-pv-claim
  storage: 20Gi
  version: v1
```

 - PersistentVolumeClaim requesting 1Mi of Storage with `rw_many` access mode

```yaml
pvc:
  access_modes:
  - rw_many
  name: nfs
  storage: 1Mi
  version: v1
```

# Skeleton

| Short Type           | Skeleton                                       |
|:---------------------|:-----------------------------------------------|
| PersistentVolumeClaim           | [skel](../skel/persistent-volume-claim.short.skel.yaml)     |

Here's a starter skeleton of a Short PersistentVolumeClaim.
```yaml
pvc:
  access_modes:
  - rw_once
  name: myclaim
  selector: release=stable&environment=dev
  storage: 8Gi
  storage_class: slow
  version: v1
```
