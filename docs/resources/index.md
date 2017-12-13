# Introduction

Short defines ops-friendly syntax for Kubernetes manifests. 

This section includes information about

 - Each of the supported resource types in Short
 - Converting from Short syntax to Kubernetes
 - Converting from Kubernetes syntax to Short
 - Skeleton for each of the resource types
 - Examples for each of the resource types

# Short approach

Short follows the same basic principles for Short'ing all Kubernetes types. 

 - Reduce boilerplate (`apiVersion`, `Kind`, `metadata`, `spec`, `status`... keys)
 - Simplify expressions and consolidate relevant information
 - Intuitive and obvious naming instead of long, programming-style keys (eg. `soft` affinity instead of `preferredDuringSchedulingIgnoredDuringExecution`)
 - Group related information
 - **DO NOT LOSE ANY INFORMATION**

We look at each resource and define Short syntax for it based on the principles above.

All Kubernetes resources have the `TypeMeta` and `ObjectMeta` structures embedded in them. The Short syntax pulls the contents of these structures to the top-level Key. 

### Type Meta

A Kubernetes structure looks like this
```yaml
apiVersion: v1
kind: Pod
...
```

where the type of the resource is inferred using the value in `kind` field, along with the apiGroup (from `apiVersion`).

The equivalent Short structure looks like this
```yaml
pod: 
...
```

### Object Meta

`ObjectMeta` in each of the Kubernetes resources is used to define metadata about the object such as `name`, `labels`, `namespace` and `annotations`. These fields are pulled up to the top-level key in Short.

A Kubernetes structure with ObjectMeta looks like this
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: Pod_Name
  labels: 
    Name: Pod_Name
  namespace: default
  annotations:
    Data: Value
...
```

The equivalent Short structure looks like this
```yaml
pod:
  name: Pod_Name
  labels:
    Name: Pod_Name
  namespace: default
  annotations: 
    Data: Value
...
```

The fields within `Spec` and `Status` are Short'ed using similar principles applied to each of their fields. 

# Resources

The following types are currently supported

| Kubernetes API Group | Kubernetes Type   | Short Type   | Skeleton   | Examples |
|:--------------------:|:-----------------:|:------------:|:----------:|:--------:|
| core/v1 | Pod            | [Pod](./pod.md)| [Pod Skeleton](./pod.md#skeleton)  | [Pod Examples](./pod.md#examples) |
| core/v1 | Service        | [Service](./service.md)| [Service Skeleton](./service.md#skeleton) | [Service Examples](./service.md#examples) |
| extensions/v1beta1 | Deployment | [Deployment](./deployment.md) | [Deployment Skeleton](./deployment.md#skeleton) | [Deployment Examples](./deployment.md#examples) |
| apps/v1beta1 | Deployment | [Deployment](./deployment.md) | [Deployment Skeleton](./deployment.md#skeleton) | [Deployment Examples](./deployment.md#examples) |
| apps/v1beta2 | Deployment |  [Deployment](./deployment.md) | [Deployment Skeleton](./deployment.md#skeleton) | [Deployment Examples](./deployment.md#examples) |
| extensions/v1beta1 | Replica Set | [Replica Set](./replica-set.md) | [Replica Set Skeleton](./replica-set.md#skeleton) | [Replica Set Examples](./replica-set.md#examples) |
| apps/v1beta2 | Replica Set | [Replica Set](./replica-set.md) | [Replica Set Skeleton](./replica-set.md#skeleton) | [Replica Set Examples](./replica-set.md#examples) |
| core/v1 | Replication Controller | [Replication Controller](./replication-controller.md) | [Replication Controller Skeleton](./replication-controller.md#skeleton) | [Replication Controller Examples](./replication-controller.md#examples) |
| batch/v1 | Job | [Job](./job.md) | [Job Skeleton](./job.md#skeleton) | [Job Examples](./job.md#examples) |
| extensions/v1beta1 | DaemonSet | [DaemonSet](./daemon-set.md) | [DaemonSet Skeleton](./daemon-set.md#skeleton) | [DaemonSet Examples](./daemon-set.md#examples) |
| apps/v1beta2 | DaemonSet | [DaemonSet](./daemon-set.md) | [DaemonSet Skeleton](./daemon-set.md#skeleton) | [DaemonSet Examples](./daemon-set.md#examples) |
| batch/v2alpha1 | CronJob | [CronJob](./cron-job.md) | [CronJob Skeleton](./cron-job.md#skeleton) | [CronJob Examples](./cron-job.md#examples) |
| batch/v1beta1 | CronJob | [CronJob](./cron-job.md) | [CronJob Skeleton](./cron-job.md#skeleton) | [CronJob Examples](./cron-job.md#examples) |
| apps/v1beta1 | StatefulSet | [StatefulSet](./stateful-set.md) | [StatefulSet Skeleton](./stateful-set.md#skeleton) | [StatefulSet Examples](./stateful-set.md#examples) |
| apps/v1beta2 | StatefulSet |  [StatefulSet](./stateful-set.md) | [StatefulSet Skeleton](./stateful-set.md#skeleton) | [StatefulSet Examples](./stateful-set.md#examples) |
| core/v1 | PersistentVolumeClaim | [PersistentVolumeClaim](./persistent-volume-claim.md) | [PersistentVolumeClaim Skeleton](./persistent-volume-claim.md#skeleton) | [PersistentVolumeClaim Examples](./persistent-volume-claim.md#examples) |
| storage/v1 | StorageClass | [StorageClass](./storage-class.md) | [StorageClass Skeleton](./storage-class.md#skeleton) | [StorageClass Examples](./storage-class.md#examples) |
| core/v1 | Endpoint | [Endpoint](./endpoint.md) | [Endpoint Skeleton](./endpoint.md#skeleton) | [Endpoint Examples](./endpoint.md#examples) |
| extensions/v1beta1 | Ingress | [Ingress](./ingress.md) | [Ingress Skeleton](./ingress.md#skeleton) | [Ingress Examples](./ingress.md#examples) |
