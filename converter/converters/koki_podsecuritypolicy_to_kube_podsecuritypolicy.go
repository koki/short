package converters

import (
	"fmt"

	"k8s.io/api/core/v1"
	exts "k8s.io/api/extensions/v1beta1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_PodSecurityPolicy_to_Kube_PodSecurityPolicy(podSecurityPolicy *types.PodSecurityPolicyWrapper) (*exts.PodSecurityPolicy, error) {
	kubePodSecurityPolicy := &exts.PodSecurityPolicy{}
	kokiPodSecurityPolicy := &podSecurityPolicy.PodSecurityPolicy

	kubePodSecurityPolicy.Name = kokiPodSecurityPolicy.Name
	kubePodSecurityPolicy.Namespace = kokiPodSecurityPolicy.Namespace
	if len(kokiPodSecurityPolicy.Version) == 0 {
		kubePodSecurityPolicy.APIVersion = "extensions/v1beta1"
	} else {
		kubePodSecurityPolicy.APIVersion = kokiPodSecurityPolicy.Version
	}
	kubePodSecurityPolicy.Kind = "PodSecurityPolicy"
	kubePodSecurityPolicy.ClusterName = kokiPodSecurityPolicy.Cluster
	kubePodSecurityPolicy.Labels = kokiPodSecurityPolicy.Labels
	kubePodSecurityPolicy.Annotations = kokiPodSecurityPolicy.Annotations

	spec, err := revertPodSecurityPolicySpec(kokiPodSecurityPolicy)
	if err != nil {
		return nil, err
	}
	kubePodSecurityPolicy.Spec = spec

	return kubePodSecurityPolicy, nil
}

func revertPodSecurityPolicySpec(kokiPodSecurityPolicy *types.PodSecurityPolicy) (exts.PodSecurityPolicySpec, error) {
	var kubeSpec exts.PodSecurityPolicySpec

	if kokiPodSecurityPolicy == nil {
		return kubeSpec, nil
	}
	kubeSpec.Privileged = kokiPodSecurityPolicy.Privileged

	kubeSpec.DefaultAddCapabilities = revertPodSecurityPolicyCapabilities(kokiPodSecurityPolicy.DefaultCapabilities)
	kubeSpec.RequiredDropCapabilities = revertPodSecurityPolicyCapabilities(kokiPodSecurityPolicy.DenyCapabilities)
	kubeSpec.AllowedCapabilities = revertPodSecurityPolicyCapabilities(kokiPodSecurityPolicy.AllowCapabilities)

	volumes, err := revertPodSecurityPolicyVolumePlugins(kokiPodSecurityPolicy.VolumePlugins)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.Volumes = volumes

	var net, pid, ipc bool
	for i := range kokiPodSecurityPolicy.HostMode {
		mode := kokiPodSecurityPolicy.HostMode[i]
		switch mode {
		case types.HostModeNet:
			net = true
		case types.HostModePID:
			pid = true
		case types.HostModeIPC:
			ipc = true
		default:
			return kubeSpec, serrors.InvalidInstanceError(mode)
		}
	}
	kubeSpec.HostNetwork = net
	kubeSpec.HostPID = pid
	kubeSpec.HostIPC = ipc

	kubeSpec.HostPorts = revertPodSecurityPolicyHostPortRanges(kokiPodSecurityPolicy.HostPortRanges)

	fsgidPolicy, err := revertPodSecurityPolicyFSGIDPolicy(kokiPodSecurityPolicy.FSGIDPolicy)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.FSGroup = fsgidPolicy

	gidPolicy, err := revertPodSecurityPolicyGIDPolicy(kokiPodSecurityPolicy.GIDPolicy)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.SupplementalGroups = gidPolicy

	uids, err := revertPodSecurityPolicyUIDPolicy(kokiPodSecurityPolicy.UIDPolicy)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.RunAsUser = uids

	selinux, err := revertPodSecurityPolicySELinuxPolicy(kokiPodSecurityPolicy.SELinux)
	if err != nil {
		return kubeSpec, err
	}
	kubeSpec.SELinux = selinux

	kubeSpec.ReadOnlyRootFilesystem = kokiPodSecurityPolicy.ReadOnlyRootFS
	kubeSpec.AllowPrivilegeEscalation = kokiPodSecurityPolicy.AllowEscalation
	kubeSpec.DefaultAllowPrivilegeEscalation = kokiPodSecurityPolicy.AllowEscalationDefault

	kubeSpec.AllowedHostPaths = revertAllowedHostPaths(kokiPodSecurityPolicy.AllowedHostPaths)
	kubeSpec.AllowedFlexVolumes = revertAllowedFlexVolumes(kokiPodSecurityPolicy.AllowedFlexVolumes)

	return kubeSpec, nil
}

func revertPodSecurityPolicySELinuxPolicy(kokiSELinuxPolicy types.SELinuxPolicy) (exts.SELinuxStrategyOptions, error) {
	var kubeSELinuxPolicy exts.SELinuxStrategyOptions

	var policyType exts.SELinuxStrategy
	switch kokiSELinuxPolicy.Policy {
	case types.SELinuxPolicyMust:
		policyType = exts.SELinuxStrategyMustRunAs
	case types.SELinuxPolicyAny:
		policyType = exts.SELinuxStrategyRunAsAny
	case "":
	default:
		return kubeSELinuxPolicy, fmt.Errorf("Invalid SELinuxStrategy %s", kokiSELinuxPolicy.Policy)
	}
	kubeSELinuxPolicy.Rule = policyType

	SELinux := kokiSELinuxPolicy.SELinux
	if SELinux.User == "" && SELinux.Role == "" && SELinux.Role == "" && SELinux.Level == "" {
		return kubeSELinuxPolicy, nil
	}

	kubeSELinuxPolicy.SELinuxOptions = &v1.SELinuxOptions{
		User:  kokiSELinuxPolicy.SELinux.User,
		Role:  kokiSELinuxPolicy.SELinux.Role,
		Type:  kokiSELinuxPolicy.SELinux.Type,
		Level: kokiSELinuxPolicy.SELinux.Level,
	}

	return kubeSELinuxPolicy, nil
}

func revertPodSecurityPolicyFSGIDPolicy(kokiFSGID types.GIDPolicy) (exts.FSGroupStrategyOptions, error) {
	var kubeFSGIDPolicy exts.FSGroupStrategyOptions

	var policyType exts.FSGroupStrategyType
	switch kokiFSGID.Policy {
	case types.GIDPolicyMust:
		policyType = exts.FSGroupStrategyMustRunAs
	case types.GIDPolicyAny:
		policyType = exts.FSGroupStrategyRunAsAny
	case "":
	default:
		return kubeFSGIDPolicy, fmt.Errorf("Invalid FSGroupsStrategy %s", kokiFSGID.Policy)
	}
	kubeFSGIDPolicy.Rule = policyType

	kubeFSGIDPolicy.Ranges = revertPodSecurityPolicyIDRanges(kokiFSGID.Ranges)

	return kubeFSGIDPolicy, nil
}

func revertPodSecurityPolicyUIDPolicy(kokiUIDPolicy types.UIDPolicy) (exts.RunAsUserStrategyOptions, error) {
	var kubeUIDPolicy exts.RunAsUserStrategyOptions

	var policyType exts.RunAsUserStrategy
	switch kokiUIDPolicy.Policy {
	case types.UIDPolicyMust:
		policyType = exts.RunAsUserStrategyMustRunAs
	case types.UIDPolicyAny:
		policyType = exts.RunAsUserStrategyRunAsAny
	case "":
	default:
		return kubeUIDPolicy, fmt.Errorf("Invalid RunAsUserStrategy %s", kokiUIDPolicy.Policy)
	}
	kubeUIDPolicy.Rule = policyType

	kubeUIDPolicy.Ranges = revertPodSecurityPolicyIDRanges(kokiUIDPolicy.Ranges)

	return kubeUIDPolicy, nil
}

func revertPodSecurityPolicyGIDPolicy(kokiGIDPolicy types.GIDPolicy) (exts.SupplementalGroupsStrategyOptions, error) {
	var kubeGIDPolicy exts.SupplementalGroupsStrategyOptions

	var policyType exts.SupplementalGroupsStrategyType
	switch kokiGIDPolicy.Policy {
	case types.GIDPolicyMust:
		policyType = exts.SupplementalGroupsStrategyMustRunAs
	case types.GIDPolicyAny:
		policyType = exts.SupplementalGroupsStrategyRunAsAny
	case "":
	default:
		return kubeGIDPolicy, fmt.Errorf("Invalid SupplementalGroupsStrategy %s", kokiGIDPolicy.Policy)
	}
	kubeGIDPolicy.Rule = policyType

	kubeGIDPolicy.Ranges = revertPodSecurityPolicyIDRanges(kokiGIDPolicy.Ranges)

	return kubeGIDPolicy, nil
}

func revertPodSecurityPolicyIDRanges(kokiIDRanges []types.IDRange) []exts.IDRange {
	var kubeIDRanges []exts.IDRange

	for i := range kokiIDRanges {
		kokiIDRange := kokiIDRanges[i]

		kubeIDRange := exts.IDRange{
			Min: kokiIDRange.Min,
			Max: kokiIDRange.Max,
		}
		kubeIDRanges = append(kubeIDRanges, kubeIDRange)
	}

	return kubeIDRanges
}

func revertAllowedHostPaths(allowedHostPaths []string) []exts.AllowedHostPath {
	var kubeAllowedHostPaths []exts.AllowedHostPath

	for i := range allowedHostPaths {
		allowedHostPath := allowedHostPaths[i]
		kubeAllowedHostPath := exts.AllowedHostPath{
			PathPrefix: allowedHostPath,
		}

		kubeAllowedHostPaths = append(kubeAllowedHostPaths, kubeAllowedHostPath)
	}
	return kubeAllowedHostPaths
}

func revertAllowedFlexVolumes(allowedFlexVolumes []string) []exts.AllowedFlexVolume {
	var kubeAllowedFlexVolumes []exts.AllowedFlexVolume

	for i := range allowedFlexVolumes {
		allowedFlexVolume := allowedFlexVolumes[i]
		kubeAllowedFlexVolume := exts.AllowedFlexVolume{
			Driver: allowedFlexVolume,
		}

		kubeAllowedFlexVolumes = append(kubeAllowedFlexVolumes, kubeAllowedFlexVolume)
	}
	return kubeAllowedFlexVolumes
}

func revertPodSecurityPolicyHostPortRanges(hostPorts []types.HostPortRange) []exts.HostPortRange {
	var kubeHostPorts []exts.HostPortRange

	for i := range hostPorts {
		hostPort := hostPorts[i]

		kubeHostPort := exts.HostPortRange{
			Min: hostPort.Min,
			Max: hostPort.Max,
		}

		kubeHostPorts = append(kubeHostPorts, kubeHostPort)
	}
	return kubeHostPorts
}

func revertPodSecurityPolicyVolumePlugins(kokiVolPlugins []string) ([]exts.FSType, error) {
	var kubeVolumes []exts.FSType

	for i := range kokiVolPlugins {
		kokiVolPlugin := kokiVolPlugins[i]

		kubeVolume, err := revertPodSecurityPolicyVolumePlugin(kokiVolPlugin)
		if err != nil {
			return nil, err
		}

		kubeVolumes = append(kubeVolumes, kubeVolume)
	}
	return kubeVolumes, nil
}

func revertPodSecurityPolicyVolumePlugin(kokiVolPlugin string) (exts.FSType, error) {
	if kokiVolPlugin == "" {
		return "", nil
	}

	switch kokiVolPlugin {
	case types.VolumeTypeAzureFile:
		return exts.AzureFile, nil
	case types.VolumeTypeFlocker:
		return exts.Flocker, nil
	case types.VolumeTypeFlex:
		return exts.FlexVolume, nil
	case types.VolumeTypeHostPath:
		return exts.HostPath, nil
	case types.VolumeTypeEmptyDir:
		return exts.EmptyDir, nil
	case types.VolumeTypeGcePD:
		return exts.GCEPersistentDisk, nil
	case types.VolumeTypeAwsEBS:
		return exts.AWSElasticBlockStore, nil
	case types.VolumeTypeGit:
		return exts.GitRepo, nil
	case types.VolumeTypeSecret:
		return exts.Secret, nil
	case types.VolumeTypeNFS:
		return exts.NFS, nil
	case types.VolumeTypeISCSI:
		return exts.ISCSI, nil
	case types.VolumeTypeGlusterfs:
		return exts.Glusterfs, nil
	case types.VolumeTypePVC:
		return exts.PersistentVolumeClaim, nil
	case types.VolumeTypeRBD:
		return exts.RBD, nil
	case types.VolumeTypeCinder:
		return exts.Cinder, nil
	case types.VolumeTypeCephFS:
		return exts.CephFS, nil
	case types.VolumeTypeDownwardAPI:
		return exts.DownwardAPI, nil
	case types.VolumeTypeFibreChannel:
		return exts.FC, nil
	case types.VolumeTypeConfigMap:
		return exts.ConfigMap, nil
	case types.VolumeTypeQuobyte:
		return exts.Quobyte, nil
	case types.VolumeTypeAzureDisk:
		return exts.AzureDisk, nil
	case types.VolumeTypeAny:
		return exts.All, nil
	}

	return "", fmt.Errorf("Invalid Koki Volume Plugin %s", kokiVolPlugin)
}

func revertPodSecurityPolicyCapabilities(kokiCaps []string) []v1.Capability {
	var kubeCaps []v1.Capability

	for i := range kokiCaps {
		kubeCap := v1.Capability(kokiCaps[i])

		kubeCaps = append(kubeCaps, kubeCap)
	}

	return kubeCaps
}
