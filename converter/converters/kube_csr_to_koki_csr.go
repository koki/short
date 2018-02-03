package converters

import (
	"k8s.io/api/certificates/v1beta1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_CSR_to_Koki_CSR(kubeCSR *v1beta1.CertificateSigningRequest) (*types.CertificateSigningRequestWrapper, error) {
	kokiWrapper := &types.CertificateSigningRequestWrapper{}
	kokiCSR := &kokiWrapper.CertificateSigningRequest

	kokiCSR.Name = kubeCSR.Name
	kokiCSR.Namespace = kubeCSR.Namespace
	kokiCSR.Version = kubeCSR.APIVersion
	kokiCSR.Cluster = kubeCSR.ClusterName
	kokiCSR.Labels = kubeCSR.Labels
	kokiCSR.Annotations = kubeCSR.Annotations

	err := convertCSRSpec(&kubeCSR.Spec, kokiCSR)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "CSR Spec")
	}

	err = convertCSRStatus(&kubeCSR.Status, kokiCSR)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "CSR Status")
	}

	return kokiWrapper, nil
}

func convertCSRSpec(kubeSpec *v1beta1.CertificateSigningRequestSpec, kokiSpec *types.CertificateSigningRequest) error {
	if kubeSpec == nil {
		return nil
	}

	kokiSpec.Request = kubeSpec.Request

	usages, err := convertCSRUsages(kubeSpec.Usages)
	if err != nil {
		return err
	}
	kokiSpec.Usages = usages

	kokiSpec.Username = kubeSpec.Username
	kokiSpec.UID = kubeSpec.UID
	kokiSpec.Groups = kubeSpec.Groups

	for k, v := range kubeSpec.Extra {
		kokiSpec.Extra[k] = []string(v)
	}

	return nil
}

func convertCSRUsages(kubeUsages []v1beta1.KeyUsage) ([]types.KeyUsage, error) {
	kokiUsages := []types.KeyUsage{}

	for i := range kubeUsages {
		kubeUsage := kubeUsages[i]
		kokiUsage, err := convertCSRUsage(kubeUsage)
		if err != nil {
			return nil, err
		}
		kokiUsages = append(kokiUsages, kokiUsage)
	}

	return kokiUsages, nil
}

func convertCSRUsage(kubeUsage v1beta1.KeyUsage) (types.KeyUsage, error) {
	if kubeUsage == "" {
		return "", nil
	}

	switch kubeUsage {
	case v1beta1.UsageSigning:
		return types.UsageSigning, nil
	case v1beta1.UsageDigitalSignature:
		return types.UsageDigitalSignature, nil
	case v1beta1.UsageContentCommittment:
		return types.UsageContentCommittment, nil
	case v1beta1.UsageKeyEncipherment:
		return types.UsageKeyEncipherment, nil
	case v1beta1.UsageKeyAgreement:
		return types.UsageKeyAgreement, nil
	case v1beta1.UsageDataEncipherment:
		return types.UsageDataEncipherment, nil
	case v1beta1.UsageCertSign:
		return types.UsageCertSign, nil
	case v1beta1.UsageCRLSign:
		return types.UsageCRLSign, nil
	case v1beta1.UsageEncipherOnly:
		return types.UsageEncipherOnly, nil
	case v1beta1.UsageDecipherOnly:
		return types.UsageDecipherOnly, nil
	case v1beta1.UsageAny:
		return types.UsageAny, nil
	case v1beta1.UsageServerAuth:
		return types.UsageServerAuth, nil
	case v1beta1.UsageClientAuth:
		return types.UsageClientAuth, nil
	case v1beta1.UsageCodeSigning:
		return types.UsageCodeSigning, nil
	case v1beta1.UsageEmailProtection:
		return types.UsageEmailProtection, nil
	case v1beta1.UsageSMIME:
		return types.UsageSMIME, nil
	case v1beta1.UsageIPsecEndSystem:
		return types.UsageIPsecEndSystem, nil
	case v1beta1.UsageIPsecTunnel:
		return types.UsageIPsecTunnel, nil
	case v1beta1.UsageIPsecUser:
		return types.UsageIPsecUser, nil
	case v1beta1.UsageTimestamping:
		return types.UsageTimestamping, nil
	case v1beta1.UsageOCSPSigning:
		return types.UsageOCSPSigning, nil
	case v1beta1.UsageMicrosoftSGC:
		return types.UsageMicrosoftSGC, nil
	case v1beta1.UsageNetscapSGC:
		return types.UsageNetscapSGC, nil
	default:
		return "", serrors.InvalidValueErrorf(kubeUsage, "Invalid KeyUsage value")
	}
}

func convertCSRStatus(kubeStatus *v1beta1.CertificateSigningRequestStatus, kokiStatus *types.CertificateSigningRequest) error {
	if kubeStatus == nil {
		return nil
	}

	kokiStatus.Certificate = kubeStatus.Certificate

	conditions, err := convertCSRConditions(kubeStatus.Conditions)
	if err != nil {
		return err
	}
	kokiStatus.Conditions = conditions

	return nil
}

func convertCSRConditions(kubeConditions []v1beta1.CertificateSigningRequestCondition) ([]types.CertificateSigningRequestCondition, error) {
	kokiConditions := []types.CertificateSigningRequestCondition{}

	for i := range kubeConditions {
		kubeCondition := kubeConditions[i]
		kokiCondition, err := convertCSRCondition(kubeCondition)
		if err != nil {
			return nil, err
		}
		kokiConditions = append(kokiConditions, kokiCondition)
	}

	return kokiConditions, nil
}

func convertCSRCondition(kubeCondition v1beta1.CertificateSigningRequestCondition) (types.CertificateSigningRequestCondition, error) {
	kokiCondition := types.CertificateSigningRequestCondition{}

	var kokiConditionType types.RequestConditionType
	if kubeCondition.Type == "" {
		kokiConditionType = ""
	}

	switch kubeCondition.Type {
	case v1beta1.CertificateApproved:
		kokiConditionType = types.CertificateApproved
	case v1beta1.CertificateDenied:
		kokiConditionType = types.CertificateDenied
	default:
		return kokiCondition, serrors.InvalidValueErrorf(kokiCondition.Type, "Invalid CertificateSigningRequestCondition Type")
	}

	kokiCondition.Type = kokiConditionType
	kokiCondition.Reason = kubeCondition.Reason
	kokiCondition.Message = kubeCondition.Message
	kokiCondition.LastUpdateTime = kubeCondition.LastUpdateTime

	return kokiCondition, nil
}
