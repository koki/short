package converters

import (
	apiext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_CRD_to_Koki(kube *apiext.CustomResourceDefinition) (*types.CRDWrapper, error) {
	var err error
	koki := &types.CustomResourceDefinition{}

	koki.Name = kube.Name
	koki.Namespace = kube.Namespace
	koki.Version = kube.APIVersion
	koki.Cluster = kube.ClusterName
	koki.Labels = kube.Labels
	koki.Annotations = kube.Annotations

	koki.CRDMeta = convertCRDMeta(kube.Spec)
	koki.Scope, err = convertCRDScope(kube.Spec.Scope)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "CRD resource scope")
	}
	koki.Validation = convertCRDValidation(kube.Spec.Validation)

	koki.Conditions, err = convertCRDConditions(kube.Status.Conditions)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "CRD conditions")
	}
	koki.Accepted = convertCRDNames(kube.Status.AcceptedNames)

	return &types.CRDWrapper{
		CRD: *koki,
	}, nil
}

func convertCRDMeta(kube apiext.CustomResourceDefinitionSpec) types.CRDMeta {
	meta := types.CRDMeta{
		Group:   kube.Group,
		Version: kube.Version,
		CRDName: convertCRDNames(kube.Names),
	}

	return meta
}

func convertCRDNames(kube apiext.CustomResourceDefinitionNames) types.CRDName {
	return types.CRDName{
		Plural:     kube.Plural,
		Singular:   kube.Singular,
		ShortNames: kube.ShortNames,
		Kind:       kube.Kind,
		ListKind:   kube.ListKind,
	}
}

func convertCRDScope(kube apiext.ResourceScope) (types.CRDResourceScope, error) {
	switch kube {
	case "":
		return "", nil
	case apiext.ClusterScoped:
		return types.CRDClusterScoped, nil
	case apiext.NamespaceScoped:
		return types.CRDNamespaceScoped, nil
	default:
		return "", serrors.InvalidInstanceError(kube)
	}
}

func convertCRDValidation(kube *apiext.CustomResourceValidation) *apiext.JSONSchemaProps {
	if kube == nil {
		return nil
	}

	return kube.OpenAPIV3Schema
}

func convertCRDConditions(kubes []apiext.CustomResourceDefinitionCondition) ([]types.CRDCondition, error) {
	if len(kubes) == 0 {
		return nil, nil
	}

	var err error
	kokis := make([]types.CRDCondition, len(kubes))
	for i, kube := range kubes {
		kokis[i], err = convertCRDCondition(kube)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "[%d]", i)
		}
	}

	return kokis, nil
}

func convertCRDCondition(kube apiext.CustomResourceDefinitionCondition) (types.CRDCondition, error) {
	var err error
	condition := types.CRDCondition{
		LastTransitionTime: kube.LastTransitionTime,
		Reason:             kube.Reason,
		Message:            kube.Message,
	}
	condition.Type, err = convertCRDConditionType(kube.Type)
	if err != nil {
		return condition, serrors.ContextualizeErrorf(err, "type")
	}
	condition.Status, err = convertCRDConditionStatus(kube.Status)
	if err != nil {
		return condition, serrors.ContextualizeErrorf(err, "status")
	}

	return condition, nil
}

func convertCRDConditionType(kube apiext.CustomResourceDefinitionConditionType) (types.CRDConditionType, error) {
	switch kube {
	case apiext.Established:
		return types.CRDEstablished, nil
	case apiext.NamesAccepted:
		return types.CRDNamesAccepted, nil
	case apiext.Terminating:
		return types.CRDTerminating, nil
	default:
		return "", serrors.InvalidInstanceError(kube)
	}
}

func convertCRDConditionStatus(status apiext.ConditionStatus) (types.ConditionStatus, error) {
	if status == "" {
		return "", nil
	}
	if status == apiext.ConditionTrue {
		return types.ConditionTrue, nil
	}
	if status == apiext.ConditionFalse {
		return types.ConditionFalse, nil
	}
	if status == apiext.ConditionUnknown {
		return types.ConditionUnknown, nil
	}
	return "", serrors.InvalidInstanceError(status)
}
