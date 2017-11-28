package converters

import (
	"reflect"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Kube_v1_PersistentVolume_to_Koki_PersistentVolume(kubePV *v1.PersistentVolume) (*types.PersistentVolumeWrapper, error) {
	var err error
	kokiPV := &types.PersistentVolume{}

	kokiPV.Name = kubePV.Name
	kokiPV.Namespace = kubePV.Namespace
	kokiPV.Version = kubePV.APIVersion
	kokiPV.Cluster = kubePV.ClusterName
	kokiPV.Labels = kubePV.Labels
	kokiPV.Annotations = kubePV.Annotations

	kubeSpec := kubePV.Spec
	kokiPV.Storage, err = convertCapacity(kubeSpec.Capacity)
	if err != nil {
		return nil, err
	}

	kokiPV.PersistentVolumeSource, err = convertPersistentVolumeSource(kubeSpec.PersistentVolumeSource)
	if err != nil {
		return nil, err
	}
	if len(kubeSpec.AccessModes) > 0 {
		kokiPV.AccessModes = &types.AccessModes{
			Modes: kubeSpec.AccessModes,
		}
	}
	kokiPV.Claim = kubeSpec.ClaimRef
	kokiPV.ReclaimPolicy = convertReclaimPolicy(kubeSpec.PersistentVolumeReclaimPolicy)
	kokiPV.StorageClass = kubeSpec.StorageClassName
	if len(kubeSpec.MountOptions) > 0 {
		kokiPV.MountOptions = strings.Join(kubeSpec.MountOptions, ",")
	}

	if !reflect.DeepEqual(kubePV.Status, v1.PersistentVolumeStatus{}) {
		kokiPV.Status = &kubePV.Status
	}

	return &types.PersistentVolumeWrapper{
		PersistentVolume: *kokiPV,
	}, nil
}

func convertPersistentVolumeSource(kubeSource v1.PersistentVolumeSource) (types.PersistentVolumeSource, error) {
	if kubeSource.GCEPersistentDisk != nil {
		return types.PersistentVolumeSource{
			GcePD: convertGcePDVolume(kubeSource.GCEPersistentDisk),
		}, nil
	}
	if kubeSource.AWSElasticBlockStore != nil {
		return types.PersistentVolumeSource{
			AwsEBS: convertAwsEBSVolume(kubeSource.AWSElasticBlockStore),
		}, nil
	}
	if kubeSource.HostPath != nil {
		source, err := convertHostPathVolume(kubeSource.HostPath)
		if err != nil {
			return types.PersistentVolumeSource{}, err
		}
		return types.PersistentVolumeSource{
			HostPath: source,
		}, nil
	}

	return types.PersistentVolumeSource{}, util.InvalidInstanceErrorf(kubeSource, "didn't find any supported volume source")
}

func convertReclaimPolicy(kubePolicy v1.PersistentVolumeReclaimPolicy) types.PersistentVolumeReclaimPolicy {
	return types.PersistentVolumeReclaimPolicy(strings.ToLower(string(kubePolicy)))
}

func convertCapacity(kubeCapacity v1.ResourceList) (*resource.Quantity, error) {
	if len(kubeCapacity) == 0 {
		return nil, nil
	}

	for res, quantity := range kubeCapacity {
		if res == v1.ResourceStorage {
			return &quantity, nil
		}
	}

	return nil, util.InvalidInstanceErrorf(kubeCapacity, "only supports Storage resource")
}
