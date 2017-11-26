package types

import (
	"encoding/json"
	"strings"

	"github.com/koki/short/util"
	"k8s.io/apimachinery/pkg/api/resource"
)

type VolumeWrapper struct {
	Volume Volume `json:"volume"`
}

type Volume struct {
	HostPath *HostPathVolume
	EmptyDir *EmptyDirVolume
	GcePD    *GcePDVolume
	AwsEBS   *AwsEBSVolume
}

const (
	VolumeTypeHostPath = "host_path"
	VolumeTypeEmptyDir = "empty_dir"
	VolumeTypeGcePD    = "gce_pd"
	VolumeTypeAwsEBS   = "aws_ebs"
)

type HostPathVolume struct {
	Path string       `json:"-"`
	Type HostPathType `json:"-"`
}

type HostPathType string

const (
	HostPathUnset             HostPathType = ""
	HostPathDirectoryOrCreate HostPathType = "dir-or-create"
	HostPathDirectory         HostPathType = "dir"
	HostPathFileOrCreate      HostPathType = "file-or-create"
	HostPathFile              HostPathType = "file"
	HostPathSocket            HostPathType = "socket"
	HostPathCharDev           HostPathType = "char-dev"
	HostPathBlockDev          HostPathType = "block-dev"
)

type EmptyDirVolume struct {
	Medium    StorageMedium      `json:"medium,omitempty"`
	SizeLimit *resource.Quantity `json:"max_size,omitempty"`
}

type GcePDVolume struct {
	PDName    string `json:"-"`
	FSType    string `json:"fs,omitempty"`
	Partition int32  `json:"partition,omitempty"`
	ReadOnly  bool   `json:"ro,omitempty"`
}

type AwsEBSVolume struct {
	VolumeID  string `json:"-"`
	FSType    string `json:"fs,omitempty"`
	Partition int32  `json:"partition,omitempty"`
	ReadOnly  bool   `json:"ro,omitempty"`
}

type StorageMedium string

const (
	StorageMediumDefault   StorageMedium = ""           // use whatever the default is for the node
	StorageMediumMemory    StorageMedium = "memory"     // use memory (tmpfs)
	StorageMediumHugepages StorageMedium = "huge-pages" // use hugepages
)

func (v *Volume) UnmarshalJSON(data []byte) error {
	var err error
	str := ""
	err = json.Unmarshal(data, &str)
	if err == nil {
		segments := strings.Split(str, ":")
		return v.Unmarshal(nil, segments[0], segments[1:])
	}

	obj := map[string]interface{}{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return util.InvalidValueErrorf(string(data), "expected either string or dictionary")
	}

	selector := []string{}
	if val, ok := obj["vol_name"]; ok {
		if volName, ok := val.(string); ok {
			selector = append(selector, volName)
		} else {
			return util.InvalidValueErrorf(string(data), "expected string for key \"vol_name\"")
		}
	}

	volType, err := util.GetStringEntry(obj, "vol_type")
	if err != nil {
		return err
	}

	return v.Unmarshal(obj, volType, selector)
}

func (v *Volume) Unmarshal(obj map[string]interface{}, volType string, selector []string) error {
	switch volType {
	case VolumeTypeHostPath:
		return v.UnmarshalHostPathVolume(selector)
	case VolumeTypeEmptyDir:
		return v.UnmarshalEmptyDirVolume(obj, selector)
	case VolumeTypeGcePD:
		return v.UnmarshalGcePDVolume(obj, selector)
	case VolumeTypeAwsEBS:
		return v.UnmarshalAwsEBSVolume(obj, selector)
	default:
		return util.InvalidValueErrorf(volType, "unsupported volume type (%s)", volType)
	}
}

type MarshalledVolume struct {
	Type        string
	Selector    []string
	ExtraFields map[string]interface{}
}

func (v Volume) MarshalJSON() ([]byte, error) {
	var marshalledVolume *MarshalledVolume
	var err error
	if v.HostPath != nil {
		marshalledVolume, err = v.HostPath.Marshal()
	}

	if v.EmptyDir != nil {
		marshalledVolume, err = v.EmptyDir.Marshal()
	}

	if v.GcePD != nil {
		marshalledVolume, err = v.GcePD.Marshal()
	}

	if v.AwsEBS != nil {
		marshalledVolume, err = v.AwsEBS.Marshal()
	}

	if err != nil {
		return nil, err
	}

	if marshalledVolume == nil {
		return nil, util.InvalidInstanceErrorf(v, "empty volume definition")
	}

	if len(marshalledVolume.ExtraFields) == 0 {
		segments := []string{marshalledVolume.Type}
		segments = append(segments, marshalledVolume.Selector...)
		return json.Marshal(strings.Join(segments, ":"))
	}

	obj := marshalledVolume.ExtraFields
	obj["vol_type"] = marshalledVolume.Type
	if len(marshalledVolume.Selector) > 0 {
		obj["vol_name"] = strings.Join(marshalledVolume.Selector, ":")
	}

	return json.Marshal(obj)
}

func (v *Volume) UnmarshalHostPathVolume(selector []string) error {
	source := HostPathVolume{}
	if len(selector) > 2 || len(selector) == 0 {
		return util.InvalidValueErrorf(selector, "expected one or two selector segments for %s", VolumeTypeHostPath)
	}

	source.Path = selector[0]

	if len(selector) > 1 {
		hostPathType := HostPathType(selector[1])
		switch hostPathType {
		case HostPathUnset:
		case HostPathDirectoryOrCreate:
		case HostPathDirectory:
		case HostPathFileOrCreate:
		case HostPathFile:
		case HostPathSocket:
		case HostPathCharDev:
		case HostPathBlockDev:
		default:
			return util.InvalidValueErrorf(hostPathType, "invalid 'vol_type' selector for %s", VolumeTypeHostPath)
		}

		source.Type = hostPathType
	}

	v.HostPath = &source
	return nil
}

func (s HostPathVolume) Marshal() (*MarshalledVolume, error) {
	var selector []string
	if len(s.Type) > 0 {
		selector = []string{s.Path, string(s.Type)}
	} else {
		selector = []string{s.Path}
	}
	return &MarshalledVolume{
		Type:     VolumeTypeHostPath,
		Selector: selector,
	}, nil
}

func (v *Volume) UnmarshalEmptyDirVolume(obj map[string]interface{}, selector []string) error {
	source := EmptyDirVolume{}
	if len(selector) > 0 {
		return util.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeEmptyDir)
	}

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeEmptyDir)
	}

	v.EmptyDir = &source
	return nil
}

func (s EmptyDirVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeEmptyDir)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeEmptyDir,
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalGcePDVolume(obj map[string]interface{}, selector []string) error {
	source := GcePDVolume{}
	if len(selector) != 1 {
		return util.InvalidValueErrorf(selector, "expected 1 selector segment (disk name) for %s", VolumeTypeGcePD)
	}
	source.PDName = selector[0]

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeGcePD)
	}

	v.GcePD = &source
	return nil
}

func (s GcePDVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeGcePD)
	}

	if len(s.PDName) == 0 {
		return nil, util.InvalidInstanceErrorf(&s, "selector must contain disk name")
	}

	return &MarshalledVolume{
		Type:        VolumeTypeGcePD,
		Selector:    []string{s.PDName},
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalAwsEBSVolume(obj map[string]interface{}, selector []string) error {
	source := AwsEBSVolume{}
	if len(selector) != 1 {
		return util.InvalidValueErrorf(selector, "expected 1 selector segment (ebs uuid) for %s", VolumeTypeAwsEBS)
	}
	source.VolumeID = selector[0]

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeAwsEBS)
	}

	v.AwsEBS = &source
	return nil
}

func (s AwsEBSVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeAwsEBS)
	}

	if len(s.VolumeID) == 0 {
		return nil, util.InvalidInstanceErrorf(&s, "selector must contain ebs uuid")
	}

	return &MarshalledVolume{
		Type:        VolumeTypeAwsEBS,
		Selector:    []string{s.VolumeID},
		ExtraFields: obj,
	}, nil
}
