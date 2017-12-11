# Introduction

 Job is a declarative interface to denote a set of Pods that are expected to run to completion.

| API group | Resource | Kube Skeleton                                   |
|:---------:|:--------:|:-----------------------------------------------:|
| batch/v1  | Job |  [skel](../skel/job.batch.v1.kube.skel.yaml)         |

Here's an example Kubernetes Job:
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    metadata:
      name: pi
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this Job is running |
|name | `string` | `metadata.name`| The name of the Job | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this Job will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the Job, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the Job |
|parallelism | `int32` | `spec.parallelism` | Maximum number of pods of this job that can run in parallel  |
|completions| `int32` | `spec.completions` | Minimum number of successfully completed pods for the job to be considered successful |
|max_retries | `int32` | `spec.backOffLimit` | Maximum number of retries before considering this job failed |
|active_deadline | `int32` | `spec.activeDeadlineSeconds` | Maximum time for a job to be active before it is terminated by the system |
|select_manually | `bool` | `spec.manualSelector` | Controls generation of pod labels and pod selectors. Defaults to false. If set, user is responsible for choosing pods correctly |
|selector | `map[string]string` or `string` | `selector` | An expression (string) or a set of key, value pairs (map) that is used to select a set of pods to manage using the job controller. See [Selector Overview](#selector-overview) |
|pod_meta | `TemplateMetadata` | `template` | Metadata of the Pod that is selected by this Job. See [Template Metadata](#template-metadata)|
|volumes | `Volume` | `spec.volumes` | Denotes the volumes that are a part of the Pod. See [Volume Overview](pod#volume-overview) |
| affinity | `[]Affinity` | `spec.affinity` and `spec.NodeSelector` | The Pod's scheduling rules, expressed as (anti-)affinities for nodes or other Pods. See [Affinity Overview](pod#affinity-overview) |
| node | `string` | `spec.nodeName` | Request that the Pod be scheduled on a specific node. | 
| containers |`Container` | `spec.containers` and `status`| Containers that run as a part of the Pod. See [Container Overview](pod#container-overview) |
| init_containers | `Container` | `spec.initContainers` and `status` | Containers that run as a part of the initialization process of the Pod. See [Container Overview](pod#container-overview) | 
| dns_policy | `DNSPolicy` | `spec.dnsPolicy` | The DNS Policy of the Pod. See [DNS Policy Overview](pod#dns-policy-overview) |
| host_aliases | `[]string` | `spec.aliases` | Set of additional records to be placed in `/etc/hosts` file inside the Pod. See [Host Aliases Overview](pod#host-aliases-overview) |
| host_mode | `[]string` | `spec.hostPID`, `spec.hostNetwork` and `spec.hostIPC`| The Pod's access to host resources. See [Host Mode Conversion](pod#host-mode-conversion) |
| hostname | `string` | `spec.hostname` and `spec.subDomain` | The fully qualified domain name of the pod|
| registry_secrets | `[]string` |`spec.ImagePullSecrets` | A list of k8s secret resource names that contain credentials to required to access private registries. |
| restart_policy | `RestartPolicy` | `spec.restartPolicy` | Behavior of a Pod when it dies. Can be "always", "on-failure" or "never" |
| scheduler_name | `string` | `spec.schedulerName` | The value from `spec.schedulerName` is stored here |
| account | `string` | `spec.serviceAccountName` and `automountService` `AccountToken` | The Pod's access to the K8s API. See [Account Conversion](pod#account-conversion) | 
| tolerations | `[]Toleration` | `spec.tolerations` | Set of host taints this Pod tolerates. See [Toleration Conversion](pod#toleration-conversion) |
| termination_ grace_period | `int64`  | `spec.termination` `GracePeriodSeconds` | Number of seconds to wait before forcefully killing the Pod. |
| active_deadline | `int64` | `spec.` `activeDeadlineSeconds`| Number of seconds the Pod is allowed to be active  |  
| priority | `Priority` | `spec.priorityClassName` and `spec.priority` | Specifies the Pod's Priority. See [Priority](pod#priority) |
| condition | `[]Pod Condition` | `status.conditions` | The list of current and previous conditions of the Pod. See [Pod Condition](pod#pod-condition) |
| node_ip | `string` | `status.hostIP` | The IP address of the Pod's host | 
| ip | `string` | `status.podIP` | The IP address of the Pod | 
| start_time | `time` | `status.startTime` | When the Pod started running | 
| msg | `string` | `status.message` | A human readable message explaining Pod's current condition |  
| phase | `string` | `status.phase` | The current phase of the Pod |
| reason | `string` | `status.reason` | Reason indicating the cause for the current state of the Pod |
| qos | `string` | `status.qosClass` | The QOS class assigned to the Pod based on resource requirements |
| fs_gid | `int64` | `spec.securityContext.` `fsGroup` | Special supplemental group that applies to all the Containers in the Pod |
| gids | `[]int64` | `spec.securityContext.` `supplementalGroups` | A list of groups applied to the first process in each of the Containers in the Pod |


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

 -  An example non-parallel job selecting pods with labels app:nginx

```yaml
job:
  containers:
  - expose:
    - 80
    image: nginx:1.7.9
    name: nginx
  name: nginx-job
  parallelism: 1
  completions: 1
  selector:
    app: nginx
  version: batch/v1
```

 - An example parallel job with fixed completion count

```yaml
job:
  containers:
  - expose:
    - 80
    image: nginx:1.7.9
    name: nginx
  name: nginx-job
  parallelism: 2
  completions: 4
  selector:
    app: nginx
  version: batch/v1
```

 - An example parallel job work queue

```yaml
job:
  containers:
  - expose:
    - 80
    image: nginx:1.7.9
    name: nginx
  name: nginx-job
  parallelism: 2
  completions: 1
  selector:
    app: nginx
  version: batch/v1
```

# Skeleton

| Short Type           | Skeleton                                       |
|:--------------------:|:----------------------------------------------:|
| Job                  | [skel](../skel/job.short.skel.yaml)            |

Here's a starter skeleton of a Short Job.
```yaml
job:
  containers:
  - command:
    - perl
    - -Mbignum=bpi
    - -wle
    - print bpi(2000)
    image: perl
    name: pi
  max_retries: 4
  name: pi
  pod_meta:
    name: pi
  restart_policy: never
  version: batch/v1
```
