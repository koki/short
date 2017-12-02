package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
)

func Convert_Kube_v1_ConfigMap_to_Koki_ConfigMap(kubeConfigMap *v1.ConfigMap) (*types.ConfigMapWrapper, error) {
	kokiConfigMapWrapper := &types.ConfigMapWrapper{}
	kokiConfigMap := types.ConfigMap{}

	kokiConfigMap.Name = kubeConfigMap.Name
	kokiConfigMap.Namespace = kubeConfigMap.Namespace
	kokiConfigMap.Version = kubeConfigMap.APIVersion
	kokiConfigMap.Cluster = kubeConfigMap.ClusterName
	kokiConfigMap.Labels = kubeConfigMap.Labels
	kokiConfigMap.Annotations = kubeConfigMap.Annotations

	kokiConfigMap.Data = kubeConfigMap.Data

	kokiConfigMapWrapper.ConfigMap = kokiConfigMap

	return kokiConfigMapWrapper, nil
}

func Convert_Koki_ConfigMap_to_Kube_v1_ConfigMap(kokiConfigMapWrapper *types.ConfigMapWrapper) (*v1.ConfigMap, error) {
	kubeConfigMap := &v1.ConfigMap{}
	kokiConfigMap := kokiConfigMapWrapper.ConfigMap

	kubeConfigMap.Name = kokiConfigMap.Name
	kubeConfigMap.Namespace = kokiConfigMap.Namespace
	kubeConfigMap.APIVersion = kokiConfigMap.Version
	kubeConfigMap.ClusterName = kokiConfigMap.Cluster
	kubeConfigMap.Kind = "ConfigMap"
	kubeConfigMap.Labels = kokiConfigMap.Labels
	kubeConfigMap.Annotations = kokiConfigMap.Annotations

	kubeConfigMap.Data = kokiConfigMap.Data

	return kubeConfigMap, nil
}
