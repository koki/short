package converters

import (
	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_StorageClass_to_Kube_StorageClass(storageClass *types.StorageClassWrapper) (interface{}, error) {
	// Perform version-agnostic conversion into storage/v1 StorageClass.
	kubeStorageClass, err := Convert_Koki_StorageClass_to_Kube_storage_v1_StorageClass(storageClass)
	if err != nil {
		return nil, err
	}

	// Serialize the "generic" kube StorageClass.
	b, err := yaml.Marshal(kubeStorageClass)
	if err != nil {
		return nil, serrors.InvalidValueContextErrorf(err, kubeStorageClass, "couldn't serialize 'generic' kube StorageClass")
	}

	// Deserialize a versioned kube StorageClass using its apiVersion.
	versionedStorageClass, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, err
	}

	switch versionedStorageClass := versionedStorageClass.(type) {
	case *storagev1.StorageClass:
		// Perform storage/v1beta1 initialization here.
	case *storagev1beta1.StorageClass:
		// Perform storage/v1beta1 initialization here.
	default:
		return nil, serrors.TypeErrorf(versionedStorageClass, "deserialized the manifest, but not as a supported kube StorageClass")
	}

	return versionedStorageClass, nil
}

func Convert_Koki_StorageClass_to_Kube_storage_v1_StorageClass(storageClass *types.StorageClassWrapper) (*storagev1.StorageClass, error) {
	var err error
	kubeStorageClass := &storagev1.StorageClass{}
	kokiStorageClass := &storageClass.StorageClass

	kubeStorageClass.Name = kokiStorageClass.Name
	kubeStorageClass.Namespace = kokiStorageClass.Namespace
	kubeStorageClass.APIVersion = kokiStorageClass.Version
	kubeStorageClass.Kind = "StorageClass"
	kubeStorageClass.ClusterName = kokiStorageClass.Cluster
	kubeStorageClass.Labels = kokiStorageClass.Labels
	kubeStorageClass.Annotations = kokiStorageClass.Annotations

	kubeStorageClass.Provisioner = kokiStorageClass.Provisioner
	kubeStorageClass.Parameters = kokiStorageClass.Parameters

	kubeStorageClass.MountOptions = kokiStorageClass.MountOptions
	kubeStorageClass.AllowVolumeExpansion = kokiStorageClass.AllowVolumeExpansion
	kubeStorageClass.VolumeBindingMode, err = revertVolumeBindingMode(kokiStorageClass.VolumeBindingMode)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "binding_mode")
	}

	if kokiStorageClass.Reclaim != nil {
		reclaimPolicy := revertReclaimPolicy(*kokiStorageClass.Reclaim)
		kubeStorageClass.ReclaimPolicy = &reclaimPolicy
	}

	return kubeStorageClass, nil
}

func revertVolumeBindingMode(mode *types.VolumeBindingMode) (*storagev1.VolumeBindingMode, error) {
	if mode == nil {
		return nil, nil
	}

	var newmode storagev1.VolumeBindingMode
	switch *mode {
	case types.VolumeBindingImmediate:
		newmode = storagev1.VolumeBindingImmediate
	case types.VolumeBindingWaitForFirstConsumer:
		newmode = storagev1.VolumeBindingWaitForFirstConsumer
	default:
		return nil, serrors.InvalidInstanceError(mode)
	}
	return &newmode, nil
}
