# Introduction

ControllerRevision is an immutable snapshot of state.
It's primarily intended for internal use by controllers.
For example, it's used by the DaemonSet and StatefulSet controllers for update and rollback.

Here's an example Kubernetes ControllerRevision:
```yaml
apiVersion: apps/v1
kind: ControllerRevision
metadata:
  name: example
data:
  key: value
revision: 1
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:------|:-----|:--------|:-----------------------|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this Job is running |
|name | `string` | `metadata.name`| The name of the Job | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this Job will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the Job, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the Job |
|data| YAML | `data` | This field can hold any valid YAML. |
|revision| `int64` | The revision number |

# Examples / Skeleton

Here's a starter skeleton of a Short ControllerRevision.
```yaml
controller_revision:
  name: example
  data:
    key: value
  revision: 1
```
