# Introduction

StorageClass describes the parameters for a class of storage for which PersistentVolumes can be dynamically provisioned. 

| API group | Resource | Kube Skeleton                                   |
|:----------|:---------|:------------------------------------------------|
| storage/v1  | StorageClass |  [skel](../skel/storage-class.storage.v1.kube.skel.yaml)         |

Here's an example Kubernetes StorageClass:
```yaml
kind: StorageClass
apiVersion: storage.k8s.io/v1
metadata:
  name: standard
provisioner: kubernetes.io/aws-ebs
parameters:
  type: gp2
reclaimPolicy: Retain
mountOptions:
  - debug
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:------|:-----|:--------|:-----------------------|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this StorageClass is running |
|name | `string` | `metadata.name`| The name of the StorageClass | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this StorageClass will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the StorageClass, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the StorageClass | 
|provisioner| `string` | `provisioner`| Indicates the type of provisioner |
|params | `map[string]string` | `spec.volumeName` | Parameters for the provisioner that should create volumes of this class |
|reclaim | `string` | `reclaimPolicy` | reclaim policy for dynamically provisioned persistent volumes. Defaults to `delete`. See [Reclaim Policy](#reclaim-policy) | 
|mount_opts | `[]string` | `mountOptions` | Mount options for dynamically provisioned persistent volumes|
|allow_expansion | `bool` | `allowVolumeExpansion` | If set, the volumes of this class are expandable |

#### Reclaim Policy

| Recalim Policy | Description |
|:----------------------|:------------|
| recycle | Recycle back volume into unbound pool on release from claim |
| delete | Delete volume on release from claim |
| retain | Leave volume in current phase (Released) for manual reclamation by an admin |

# Examples 

 - StorageClass for AWS EBS volume with IOPS and zone requests

```yaml
storage_class:
  name: slow
  params:
    iopsPerGB: "10"
    type: io1
    zones: us-east-1d, us-east-1c
  provisioner: kubernetes.io/aws-ebs
  version: storage.k8s.io/v1
```

 - GlusterFS StorageClass

```yaml
storage_class:
  name: slow
  params:
    clusterid: 630372ccdc720a92c681fb928f27b53f
    gidMax: "50000"
    gidMin: "40000"
    restauthenabled: "true"
    resturl: http://127.0.0.1:8081
    restuser: admin
    secretName: heketi-secret
    secretNamespace: default
    volumetype: replicate:3
  provisioner: kubernetes.io/glusterfs
  version: storage.k8s.io/v1
```

# Skeleton

| Short Type           | Skeleton                                       |
|:---------------------|:-----------------------------------------------|
| StorageClass           | [skel](../skel/storage-class.short.skel.yaml)     |

Here's a starter skeleton of a Short StorageClass.
```yaml
storage_class:
  name: slow
  params:
    type: pd-standard
    zones: us-central1-a, us-central1-b
  provisioner: kubernetes.io/gce-pd
  version: storage.k8s.io/v1
```
