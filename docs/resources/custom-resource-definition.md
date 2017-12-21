# Introduction

Custom Resource Definition (CRD) represents a resource type to expose in the API server.

```yaml
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  # name must match the spec fields below, and be in the form: <plural>.<group>
  name: crontabs.stable.example.com
spec:
  group: stable.example.com
  version: v1
  scope: Namespaced
  names:
    plural: crontabs
    singular: crontab
    kind: CronTab
    shortNames:
    - ct
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview


| Field | Type | K8s counterpart(s) | Description         |
|:------|:-----|:--------|:-----------------------|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this CRD is running |
|name | `string` | `metadata.name`| The name of the CRD. It must be of the form `<plural>.<group>` | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this CRD will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the CRD, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the CRD |
|meta| `CRD Meta` | `spec.group` `spec.version` `spec.names` | Metadata about the resource defined by the CRD. See [CRD Meta](#crd-meta) |
|scope| `"namespaced"` or `"cluster"` | `spec.scope` | Whether the resource is namespaced or cluster-scoped. Defaults to `"namespaced"` if omitted. |
|validation| `JSONSchemaProps` | Optional OpenAPI schema for the defined resource type. |
|conditions| `[]CRD Condition`| `status.conditions` | The list of current and previous conditions of the CRD. See [CRD Condition](#crd-condition) |
|accepted| `CRD Names` | `spec.names` | The names actually being used for the discovery service. See [CRD Names](#crd-meta) |

#### CRD Meta

`CRD Meta` is `CRD Names` + `group` and `version`.

| Field | Type |  Description         |
|:------|:-----|:--------|
|group| `string` | The API group for the API resource. e.g. `core` |
|version| `string` | The API version for the API resource definition. e.g. `v1` |
|plural| `string` | Lowercase plural name for the API resource. |
|singular| `string` | Lowercase singular name for the API resource. Defaults to lowercase of `kind`. |
|short| `[]string` | Lowercase abbreviated names for use in the command line. |
|kind| `string` | Capitalized camel-case name for the resource. Usually singular. e.g. `Pod` |
|list| `string`| Defaults to `<kind>List`. e.g. `PodList` |

#### CRD Condition

| Field | Type |  Description         |
|:------|:-----|:--------|
| reason| `string` | One word camel case reason for CRD's last transition |
| msg | `string` | Human readable message about the CRD's last transition |
| status | `ConditionStatus` | String value that represents the status of the condition. Can be "True", "False" or "Unknown" |
| type | `CRDConditionType` | String value that represents the type of condition. Can be "established", "names-accepted" or "terminating" |
| last_change | `time` | Last time the condition status changed |

