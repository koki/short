# Introduction

 Deployment is a declarative interface to upgrade Pods and ReplicaSets.

| API group | Resource | Kube Skeleton                                   |
|:---------:|:--------:|:-----------------------------------------------:|
| extensions/v1beta1  | Deployment |  [skel](../skel/deployment.extensions.v1beta1.kube.skel.yaml)         |
| apps/v1beta1  | Deployment |  [skel](../skel/deployment.apps.v1beta1.kube.skel.yaml)         |
| apps/v1beta2  | Deployment |  [skel](../skel/deployment.apps.v1beta2.kube.skel.yaml)         |

Here's an example Kubernetes Deployment:
```yaml
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```

The following sections contain detailed information about the process of converting each of the fields within the Kubernetes Deployment definition to Short spec and back.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this Pod is running |
|name | `string` | `metadata.name`| The name of the Pod | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this Pod will be a member of | 
|labels | `string` | `metadata.labels`| Metadata that could be identifying information about the Pod | 
|annotations| `string` | `metadata.annotations`| Non identifying information about the Pod| 
|replicas| `int32` | `replicas`| The number of replicas of the selected pod|
|recreate | `bool` | `strategy` | The strategy for performing upgrades. If recreate is set to `true`, then the strategy is `Recreate`. If not, by default the strategy is `RollingUpdate` |
|max_unavailable | `int` or `string` | `strategy.rollingUpdate` | Maximum number of pods that can be unavailable during update. More information below| 
|max_extra |   `int` or `string` | `strategy.rollingUpdate` | Maximum number of pods that can exceed the replica count during update. More information below|
|min_ready | `int32` | `minReadySeconds` | Minimum number of seconds that your pod should be ready before it is considered available |
|max_revs | `int32` | `revisionHistoryLimit` | Number of old replica sets to retain to allow rollback|
|paused | `bool` | `paused` | Prevent deployment from being managed by Kubernetes by setting this to true|
|progress_deadline | `int32` | `progressDeadlineSeconds` | Maximum time for a deployment to make progress before it is considered unavailable|
|selector | `map[string]string` or `string` | `selector` | An expression (string) or a set of key, value pairs (map) that is used to select a set of pods to manage using the deployment controller. More information available in [Selector Overview](#selector-overview) |
|pod_meta | `TemplateMetadata` | `template` | Metadata of the pod that is selected by this deployment. More details in [Template Metadata](#template-metadata)|
|volumes | [`Volume`](./pod.md#volume-overview) | `spec.volumes` | Denotes the volumes that are a part of the Pod. More information is available in [Volume Overview](./pod.md#volume-overview) |
| affinity | [`[]Affinity`](#affinity-overview) | `spec.affinity` and `spec.NodeSelector` | (Anti-) Affinity of the Pod to nodes or other Pods. More information is available in [Affinity Overview](./pod.md#affinity-overview) |
| containers |[`Container`](../skel/container.short.skel.yaml) | `spec.containers` and `status`| Containers that run as a part of the Pod. More information is available in [Container Overview](./pod.md#container-overview) |
| init_containers | [`Container`](../skel/container.short.skel.yaml) | `spec.initContainers` and `status` | Containers that run as a part of the initialization process of the Pod. More information is available in [Container Overview](./pod.md#container-overview) | 
| dns_policy | [`DNSPolicy`](./pod.md#dns-policy-overview) | `spec.dnsPolicy` | The DNS Policy of the Pod. The conversion function is explained in [DNS Policy Overview](./pod.md#dns-policy-overview) |
| host_aliases | `[]string` | `spec.aliases` | Set of additional records to be placed in `/etc/hosts` file inside the Pod. More information is available in [Host Aliases Overview](./pod.md#host-aliases-overview) |
| host_mode | `[]string` | `spec.hostPID`, `spec.hostNetwork` and `spec.hostIPC`| The access the Pod has to host resources. The conversion function is explained in [Host Mode Conversion](./pod.md#host-mode-conversion) |
| hostname | `string` | `spec.hostname` and `spec.subDomain` | The fully qualified domain name of the pod|
| registry_secrets | `[]string` |`spec.ImagePullSecrets` | A list of k8s secret resource names that contain credentials to required to access private registries. |
| restart_policy | [`RestartPolicy`](./pod.md#restart-policy) | `spec.restartPolicy` | Behavior of a Pod when it dies. Can be "Always", "OnFailure" or "Never" |
| scheduler_name | `string` | The value from `spec.schedulerName` is stored here | The value from `spec.schedulerName` is stored here |
| account | `string` | `spec.serviceAccountName` and `spec.automountServiceAccountToken` | The access the Pod gets to the K8s API. More information is available in [Account Conversion](./pod.md#account-conversion) | 
| tolerations | [`[]Toleration`](../skel/toleration.short.skel.yaml) | `spec.tolerations` | Set of tainted hosts to tolerate on scheduling the Pod. The conversion function is explained in [Toleration Conversion](./pod.md#toleration-conversion) |
| termination_grace_period | `int64`  | `spec.terminationGracePeriodSeconds` | Number of seconds to wait before forcefully killing the Pod. |
| active_deadline | `int64` | `spec.activeDeadlineSeconds`| Number of seconds the Pod is allowed to be active  |  
| node | `string` | `spec.nodeName` | Request Pod to be scheduled on node with the name specified in this field| 
| priority | `Priority` | `spec.priorityClassName` and `spec.priority` | Specifies the Pod's Priority. More information in [Priority](./pod.md#priority) |
| condition | `[]Pod Condition` | `status.conditions` | The list of current and previous conditions of the Pod. More information in [Pod Condition](./pod.md#pod-condition) |
| node_ip | `string` | `status.hostIP` | The IP address of the host on which the Pod is running | 
| ip | `string` | `status.podIP` | The IP address of the Pod | 
| start_time | `time` | `status.startTime` | The time at which the Pod started running | 
| msg | `string` | `status.message` | A human readable message explaining Pod's current condition |  
| phase | `string` | `status.phase` | The current phase of the Pod |
| reason | `string` | `status.reason` | Reason indicating the cause for the current state of the Pod |
| qos | `string` | `status.qosClass` | The QOS class assigned to the Pod based on resource requirements |
| fs_gid | `int64` | `spec.securityContext.fsGroup` | Special supplemental group that apply to all the Containers in the Pod |
| gids | `[]int64` | `spec.securityContext.supplementalGroups` | A list of groups applied to the first process in each of the containers of the Pod |

`max_unavailable` and `max_extra` are used to configure `RollingUpdate` deployment strategy. `max_unavailable` indicates the maximum number of pods that can be unavailable during update. The value can be number or a percentage value of the total number of replicas. Percentage values are represented using a `%` symbol at the end of the value.

```yaml
max_extra: 3
max_unavailable: 30%
```

#### Selector Overview

Selector can be a map value or a string value. If it is a string value, then it can be an expression of type

 - `Key=Value`

Valid Operators are 

| Operator | Syntax | Description         |
|:-----:|:----:|:----------------------:|
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
|:-----:|:----:|:-------:|:----------------------:|
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this Pod is running |
|name | `string` | `metadata.name`| The name of the Pod | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this Pod will be a member of | 
|labels | `string` | `metadata.labels`| Metadata that could be identifying information about the Pod | 
|annotations| `string` | `metadata.annotations`| Non identifying information about the Pod| 

# Examples 

 -  An example deployment with 1 replica selecting app:nginx

```yaml
deployment:
  containers:
  - expose:
    - 80
    image: nginx:1.7.9
    name: nginx
  name: nginx-deployment
  replicas: 1
  selector:
    app: nginx
  version: apps/v1beta1
```

 - An example deployment with deployment strategy set to Recreate and replicas set to 2

```yaml
deployment:
  containers:
  - expose:
    - 3001
    image: gcr.io/my-project/foo-service:3adfa3e
    name: foo-web
  labels:
    app: foo
  name: foo-web
  replicas: 2
  recreate: true   # if this is not set, then default of Rolling Update strategy is used
  selector:
    app: foo
    tier: backend
  version: extensions/v1beta1
```

 - An example deployment that selects on labels app=nginx and app=haproxy

```yaml
deployment:
  containers:
  - expose:
    - 80
    image: nginx:1.7.9
    name: nginx
  labels:
    app: nginx
  name: nginx-deployment
  replicas: 3
  selector: app=nginx,haproxy  # string selector (expression)
  version: apps/v1beta2
```

# Skeleton

| Short Type           | Skeleton                                       |
|:--------------------:|:----------------------------------------------:|
| Deployment           | [skel](../skel/deployment.short.skel.yaml)     |

Here's a starter skeleton of a Short Deployment.
```yaml
deployment:
  containers:
  - expose:
    - 80
    image: nginx:1.7.9
    name: nginx
  name: nginx-deployment
  replicas: 3
  selector:
    app: nginx
  version: apps/v1beta1
```
