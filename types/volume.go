package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/koki/short/util"
	"k8s.io/apimachinery/pkg/api/resource"
)

type VolumeWrapper struct {
	Volume Volume `json:"volume"`
}

type Volume struct {
	HostPath     *HostPathVolume
	EmptyDir     *EmptyDirVolume
	GcePD        *GcePDVolume
	AwsEBS       *AwsEBSVolume
	AzureDisk    *AzureDiskVolume
	AzureFile    *AzureFileVolume
	CephFS       *CephFSVolume
	Cinder       *CinderVolume
	FibreChannel *FibreChannelVolume
	Flex         *FlexVolume
	Flocker      *FlockerVolume
	Glusterfs    *GlusterfsVolume
	ISCSI        *ISCSIVolume
	NFS          *NFSVolume
	PhotonPD     *PhotonPDVolume
	Portworx     *PortworxVolume
	PVC          *PVCVolume
	Quobyte      *QuobyteVolume
	ScaleIO      *ScaleIOVolume
}

const (
	VolumeTypeHostPath     = "host_path"
	VolumeTypeEmptyDir     = "empty_dir"
	VolumeTypeGcePD        = "gce_pd"
	VolumeTypeAwsEBS       = "aws_ebs"
	VolumeTypeAzureDisk    = "azure_disk"
	VolumeTypeAzureFile    = "azure_file"
	VolumeTypeCephFS       = "cephfs"
	VolumeTypeCinder       = "cinder"
	VolumeTypeFibreChannel = "fc"
	VolumeTypeFlex         = "flex"
	VolumeTypeFlocker      = "flocker"
	VolumeTypeGlusterfs    = "glusterfs"
	VolumeTypeISCSI        = "iscsi"
	VolumeTypeNFS          = "nfs"
	VolumeTypePhotonPD     = "photon"
	VolumeTypePortworx     = "portworx"
	VolumeTypePVC          = "pvc"
	VolumeTypeQuobyte      = "quobyte"
	VolumeTypeScaleIO      = "scaleio"

	SelectorSegmentReadOnly = "ro"
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

type StorageMedium string

const (
	StorageMediumDefault   StorageMedium = ""           // use whatever the default is for the node
	StorageMediumMemory    StorageMedium = "memory"     // use memory (tmpfs)
	StorageMediumHugepages StorageMedium = "huge-pages" // use hugepages
)

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

type AzureDiskVolume struct {
	DiskName string `json:"-"`
	// DataDiskURI is required
	DataDiskURI string                    `json:"disk_uri"`
	CachingMode *AzureDataDiskCachingMode `json:"cache,omitempty"`
	FSType      string                    `json:"fs,omitempty"`
	ReadOnly    bool                      `json:"ro,omitempty"`
	Kind        *AzureDataDiskKind        `json:"kind,omitempty"`
}

type AzureDataDiskCachingMode string
type AzureDataDiskKind string

const (
	AzureDataDiskCachingNone      AzureDataDiskCachingMode = "none"
	AzureDataDiskCachingReadOnly  AzureDataDiskCachingMode = "ro"
	AzureDataDiskCachingReadWrite AzureDataDiskCachingMode = "rw"

	AzureSharedBlobDisk    AzureDataDiskKind = "shared"
	AzureDedicatedBlobDisk AzureDataDiskKind = "dedicated"
	AzureManagedDisk       AzureDataDiskKind = "managed"
)

type AzureFileVolume struct {
	SecretName string `json:"-"`
	ShareName  string `json:"-"`
	ReadOnly   bool   `json:"-"`
}

type CephFSVolume struct {
	Monitors        []string               `json:"monitors"`
	Path            string                 `json:"path, omitempty"`
	User            string                 `json:"user,omitempty"`
	SecretFileOrRef *CephFSSecretFileOrRef `json:"secret,omitempty"`
	ReadOnly        bool                   `json:"ro,omitempty"`
}

type CephFSSecretFileOrRef struct {
	File string `json:"-"`
	Ref  string `json:"-"`
}

type CinderVolume struct {
	VolumeID string `json:"-"`
	FSType   string `json:"fs,omitempty"`
	ReadOnly bool   `json:"ro,omitempty"`
}

type FibreChannelVolume struct {
	TargetWWNs []string `json:"wwn,omitempty"`
	Lun        *int32   `json:"lun,omitempty"`
	FSType     string   `json:"fs,omitempty"`
	ReadOnly   bool     `json:"ro,omitempty"`
	WWIDs      []string `json:"wwid,omitempty"`
}

type FlexVolume struct {
	Driver    string            `json:"-"`
	FSType    string            `json:"fs,omitempty"`
	SecretRef string            `json:"secret,omitempty"`
	ReadOnly  bool              `json:"ro,omitempty"`
	Options   map[string]string `json:"options,omitempty"`
}

type FlockerVolume struct {
	DatasetUUID string `json:"-"`
}

type GlusterfsVolume struct {
	EndpointsName string `json:"-"`
	Path          string `json:"path"`
	ReadOnly      bool   `json:"ro,omitempty"`
}

type ISCSIVolume struct {
	TargetPortal   string   `json:"target_portal"`
	IQN            string   `json:"iqn"`
	Lun            int32    `json:"lun"`
	ISCSIInterface string   `json:"iscsi_interface,omitempty"`
	FSType         string   `json:"fs,omitempty"`
	ReadOnly       bool     `json:"ro,omitempty"`
	Portals        []string `json:"portals,omitempty"`
	// TODO: should this actually be "chap_auth"?
	DiscoveryCHAPAuth bool   `json:"chap_discovery,omitempty"`
	SessionCHAPAuth   bool   `json:"chap_session,omitempty"`
	SecretRef         string `json:"secret,omitempty"`
	// NOTE: InitiatorName is a pointer in k8s
	InitiatorName string `json:"initiator,omitempty"`
}

type NFSVolume struct {
	Server   string `json:"-"`
	Path     string `json:"-"`
	ReadOnly bool   `json:"-"`
}

type PhotonPDVolume struct {
	PdID   string `json:"-"`
	FSType string `json:"-"`
}

type PortworxVolume struct {
	VolumeID string `json:"-"`
	FSType   string `json:"fs,omitempty"`
	ReadOnly bool   `json:"ro,omitempty"`
}

type PVCVolume struct {
	ClaimName string `json:"-"`
	ReadOnly  bool   `json:"-"`
}

type QuobyteVolume struct {
	Registry string `json:"registry"`
	Volume   string `json:"-"`
	ReadOnly bool   `json:"ro,omitempty"`
	User     string `json:"user,omitempty"`
	Group    string `json:"group,omitempty"`
}

type ScaleIOVolume struct {
	Gateway          string `json:"gateway"`
	System           string `json:"system"`
	SecretRef        string `json:"secret"`
	SSLEnabled       bool   `json:"ssl,omitempty"`
	ProtectionDomain string `json:"protection_domain,omitempty"`
	StoragePool      string `json:"storage_pool,omitempty"`
	StorageMode      string `json:"storage_mode,omitempty"`
	VolumeName       string `json:"-"`
	FSType           string `json:"fs,omitempty"`
	ReadOnly         bool   `json:"ro,omitempty"`
}

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
	case VolumeTypeAzureDisk:
		return v.UnmarshalAzureDiskVolume(obj, selector)
	case VolumeTypeAzureFile:
		return v.UnmarshalAzureFileVolume(selector)
	case VolumeTypeCephFS:
		return v.UnmarshalCephFSVolume(obj, selector)
	case VolumeTypeCinder:
		return v.UnmarshalCinderVolume(obj, selector)
	case VolumeTypeFibreChannel:
		return v.UnmarshalFibreChannelVolume(obj, selector)
	case VolumeTypeFlex:
		return v.UnmarshalFlexVolume(obj, selector)
	case VolumeTypeFlocker:
		return v.UnmarshalFlockerVolume(selector)
	case VolumeTypeGlusterfs:
		return v.UnmarshalGlusterfsVolume(obj, selector)
	case VolumeTypeISCSI:
		return v.UnmarshalISCSIVolume(obj, selector)
	case VolumeTypeNFS:
		return v.UnmarshalNFSVolume(selector)
	case VolumeTypePhotonPD:
		return v.UnmarshalPhotonPDVolume(selector)
	case VolumeTypePortworx:
		return v.UnmarshalPortworxVolume(obj, selector)
	case VolumeTypePVC:
		return v.UnmarshalPVCVolume(selector)
	case VolumeTypeQuobyte:
		return v.UnmarshalQuobyteVolume(obj, selector)
	case VolumeTypeScaleIO:
		return v.UnmarshalScaleIOVolume(obj, selector)
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
	if v.AzureDisk != nil {
		marshalledVolume, err = v.AzureDisk.Marshal()
	}
	if v.AzureFile != nil {
		marshalledVolume, err = v.AzureFile.Marshal()
	}
	if v.CephFS != nil {
		marshalledVolume, err = v.CephFS.Marshal()
	}
	if v.Cinder != nil {
		marshalledVolume, err = v.Cinder.Marshal()
	}
	if v.FibreChannel != nil {
		marshalledVolume, err = v.FibreChannel.Marshal()
	}
	if v.Flex != nil {
		marshalledVolume, err = v.Flex.Marshal()
	}
	if v.Flocker != nil {
		marshalledVolume, err = v.Flocker.Marshal()
	}
	if v.Glusterfs != nil {
		marshalledVolume, err = v.Glusterfs.Marshal()
	}
	if v.ISCSI != nil {
		marshalledVolume, err = v.ISCSI.Marshal()
	}
	if v.NFS != nil {
		marshalledVolume, err = v.NFS.Marshal()
	}
	if v.PhotonPD != nil {
		marshalledVolume, err = v.PhotonPD.Marshal()
	}
	if v.Portworx != nil {
		marshalledVolume, err = v.Portworx.Marshal()
	}
	if v.PVC != nil {
		marshalledVolume, err = v.PVC.Marshal()
	}
	if v.Quobyte != nil {
		marshalledVolume, err = v.Quobyte.Marshal()
	}
	if v.ScaleIO != nil {
		marshalledVolume, err = v.ScaleIO.Marshal()
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

func (v *Volume) UnmarshalAzureDiskVolume(obj map[string]interface{}, selector []string) error {
	source := AzureDiskVolume{}
	if len(selector) != 1 {
		return util.InvalidValueErrorf(selector, "expected 1 selector segment (disk name) for %s", VolumeTypeAzureDisk)
	}
	source.DiskName = selector[0]

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeAzureDisk)
	}

	v.AzureDisk = &source
	return nil
}

func (s AzureDiskVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeAzureDisk)
	}

	if len(s.DiskName) == 0 {
		return nil, util.InvalidInstanceErrorf(&s, "selector must contain disk name")
	}

	return &MarshalledVolume{
		Type:        VolumeTypeAzureDisk,
		Selector:    []string{s.DiskName},
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalAzureFileVolume(selector []string) error {
	source := AzureFileVolume{}
	if len(selector) > 3 || len(selector) < 2 {
		return util.InvalidValueErrorf(selector, "expected two or three selector segments for %s", VolumeTypeAzureFile)
	}

	source.SecretName = selector[0]
	source.ShareName = selector[1]

	if len(selector) > 2 {
		switch selector[2] {
		case SelectorSegmentReadOnly:
			source.ReadOnly = true
		default:
			return util.InvalidValueErrorf(selector[2], "invalid selector segment for %s", VolumeTypeAzureFile)
		}
	}

	v.AzureFile = &source
	return nil
}

func (s AzureFileVolume) Marshal() (*MarshalledVolume, error) {
	var selector []string
	if s.ReadOnly {
		selector = []string{s.SecretName, s.ShareName, SelectorSegmentReadOnly}
	} else {
		selector = []string{s.SecretName, s.ShareName}
	}
	return &MarshalledVolume{
		Type:     VolumeTypeAzureFile,
		Selector: selector,
	}, nil
}

func (v *Volume) UnmarshalCephFSVolume(obj map[string]interface{}, selector []string) error {
	source := CephFSVolume{}
	if len(selector) != 0 {
		return util.InvalidValueErrorf(selector, "expected 0 selector segments for %s", VolumeTypeCephFS)
	}

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeCephFS)
	}

	v.CephFS = &source
	return nil
}

func (s CephFSVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeCephFS)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeCephFS,
		ExtraFields: obj,
	}, nil
}

var fileOrRefRegexp = regexp.MustCompile(`^(file|ref):(.*)$`)

func (s *CephFSSecretFileOrRef) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return util.ContextualizeErrorf(err, "cephfs secret should be a string")
	}

	matches := fileOrRefRegexp.FindStringSubmatch(str)
	if len(matches) > 0 {
		if matches[1] == "file" {
			s.File = matches[2]
		} else {
			s.Ref = matches[2]
		}
	} else {
		return util.InvalidValueErrorf(string(data), "unrecognized format for cephfs secret")
	}

	return nil
}

func (s CephFSSecretFileOrRef) MarshalJSON() ([]byte, error) {
	if len(s.Ref) > 0 {
		return json.Marshal(fmt.Sprintf("ref:%s", s.Ref))
	}

	return json.Marshal(fmt.Sprintf("file:%s", s.File))
}

func (v *Volume) UnmarshalCinderVolume(obj map[string]interface{}, selector []string) error {
	source := CinderVolume{}
	if len(selector) != 1 {
		return util.InvalidValueErrorf(selector, "expected 1 (volume ID) selector segment for %s", VolumeTypeCinder)
	}

	source.VolumeID = selector[0]

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeCinder)
	}

	v.Cinder = &source
	return nil
}

func (s CinderVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeCinder)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeCinder,
		Selector:    []string{s.VolumeID},
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalFibreChannelVolume(obj map[string]interface{}, selector []string) error {
	source := FibreChannelVolume{}
	if len(selector) != 0 {
		return util.InvalidValueErrorf(selector, "expected 0 selector segments for %s", VolumeTypeFibreChannel)
	}

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeFibreChannel)
	}

	v.FibreChannel = &source
	return nil
}

func (s FibreChannelVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeFibreChannel)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeFibreChannel,
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalFlexVolume(obj map[string]interface{}, selector []string) error {
	source := FlexVolume{}
	if len(selector) != 1 {
		return util.InvalidValueErrorf(selector, "expected 1 selector segment (driver) for %s", VolumeTypeFlex)
	}
	source.Driver = selector[0]

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeFlex)
	}

	v.Flex = &source
	return nil
}

func (s FlexVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeFlex)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeFlex,
		Selector:    []string{s.Driver},
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalFlockerVolume(selector []string) error {
	source := FlockerVolume{}
	if len(selector) != 1 {
		return util.InvalidValueErrorf(selector, "expected exactly one selector segment (dataset UUID) for %s", VolumeTypeFlocker)
	}

	source.DatasetUUID = selector[0]

	v.Flocker = &source
	return nil
}

func (s FlockerVolume) Marshal() (*MarshalledVolume, error) {
	return &MarshalledVolume{
		Type:     VolumeTypeFlocker,
		Selector: []string{s.DatasetUUID},
	}, nil
}

func (v *Volume) UnmarshalGlusterfsVolume(obj map[string]interface{}, selector []string) error {
	source := GlusterfsVolume{}
	if len(selector) != 1 {
		return util.InvalidValueErrorf(selector, "expected 1 selector segment (endpoints name) for %s", VolumeTypeGlusterfs)
	}
	source.EndpointsName = selector[0]

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeGlusterfs)
	}

	v.Glusterfs = &source
	return nil
}

func (s GlusterfsVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeGlusterfs)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeGlusterfs,
		Selector:    []string{s.EndpointsName},
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalISCSIVolume(obj map[string]interface{}, selector []string) error {
	source := ISCSIVolume{}
	if len(selector) != 0 {
		return util.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeISCSI)
	}

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeISCSI)
	}

	v.ISCSI = &source
	return nil
}

func (s ISCSIVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeISCSI)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeISCSI,
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalNFSVolume(selector []string) error {
	source := NFSVolume{}
	if len(selector) > 3 || len(selector) < 2 {
		return util.InvalidValueErrorf(selector, "expected two or three selector segments for %s", VolumeTypeNFS)
	}

	source.Server = selector[0]
	source.Path = selector[1]

	if len(selector) > 2 {
		switch selector[2] {
		case SelectorSegmentReadOnly:
			source.ReadOnly = true
		default:
			return util.InvalidValueErrorf(selector[2], "invalid selector segment for %s", VolumeTypeNFS)
		}
	}

	v.NFS = &source
	return nil
}

func (s NFSVolume) Marshal() (*MarshalledVolume, error) {
	var selector []string
	if s.ReadOnly {
		selector = []string{s.Server, s.Path, SelectorSegmentReadOnly}
	} else {
		selector = []string{s.Server, s.Path}
	}
	return &MarshalledVolume{
		Type:     VolumeTypeNFS,
		Selector: selector,
	}, nil
}

func (v *Volume) UnmarshalPhotonPDVolume(selector []string) error {
	source := PhotonPDVolume{}
	if len(selector) > 2 || len(selector) < 1 {
		return util.InvalidValueErrorf(selector, "expected one or two selector segments for %s", VolumeTypePhotonPD)
	}

	source.PdID = selector[0]

	if len(selector) > 1 {
		source.FSType = selector[1]
	}

	v.PhotonPD = &source
	return nil
}

func (s PhotonPDVolume) Marshal() (*MarshalledVolume, error) {
	var selector []string
	if len(s.FSType) > 0 {
		selector = []string{s.PdID, s.FSType}
	} else {
		selector = []string{s.PdID}
	}
	return &MarshalledVolume{
		Type:     VolumeTypePhotonPD,
		Selector: selector,
	}, nil
}

func (v *Volume) UnmarshalPortworxVolume(obj map[string]interface{}, selector []string) error {
	source := PortworxVolume{}
	if len(selector) != 1 {
		return util.InvalidValueErrorf(selector, "expected 1 selector segment (volume ID) for %s", VolumeTypePortworx)
	}
	source.VolumeID = selector[0]

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypePortworx)
	}

	v.Portworx = &source
	return nil
}

func (s PortworxVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypePortworx)
	}

	return &MarshalledVolume{
		Type:        VolumeTypePortworx,
		Selector:    []string{s.VolumeID},
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalPVCVolume(selector []string) error {
	source := PVCVolume{}
	if len(selector) > 2 || len(selector) < 1 {
		return util.InvalidValueErrorf(selector, "expected one or two selector segments for %s", VolumeTypePVC)
	}

	source.ClaimName = selector[0]

	if len(selector) > 1 {
		switch selector[1] {
		case SelectorSegmentReadOnly:
			source.ReadOnly = true
		default:
			return util.InvalidValueErrorf(selector[2], "invalid selector segment for %s", VolumeTypePVC)
		}
	}

	v.PVC = &source
	return nil
}

func (s PVCVolume) Marshal() (*MarshalledVolume, error) {
	var selector []string
	if s.ReadOnly {
		selector = []string{s.ClaimName, SelectorSegmentReadOnly}
	} else {
		selector = []string{s.ClaimName}
	}
	return &MarshalledVolume{
		Type:     VolumeTypePVC,
		Selector: selector,
	}, nil
}

func (v *Volume) UnmarshalQuobyteVolume(obj map[string]interface{}, selector []string) error {
	source := QuobyteVolume{}
	if len(selector) != 1 {
		return util.InvalidValueErrorf(selector, "expected 1 selector segment (volume ID) for %s", VolumeTypeQuobyte)
	}
	source.Volume = selector[0]

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeQuobyte)
	}

	v.Quobyte = &source
	return nil
}

func (s QuobyteVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeQuobyte)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeQuobyte,
		Selector:    []string{s.Volume},
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalScaleIOVolume(obj map[string]interface{}, selector []string) error {
	source := ScaleIOVolume{}
	if len(selector) != 1 {
		return util.InvalidValueErrorf(selector, "expected 1 selector segment (volume name) for %s", VolumeTypeScaleIO)
	}
	source.VolumeName = selector[0]

	err := util.UnmarshalMap(obj, &source)
	if err != nil {
		return util.ContextualizeErrorf(err, VolumeTypeScaleIO)
	}

	v.ScaleIO = &source
	return nil
}

func (s ScaleIOVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := util.MarshalMap(&s)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, VolumeTypeScaleIO)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeScaleIO,
		Selector:    []string{s.VolumeName},
		ExtraFields: obj,
	}, nil
}
