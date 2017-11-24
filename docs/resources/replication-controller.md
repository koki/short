# Introduction

 ReplicationController ensures that a specified number of Pod Replicas are running at any one given time.

| API group | Resource | Kube Skeleton                                   |
|:---------:|:--------:|:-----------------------------------------------:|
| core/v1 | ReplicationController |  [skel](../skel/replication-controller.kube.skel.yaml)         |

Here's an example Kubernetes ReplicationController:
```yaml
apiVersion: v1
kind: ReplicationController
metadata:
  name: nginx
spec:
  replicas: 3
  selector:
    app: nginx
  template:
    metadata:
      name: nginx
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
```

The following sections contain detailed information about the process of converting each of the fields within the Kubernetes ReplicationController definition to Short spec and back.

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
|min_ready | `int32` | `minReadySeconds` | Minimum number of seconds that your pod should be ready before it is considered available |
|selector | `map[string]string` | `selector` | A set of key, value pairs (map) that is used to select a set of pods to manage using the replication controller. ReplicationController cannot use expressions like ReplicaSet, and that is the only differnce between them |
|pod_meta | `TemplateMetadata` | `template` | Metadata of the pod that is selected by this replication controller. More details in [Template Metadata](./deployment.md#template-metadata)|
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

# Examples 

 -  An example replication_controller with 1 replica selecting app:nginx

```yaml
replication_controller:
  containers:
  - expose:
    - 80
    image: nginx:1.7.9
    name: nginx
  name: nginx-rc
  replicas: 1
  selector:
    app: nginx
  version: v1
```

 - An example replication_controller with 3 replicas that selects on labels app=nginx

```yaml
replication_controller:
  containers:
  - expose:
    - 80
    image: nginx:1.7.9
    name: nginx
  labels:
    app: nginx
  name: nginx-rc
  replicas: 3
  selector: 
    app: nginx
  version: v1
```

# Skeleton

| Short Type           | Skeleton                                       |
|:--------------------:|:----------------------------------------------:|
| ReplicationController           | [skel](../skel/replication-controller.short.skel.yaml)     |

Here's a starter skeleton of a ReplicationController
```yaml
replication_controller:
  containers:
  - cpu:
      min: 100m
    env:
    - GET_HOSTS_FROM=dns
    expose:
    - 80
    image: gcr.io/google_samples/gb-frontend:v3
    mem:
      min: 100Mi
    name: php-redis
  labels:
    app: guestbook
    tier: frontend
  name: frontend
  pod_meta:
    labels:
      app: guestbook
      tier: frontend
  replicas: 3
  selector: 
    tier: frontend
  version: v1
```
