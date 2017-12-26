package converters

import (
	"fmt"

	apiregistrationv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_APIService_to_Koki_APIService(kubeAPIService *apiregistrationv1beta1.APIService) (*types.APIServiceWrapper, error) {
	var err error
	kokiWrapper := &types.APIServiceWrapper{}
	kokiAPIService := &kokiWrapper.APIService

	kokiAPIService.Name = kubeAPIService.Name
	kokiAPIService.Namespace = kubeAPIService.Namespace
	kokiAPIService.Version = kubeAPIService.APIVersion
	kokiAPIService.Cluster = kubeAPIService.ClusterName
	kokiAPIService.Labels = kubeAPIService.Labels
	kokiAPIService.Annotations = kubeAPIService.Annotations

	err = convertAPIServiceSpec(kubeAPIService.Spec, kokiAPIService)
	if err != nil {
		return nil, err
	}

	err = convertAPIServiceStatus(kubeAPIService.Status, kokiAPIService)
	if err != nil {
		return nil, err
	}

	return kokiWrapper, nil
}

func convertAPIServiceSpec(kubeSpec apiregistrationv1beta1.APIServiceSpec, kokiSpec *types.APIService) error {
	if kubeSpec.Service != nil {
		kokiSpec.Service = fmt.Sprintf("%s/%s", kubeSpec.Service.Namespace, kubeSpec.Service.Name)
	}

	kokiSpec.GroupVersion = fmt.Sprintf("%s/%s", kubeSpec.Group, kubeSpec.Version)
	kokiSpec.TLSVerify = !kubeSpec.InsecureSkipTLSVerify
	kokiSpec.CABundle = kubeSpec.CABundle
	kokiSpec.MinGroupPriority = kubeSpec.GroupPriorityMinimum
	kokiSpec.VersionPriority = kubeSpec.VersionPriority

	return nil
}

func convertAPIServiceStatus(kubeStatus apiregistrationv1beta1.APIServiceStatus, kokiAPIService *types.APIService) error {
	if kubeStatus.Conditions == nil {
		return nil
	}

	var kokiConditions []types.APIServiceCondition
	for i := range kubeStatus.Conditions {
		kubeCondition := kubeStatus.Conditions[i]
		kokiCondition, err := convertKubeAPIServiceCondition(kubeCondition)
		if err != nil {
			return err
		}
		kokiConditions = append(kokiConditions, kokiCondition)
	}
	kokiAPIService.Conditions = kokiConditions

	return nil
}

func convertKubeAPIServiceCondition(kubeCondition apiregistrationv1beta1.APIServiceCondition) (types.APIServiceCondition, error) {
	var kokiCondition types.APIServiceCondition

	kokiAPIServiceConditionType, err := convertAPIServiceConditionType(kubeCondition.Type)
	if err != nil {
		return kokiCondition, err
	}
	kokiCondition.Type = kokiAPIServiceConditionType

	kokiAPIServiceConditionStatus, err := convertAPIServiceConditionStatus(kubeCondition.Status)
	if err != nil {
		return kokiCondition, err
	}
	kokiCondition.Status = kokiAPIServiceConditionStatus

	kokiCondition.LastTransitionTime = kubeCondition.LastTransitionTime
	kokiCondition.Reason = kubeCondition.Reason
	kokiCondition.Message = kubeCondition.Message

	return kokiCondition, nil
}

func convertAPIServiceConditionType(kubeType apiregistrationv1beta1.APIServiceConditionType) (types.APIServiceConditionType, error) {
	if kubeType == "" {
		return "", nil
	}

	switch kubeType {
	case apiregistrationv1beta1.Available:
		return types.APIServiceAvailable, nil
	}

	return "", serrors.InvalidValueErrorf(kubeType, "unrecognized api_service condition type")
}

func convertAPIServiceConditionStatus(kubeStatus apiregistrationv1beta1.ConditionStatus) (types.APIServiceConditionStatus, error) {
	if kubeStatus == "" {
		return "", nil
	}

	switch kubeStatus {
	case apiregistrationv1beta1.ConditionTrue:
		return types.APIServiceConditionTrue, nil
	case apiregistrationv1beta1.ConditionFalse:
		return types.APIServiceConditionFalse, nil
	case apiregistrationv1beta1.ConditionUnknown:
		return types.APIServiceConditionUnknown, nil
	}

	return "", serrors.InvalidValueErrorf(kubeStatus, "unrecognized api_service condition status")
}
