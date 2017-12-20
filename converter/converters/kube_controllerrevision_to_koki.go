package converters

import (
	apps "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_ControllerRevision_to_Koki(kubeRev runtime.Object) (*types.ControllerRevisionWrapper, error) {
	groupVersionKind := kubeRev.GetObjectKind().GroupVersionKind()
	groupVersionString := groupVersionKind.GroupVersion().String()
	groupVersionKind.Version = "v1"
	groupVersionKind.Group = "apps"
	kubeRev.GetObjectKind().SetGroupVersionKind(groupVersionKind)

	// Serialize as v1
	b, err := yaml.Marshal(kubeRev)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, kubeRev, "couldn't serialize kube ControllerRevision after setting apiVersion to apps/v1")
	}

	// Deserialize the "generic" kube rev
	genericRev, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, serrors.InvalidInstanceContextErrorf(err, string(b), "couldn't deserialize 'generic' kube ControllerRevision")
	}

	if genericRev, ok := genericRev.(*apps.ControllerRevision); ok {
		kokiWrapper, err := Convert_Kube_v1_ControllerRevision_to_Koki(genericRev)
		if err != nil {
			return nil, err
		}

		kokiRev := &kokiWrapper.ControllerRevision
		kokiRev.Version = groupVersionString
		return kokiWrapper, nil
	}

	return nil, serrors.InvalidInstanceErrorf(genericRev, "didn't deserialize 'generic' kube ControllerRevision as apps/v1.ControllerRevision")
}

func Convert_Kube_v1_ControllerRevision_to_Koki(kubeRev *apps.ControllerRevision) (*types.ControllerRevisionWrapper, error) {
	kokiRev := &types.ControllerRevision{}

	kokiRev.Name = kubeRev.Name
	kokiRev.Namespace = kubeRev.Namespace
	kokiRev.Version = kubeRev.APIVersion
	kokiRev.Cluster = kubeRev.ClusterName
	kokiRev.Labels = kubeRev.Labels
	kokiRev.Annotations = kubeRev.Annotations

	kokiRev.Data = kubeRev.Data
	kokiRev.Revision = kubeRev.Revision

	return &types.ControllerRevisionWrapper{
		ControllerRevision: *kokiRev,
	}, nil
}
