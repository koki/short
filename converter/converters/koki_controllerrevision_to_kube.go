package converters

import (
	apps "k8s.io/api/apps/v1"
	appsv1beta1 "k8s.io/api/apps/v1beta1"
	appsv1beta2 "k8s.io/api/apps/v1beta2"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/yaml"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_ControllerRevision_to_Kube(kokiRev *types.ControllerRevisionWrapper) (interface{}, error) {
	// Perform version-agnostic conversion into apps/v1beta2 ControllerRevision.
	kubeControllerRevision, err := Convert_Koki_ControllerRevision_to_Kube_v1(kokiRev)
	if err != nil {
		return nil, err
	}

	// Serialize the "generic" kube ControllerRevision.
	b, err := yaml.Marshal(kubeControllerRevision)
	if err != nil {
		return nil, serrors.InvalidValueContextErrorf(err, kubeControllerRevision, "couldn't serialize 'generic' kube ControllerRevision")
	}

	// Deserialize a versioned kube ControllerRevision using its apiVersion.
	versionedControllerRevision, err := parser.ParseSingleKubeNativeFromBytes(b)
	if err != nil {
		return nil, err
	}

	switch versionedControllerRevision := versionedControllerRevision.(type) {
	case *appsv1beta1.ControllerRevision:
		// Perform apps/v1beta1-specific initialization here.
	case *appsv1beta2.ControllerRevision:
		// Perform apps/v1beta2-specific initialization here.
	case *apps.ControllerRevision:
		// Perform apps/v1-specific initialization here.
	default:
		return nil, serrors.TypeErrorf(versionedControllerRevision, "deserialized the manifest, but not as a supported kube ControllerRevision")
	}

	return versionedControllerRevision, nil
}

func Convert_Koki_ControllerRevision_to_Kube_v1(rev *types.ControllerRevisionWrapper) (*apps.ControllerRevision, error) {
	kubeRev := &apps.ControllerRevision{}
	kokiRev := rev.ControllerRevision

	kubeRev.Name = kokiRev.Name
	kubeRev.Namespace = kokiRev.Namespace
	if len(kokiRev.Version) == 0 {
		kubeRev.APIVersion = "apps/v1"
	} else {
		kubeRev.APIVersion = kokiRev.Version
	}
	kubeRev.Kind = "ControllerRevision"
	kubeRev.ClusterName = kokiRev.Cluster
	kubeRev.Labels = kokiRev.Labels
	kubeRev.Annotations = kokiRev.Annotations

	kubeRev.Data = kokiRev.Data
	kubeRev.Revision = kokiRev.Revision

	return kubeRev, nil
}
