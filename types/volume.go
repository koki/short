package types

import (
	"k8s.io/api/core/v1"
)

type VolumeWrapper struct {
	Volume v1.Volume `json:"volume"`
}
