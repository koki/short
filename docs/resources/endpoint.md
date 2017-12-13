# Introduction

Endpoints is a collection of endpoints that implement the actual service

| API group | Resource | Kube Skeleton                                   |
|:---------:|:--------:|:-----------------------------------------------:|
| core/v1  | Endpoints |  [skel](../skel/endpoints.kube.skel.yaml)         |

Here's an example Kubernetes Endpoints:
```yaml
kind: Endpoints
apiVersion: v1
metadata:
  name: my-service
subsets:
  - addresses:
      - ip: 1.2.3.4
    ports:
      - port: 9376
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:-----:|:----:|:-------:|:----------------------:|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this Endpoints is present |
|name | `string` | `metadata.name`| The name of the Endpoints| 
|namespace | `string` | `metadata.namespace` | The K8s namespace this Endpoints will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the Endpoints, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the Endpoints | 
|subsets| `[]Subsets` | `subsets`| The set of all endpoints is derived by taking the union of all subsets. See [Subsets](#subsets)|

#### Subsets

| Field | Type | K8s counterpart(s) | Description |
|:-----:|:----:|:------------------:|:-----------:|
|addrs| `[]EndpointAddress` | `addresses` | IP addresses which offer the related ports that are marked ready. See [Endpoint Address](#endpoint-address) |
|unready_addrs | `[]EndpointAddress` | `notReadyAddresses` | IP adddresses which offer the related ports that are NOT marked ready. See [Endpoint Addres](#endpoint-address)|
|ports | `[]string` | `ports` | Port numbers available on the related IP addresses. See [Endpoint Ports](#endpoint-ports) |


#### Endpoint Address

| Field    | Type       | K8s counterpart(s) | Description | 
|:--------:|:----------:|:------------------:|:-----------:|
| ip       |`string`    | `ip`               | IP of this endpoint        |
| hostname |`string`    | `hostname`         | hostname of this endpoint  |
| nodename |`string`    | `nodename`         | node hosting this endpoint |
| target   |`ObjectReference` | `targetRef`        | Reference to object providing the endpoint. See [Object Reference](./persistent-volume.md#object-reference) |

#### Endpoint Port

The representation of endpoint port in short syntax is done using a string of the following format

`PROTOCOL://{PORT_NUM}:{NAME}`

where `PROTOCOL` defaults to TCP
and `PORT_NUM` is mandatory and `NAME` is optional

# Examples 

 - Endpoints for a glusterfs cluster

```yaml
endpoints:
  name: glusterfs-cluster
  namespace: spark-cluster
  subsets:
  - addrs:
    - ip: 192.168.30.104
    ports:
    - tcp://1
  - addrs:
    - ip: 192.168.30.105
    ports:
    - tcp://1
  version: v1
```

# Skeleton

| Short Type           | Skeleton                                       |
|:--------------------:|:----------------------------------------------:|
| Endpoints           | [skel](../skel/endpoints.short.skel.yaml)     |

Here's a starter skeleton of a Short Endpoints.
```yaml
endpoints:
  name: glusterfs-cluster
  namespace: spark-cluster
  subsets:
  - addrs:
    - ip: 192.168.30.104
    ports:
    - tcp://1
  - addrs:
    - ip: 192.168.30.105
    ports:
    - tcp://1
  version: v1
```
