package converters

import (
	storagev1 "k8s.io/api/storage/v1"
	storagev1beta1 "k8s.io/api/storage/v1beta1"

	"github.com/ghodss/yaml"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_StorageClass_to_Kube_StorageClass(storageClass *types.StorageClassWrapper) (interface{}, error) {
	// Perform version-agnostic conversion into storage/v1 StorageClass.
	kubeStorageClass := Convert_Koki_StorageClass_to_Kube_storage_v1_StorageClass(storageClass)

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

func Convert_Koki_StorageClass_to_Kube_storage_v1_StorageClass(storageClass *types.StorageClassWrapper) *storagev1.StorageClass {
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

	if kokiStorageClass.Reclaim != nil {
		reclaimPolicy := revertReclaimPolicy(*kokiStorageClass.Reclaim)
		kubeStorageClass.ReclaimPolicy = &reclaimPolicy
	}

	return kubeStorageClass
}
