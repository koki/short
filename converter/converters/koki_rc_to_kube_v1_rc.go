package converters

import (
	"reflect"

	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Koki_ReplicationController_to_Kube_v1_ReplicationController(rc *types.ReplicationControllerWrapper) (*v1.ReplicationController, error) {
	var err error
	kubeRC := &v1.ReplicationController{}
	kokiRC := rc.ReplicationController

	kubeRC.Name = kokiRC.Name
	kubeRC.Namespace = kokiRC.Namespace
	if len(kokiRC.Version) == 0 {
		kubeRC.APIVersion = "v1"
	} else {
		kubeRC.APIVersion = kokiRC.Version
	}
	kubeRC.Kind = "ReplicationController"
	kubeRC.ClusterName = kokiRC.Cluster
	kubeRC.Labels = kokiRC.Labels
	kubeRC.Annotations = kokiRC.Annotations

	kubeSpec := &kubeRC.Spec
	kubeSpec.Replicas = kokiRC.Replicas
	kubeSpec.MinReadySeconds = kokiRC.MinReadySeconds

	kubeSpec.Selector = kokiRC.Selector
	kubeSpec.Template, err = revertTemplate(kokiRC.TemplateMetadata, kokiRC.PodTemplate)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, "pod template")
	}

	// Make sure there's at least one Label in the Template and the Selector.
	if kubeSpec.Template != nil {
		if len(kubeSpec.Template.Labels) == 0 {
			if len(kubeSpec.Selector) > 0 {
				kokiRC.TemplateMetadata.Labels = kubeSpec.Selector
			} else {
				kokiRC.TemplateMetadata.Labels = map[string]string{
					"koki.io/selector.name": kokiRC.Name,
				}
				kokiRC.Selector = map[string]string{
					"koki.io/selector.name": kokiRC.Name,
				}
			}
		}
	}

	if kokiRC.Status != nil {
		kubeRC.Status = *kokiRC.Status
	}

	return kubeRC, nil
}

func revertTemplate(kokiMeta *types.PodTemplateMeta, kokiSpec types.PodTemplate) (*v1.PodTemplateSpec, error) {
	var hasMeta = kokiMeta != nil
	var hasSpec = !reflect.DeepEqual(kokiSpec, types.PodTemplate{})
	if !hasMeta && !hasSpec {
		return nil, nil
	}

	template := v1.PodTemplateSpec{}

	if hasMeta {
		template.ObjectMeta = revertPodObjectMeta(*kokiMeta)
	}

	if hasSpec {
		spec, err := revertPodSpec(kokiSpec)
		if err != nil {
			return nil, err
		}
		template.Spec = *spec
	}

	return &template, nil
}
