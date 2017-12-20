package converters

import (
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_StorageClass_to_Koki_StorageClass(kubeStorageClass runtime.Object) (*types.StorageClassWrapper, error) {
	groupVersionKind := kubeStorageClass.GetObjectKind().GroupVersionKind()
	groupVersionString := groupVersionKind.GroupVersion().String()
	groupVersionKind.Version = "v1"
	groupVersionKind.Group = "storage.k8s.io"
	kubeStorageClass.GetObjectKind().SetGroupVersionKind(groupVersionKind)

	// Serialize as v1
	b, err := yaml.Marshal(kubeStorageClass)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, kubeStorageClass, "couldn't serialize kube StorageClass after setting apiVersion to storage/v1")
	}

	// Deserialize the "generic" kube StorageClass
	genericStorageClass, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, string(b), "couldn't deserialize 'generic' kube StorageClass")
	}

	if genericStorageClass, ok := genericStorageClass.(*storagev1.StorageClass); ok {
		kokiWrapper, err := Convert_Kube_storage_v1_StorageClass_to_Koki_StorageClass(genericStorageClass)
		if err != nil {
			return nil, err
		}

		kokiStorageClass := &kokiWrapper.StorageClass

		kokiStorageClass.Version = groupVersionString

		return kokiWrapper, nil
	}

	return nil, serrors.InvalidInstanceErrorf(genericStorageClass, "didn't deserialize 'generic' kube Deployment as storage/v1.StorageClass")
}

func Convert_Kube_storage_v1_StorageClass_to_Koki_StorageClass(kubeStorageClass *storagev1.StorageClass) (*types.StorageClassWrapper, error) {
	var err error
	kokiStorageClass := &types.StorageClass{}

	kokiStorageClass.Name = kubeStorageClass.Name
	kokiStorageClass.Namespace = kubeStorageClass.Namespace
	kokiStorageClass.Version = kubeStorageClass.APIVersion
	kokiStorageClass.Cluster = kubeStorageClass.ClusterName
	kokiStorageClass.Labels = kubeStorageClass.Labels
	kokiStorageClass.Annotations = kubeStorageClass.Annotations

	kokiStorageClass.Provisioner = kubeStorageClass.Provisioner
	kokiStorageClass.Parameters = kubeStorageClass.Parameters
	kokiStorageClass.MountOptions = kubeStorageClass.MountOptions
	kokiStorageClass.AllowVolumeExpansion = kubeStorageClass.AllowVolumeExpansion
	kokiStorageClass.VolumeBindingMode, err = convertVolumeBindingMode(kubeStorageClass.VolumeBindingMode)
	if err != nil {
		return nil, err
	}

	if kubeStorageClass.ReclaimPolicy != nil {
		reclaimPolicy := convertReclaimPolicy(*kubeStorageClass.ReclaimPolicy)
		kokiStorageClass.Reclaim = &reclaimPolicy
	}

	return &types.StorageClassWrapper{
		StorageClass: *kokiStorageClass,
	}, nil
}

func convertVolumeBindingMode(mode *storagev1.VolumeBindingMode) (*types.VolumeBindingMode, error) {
	if mode == nil {
		return nil, nil
	}

	var newmode types.VolumeBindingMode
	switch *mode {
	case storagev1.VolumeBindingImmediate:
		newmode = types.VolumeBindingImmediate
	case storagev1.VolumeBindingWaitForFirstConsumer:
		newmode = types.VolumeBindingWaitForFirstConsumer
	default:
		return nil, serrors.InvalidInstanceError(mode)
	}
	return &newmode, nil
}
