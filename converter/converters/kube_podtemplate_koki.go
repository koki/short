package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_PodTemplate_to_Koki(kubeTemplate *v1.PodTemplate) (*types.PodTemplateWrapper, error) {
	var err error
	kokiTemplate := &types.PodTemplateResource{}

	kokiTemplate.Name = kubeTemplate.Name
	kokiTemplate.Namespace = kubeTemplate.Namespace
	kokiTemplate.Version = kubeTemplate.APIVersion
	kokiTemplate.Cluster = kubeTemplate.ClusterName
	kokiTemplate.Labels = kubeTemplate.Labels
	kokiTemplate.Annotations = kubeTemplate.Annotations

	meta, template, err := convertTemplate(kubeTemplate.Template)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "")
	}
	kokiTemplate.TemplateMetadata = *meta
	kokiTemplate.PodTemplate = template

	return &types.PodTemplateWrapper{
		PodTemplate: *kokiTemplate,
	}, nil
}
