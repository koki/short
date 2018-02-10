package converters

import (
	v1 "k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_PodTemplate_to_Kube(template *types.PodTemplateWrapper) (interface{}, error) {
	var err error
	kubeTemplate := &v1.PodTemplate{}
	kokiTemplate := &template.PodTemplate

	kubeTemplate.Name = kokiTemplate.Name
	kubeTemplate.Namespace = kokiTemplate.Namespace
	if len(kokiTemplate.Version) == 0 {
		kubeTemplate.APIVersion = "v1"
	} else {
		kubeTemplate.APIVersion = kokiTemplate.Version
	}
	kubeTemplate.Kind = "PodTemplate"
	kubeTemplate.ClusterName = kokiTemplate.Cluster
	kubeTemplate.Labels = kokiTemplate.Labels
	kubeTemplate.Annotations = kokiTemplate.Annotations

	kubeSpec, err := revertTemplate(&kokiTemplate.TemplateMetadata, kokiTemplate.PodTemplate)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "template spec")
	}
	kubeTemplate.Template = *kubeSpec

	return kubeTemplate, nil
}
