# Introduction

A Pod is the unit of execution in Kubernetes. It consists of a set of co-located containers that share the same fate. The Pod definition in Kubernetes includes information about the containers, their runtime characteristics, and metadata about the pod.

| API group | Resource | Kube Skeleton                                   |
|:---------:|:--------:|:-----------------------------------------------:|
| core/v1   | Pod      |  [skel](../skel/pod.kube.skel.yaml)             |

Here's an example Kubernetes Pod spec:
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: mysql-pod
  labels:
    name: mysql-pod
spec:
  containers:
    -
      name: mysql
      image: mysql:latest
      env:
        -
          name: "MYSQL_USER"
          value: "mysql"
        -
          name: "MYSQL_PASSWORD"
          value: "mysql"
        -
          name: "MYSQL_DATABASE"
          value: "sample"
        -
          name: "MYSQL_ROOT_PASSWORD"
          value: "supersecret"
      ports:
        -
          containerPort: 3306
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this Pod is running |
|name | `string` | `metadata.name`| The name of the Pod | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this Pod will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the Pod, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the Pod | 
|volumes | `Volume` | `spec.volumes` | Denotes the volumes that are a part of the Pod. See [Volume Overview](#volume-overview) |
| affinity | `[]Affinity` | `spec.affinity` and `spec.NodeSelector` | The Pod's scheduling rules, expressed as (anti-)affinities for nodes or other Pods. See [Affinity Overview](#affinity-overview) |
| node | `string` | `spec.nodeName` | Request that the Pod be scheduled on a specific node. | 
| containers |`Container` | `spec.containers` and `status`| Containers that run as a part of the Pod. See [Container Overview](#container-overview) |
| init_containers | `Container` | `spec.initContainers` and `status` | Containers that run as a part of the initialization process of the Pod. See [Container Overview](#container-overview) | 
| dns_policy | `DNSPolicy` | `spec.dnsPolicy` | The DNS Policy of the Pod. See [DNS Policy Overview](#dns-policy-overview) |
| host_aliases | `[]string` | `spec.aliases` | Set of additional records to be placed in `/etc/hosts` file inside the Pod. See [Host Aliases Overview](#host-aliases-overview) |
| host_mode | `[]string` | `spec.hostPID`, `spec.hostNetwork` and `spec.hostIPC`| The Pod's access to host resources. See [Host Mode Conversion](#host-mode-conversion) |
| hostname | `string` | `spec.hostname` and `spec.subDomain` | The fully qualified domain name of the pod|
| registry_secrets | `[]string` |`spec.ImagePullSecrets` | A list of k8s secret resource names that contain credentials to required to access private registries. |
| restart_policy | `RestartPolicy` | `spec.restartPolicy` | Behavior of a Pod when it dies. Can be "always", "on-failure" or "never" |
| scheduler_name | `string` | `spec.schedulerName` | The value from `spec.schedulerName` is stored here |
| account | `string` | `spec.serviceAccountName` and `automountService` `AccountToken` | The Pod's access to the K8s API. See [Account Conversion](#account-conversion) | 
| tolerations | `[]Toleration` | `spec.tolerations` | Set of host taints this Pod tolerates. See [Toleration Conversion](#toleration-conversion) |
| termination_ grace_period | `int64`  | `spec.termination` `GracePeriodSeconds` | Number of seconds to wait before forcefully killing the Pod. |
| active_deadline | `int64` | `spec.` `activeDeadlineSeconds`| Number of seconds the Pod is allowed to be active  |  
| priority | `Priority` | `spec.priorityClassName` and `spec.priority` | Specifies the Pod's Priority. See [Priority](#priority) |
| condition | `[]Pod Condition` | `status.conditions` | The list of current and previous conditions of the Pod. See [Pod Condition](#pod-condition) |
| node_ip | `string` | `status.hostIP` | The IP address of the Pod's host | 
| ip | `string` | `status.podIP` | The IP address of the Pod | 
| start_time | `time` | `status.startTime` | When the Pod started running | 
| msg | `string` | `status.message` | A human readable message explaining Pod's current condition |  
| phase | `string` | `status.phase` | The current phase of the Pod |
| reason | `string` | `status.reason` | Reason indicating the cause for the current state of the Pod |
| qos | `string` | `status.qosClass` | The QOS class assigned to the Pod based on resource requirements |
| fs_gid | `int64` | `spec.securityContext.` `fsGroup` | Special supplemental group that applies to all the Containers in the Pod |
| gids | `[]int64` | `spec.securityContext.` `supplementalGroups` | A list of groups applied to the first process in each of the Containers in the Pod |

#### Affinity Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| node | `string` | `affinity.nodeAffinity` | The Pod's affinity for certain nodes. More information below | 
| pod | `string` | `affinity.podAffinity` | The Pod's affinity for certain other other Pods in the cluster. More information below |
| anti_pod | `string` | `affinity.podAntiAffinity` | The Pod's anti-affinity for certain other Pods in the cluster. More information below |
| topology | `string` | `affinity.pod*.` `podAffinityTerm.topologyKey` | A node label key, e.g. "kubernetes.io/hostname". Determines the scope (same host vs same region vs ...) of the Pod's (anti-)affinity for certain Pods. More information below| 
| namespaces | `[]string` | `affinity.pod*.` `podAffinityTerm.namespaces` | A list of namespaces in which the `pod` and `anti_pod` affinities are applied |

`node`, `pod` and `anti_pod` are string fields that expect `expressions` that denote affinities of the Pod to nodes, other pods and anti-affinity to other pods respectively.

`expressions` are label selectors, i.e. node or pods that contain labels that match these expressions are used to make scheduling decisions for the Pod.

An expression is a set of sub-expressions that are ANDed(`&`) together. Sub-expressions select on labels using `=`, `!=`, `Exists`, `Does not Exist`, `Greater than` and `Less than` operators.

| Operator | Symbol | Validity | Example |
|:--------:|:------:|:--------:|:--------:|
| Equals   | `=`   | node, pod and anti_pod | `k8s.io/failure-domain=us-east1` |
| Not Equal | `!=`  | node, pod and anti_pod | `k8s.io/failure-domain!=us-east-1`  |
| Exists | N/A  |  node, pod and anti_pod | `k8s.io/cloud-provider` |
| Does Not Exist | N/A | node, pod and anti_pod | `!k8s.io/bare-metal`  |
| Greater Than | '>' | node | `k8s.io/cpus>1` |
| Less Than | '<' | node | `k8s.io/cpus < 1`|

Expressions also have qualifiers at the end of the composite sub-expressions. Qualifiers can be used to set `soft` affinity and (weight) of the soft affinity. `soft` affinities have weights ranging from 1 to 100, where 1 is the default weight.

Pods accept multiple affinity items, and the entire set of affinity items is considered for its scheduling. 

*It is not valid to include more than one of (node, pod, anti_pod) in a single affinity item. They should be specified in separate items. `topology` and `namespaces` are ignored if the affinity item is a `node` selector affinity item*

#### Node Affinity
`node` values denote either `soft` or `hard` affinities to nodes.
`hard` affinities are expressions that must be satisfied for a Pod to be scheduled on a node.
`soft` affinities are expressions that we prefer to satisfy,
but Pods may be scheduled on nodes that do not satisfy these expressions.

If the list of `affinity` items contains multiple `hard` node affinity selectors, only one `hard` node affinity selector needs to be satisfied.

If the list of `affinity` items contains multiple `soft` node affinity selectors, the scheduler chooses the node that satisfies the most `soft` node affinity selectors, where "most" means the greatest sum of weights.

Here are some example node affinity expressions

| Expression | Affinity Type | Description |
|:----------|:-------------:|:-----------:|
|-node:`failure-domain=us-east1&instance-type=t2.large` | `node` hard affinity  | run the pod on a node that is in failure domain `us-east1` and whose instance type is `t2.large` |
|-node:`failure-domain=us-east1&instance-type=t2.large`<br/>-node:`failure-domain=us-east2&instance-type=t2.large` | `node` hard affinity | run the pod on a node that is in failure-domain `us-east1` and the instance type is `t2.large` <br/> or <br/> on a node in `us-east2` and instance type is `t2.large`|
|-node:`failure-domain=us-east1&instance-type=t2.large:soft` | `node` soft affinity | prefer to run the pod on a node that is in failure-domain `us-east1` and whose instance type is `t2.large` |
|-node:`failure-domain=us-east1:soft:10`<br/>-node:`failure-domain=us-east2:soft:20` | `node` soft affinity | run the pod preferrably on a node in failure-domain `us-east2`, less preferrably on a node in `us-east1`, some other node if none of those options are available|

#### Pod Affinity 
`pod` values denote either `soft` or `hard` affinities to other pods.

A `hard` pod affinity selector indicates that the Pod must run alongside a Pod that satisfies the selector.
If the list of `affinity` items contains multiple `hard` pod affinity selectors, all `hard` pod affinity selectors must be satisfied.

A `soft` pod affinity selector indicates that the Pod prefers to run alongside a Pod that satisfies the selector.
If the list of `affinity` items contains multiple `soft` pod affinity selectors, the scheduler tries to satisfy as many `soft` pod affinity selectors
as possible, where "many" is measured by the sum of the selectors' weights.

Here are some example pod affinity expressions

| Expression | Affinity Type | Description |
|:----------|:-------------:|:-----------:|
|-pod:`app=front-end` | `pod` hard affinity  | run the pod on alongside another pod which has label `app=front-end` |
|-pod:`app=front-end`<br/>-pod:`name=react` | `pod` hard affinity | run the pod alongside another pod that has labels `app=front-end` and `name=react` |
|-pod:`app=front-end&name=react:soft` | `pod` soft affinity | prefer to run the pod alongside another pod that has labels `app=front-end` and `name=react` |
|-pod:`app=front-end&name=react:soft:10`<br/>-pod:`app=front-end&name=flux:soft:20` | `pod` soft affinity | run the pod preferrably alongside another pod that has labels `app=front-end` and `name=flux`, less preferrably alongside another pod that has labels `app=front-end` and `name=react`, some other node if none of those options are available|
|-pod:`app=front-end`<br/> topology:`k8s.io/failure-domain` |`pod` hard affinity | run the pod on a node whose label value for the key `k8s.io/failure-domain` matches the value of the label in the node on which a pod with label `app-front-end` is running | 

#### Pod Anti Affinity
The syntax and mechanism of pod anti affinity is the same as pod affinity, except that if another Pod matches the selector, this Pod should *not* be scheduled alongside it.

Here are some example pod anti affinity expressions

| Expression | Affinity Type | Description |
|:----------|:-------------:|:-----------:|
|-anti_pod:`app=front-end` | `pod` hard anti-affinity  | never run the pod on alongside another pod which has label `app=front-end` |
|-anti_pod:`app=front-end`<br/>-anti_pod:`name=react` | `pod` hard anti-affinity | never run the pod alongside another pod that has labels `app=front-end` and `name=react` |
|-anti_pod:`app=front-end&name=react:soft` | `pod` soft anti-affinity | prefer to NOT run the pod alongside another pod that has labels `app=front-end` and `name=react` |
|-anti_pod:`app=front-end&name=react:soft:10`<br/>-anti_pod:`app=front-end&name=flux:soft:20` | `pod` soft anti-affinity | run the pod preferrably not alongside another pod that has labels `app=front-end` and `name=flux`, less preferrably not alongside another pod that has labels `app=front-end` and `name=react`, some other node if none of those options are available|
|-anti_pod:`app=front-end`<br/> topology:`k8s.io/hostname` |`pod` hard anti-affinity | never run the pod on a node whose label value for the key `k8s.io/hostname` matches the value of the label in the node on which a pod with label `app-front-end` is running. i.e. never run these two pods on the same host | 

#### Container Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| command | `[]string` | `command` | The command that runs as the entrypoint to the container | 
| args | `[]floatOrString` | `args` | The arguments to the command. Accepts both float and string values |
| env | `[]Env` | `env` or `envFrom` | The environment variables that get set in the container. See [Environment Overview](#environment-overview) |
| image | `string` | `image` | The Image of the container |
| pull | `string` | `imagePullPolicy` | The image pull policy of the container. It can be "always", "never" or "if-not-present" |
| on_start | `action` | `postStart` | Action to be taken right after container start. See [Actions Overview](#action-overview) |
| pre_stop | `action` | `preStop` | Action to be taken right before container termination. See [Actions Overview](#action-overview) |
| cpu | `CPU` | `resources` | The minimum and the maximum amount of CPUs for this container. More information below | 
| mem | `Mem` | `resources` | The minimum and the maximum amount of memory for this container. More information below |
| cap_add | `[]string` | `capabilites` | The linux capabilities to add to the container | 
| cap_drop | `[]string` | `capabilities` | The linux capabilities to drop from the container |
| privileged | `bool` | `privileged` | Run container in privileged mode | 
| allow_escalation | `bool` | `allowPrivilegeEscalation`| Denotes if processes can gain more privileges than its parents| 
| rw and ro | `bool` | `readOnlyFileSystem` | Mutually inverse flags that denote if the file system is read-only or read-write|
| force_non_root | `bool` | `runAsNonRoot` | Indicates that the container must run as non-root user |
| uid | `int64` | `runAsUser` | Indicates that the container must run as a particular user |
| selinux | `Selinux` | `seLinuxOptions` | SELinux context for the container. More information below |
| liveness_probe | `Probe`| `livenessProbe`| A probe to check if the container is running and alive. See [Probe Overview](#probe-overview)|
| readiness_probe| `Probe` | `readinessProbe` | A probe to check if the container is ready. See [Probe Overview](#probe-overview)|  
| expose | `[]Port` | `Ports` | The set of ports to be exposed by the container. See [Expose Overview](#expose-overview) | 
| stdin | `bool` | `stdin` | Allocate a buffer for stdin |
| stdin_once | `bool` | `stdinOnce` | Close stdin after first attach |
| tty | `bool`| `tty` | Allocate a TTY for container |
| wd | `string` | `workingDir` | Working directory of the container | 
| termination_msg_path | `string` | `terminationMessagePath` | Path where container's termination msg will be read from|
| termination_msg_policy | `string` | `terminationMessagePolicy` | The policy for handling the termination message. See [Termination Message Policy Overview](#termination-message-policy-overview)|
| volume | `[]VolumeMount` | `volumeMounts` | Mount volumes into the container. See [VolumeMounts](#volume-mounts)|

The following fields are status fields and cannot be set

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| container_id | `string` | `status.containerStatus` | The ID of the container as UUID |
| image_id | `string` | `status.imageId` | The ID of the image as UUID |
| ready | `bool` | `status.ready` | Whether the container is ready or not|
| restarts | `int32` | `status.restartCount` | Number of times this container restarted |
| last_state | `ContainerState` | `status.lastTerminationState` | Conditions of the container's last termination. See [Container State](#container-state) | 
| current_state | `ContainerState` | `status.state` | Current condition of the container. See [Container State](#container-state) |

cpu and mem both can contain two fields

| Field | Type | Description         |
|:-----:|:----:|:-------:|
| min | `string` | The minimum amount of resource for the container to run| 
| max | `string` | The maximum amount of resource that the container can use|

cpu is measured in *cpu units*. A cpu unit roughly corresponds to 1 core in any machine. The smallest granule of cpu unit is `one milli core` denoted with an `m` succeeding the cpu unit value. `1m` is a thousanth of 1 core. If no suffix is specified, then it is considered to be whole CPU units. Floating point values are accepted.

mem is measured as an integer or as fixed point integer with units such as `E, P, T, G, M, k` or power of two equivalents `Ei, Pi, Ti, Gi, Mi, and ki`.

Here's an example

```yaml
cpu:
  min: 100m
  max: 0.5 # denotes half of a whole core

mem:
  min: 500M
  max: 1G
```

SELinux options can take these following fields

| Field | Type | Description         |
|:-----:|:----:|:-------:|
| user| `string` | SELinux user label |
| role| `string` | SELinux role label | 
| type | `string` | SELinux type label |
| level | `string` | SELinux level label | 

#### Container State

| Field | Type | Description         |
|:-----:|:----:|:-------:|
| waiting| `ContainerStateWaiting` | Details about a waiting container |
| running| `ContainerStateRunning` | Details about a running container |
| terminated| `ContainerStateTerminated` | Details about a terminated container |


ContainerStateWaiting

| Field | Type | Description         |
|:-----:|:----:|:-------:|
| reason| `string` | Reason for waiting |
| msg | `string` | Message |

ContainerStateRunning

| Field | Type | Description         |
|:-----:|:----:|:-------:|
| start_time| `time` | Time of the last start of the container |

ContainerStateTerminated

| Field | Type | Description         |
|:-----:|:----:|:-------:|
| reason| `string` | Reason for termination |
| msg | `string` | Message |
| exit_code | `int32` | Exit code from the container process |
| signal | `int32` | Signal from the container process | 
| start_time| `time` | Time of the last start of the container |
| finish_time| `time` | Time of the termination of the container |

#### Volume Mounts

| Field | Type | Description         |
|:-----:|:----:|:-------:|
| mount | `string` | Path at which the volume should be mounted |
| store | `string` | Name of the volume to be mounted |
| propagation| `MountPropagation` | Directionality of the mount propagation between host and container (See below.)| 

MountPropagation


| MountPropagationType | K8s counterpart | Description         |
|:-----:|:----:|:-------:|
| host-to-container| HostToContainer| Mounts from host are propagated into container. Not the other way around|
| bidirectional | Bidirectional | Mounts from host are propagated into container and mounts from container are propagated to host|

#### Expose Overview
The expose syntax in Short can be of two types. 

- String
- Struct

If it is a string, then the value is of the format

```yaml
- $protocol://$ip:$host_port:$container_port
```

where `protocol` can take values `TCP` or `UDP`

The expose format in short allows any of the left sub-components to be omitted:

```yaml
# valid expose syntax
- 8080 # expose container port 8080
- 80:8080 # expose host port 80 to container port 8080
- 10.10.1.10:80:8080 # expose host port 80 on nic with ip 10.10.1.10 to container port 8080
- UDP://10.10.1.10:80:8080 # expose host port 80 on 10.10.1.10 to container port 8080 as a UDP port
```
If the struct syntax is used, then the port can be named.

```yaml
- port_name: UDP://10.10.1.10:80:8080 # port is named "port_name"
```

*Note: struct and string syntax can be mixed and matched arbitrarily in the expose array.*

#### Probe Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| action | `Action` | `probe.handler` | An action that determines the state of the container. See [Action Overview](#action-overview) |
| delay | `int32`| `probe.initialDelaySeconds`| Number of seconds to wait before probing initially |
| timeout | `int32` | `probe.timeoutSeconds` | Number of seconds after which the probe times out (default 1)|
| interval | `int32` | `probe.periodSeconds`| Interval of time between two probes (default 10)|
| min_count_success | `int32` | `probe.successThreshold`| Minimum consecutive successful probes to be considered a success (default 1)|
| min_count_failure | `int32` | `probe.failureThreshold` | Minimum consecutive failed probes to be considered a failure (default 1)|

#### Environment Overview

Env variables in Short can be a string or a struct. If it is a struct, then the keys are:

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| from | `string` | `container.EnvFrom` | Obtain environment from k8s resource. Can start with `config:` or `secret:` |
| key | `string` | `container.EnvFrom` | Key of the environment variable or prefix to keys in k8s resource |
| required | `bool` | `container.EnvFrom` | States whether the resource should exist before the creation of container | 

The list of `env` items can mix and match a combination of plain string or struct values. If plain string is used, then it can be one of two formats

 - Key=Value
 - Key


If an env item is a struct, the `from` field is a string with either 2 or 3 `:`-separated sections.
The first section indicates what kind of resource to extract value(s) from: Config Map (`config`) or Secret (`secret`).
The second section is the name of the resource.
The third (optional) section is a specific field to extract from the resource.

If a specific field is specified, its value is used for the env variable in the `key` field.
Otherwise, each field in the named resource is added to the environment, and the `key` field is used as a prefix for each env variable name.

`env` values in Short can be of the following types

```
 # Set key
 - Key

 # Set key = value
 - Key=Value

 # Set environment from config map
 - from:config:$config_map_name
   key: Key  #This is prefixed to every key in the config_map
   required: true # denotes that the config map MUST exist before container creation

 # Set environment from config map key
 - from:config:$config_map_name:$key_in_config_map
   key: Key # This is the name of the env variable inside the container
   required: true

 # Set environment from secret
 - from: secret:$secret_name
   key: Key  #This is prefixed to every key in the secret
   required: true

 # Set environment from secret key
 - from: secret:$secret_name:$key_in_secret
   key: Key  #This is the name of the env variable inside the container
   required: true

```

#### Action Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| command | `[]string` | `container.lifecycle.postStar.exec` | The command to execute as the action |
| net | `NetAction`  | `container.lifecycle.postStart.httpGet` and `container.lifecycle.postStart.tcpSocket` | The network call to make as the action. See [NetAction Overview](#netaction-overview) | 

Here an example action

```yaml
# command action
action:
  command: 
  - path
  - to
  - command
```

*Note: NetAction examples provided in the [NetAction Section](#netaction-overview)*

#### NetAction Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| headers | `[]string` | `container.lifecycle.postStart.httpGet.headers` | The headers that get sent as a part of the network call |
| url | `string` | `container.lifecycle.postStart.httpGet.(path|port|host|scheme)` | The url of the network call |

The URL should be of this form:

`$SCHEME://$HOST:$PORT`

where `$SCHEME` can be `HTTP` (default), `HTTPS` or `TCP`

Here's few examples of net actions

```yaml
# net action
action:
  net:
   url: HTTP://localhost:8080/healthz
   headers:
   - X-CACHE-INVALIDATE:true
   - X-Authhorization:headervalue

# net tcp action
action:
  net:
   url: TCP://localhost:34312
```

#### Termination Message Policy Overview

The Termination Policy to use when handling Container Termination

| Short Termination Policy | K8s counterpart(s) | Description            |
|:----------------:|:------------------:|:----------------------:|
| file | File | Read the container's status message from the file in termination_msg_path |
| fallback-to-logs-on-error | FallbackToLogsOnError | Read the container's status message from logs if file in termination_msg_path is empty |

#### DNS Policy Overview

The DNS Policy supported by Short are the same DNS Policies as Kubernetes. 

| Short DNS Policy | K8s counterpart(s) | Description            |
|:----------------:|:------------------:|:----------------------:|
| cluster-first | ClusterFirst | Pod uses cluster DNS unless HostNetwork is true, then fallback to default DNS |
| cluster-first-with-host-net | ClusterFirstWithHostNet | Pod uses cluster DNS first, then fallback to default DNS |
| default | Default | Pod should use default DNS settings, as set in Kubelet |

#### Host Aliases Overview

Host Aliases are entries to add to the /etc/hosts file.

The /etc/hosts file has multiple lines with data in each line of this format

`127.0.0.1 localhost localhost.1`

i.e the format can be summarized as `$IP $HOST_NAME1 $HOST_NAME2 $HOST_NAME3...`

This is the format followed in Short syntax. The following are valid Host Aliases in Short.

```yaml
 - 127.0.0.1 localhost
 - 10.10.0.10 hp-printer hp-printer.office
 - 10.10.0.11 canon-printer cannon-printer.home
```

#### Host Mode Conversion

`host_mode` expects a list of strings. The list items can take any of the following values.

| Value | K8s counterpart(s) | Description         |
|:-----:|:-----------:|:--------------------------:|
| net | `spec.hostNetwork=true` | Use the host's network namespace for the pod|
| pid | `spec.hostPID=true` | Use the host's PID namespace for the pod|
| ipc | `spec.hostIPC=true` | Use the host's IPC namespace for the pod|
 
#### Account Conversion

Account in Short syntax correspond to the name of the ServiceAccount resource. The Short syntax also indicates if this should be automounted using the suffix `:auto`

Here's are example account values

```yaml
account: default # service account default. Do not automount
account: apiAccess:auto  # service account apiAccess. Automount it.
```

#### Toleration Conversion

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| selector | `string` | `spec.toleration`| A string that selects the taint to tolerate. More information below |
| expiry_after | `int64` | `spec.toleration` | The number of seconds after which the toleration tolerates the taint |

The selector string selects Taints using the following formats

```yaml
- selector: TaintKey=TaintValue:Effect # tolerate the taint if key exists on node, matches value and effect
- selector: TaintKey:Effect  # tolerate the taint if key exists on node and effect matches
- selector: TaintKey  # tolerate the taint if the key exists on node
```

#### Priority

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| class | `string` | `spec.PriorityClassName`| Indicates the Pod's priority. `SYSTEM` is a reserved keyword with the highest priority  |
| value | `int32` | `spec.Priority` | The priority value |

#### Pod Condition

| Field | Type |  Description         |
|:-----:|:----:|:-------:|
| reason| `string` | One word camel case reason for pod's last transition |
| msg | `string` | Human readable message about the pod's last transition |
| status | `ConditionStatus` | String value that represents the status of the condition. Can be "True", "False" or "Unknown" |
| type | `PodConditionType` | String value that represents the type of condition. Can be "PodScheduled", "Ready" or "Initialized" |
| last_probe_time | `time` | Last time the condition was probed |
| last_transition_time | `time` | Last time the condition transitioned |

#### Volume Overview

This section describes the mechanism for specifying volume sources for the pod. In short syntax, Volume sources are always defined as a map. 

Each of the keys of this map is considered to be the local reference(store) of the volume source, and this key should be referenced in the `volume.store` key in the *container* definition. Here's a simple example showing this behavior

```yaml
# Pod with simple volume source
apiVersion: v1
kind: Pod
metadata:
  name: redis
spec:
  containers:
  - name: redis
    image: redis
    volumeMounts:
    - name: redis-storage
      mountPath: /data/redis
  volumes:
  - name: redis-storage
    emptyDir: {}
```

The equivalent short syntax looks like this

```yaml
# short syntax for volume source with string as the value's type
pod:
  containers:
  - image: redis
    name: redis
    volume:
    - mount: /data/redis
      store: redis-storage  # The store is the local reference to the volume
  name: redis
  version: v1
  volumes:
    redis-storage: empty_dir # The key is the name of the volume 
```

In case this information about the volume cannot be encoded in a single string easily, then the entry's value can be a map instead of a string. This is the case for most of the volume sources in Kubernetes. Here is an example depicting this syntax

```yaml
# Pod with complex volume source
apiVersion: v1
kind: Pod
metadata:
  name: test-ebs
spec:
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volumeMounts:
    - mountPath: /test-ebs
      name: test-volume
  volumes:
  - name: test-volume
    # This AWS EBS volume must already exist.
    awsElasticBlockStore:
      volumeID: <volume-id>
      fsType: ext4
```

```yaml
# short syntax for volume source with map as the value's type
pod:
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-ebs
      store: test-volume   # The local reference is used to denote that the volume is used by this container
  name: test-ebs
  version: v1
  volumes:
    test-volume:    # The key is the local reference to the volume 
      fs: ext4      # parameters for the volume...
      vol_id: <volume-id> 
      vol_type: aws_ebs
```

Here's an example pod that has multiple volume sources (in short syntax)

```yaml
pod:
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-ebs
      store: test-volume
    - mount: /test-empty-dir
      store: test-empty-dir-vol
  name: test-ebs
  version: v1
  volumes:
    test-empty-dir-vol: empty_dir    # Each key is a local reference to a volume
    test-volume:
      fs: ext4
      vol_id: <volume-id>
      vol_type: aws_ebs
```

The following volume sources are defined in Short syntax

| Volume Source | Link | 
|:-------------:|:----:|
| Empty Directory | [empty_dir](#empty-directory) |
| AWS Elastic Block Store | [aws_ebs](#aws-elastic-block-store) |
| Azure Disk | [azure_disk](#azure-disk) |
| Azure File | [azure_file](#azure-file) |
| Ceph FS | [cephfs](#ceph-fs) |
| Cinder | [cinder](#cinder) |
| Downward API | [downward_api](#downward-api) |
| Fibre Channel | [fc](#fibre-channel) |
| Flex  | [flex](#flex) |
| Flocker | [flocker](#flocker) |
| GCE Persistent Disk | [gce_pd](#gce-persistent-disk) |
| Git Repository | [git](#git-repository) |
| GlusterFS | [glusterfs](#gluster-fs) |
| Host Path | [host_path](#host-path) |
| ISCSI | [iscsi](#iscsi) |
| NFS | [nfs](#nfs) |
| Persistent Volume Claim | [pvc](#persistent-volume-claim) |
| Photon Persistent Disk | [photon](#photon-persistent-disk) |
| Projected | [projected](#projected) |
| Portworx | [portworx](#portworx) |
| QuoByte | [quobyte](#quobyte) |
| RBD | [rbd](#rbd) |
| ScaleIO | [scaleio](#scaleio) |
| Secret | [secret](#secret) |
| Storage OS | [storageos](#storage-os) |
| VSphere Volume | [vsphere](#vsphere-volume) |

The next section describes the short syntax for each of the volume source types

##### Empty Directory

| Field | Type | K8s counterpart(s) | Description | 
|:-----:|:----:|:------------------:|:-----------:|
|max_size|`string`|`SizeLimit`| The maximum amount of data that can be stored in this volume|
|medium| `string` | `Medium` | The mechanism by which the storage is allocated to the volume|
|vol_type| `string` | - | This should always be set to `empty_dir` for volumes of type `empty_dir` |

Medium can take the following values

| Medium Type | Description | 
|:-----------:|:-----------:|
|      | `<empty value>` Use the node default| 
| memory      | Store the data in memory (tmpfs) |
| huge-pages   | Store the data in huge pages |

Here's an example pod with empty directory volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  volumes:
    test_volume:
      max_size: 500G
      medium: huge-pages
      vol_type: empty_dir
```

##### AWS Elastic Block Store

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| ro | `bool`  | `ReadOnly` | Make the volume read only |
| vol_id | `string` | `VolumeID` |  Unique ID of the persistent disk resource in AWS (Amazon EBS volume) |
| partition | `int32` | `Partition` | The partition in the volume that you want to mount |
| vol_type| `string` | - | This should always be set to `aws_ebs` for volumes of type `aws_ebs` |

Here's an example pod with aws ebs volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      fs: xfs
      ro: true
      vol_id: i4054053
      vol_type: aws_ebs
      partition: 2
```

##### Azure Disk

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| cache | `string` | `CachingMode` | Host Caching mode |
| disk_uri | `string` | `DataDiskURI` | The URI the data disk in the blob storage |
| disk_name | `string` | `DiskName` | The Name of the data disk in the blob storage |
| kind | `string` | `Kind` | The kind of azure disk volume | 
| vol_type| `string` | - | This should always be set to `azure_disk` for volumes of type `azure_disk` |

Azure disks can be of the following Kinds

| Azure disk Kind | Description | 
|:---------------:|:-----------:|
| dedicated | This denotes single blob disks per storage account |
| shared    | This is the default value. This denotes multiple blob disks per storage account |
| managed   | Azure data disk (only in managed availability set) | 

Azure disk cache can take the following values 

| Azure Disk Caching Mode | Description | 
|:-----------------------:|:-----------:|
| none | No caching |
| ro | Cache reads |
| rw | Cache reads and writes |

Here's an example pod with azure disk volume source 

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      cache: rw
      disk_name: azure_disk_name
      disk_uri: disk://uri.azure.disk
      fs: xfs
      kind: dedicated
      vol_type: azure_disk
```

##### Azure File

Azure file is a simple volume source. It doesn't have a map representation. It just uses plain string

The format of the string is

`azure_file:{secret_name}:{share_name}:{ro}`

where, `secret_name` and `share_name` are required fields and `ro` is optional

The description for each of the fields are provided in this table

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| secret_name | `string` | `SecretName` | The name of secret that contains Azure Storage Account Name and Key  |
| share_name | `string` | `ShareName` | Name of the share  |
| ro | `string` | `ReadOnly` | This is an **optional** field and should only be set if this volume is read-only. When it is not set, then the volume defaults to read-write|
| azure_file| `constant` | - | This prefix is used to denote that this is an `azure_file` volume |

**Please note that this is a simple volume source with a string representation only**

Here's an example pod with azure file volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume: azure_file:azure_file_secret_name:azure_file_share_name:ro
```

##### Ceph FS

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| monitors | `[]string` | `Monitors` | A list of Ceph monitors |
| ro | `bool` | `ReadOnly` | Denote that the volume should be read only |
| secret | `string` | `SecretFile`, `SecretRef`  | Authentication secret for the user |
| user | `string` | `User` | Rados user name. Defaults to `admin` | 
| vol_type| `string` | - | This should always be set to `cephfs` for volumes of type `cephfs` |

The `secret` field is used to denote both `SecretFile` and `SecretRef` in the Kubernetes API resource. This is done by using one of the two prefixes for the value of the `secret` field.

| Secret Field Prefix | Description | 
|:-------------------:|:-----------:|
| `ref:${secret_name}`| If the `ref` prefix is used, then it is a reference to a secret name |
| `file:${/path/to/file}` | A path to keyring for the user. Defaults to `/etc/ceph/user.secret` |

Here's an example pod with ceph fs volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      monitors:
      - monitor1
      - monitor2
      path: /path/to/nowhere
      ro: true
      secret: file:/path/to/secret   # can also be ref:secretname
      user: username
      vol_type: cephfs
```

##### Cinder

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| ro | `bool`  | `ReadOnly` | Make the volume read only |
| vol_id | `string` | `VolumeID` | Identifier to reference the volume in Cinder |
| vol_type| `string` | - | This should always be set to `cinder` for volumes of type `cinder` |

Here's an example pod with cinder volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      fs: xfs
      ro: true
      vol_id: cinderID
      vol_type: cinder
```

##### Downward API

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| items | `map[string]items` | `Items` | A list of downward api sources |
| mode | `int32` | `DefaultMode` | Mode bits to use for created files that don't have mode set. Defaults to `0644` |
| vol_type| `string` | - | This should always be set to `downward_api` for volumes of type `downward_api` |

The items map is a special map. The keys to the map are file paths of the files that will be created from the downward_api source. This key directly corresponds to the `Path` key in the `DownwardAPIVolumeFile` resource in the Kubernetes API.

The value to the file path key is a resource with the following fields

| Field | Type | K8s counterpart(s) | Description | 
|:-----:|:----:|:------------------:|:-----------:|
| field | `string` | `FieldRef` | Selects a field of the pod: only annotations, labels, name and namespace are supported |
| resource | `string` | `ResourceFieldRef` | Selects a resource of the container: only resources limits and requests (limits.cpu, limits.memory, requests.cpu and requests.memory) are currently supported |
| mode | `int32` | `Mode` | Mode bits to use on this file. Overrides default mode |

The volume for the resource field should satisfy the format

`${container-name}:{resource}:{requirement}`

where `container-name` is the name of the container, and `resource` is one of `limits.cpu`, `limits.memory`, `request.cpu` and `requests.memory` and `requirement` is the value for the resource.

Here's an example pod with downward api volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      items:
        path/to/file:
          field: metadata.annotation
          mode: "0644"
        path/to/file1:
          mode: "0644"
          resource: container-name:limits.cpu:1m
      mode: "0644"
      vol_type: downward_api
```

##### Fibre Channel

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| lun | `int32` | `Lun` | FC target lun number |
| ro | `bool` | `ReadOnly` | Denote that the volume is read only. Defaults to `false` | 
| wwid | `[]string` | `WWIDs` |  FC volume world wide identifiers  |
| wwn | `[]string` | `TargetWWNs` | FC target worldwide names |
| vol_type| `string` | - | This should always be set to `azure_disk` for volumes of type `azure_disk` |

Here's an example pod with fibre channel volume source 

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      fs: xfs
      lun: 3
      ro: true
      vol_type: fc
      wwid:
      - wwid1
      - wwid2
      wwn:
      - wwn1
      - wwn2
```

##### Flex

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| ro | `bool`  | `ReadOnly` | Make the volume read only |
| driver | `string` | `Driver` |  Name of the flex driver for this volume |
| options | `map[string]string` | `Options` | Additional parameters to send to the flex drive, if any |
| secret | `string` | `SecretRef` | The name of the secret with sensitive information to pass to the flex driver |
| vol_type| `string` | - | This should always be set to `flex` for volumes of type `flex` |

Here's an example pod with flex volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      fs: xfs
      options:
        key: value
        key2: value2
      ro: true
      secret: secret_name
      driver: flex-driver
      vol_type: flex
```

##### Flocker

Flocker is a simple volume source. It doesn't have a map representation. It just uses plain string

The format of the string is

`flocker:{uuid}`

where, `uuid` is the unique identifier of the flocker dataset

and `flocker` is the volume type (mandatory constant)

**Please note that this is a simple volume source with a string representation only**

Here's an example pod with flocker volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume: flocker:uuid
```

##### GCE Persistent Disk

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| ro | `bool`  | `ReadOnly` | Make the volume read only |
| vol_id | `string` | `PDName` |  Unique ID of the persistent disk resource in GCE (GCE PD volume) |
| partition | `int32` | `Partition` | The partition in the volume that you want to mount |
| vol_type| `string` | - | This should always be set to `gce_pd` for volumes of type `gce_pd` |

Here's an example pod with gce persistent disk volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      fs: xfs
      partition: 1
      ro: true
      vol_id: name_of_pd
      vol_type: gce_pd
```

##### GIT Repository 

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| dir | `string`| `Directory`   | Target directory name |
| rev | `string`  | `Revision` | Commit hash for the specfied revision |
| vol_id | `string` | `Repository` |  URL of the git repository |
| vol_type| `string` | - | This should always be set to `git` for volumes of type `git` |

Here's an example pod with git repository volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      dir: testdata
      rev: master
      vol_id: github.com/koki
      vol_type: git
```

##### Gluster FS

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| ro | `bool`| `ReadOnly`   | Denotes that the volume is read-only |
| endpoints | `string`  | `Endpoints` | Endpoint name that details GlusterFS topology |
| path | `string` | `Path` |  GlusterFS volume path |
| vol_type| `string` | - | This should always be set to `glusterfs` for volumes of type `glusterfs` |

Here's an example pod with gluster fs volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      endpoints: endpointsName
      ro: true
      path: /path/to/gluster/vol
      vol_type: glusterfs
```

##### Host Path

Host Path is a simple volume source. It doesn't have a map representation. It just uses plain string

The format of the string is

`host_path:{/path/to/dir}:{host_path_type}`

where, `/path/to/dir` is a required field and `host_path_type` is optional

and `host_path` is the volume type (mandatory constant)

Host Path Type can take the following values

| Host Path Type | Description |
|:--------------:|:-----------:|
|                | `<empty value>` for backwards compatibility. This is also the default value |
| `dir-or-create` | If nothing exists at the given path, an empty directory will be created there as needed with file mode 0755, having the same group and ownership with Kubelet |
| `dir` | A directory must exist at the given path |
| `file-or-create` | If nothing exists at the given path, an empty file will be created there as needed with file mode 0644, having the same group and ownership with Kubelet |
| `file` | A file must exist at the given path |
| `socket` | A UNIX socket must exist at the given path |
| `char-dev` | A character device must exist at the given path |
| `block-dev` | A block device must exist at the given path |

**Please note that this is a simple volume source with a string representation only**

Here's an example pod with host path volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume: host_path:/path/to/dir:block-dev
```

##### ISCSI

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| ro | `bool`  | `ReadOnly` | Make the volume read only |
| portals | `[]string` | `Portals` | iSCSI Target Portal List. The portal is either an IP or ip_addr:port if the port is other than default (typically TCP ports 860 and 3260) |
| chap_discovery | `bool` | `DiscoveryCHAPAuth` | Denotes the support for iSCSI Discovery CHAP authentication |
| chap_session | `bool` | `SessionCHAPAuth` | Denotes the support for iSCSI Session CHAP authentication |
| secret | `string` | `SecretRef` | The name of the secret  for iSCSI target and initiator authentication |
| initiator | `string` | `InitiatorName` | Custom iSCSI Initiator Name. If initiatorName is specified with iscsiInterface simultaneously, new iSCSI interface. `<target portal>:<volume name>` will be created for the connection |
| target_portal | `string` | `TargetPortal` | iSCSI Target portal. The Portal is either an IP or ip_addr:port if the port is other than default (typically TCP ports 860 and 3260) |
| iqn | `string` | `IQN` | Target iSCSI qualified Name |
| lun | `int32` | `Lun` | Target iSCSI Lun number |
| iscsi_interface | `string` | `ISCSIInterface` | iSCSI Interface Name that uses an iSCSI transport. Defaults to 'default' (tcp)|
| vol_type| `string` | - | This should always be set to `iscsi` for volumes of type `iscsi` |

Here's an example pod with iscsi volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      chap_discovery: true
      chap_session: true
      fs: xfs
      initiator: initiator.name
      iqn: iqn.iscsi.hello.world
      iscsi_interface: default
      lun: 45
      portals:
      - 127.0.0.1
      - 192.168.1.132:861
      ro: true
      secret: secretName
      target_portal: 10.10.0.10
      vol_type: iscsi
```

##### NFS

NFS is a simple volume source. It doesn't have a map representation. It just uses plain string

The format of the string is

`nfs:{server_addr}:{/path/to/dir}:{ro}`

where, `server_addr` and `/path/to/dir` are required fields and `ro` is optional

The descriptions for the above fields are provided in the table below

| Fields | Type | K8s counterpart(s) | Description |
|:-----:|:------:|:-----------------:|:-----------:|
| server_addr | `string` | `Server` | Server is the hostname or IP address of the NFS server |
| /path/to/vol | `string` | `Path` | Path that is exported by the NFS server |
| ro | `bool`  | `ReadOnly` | Make the volume read only. Defaults to `false` |
| vol_type| `string` | - | This should always be set to `nfs` for volumes of type `nfs` |

**Please note that this is a simple volume source with a string representation only**

Here's an example pod with nfs volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume: nfs:server_addr:/path/to/data
```

##### Persistent Volume Claim

Persistent Volume Claim (PVC) is a simple volume source. It doesn't have a map representation. It just uses plain string

The format of the string is

`pvc:{pvc_name}:{ro}`

where, `pvc_name` is a required field and `ro` is optional

The descriptions for the above fields are provided in the table below

| Fields | Type | K8s counterpart(s) | Description |
|:-----:|:------:|:-----------------:|:-----------:|
| pvc_name | `string` | `ClaimName` | Name of a PersistentVolumeClaim in the same namespace as the pod using this volume |
| ro | `bool`  | `ReadOnly` | Make the volume read only. Defaults to `false` |
| vol_type| `string` | - | This should always be set to `pvc` for volumes of type `pvc` |

**Please note that this is a simple volume source with a string representation only**

Here's an example pod with persistent volume claim volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume: pvc:pvc_name:ro
```

##### Photon Persistent Disk

Photon Persistent Disk is a simple volume source. It doesn't have a map representation. It just uses plain string

The format of the string is

`photon:{photon_id}:{fs}`

where, `photon_id` is a required field and `fs` is optional

The descriptions for the above fields are provided in the table below

| Fields | Type | K8s counterpart(s) | Description |
|:-----:|:------:|:-----------------:|:-----------:|
| photon_id | `string` | `pdID` | Identifier for Photon Controller Persistent Disk |
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| vol_type| `string` | - | This should always be set to `photon` for volumes of type `photon` |

**Please note that this is a simple volume source with a string representation only**

Here's an example pod with photon persistent disk volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume: photon:pdname123:xfs
```

##### Projected

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| sources | `[]source` | `Sources` | A list of Projected Sources |
| mode | `int32` | `DefaultMode` | Mode bits to use for created files that don't have mode set. Defaults to `0644` |
| vol_type| `string` | - | This should always be set to `projected` for volumes of type `projected` |

Source is a special structure that takes on different forms based on the type of `source`. There are three types of sources

| Source Type | Description | 
|:-----------:|:-----------:|
| Secret Source | secret to project |
| Config Map Source | configMap to project |
| Downward API Source | downwardAPI to project |

*Note that this source type is used for conceptual understanding only, and is not used anywhere in short syntax*

###### Secret Source

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| items | `map[string]item` | `Items` | A list of secret projections |
| secret | `string` | `LocalObjectReference` | Name of the secret to project from |

The items map for secret source is a special map. The keys to the map are the file paths of the files that will be created from the secret source. The key directly corresponds to the Key in the `KeyToPath` resource in Kubernetes API.

The value to the file path key for secret source is a string of the following format.

`{key}:{mode}`

where `key` is a key that should exist in the secret object being referenced, and `mode` is an optional field to override the default mode bits for the projected volume source 

###### Config Map Source

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| items | `map[string]item` | `Items` | A list of secret projections |
| config | `string` | `LocalObjectReference` | Name of the Config Map to project from |

The items map for config map source is a special map. The keys to the map are the file paths of the files that will be created from the config map source. The key directly corresponds to the Key in the `KeyToPath` resource in Kubernetes API.

The value to the file path key for secret source is a string of the following format.

`{key}:{mode}`

where `key` is a key that should exist in the config map object being referenced, and `mode` is an optional field to override the default mode bits for the projected volume source 

###### Downward API Source

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| items | `map[string]item` | `Items` | A list of downward api projections |

The items map for downward api source is a special map, and different from the items map for secret and configmap. The keys to the map are file paths of the files that will be created from the downward_api source. This key directly corresponds to the `Path` key in the `DownwardAPIVolumeFile` resource in the Kubernetes API.

The value to the file path key for downward api source is a resource with the following fields

| Field | Type | K8s counterpart(s) | Description | 
|:-----:|:----:|:------------------:|:-----------:|
| field | `string` | `FieldRef` | Selects a field of the pod: only annotations, labels, name and namespace are supported |
| resource | `string` | `ResourceFieldRef` | Selects a resource of the container: only resources limits and requests (limits.cpu, limits.memory, requests.cpu and requests.memory) are currently supported |
| mode | `int32` | `Mode` | Mode bits to use on this file. Overrides default mode |

The volume for the resource field should satisfy the format

`${container-name}:{resource}:{requirement}`

where `container-name` is the name of the container, and `resource` is one of `limits.cpu`, `limits.memory`, `request.cpu` and `requests.memory` and `requirement` is the value for the resource.

Here's an example pod with projected volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      sources:
      - items:
          path/to/key: key:0644
        secret: secret_name
      - config: config_map_name
        items:
          path/to/key1: key:0644
      - items:
          path/to/file:
            field: metadata.annotation
            mode: "0644"
          path/to/file1:
            resource: container-name:limits.cpu:1m
      vol_type: projected
```

##### Portworx

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| ro | `bool`  | `ReadOnly` | Make the volume read only |
| vol_id | `string` | `VolumeID` | Unique identifier a Portworx volume |  
| vol_type| `string` | - |  This should always be set to `portworx` for volumes of type `portworx` |

Here's an example pod with portworx volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      fs: xfs
      ro: true
      vol_id: 3094ejer
      vol_type: portworx
```

##### QuoByte

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| ro | `bool`  | `ReadOnly` | Make the volume read only |
| vol_id | `string` | `Volume` | Reference to already created QuoByte volume by name |  
| user | `string` | `User` | User to map volume access to. Defaults to serviceaccount user | 
| group| `string` | `Group` | Group to map volume access to. Defaults to no group |
| registry| `string` | `Registry` | `registry` represents a single or multiple Quobyte Registry services specified as a string as host:port pair (multiple entries are separated with commas) which acts as the central registry for volumes |
| vol_type| `string` | - |  This should always be set to `quobyte` for volumes of type `quobyte` |

Here's an example pod with quobyte volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      group: quobyte.group
      registry: registry.quobyte
      ro: true
      user: quobyte.user
      vol_id: volume.quobyte
      vol_type: quobyte
```

##### RBD 

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| ro | `bool`  | `ReadOnly` | Make the volume read only |
| user | `string` | `RadosUser` | Rados user name. Default is `admin` | 
| pool| `string` | `RBDPool` | Rados Pool name. Defaults to `rbd` |
| image| `string` | `RBDImage` | Rados image name |
| monitors | `[]string` | `CephMonitors` | A list of ceph monitors |
| keyring | `string` | `Keyring` | Keyring is the path to key ring for RBDUser. Defaults to `/etc/ceph/keyring` |
| secret | `string` | `SecretRef` |  Name of the authentication secret for RBDUser. If provided overrides keyring |
| vol_type| `string` | - |  This should always be set to `rbd` for volumes of type `rbd` |

Here's an example pod with rbd volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      fs: ext4
      image: foo
      keyring: /etc/ceph/keyring
      monitors:
      - 1.2.3.4:6789
      - 1.2.3.5:6789
      pool: kube
      ro: true
      secret: secret-name
      user: admin
      vol_type: rbd
```

##### ScaleIO

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| ro | `bool`  | `ReadOnly` | Make the volume read only |
| gateway | `string` | `Gateway` | The host address of the ScaleIO API Gateway | 
| system | `string` | `System` | The name of the storage system as configured in ScaleIO |
| secret | `string` | `SecretRef` | Name of the secret for ScaleIO user and other sensitive information. If this is not provided, Login operation will fail |
| ssl | `bool` | `SSLEnabled` | Flag to enable/disable SSL communication with Gateway, default false |
| protection_domain | `string` | `ProtectionDomain` | The name of the ScaleIO Protection Domain for the configured storage |
| storage_pool | `string` | `StoragePool` | The ScaleIO Storage Pool associated with the protection domain |
| storage_mode | `string` | `StorageMode` | Indicates the storage mode for a volume |
| vol_id | `string` | `VolumeName` | The name of a volume already created in the ScaleIO system that is associated with this volume source |
| vol_type| `string` | - |  This should always be set to `scaleio` for volumes of type `scaleio` |

Where, the storage mode can be

| Storage Mode Type |  Description |
|:-----------------:|:------------:|
| thick-provisioned | Disk is allocated on creation |
| thin-provisioned | Space required is allocated on demand |

Here's an example pod with scaleio volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      fs: xfs
      gateway: gateway
      protection_domain: domain
      ro: true
      secret: secret
      ssl: true
      storage_mode: thick-provisioned
      storage_pool: storage.pool
      system: system
      vol_id: volName
      vol_type: scaleio
```

##### Secret 

| Field | Type| K8s counterpart(s) | Description |
|:-----:|:---:|:------------------:|:-----------:|
| required | `bool` | `Optional` | Denotes that the secret should exist if set to `true` |
| mode | `int32` | `DefaultMode` | Mode bits to use for created files that don't have mode set. Defaults to `0644` |
| items | `map[string]item` | `Items` | A list of secret projections |
| vol_id | `string` | `SecretName` | Name of the secret to project from |
| vol_type| `string` | - | This should always be set to `secret` for volumes of type `secret` |


The items map for secret source is a special map. The keys to the map are the file paths of the files that will be created from the secret source. The key directly corresponds to the Key in the `KeyToPath` resource in Kubernetes API.

The value to the file path key for secret source is a string of the following format.

`{key}:{mode}`

where `key` is a key that should exist in the secret object being referenced, and `mode` is an optional field to override the default mode bits for the secret volume source 

Here's an example pod with secret volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      items:
        /path/to/key: key:0644
      mode: "0644"
      required: false
      vol_id: secret_name
      vol_type: secret
```

##### Storage OS

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| ro | `bool`  | `ReadOnly` | Make the volume read only |
| secret | `string` | `SecretRef` | Name of the secret to use for obtaining the StorageOS API credentials.  If not specified, default values will be attempted |
| vol_namespace | `string` | `VolumeNamespace` | Scope of the volume in storageOS. Defaults to pod namespace if unspecified |
| vol_id | `string` | `VolumeName` | The name of a volume in the storageOS system. Unique withing namespace |
| vol_type| `string` | - |  This should always be set to `storageos` for volumes of type `storageos` |

Here's an example pod with storage os volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      fs: xfs
      ro: true
      secret: secret_name
      vol_id: storageOS.VolName
      vol_namespace: storageOS.VolNamespace
      vol_type: storageos
```
##### VSphere Volume

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| fs | `string`| `FSType`   | Filesystem type of the volume that you want to mount |
| policy | `policy`  | `StoragePolicyName`, `StoragePolicyID` | Policy for policy based management |
| vol_id | `string` | `VolumeName` | The path that identifies vsphere volume vmdk |
| vol_type| `string` | - |  This should always be set to `vsphere` for volumes of type `vsphere` |

where the policy struct has the following fields

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
| id | `string`| `StoragePolicyID`   | Storage Policy Based Management (SPBM) profile ID associated with the StoragePolicyName |
| name | `string` | `StoragePolictName` |  Storage Policy Based Management (SPBM) profile name |

Here's an example pod with vsphere volume source

```yaml
pod:
  annotations:
    meta: _test
  cluster: test_cluster
  labels:
    app: meta_test
  name: meta_test
  namespace: test
  version: v1
  containers:
  - image: gcr.io/google_containers/test-webserver
    name: test-container
    volume:
    - mount: /test-volume
      store: test-volume
  volumes:
    test_volume:
      fs: xfs
      policy:
        id: id
        name: policy
      vol_id: /path/to/vsphere/volume
      vol_type: vsphere
```

# Examples 

 - Nginx Pod that exposes container Port 80, and defines a readiness_probe

```yaml
pod:
  containers:
  - expose:
    - 80
    image: nginx
    name: nginx
    readiness_probe:
      delay: 5
      timeout: 5
      net:
        url: HTTP://localhost:80/  #localhost denotes the localhost of the pod
  labels:
    name: nginx
  name: nginx
  version: v1
```

 - Pod that defines environment values
     1. directly (as KEY=VALUE)
     2. from secret
     3. from configmap key
     4. from resource field

```yaml
pod:
  containers:
  - command:
    - /bin/sh
    - -c
    - env
    env:
    - CONFIG_DIR=/config
    - from: secret:test-secret
      key: TEST_CMD_1
    - from: config:test-configmap:key-2
      key: TEST_CMD_2
    - from: metadata.name
      key: TEST_CMD_3
    image: gcr.io/google_containers/busybox
    name: test-container
  name: env-test-pod
  restart_policy: never
  version: v1

```
 - Pod that runs in either failure domain `us-east1` or `us-east2`

```yaml
pod:
  affinity:
  - node: node:k8s.io/failure-domain=us-east1
  - node: node:k8s.io/failure-domain=us-east2
  containers:
  - image: nginx:latest
    name: nginx_container
  labels:
    app: nginx
  name: nginx
  version: v1
```

 - Pod that has multiple containers and defines CPU limit and MEM request

```yaml
pod:
  containers:
  - cpu:
      max: 500m
    mem:
      min: 512m
    env:
    - MASTER=true
    expose:
    - 6379
    image: kubernetes/redis:v1
    name: master
    volume:
    - mount: /redis-master-data
      store: data
  - env:
    - SENTINEL=true
    expose:
    - 26379
    image: kubernetes/redis:v1
    name: sentinel
  labels:
    name: redis
    redis-sentinel: "true"
    role: master
  name: redis-master
  version: v1
  volumes:
  - type: empty-dir
    name: data
```

# Skeleton

| Short Type           | Skeleton                                       |
|:--------------------:|:----------------------------------------------:|
| Pod                  | [skel](../skel/pod.short.skel.yaml)            |

Here's a starter skeleton of a Short Pod.
```yaml
pod:
  name: $name
  labels: 
    app: $name
  namespace: default
  version: v1
  containers:
    image: $image
    name: $container_name
    env:
    - $key=$val
    expose:
    - 80:8080 # hostPort:containerPort
  restart_policy: always
```
