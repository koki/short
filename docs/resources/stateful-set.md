# Introduction

 StatefulSet is a declarative interface to deploy and manage a set of stateful pod which need guarantees about uniqueness and ordering.

| API group | Resource | Kube Skeleton                                   |
|:----------|:---------|:------------------------------------------------|
| apps/v1beta1  | StatefulSet |  [skel](../skel/stateful-set.apps.v1beta1.kube.skel.yaml)         |
| apps/v1beta2  | StatefulSet |  [skel](../skel/stateful-set.apps.v1beta2.kube.skel.yaml)         |

Here's an example Kubernetes StatefulSet along with its headless service:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  ports:
  - port: 80
    name: web
  clusterIP: None
  selector:
    app: nginx
---
apiVersion: apps/v1beta2
kind: StatefulSet
metadata:
  name: web
spec:
  selector:
    matchLabels:
      app: nginx # has to match .spec.template.metadata.labels
  serviceName: "nginx"
  replicas: 3 # by default is 1
  template:
    metadata:
      labels:
        app: nginx # has to match .spec.selector.matchLabels
    spec:
      terminationGracePeriodSeconds: 10
      containers:
      - name: nginx
        image: gcr.io/google_containers/nginx-slim:0.8
        ports:
        - containerPort: 80
          name: web
        volumeMounts:
        - name: www
          mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
  - metadata:
      name: www
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: my-storage-class
      resources:
        requests:
          storage: 1Gi
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:------|:-----|:--------|:-----------------------|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this StatefulSet is running |
|name | `string` | `metadata.name`| The name of the StatefulSet | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this StatefulSet will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the StatefulSet, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the StatefulSet | 
|replicas| `int32` | `replicas`| The number of replicas of the selected Pod |
|replace_on_delete | `bool` | `strategy` | The strategy for performing upgrades. If set to `true`, then the strategy is `OnDelete`. If not, by default the strategy is `RollingUpdate` |
|partition | `int` | `strategy.rollingUpdate` | Ordinal at which the statefulset should be partitioned during upgrade | 
|service  | `string` | `spec.serviceName`  | Name of the service that governs this Statefulset |
|pod_policy| `string` | `spec.podManagementPolicy` | Policy for creating pods under a StatefulSet. See [Pod Management Policy](#pod-management-policy) |
|pvcs | `[]PersistentVolumeClaim` | `spec.VolumeClaimTemplates` | List of claims the pods are allowed to reference. See [Persistent Volume Claim](./persistent-volume-claim.md) |
|max_revs | `int32` | `revisionHistoryLimit` | Number of old replica sets to retain to allow rollback|
|selector | `map[string]string` or `string` | `selector` | An expression (string) or a set of key, value pairs (map) that is used to select a set of pods to manage using the StatefulSet controller. See [Selector Overview](#selector-overview) |
|pod_meta | `TemplateMetadata` | `template` | Metadata of the Pod that is selected by this StatefulSet. See [Template Metadata](#template-metadata)|
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

#### Pod Management Policy 

| Pod Management Policy | Description |
|:----------------------|:------------|
| ordered | Create pods in strictly increasing order on scale up and strictly decreasing order on scale down  |
| parallel | Create all pods in parallel, and update or delete in parallel |

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

 - Cassandra deployed as StatefulSet 

```yaml
stateful_set:
  containers:
  - cap_add:
    - IPC_LOCK
    cpu:
      max: 500m
      min: 500m
    env:
    - MAX_HEAP_SIZE=512M
    - HEAP_NEWSIZE=100M
    - CASSANDRA_SEEDS=cassandra-0.cassandra.default.svc.cluster.local
    - CASSANDRA_CLUSTER_NAME=K8Demo
    - CASSANDRA_DC=DC1-K8Demo
    - CASSANDRA_RACK=Rack1-K8Demo
    - CASSANDRA_AUTO_BOOTSTRAP=false
    - from: status.podIP
      key: POD_IP
    expose:
    - intra-node: 7000
    - tls-intra-node: 7001
    - jmx: 7199
    - cql: 9042
    image: gcr.io/google-samples/cassandra:v12
    mem:
      max: 1Gi
      min: 1Gi
    name: cassandra
    pre_stop:
      command:
      - /bin/sh
      - -c
      - PID=$(pidof java) && kill $PID && while ps -p $PID > /dev/null; do sleep 1;
        done
    pull: always
    readiness_probe:
      command:
      - /bin/bash
      - -c
      - /ready-probe.sh
      delay: 15
      timeout: 5
    volume:
    - mount: /cassandra_data
      store: cassandra-data
  labels:
    app: cassandra
  name: cassandra
  pvcs:
  - access_modes:
    - rw_once
    annotations:
      volume.beta.kubernetes.io/storage-class: fast
    name: cassandra-data
    storage: 1Gi
  replicas: 3
  selector:
    app: cassandra
  service: cassandra
  version: apps/v1beta2
---
storage_class:
  name: fast
  params:
    type: pd-ssd
  provisioner: k8s.io/minikube-hostpath
  version: storage.k8s.io/v1
```

# Skeleton

| Short Type           | Skeleton                                       |
|:---------------------|:-----------------------------------------------|
| StatefulSet           | [skel](../skel/stateful-set.short.skel.yaml)     |

Here's a starter skeleton of a Short StatefulSet.
```yaml
service:
  cluster_ip: None
  labels:
    app: nginx
  name: nginx
  ports:
  - web: 80
  selector:
    app: nginx
  version: v1
---
stateful_set:
  containers:
  - expose:
    - web: 80
    image: gcr.io/google_containers/nginx-slim:0.8
    name: nginx
    volume:
    - mount: /usr/share/nginx/html
      store: www
  name: web
  pvcs:
  - access_modes:
    - rw_once
    name: www
    storage: 1Gi
    storage_class: my-storage-class
  replicas: 3
  selector:
    app: nginx
  service: nginx
  termination_grace_period: 10
  version: apps/v1beta2
```
