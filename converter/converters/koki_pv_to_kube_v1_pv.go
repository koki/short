package converters

import (
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Koki_PersistentVolume_to_Kube_v1_PersistentVolume(pv *types.PersistentVolumeWrapper) (*v1.PersistentVolume, error) {
	var err error
	kubePV := &v1.PersistentVolume{}
	kokiPV := pv.PersistentVolume

	kubePV.Name = kokiPV.Name
	kubePV.Namespace = kokiPV.Namespace
	if len(kokiPV.Version) == 0 {
		kubePV.APIVersion = "v1"
	} else {
		kubePV.APIVersion = kokiPV.Version
	}
	kubePV.Kind = "PersistentVolume"
	kubePV.ClusterName = kokiPV.Cluster
	kubePV.Labels = kokiPV.Labels
	kubePV.Annotations = kokiPV.Annotations

	kubeSpec := &kubePV.Spec
	kubeSpec.Capacity = revertCapacity(kokiPV.Storage)
	kubeSpec.PersistentVolumeSource, err = revertPersistentVolumeSource(kokiPV.PersistentVolumeSource)
	if err != nil {
		return nil, err
	}
	if kokiPV.AccessModes != nil {
		kubeSpec.AccessModes = kokiPV.AccessModes.Modes
	}
	kubeSpec.ClaimRef = kokiPV.Claim
	kubeSpec.PersistentVolumeReclaimPolicy = revertReclaimPolicy(kokiPV.ReclaimPolicy)
	kubeSpec.StorageClassName = kokiPV.StorageClass
	if len(kokiPV.MountOptions) > 0 {
		kubeSpec.MountOptions = strings.Split(kokiPV.MountOptions, ",")
	}

	if kokiPV.Status != nil {
		kubePV.Status = *kokiPV.Status
	}

	return kubePV, nil
}

func revertPersistentVolumeSource(kokiSource types.PersistentVolumeSource) (v1.PersistentVolumeSource, error) {
	if kokiSource.GcePD != nil {
		return v1.PersistentVolumeSource{
			GCEPersistentDisk: revertGcePDVolume(kokiSource.GcePD),
		}, nil
	}

	return v1.PersistentVolumeSource{}, util.InvalidInstanceErrorf(kokiSource, "didn't find any supported volume source")
}

func revertReclaimPolicy(kokiPolicy types.PersistentVolumeReclaimPolicy) v1.PersistentVolumeReclaimPolicy {
	return v1.PersistentVolumeReclaimPolicy(strings.Title(string(kokiPolicy)))
}

func revertCapacity(kokiStorage *resource.Quantity) v1.ResourceList {
	if kokiStorage == nil {
		return nil
	}

	kubeCapacity := v1.ResourceList{}
	kubeCapacity[v1.ResourceStorage] = *kokiStorage

	return kubeCapacity
}
