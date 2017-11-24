# Introduction

Short defines ops-friendly resource types that correspond to the Kubernetes types. 

This section includes information about

 - Each of the supported types in Short
 - Type conversion from Short to Kubernetes
 - Type conversion from Kubernetes to Short
 - Skeleton for each of the types
 - Examples for each of the types

# Short approach

Short follows the same basic principles for Short'ing all Kubernetes types. 

 - Remove the need to specify boiler plate (`apiVersion`, `Kind`, `metadata`, `spec`, `status`... keys)
 - Simplify expressions and consolidate relevant information
 - Intuitive and obvious names instead of long winded code style Key names (eg. `soft` affinity instead of `preferredDuringSchedulingIgnoredDuringExecution`)
 - Group Information that are related to one another
 - **DO NOT LOSE ANY INFORMATION**

We look at each resource and define the short syntax for it based on the above principles.

All Kubernetes resource have the `TypeMeta` and `ObjectMeta` structures embedded in them. The Short syntax pulls the contents of these structures to the top level Key. 

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

`ObjectMeta` in each of the Kubernetes resources is used to define metadata about the object such as `name`, `labels`, `namespace` and `annotations`. These fields are pulled up to the top level key in Short.

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
| core/v1 | PersistentVolume | [Volume](./volume.md) | [Volume Skeleton](./volume.md#skeleton) | [Volume Examples](./volume.md#examples) |
