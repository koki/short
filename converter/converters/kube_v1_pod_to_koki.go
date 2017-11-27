package converters

import (
	"fmt"
	"net/url"
	"strings"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/short/parser/expressions"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
	"github.com/koki/short/util/floatstr"
)

func Convert_Kube_v1_Pod_to_Koki_Pod(pod *v1.Pod) (*types.PodWrapper, error) {
	var err error
	kokiPod := &types.Pod{}

	kokiPod.Name = pod.Name
	kokiPod.Namespace = pod.Namespace
	kokiPod.Version = pod.APIVersion
	kokiPod.Cluster = pod.ClusterName
	kokiPod.Labels = pod.Labels
	kokiPod.Annotations = pod.Annotations

	kokiPod.Volumes, err = convertVolumes(pod.Spec.Volumes)
	if err != nil {
		return nil, util.ContextualizeErrorf(err, "pod volumes")
	}
	affinity, err := convertAffinity(pod.Spec)
	if err != nil {
		return nil, err
	}
	kokiPod.Affinity = affinity

	var initContainers []types.Container
	for i := range pod.Spec.InitContainers {
		container := pod.Spec.InitContainers[i]
		kokiContainer, err := convertContainer(&container)
		if err != nil {
			return nil, err
		}
		initContainers = append(initContainers, *kokiContainer)
	}
	kokiPod.InitContainers = initContainers

	var kokiContainers []types.Container
	for i := range pod.Spec.Containers {
		container := pod.Spec.Containers[i]
		kokiContainer, err := convertContainer(&container)
		if err != nil {
			return nil, err
		}
		kokiContainers = append(kokiContainers, *kokiContainer)
	}
	kokiPod.Containers = kokiContainers

	dnsPolicy, err := convertDNSPolicy(pod.Spec.DNSPolicy)
	if err != nil {
		return nil, err
	}
	kokiPod.DNSPolicy = dnsPolicy

	kokiPod.HostAliases = convertHostAliases(pod.Spec.HostAliases)
	kokiPod.HostMode = convertHostMode(pod.Spec)
	kokiPod.Hostname = convertHostname(pod.Spec)
	kokiPod.Registries = convertRegistries(pod.Spec.ImagePullSecrets)

	restartPolicy, err := convertRestartPolicy(pod.Spec.RestartPolicy)
	if err != nil {
		return nil, err
	}
	kokiPod.RestartPolicy = restartPolicy

	kokiPod.SchedulerName = pod.Spec.SchedulerName
	kokiPod.Account = pod.Spec.ServiceAccountName

	if pod.Spec.AutomountServiceAccountToken != nil && *pod.Spec.AutomountServiceAccountToken {
		kokiPod.Account = fmt.Sprintf("%s:auto", kokiPod.Account)
	}

	tolerations, err := convertTolerations(pod.Spec.Tolerations)
	if err != nil {
		return nil, err
	}
	kokiPod.Tolerations = tolerations

	kokiPod.TerminationGracePeriod = pod.Spec.TerminationGracePeriodSeconds
	kokiPod.ActiveDeadline = pod.Spec.ActiveDeadlineSeconds
	kokiPod.Node = pod.Spec.NodeName
	kokiPod.Priority = convertPriority(pod.Spec)

	if pod.Spec.SecurityContext != nil {
		securityContext := pod.Spec.SecurityContext
		kokiPod.GIDs = securityContext.SupplementalGroups
		kokiPod.FSGID = securityContext.FSGroup
		for i := range kokiPod.Containers {
			container := kokiPod.Containers[i]
			if container.SELinux == nil {
				container.SELinux = convertSELinux(securityContext.SELinuxOptions)
			}
			if container.UID == nil {
				container.UID = securityContext.RunAsUser
			}
			if container.ForceNonRoot == nil {
				container.ForceNonRoot = securityContext.RunAsNonRoot
			}
		}
	}

	kokiPod.Msg = pod.Status.Message
	kokiPod.Reason = pod.Status.Reason
	phase, err := convertPhase(pod.Status.Phase)
	if err != nil {
		return nil, err
	}
	kokiPod.Phase = phase
	kokiPod.IP = pod.Status.PodIP
	kokiPod.NodeIP = pod.Status.HostIP
	kokiPod.StartTime = pod.Status.StartTime

	qosClass, err := convertPodQOSClass(pod.Status.QOSClass)
	if err != nil {
		return nil, err
	}
	kokiPod.QOS = qosClass

	conditions, err := convertPodConditions(pod.Status.Conditions)
	if err != nil {
		return nil, err
	}
	kokiPod.Conditions = conditions

	err = convertContainerStatuses(pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses, kokiPod.Containers)
	if err != nil {
		return nil, err
	}

	return &types.PodWrapper{Pod: *kokiPod}, nil
}

func convertVolumes(kubeVolumes []v1.Volume) (map[string]types.Volume, error) {
	kokiVolumes := map[string]types.Volume{}
	for _, kubeVolume := range kubeVolumes {
		name, kokiVolume, err := convertVolume(kubeVolume)
		if err != nil {
			return nil, util.ContextualizeErrorf(err, "volume (%s)", name)
		}
		kokiVolumes[name] = *kokiVolume
	}

	return kokiVolumes, nil
}

func convertStorageMedium(kubeMedium v1.StorageMedium) (types.StorageMedium, error) {
	switch kubeMedium {
	case v1.StorageMediumDefault:
		return types.StorageMediumDefault, nil
	case v1.StorageMediumMemory:
		return types.StorageMediumMemory, nil
	case v1.StorageMediumHugepages:
		return types.StorageMediumHugepages, nil
	default:
		return types.StorageMediumDefault, util.InvalidValueErrorf(kubeMedium, "unrecognized storage medium")
	}
}

func convertHostPathType(kubeType *v1.HostPathType) (types.HostPathType, error) {
	if kubeType == nil {
		return types.HostPathUnset, nil
	}

	switch *kubeType {
	case v1.HostPathUnset:
		return types.HostPathUnset, nil
	case v1.HostPathDirectoryOrCreate:
		return types.HostPathDirectoryOrCreate, nil
	case v1.HostPathDirectory:
		return types.HostPathDirectory, nil
	case v1.HostPathFileOrCreate:
		return types.HostPathFileOrCreate, nil
	case v1.HostPathFile:
		return types.HostPathFile, nil
	case v1.HostPathSocket:
		return types.HostPathSocket, nil
	case v1.HostPathCharDev:
		return types.HostPathCharDev, nil
	case v1.HostPathBlockDev:
		return types.HostPathBlockDev, nil
	default:
		return types.HostPathUnset, util.InvalidValueErrorf(kubeType, "unrecognized host_path type")
	}
}

func convertAzureDiskKind(kubeKind *v1.AzureDataDiskKind) (*types.AzureDataDiskKind, error) {
	if kubeKind == nil {
		return nil, nil
	}

	var kind types.AzureDataDiskKind
	switch *kubeKind {
	case v1.AzureDedicatedBlobDisk:
		kind = types.AzureDedicatedBlobDisk
	case v1.AzureSharedBlobDisk:
		kind = types.AzureSharedBlobDisk
	case v1.AzureManagedDisk:
		kind = types.AzureManagedDisk
	default:
		return nil, util.InvalidValueErrorf(kubeKind, "unrecognized kind")
	}

	return &kind, nil
}

func convertAzureDiskCachingMode(kubeMode *v1.AzureDataDiskCachingMode) (*types.AzureDataDiskCachingMode, error) {
	if kubeMode == nil {
		return nil, nil
	}

	var mode types.AzureDataDiskCachingMode
	switch *kubeMode {
	case v1.AzureDataDiskCachingNone:
		mode = types.AzureDataDiskCachingNone
	case v1.AzureDataDiskCachingReadOnly:
		mode = types.AzureDataDiskCachingReadOnly
	case v1.AzureDataDiskCachingReadWrite:
		mode = types.AzureDataDiskCachingReadWrite
	default:
		return nil, util.InvalidValueErrorf(kubeMode, "unrecognized cache")
	}

	return &mode, nil
}

func convertCephFSSecretFileOrRef(kubeFile string, kubeRef *v1.LocalObjectReference) *types.CephFSSecretFileOrRef {
	if len(kubeFile) > 0 {
		return &types.CephFSSecretFileOrRef{
			File: kubeFile,
		}
	}

	if kubeRef != nil {
		return &types.CephFSSecretFileOrRef{
			Ref: kubeRef.Name,
		}
	}

	return nil
}

func convertLocalObjectRef(kubeRef *v1.LocalObjectReference) string {
	if kubeRef == nil {
		return ""
	}

	return kubeRef.Name
}

func convertVolume(kubeVolume v1.Volume) (string, *types.Volume, error) {
	name := kubeVolume.Name
	if kubeVolume.VolumeSource.EmptyDir != nil {
		medium, err := convertStorageMedium(kubeVolume.VolumeSource.EmptyDir.Medium)
		if err != nil {
			return name, nil, err
		}
		return name, &types.Volume{
			EmptyDir: &types.EmptyDirVolume{
				Medium:    medium,
				SizeLimit: kubeVolume.VolumeSource.EmptyDir.SizeLimit,
			},
		}, nil
	}
	if kubeVolume.VolumeSource.HostPath != nil {
		kokiType, err := convertHostPathType(kubeVolume.VolumeSource.HostPath.Type)
		if err != nil {
			return name, nil, util.ContextualizeErrorf(err, "volume (%s)", name)
		}
		return name, &types.Volume{
			HostPath: &types.HostPathVolume{
				Path: kubeVolume.VolumeSource.HostPath.Path,
				Type: kokiType,
			},
		}, nil
	}
	if kubeVolume.VolumeSource.GCEPersistentDisk != nil {
		source := kubeVolume.VolumeSource.GCEPersistentDisk
		return name, &types.Volume{
			GcePD: &types.GcePDVolume{
				PDName:    source.PDName,
				FSType:    source.FSType,
				Partition: source.Partition,
				ReadOnly:  source.ReadOnly,
			},
		}, nil
	}
	if kubeVolume.VolumeSource.AWSElasticBlockStore != nil {
		source := kubeVolume.VolumeSource.AWSElasticBlockStore
		return name, &types.Volume{
			AwsEBS: &types.AwsEBSVolume{
				VolumeID:  source.VolumeID,
				FSType:    source.FSType,
				Partition: source.Partition,
				ReadOnly:  source.ReadOnly,
			},
		}, nil
	}
	if kubeVolume.VolumeSource.AzureDisk != nil {
		source := kubeVolume.VolumeSource.AzureDisk
		fstype := util.FromStringPtr(source.FSType)
		readOnly := util.FromBoolPtr(source.ReadOnly)
		kind, err := convertAzureDiskKind(source.Kind)
		if err != nil {
			return name, nil, err
		}
		cachingMode, err := convertAzureDiskCachingMode(source.CachingMode)
		if err != nil {
			return name, nil, err
		}
		return name, &types.Volume{
			AzureDisk: &types.AzureDiskVolume{
				DiskName:    source.DiskName,
				DataDiskURI: source.DataDiskURI,
				FSType:      fstype,
				ReadOnly:    readOnly,
				Kind:        kind,
				CachingMode: cachingMode,
			},
		}, nil
	}
	if kubeVolume.VolumeSource.AzureFile != nil {
		source := kubeVolume.VolumeSource.AzureFile
		return name, &types.Volume{
			AzureFile: &types.AzureFileVolume{
				SecretName: source.SecretName,
				ShareName:  source.ShareName,
				ReadOnly:   source.ReadOnly,
			},
		}, nil
	}
	if kubeVolume.VolumeSource.CephFS != nil {
		source := kubeVolume.VolumeSource.CephFS
		secretFileOrRef := convertCephFSSecretFileOrRef(source.SecretFile, source.SecretRef)
		return name, &types.Volume{
			CephFS: &types.CephFSVolume{
				Monitors:        source.Monitors,
				Path:            source.Path,
				User:            source.User,
				SecretFileOrRef: secretFileOrRef,
				ReadOnly:        source.ReadOnly,
			},
		}, nil
	}
	if kubeVolume.VolumeSource.Cinder != nil {
		source := kubeVolume.VolumeSource.Cinder
		return name, &types.Volume{
			Cinder: &types.CinderVolume{
				VolumeID: source.VolumeID,
				FSType:   source.FSType,
				ReadOnly: source.ReadOnly,
			},
		}, nil
	}
	if kubeVolume.VolumeSource.FC != nil {
		source := kubeVolume.VolumeSource.FC
		return name, &types.Volume{
			FibreChannel: &types.FibreChannelVolume{
				TargetWWNs: source.TargetWWNs,
				Lun:        source.Lun,
				ReadOnly:   source.ReadOnly,
				WWIDs:      source.WWIDs,
			},
		}, nil
	}
	if kubeVolume.VolumeSource.FlexVolume != nil {
		source := kubeVolume.VolumeSource.FlexVolume
		return name, &types.Volume{
			Flex: &types.FlexVolume{
				Driver:    source.Driver,
				FSType:    source.FSType,
				SecretRef: convertLocalObjectRef(source.SecretRef),
				ReadOnly:  source.ReadOnly,
				Options:   source.Options,
			},
		}, nil
	}
	if kubeVolume.VolumeSource.Flocker != nil {
		source := kubeVolume.VolumeSource.Flocker
		var dataset string
		if len(source.DatasetUUID) > 0 {
			dataset = source.DatasetUUID
		} else {
			dataset = source.DatasetName
		}
		return name, &types.Volume{
			Flocker: &types.FlockerVolume{
				DatasetUUID: dataset,
			},
		}, nil
	}

	return name, nil, util.InvalidInstanceErrorf(kubeVolume, "empty volume definition")
}

func convertContainer(container *v1.Container) (*types.Container, error) {
	kokiContainer := &types.Container{}

	kokiContainer.Name = container.Name
	kokiContainer.Command = container.Command
	kokiContainer.Image = container.Image
	kokiContainer.Args = convertContainerArgs(container.Args)
	kokiContainer.WorkingDir = container.WorkingDir

	pullPolicy, err := convertPullPolicy(container.ImagePullPolicy)
	if err != nil {
		return nil, err
	}
	kokiContainer.Pull = pullPolicy

	onStart, preStop, err := convertLifecycle(container.Lifecycle)
	if err != nil {
		return nil, err
	}
	kokiContainer.OnStart = onStart
	kokiContainer.PreStop = preStop

	kokiContainer.CPU = convertCPU(container.Resources)
	kokiContainer.Mem = convertMem(container.Resources)

	if container.SecurityContext != nil {
		kokiContainer.Privileged = container.SecurityContext.Privileged
		kokiContainer.AllowEscalation = container.SecurityContext.AllowPrivilegeEscalation
		if container.SecurityContext.ReadOnlyRootFilesystem != nil {
			kokiContainer.RO = container.SecurityContext.ReadOnlyRootFilesystem
			rw := !(*kokiContainer.RO)
			kokiContainer.RW = &rw
		}
		kokiContainer.ForceNonRoot = container.SecurityContext.RunAsNonRoot
		kokiContainer.UID = container.SecurityContext.RunAsUser
		kokiContainer.SELinux = convertSELinux(container.SecurityContext.SELinuxOptions)
		kokiContainer.AddCapabilities = convertCapabilitiesAdds(container.SecurityContext.Capabilities)
		kokiContainer.DelCapabilities = convertCapabilitiesDels(container.SecurityContext.Capabilities)
	}

	livenessProbe, err := convertProbe(container.LivenessProbe)
	if err != nil {
		return nil, err
	}
	kokiContainer.LivenessProbe = livenessProbe

	readinessProbe, err := convertProbe(container.ReadinessProbe)
	if err != nil {
		return nil, err
	}
	kokiContainer.ReadinessProbe = readinessProbe

	ports, err := convertContainerPorts(container.Ports)
	if err != nil {
		return nil, err
	}

	kokiContainer.Expose = ports

	kokiContainer.Stdin = container.Stdin
	kokiContainer.StdinOnce = container.StdinOnce
	kokiContainer.TTY = container.TTY

	kokiContainer.TerminationMsgPath = container.TerminationMessagePath

	policy, err := convertTerminationMsgPolicy(container.TerminationMessagePolicy)
	if err != nil {
		return nil, err
	}
	kokiContainer.TerminationMsgPolicy = policy

	kokiContainer.Env = convertEnvVars(container.Env, container.EnvFrom)

	volumeMounts, err := convertVolumeMounts(container.VolumeMounts)
	if err != nil {
		return nil, err
	}

	kokiContainer.VolumeMounts = volumeMounts

	return kokiContainer, nil
}

func convertContainerArgs(kubeArgs []string) []floatstr.FloatOrString {
	if kubeArgs == nil {
		return nil
	}
	kokiArgs := make([]floatstr.FloatOrString, len(kubeArgs))
	for i, kubeArg := range kubeArgs {
		kokiArgs[i] = *floatstr.Parse(kubeArg)
	}

	return kokiArgs
}

func convertPullPolicy(pullPolicy v1.PullPolicy) (types.PullPolicy, error) {
	if pullPolicy == "" {
		return "", nil
	}
	if pullPolicy == v1.PullAlways {
		return types.PullAlways, nil
	}
	if pullPolicy == v1.PullNever {
		return types.PullNever, nil
	}
	if pullPolicy == v1.PullIfNotPresent {
		return types.PullNever, nil
	}
	return "", util.InvalidInstanceError(pullPolicy)
}

func convertLifecycle(lifecycle *v1.Lifecycle) (onStart *types.Action, preStop *types.Action, err error) {
	if lifecycle == nil {
		return nil, nil, nil
	}

	actionOnStart, err := convertLifecycleAction(lifecycle.PostStart)
	if err != nil {
		return nil, nil, err
	}
	onStart = actionOnStart

	actionPreStop, err := convertLifecycleAction(lifecycle.PreStop)
	if err != nil {
		return nil, nil, err
	}
	preStop = actionPreStop

	return onStart, preStop, nil
}

func convertLifecycleAction(lcHandler *v1.Handler) (*types.Action, error) {
	if lcHandler == nil {
		return nil, nil
	}
	var act *types.Action
	ps := lcHandler
	if ps.Exec != nil {
		act = &types.Action{}
		act.Command = ps.Exec.Command
	}
	if ps.HTTPGet != nil {
		if act == nil {
			act = &types.Action{}
			scheme := "HTTP"
			hostPort := ""
			if ps.HTTPGet.Scheme != "" {
				scheme = string(ps.HTTPGet.Scheme)
			}

			if ps.HTTPGet.Port.String() == "" {
				return nil, util.InvalidInstanceErrorf(ps, "URL Port is missing")
			}

			host := "localhost"
			if ps.HTTPGet.Host != "" {
				host = ps.HTTPGet.Host
			}
			port := "80"
			if ps.HTTPGet.Port.String() != "" {
				port = ps.HTTPGet.Port.String()
			}
			hostPort = fmt.Sprintf("%s:%s", host, port)

			var headers []string

			if ps.HTTPGet.HTTPHeaders != nil {
				headers = []string{}
				for i := range ps.HTTPGet.HTTPHeaders {
					inHeader := ps.HTTPGet.HTTPHeaders[i]
					outHeader := fmt.Sprintf("%s:%s", inHeader.Name, inHeader.Value)
					headers = append(headers, outHeader)
				}
			}

			url := &url.URL{
				Scheme: scheme,
				Host:   hostPort,
				Path:   ps.HTTPGet.Path,
			}
			act.Net = &types.NetAction{
				URL: url.String(),

				Headers: headers,
			}
		}
	}
	if ps.TCPSocket != nil {
		if act == nil {
			url := &url.URL{
				Scheme: "TCP",
				Host:   fmt.Sprintf("%s:%s", ps.TCPSocket.Host, ps.TCPSocket.Port.String()),
			}
			act = &types.Action{
				Net: &types.NetAction{
					URL: url.String(),
				},
			}
		}
	}
	return act, nil
}

func convertCPU(resources v1.ResourceRequirements) *types.CPU {
	cpu := &types.CPU{}
	mark := false
	if resources.Limits != nil {
		max := ""
		if q, ok := resources.Limits["cpu"]; ok {
			mark = true
			max = q.String()
		}
		cpu.Max = max
	}
	if resources.Requests != nil {
		min := ""
		if q, ok := resources.Requests["cpu"]; ok {
			mark = true
			min = q.String()
		}
		cpu.Min = min
	}
	if mark {
		return cpu
	}
	return nil
}

func convertMem(resources v1.ResourceRequirements) *types.Mem {
	mem := &types.Mem{}
	mark := false
	if resources.Limits != nil {
		max := ""
		if q, ok := resources.Limits["memory"]; ok {
			mark = true
			max = q.String()
		}
		mem.Max = max
	}
	if resources.Requests != nil {
		min := ""
		if q, ok := resources.Requests["memory"]; ok {
			mark = true
			min = q.String()
		}
		mem.Min = min
	}
	if mark {
		return mem
	}
	return nil
}

func convertSELinux(opts *v1.SELinuxOptions) *types.SELinux {
	if opts == nil {
		return nil
	}
	return &types.SELinux{
		User:  opts.User,
		Level: opts.Level,
		Role:  opts.Role,
		Type:  opts.Type,
	}
}

func convertCapabilitiesAdds(caps *v1.Capabilities) []string {
	if caps == nil {
		return nil
	}
	var kokiCaps []string
	if caps.Add != nil {
		for i := range caps.Add {
			cap := string(caps.Add[i])
			kokiCaps = append(kokiCaps, cap)
		}
	}
	return kokiCaps
}

func convertCapabilitiesDels(caps *v1.Capabilities) []string {
	if caps == nil {
		return nil
	}
	var kokiCaps []string
	if caps.Drop != nil {
		for i := range caps.Drop {
			cap := string(caps.Drop[i])
			kokiCaps = append(kokiCaps, cap)
		}
	}
	return kokiCaps

}

func convertProbe(probe *v1.Probe) (*types.Probe, error) {
	if probe == nil {
		return nil, nil
	}

	action, err := convertLifecycleAction(&probe.Handler)
	if err != nil {
		return nil, err
	}

	p := &types.Probe{
		Action: *action,
	}
	p.Delay = probe.InitialDelaySeconds
	p.MinCountSuccess = probe.SuccessThreshold
	p.MinCountFailure = probe.FailureThreshold
	p.Interval = probe.PeriodSeconds
	p.Timeout = probe.TimeoutSeconds

	return p, nil
}

func convertContainerPorts(ports []v1.ContainerPort) ([]types.Port, error) {
	if ports == nil {
		return nil, nil
	}

	var p []types.Port
	for i := range ports {
		port := ports[i]
		kokiPort := types.Port{}

		kokiPort.Name = port.Name
		kokiPort.Protocol = convertProtocol(port.Protocol)
		kokiPort.IP = port.HostIP
		if port.HostPort != 0 {
			kokiPort.HostPort = fmt.Sprintf("%d", port.HostPort)
		}
		if port.ContainerPort != 0 {
			kokiPort.ContainerPort = fmt.Sprintf("%d", port.ContainerPort)
		}
		p = append(p, kokiPort)
	}
	return p, nil
}

func convertProtocol(kubeProtocol v1.Protocol) types.Protocol {
	return types.Protocol(strings.ToLower(string(kubeProtocol)))
}

func convertTerminationMsgPolicy(p v1.TerminationMessagePolicy) (types.TerminationMessagePolicy, error) {
	if p == "" {
		return "", nil
	}
	if p == v1.TerminationMessageReadFile {
		return types.TerminationMessageReadFile, nil
	}
	if p == v1.TerminationMessageFallbackToLogsOnError {
		return types.TerminationMessageFallbackToLogsOnError, nil
	}
	return "", util.InvalidInstanceError(p)
}

func convertEnvVars(env []v1.EnvVar, envFromSrc []v1.EnvFromSource) []types.Env {
	var kokiEnvs []types.Env
	for i := range env {
		v := env[i]
		if v.ValueFrom == nil {
			kokiEnvs = append(kokiEnvs, types.EnvWithVal(types.EnvVal{
				Key: v.Name,
				Val: v.Value,
			}))
			continue
		}

		e := types.EnvFrom{}
		e.Key = v.Name
		if v.ValueFrom.FieldRef != nil {
			e.From = v.ValueFrom.FieldRef.FieldPath
		}
		if v.ValueFrom.ResourceFieldRef != nil {
			//This might be losing some information
			e.From = v.ValueFrom.ResourceFieldRef.Resource
		}
		if v.ValueFrom.ConfigMapKeyRef != nil {
			e.From = fmt.Sprintf("config:%s:%s", v.ValueFrom.ConfigMapKeyRef.Name, v.ValueFrom.ConfigMapKeyRef.Key)
			required := v.ValueFrom.ConfigMapKeyRef.Optional
			e.Required = required
		}
		if v.ValueFrom.SecretKeyRef != nil {
			e.From = fmt.Sprintf("secret:%s:%s", v.ValueFrom.SecretKeyRef.Name, v.ValueFrom.SecretKeyRef.Key)
			required := v.ValueFrom.SecretKeyRef.Optional
			e.Required = required
		}
		kokiEnvs = append(kokiEnvs, types.EnvWithFrom(e))
	}
	for i := range envFromSrc {
		v := envFromSrc[i]
		e := types.EnvFrom{}
		e.Key = v.Prefix
		if v.ConfigMapRef != nil {
			e.From = fmt.Sprintf("config:%s", v.ConfigMapRef.Name)
			required := v.ConfigMapRef.Optional
			e.Required = required
		}
		if v.SecretRef != nil {
			e.From = fmt.Sprintf("secret:%s", v.SecretRef.Name)
			required := v.SecretRef.Optional
			e.Required = required
		}
		kokiEnvs = append(kokiEnvs, types.EnvWithFrom(e))
	}
	return kokiEnvs
}

func convertVolumeMounts(mounts []v1.VolumeMount) ([]types.VolumeMount, error) {
	var kokiMounts []types.VolumeMount
	for i := range mounts {
		mount := mounts[i]
		km := types.VolumeMount{
			MountPath: mount.MountPath,
		}
		if mount.MountPropagation != nil {
			propagation, err := convertMountPropagation(*mount.MountPropagation)
			if err != nil {
				return nil, err
			}
			km.Propagation = propagation
		}
		access := "rw"
		if mount.ReadOnly {
			access = "ro"
		}
		trailer := ""
		if mount.SubPath == "" {
			if access == "ro" {
				trailer = fmt.Sprintf(access)
			}
		} else {
			trailer = fmt.Sprintf("%s", mount.SubPath)
			if access == "ro" {
				trailer = fmt.Sprintf("%s:%s", trailer, access)
			}
		}
		if trailer != "" {
			km.Store = fmt.Sprintf("%s:%s", mount.Name, trailer)
		} else {
			km.Store = mount.Name
		}
		kokiMounts = append(kokiMounts, km)
	}
	return kokiMounts, nil
}

func convertMountPropagation(p v1.MountPropagationMode) (types.MountPropagation, error) {
	if p == "" {
		return "", nil
	}
	if p == v1.MountPropagationHostToContainer {
		return types.MountPropagationHostToContainer, nil
	} else if p == v1.MountPropagationBidirectional {
		return types.MountPropagationBidirectional, nil
	}
	return "", util.InvalidInstanceError(p)
}

func convertAffinity(spec v1.PodSpec) ([]types.Affinity, error) {
	var affinity []types.Affinity
	affinityString := ""
	for k, v := range spec.NodeSelector {
		expr := fmt.Sprintf("%s=%s", k, v)
		if affinityString == "" {
			affinityString = fmt.Sprintf("node:%s", expr)
			continue
		}
		affinityString = fmt.Sprintf("%s&%s", affinityString, expr)
	}

	if affinityString != "" {
		affinity = append(affinity, types.Affinity{NodeAffinity: affinityString})
	}

	if spec.Affinity == nil {
		return affinity, nil
	}

	// Node Affinity
	nodeAffinity, err := convertNodeAffinity(spec.Affinity.NodeAffinity)
	if err != nil {
		return nil, err
	}

	affinity = append(affinity, nodeAffinity...)

	// Pod affinity
	podAffinity, err := convertPodAffinity(spec.Affinity.PodAffinity)
	if err != nil {
		return nil, err
	}

	// Pod Anti Affinity
	podAntiAffinity, err := convertPodAntiAffinity(spec.Affinity.PodAntiAffinity)
	if err != nil {
		return nil, err
	}

	affinity = append(affinity, podAffinity...)
	affinity = append(affinity, podAntiAffinity...)

	return affinity, nil
}

func convertNodeAffinity(nodeAffinity *v1.NodeAffinity) ([]types.Affinity, error) {
	if nodeAffinity == nil {
		return nil, nil
	}

	var affinity []types.Affinity
	if nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution != nil {
		nodeHardAffinity := nodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution
		for i := range nodeHardAffinity.NodeSelectorTerms {
			selectorTerm := nodeHardAffinity.NodeSelectorTerms[i]
			affinityString := ""
			for i := range selectorTerm.MatchExpressions {

				expr := selectorTerm.MatchExpressions[i]
				value := strings.Join(expr.Values, ",")
				op, err := convertOperator(expr.Operator)
				if err != nil {
					return nil, util.InvalidInstanceErrorf(nodeHardAffinity, "unsupported Operator: %s", err.Error())
				}
				kokiExpr := fmt.Sprintf("%s%s%s", expr.Key, op, value)
				if expr.Operator == v1.NodeSelectorOpExists {
					kokiExpr = fmt.Sprintf("%s", expr.Key)
				}
				if expr.Operator == v1.NodeSelectorOpDoesNotExist {
					kokiExpr = fmt.Sprintf("!%s", expr.Key)
				}
				if len(affinityString) == 0 {
					affinityString = kokiExpr
					continue
				}
				affinityString = fmt.Sprintf("%s&%s", affinityString, kokiExpr)
			}
			if len(affinityString) > 0 {
				affinity = append(affinity, types.Affinity{NodeAffinity: affinityString})
			}
		}
	}

	// Node soft affinities
	if nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution != nil {
		nodeSoftAffinity := nodeAffinity.PreferredDuringSchedulingIgnoredDuringExecution
		for i := range nodeSoftAffinity {
			selectorTerm := nodeSoftAffinity[i]
			affinityString := ""
			weight := selectorTerm.Weight
			for i := range selectorTerm.Preference.MatchExpressions {
				expr := selectorTerm.Preference.MatchExpressions[i]
				value := strings.Join(expr.Values, ",")
				op, err := convertOperator(expr.Operator)
				if err != nil {
					return nil, util.InvalidInstanceErrorf(nodeSoftAffinity, "unsupported Operator: %s", err.Error())
				}
				kokiExpr := fmt.Sprintf("%s%s%s", expr.Key, op, value)
				if expr.Operator == v1.NodeSelectorOpExists {
					kokiExpr = fmt.Sprintf("%s", expr.Key)
				}
				if expr.Operator == v1.NodeSelectorOpDoesNotExist {
					kokiExpr = fmt.Sprintf("!%s", expr.Key)
				}
				if len(affinityString) == 0 {
					affinityString = kokiExpr
					continue
				}
				affinityString = fmt.Sprintf("%s&%s", affinityString, kokiExpr)
			}
			if len(affinityString) > 0 {
				affinityString = fmt.Sprintf("%s:soft", affinityString)
				// The default value for Weight is 1. 0 means "unspecified".
				if weight != 0 {
					affinityString = fmt.Sprintf("%s:%d", affinityString, weight)
				}
				affinity = append(affinity, types.Affinity{NodeAffinity: affinityString})
			}
		}
	}
	return affinity, nil
}

func convertPodAffinity(podAffinity *v1.PodAffinity) ([]types.Affinity, error) {
	if podAffinity == nil {
		return nil, nil
	}

	hardAffinity, err := convertPodAffinityTerms(false, podAffinity.RequiredDuringSchedulingIgnoredDuringExecution)
	if err != nil {
		return nil, err
	}

	softAffinity, err := convertPodWeightedAffinityTerms(false, podAffinity.PreferredDuringSchedulingIgnoredDuringExecution)
	if err != nil {
		return nil, err
	}

	return append(hardAffinity, softAffinity...), nil
}

func convertPodAntiAffinity(podAntiAffinity *v1.PodAntiAffinity) ([]types.Affinity, error) {
	if podAntiAffinity == nil {
		return nil, nil
	}

	hardAffinity, err := convertPodAffinityTerms(true, podAntiAffinity.RequiredDuringSchedulingIgnoredDuringExecution)
	if err != nil {
		return nil, err
	}

	softAffinity, err := convertPodWeightedAffinityTerms(true, podAntiAffinity.PreferredDuringSchedulingIgnoredDuringExecution)
	if err != nil {
		return nil, err
	}

	return append(hardAffinity, softAffinity...), nil
}

func convertPodWeightedAffinityTerms(isAntiAffinity bool, podSoftAffinity []v1.WeightedPodAffinityTerm) ([]types.Affinity, error) {
	var affinity []types.Affinity
	// Pod soft affinity
	for i := range podSoftAffinity {
		selectorTerm := podSoftAffinity[i]
		weight := selectorTerm.Weight
		affinityString := ""
		if selectorTerm.PodAffinityTerm.LabelSelector != nil {
			// parse through match labels first
			for k, v := range selectorTerm.PodAffinityTerm.LabelSelector.MatchLabels {
				kokiExpr := fmt.Sprintf("%s=%s", k, v)
				if len(affinityString) == 0 {
					affinityString = kokiExpr
					continue
				}
				affinityString = fmt.Sprintf("%s&%s", affinityString, kokiExpr)
			}

			// parse through match expressions now
			for i := range selectorTerm.PodAffinityTerm.LabelSelector.MatchExpressions {
				expr := selectorTerm.PodAffinityTerm.LabelSelector.MatchExpressions[i]
				value := strings.Join(expr.Values, ",")
				op, err := expressions.ConvertOperatorLabelSelector(expr.Operator)
				if err != nil {
					return nil, util.InvalidInstanceErrorf(selectorTerm.PodAffinityTerm, "unsupported Operator: %s", err.Error())
				}
				kokiExpr := fmt.Sprintf("%s%s%s", expr.Key, op, value)
				if expr.Operator == metav1.LabelSelectorOpExists {
					kokiExpr = fmt.Sprintf("%s", expr.Key)
				}
				if expr.Operator == metav1.LabelSelectorOpDoesNotExist {
					kokiExpr = fmt.Sprintf("!%s", expr.Key)
				}
				if len(affinityString) == 0 {
					affinityString = kokiExpr
					continue
				}
				affinityString = fmt.Sprintf("%s&%s", affinityString, kokiExpr)
			}
		}
		if len(affinityString) > 0 {
			affinityString = fmt.Sprintf("%s:soft", affinityString)
			if weight != 0 {
				affinityString = fmt.Sprintf("%s:%d", affinityString, weight)
			}
			a := types.Affinity{
				PodAffinity: affinityString,
				Namespaces:  selectorTerm.PodAffinityTerm.Namespaces,
				Topology:    selectorTerm.PodAffinityTerm.TopologyKey,
			}
			if isAntiAffinity {
				a.PodAntiAffinity = a.PodAffinity
				a.PodAffinity = ""
			}
			affinity = append(affinity, a)
		}
	}
	return affinity, nil
}

func convertPodAffinityTerms(isAntiAffinity bool, podHardAffinity []v1.PodAffinityTerm) ([]types.Affinity, error) {
	var affinity []types.Affinity
	// Pod hard affinity
	for i := range podHardAffinity {
		selectorTerm := podHardAffinity[i]
		affinityString := ""

		if selectorTerm.LabelSelector != nil {
			// parse through match labels first
			for k, v := range selectorTerm.LabelSelector.MatchLabels {
				kokiExpr := fmt.Sprintf("%s=%s", k, v)
				if len(affinityString) == 0 {
					affinityString = kokiExpr
					continue
				}
				affinityString = fmt.Sprintf("%s&%s", affinityString, kokiExpr)
			}

			// parse through match expressions now
			for i := range selectorTerm.LabelSelector.MatchExpressions {
				expr := selectorTerm.LabelSelector.MatchExpressions[i]
				value := strings.Join(expr.Values, ",")
				op, err := expressions.ConvertOperatorLabelSelector(expr.Operator)
				if err != nil {
					return nil, util.InvalidInstanceErrorf(selectorTerm, "unsupported Operator: %s", err.Error())
				}
				kokiExpr := fmt.Sprintf("%s%s%s", expr.Key, op, value)
				if expr.Operator == metav1.LabelSelectorOpExists {
					kokiExpr = fmt.Sprintf("%s", expr.Key)
				}
				if expr.Operator == metav1.LabelSelectorOpDoesNotExist {
					kokiExpr = fmt.Sprintf("!%s", expr.Key)
				}
				if len(affinityString) == 0 {
					affinityString = kokiExpr
					continue
				}
				affinityString = fmt.Sprintf("%s&%s", affinityString, kokiExpr)
			}
		}

		if len(affinityString) > 0 {
			a := types.Affinity{
				PodAffinity: affinityString,
				Namespaces:  selectorTerm.Namespaces,
				Topology:    selectorTerm.TopologyKey,
			}
			if isAntiAffinity {
				a.PodAntiAffinity = a.PodAffinity
				a.PodAffinity = ""
			}
			affinity = append(affinity, a)
		}
	}
	return affinity, nil
}

func convertOperator(op v1.NodeSelectorOperator) (string, error) {
	if op == "" {
		return "", nil
	}
	if op == v1.NodeSelectorOpIn {
		return "=", nil
	}
	if op == v1.NodeSelectorOpNotIn {
		return "!=", nil
	}
	if op == v1.NodeSelectorOpExists {
		return "", nil
	}
	if op == v1.NodeSelectorOpDoesNotExist {
		return "", nil
	}
	if op == v1.NodeSelectorOpGt {
		return ">", nil
	}
	if op == v1.NodeSelectorOpLt {
		return "<", nil
	}
	return "", util.InvalidInstanceError(op)
}

func convertDNSPolicy(dnsPolicy v1.DNSPolicy) (types.DNSPolicy, error) {
	if dnsPolicy == "" {
		return "", nil
	}
	if dnsPolicy == v1.DNSClusterFirstWithHostNet {
		return types.DNSClusterFirstWithHostNet, nil
	}
	if dnsPolicy == v1.DNSClusterFirst {
		return types.DNSClusterFirst, nil
	}
	if dnsPolicy == v1.DNSDefault {
		return types.DNSDefault, nil
	}
	return "", util.InvalidInstanceError(dnsPolicy)
}

func convertHostAliases(aliases []v1.HostAlias) []string {
	var kokiAliases []string
	for i := range aliases {
		alias := aliases[i]
		aliasStr := fmt.Sprintf("%s", alias.IP)
		// Do not add empty/invalid entries
		if aliasStr == "" || len(alias.Hostnames) == 0 {
			continue
		}
		kokiAliases = append(kokiAliases, fmt.Sprintf("%s %s", aliasStr, strings.Join(alias.Hostnames, " ")))
	}
	return kokiAliases
}

func convertHostMode(spec v1.PodSpec) []types.HostMode {
	var hostMode []types.HostMode
	if spec.HostNetwork {
		hostMode = append(hostMode, types.HostModeNet)
	}
	if spec.HostPID {
		hostMode = append(hostMode, types.HostModePID)
	}
	if spec.HostIPC {
		hostMode = append(hostMode, types.HostModeIPC)
	}
	return hostMode

}

func convertHostname(spec v1.PodSpec) string {
	hostName := ""
	if spec.Hostname != "" {
		hostName = fmt.Sprintf("%s", spec.Hostname)
	}
	// TODO: verify that .subdomain is a valid input. i.e. without hostname
	if spec.Subdomain != "" {
		hostName = fmt.Sprintf("%s.%s", hostName, spec.Subdomain)
	}
	return hostName
}

func convertRegistries(ref []v1.LocalObjectReference) []string {
	var registries []string
	for i := range ref {
		r := ref[i]
		registries = append(registries, r.Name)
	}
	return registries
}

func convertRestartPolicy(policy v1.RestartPolicy) (types.RestartPolicy, error) {
	if policy == "" {
		return "", nil
	}
	if policy == v1.RestartPolicyAlways {
		return types.RestartPolicyAlways, nil
	}
	if policy == v1.RestartPolicyOnFailure {
		return types.RestartPolicyOnFailure, nil
	}
	if policy == v1.RestartPolicyNever {
		return types.RestartPolicyNever, nil
	}
	return "", util.InvalidInstanceError(policy)
}

func convertTolerations(tolerations []v1.Toleration) ([]types.Toleration, error) {
	var tols []types.Toleration
	for i := range tolerations {
		toleration := tolerations[i]
		tol := types.Toleration{}
		tol.ExpiryAfter = toleration.TolerationSeconds
		tolExpr := ""
		if toleration.Operator == v1.TolerationOpEqual {
			tolExpr = fmt.Sprintf("%s=%s", tolExpr, toleration.Value)
		} else if toleration.Operator == v1.TolerationOpExists {
			tolExpr = fmt.Sprintf("%s", toleration.Key)
		} else {
			return nil, util.InvalidInstanceErrorf(toleration, "unsupported operator")
		}
		if tolExpr != "" {
			if toleration.Effect != "" {
				tol.Selector = types.Selector(fmt.Sprintf("%s:%s", tolExpr, toleration.Effect))
			} else {
				tol.Selector = types.Selector(tolExpr)
			}
			tols = append(tols, tol)
		}
	}
	return tols, nil
}

func convertPriority(spec v1.PodSpec) *types.Priority {
	if spec.PriorityClassName == "" || spec.Priority == nil {
		return nil
	}
	return &types.Priority{
		Class: spec.PriorityClassName,
		Value: spec.Priority,
	}
}

func convertPhase(phase v1.PodPhase) (types.PodPhase, error) {
	if phase == "" {
		return "", nil
	}
	if phase == v1.PodPending {
		return types.PodPending, nil
	}
	if phase == v1.PodRunning {
		return types.PodRunning, nil
	}
	if phase == v1.PodSucceeded {
		return types.PodSucceeded, nil
	}
	if phase == v1.PodFailed {
		return types.PodFailed, nil
	}
	if phase == v1.PodUnknown {
		return types.PodUnknown, nil
	}
	return "", util.InvalidInstanceError(phase)
}

func convertPodQOSClass(class v1.PodQOSClass) (types.PodQOSClass, error) {
	if class == "" {
		return "", nil
	}
	if class == v1.PodQOSGuaranteed {
		return types.PodQOSGuaranteed, nil
	}
	if class == v1.PodQOSBurstable {
		return types.PodQOSBurstable, nil
	}
	if class == v1.PodQOSBestEffort {
		return types.PodQOSBestEffort, nil
	}
	return "", util.InvalidInstanceError(class)
}

func convertPodConditions(conditions []v1.PodCondition) ([]types.PodCondition, error) {
	var kConds []types.PodCondition
	for i := range conditions {
		cond := conditions[i]
		kCond := types.PodCondition{}
		typ, err := convertPodConditionType(cond.Type)
		if err != nil {
			return nil, err
		}
		kCond.Type = typ
		status, err := convertConditionStatus(cond.Status)
		if err != nil {
			return nil, err
		}
		kCond.Status = status
		kCond.Msg = cond.Message
		kCond.Reason = cond.Reason
		kCond.LastProbeTime = cond.LastProbeTime
		kCond.LastTransitionTime = cond.LastTransitionTime
		kConds = append(kConds, kCond)
	}
	return kConds, nil
}

func convertPodConditionType(typ v1.PodConditionType) (types.PodConditionType, error) {
	if typ == "" {
		return "", nil
	}
	if typ == v1.PodScheduled {
		return types.PodScheduled, nil
	}
	if typ == v1.PodReady {
		return types.PodReady, nil
	}
	if typ == v1.PodInitialized {
		return types.PodInitialized, nil
	}
	if typ == v1.PodReasonUnschedulable {
		return types.PodReasonUnschedulable, nil
	}
	return "", util.InvalidInstanceError(typ)
}

func convertConditionStatus(status v1.ConditionStatus) (types.ConditionStatus, error) {
	if status == "" {
		return "", nil
	}
	if status == v1.ConditionTrue {
		return types.ConditionTrue, nil
	}
	if status == v1.ConditionFalse {
		return types.ConditionFalse, nil
	}
	if status == v1.ConditionUnknown {
		return types.ConditionUnknown, nil
	}
	return "", util.InvalidInstanceError(status)
}

func convertContainerStatuses(initContainerStatuses, containerStatuses []v1.ContainerStatus, kokiContainers []types.Container) error {
	allContainerStatuses := append(initContainerStatuses, containerStatuses...)

	for i := range allContainerStatuses {
		status := allContainerStatuses[i]
		for i := range kokiContainers {
			container := kokiContainers[i]
			if container.Name == status.Name {
				container.Restarts = status.RestartCount
				container.Ready = status.Ready
				container.ImageID = status.ImageID
				container.ContainerID = status.ContainerID
				container.CurrentState = convertContainerState(status.State)
				container.LastState = convertContainerState(status.LastTerminationState)
			}
		}
	}
	return nil
}

func convertContainerState(state v1.ContainerState) *types.ContainerState {
	s := &types.ContainerState{}
	if state.Waiting != nil {
		s.Waiting = &types.ContainerStateWaiting{
			Reason: state.Waiting.Reason,
			Msg:    state.Waiting.Message,
		}
	}
	if state.Running != nil {
		s.Running = &types.ContainerStateRunning{
			StartTime: state.Running.StartedAt,
		}
	}
	if state.Terminated != nil {
		s.Terminated = &types.ContainerStateTerminated{
			StartTime:  state.Terminated.StartedAt,
			FinishTime: state.Terminated.FinishedAt,
			Reason:     state.Terminated.Reason,
			Msg:        state.Terminated.Message,
			Signal:     state.Terminated.Signal,
			ExitCode:   state.Terminated.ExitCode,
		}
	}
	return s
}
