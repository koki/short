package types

import (
	"reflect"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/ghodss/yaml"
	"github.com/kr/pretty"
)

var kokiHostPath0 = Volume{
	HostPath: &HostPathVolume{
		Path: "/path",
		Type: HostPathDirectoryOrCreate,
	},
}

var sizeLimit0 = resource.MustParse("10Gi")
var kokiEmptyDir0 = Volume{
	EmptyDir: &EmptyDirVolume{
		Medium:    "memory",
		SizeLimit: &sizeLimit0,
	},
}
var kokiEmptyDir1 = Volume{
	EmptyDir: &EmptyDirVolume{},
}

var kokiGcePD0 = Volume{
	GcePD: &GcePDVolume{
		PDName:    "data-disk",
		FSType:    "ext4",
		Partition: 1,
		ReadOnly:  true,
	},
}
var kokiGcePD1 = Volume{
	GcePD: &GcePDVolume{
		PDName: "data-disk",
	},
}

var kokiAwsEBS0 = Volume{
	AwsEBS: &AwsEBSVolume{
		VolumeID:  "ebs_uuid",
		FSType:    "xfs",
		Partition: 1,
		ReadOnly:  true,
	},
}
var kokiAwsEBS1 = Volume{
	AwsEBS: &AwsEBSVolume{
		VolumeID: "ebs_uuid",
	},
}

func TestVolume(t *testing.T) {
	testVolumeSource(kokiHostPath0, t, true)
	testVolumeSource(kokiEmptyDir0, t, false)
	testVolumeSource(kokiEmptyDir1, t, true)
	testVolumeSource(kokiGcePD0, t, false)
	testVolumeSource(kokiGcePD1, t, true)
	testVolumeSource(kokiAwsEBS0, t, false)
	testVolumeSource(kokiAwsEBS1, t, true)

}

func isString(data []byte, t *testing.T) bool {
	str := ""
	err := yaml.Unmarshal(data, &str)
	return err == nil
}

func testVolumeSource(kokiVolume Volume, t *testing.T, expectString bool) {
	b, err := yaml.Marshal(kokiVolume)
	if err != nil {
		t.Error(pretty.Sprintf("%s\n%# v", err.Error(), kokiVolume))
		return
	}

	if expectString != isString(b, t) {
		t.Error(pretty.Sprintf("unexpected serialization\n(%s)\n(%# v)", string(b), kokiVolume))
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
