package types

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/ghodss/yaml"
	"github.com/kr/pretty"

	"github.com/koki/short/util"
)

var persistentClaim0 = v1.VolumeSource{
	PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
		ClaimName: "claim-0",
		ReadOnly:  true,
	},
}

var hostPathType0 = v1.HostPathDirectoryOrCreate
var hostPath0 = v1.VolumeSource{
	HostPath: &v1.HostPathVolumeSource{
		Path: "/path",
		Type: &hostPathType0,
	},
}
var persistentHostPath0 = v1.PersistentVolumeSource{
	HostPath: hostPath0.HostPath,
}

var sizeLimit0 = resource.MustParse("10Gi")
var emptyDir0 = v1.VolumeSource{
	EmptyDir: &v1.EmptyDirVolumeSource{
		Medium:    v1.StorageMediumMemory,
		SizeLimit: &sizeLimit0,
	},
}

var glusterfs0 = v1.VolumeSource{
	Glusterfs: &v1.GlusterfsVolumeSource{
		EndpointsName: "glusterfs-cluster",
		Path:          "kube_vol",
		ReadOnly:      true,
	},
}
var persistentGlusterfs0 = v1.PersistentVolumeSource{
	Glusterfs: glusterfs0.Glusterfs,
}

var secretRef0 = v1.LocalObjectReference{
	Name: "secretName",
}
var rbd0 = v1.VolumeSource{
	RBD: &v1.RBDVolumeSource{
		CephMonitors: []string{
			"1.2.3.4:6789",
			"1.2.3.5:6789",
		},
		RBDImage:  "foo",
		FSType:    "ext4",
		RBDPool:   "kube",
		RadosUser: "admin",
		Keyring:   "/etc/ceph/keyring",
		SecretRef: &secretRef0,
		ReadOnly:  true,
	},
}
var persistentRBD0 = v1.PersistentVolumeSource{
	RBD: &v1.RBDPersistentVolumeSource{
		CephMonitors: []string{
			"1.2.3.4:6789",
			"1.2.3.5:6789",
		},
		RBDImage:  "foo",
		FSType:    "ext4",
		RBDPool:   "kube",
		RadosUser: "admin",
		Keyring:   "/etc/ceph/keyring",
		SecretRef: &secretRef1,
		ReadOnly:  true,
	},
}

var cinder0 = v1.VolumeSource{
	Cinder: &v1.CinderVolumeSource{
		VolumeID: "bd82f7e2-wece-4c01-a505-4acf60b07f4a",
		FSType:   "ext4",
		ReadOnly: true,
	},
}
var persistentCinder0 = v1.PersistentVolumeSource{
	Cinder: cinder0.Cinder,
}

var cephfs0 = v1.VolumeSource{
	CephFS: &v1.CephFSVolumeSource{
		Monitors: []string{
			"1.2.3.4:6789",
			"1.2.3.5:6789",
		},
		Path:       "/path",
		User:       "admin",
		SecretFile: "/etc/ceph/admin.secret",
		SecretRef:  &secretRef0,
		ReadOnly:   true,
	},
}

var secretRef1 = v1.SecretReference{
	Name:      "secretName",
	Namespace: "secretNamespace",
}
var persistentCephfs0 = v1.PersistentVolumeSource{
	CephFS: &v1.CephFSPersistentVolumeSource{
		Monitors: []string{
			"1.2.3.4:6789",
			"1.2.3.5:6789",
		},
		Path:       "/path",
		User:       "admin",
		SecretFile: "/etc/ceph/admin.secret",
		SecretRef:  &secretRef1,
		ReadOnly:   true,
	},
}

var flocker0 = v1.VolumeSource{
	Flocker: &v1.FlockerVolumeSource{
		DatasetName: "flocker-vol",
		DatasetUUID: "d2d35de3-25c2-45fb-93a3-d47bb5861d46",
	},
}
var persistentFlocker0 = v1.PersistentVolumeSource{
	Flocker: flocker0.Flocker,
}

var gcePersistentDisk0 = v1.VolumeSource{
	GCEPersistentDisk: &v1.GCEPersistentDiskVolumeSource{
		PDName:    "data-disk",
		FSType:    "ext4",
		Partition: 1,
		ReadOnly:  true,
	},
}
var persistentGCEPersistentDisk = v1.PersistentVolumeSource{
	GCEPersistentDisk: gcePersistentDisk0.GCEPersistentDisk,
}

var quobyte0 = v1.VolumeSource{
	Quobyte: &v1.QuobyteVolumeSource{
		Registry: "registry:6789",
		Volume:   "testVolume",
		ReadOnly: true,
		User:     "root",
		Group:    "root",
	},
}
var persistentQuobyte0 = v1.PersistentVolumeSource{
	Quobyte: quobyte0.Quobyte,
}

var flex0 = v1.VolumeSource{
	FlexVolume: &v1.FlexVolumeSource{
		Driver:    "kubernetes.io/lvm",
		FSType:    "ext4",
		SecretRef: &secretRef0,
		ReadOnly:  true,
		Options: map[string]string{
			"volumeID":    "vol1",
			"size":        "1000m",
			"volumegroup": "kube_vg",
		},
	},
}
var persistentFlex0 = v1.PersistentVolumeSource{
	FlexVolume: flex0.FlexVolume,
}

var awsebs0 = v1.VolumeSource{
	AWSElasticBlockStore: &v1.AWSElasticBlockStoreVolumeSource{
		VolumeID:  "volume-id",
		FSType:    "ext4",
		Partition: 1,
		ReadOnly:  true,
	},
}
var persistentAWSEBS0 = v1.PersistentVolumeSource{
	AWSElasticBlockStore: awsebs0.AWSElasticBlockStore,
}

var gitrepo0 = v1.VolumeSource{
	GitRepo: &v1.GitRepoVolumeSource{
		Repository: "git@github.com:koki/short.git",
		Revision:   "425cf6991e957446c2bd09db9ef7baf154d19b23",
		Directory:  "./types",
	},
}

var secretItems0 = []v1.KeyToPath{
	v1.KeyToPath{
		Key:  "username",
		Path: "my-group/my-username",
	},
}
var secret0 = v1.VolumeSource{
	Secret: &v1.SecretVolumeSource{
		SecretName:  "secretName",
		Items:       secretItems0,
		DefaultMode: util.Int32Ptr(0777),
		Optional:    util.BoolPtr(true),
	},
}

var nfs0 = v1.VolumeSource{
	NFS: &v1.NFSVolumeSource{
		Server:   "server-hostname",
		Path:     "/path",
		ReadOnly: true,
	},
}
var persistentNFS0 = v1.PersistentVolumeSource{
	NFS: nfs0.NFS,
}

var iscsi0 = v1.VolumeSource{
	ISCSI: &v1.ISCSIVolumeSource{
		TargetPortal:   "1.2.3.4:3260",
		IQN:            "iqn.2015-02.example.com:test",
		Lun:            0,
		ISCSIInterface: "default",
		FSType:         "ext4",
		ReadOnly:       true,
		Portals: []string{
			"1.2.3.5:3260",
			"1.2.3.6:3260",
		},
		DiscoveryCHAPAuth: true,
		SessionCHAPAuth:   true,
		SecretRef:         &secretRef0,
		InitiatorName:     util.StringPtr("iqn.1996-04.de.suse:linux-host1"),
	},
}
var persistentISCSI0 = v1.PersistentVolumeSource{
	ISCSI: iscsi0.ISCSI,
}

var fc0 = v1.VolumeSource{
	FC: &v1.FCVolumeSource{
		TargetWWNs: []string{
			"500a0982991b8dc5",
			"500a0982891b8dc5",
		},
		Lun:      util.Int32Ptr(2),
		FSType:   "ext4",
		ReadOnly: true,
		WWIDs: []string{
			"actually, WWIDs should not be specified here",
			"either wwids or (targetwwns + lun), not both",
		},
	},
}
var persistentFC0 = v1.PersistentVolumeSource{
	FC: fc0.FC,
}

var azurefile0 = v1.VolumeSource{
	AzureFile: &v1.AzureFileVolumeSource{
		SecretName: "azure-secret",
		ShareName:  "k8stest",
		ReadOnly:   true,
	},
}
var persistentAzurefile0 = v1.PersistentVolumeSource{
	AzureFile: &v1.AzureFilePersistentVolumeSource{
		SecretName:      "azure-secret",
		ShareName:       "k8stest",
		ReadOnly:        true,
		SecretNamespace: util.StringPtr("secret-namespace"),
	},
}

var vspherevdisk0 = v1.VolumeSource{
	VsphereVolume: &v1.VsphereVirtualDiskVolumeSource{
		VolumePath:        "[datastore1] volumes/myDisk",
		FSType:            "ext4",
		StoragePolicyName: "policy-name",
		StoragePolicyID:   "policy-id",
	},
}
var persistentVspherevdisk0 = v1.PersistentVolumeSource{
	VsphereVolume: vspherevdisk0.VsphereVolume,
}

var photonPersistentDisk0 = v1.VolumeSource{
	PhotonPersistentDisk: &v1.PhotonPersistentDiskVolumeSource{
		PdID:   "some-pdid",
		FSType: "ext4",
	},
}
var persistentPhotonPersistentDisk0 = v1.PersistentVolumeSource{
	PhotonPersistentDisk: photonPersistentDisk0.PhotonPersistentDisk,
}

var azureCachingMode0 = v1.AzureDataDiskCachingReadWrite
var azureKind0 = v1.AzureSharedBlobDisk
var azuredisk0 = v1.VolumeSource{
	AzureDisk: &v1.AzureDiskVolumeSource{
		DiskName:    "test.vhd",
		DataDiskURI: "https://someaccount.blob.microsoft.net/vhds/test.vhd",
		CachingMode: &azureCachingMode0,
		FSType:      util.StringPtr("ext4"),
		ReadOnly:    util.BoolPtr(true),
		Kind:        &azureKind0,
	},
}
var persistentAzuredisk0 = v1.PersistentVolumeSource{
	AzureDisk: azuredisk0.AzureDisk,
}

var portworx0 = v1.VolumeSource{
	PortworxVolume: &v1.PortworxVolumeSource{
		VolumeID: "volume-id",
		FSType:   "ext4",
		ReadOnly: true,
	},
}
var persistentPortworx0 = v1.PersistentVolumeSource{
	PortworxVolume: portworx0.PortworxVolume,
}

var scaleio0 = v1.VolumeSource{
	ScaleIO: &v1.ScaleIOVolumeSource{
		Gateway:          "https://localhost:443/api",
		System:           "scaleio",
		SecretRef:        &secretRef0,
		SSLEnabled:       true,
		ProtectionDomain: "pd01",
		StoragePool:      "sp01",
		StorageMode:      "ThickProvisioned",
		VolumeName:       "vol-0",
		FSType:           "xfs",
		ReadOnly:         true,
	},
}
var persistentScalio0 = v1.PersistentVolumeSource{
	ScaleIO: &v1.ScaleIOPersistentVolumeSource{
		Gateway:          "https://localhost:443/api",
		System:           "scaleio",
		SecretRef:        &secretRef1,
		SSLEnabled:       true,
		ProtectionDomain: "pd01",
		StoragePool:      "sp01",
		StorageMode:      "ThickProvisioned",
		VolumeName:       "vol-0",
		FSType:           "xfs",
		ReadOnly:         true,
	},
}

var storageos0 = v1.VolumeSource{
	StorageOS: &v1.StorageOSVolumeSource{
		VolumeName:      "vol-0",
		VolumeNamespace: "namespace-0",
		FSType:          "ext4",
		ReadOnly:        true,
		SecretRef:       &secretRef0,
	},
}

var secretRef2 = v1.ObjectReference{
	Namespace: "secret-namespace",
	Name:      "secret-name",
}
var persistentStorageos0 = v1.PersistentVolumeSource{
	StorageOS: &v1.StorageOSPersistentVolumeSource{
		VolumeName:      "vol-0",
		VolumeNamespace: "namespace-0",
		FSType:          "ext4",
		ReadOnly:        true,
		SecretRef:       &secretRef2,
	},
}

var configmapRef0 = v1.LocalObjectReference{
	Name: "cm-name",
}
var configmapItems0 = []v1.KeyToPath{
	v1.KeyToPath{
		Key:  "config",
		Path: "my-group/my-config",
	},
}
var configmap0 = v1.VolumeSource{
	ConfigMap: &v1.ConfigMapVolumeSource{
		LocalObjectReference: configmapRef0,
		Items:                configmapItems0,
		DefaultMode:          util.Int32Ptr(0777),
		Optional:             util.BoolPtr(true),
	},
}

var secretProjection0 = v1.SecretProjection{
	LocalObjectReference: secretRef0,
	Items:                secretItems0,
	Optional:             util.BoolPtr(true),
}

var downwardAPIItems0 = []v1.DownwardAPIVolumeFile{
	v1.DownwardAPIVolumeFile{
		Path: "labels",
		FieldRef: &v1.ObjectFieldSelector{
			FieldPath: "metadata.labels",
		},
	},
	v1.DownwardAPIVolumeFile{
		Path: "cpu_limit",
		ResourceFieldRef: &v1.ResourceFieldSelector{
			ContainerName: "client-container",
			Resource:      "limits.cpu",
			Divisor:       resource.MustParse("1m"),
		},
	},
}
var downwardAPIProjection0 = v1.DownwardAPIProjection{
	Items: downwardAPIItems0,
}

var configmapProjection0 = v1.ConfigMapProjection{
	LocalObjectReference: configmapRef0,
	Items:                configmapItems0,
	Optional:             util.BoolPtr(true),
}

var projected0 = v1.VolumeSource{
	Projected: &v1.ProjectedVolumeSource{
		Sources: []v1.VolumeProjection{
			v1.VolumeProjection{
				Secret: &secretProjection0,
			},
			v1.VolumeProjection{
				ConfigMap: &configmapProjection0,
			},
		},
	},
}

var local0 = v1.PersistentVolumeSource{
	Local: &v1.LocalVolumeSource{
		Path: "/path",
	},
}

var downwardAPI0 = v1.VolumeSource{
	DownwardAPI: &v1.DownwardAPIVolumeSource{
		Items:       downwardAPIItems0,
		DefaultMode: util.Int32Ptr(0777),
	},
}

func TestVolume(t *testing.T) {
	testVolumeSource(persistentClaim0, t)
	testVolumeSource(hostPath0, t)
	testVolumeSource(emptyDir0, t)
	testVolumeSource(glusterfs0, t)
	testVolumeSource(rbd0, t)
	testVolumeSource(cinder0, t)
	testVolumeSource(cephfs0, t)
	testVolumeSource(flocker0, t)
	testVolumeSource(gcePersistentDisk0, t)
	testVolumeSource(quobyte0, t)
	testVolumeSource(flex0, t)
	testVolumeSource(awsebs0, t)
	testVolumeSource(gitrepo0, t)
	testVolumeSource(secret0, t)
	testVolumeSource(nfs0, t)
	testVolumeSource(iscsi0, t)
	testVolumeSource(fc0, t)
	testVolumeSource(azurefile0, t)
	testVolumeSource(vspherevdisk0, t)
	testVolumeSource(photonPersistentDisk0, t)
	testVolumeSource(azuredisk0, t)
	testVolumeSource(portworx0, t)
	testVolumeSource(scaleio0, t)
	testVolumeSource(storageos0, t)
	testVolumeSource(configmap0, t)
	testVolumeSource(projected0, t)
	testVolumeSource(downwardAPI0, t)

	testPersistentVolumeSource(local0, t)
	testPersistentVolumeSource(persistentAWSEBS0, t)
	testPersistentVolumeSource(persistentAzuredisk0, t)
	testPersistentVolumeSource(persistentAzurefile0, t)
	testPersistentVolumeSource(persistentCephfs0, t)
	testPersistentVolumeSource(persistentCinder0, t)
	testPersistentVolumeSource(persistentFC0, t)
	testPersistentVolumeSource(persistentFlex0, t)
	testPersistentVolumeSource(persistentFlocker0, t)
	testPersistentVolumeSource(persistentGCEPersistentDisk, t)
}

func testVolumeSource(v v1.VolumeSource, t *testing.T) {
	kokiVolume := Volume{
		VolumeMeta: VolumeMeta{
			Name: "vol-name",
		},
		VolumeSource: VolumeSource{
			VolumeSource: v,
		},
	}

	b, err := yaml.Marshal(kokiVolume)
	if err != nil {
		t.Error(pretty.Sprintf("%s\n%# v", err.Error(), kokiVolume))
		return
	}

	newVolume := Volume{}

	err = yaml.Unmarshal(b, &newVolume)
	if err != nil {
		t.Error(pretty.Sprintf("%s\n(%s)\n(%# v)", err.Error(), string(b), kokiVolume))
		return
	}

	newB, err := yaml.Marshal(newVolume)
	if err != nil {
		t.Error(pretty.Sprintf("%s\n(%# v)\n(%# v)\n(%s)", err.Error(), newVolume, kokiVolume, string(b)))
		return
	}

	if !reflect.DeepEqual(kokiVolume, newVolume) {
		t.Error(pretty.Sprintf("failed round-trip\n(%# v)\n(%# v)\n(%s)\n(%s)", kokiVolume, newVolume, string(b), string(newB)))
		return
	}
}

func testPersistentVolumeSource(v v1.PersistentVolumeSource, t *testing.T) {
	kokiVolume := PersistentVolume{
		PersistentVolumeMeta: PersistentVolumeMeta{
			Version:   "v1",
			Cluster:   "cluster",
			Name:      "vol-name",
			Namespace: "namespace",
			Labels: map[string]string{
				"labelKey": "labelValue",
			},
			Annotations: map[string]string{
				"annotationKey": "annotationValue",
			},
			Storage: &sizeLimit0,
			AccessModes: &AccessModes{
				Modes: []v1.PersistentVolumeAccessMode{
					v1.ReadWriteOnce,
				},
			},
			Claim: &v1.ObjectReference{
				Name:      "claimName",
				Namespace: "claimNamespace",
			},
			ReclaimPolicy: v1.PersistentVolumeReclaimRecycle,
			StorageClass:  "storageClass",
			MountOptions:  "option 1,option 2,option 3",
			Status:        &v1.PersistentVolumeStatus{},
		},
		PersistentVolumeSource: PersistentVolumeSource{
			VolumeSource: v,
		},
	}

	b, err := yaml.Marshal(kokiVolume)
	if err != nil {
		t.Error(pretty.Sprintf("%s\n%# v", err.Error(), kokiVolume))
		return
	}

	newVolume := PersistentVolume{}

	err = yaml.Unmarshal(b, &newVolume)
	if err != nil {
		t.Error(pretty.Sprintf("%s\n(%s)\n(%# v)", err.Error(), string(b), kokiVolume))
		return
	}

	newB, err := yaml.Marshal(newVolume)
	if err != nil {
		t.Error(pretty.Sprintf("%s\n(%# v)\n(%# v)\n(%s)", err.Error(), newVolume, kokiVolume, string(b)))
		return
	}

	if !reflect.DeepEqual(kokiVolume, newVolume) {
		t.Error(pretty.Sprintf("failed round-trip\n(%# v)\n(%# v)\n(%s)\n(%s)", kokiVolume, newVolume, string(b), string(newB)))
		return
	}
}
