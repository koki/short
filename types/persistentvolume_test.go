package types

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"

	"github.com/kr/pretty"

	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

var kokiPersistentGcePDVolume0 = PersistentVolumeSource{
	GcePD: kokiGcePD0.GcePD,
}

var kokiPersistentAwsEBSVolume0 = PersistentVolumeSource{
	AwsEBS: kokiAwsEBS0.AwsEBS,
}

var kokiPersistentHostPathVolume0 = PersistentVolumeSource{
	HostPath: kokiHostPath0.HostPath,
}

var kokiPersistentGlusterfsVolume0 = PersistentVolumeSource{
	Glusterfs: kokiGlusterfsVolume0.Glusterfs,
}

var kokiPersistentNFSVolume0 = PersistentVolumeSource{
	NFS: kokiNFSVolume0.NFS,
}

var kokiPersistentISCSIVolume0 = PersistentVolumeSource{
	ISCSI: &ISCSIPersistentVolume{
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
		SecretRef: &SecretReference{
			Namespace: "secret-ns",
			Name:      "secret-name",
		},
		InitiatorName: "iqn.1996-04.de.suse:linux-host1",
	},
}

var kokiPersistentCinderVolume0 = PersistentVolumeSource{
	Cinder: kokiCinder0.Cinder,
}

var kokiPersistentFibreChannelVolume0 = PersistentVolumeSource{
	FibreChannel: kokiFibreChannel0.FibreChannel,
}

var kokiPersistentFlockerVolume0 = PersistentVolumeSource{
	Flocker: kokiFlockerVolume0.Flocker,
}

var kokiPersistentFlexVolume0 = PersistentVolumeSource{
	Flex: kokiFlexVolume0.Flex,
}

var kokiPersistentVsphereVolume0 = PersistentVolumeSource{
	Vsphere: kokiVsphereVolume0.Vsphere,
}

var kokiPersistentQuobyteVolume0 = PersistentVolumeSource{
	Quobyte: kokiQuobyteVolume0.Quobyte,
}

var kokiPersistentAzureDiskVolume0 = PersistentVolumeSource{
	AzureDisk: kokiAzureDisk0.AzureDisk,
}

var kokiPersistentPhotonPDVolume0 = PersistentVolumeSource{
	PhotonPD: kokiPhotonPDVolume0.PhotonPD,
}

var kokiPersistentPortworxVolume0 = PersistentVolumeSource{
	Portworx: kokiPortworxVolume0.Portworx,
}

var kokiPersistentRBDVolume0 = PersistentVolumeSource{
	RBD: &RBDPersistentVolume{
		CephMonitors: []string{
			"1.2.3.4:6789",
			"1.2.3.5:6789",
		},
		RBDImage:  "foo",
		FSType:    "ext4",
		RBDPool:   "kube",
		RadosUser: "admin",
		Keyring:   "/etc/ceph/keyring",
		SecretRef: &SecretReference{
			Namespace: "secret-namespace",
			Name:      "secret-name",
		},
		ReadOnly: true,
	},
}
var kokiPersistentRBDVolume1 = PersistentVolumeSource{
	RBD: &RBDPersistentVolume{
		CephMonitors: []string{
			"1.2.3.4:6789",
			"1.2.3.5:6789",
		},
		RBDImage: "foo",
		SecretRef: &SecretReference{
			Name: "secret-name",
		},
	},
}

var kokiPersistentCephFS0 = PersistentVolumeSource{
	CephFS: &CephFSPersistentVolume{
		Monitors: []string{
			"1.2.3.4:6789",
			"1.2.3.5:6789",
		},
		Path: "/path",
		User: "admin",
		SecretFileOrRef: &CephFSPersistentSecretFileOrRef{
			File: "/etc/ceph/admin.secret",
		},
		ReadOnly: true,
	},
}
var kokiPersistentCephFS1 = PersistentVolumeSource{
	CephFS: &CephFSPersistentVolume{
		Monitors: []string{
			"1.2.3.4:6789",
			"1.2.3.5:6789",
		},
		Path: "/path",
		User: "admin",
		SecretFileOrRef: &CephFSPersistentSecretFileOrRef{
			Ref: &SecretReference{
				Namespace: "secret-namespace",
				Name:      "secret-name",
			},
		},
		ReadOnly: true,
	},
}

var kokiPersistentAzureFile0 = PersistentVolumeSource{
	AzureFile: &AzureFilePersistentVolume{
		Secret: SecretReference{
			Namespace: "secret-namespace",
			Name:      "secret-name",
		},
		ShareName: "k8stest",
		ReadOnly:  true,
	},
}
var kokiPersistentAzureFile1 = PersistentVolumeSource{
	AzureFile: &AzureFilePersistentVolume{
		Secret: SecretReference{
			Name: "secret-name",
		},
		ShareName: "k8stest",
	},
}

var kokiPersistentScaleIOVolume0 = PersistentVolumeSource{
	ScaleIO: &ScaleIOPersistentVolume{
		Gateway: "https://localhost:443/api",
		System:  "scaleio",
		SecretRef: SecretReference{
			Name:      "secret-name",
			Namespace: "secret-namespace",
		},
		SSLEnabled:       true,
		ProtectionDomain: "pd01",
		StoragePool:      "sp01",
		StorageMode:      "ThickProvisioned",
		VolumeName:       "vol-0",
		FSType:           "xfs",
		ReadOnly:         true,
	},
}

var kokiPersistentLocalVolume0 = PersistentVolumeSource{
	Local: &LocalVolume{
		Path: "/some/path",
	},
}

var kokiPersistentStorageOSVolume0 = PersistentVolumeSource{
	StorageOS: &StorageOSPersistentVolume{
		VolumeName:      "vol-0",
		VolumeNamespace: "namespace-0",
		FSType:          "ext4",
		ReadOnly:        true,
		SecretRef: &SecretReference{
			Name:      "secret-name",
			Namespace: "secret-namespace",
		},
	},
}

func TestPersistentVolume(t *testing.T) {
	testPersistentVolumeSource(kokiPersistentGcePDVolume0, t)
	testPersistentVolumeSource(kokiPersistentAwsEBSVolume0, t)
	testPersistentVolumeSource(kokiPersistentHostPathVolume0, t)
	testPersistentVolumeSource(kokiPersistentGlusterfsVolume0, t)
	testPersistentVolumeSource(kokiPersistentNFSVolume0, t)
	testPersistentVolumeSource(kokiPersistentISCSIVolume0, t)
	testPersistentVolumeSource(kokiPersistentCinderVolume0, t)
	testPersistentVolumeSource(kokiPersistentFibreChannelVolume0, t)
	testPersistentVolumeSource(kokiPersistentFlockerVolume0, t)
	testPersistentVolumeSource(kokiPersistentFlexVolume0, t)
	testPersistentVolumeSource(kokiPersistentVsphereVolume0, t)
	testPersistentVolumeSource(kokiPersistentQuobyteVolume0, t)
	testPersistentVolumeSource(kokiPersistentAzureDiskVolume0, t)
	testPersistentVolumeSource(kokiPersistentPhotonPDVolume0, t)
	testPersistentVolumeSource(kokiPersistentPortworxVolume0, t)
	testPersistentVolumeSource(kokiPersistentRBDVolume0, t)
	testPersistentVolumeSource(kokiPersistentRBDVolume1, t)
	testPersistentVolumeSource(kokiPersistentCephFS0, t)
	testPersistentVolumeSource(kokiPersistentCephFS1, t)
	testPersistentVolumeSource(kokiPersistentAzureFile0, t)
	testPersistentVolumeSource(kokiPersistentAzureFile1, t)
	testPersistentVolumeSource(kokiPersistentScaleIOVolume0, t)
	testPersistentVolumeSource(kokiPersistentLocalVolume0, t)
	testPersistentVolumeSource(kokiPersistentStorageOSVolume0, t)
}

func testPersistentVolumeSource(v PersistentVolumeSource, t *testing.T) {
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
			ReclaimPolicy: PersistentVolumeReclaimRecycle,
			StorageClass:  "storageClass",
			MountOptions:  "option 1,option 2,option 3",
			PersistentVolumeStatus: PersistentVolumeStatus{
				Phase:   VolumeAvailable,
				Message: "user-friendly message about the status",
				Reason:  "machineFriendlyReasonForStatus",
			},
		},
		PersistentVolumeSource: v,
	}

	b, err := yaml.Marshal(kokiVolume)
	if err != nil {
		t.Error(pretty.Sprintf("%s\n%# v", serrors.PrettyError(err), kokiVolume))
		return
	}

	newVolume := PersistentVolume{}

	err = yaml.Unmarshal(b, &newVolume)
	if err != nil {
		t.Error(pretty.Sprintf("%s\n(%s)\n(%# v)", serrors.PrettyError(err), string(b), kokiVolume))
		return
	}

	newB, err := yaml.Marshal(newVolume)
	if err != nil {
		t.Error(pretty.Sprintf("%s\n(%# v)\n(%# v)\n(%s)", serrors.PrettyError(err), newVolume, kokiVolume, string(b)))
		return
	}

	if !reflect.DeepEqual(kokiVolume, newVolume) {
		t.Error(pretty.Sprintf("failed round-trip\n(%# v)\n(%# v)\n(%s)\n(%s)", kokiVolume, newVolume, string(b), string(newB)))
		return
	}
}
