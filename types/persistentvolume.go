package types

import (
	"fmt"
	"regexp"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/koki/json"
	"github.com/koki/json/jsonutil"
	serrors "github.com/koki/structurederrors"
)

type PersistentVolumeWrapper struct {
	PersistentVolume PersistentVolume `json:"persistent_volume"`
}

type PersistentVolume struct {
	PersistentVolumeMeta
	PersistentVolumeSource
}

type PersistentVolumeMeta struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Storage       *resource.Quantity            `json:"storage,omitempty"`
	AccessModes   *AccessModes                  `json:"modes,omitempty"`
	Claim         *v1.ObjectReference           `json:"claim,omitempty"`
	ReclaimPolicy PersistentVolumeReclaimPolicy `json:"reclaim,omitempty"`
	StorageClass  string                        `json:"storage_class,omitempty"`

	// comma-separated list of options
	MountOptions string `json:"mount_opts,omitempty"`

	PersistentVolumeStatus `json:",inline"`
}

type PersistentVolumeStatus struct {
	Phase   PersistentVolumePhase `json:"status,omitempty"`
	Message string                `json:"status_message,omitempty"`
	Reason  string                `json:"status_reason,omitempty"`
}

type PersistentVolumePhase string

const (
	VolumePending   PersistentVolumePhase = "pending"
	VolumeAvailable PersistentVolumePhase = "available"
	VolumeBound     PersistentVolumePhase = "bound"
	VolumeReleased  PersistentVolumePhase = "released"
	VolumeFailed    PersistentVolumePhase = "failed"
)

type PersistentVolumeReclaimPolicy string

const (
	PersistentVolumeReclaimRecycle PersistentVolumeReclaimPolicy = "recycle"
	PersistentVolumeReclaimDelete  PersistentVolumeReclaimPolicy = "delete"
	PersistentVolumeReclaimRetain  PersistentVolumeReclaimPolicy = "retain"
)

type PersistentVolumeSource struct {
	GcePD        *GcePDVolume
	AwsEBS       *AwsEBSVolume
	HostPath     *HostPathVolume
	Glusterfs    *GlusterfsVolume
	NFS          *NFSVolume
	ISCSI        *ISCSIPersistentVolume
	Cinder       *CinderVolume
	FibreChannel *FibreChannelVolume
	Flocker      *FlockerVolume
	Flex         *FlexVolume
	Vsphere      *VsphereVolume
	Quobyte      *QuobyteVolume
	AzureDisk    *AzureDiskVolume
	PhotonPD     *PhotonPDVolume
	Portworx     *PortworxVolume
	RBD          *RBDPersistentVolume
	CephFS       *CephFSPersistentVolume
	AzureFile    *AzureFilePersistentVolume
	ScaleIO      *ScaleIOPersistentVolume
	Local        *LocalVolume
	StorageOS    *StorageOSPersistentVolume
	CSI          *CSIPersistentVolume
}

const (
	VolumeTypeLocal = "local"
	VolumeTypeCSI   = "csi"
)

type ISCSIPersistentVolume struct {
	TargetPortal   string   `json:"target_portal"`
	IQN            string   `json:"iqn"`
	Lun            int32    `json:"lun"`
	ISCSIInterface string   `json:"iscsi_interface,omitempty"`
	FSType         string   `json:"fs,omitempty"`
	ReadOnly       bool     `json:"ro,omitempty"`
	Portals        []string `json:"portals,omitempty"`
	// TODO: should this actually be "chap_auth"?
	DiscoveryCHAPAuth bool             `json:"chap_discovery,omitempty"`
	SessionCHAPAuth   bool             `json:"chap_session,omitempty"`
	SecretRef         *SecretReference `json:"secret,omitempty"`
	// NOTE: InitiatorName is a pointer in k8s
	InitiatorName string `json:"initiator,omitempty"`
}

type RBDPersistentVolume struct {
	CephMonitors []string         `json:"monitors"`
	RBDImage     string           `json:"image"`
	FSType       string           `json:"fs,omitempty"`
	RBDPool      string           `json:"pool,omitempty"`
	RadosUser    string           `json:"user,omitempty"`
	Keyring      string           `json:"keyring,omitempty"`
	SecretRef    *SecretReference `json:"secret,omitempty"`
	ReadOnly     bool             `json:"ro,omitempty"`
}

type SecretReference struct {
	Name      string `json:"-"`
	Namespace string `json:"-"`
}

type CephFSPersistentVolume struct {
	Monitors        []string                         `json:"monitors"`
	Path            string                           `json:"path, omitempty"`
	User            string                           `json:"user,omitempty"`
	SecretFileOrRef *CephFSPersistentSecretFileOrRef `json:"secret,omitempty"`
	ReadOnly        bool                             `json:"ro,omitempty"`
}

type CephFSPersistentSecretFileOrRef struct {
	File string           `json:"-"`
	Ref  *SecretReference `json:"-"`
}

type AzureFilePersistentVolume struct {
	Secret    SecretReference `json:"secret"`
	ShareName string          `json:"share"`
	ReadOnly  bool            `json:"ro,omitempty"`
}

type ScaleIOPersistentVolume struct {
	Gateway          string             `json:"gateway"`
	System           string             `json:"system"`
	SecretRef        SecretReference    `json:"secret"`
	SSLEnabled       bool               `json:"ssl,omitempty"`
	ProtectionDomain string             `json:"protection_domain,omitempty"`
	StoragePool      string             `json:"storage_pool,omitempty"`
	StorageMode      ScaleIOStorageMode `json:"storage_mode,omitempty"`
	VolumeName       string             `json:"-"`
	FSType           string             `json:"fs,omitempty"`
	ReadOnly         bool               `json:"ro,omitempty"`
}

type LocalVolume struct {
	Path string `json:"path"`
}

type StorageOSPersistentVolume struct {
	VolumeName      string           `json:"-"`
	VolumeNamespace string           `json:"vol_ns,omitempty"`
	FSType          string           `json:"fs,omitempty"`
	ReadOnly        bool             `json:"ro,omitempty"`
	SecretRef       *SecretReference `json:"secret,omitempty"`
}

type CSIPersistentVolume struct {
	Driver       string `json:"driver"`
	VolumeHandle string `json:"-"`
	ReadOnly     bool   `json:"ro,omitempty"`
}

// comma-separated list of modes
type AccessModes struct {
	Modes []v1.PersistentVolumeAccessMode
}

func (a *AccessModes) ToString() (string, error) {
	if a == nil {
		return "", nil
	}

	if len(a.Modes) == 0 {
		return "", nil
	}

	modes := make([]string, len(a.Modes))
	for i, mode := range a.Modes {
		switch mode {
		case v1.ReadOnlyMany:
			modes[i] = "ro"
		case v1.ReadWriteMany:
			modes[i] = "rw"
		case v1.ReadWriteOnce:
			modes[i] = "rw-once"
		default:
			return "", serrors.InvalidInstanceError(mode)
		}
	}

	return strings.Join(modes, ","), nil
}

func (a *AccessModes) InitFromString(s string) error {
	modes := strings.Split(s, ",")
	if len(modes) == 0 {
		a.Modes = nil
		return nil
	}

	a.Modes = make([]v1.PersistentVolumeAccessMode, len(modes))
	for i, mode := range modes {
		switch mode {
		case "ro":
			a.Modes[i] = v1.ReadOnlyMany
		case "rw":
			a.Modes[i] = v1.ReadWriteMany
		case "rw-once":
			a.Modes[i] = v1.ReadWriteOnce
		default:
			return serrors.InvalidValueErrorf(a, "couldn't parse (%s)", s)
		}
	}

	return nil
}

func (a AccessModes) MarshalJSON() ([]byte, error) {
	str, err := a.ToString()
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(&str)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, a, "marshalling to JSON")
	}

	return b, nil
}

func (a *AccessModes) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return serrors.InvalidInstanceContextErrorf(err, a, "unmarshalling from JSON")
	}

	return a.InitFromString(str)
}

func (v *PersistentVolume) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &v.PersistentVolumeSource)
	if err != nil {
		return serrors.InvalidValueForTypeContextErrorf(err, string(data), v, "unmarshalling volume source from JSON")
	}

	err = json.Unmarshal(data, &v.PersistentVolumeMeta)
	if err != nil {
		return serrors.InvalidValueForTypeContextErrorf(err, string(data), v, "unmarshalling metadata from JSON")
	}

	return nil
}

func (v PersistentVolume) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(v.PersistentVolumeMeta)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, v, "marshalling metadata to JSON")
	}

	bb, err := json.Marshal(v.PersistentVolumeSource)
	if err != nil {
		return nil, err
	}

	metaObj := map[string]interface{}{}
	err = json.Unmarshal(b, &metaObj)
	if err != nil {
		return nil, serrors.InvalidValueForTypeContextErrorf(err, string(b), v.PersistentVolumeMeta, "converting metadata to dictionary")
	}

	sourceObj := map[string]interface{}{}
	err = json.Unmarshal(bb, &sourceObj)
	if err != nil {
		return nil, serrors.InvalidValueForTypeContextErrorf(err, string(bb), v.PersistentVolumeSource, "converting volume source to dictionary")
	}

	// Merge metadata with volume-source
	for key, val := range metaObj {
		sourceObj[key] = val
	}

	result, err := json.Marshal(sourceObj)
	if err != nil {
		return nil, serrors.InvalidValueForTypeContextErrorf(err, result, v, "marshalling PersistentVolume-as-dictionary to JSON")
	}

	return result, nil
}

func (v *PersistentVolumeSource) UnmarshalJSON(data []byte) error {
	var err error
	obj := map[string]interface{}{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return serrors.InvalidValueErrorf(string(data), "expected dictionary for persistent volume")
	}

	var selector []string
	if val, ok := obj["vol_id"]; ok {
		if volName, ok := val.(string); ok {
			selector = strings.Split(volName, ":")
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

func (v *PersistentVolumeSource) Unmarshal(obj map[string]interface{}, volType string, selector []string) error {
	switch volType {
	case VolumeTypeGcePD:
		v.GcePD = &GcePDVolume{}
		return v.GcePD.Unmarshal(obj, selector)
	case VolumeTypeAwsEBS:
		v.AwsEBS = &AwsEBSVolume{}
		return v.AwsEBS.Unmarshal(obj, selector)
	case VolumeTypeHostPath:
		v.HostPath = &HostPathVolume{}
		return v.HostPath.Unmarshal(selector)
	case VolumeTypeGlusterfs:
		v.Glusterfs = &GlusterfsVolume{}
		return v.Glusterfs.Unmarshal(obj, selector)
	case VolumeTypeNFS:
		v.NFS = &NFSVolume{}
		return v.NFS.Unmarshal(selector)
	case VolumeTypeISCSI:
		v.ISCSI = &ISCSIPersistentVolume{}
		return v.ISCSI.Unmarshal(obj, selector)
	case VolumeTypeCinder:
		v.Cinder = &CinderVolume{}
		return v.Cinder.Unmarshal(obj, selector)
	case VolumeTypeFibreChannel:
		v.FibreChannel = &FibreChannelVolume{}
		return v.FibreChannel.Unmarshal(obj, selector)
	case VolumeTypeFlocker:
		v.Flocker = &FlockerVolume{}
		return v.Flocker.Unmarshal(selector)
	case VolumeTypeFlex:
		v.Flex = &FlexVolume{}
		return v.Flex.Unmarshal(obj, selector)
	case VolumeTypeVsphere:
		v.Vsphere = &VsphereVolume{}
		return v.Vsphere.Unmarshal(obj, selector)
	case VolumeTypeQuobyte:
		v.Quobyte = &QuobyteVolume{}
		return v.Quobyte.Unmarshal(obj, selector)
	case VolumeTypeAzureDisk:
		v.AzureDisk = &AzureDiskVolume{}
		return v.AzureDisk.Unmarshal(obj, selector)
	case VolumeTypePhotonPD:
		v.PhotonPD = &PhotonPDVolume{}
		return v.PhotonPD.Unmarshal(selector)
	case VolumeTypePortworx:
		v.Portworx = &PortworxVolume{}
		return v.Portworx.Unmarshal(obj, selector)
	case VolumeTypeRBD:
		v.RBD = &RBDPersistentVolume{}
		return v.RBD.Unmarshal(obj, selector)
	case VolumeTypeCephFS:
		v.CephFS = &CephFSPersistentVolume{}
		return v.CephFS.Unmarshal(obj, selector)
	case VolumeTypeAzureFile:
		v.AzureFile = &AzureFilePersistentVolume{}
		return v.AzureFile.Unmarshal(obj, selector)
	case VolumeTypeScaleIO:
		v.ScaleIO = &ScaleIOPersistentVolume{}
		return v.ScaleIO.Unmarshal(obj, selector)
	case VolumeTypeLocal:
		v.Local = &LocalVolume{}
		return v.Local.Unmarshal(obj, selector)
	case VolumeTypeStorageOS:
		v.StorageOS = &StorageOSPersistentVolume{}
		return v.StorageOS.Unmarshal(obj, selector)
	case VolumeTypeCSI:
		v.CSI = &CSIPersistentVolume{}
		return v.CSI.Unmarshal(obj, selector)
	default:
		return serrors.InvalidValueErrorf(volType, "unsupported volume type (%s)", volType)
	}
}

func (v PersistentVolumeSource) MarshalJSON() ([]byte, error) {
	var marshalledVolume *MarshalledVolume
	var err error
	if v.GcePD != nil {
		marshalledVolume, err = v.GcePD.Marshal()
	}
	if v.AwsEBS != nil {
		marshalledVolume, err = v.AwsEBS.Marshal()
	}
	if v.HostPath != nil {
		marshalledVolume, err = v.HostPath.Marshal()
	}
	if v.Glusterfs != nil {
		marshalledVolume, err = v.Glusterfs.Marshal()
	}
	if v.NFS != nil {
		marshalledVolume, err = v.NFS.Marshal()
	}
	if v.ISCSI != nil {
		marshalledVolume, err = v.ISCSI.Marshal()
	}
	if v.Cinder != nil {
		marshalledVolume, err = v.Cinder.Marshal()
	}
	if v.FibreChannel != nil {
		marshalledVolume, err = v.FibreChannel.Marshal()
	}
	if v.Flocker != nil {
		marshalledVolume, err = v.Flocker.Marshal()
	}
	if v.Flex != nil {
		marshalledVolume, err = v.Flex.Marshal()
	}
	if v.Vsphere != nil {
		marshalledVolume, err = v.Vsphere.Marshal()
	}
	if v.Quobyte != nil {
		marshalledVolume, err = v.Quobyte.Marshal()
	}
	if v.AzureDisk != nil {
		marshalledVolume, err = v.AzureDisk.Marshal()
	}
	if v.PhotonPD != nil {
		marshalledVolume, err = v.PhotonPD.Marshal()
	}
	if v.Portworx != nil {
		marshalledVolume, err = v.Portworx.Marshal()
	}
	if v.RBD != nil {
		marshalledVolume, err = v.RBD.Marshal()
	}
	if v.CephFS != nil {
		marshalledVolume, err = v.CephFS.Marshal()
	}
	if v.AzureFile != nil {
		marshalledVolume, err = v.AzureFile.Marshal()
	}
	if v.ScaleIO != nil {
		marshalledVolume, err = v.ScaleIO.Marshal()
	}
	if v.Local != nil {
		marshalledVolume, err = v.Local.Marshal()
	}
	if v.StorageOS != nil {
		marshalledVolume, err = v.StorageOS.Marshal()
	}
	if v.CSI != nil {
		marshalledVolume, err = v.CSI.Marshal()
	}

	if err != nil {
		return nil, err
	}

	if marshalledVolume == nil {
		return nil, serrors.InvalidInstanceErrorf(v, "empty volume definition")
	}

	if len(marshalledVolume.ExtraFields) == 0 {
		marshalledVolume.ExtraFields = map[string]interface{}{}
	}

	obj := marshalledVolume.ExtraFields
	obj["vol_type"] = marshalledVolume.Type
	if len(marshalledVolume.Selector) > 0 {
		obj["vol_id"] = strings.Join(marshalledVolume.Selector, ":")
	}

	return json.Marshal(obj)
}

var secretRefRegexp = regexp.MustCompile(`^(.*):([^:]*)`)

func (s *SecretReference) UnmarshalString(str string) {
	matches := secretRefRegexp.FindStringSubmatch(str)
	if len(matches) > 0 {
		s.Namespace = matches[1]
		s.Name = matches[2]
	} else {
		s.Name = str
	}
}

func (s *SecretReference) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return serrors.ContextualizeErrorf(err, "secret ref should be a string")
	}
	s.UnmarshalString(str)

	return nil
}

func (s SecretReference) MarshalString() string {
	if len(s.Namespace) > 0 {
		return fmt.Sprintf("%s:%s", s.Namespace, s.Name)
	}

	return s.Name
}

func (s SecretReference) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.MarshalString())
}

func (s *ISCSIPersistentVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeISCSI)
	}

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeISCSI)
	}

	return nil
}

func (s ISCSIPersistentVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeISCSI)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeISCSI,
		ExtraFields: obj,
	}, nil
}

func (s *RBDPersistentVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeRBD)
	}

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeRBD)
	}

	return nil
}

func (s RBDPersistentVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeRBD)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeRBD,
		ExtraFields: obj,
	}, nil
}

func (s *CephFSPersistentSecretFileOrRef) UnmarshalJSON(data []byte) error {
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
			s.Ref = &SecretReference{}
			s.Ref.UnmarshalString(matches[2])
		}
	} else {
		return serrors.InvalidValueErrorf(string(data), "unrecognized format for cephfs secret")
	}

	return nil
}

func (s CephFSPersistentSecretFileOrRef) MarshalJSON() ([]byte, error) {
	if s.Ref != nil {
		return json.Marshal(fmt.Sprintf("ref:%s", s.Ref.MarshalString()))
	}

	return json.Marshal(fmt.Sprintf("file:%s", s.File))
}

func (s *CephFSPersistentVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected 0 selector segments for %s", VolumeTypeCephFS)
	}

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeCephFS)
	}

	return nil
}

func (s CephFSPersistentVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeCephFS)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeCephFS,
		ExtraFields: obj,
	}, nil
}

func (s *AzureFilePersistentVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected 0 selector segments for %s", VolumeTypeAzureFile)
	}

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeAzureFile)
	}

	return nil
}

func (s AzureFilePersistentVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeAzureFile)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeAzureFile,
		ExtraFields: obj,
	}, nil
}

func (s *ScaleIOPersistentVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (volume name) for %s", VolumeTypeScaleIO)
	}
	s.VolumeName = selector[0]

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeScaleIO)
	}

	return nil
}

func (s ScaleIOPersistentVolume) Marshal() (*MarshalledVolume, error) {
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

func (s *LocalVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 0 {
		return serrors.InvalidValueErrorf(selector, "expected zero selector segments for %s", VolumeTypeLocal)
	}

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeLocal)
	}

	return nil
}

func (s LocalVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeLocal)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeLocal,
		ExtraFields: obj,
	}, nil
}

func (s *StorageOSPersistentVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
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

func (s StorageOSPersistentVolume) Marshal() (*MarshalledVolume, error) {
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

func (s *CSIPersistentVolume) Unmarshal(obj map[string]interface{}, selector []string) error {
	if len(selector) != 1 {
		return serrors.InvalidValueErrorf(selector, "expected 1 selector segment (volume handle) for %s", VolumeTypeCSI)
	}
	s.VolumeHandle = selector[0]

	err := jsonutil.UnmarshalMap(obj, &s)
	if err != nil {
		return serrors.ContextualizeErrorf(err, VolumeTypeCSI)
	}

	return nil
}

func (s CSIPersistentVolume) Marshal() (*MarshalledVolume, error) {
	obj, err := jsonutil.MarshalMap(&s)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, VolumeTypeCSI)
	}

	return &MarshalledVolume{
		Type:        VolumeTypeCSI,
		Selector:    []string{s.VolumeHandle},
		ExtraFields: obj,
	}, nil
}
