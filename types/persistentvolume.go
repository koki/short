package types

import (
	"encoding/json"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/koki/short/util"
)

type PersistentVolumeWrapper struct {
	PersistentVolume PersistentVolume `json:"persistentVolume"`
}

type PersistentVolume struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Storage                   *resource.Quantity `json:"storage,omitempty"`
	v1.PersistentVolumeSource `json:",inline"`
	AccessModes               *AccessModes                     `json:"accessModes,omitempty"`
	Claim                     *v1.ObjectReference              `json:"claim,omitempty"`
	ReclaimPolicy             v1.PersistentVolumeReclaimPolicy `json:"reclaimPolicy,omitempty"`
	StorageClass              string                           `json:"storageClass,omitempty"`

	// comma-separated list of options
	MountOptions string `json:"mountOptions,omitempty" protobuf:"bytes,7,opt,name=mountOptions"`

	Status *v1.PersistentVolumeStatus `json:"status,omitempty"`
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
			modes[i] = "rw_once"
		default:
			return "", util.TypeValueErrorf(mode, mode)
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
		case "rw_once":
			a.Modes[i] = v1.ReadWriteOnce
		default:
			return util.TypeValueErrorf(a, s)
		}
	}

	return nil
}

func (a AccessModes) MarshalJSON() ([]byte, error) {
	str, err := a.ToString()
	if err != nil {
		return nil, err
	}

	return json.Marshal(&str)
}

func (a *AccessModes) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	return a.InitFromString(str)
}
