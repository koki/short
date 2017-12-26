package types

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/golang/glog"

	"github.com/koki/json"
	"github.com/koki/json/jsonutil"
	serrors "github.com/koki/structurederrors"
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
	Vsphere      *VsphereVolume
	ConfigMap    *ConfigMapVolume
	Secret       *SecretVolume
	DownwardAPI  *DownwardAPIVolume
	Projected    *ProjectedVolume
	Git          *GitVolume
	RBD          *RBDVolume
	StorageOS    *StorageOSVolume
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
	VolumeTypeVsphere      = "vsphere"
	VolumeTypeConfigMap    = "config-map"
	VolumeTypeSecret       = "secret"
	VolumeTypeDownwardAPI  = "downward_api"
	VolumeTypeProjected    = "projected"
	VolumeTypeGit          = "git"
	VolumeTypeRBD          = "rbd"
	VolumeTypeStorageOS    = "storageos"
	VolumeTypeAny          = "*"

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
	StorageMediumHugePages StorageMedium = "huge-pages" // use hugepages
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
	DiskName    string                    `json:"disk_name"`
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
	EndpointsName string `json:"endpoints"`

	// Path is the Glusterfs volume name.
	Path     string `json:"path"`
	ReadOnly bool   `json:"ro,omitempty"`
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
	Gateway          string             `json:"gateway"`
	System           string             `json:"system"`
	SecretRef        string             `json:"secret"`
	SSLEnabled       bool               `json:"ssl,omitempty"`
	ProtectionDomain string             `json:"protection_domain,omitempty"`
	StoragePool      string             `json:"storage_pool,omitempty"`
	StorageMode      ScaleIOStorageMode `json:"storage_mode,omitempty"`
	VolumeName       string             `json:"-"`
	FSType           string             `json:"fs,omitempty"`
	ReadOnly         bool               `json:"ro,omitempty"`
}

type ScaleIOStorageMode string

const (
	ScaleIOStorageModeThick ScaleIOStorageMode = "thick"
	ScaleIOStorageModeThin  ScaleIOStorageMode = "thin"
)

type VsphereVolume struct {
	VolumePath    string                `json:"-"`
	FSType        string                `json:"fs,omitempty"`
	StoragePolicy *VsphereStoragePolicy `json:"policy,omitempty"`
}

type VsphereStoragePolicy struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

type ConfigMapVolume struct {
	Name string `json:"-"`

	Items       map[string]KeyAndMode `json:"items,omitempty"`
	DefaultMode *FileMode             `json:"mode,omitempty"`

	// NOTE: opposite of Optional
	Required *bool `json:"required,omitempty"`
}

type KeyAndMode struct {
	Key  string    `json:"-"`
	Mode *FileMode `json:"-"`
}

type SecretVolume struct {
	SecretName string `json:"-"`

	Items       map[string]KeyAndMode `json:"items,omitempty"`
	DefaultMode *FileMode             `json:"mode,omitempty"`

	// NOTE: opposite of Optional
	Required *bool `json:"required,omitempty"`
}

// FileMode can be unmarshalled from either a number (octal is supported) or a string.
// The json library doesn't allow serializing numbers as octal, so FileMode always marshals to a string.
type FileMode int32

type DownwardAPIVolume struct {
	Items       map[string]DownwardAPIVolumeFile `json:"items,omitempty"`
	DefaultMode *FileMode                        `json:"mode,omitempty"`
}

type DownwardAPIVolumeFile struct {
	FieldRef         *ObjectFieldSelector         `json:"field,omitempty"`
	ResourceFieldRef *VolumeResourceFieldSelector `json:"resource,omitempty"`
	Mode             *FileMode                    `json:"mode,omitempty"`
}

type ObjectFieldSelector struct {
	// required
	FieldPath string `json:"-"`

	// optional
	APIVersion string `json:"-"`
}

type VolumeResourceFieldSelector struct {
	// required
	ContainerName string `json:"-"`

	// required
	Resource string `json:"-"`

	// optional
	Divisor resource.Quantity `json:"-"`
}

type ProjectedVolume struct {
	Sources     []VolumeProjection `json:"sources"`
	DefaultMode *FileMode          `json:"mode,omitempty"`
}

type VolumeProjection struct {
	Secret      *SecretProjection      `json:"-"`
	DownwardAPI *DownwardAPIProjection `json:"-"`
	ConfigMap   *ConfigMapProjection   `json:"-"`
}

type SecretProjection struct {
	Name string `json:"secret"`

	Items map[string]KeyAndMode `json:"items,omitempty"`

	// NOTE: opposite of Optional
	Required *bool `json:"required,omitempty"`
}

type ConfigMapProjection struct {
	Name string `json:"config"`

	Items map[string]KeyAndMode `json:"items,omitempty"`

	// NOTE: opposite of Optional
	Required *bool `json:"required,omitempty"`
}

type DownwardAPIProjection struct {
	Items map[string]DownwardAPIVolumeFile `json:"items,omitempty"`
}

type GitVolume struct {
	Repository string `json:"-"`
	Revision   string `json:"rev,omitempty"`
	Directory  string `json:"dir,omitempty"`
}

type RBDVolume struct {
	CephMonitors []string `json:"monitors"`
	RBDImage     string   `json:"image"`
	FSType       string   `json:"fs,omitempty"`
	RBDPool      string   `json:"pool,omitempty"`
	RadosUser    string   `json:"user,omitempty"`
	Keyring      string   `json:"keyring,omitempty"`
	SecretRef    string   `json:"secret,omitempty"`
	ReadOnly     bool     `json:"ro,omitempty"`
}

type StorageOSVolume struct {
	VolumeName      string `json:"-"`
	VolumeNamespace string `json:"vol_ns,omitempty"`
	FSType          string `json:"fs,omitempty"`
	ReadOnly        bool   `json:"ro,omitempty"`
	SecretRef       string `json:"secret,omitempty"`
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
		return serrors.InvalidValueErrorf(string(data), "expected either string or dictionary")
	}

	selector := []string{}
	if val, ok := obj["vol_id"]; ok {
		if volName, ok := val.(string); ok {
			selector = append(selector, volName)
		} else {
			return serrors.InvalidValueErrorf(string(data), "expected string for key \"vol_id\"")
		}
	}

	volType, err := jsonutil.GetStringEntry(obj, "vol_type")
	if err != nil {
		return err
	}

	return v.Unmarshal(obj, volType, selector)
}

func (v *Volume) Unmarshal(obj map[string]interface{}, volType string, selector []string) error {
	switch volType {
	case VolumeTypeHostPath:
		v.HostPath = &HostPathVolume{}
		return v.HostPath.Unmarshal(selector)
	case VolumeTypeEmptyDir:
		return v.UnmarshalEmptyDirVolume(obj, selector)
	case VolumeTypeGcePD:
		v.GcePD = &GcePDVolume{}
		return v.GcePD.Unmarshal(obj, selector)
	case VolumeTypeAwsEBS:
		v.AwsEBS = &AwsEBSVolume{}
		return v.AwsEBS.Unmarshal(obj, selector)
	case VolumeTypeAzureDisk:
		v.AzureDisk = &AzureDiskVolume{}
		return v.AzureDisk.Unmarshal(obj, selector)
	case VolumeTypeAzureFile:
		return v.UnmarshalAzureFileVolume(selector)
	case VolumeTypeCephFS:
		return v.UnmarshalCephFSVolume(obj, selector)
	case VolumeTypeCinder:
		v.Cinder = &CinderVolume{}
		return v.Cinder.Unmarshal(obj, selector)
	case VolumeTypeFibreChannel:
		v.FibreChannel = &FibreChannelVolume{}
		return v.FibreChannel.Unmarshal(obj, selector)
	case VolumeTypeFlex:
		v.Flex = &FlexVolume{}
		return v.Flex.Unmarshal(obj, selector)
	case VolumeTypeFlocker:
		v.Flocker = &FlockerVolume{}
		return v.Flocker.Unmarshal(selector)
	case VolumeTypeGlusterfs:
		v.Glusterfs = &GlusterfsVolume{}
		return v.Glusterfs.Unmarshal(obj, selector)
	case VolumeTypeISCSI:
		v.ISCSI = &ISCSIVolume{}
		return v.ISCSI.Unmarshal(obj, selector)
	case VolumeTypeNFS:
		v.NFS = &NFSVolume{}
		return v.NFS.Unmarshal(selector)
	case VolumeTypePhotonPD:
		v.PhotonPD = &PhotonPDVolume{}
		return v.PhotonPD.Unmarshal(selector)
	case VolumeTypePortworx:
		v.Portworx = &PortworxVolume{}
		return v.Portworx.Unmarshal(obj, selector)
	case VolumeTypePVC:
		return v.UnmarshalPVCVolume(selector)
	case VolumeTypeQuobyte:
		v.Quobyte = &QuobyteVolume{}
		return v.Quobyte.Unmarshal(obj, selector)
	case VolumeTypeScaleIO:
		return v.UnmarshalScaleIOVolume(obj, selector)
	case VolumeTypeVsphere:
		v.Vsphere = &VsphereVolume{}
		return v.Vsphere.Unmarshal(obj, selector)
	case VolumeTypeConfigMap:
		return v.UnmarshalConfigMapVolume(obj, selector)
	case VolumeTypeSecret:
		return v.UnmarshalSecretVolume(obj, selector)
	case VolumeTypeDownwardAPI:
		return v.UnmarshalDownwardAPIVolume(obj, selector)
	case VolumeTypeProjected:
		return v.UnmarshalProjectedVolume(obj, selector)
	case VolumeTypeGit:
		return v.UnmarshalGitVolume(obj, selector)
	case VolumeTypeRBD:
		v.RBD = &RBDVolume{}
		return v.RBD.Unmarshal(obj, selector)
	case VolumeTypeStorageOS:
		v.StorageOS = &StorageOSVolume{}
		return v.StorageOS.Unmarshal(obj, selector)
	default:
		return serrors.InvalidValueErrorf(volType, "unsupported volume type (%s)", volType)
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
	if v.Vsphere != nil {
		marshalledVolume, err = v.Vsphere.Marshal()
	}
	if v.ConfigMap != nil {
		marshalledVolume, err = v.ConfigMap.Marshal()
	}
	if v.Secret != nil {
		marshalledVolume, err = v.Secret.Marshal()
	}
	if v.DownwardAPI != nil {
		marshalledVolume, err = v.DownwardAPI.Marshal()
	}
	if v.Projected != nil {
		marshalledVolume, err = v.Projected.Marshal()
	}
	if v.Git != nil {
		marshalledVolume, err = v.Git.Marshal()
	}
	if v.RBD != nil {
		marshalledVolume, err = v.RBD.Marshal()
	}
	if v.StorageOS != nil {
		marshalledVolume, err = v.StorageOS.Marshal()
	}

	if err != nil {
		return nil, err
	}

	if marshalledVolume == nil {
		return nil, serrors.InvalidInstanceErrorf(v, "empty volume definition")
	}

	if len(marshalledVolume.ExtraFields) == 0 {
		segments := []string{marshalledVolume.Type}
		segments = append(segments, marshalledVolume.Selector...)
		return json.Marshal(strings.Join(segments, ":"))
	}

	obj := marshalledVolume.ExtraFields
	obj["vol_type"] = marshalledVolume.Type
	if len(marshalledVolume.Selector) > 0 {
		obj["vol_id"] = strings.Join(marshalledVolume.Selector, ":")
	}

	return json.Marshal(obj)
}

func (s *HostPathVolume) Unmarshal(selector []string) error {
	if len(selector) > 2 || len(selector) == 0 {
		return serrors.InvalidValueErrorf(selector, "expected one or two selector segments for %s", VolumeTypeHostPath)
	}

	s.Path = selector[0]

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
			return serrors.InvalidValueErrorf(hostPathType, "invalid 'vol_type' selector for %s", VolumeTypeHostPath)
		}

		s.Type = hostPathType
	}

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
		return serrors.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeEmptyDir)
	}

	err := jsonutil.UnmarshalMap(obj, &source)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeEmptyDir)
	}

	v.EmptyDir = &source
	return nil
}

func (s EmptyDirVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeEmptyDir)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeEmptyDir,
		ExtraFields: obj,
	}, nil
}

func (s *GcePDVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (disk name) for %s", VolumeTypeGcePD)
	}
	s.PDName = selector[0]

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeGcePD)
	}

	return nil
}

func (s GcePDVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeGcePD)
	}

	if len(s.PDName) == 0 {
		return nil, serrors.InvalidInstanceErrorf(&s, "selector must contain disk name")
	}

	return &MarshalledVolume{
		Type:        VolumeTypeGcePD,
		Selector:    []string{s.PDName},
		ExtraFields: obj,
	}, nil
}

func (s *AwsEBSVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (ebs uuid) for %s", VolumeTypeAwsEBS)
	}
	s.VolumeID = selector[0]

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeAwsEBS)
	}

	return nil
}

func (s AwsEBSVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeAwsEBS)
	}

	if len(s.VolumeID) == 0 {
		return nil, serrors.InvalidInstanceErrorf(&s, "selector must contain ebs uuid")
	}

	return &MarshalledVolume{
		Type:        VolumeTypeAwsEBS,
		Selector:    []string{s.VolumeID},
		ExtraFields: obj,
	}, nil
}

func (s *AzureDiskVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeAzureDisk)
	}

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeAzureDisk)
	}

	return nil
}

func (s AzureDiskVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeAzureDisk)
	}

	if len(s.DiskName) == 0 {
		return nil, serrors.InvalidInstanceErrorf(&s, "disk_name is required for %s", VolumeTypeAzureDisk)
	}

	if len(s.DataDiskURI) == 0 {
		return nil, serrors.InvalidInstanceErrorf(&s, "disk_uri is required for %s", VolumeTypeAzureDisk)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeAzureDisk,
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalAzureFileVolume(selector []string) error {
	source := AzureFileVolume{}
	if len(selector) > 3 || len(selector) < 2 {
		return serrors.InvalidValueErrorf(selector, "expected two or three selector segments for %s", VolumeTypeAzureFile)
	}

	source.SecretName = selector[0]
	source.ShareName = selector[1]

	if len(selector) > 2 {
		switch selector[2] {
		case SelectorSegmentReadOnly:
			source.ReadOnly = true
		default:
			return serrors.InvalidValueErrorf(selector[2], "invalid selector segment for %s", VolumeTypeAzureFile)
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
		return serrors.InvalidValueErrorf(selector, "expected 0 selector segments for %s", VolumeTypeCephFS)
	}

	err := jsonutil.UnmarshalMap(obj, &source)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeCephFS)
	}

	v.CephFS = &source
	return nil
}

func (s CephFSVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeCephFS)
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
		return serrors.ContextualizeErrorf(err, "cephfs secret should be a string")
	}

	matches := fileOrRefRegexp.FindStringSubmatch(str)
	if len(matches) > 0 {
		if matches[1] == "file" {
			s.File = matches[2]
		} else {
			s.Ref = matches[2]
		}
	} else {
		return serrors.InvalidValueErrorf(string(data), "unrecognized format for cephfs secret")
	}

	return nil
}

func (s CephFSSecretFileOrRef) MarshalJSON() ([]byte, error) {
	if len(s.Ref) > 0 {
		return json.Marshal(fmt.Sprintf("ref:%s", s.Ref))
	}

	return json.Marshal(fmt.Sprintf("file:%s", s.File))
}

func (s *CinderVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 (volume ID) selector segment for %s", VolumeTypeCinder)
	}

	s.VolumeID = selector[0]

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeCinder)
	}

	return nil
}

func (s CinderVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeCinder)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeCinder,
		Selector:    []string{s.VolumeID},
		ExtraFields: obj,
	}, nil
}

func (s *FibreChannelVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected 0 selector segments for %s", VolumeTypeFibreChannel)
	}

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeFibreChannel)
	}

	return nil
}

func (s FibreChannelVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeFibreChannel)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeFibreChannel,
		ExtraFields: obj,
	}, nil
}

func (s *FlexVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (driver) for %s", VolumeTypeFlex)
	}
	s.Driver = selector[0]

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeFlex)
	}

	return nil
}

func (s FlexVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeFlex)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeFlex,
		Selector:    []string{s.Driver},
		ExtraFields: obj,
	}, nil
}

func (s *FlockerVolume) Unmarshal(selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected exactly one selector segment (dataset UUID) for %s", VolumeTypeFlocker)
	}

	s.DatasetUUID = selector[0]

	return nil
}

func (s FlockerVolume) Marshal() (*MarshalledVolume, error) {
	return &MarshalledVolume{
		Type:     VolumeTypeFlocker,
		Selector: []string{s.DatasetUUID},
	}, nil
}

func (s *GlusterfsVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeGlusterfs)
	}

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeGlusterfs)
	}

	return nil
}

func (s GlusterfsVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeGlusterfs)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeGlusterfs,
		ExtraFields: obj,
	}, nil
}

func (s *ISCSIVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeISCSI)
	}

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeISCSI)
	}

	return nil
}

func (s ISCSIVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeISCSI)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeISCSI,
		ExtraFields: obj,
	}, nil
}

func (s *NFSVolume) Unmarshal(selector []string) error {
	if len(selector) > 3 || len(selector) < 2 {
		return serrors.InvalidValueErrorf(selector, "expected two or three selector segments for %s", VolumeTypeNFS)
	}

	s.Server = selector[0]
	s.Path = selector[1]

	if len(selector) > 2 {
		switch selector[2] {
		case SelectorSegmentReadOnly:
			s.ReadOnly = true
		default:
			return serrors.InvalidValueErrorf(selector[2], "invalid selector segment for %s", VolumeTypeNFS)
		}
	}

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

func (s *PhotonPDVolume) Unmarshal(selector []string) error {
	if len(selector) > 2 || len(selector) < 1 {
		return serrors.InvalidValueErrorf(selector, "expected one or two selector segments for %s", VolumeTypePhotonPD)
	}

	s.PdID = selector[0]

	if len(selector) > 1 {
		s.FSType = selector[1]
	}

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

func (s *PortworxVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (volume ID) for %s", VolumeTypePortworx)
	}
	s.VolumeID = selector[0]

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypePortworx)
	}

	return nil
}

func (s PortworxVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypePortworx)
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
		return serrors.InvalidValueErrorf(selector, "expected one or two selector segments for %s", VolumeTypePVC)
	}

	source.ClaimName = selector[0]

	if len(selector) > 1 {
		switch selector[1] {
		case SelectorSegmentReadOnly:
			source.ReadOnly = true
		default:
			return serrors.InvalidValueErrorf(selector[2], "invalid selector segment for %s", VolumeTypePVC)
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

func (s *QuobyteVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (volume ID) for %s", VolumeTypeQuobyte)
	}
	s.Volume = selector[0]

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeQuobyte)
	}

	return nil
}

func (s QuobyteVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeQuobyte)
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
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (volume name) for %s", VolumeTypeScaleIO)
	}
	source.VolumeName = selector[0]

	err := jsonutil.UnmarshalMap(obj, &source)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeScaleIO)
	}

	v.ScaleIO = &source
	return nil
}

func (s ScaleIOVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeScaleIO)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeScaleIO,
		Selector:    []string{s.VolumeName},
		ExtraFields: obj,
	}, nil
}

func (s *VsphereVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (volume path) for %s", VolumeTypeVsphere)
	}
	s.VolumePath = selector[0]

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeVsphere)
	}

	return nil
}

func (s VsphereVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeVsphere)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeVsphere,
		Selector:    []string{s.VolumePath},
		ExtraFields: obj,
	}, nil
}

func (m *FileMode) UnmarshalJSON(data []byte) error {
	var i int32
	err := json.Unmarshal(data, &i)
	if err == nil {
		*m = FileMode(i)
		return nil
	}

	var str string
	err = json.Unmarshal(data, &str)
	mode, err := strconv.ParseInt(str, 8, 32)
	if err != nil {
		return serrors.InvalidValueErrorf(str, "file mode should be an octal integer, written either as string or number")
	}

	*m = FileMode(int32(mode))
	return nil
}

func (m FileMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(fmt.Sprintf("0%o", m))
}

func FileModePtr(m FileMode) *FileMode {
	return &m
}

var keyAndModeRegexp = regexp.MustCompile(`^(.*):(0[0-7][0-7][0-7])$`)

func (k *KeyAndMode) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return serrors.InvalidValueErrorf(string(data), "expected string for key:mode")
	}

	matches := keyAndModeRegexp.FindStringSubmatch(str)
	if len(matches) == 0 {
		k.Key = str
		return nil
	}

	k.Key = matches[1]

	// The regexp should ensure that this always succeeds.
	i, err := strconv.ParseInt(matches[2], 8, 32)
	if err != nil {
		glog.V(0).Info("KeyAndMode regexp is matching non-integer file modes.")
		return serrors.InvalidValueErrorf(str, "expected integer for file mode in key:mode")
	}
	mode := FileMode(int32(i))
	k.Mode = &mode
	return nil
}

func (k KeyAndMode) MarshalJSON() ([]byte, error) {
	if k.Mode != nil {
		return json.Marshal(fmt.Sprintf("%s:0%o", k.Key, *k.Mode))
	}

	return json.Marshal(k.Key)
}

func (v *Volume) UnmarshalConfigMapVolume(obj map[string]interface{}, selector []string) error {
	source := ConfigMapVolume{}
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (config name) for %s", VolumeTypeConfigMap)
	}
	source.Name = selector[0]

	err := jsonutil.UnmarshalMap(obj, &source)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeConfigMap)
	}

	v.ConfigMap = &source
	return nil
}

func (s ConfigMapVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeConfigMap)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeConfigMap,
		Selector:    []string{s.Name},
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalSecretVolume(obj map[string]interface{}, selector []string) error {
	source := SecretVolume{}
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (secret name) for %s", VolumeTypeSecret)
	}
	source.SecretName = selector[0]

	err := jsonutil.UnmarshalMap(obj, &source)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeSecret)
	}

	v.Secret = &source
	return nil
}

func (s SecretVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeSecret)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeSecret,
		Selector:    []string{s.SecretName},
		ExtraFields: obj,
	}, nil
}

func (s *ObjectFieldSelector) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return serrors.ContextualizeErrorf(err, "field selector should be written as a string")
	}

	segments := strings.Split(str, ":")
	if len(segments) > 2 {
		return serrors.InvalidValueErrorf(str, "field selector should contain one or two segments")
	}

	s.FieldPath = segments[0]
	if len(segments) > 1 {
		s.APIVersion = segments[1]
	}

	return nil
}

func (s ObjectFieldSelector) MarshalJSON() ([]byte, error) {
	if len(s.APIVersion) == 0 {
		b, err := json.Marshal(s.FieldPath)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "field selector path")
		}

		return b, nil
	}

	b, err := json.Marshal(fmt.Sprintf("%s:%s", s.FieldPath, s.APIVersion))
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "field selector")
	}

	return b, nil
}

func (s *VolumeResourceFieldSelector) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return serrors.ContextualizeErrorf(err, "resource selector should be written as a string")
	}

	segments := strings.Split(str, ":")
	if len(segments) > 3 || len(segments) < 2 {
		return serrors.InvalidValueErrorf(str, "resource selector should contain two or three segments")
	}

	s.ContainerName = segments[0]
	s.Resource = segments[1]
	if len(segments) > 2 {
		divisor, err := resource.ParseQuantity(segments[2])
		if err != nil {
			return serrors.ContextualizeErrorf(err, "resource selector divisor")
		}
		s.Divisor = divisor
	}

	return nil
}

func (s VolumeResourceFieldSelector) MarshalJSON() ([]byte, error) {
	if reflect.DeepEqual(s.Divisor, resource.Quantity{}) {
		b, err := json.Marshal(fmt.Sprintf("%s:%s", s.ContainerName, s.Resource))
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "resource selector")
		}

		return b, nil
	}

	b, err := json.Marshal(fmt.Sprintf("%s:%s:%s", s.ContainerName, s.Resource, s.Divisor.String()))
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "resource selector")
	}

	return b, nil
}

func (v *Volume) UnmarshalDownwardAPIVolume(obj map[string]interface{}, selector []string) error {
	source := DownwardAPIVolume{}
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeDownwardAPI)
	}

	err := jsonutil.UnmarshalMap(obj, &source)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeDownwardAPI)
	}

	v.DownwardAPI = &source
	return nil
}

func (s DownwardAPIVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeDownwardAPI)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeDownwardAPI,
		ExtraFields: obj,
	}, nil
}

func (v *Volume) UnmarshalProjectedVolume(obj map[string]interface{}, selector []string) error {
	source := ProjectedVolume{}
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeProjected)
	}

	err := jsonutil.UnmarshalMap(obj, &source)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeProjected)
	}

	v.Projected = &source
	return nil
}

func (s ProjectedVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeProjected)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeProjected,
		ExtraFields: obj,
	}, nil
}

func (p *VolumeProjection) UnmarshalJSON(data []byte) error {
	obj := map[string]interface{}{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return serrors.ContextualizeErrorf(err, "projected volume source item")
	}

	if _, ok := obj["secret"]; ok {
		p.Secret = &SecretProjection{}
		return json.Unmarshal(data, p.Secret)
	}
	if _, ok := obj["config"]; ok {
		p.ConfigMap = &ConfigMapProjection{}
		return json.Unmarshal(data, p.ConfigMap)
	}
	p.DownwardAPI = &DownwardAPIProjection{}
	return json.Unmarshal(data, p.DownwardAPI)
}

func (p VolumeProjection) MarshalJSON() ([]byte, error) {
	if p.Secret != nil {
		return json.Marshal(p.Secret)
	}

	if p.DownwardAPI != nil {
		return json.Marshal(p.DownwardAPI)
	}

	if p.ConfigMap != nil {
		return json.Marshal(p.ConfigMap)
	}

	return nil, serrors.InvalidInstanceErrorf(p, "empty volume projection")
}

func (v *Volume) UnmarshalGitVolume(obj map[string]interface{}, selector []string) error {
	source := GitVolume{}
	source.Repository = strings.Join(selector, ":")

	err := jsonutil.UnmarshalMap(obj, &source)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeGit)
	}

	v.Git = &source
	return nil
}

func (s GitVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeGit)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeGit,
		Selector:    []string{s.Repository},
		ExtraFields: obj,
	}, nil
}

func (s *RBDVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeRBD)
	}

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeRBD)
	}

	return nil
}

func (s RBDVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeRBD)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeRBD,
		ExtraFields: obj,
	}, nil
}

func (s *StorageOSVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (volume name) for %s", VolumeTypeStorageOS)
	}
	s.VolumeName = selector[0]

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeStorageOS)
	}

	return nil
}

func (s StorageOSVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeStorageOS)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeStorageOS,
		Selector:    []string{s.VolumeName},
		ExtraFields: obj,
	}, nil
}
