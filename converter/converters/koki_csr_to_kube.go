package converters

import (
	"k8s.io/api/certificates/v1beta1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_CSR_to_Kube_CSR(wrapper *types.CertificateSigningRequestWrapper) (*v1beta1.CertificateSigningRequest, error) {
	var err error
	kubeCSR := &v1beta1.CertificateSigningRequest{}
	kokiCSR := wrapper.CertificateSigningRequest

	kubeCSR.Name = kokiCSR.Name
	kubeCSR.Namespace = kokiCSR.Namespace
	if len(kokiCSR.Version) == 0 {
		kubeCSR.APIVersion = "v1"
	} else {
		kubeCSR.APIVersion = kokiCSR.Version
	}
	kubeCSR.Kind = "CertificateSigningRequest"
	kubeCSR.ClusterName = kokiCSR.Cluster
	kubeCSR.Labels = kokiCSR.Labels
	kubeCSR.Annotations = kokiCSR.Annotations

	kubeCSR.Spec, err = revertCSRSpec(kokiCSR)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "CSR.Spec")
	}

	kubeCSR.Status, err = revertCSRStatus(kokiCSR)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "CSR.Status")
	}

	return kubeCSR, nil
}

func revertCSRSpec(kokiCSR types.CertificateSigningRequest) (v1beta1.CertificateSigningRequestSpec, error) {
	kubeSpec := v1beta1.CertificateSigningRequestSpec{}

	kubeSpec.Request = kokiCSR.Request
	usages, err := revertCSRUsages(kokiCSR.Usages)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.Usages = usages

	kubeSpec.Username = kokiCSR.Username
	kubeSpec.UID = kokiCSR.UID
	kubeSpec.Groups = kokiCSR.Groups
	for k, v := range kokiCSR.Extra {
		kubeSpec.Extra[k] = v1beta1.ExtraValue(v)
	}

	return kubeSpec, nil
}

func revertCSRUsages(kokiUsages []types.KeyUsage) ([]v1beta1.KeyUsage, error) {
	kubeUsages := []v1beta1.KeyUsage{}

	for i := range kokiUsages {
		kokiUsage := kokiUsages[i]
		kubeUsage, err := revertCSRUsage(kokiUsage)
		if err != nil {
			return nil, err
		}
		kubeUsages = append(kubeUsages, kubeUsage)
	}

	return kubeUsages, nil
}

func revertCSRUsage(kokiUsage types.KeyUsage) (v1beta1.KeyUsage, error) {
	if kokiUsage == "" {
		return "", nil
	}

	switch kokiUsage {
	case types.UsageSigning:
		return v1beta1.UsageSigning, nil
	case types.UsageDigitalSignature:
		return v1beta1.UsageDigitalSignature, nil
	case types.UsageContentCommittment:
		return v1beta1.UsageContentCommittment, nil
	case types.UsageKeyEncipherment:
		return v1beta1.UsageKeyEncipherment, nil
	case types.UsageKeyAgreement:
		return v1beta1.UsageKeyAgreement, nil
	case types.UsageDataEncipherment:
		return v1beta1.UsageDataEncipherment, nil
	case types.UsageCertSign:
		return v1beta1.UsageCertSign, nil
	case types.UsageCRLSign:
		return v1beta1.UsageCRLSign, nil
	case types.UsageEncipherOnly:
		return v1beta1.UsageEncipherOnly, nil
	case types.UsageDecipherOnly:
		return v1beta1.UsageDecipherOnly, nil
	case types.UsageAny:
		return v1beta1.UsageAny, nil
	case types.UsageServerAuth:
		return v1beta1.UsageServerAuth, nil
	case types.UsageClientAuth:
		return v1beta1.UsageClientAuth, nil
	case types.UsageCodeSigning:
		return v1beta1.UsageCodeSigning, nil
	case types.UsageEmailProtection:
		return v1beta1.UsageEmailProtection, nil
	case types.UsageSMIME:
		return v1beta1.UsageSMIME, nil
	case types.UsageIPsecEndSystem:
		return v1beta1.UsageIPsecEndSystem, nil
	case types.UsageIPsecTunnel:
		return v1beta1.UsageIPsecTunnel, nil
	case types.UsageIPsecUser:
		return v1beta1.UsageIPsecUser, nil
	case types.UsageTimestamping:
		return v1beta1.UsageTimestamping, nil
	case types.UsageOCSPSigning:
		return v1beta1.UsageOCSPSigning, nil
	case types.UsageMicrosoftSGC:
		return v1beta1.UsageMicrosoftSGC, nil
	case types.UsageNetscapSGC:
		return v1beta1.UsageNetscapSGC, nil
	default:
		return "", serrors.InvalidValueErrorf(kokiUsage, "Invalid KeyUsage")
	}
}

func revertCSRStatus(kokiCSR types.CertificateSigningRequest) (v1beta1.CertificateSigningRequestStatus, error) {
	kubeStatus := v1beta1.CertificateSigningRequestStatus{}

	kubeStatus.Certificate = kokiCSR.Certificate

	csrConditions, err := revertCSRConditions(kokiCSR.Conditions)
	if err != nil {
		return kubeStatus, err
	}
	kubeStatus.Conditions = csrConditions

	return kubeStatus, nil
}

func revertCSRConditions(kokiConditions []types.CertificateSigningRequestCondition) ([]v1beta1.CertificateSigningRequestCondition, error) {
	kubeConditions := []v1beta1.CertificateSigningRequestCondition{}

	for i := range kokiConditions {
		kokiCondition := kokiConditions[i]
		kubeCondition, err := revertCSRCondition(kokiCondition)
		if err != nil {
			return nil, err
		}
		kubeConditions = append(kubeConditions, kubeCondition)
	}
	return kubeConditions, nil
}

func revertCSRCondition(kokiCondition types.CertificateSigningRequestCondition) (v1beta1.CertificateSigningRequestCondition, error) {
	kubeCondition := v1beta1.CertificateSigningRequestCondition{}

	var kubeConditionType v1beta1.RequestConditionType
	if kokiCondition.Type == "" {
		kubeConditionType = ""
	}

	switch kokiCondition.Type {
	case types.CertificateApproved:
		kubeConditionType = v1beta1.CertificateApproved
	case types.CertificateDenied:
		kubeConditionType = v1beta1.CertificateDenied
	default:
		return kubeCondition, serrors.InvalidValueErrorf(kokiCondition.Type, "Invalid CertificateSigningRequest Condition Type")
	}

	kubeCondition.Type = kubeConditionType
	kubeCondition.Reason = kokiCondition.Reason
	kubeCondition.Message = kokiCondition.Message
	kubeCondition.LastUpdateTime = kokiCondition.LastUpdateTime

	return kubeCondition, nil
}
