package converters

import (
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_CRD_to_Kube(kokiWrapper *types.CRDWrapper) (*apiext.CustomResourceDefinition, error) {
	var err error
	kube := &apiext.CustomResourceDefinition{}
	koki := kokiWrapper.CRD

	kube.Name = koki.Name
	kube.Namespace = koki.Namespace
	if len(koki.Version) == 0 {
		kube.APIVersion = "apiextensions/v1beta1"
	} else {
		kube.APIVersion = koki.Version
	}
	kube.Kind = "CustomResourceDefinition"
	kube.ClusterName = koki.Cluster
	kube.Labels = koki.Labels
	kube.Annotations = koki.Annotations

	kubeSpec := &kube.Spec
	kubeSpec.Group = koki.CRDMeta.Group
	kubeSpec.Version = koki.CRDMeta.Version
	kubeSpec.Names = revertCRDNames(koki.CRDMeta.CRDName)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "CRD names")
	}
	kubeSpec.Scope, err = revertCRDScope(koki.Scope)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "CRD resource scope")
	}
	kubeSpec.Validation = revertCRDValidation(koki.Validation)

	kube.Status.Conditions, err = revertCRDConditions(koki.Conditions)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "CRD conditions")
	}
	kube.Status.AcceptedNames = revertCRDNames(koki.Accepted)

	return kube, nil
}

func revertCRDNames(koki types.CRDName) apiext.CustomResourceDefinitionNames {
	return apiext.CustomResourceDefinitionNames{
		Plural:     koki.Plural,
		Singular:   koki.Singular,
		ShortNames: koki.ShortNames,
		Kind:       koki.Kind,
		ListKind:   koki.ListKind,
	}
}

func revertCRDScope(koki types.CRDResourceScope) (apiext.ResourceScope, error) {
	switch koki {
	case "":
		return "", nil
	case types.CRDClusterScoped:
		return apiext.ClusterScoped, nil
	case types.CRDNamespaceScoped:
		return apiext.NamespaceScoped, nil
	default:
		return "", serrors.InvalidInstanceError(koki)
	}
}

func revertCRDValidation(koki *apiext.JSONSchemaProps) *apiext.CustomResourceValidation {
	if koki == nil {
		return nil
	}

	return &apiext.CustomResourceValidation{
		OpenAPIV3Schema: koki,
	}
}

func revertCRDConditions(kokis []types.CRDCondition) ([]apiext.CustomResourceDefinitionCondition, error) {
	if len(kokis) == 0 {
		return nil, nil
	}

	var err error
	kubes := make([]apiext.CustomResourceDefinitionCondition, len(kokis))
	for i, koki := range kokis {
		kubes[i], err = revertCRDCondition(koki)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "[%d]", i)
		}
	}

	return kubes, nil
}

func revertCRDCondition(koki types.CRDCondition) (apiext.CustomResourceDefinitionCondition, error) {
	var err error
	condition := apiext.CustomResourceDefinitionCondition{
		LastTransitionTime: koki.LastTransitionTime,
		Reason:             koki.Reason,
		Message:            koki.Message,
	}
	condition.Type, err = revertCRDConditionType(koki.Type)
	if err != nil {
		return condition, serrors.ContextualizeErrorf(err, "type")
	}
	condition.Status, err = revertCRDConditionStatus(koki.Status)
	if err != nil {
		return condition, serrors.ContextualizeErrorf(err, "status")
	}

	return condition, nil
}

func revertCRDConditionType(koki types.CRDConditionType) (apiext.CustomResourceDefinitionConditionType, error) {
	switch koki {
	case types.CRDEstablished:
		return apiext.Established, nil
	case types.CRDNamesAccepted:
		return apiext.NamesAccepted, nil
	case types.CRDTerminating:
		return apiext.Terminating, nil
	default:
		return "", serrors.InvalidInstanceError(koki)
	}
}

func revertCRDConditionStatus(status types.ConditionStatus) (apiext.ConditionStatus, error) {
	if status == "" {
		return "", nil
	}
	if status == types.ConditionTrue {
		return apiext.ConditionTrue, nil
	}
	if status == types.ConditionFalse {
		return apiext.ConditionFalse, nil
	}
	if status == types.ConditionUnknown {
		return apiext.ConditionUnknown, nil
	}
	return "", serrors.InvalidInstanceError(status)
}
