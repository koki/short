# Introduction

A Pod is the unit of execution in Kubernetes. It consists of a set of co-located containers that share the same fate. The Pod definition in Kubernetes includes information about the containers, their runtime characteristics and metadata about the pod.

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

The following sections contain detailed information about the process of converting each of the fields within the Kubernetes Pod definition to Short spec and back.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this Pod is running |
|name | `string` | `metadata.name`| The name of the Pod | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this Pod will be a member of | 
|labels | `string` | `metadata.labels`| Metadata that could be identifying information about the Pod | 
|annotations| `string` | `metadata.annotations`| Non identifying information about the Pod| 
|volumes | [`Volume`](#volume-overview) | `spec.volumes` | Denotes the volumes that are a part of the Pod. More information is available in [Volume Overview](#volume-overview) |
| affinity | [`[]Affinity`](#affinity-overview) | `spec.affinity` and `spec.NodeSelector` | (Anti-) Affinity of the Pod to nodes or other Pods. More information is available in [Affinity Overview](#affinity-overview) |
| containers |[`Container`](../skel/container.short.skel.yaml) | `spec.containers` and `status`| Containers that run as a part of the Pod. More information is available in [Container Overview](#container-overview) |
| init_containers | [`Container`](../skel/container.short.skel.yaml) | `spec.initContainers` and `status` | Containers that run as a part of the initialization process of the Pod. More information is available in [Container Overview](#container-overview) | 
| dns_policy | [`DNSPolicy`](#dns-policy-overview) | `spec.dnsPolicy` | The DNS Policy of the Pod. The conversion function is explained in [DNS Policy Overview](#dns-policy-overview) |
| host_aliases | `[]string` | `spec.aliases` | Set of additional records to be placed in `/etc/hosts` file inside the Pod. More information is available in [Host Aliases Overview](#host-aliases-overview) |
| host_mode | `[]string` | `spec.hostPID`, `spec.hostNetwork` and `spec.hostIPC`| The access the Pod has to host resources. The conversion function is explained in [Host Mode Conversion](#host-mode-conversion) |
| hostname | `string` | `spec.hostname` and `spec.subDomain` | The fully qualified domain name of the pod|
| registry_secrets | `[]string` |`spec.ImagePullSecrets` | A list of k8s secret resource names that contain credentials to required to access private registries. |
| restart_policy | [`RestartPolicy`](#restart-policy) | `spec.restartPolicy` | Behavior of a Pod when it dies. Can be "Always", "OnFailure" or "Never" |
| scheduler_name | `string` | The value from `spec.schedulerName` is stored here | The value from `spec.schedulerName` is stored here |
| account | `string` | `spec.serviceAccountName` and `spec.automountServiceAccountToken` | The access the Pod gets to the K8s API. More information is available in [Account Conversion](#account-conversion) | 
| tolerations | [`[]Toleration`](../skel/toleration.short.skel.yaml) | `spec.tolerations` | Set of tainted hosts to tolerate on scheduling the Pod. The conversion function is explained in [Toleration Conversion](#toleration-conversion) |
| termination_grace_period | `int64`  | `spec.terminationGracePeriodSeconds` | Number of seconds to wait before forcefully killing the Pod. |
| active_deadline | `int64` | `spec.activeDeadlineSeconds`| Number of seconds the Pod is allowed to be active  |  
| node | `string` | `spec.nodeName` | Request Pod to be scheduled on node with the name specified in this field| 
| priority | `Priority` | `spec.priorityClassName` and `spec.priority` | Specifies the Pod's Priority. More information in [Priority](#priority) |
| condition | `[]Pod Condition` | `status.conditions` | The list of current and previous conditions of the Pod. More information in [Pod Condition](#pod-condition) |
| node_ip | `string` | `status.hostIP` | The IP address of the host on which the Pod is running | 
| ip | `string` | `status.podIP` | The IP address of the Pod | 
| start_time | `time` | `status.startTime` | The time at which the Pod started running | 
| msg | `string` | `status.message` | A human readable message explaining Pod's current condition |  
| phase | `string` | `status.phase` | The current phase of the Pod |
| reason | `string` | `status.reason` | Reason indicating the cause for the current state of the Pod |
| qos | `string` | `status.qosClass` | The QOS class assigned to the Pod based on resource requirements |
| fs_gid | `int64` | `spec.securityContext.fsGroup` | Special supplemental group that apply to all the Containers in the Pod |
| gids | `[]int64` | `spec.securityContext.supplementalGroups` | A list of groups applied to the first process in each of the containers of the Pod |

#### Affinity Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| node | `string` | `affinity.nodeAffinity` | The affinity of the Pod to the node. More information available below | 
| pod | `string` | `affinity.podAffinity` | The affinity of a Pod to other Pod(s) in the cluster. More information below |
| anti_pod | `string` | `affinity.podAntiAffinity` | The anti-affinity of a Pod to other Pod(s) in the cluster. More information below |
| topology | `string` | `affinity.pod*.podAffinityTerm.topologyKey` | A key to discern the membership of a Pod to a particular group of nodes in the cluster|
| namespaces | `[]string` | `affinity.pod*.podAffinityTerm.namespaces` | A list of namespaces in which the pod and anti_pod affinities are applied |

`node`, `pod` and `anti_pod` are string fields that expect `expressions` that denote affinities of the Pod to nodes, other pods and anti-affinity to other pods respectively.

`expressions` are label selectors, i.e. node or pods that contain labels that match these expressions are used to make scheduling decisions for the Pod.

An expression is a set of sub-expressions that are ANDed(`&`) together. Sub expressions select on labels using `=`, `!=`, `Exists`, `Does not Exist`, `Greater than` and `Less than` operators.

| Operator | Symbol | Validity | Example |
|:--------:|:------:|:--------:|:--------:|
| Equals   | `=`   | node, pod and anti_pod | `k8s.io/failure-domain=us-east1` |
| Not Equal | `!=`  | node, pod and anti_pod | `k8s.io/failure-domain!=us-east-1`  |
| Exists | N/A  |  node, pod and anti_pod | `k8s.io/cloud-provider` |
| Does Not Exist | N/A | node, pod and anti_pod | `!k8s.io/bare-metal`  |
| Greater Than | '>' | node | `k8s.io/cpus>1` |
| Less Than | '<' | node | `k8s.io/cpus < 1`|

Expressions also have qualifiers at the end of the composite sub-expressions. Qualifiers can be used to set `soft` affinity and (weight) of the soft affinity. 

Pods accept multiple affinity items, and the entire set of affinity items is considered for its scheduling. 

*It is not valid to include more than one of (node, pod, anti_pod) in a single affinity item. They should be specified in separate items. `topology` and `namespaces` are ignored if the affinity item is a `node` selector affinity item*

#### Node Affinity
`node` values can be used to denote `soft` as well as `hard` affinities to nodes. `hard` affinities are expressions that must be satisfied for a Pod to be scheduled on a node.`soft` affinities are expressions that should be satisfied to the best extent possible, but Pods maybe scheduled on nodes that do not completely satisfy these expressions.

In the entire list of `affinity` items, if there are multiple `hard` node affinity selectors, then a node which satisfies any one of the `hard` node affinity selectors is chosen.

In case of `soft` node affinity, multiple affinity items can have `soft` node affinites, and the node which satisfies the sub-set of node `soft` affinity items with the greatest sum of weights is chosen to run the Pod. 

Here are some example node affinity expressions

| Expression | Affinity Type | Description |
|:----------|:-------------:|:-----------:|
|-node:`failure-domain=us-east1&instance-type=t2.large` | `node` hard affinity  | run the pod on a node that is in failure domain `us-east1` and whose instance type is `t2.large` |
|-node:`failure-domain=us-east1&instance-type=t2.large`<br/>-node:`failure-domain=us-east2&instance-type=t2.large` | `node` hard affinity | run the pod on a node that is in failure-domain `us-east1` and the instance type is `t2.large` <br/> or <br/> on a node in `us-east2` and instance type is `t2.large`|
|-node:`failure-domain=us-east1&instance-type=t2.large:soft` | `node` soft affinity | prefer to run the pod on a node that is in failure-domain `us-east1` and whose instance type is `t2.large` |
|-node:`failure-domain=us-east1:soft:10`<br/>-node:`failure-domain=us-east2:soft:20` | `node` soft affinity | run the pod preferrably on a node in failure-domain `us-east2`, less preferrably on a node in `us-east1`, some other node if none of those options are available|

#### Pod Affinity 
`pod` values can be used to denote `soft` as well as `hard` affinities to other pods. `hard` affinities are expressions that must be satisfied for a Pod alongside another Pod.`soft` affinities are expressions that should be satisfied to the best extent possible, but Pods maybe schedules alongside other Pods that do not completely satisfy these expressions. 

In the entire list of `affinity` items, if there are multiple `hard` pod affinity selectors, then the pod is scheduled along another pod which satisfies ALL of the `hard` node affinity selectors. If there is no such pod, then the pod is not scheduled.

In case of `soft` pod affinity, multiple affinity items can have `soft` node affinites, and the pod is scheduled alongside another pod which satisfies the sub set of pod `soft` affinity items with the greatest sum of weights. 

Here are some example pod affinity expressions

| Expression | Affinity Type | Description |
|:----------|:-------------:|:-----------:|
|-pod:`app=front-end` | `pod` hard affinity  | run the pod on alongside another pod which has label `app=front-end` |
|-pod:`app=front-end`<br/>-pod:`name=react` | `pod` hard affinity | run the pod alongside another pod that has labels `app=front-end` and `name=react` |
|-pod:`app=front-end&name=react:soft` | `pod` soft affinity | prefer to run the pod alongside another pod that has labels `app=front-end` and `name=react` |
|-pod:`app=front-end&name=react:soft:10`<br/>-pod:`app=front-end&name=flux:soft:20` | `pod` soft affinity | run the pod preferrably alongside another pod that has labels `app=front-end` and `name=flux`, less preferrably alongside another pod that has labels `app=front-end` and `name=react`, some other node if none of those options are available|
|-pod:`app=front-end`<br/> topology:`k8s.io/failure-domain` |`pod` hard affinity | run the pod on a node whose label value for the key `k8s.io/failure-domain` matches the value of the label in the node on which a pod with label `app-front-end` is running | 

#### Pod Anti Affinity
The syntax and mechanism of pod anti affinity is the same as pod affinity, except that whenever a match of expression occurs, then this pod is not scheduled alongside another pod that matches the expression.

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
| env | `[]Env` | `env` or `envFrom` | The environment variables that get set in the container. More information available in [Environment Overview](#environment-overview) |
| image | `string` | `image` | The Image of the container |
| pull | `string` | `imagePullPolicy` | The image pull policy of the container. It can be "Always", "Never" or "IfNotPresent" |
| on_start | `action` | `postStart` | Action to be taken right after container start. More information is available in [Actions Overview](#action-overview) |
| pre_stop | `action` | `preStop` | Action to be taken right before container termination. More information is available in [Actions Overview](#action-overview) |
| cpu | `CPU` | `resources` | The minimum and the maximum amount of CPUs for this container. More information below | 
| mem | `Mem` | `resources` | The minimum and the maximum amount of memory for this container. More information below |
| cap_add | `[]string` | `capabilites` | The linux capabilities to add to the container | 
| cap_drop | `[]string` | `capabilities` | The linux capabilities to drop from the container |
| privileged | `bool` | `privileged` | Run container in privileged mode | 
| allow_escalation | `bool` | `allowPrivilegeEscalation`| Denotes if processes can gain more privileges than its parents| 
| rw and ro | `bool` | `readOnlyFileSystem` | Mutually inverse flags that denote if the file system is read-only or read-write|
| force_non_root | `bool` | `runAsNonRoot` | Indicates that the container must run as non-root user |
| uid | `int64` | `runAsUser` | Indicates that the container must run as a particular user |
| selinux | `Selinux` | `seLinuxOptions` | SELinux context for the container. More information is available below |
| liveness_probe | `Probe`| `livenessProbe`| A probe to check if the container is running and alive. More information is available in [Probe Overview](#probe-overview)|
| readiness_probe| `Probe` | `readinessProbe` | A probe to check if the container is ready. More information is available in [Probe Overview](#probe-overview)|  
| expose | `[]Port` | `Ports` | The set of ports to be exposed by the container. More information is available in [Expose Overview](#expose-overview) | 
| stdin | `bool` | `stdin` | Allocate a buffer for stdin |
| stdin_once | `bool` | `stdinOnce` | Close stdin after first attach |
| tty | `bool`| `tty` | Allocate a TTY for container |
| wd | `string` | `workingDir` | Working directory of the container | 
| termination_msg_path | `string` | `terminationMessagePath` | Path where container's termination msg will be read from|
| termination_msg_policy | `string` | `termintaionMessagePolicy` | The policy based on which termination message is handled. More information in [Termination Message Poicy Overview](#termination-message-policy-overview)
| volume | `[]VolumeMount` | `volumeMounts` | Mount volumes into the container. More information in [VolumeMounts](#volume-mounts)|

The following fields are status fields and cannot be set

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| container_id | `string` | `status.containerStatus` | The ID of the container as UUID |
| image_id | `string` | `status.imageId` | The ID of the image as UUID |
| ready | `bool` | `status.ready` | States if the container is ready or not|
| restarts | `int32` | `status.restartCount` | Number of times this container restarted |
| last_state | `ContainerState` | `status.lastTerminationState` | Conditions of last termination of container. More information in [Container State](#container-state) | 
| current_state | `ContainerState` | `status.state` | Current condition of the container. More information in [Container State](#container-state) |

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
| propogation| `MountPropagation` | Directionality of the mount propagation between host and container | 

MountPropagation


| MountPropagationType | K8s counterpart | Description         |
|:-----:|:----:|:-------:|
| HostToContainer| HostToContainer| Mounts from host are propagated into container. Not the other way around|
| Bidirectional | Bidirectional | Mounts from host are propagated into container and mounts from container are propagated to host|

#### Expose Overview
The expose syntax in Short can be of two types. 

- String
- Struct

If it is a string, then the value is of the format

```yaml
- $protocol://$ip:$host_port:$container_port
```

where `protocol` can take values `TCP` or `UDP`

The expose format in short allows any of the left sub-components to be omitted. i.e.

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
| action | `Action` | `probe.handler` | An action that determines the state of the container. More information is available in [Action Overview](#action-overview) |
| delay | `int32`| `probe.initialDelaySeconds`| Number of seconds to wait before probing initially |
| timeout | `int32` | `probe.timeoutSeconds` | Number of seconds after which the probe times out (default 1)|
| interval | `int32` | `probe.periodSeconds`| Interval of time between two probes (default 10)|
| min_count_success | `int32` | `probe.successThreshold`| Minimum consecutive successful probes to be considered a success (default 1)|
| min_count_failure | `int32` | `probe.failureThreshold` | Minimum consecutive failed probes to be considered a failure (default 1)|

#### Environment Overview

Env variables in Short can be a string or a struct. If it is a struct, then the keys are

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| from | `string` | `container.EnvFrom` | Obtain environment from k8s resource. Can start with `config:` or `secret:` |
| key | `string` | `container.EnvFrom` | Key of the environment variable or prefix to keys in k8s resource |
| required | `bool` | `container.EnvFrom` | States whether the resource should exist before the creation of container | 

The list of `env` items can mix and match a combination of plain string or struct values. If plain string is used, then it can be one of two formats

 - Key=Value
 - Key

If the env items is a struct, then the `from` key determines if the value is from a Config Map (prefix `config:`) or Secret (prefix `secret:`). This prefix is then followed by the name of the resource. It can additionally have another value added with a `:` delimiter. If this second instance of the delimiter is present, then the following string is considered the Key of the value to be extracted from the Config Map or Secret.

If using Resource:Key mode, then the `key` value in the list item struct denotes the Key of the env variable within the container.

If using simple Resource mode, then the `key` value is applied to each of the items in the Resource as a prefix and then added to the container.

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

 # Set the environment from secret
 - from: secret:$secret_name
   key: Key  #This is prefixed to every key in the secret
   required: true

 # Set the environment from secret key
 - from: secret:$secret_name:$key_in_secret
   key: Key  #This is the name of the env variable inside the container
   required: true

```

#### Action Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
| command | `[]string` | `container.lifecycle.postStar.exec` | The command that gets run as the action |
| net | `NetAction`  | `container.lifecycle.postStart.httpGet` and `container.lifecycle.postStart.tcpSocket` | The network call that gets made as the action. More information is available in [NetAction Overview](#netaction-overview) | 

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

The URL should be of the form

`$SCHEME://$HOST:$PORT`

Where,
 
`$SCHEME` can be `HTTP` (default), `HTTPS` or `TCP`

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
| File | File | Read the container's status message from the file in termination_msg_path |
| FallbackToLogsOnError | FallbackToLogsOnError | Read the container's status message from logs if file in termination_msg_path is empty |

#### DNS Policy Overview

The DNS Policy supported by Short are the same DNS Policies as Kubernetes. 

| Short DNS Policy | K8s counterpart(s) | Description            |
|:----------------:|:------------------:|:----------------------:|
| ClusterFirst | ClusterFirst | Pod uses cluster DNS unless HostNetwork is true, then fallback to default DNS |
| ClusterFirstWithHostNet | ClusterFirstWithHostNet | Pod uses cluster DNS first, then fallback to default DNS |
| Default | Default | Pod should use default DNS settings, as set in Kubelet |

#### Host Aliases Overview

Host Aliases are entires that are added to the /etc/hosts file.

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
  restart_policy: Never
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
  restart_policy: Always
```
