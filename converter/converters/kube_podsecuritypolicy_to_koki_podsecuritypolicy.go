package converters

import (
	"fmt"

	"k8s.io/api/core/v1"
	exts "k8s.io/api/extensions/v1beta1"

	"github.com/koki/short/types"
)

func Convert_Kube_PodSecurityPolicy_to_Koki_PodSecurityPolicy(kubePodSecurityPolicy *exts.PodSecurityPolicy) (*types.PodSecurityPolicyWrapper, error) {
	kokiWrapper := &types.PodSecurityPolicyWrapper{}
	kokiPodSecurityPolicy := &kokiWrapper.PodSecurityPolicy

	kokiPodSecurityPolicy.Name = kubePodSecurityPolicy.Name
	kokiPodSecurityPolicy.Namespace = kubePodSecurityPolicy.Namespace
	kokiPodSecurityPolicy.Version = kubePodSecurityPolicy.APIVersion
	kokiPodSecurityPolicy.Cluster = kubePodSecurityPolicy.ClusterName
	kokiPodSecurityPolicy.Labels = kubePodSecurityPolicy.Labels
	kokiPodSecurityPolicy.Annotations = kubePodSecurityPolicy.Annotations

	err := convertPodSecurityPolicySpec(kubePodSecurityPolicy.Spec, kokiPodSecurityPolicy)
	if err != nil {
		return nil, err
	}
	return kokiWrapper, nil
}

func convertPodSecurityPolicySpec(kubeSpec exts.PodSecurityPolicySpec, kokiPodSecurityPolicy *types.PodSecurityPolicy) error {
	if kokiPodSecurityPolicy == nil {
		return fmt.Errorf("Writing to uninitialized koki Pointer")
	}

	kokiPodSecurityPolicy.Privileged = kubeSpec.Privileged
	kokiPodSecurityPolicy.AllowCapabilities = convertPodSecurityPolicyCapabilities(kubeSpec.AllowedCapabilities)
	kokiPodSecurityPolicy.DenyCapabilities = convertPodSecurityPolicyCapabilities(kubeSpec.RequiredDropCapabilities)
	kokiPodSecurityPolicy.DefaultCapabilities = convertPodSecurityPolicyCapabilities(kubeSpec.DefaultAddCapabilities)

	volumePlugins, err := convertPodSecurityPolicyVolumePlugins(kubeSpec.Volumes)
	if err != nil {
		return err
	}
	kokiPodSecurityPolicy.VolumePlugins = volumePlugins

	var hostMode []types.HostMode
	if kubeSpec.HostNetwork {
		hostMode = append(hostMode, types.HostModeNet)
	}
	if kubeSpec.HostPID {
		hostMode = append(hostMode, types.HostModePID)
	}
	if kubeSpec.HostIPC {
		hostMode = append(hostMode, types.HostModeIPC)
	}
	kokiPodSecurityPolicy.HostMode = hostMode

	kokiPodSecurityPolicy.HostPortRanges = convertHostPortRange(kubeSpec.HostPorts)

	selinux, err := convertPodSecurityPolicySELinux(kubeSpec.SELinux)
	if err != nil {
		return err
	}
	kokiPodSecurityPolicy.SELinux = selinux

	uids, err := convertPodSecurityPolicyUIDPolicy(kubeSpec.RunAsUser)
	if err != nil {
		return err
	}
	kokiPodSecurityPolicy.UIDPolicy = uids

	gids, err := convertPodSecurityPolicyGIDPolicy(kubeSpec.SupplementalGroups)
	if err != nil {
		return err
	}
	kokiPodSecurityPolicy.GIDPolicy = gids

	fsgids, err := convertPodSecurityPolicyFSGIDPolicy(kubeSpec.FSGroup)
	if err != nil {
		return err
	}
	kokiPodSecurityPolicy.FSGIDPolicy = fsgids

	kokiPodSecurityPolicy.ReadOnlyRootFS = kubeSpec.ReadOnlyRootFilesystem
	kokiPodSecurityPolicy.AllowEscalation = kubeSpec.AllowPrivilegeEscalation
	kokiPodSecurityPolicy.AllowEscalationDefault = kubeSpec.DefaultAllowPrivilegeEscalation

	kokiPodSecurityPolicy.AllowedHostPaths = convertAllowedHostPaths(kubeSpec.AllowedHostPaths)
	kokiPodSecurityPolicy.AllowedFlexVolumes = convertAllowedFlexVolumes(kubeSpec.AllowedFlexVolumes)

	return nil
}

func convertAllowedHostPaths(kubeHostPaths []exts.AllowedHostPath) []string {
	var kokiHostPaths []string

	for i := range kubeHostPaths {
		kubeHostPath := kubeHostPaths[i]

		if kubeHostPath.PathPrefix != "" {
			kokiHostPaths = append(kokiHostPaths, kubeHostPath.PathPrefix)
		}
	}
	return kokiHostPaths
}

func convertAllowedFlexVolumes(kubeFlexVolumes []exts.AllowedFlexVolume) []string {
	var kokiFlexVolumes []string

	for i := range kubeFlexVolumes {
		kubeFlexVolume := kubeFlexVolumes[i]

		if kubeFlexVolume.Driver != "" {
			kokiFlexVolumes = append(kokiFlexVolumes, kubeFlexVolume.Driver)
		}
	}
	return kokiFlexVolumes
}

func convertPodSecurityPolicyFSGIDPolicy(kubeFSGIDPolicy exts.FSGroupStrategyOptions) (types.GIDPolicy, error) {
	var kokiGIDPolicy types.GIDPolicy

	var policyType types.GIDPolicyType
	switch kubeFSGIDPolicy.Rule {
	case exts.FSGroupStrategyMustRunAs:
		policyType = types.GIDPolicyMust
	case exts.FSGroupStrategyRunAsAny:
		policyType = types.GIDPolicyAny
	case "":
	default:
		return kokiGIDPolicy, fmt.Errorf("Invalid FSGroupsStrategy %s", kubeFSGIDPolicy.Rule)
	}
	kokiGIDPolicy.Policy = policyType

	kokiGIDPolicy.Ranges = convertPodSecurityPolicyIDRanges(kubeFSGIDPolicy.Ranges)

	return kokiGIDPolicy, nil
}

func convertPodSecurityPolicyUIDPolicy(kubeUIDPolicy exts.RunAsUserStrategyOptions) (types.UIDPolicy, error) {
	var kokiUIDPolicy types.UIDPolicy

	var policyType types.UIDPolicyType
	switch kubeUIDPolicy.Rule {
	case exts.RunAsUserStrategyMustRunAs:
		policyType = types.UIDPolicyMust
	case exts.RunAsUserStrategyMustRunAsNonRoot:
		policyType = types.UIDPolicyNonRoot
	case exts.RunAsUserStrategyRunAsAny:
		policyType = types.UIDPolicyAny
	case "":
	default:
		return kokiUIDPolicy, fmt.Errorf("Invalid RunAsUserStrategy %s", kubeUIDPolicy.Rule)
	}
	kokiUIDPolicy.Policy = policyType

	kokiUIDPolicy.Ranges = convertPodSecurityPolicyIDRanges(kubeUIDPolicy.Ranges)

	return kokiUIDPolicy, nil
}

func convertPodSecurityPolicyGIDPolicy(kubeGIDPolicy exts.SupplementalGroupsStrategyOptions) (types.GIDPolicy, error) {
	var kokiGIDPolicy types.GIDPolicy

	var policyType types.GIDPolicyType
	switch kubeGIDPolicy.Rule {
	case exts.SupplementalGroupsStrategyMustRunAs:
		policyType = types.GIDPolicyMust
	case exts.SupplementalGroupsStrategyRunAsAny:
		policyType = types.GIDPolicyAny
	case "":
	default:
		return kokiGIDPolicy, fmt.Errorf("Invalid SupplementalGroupsStrategy %s", kubeGIDPolicy.Rule)
	}
	kokiGIDPolicy.Policy = policyType

	kokiGIDPolicy.Ranges = convertPodSecurityPolicyIDRanges(kubeGIDPolicy.Ranges)

	return kokiGIDPolicy, nil
}

func convertPodSecurityPolicyIDRanges(kubeIDRanges []exts.IDRange) []types.IDRange {
	var kokiIDRanges []types.IDRange

	for i := range kubeIDRanges {
		kubeIDRange := kubeIDRanges[i]

		kokiIDRange := types.IDRange{
			Min: kubeIDRange.Min,
			Max: kubeIDRange.Max,
		}

		kokiIDRanges = append(kokiIDRanges, kokiIDRange)
	}
	return kokiIDRanges
}

func convertPodSecurityPolicySELinux(kubeSELinux exts.SELinuxStrategyOptions) (types.SELinuxPolicy, error) {
	var kokiSELinux types.SELinuxPolicy

	policy, err := convertPodSecurityPolicySELinuxPolicy(kubeSELinux.Rule)
	if err != nil {
		return kokiSELinux, err
	}
	seLinux := convertSELinux(kubeSELinux.SELinuxOptions)
	if seLinux != nil {
		kokiSELinux.SELinux = *seLinux
	}
	kokiSELinux.Policy = policy

	return kokiSELinux, nil
}

func convertPodSecurityPolicySELinuxPolicy(rule exts.SELinuxStrategy) (types.SELinuxPolicyType, error) {
	if rule == "" {
		return "", nil
	}

	switch rule {
	case exts.SELinuxStrategyMustRunAs:
		return types.SELinuxPolicyMust, nil
	case exts.SELinuxStrategyRunAsAny:
		return types.SELinuxPolicyAny, nil
	}

	return "", fmt.Errorf("Invalid SELinux Strategy %s", rule)
}

func convertHostPortRange(kubeHostPorts []exts.HostPortRange) []types.HostPortRange {
	var kokiHostPorts []types.HostPortRange

	for i := range kubeHostPorts {
		kubeHostPort := kubeHostPorts[i]

		kokiHostPort := types.HostPortRange{
			Min: kubeHostPort.Min,
			Max: kubeHostPort.Max,
		}

		kokiHostPorts = append(kokiHostPorts, kokiHostPort)
	}

	return kokiHostPorts
}

func convertPodSecurityPolicyCapabilities(kubeCaps []v1.Capability) []string {
	var kokiCaps []string

	for i := range kubeCaps {
		kubeCap := kubeCaps[i]

		kokiCaps = append(kokiCaps, string(kubeCap))
	}

	return kokiCaps
}

func convertPodSecurityPolicyVolumePlugins(kubeVolPlugins []exts.FSType) ([]string, error) {
	var kokiVolPlugins []string

	for i := range kubeVolPlugins {
		kubeVolPlugin := kubeVolPlugins[i]

		kokiVolPlugin, err := convertPodSecurityPolicyVolumePlugin(kubeVolPlugin)
		if err != nil {
			return nil, err
		}
		kokiVolPlugins = append(kokiVolPlugins, kokiVolPlugin)
	}

	return kokiVolPlugins, nil
}

func convertPodSecurityPolicyVolumePlugin(kubeVolPlugin exts.FSType) (string, error) {
	if len(kubeVolPlugin) == 0 {
		return "", nil
	}

	switch kubeVolPlugin {
	case exts.AzureFile:
		return types.VolumeTypeAzureFile, nil
	case exts.Flocker:
		return types.VolumeTypeFlocker, nil
	case exts.FlexVolume:
		return types.VolumeTypeFlex, nil
	case exts.HostPath:
		return types.VolumeTypeHostPath, nil
	case exts.EmptyDir:
		return types.VolumeTypeEmptyDir, nil
	case exts.GCEPersistentDisk:
		return types.VolumeTypeGcePD, nil
	case exts.AWSElasticBlockStore:
		return types.VolumeTypeAwsEBS, nil
	case exts.GitRepo:
		return types.VolumeTypeGit, nil
	case exts.Secret:
		return types.VolumeTypeSecret, nil
	case exts.NFS:
		return types.VolumeTypeNFS, nil
	case exts.ISCSI:
		return types.VolumeTypeISCSI, nil
	case exts.Glusterfs:
		return types.VolumeTypeGlusterfs, nil
	case exts.PersistentVolumeClaim:
		return types.VolumeTypePVC, nil
	case exts.RBD:
		return types.VolumeTypeRBD, nil
	case exts.Cinder:
		return types.VolumeTypeCinder, nil
	case exts.CephFS:
		return types.VolumeTypeCephFS, nil
	case exts.DownwardAPI:
		return types.VolumeTypeDownwardAPI, nil
	case exts.FC:
		return types.VolumeTypeFibreChannel, nil
	case exts.ConfigMap:
		return types.VolumeTypeConfigMap, nil
	case exts.Quobyte:
		return types.VolumeTypeQuobyte, nil
	case exts.AzureDisk:
		return types.VolumeTypeAzureDisk, nil
	case exts.All:
		return types.VolumeTypeAny, nil
	}

	return "", fmt.Errorf("Invalid Kube Volume Plugin")
}
