package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CertificateSigningRequestWrapper struct {
	CertificateSigningRequest CertificateSigningRequest `json:"csr"`
}

type CertificateSigningRequest struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	//Spec fields
	Request  []byte              `json:"request,omitempty"`
	Usages   []KeyUsage          `json:"usages,omitempty"`
	Username string              `json:"username,omitempty"`
	UID      string              `json:"uid,omitempty"`
	Groups   []string            `json:"groups,omitempty"`
	Extra    map[string][]string `json:"extra,omitempty"`

	//Status fields
	Certificate []byte                               `json:"cert,omitempty"`
	Conditions  []CertificateSigningRequestCondition `json:"conditions,omitempty"`
}

type RequestConditionType string

const (
	CertificateApproved RequestConditionType = "approved"
	CertificateDenied   RequestConditionType = "denied"
)

type CertificateSigningRequestCondition struct {
	Type           RequestConditionType `json:"type,omitempty"`
	Reason         string               `json:"reason,omitempty"`
	Message        string               `json:"message,omitempty"`
	LastUpdateTime metav1.Time          `json:"last_update,omitempty"`
}

type KeyUsage string

const (
	UsageSigning            KeyUsage = "signing"
	UsageDigitalSignature   KeyUsage = "digital signature"
	UsageContentCommittment KeyUsage = "content committment"
	UsageKeyEncipherment    KeyUsage = "key encipherment"
	UsageKeyAgreement       KeyUsage = "key agreement"
	UsageDataEncipherment   KeyUsage = "data encipherment"
	UsageCertSign           KeyUsage = "cert sign"
	UsageCRLSign            KeyUsage = "crl sign"
	UsageEncipherOnly       KeyUsage = "encipher only"
	UsageDecipherOnly       KeyUsage = "decipher only"
	UsageAny                KeyUsage = "any"
	UsageServerAuth         KeyUsage = "server auth"
	UsageClientAuth         KeyUsage = "client auth"
	UsageCodeSigning        KeyUsage = "code signing"
	UsageEmailProtection    KeyUsage = "email protection"
	UsageSMIME              KeyUsage = "s/mime"
	UsageIPsecEndSystem     KeyUsage = "ipsec end system"
	UsageIPsecTunnel        KeyUsage = "ipsec tunnel"
	UsageIPsecUser          KeyUsage = "ipsec user"
	UsageTimestamping       KeyUsage = "timestamping"
	UsageOCSPSigning        KeyUsage = "ocsp signing"
	UsageMicrosoftSGC       KeyUsage = "microsoft sgc"
	UsageNetscapSGC         KeyUsage = "netscape sgc"
)
