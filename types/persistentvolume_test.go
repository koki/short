package types

import (
	"reflect"
	"testing"

	"k8s.io/api/core/v1"

	"github.com/ghodss/yaml"
	"github.com/kr/pretty"
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
	ISCSI: kokiISCSIVolume0.ISCSI,
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
			Status:        &v1.PersistentVolumeStatus{},
		},
		PersistentVolumeSource: v,
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
