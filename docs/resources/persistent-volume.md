# Introduction

 PersistentVolume is a storage resource

| API group | Resource | Kube Skeleton                                   |
|:----------|:---------|:------------------------------------------------|
| core/v1  | PersistentVolume |  [skel](../skel/persistent-volume.kube.skel.yaml)         |

Here's an example Kubernetes PersistentVolume:
```yaml
  apiVersion: v1
  kind: PersistentVolume
  metadata:
    name: pv0003
  spec:
    capacity:
      storage: 5Gi
    accessModes:
      - ReadWriteOnce
    persistentVolumeReclaimPolicy: Recycle
    storageClassName: slow
    mountOptions:
      - hard
      - nfsvers=4.1
    nfs:
      path: /tmp
      server: 172.17.0.2
```

The following sections contain detailed information about each field in Short syntax, including how the field translates to and from Kubernetes syntax.

# API Overview

| Field | Type | K8s counterpart(s) | Description         |
|:------|:-----|:--------|:-----------------------|
|version| `string` | `apiVersion` | The version of the resource object | 
|cluster| `string` | `metadata.clusterName` | The name of the cluster on which this PersistentVolumeClaim is running |
|name | `string` | `metadata.name`| The name of the PersistentVolumeClaim | 
|namespace | `string` | `metadata.namespace` | The K8s namespace this PersistentVolumeClaim will be a member of | 
|labels | `string` | `metadata.labels`| Metadata about the PersistentVolumeClaim, including identifying information | 
|annotations| `string` | `metadata.annotations`| Non-identifying information about the PersistentVolumeClaim | 
|storage_class| `string` | `spec.storageClassName`| The number of storageclass required by the claim |
|access_modes | `[]string` | `spec.accessModes` | Desired access mode the volume should have. See [Access Modes](#access-modes) | 
|storage | `string` | `spec.resources.requests.limit` | Amount of storage the volume should have (eg. 4Gi)|
|reclaim | `string` | `reclaimPolicy` | reclaim policy for dynamically provisioned persistent volumes. Defaults to `delete`. See [Reclaim Policy](./storage-class.md#reclaim-policy) | 
|mount_opts | `[]string` | `mountOptions` | Mount options for dynamically provisioned persistent volumes|
|claim | `ObjectReference` | `spec.claimRef` | Binding reference to persistent volume claim holding this reference |
|vol_type| `string` | - | Reference to the backend volume resource. See [Volume Sources](#volume-sources)|
|... | - | - | Based on the volume type chosen, the appropriate fields for that volume type should be filled into the resource |

#### Volume Sources

| Volume Source | Volume Type | Link | 
|:--------------|:------------|:-----|
| AWS Elastic Block Store | aws_ebs | [aws_ebs](pod#aws-elastic-block-store) |
| Azure Disk | azure_disk | [azure_disk](pod#azure-disk) |
| Azure File | azure_file | [azure_file](pod#azure-file) |
| Ceph FS | cephfs| [cephfs](pod#ceph-fs) |
| Cinder | cinder | [cinder](pod#cinder) |
| Fibre Channel | fc | [fc](pod#fibre-channel) |
| Flex  | flex | [flex](pod#flex) |
| Flocker | flocker | [flocker](pod#flocker) |
| GCE Persistent Disk | gce_pd | [gce_pd](pod#gce-persistent-disk) |
| GlusterFS | glusterfs | [glusterfs](pod#gluster-fs) |
| Host Path | host_path | [host_path](pod#host-path) |
| ISCSI | iscsi | [iscsi](pod#iscsi) |
| NFS | nfs | [nfs](pod#nfs) |
| Photon Persistent Disk | photon | [photon](pod#photon-persistent-disk) |
| Portworx | portworx | [portworx](pod#portworx) |
| QuoByte | quobyte | [quobyte](pod#quobyte) |
| RBD | rbd | [rbd](pod#rbd) |
| ScaleIO | scaleio | [scaleio](pod#scaleio) |
| Storage OS | storageos | [storageos](pod#storage-os) |
| VSphere Volume | vsphere | [vsphere](pod#vsphere-volume) |

The next section describes the short syntax for each of the volume source types

#### Access Modes 

| Access Mode | Description |
|:----------------------|:------------|
| rw_once | Can be mounted read/write mode to exactly 1 host |
| ro_many | Can be mounted read only mode to many hosts |
| rw_many | Can be mounted read/write mode to many hosts |

#### Object Reference

| Field            | Type   | K8s counterpart(s) |
|:-----------------|:-------|:-------------------|
| kind             |`string`| `kind`             |
| namespace        |`string`| `namespace`        | 
| name             |`string`| `name`             |
| uid              |`string`| `uid`              |
| version          |`string`| `version`          |
| resource_version |`string`| `resourceVersion`  |
| field_path       |`string`| `fieldPath`        |

# Examples 

 - PersistentVolume representing nfs storage

```yaml
persistent_volume:
  modes: rw-once
  mount_opts: hard,nfsvers=4.1
  name: pv0003
  reclaim: recycle
  storage: 5Gi
  storage_class: slow
  version: v1
  vol_id: 172.17.0.2:/tmp
  vol_type: nfs
```

 - PersistentVolume representing AWS EBS volume resource

```yaml
persistent_volume:
  fs: ext4
  labels:
    type: amazonEBS
  modes: rw-once
  name: couchbase-pv
  storage: 5Gi
  version: v1
  vol_id: vol-47f59cce
  vol_type: aws_ebs
```

# Skeleton

| Short Type           | Skeleton                                       |
|:---------------------|:-----------------------------------------------|
| PersistentVolume           | [skel](../skel/persistent-volume.short.skel.yaml)     |

Here's a starter skeleton of a Short PersistentVolume.
```yaml
persistent_volume:
  fs: ext4
  modes: rw-once
  name: pv0001
  storage: 5Gi
  version: v1
  vol_id: pd-disk-1
  vol_type: gce_pd
```
