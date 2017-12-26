package converters

import (
	"strings"

	apiregistrationv1beta1 "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_APIService_to_Kube_APIService(apiService *types.APIServiceWrapper) (*apiregistrationv1beta1.APIService, error) {
	kubeAPIService := &apiregistrationv1beta1.APIService{}
	kokiAPIService := &apiService.APIService

	kubeAPIService.Name = kokiAPIService.Name
	kubeAPIService.Namespace = kokiAPIService.Namespace
	if len(kokiAPIService.Version) == 0 {
		kubeAPIService.APIVersion = apiregistrationv1beta1.GroupName
	} else {
		kubeAPIService.APIVersion = kokiAPIService.Version
	}
	kubeAPIService.Kind = "APIService"
	kubeAPIService.ClusterName = kokiAPIService.Cluster
	kubeAPIService.Labels = kokiAPIService.Labels
	kubeAPIService.Annotations = kokiAPIService.Annotations

	spec, err := revertAPIServiceSpec(kokiAPIService)
	if err != nil {
		return nil, err
	}
	kubeAPIService.Spec = spec

	status, err := revertAPIServiceStatus(kokiAPIService)
	if err != nil {
		return nil, err
	}
	kubeAPIService.Status = status

	return kubeAPIService, nil
}

func revertAPIServiceStatus(kokiAPIService *types.APIService) (apiregistrationv1beta1.APIServiceStatus, error) {
	var kubeStatus apiregistrationv1beta1.APIServiceStatus

	if kokiAPIService == nil {
		return kubeStatus, nil
	}

	if kokiAPIService.Conditions == nil {
		return kubeStatus, nil
	}

	var kubeConditions []apiregistrationv1beta1.APIServiceCondition
	for i := range kokiAPIService.Conditions {
		kokiCondition := kokiAPIService.Conditions[i]
		kubeCondition, err := revertAPIServiceCondition(kokiCondition)
		if err != nil {
			return kubeStatus, err
		}
		kubeConditions = append(kubeConditions, kubeCondition)
	}

	kubeStatus.Conditions = kubeConditions
	return kubeStatus, nil
}

func revertAPIServiceCondition(kokiCondition types.APIServiceCondition) (apiregistrationv1beta1.APIServiceCondition, error) {
	var kubeStatus apiregistrationv1beta1.APIServiceCondition

	conditionType, err := revertAPIServiceConditionType(kokiCondition.Type)
	if err != nil {
		return kubeStatus, err
	}
	kubeStatus.Type = conditionType

	status, err := revertAPIServiceServiceStatus(kokiCondition.Status)
	if err != nil {
		return kubeStatus, err
	}
	kubeStatus.Status = status

	kubeStatus.LastTransitionTime = kokiCondition.LastTransitionTime
	kubeStatus.Reason = kokiCondition.Reason
	kubeStatus.Message = kokiCondition.Message

	return kubeStatus, nil
}

func revertAPIServiceConditionType(kokiConditionType types.APIServiceConditionType) (apiregistrationv1beta1.APIServiceConditionType, error) {
	if kokiConditionType == "" {
		return "", nil
	}
	var kubeConditionType apiregistrationv1beta1.APIServiceConditionType

	switch kokiConditionType {
	case types.APIServiceAvailable:
		return apiregistrationv1beta1.Available, nil
	}

	return kubeConditionType, serrors.InvalidValueErrorf(kokiConditionType, "Invalid value for Condition Type")
}

func revertAPIServiceServiceStatus(kokiConditionStatus types.APIServiceConditionStatus) (apiregistrationv1beta1.ConditionStatus, error) {
	if kokiConditionStatus == "" {
		return "", nil
	}
	var kubeStatus apiregistrationv1beta1.ConditionStatus

	switch kokiConditionStatus {
	case types.APIServiceConditionTrue:
		return apiregistrationv1beta1.ConditionTrue, nil
	case types.APIServiceConditionFalse:
		return apiregistrationv1beta1.ConditionFalse, nil
	case types.APIServiceConditionUnknown:
		return apiregistrationv1beta1.ConditionUnknown, nil
	}

	return kubeStatus, serrors.InvalidValueErrorf(kokiConditionStatus, "Invalid value for Condition Status")
}

func revertAPIServiceSpec(kokiAPIService *types.APIService) (apiregistrationv1beta1.APIServiceSpec, error) {
	var kubeSpec apiregistrationv1beta1.APIServiceSpec

	if kokiAPIService == nil {
		return kubeSpec, nil
	}

	service, err := revertAPIServiceService(kokiAPIService.Service)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.Service = service

	group, version, err := revertAPIServiceGroupVersion(kokiAPIService.GroupVersion)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.Group = group
	kubeSpec.Version = version

	kubeSpec.InsecureSkipTLSVerify = !kokiAPIService.TLSVerify

	kubeSpec.CABundle = kokiAPIService.CABundle
	kubeSpec.GroupPriorityMinimum = kokiAPIService.MinGroupPriority
	kubeSpec.VersionPriority = kokiAPIService.VersionPriority

	return kubeSpec, nil
}

func revertAPIServiceService(kokiAPIServiceService string) (*apiregistrationv1beta1.ServiceReference, error) {
	if kokiAPIServiceService == "" {
		return nil, nil
	}

	parts := strings.Split(kokiAPIServiceService, "/")
	if len(parts) != 2 {
		return nil, serrors.InvalidValueErrorf(kokiAPIServiceService, "Invalid value for APIService Service Reference")
	}

	service := &apiregistrationv1beta1.ServiceReference{
		Namespace: parts[0],
		Name:      parts[1],
	}

	return service, nil
}

func revertAPIServiceGroupVersion(groupVersion string) (group string, version string, err error) {
	if groupVersion == "" {
		return "", "", nil
	}

	parts := strings.Split(groupVersion, "/")
	if len(parts) != 2 {
		return "", "", serrors.InvalidValueErrorf(groupVersion, "Invalid value for APIService Service Reference")
	}

	group = parts[0]
	version = parts[1]

	return group, version, nil
}
