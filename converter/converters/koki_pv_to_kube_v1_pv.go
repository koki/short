package converters

import (
	"strings"

	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	"k8s.io/apimachinery/pkg/api/resource"
)

func Convert_Koki_PersistentVolume_to_Kube_v1_PersistentVolume(pv *types.PersistentVolumeWrapper) (*v1.PersistentVolume, error) {
	kubePV := &v1.PersistentVolume{}
	kokiPV := pv.PersistentVolume

	kubePV.Name = kokiPV.Name
	kubePV.Namespace = kokiPV.Namespace
	kubePV.APIVersion = kokiPV.Version
	kubePV.Kind = "PersistentVolume"
	kubePV.ClusterName = kokiPV.Cluster
	kubePV.Labels = kokiPV.Labels
	kubePV.Annotations = kokiPV.Annotations

	kubeSpec := &kubePV.Spec
	kubeSpec.Capacity = revertCapacity(kokiPV.Storage)
	kubeSpec.PersistentVolumeSource = kokiPV.PersistentVolumeSource.VolumeSource
	if kokiPV.AccessModes != nil {
		kubeSpec.AccessModes = kokiPV.AccessModes.Modes
	}
	kubeSpec.ClaimRef = kokiPV.Claim
	kubeSpec.PersistentVolumeReclaimPolicy = kokiPV.ReclaimPolicy
	kubeSpec.StorageClassName = kokiPV.StorageClass
	if len(kokiPV.MountOptions) > 0 {
		kubeSpec.MountOptions = strings.Split(kokiPV.MountOptions, ",")
	}

	if kokiPV.Status != nil {
		kubePV.Status = *kokiPV.Status
	}

	return kubePV, nil
}

func revertCapacity(kokiStorage *resource.Quantity) v1.ResourceList {
	if kokiStorage == nil {
		return nil
	}

	kubeCapacity := v1.ResourceList{}
	kubeCapacity[v1.ResourceStorage] = *kokiStorage

	return kubeCapacity
}
