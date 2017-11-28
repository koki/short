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

func TestPersistentVolume(t *testing.T) {
	testPersistentVolumeSource(kokiPersistentGcePDVolume0, t)
	testPersistentVolumeSource(kokiPersistentAwsEBSVolume0, t)
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
