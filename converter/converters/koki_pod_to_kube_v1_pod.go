package converters

import (
	"net/url"
	"strconv"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	intstr "k8s.io/apimachinery/pkg/util/intstr"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func Convert_Koki_Pod_to_Kube_v1_Pod(pod *types.PodWrapper) (*v1.Pod, error) {
	kubePod := &v1.Pod{}
	kokiPod := pod.Pod

	kubePod.Name = kokiPod.Name
	kubePod.Namespace = kokiPod.Namespace
	kubePod.APIVersion = kokiPod.Version
	kubePod.ClusterName = kokiPod.Cluster
	kubePod.Labels = kokiPod.Labels
	kubePod.Annotations = kokiPod.Annotations

	kubePod.Spec = v1.PodSpec{}

	fields := strings.SplitN(kokiPod.Hostname, ".", 2)
	if len(fields) == 1 {
		kubePod.Spec.Hostname = kokiPod.Hostname
	} else {
		kubePod.Spec.Hostname = fields[0]
		kubePod.Spec.Subdomain = fields[1]
	}

	var initContainers []v1.Container
	for i := range kokiPod.InitContainers {
		container := kokiPod.InitContainers[i]
		kubeContainer, err := revertKokiContainer(container)
		if err != nil {
			return nil, err
		}
		initContainers = append(initContainers, kubeContainer)
	}
	kubePod.Spec.InitContainers = initContainers

	var kubeContainers []v1.Container
	for i := range kokiPod.Containers {
		container := kokiPod.Containers[i]
		kubeContainer, err := revertKokiContainer(container)
		if err != nil {
			return nil, err
		}
		kubeContainers = append(kubeContainers, kubeContainer)
	}
	kubePod.Spec.Containers = kubeContainers

	hostAliases, err := revertHostAliases(kokiPod.HostAliases)
	if err != nil {
		return nil, err
	}
	kubePod.Spec.HostAliases = hostAliases

	return kubePod, nil
}

func revertHostAliases(aliases []string) ([]v1.HostAlias, error) {
	var hostAliases []v1.HostAlias
	for i := range aliases {
		alias := aliases[i]
		hostAlias := v1.HostAlias{}

		fields := strings.SplitN(alias, " ", 2)
		if len(fields) == 2 {
			hostAlias.IP = strings.TrimSpace(fields[0])
			hostNames := strings.Split(fields[1], " ")
			for i := range hostNames {
				hostAlias.Hostnames = append(hostAlias.Hostnames, hostNames[i])
			}
		} else {
			return nil, util.TypeValueErrorf(alias, "Unexpected value %s", alias)
		}
		hostAliases = append(hostAliases, hostAlias)
	}
	return hostAliases, nil
}

func revertKokiContainer(container types.Container) (v1.Container, error) {
	kubeContainer := v1.Container{}

	kubeContainer.Name = container.Name
	kubeContainer.Args = container.Args
	kubeContainer.Command = container.Command
	kubeContainer.Image = container.Image
	kubeContainer.WorkingDir = container.WorkingDir

	kubeContainerPorts, err := revertExpose(container.Expose)
	if err != nil {
		return v1.Container{}, err
	}
	kubeContainer.Ports = kubeContainerPorts

	envs, envFroms, err := revertEnv(container.Env)
	if err != nil {
		return v1.Container{}, err
	}
	kubeContainer.Env = envs
	kubeContainer.EnvFrom = envFroms

	resources, err := revertResources(container.CPU, container.Mem)
	if err != nil {
		return v1.Container{}, err
	}
	kubeContainer.Resources = resources

	livenessProbe, err := revertProbe(container.LivenessProbe)
	if err != nil {
		return v1.Container{}, err
	}
	kubeContainer.LivenessProbe = livenessProbe

	readinessProbe, err := revertProbe(container.ReadinessProbe)
	if err != nil {
		return v1.Container{}, err
	}
	kubeContainer.ReadinessProbe = readinessProbe

	kubeContainer.TerminationMessagePath = container.TerminationMsgPath
	kubeContainer.TerminationMessagePolicy = revertTerminationMsgPolicy(container.TerminationMsgPolicy)
	kubeContainer.ImagePullPolicy = revertImagePullPolicy(container.Pull)
	kubeContainer.VolumeMounts = revertVolumeMounts(container.VolumeMounts)

	kubeContainer.Stdin = container.Stdin
	kubeContainer.StdinOnce = container.StdinOnce
	kubeContainer.TTY = container.TTY

	// TODO: LifeCycle
	// TODO: SecurityContext

	return kubeContainer, nil
}

func revertVolumeMounts(mounts []types.VolumeMount) []v1.VolumeMount {
	var kubeMounts []v1.VolumeMount
	for i := range mounts {
		mount := mounts[i]
		kubeMount := v1.VolumeMount{}
		kubeMount.MountPropagation = revertMountPropagation(mount.Propagation)
		kubeMount.MountPath = mount.MountPath

		fields := strings.Split(mount.Store, ":")
		if len(fields) == 1 {
			kubeMount.Name = mount.Store
		} else if len(fields) == 2 {
			kubeMount.Name = fields[0]
			if fields[1] == "ro" {
				kubeMount.ReadOnly = true
			} else {
				kubeMount.SubPath = fields[2]
			}
		} else if len(fields) == 3 {
			kubeMount.Name = fields[0]
			kubeMount.SubPath = fields[1]
			kubeMount.ReadOnly = true
		}
		kubeMounts = append(kubeMounts, kubeMount)
	}
	return kubeMounts
}

func revertMountPropagation(prop types.MountPropagation) *v1.MountPropagationMode {
	var mode v1.MountPropagationMode

	if prop == types.MountPropagationHostToContainer {
		mode = v1.MountPropagationHostToContainer
	}
	if prop == types.MountPropagationBidirectional {
		mode = v1.MountPropagationBidirectional
	}
	return &mode
}

func revertImagePullPolicy(policy types.PullPolicy) v1.PullPolicy {
	if policy == types.PullAlways {
		return v1.PullAlways
	}
	if policy == types.PullNever {
		return v1.PullNever
	}
	if policy == types.PullIfNotPresent {
		return v1.PullIfNotPresent
	}
	return ""
}

func revertTerminationMsgPolicy(policy types.TerminationMessagePolicy) v1.TerminationMessagePolicy {
	if policy == types.TerminationMessageReadFile {
		return v1.TerminationMessageReadFile
	}
	if policy == types.TerminationMessageFallbackToLogsOnError {
		return v1.TerminationMessageFallbackToLogsOnError
	}
	return ""
}

func revertProbe(probe *types.Probe) (*v1.Probe, error) {
	if probe == nil {
		return nil, nil
	}
	kubeProbe := &v1.Probe{
		InitialDelaySeconds: probe.Delay,
		TimeoutSeconds:      probe.Timeout,
		PeriodSeconds:       probe.Interval,
		SuccessThreshold:    probe.MinCountSuccess,
		FailureThreshold:    probe.MinCountFailure,
	}

	if len(probe.Command) != 0 {
		kubeProbe.Exec = &v1.ExecAction{
			Command: probe.Command,
		}
	}

	if probe.Net != nil {
		urlStruct, err := url.Parse(probe.Net.URL)
		if err != nil {
			return nil, err
		}
		if urlStruct.Scheme == "TCP" {
			hostPort := urlStruct.Host
			fields := strings.Split(hostPort, ":")
			if len(fields) != 2 && len(fields) != 1 {
				return nil, util.TypeValueErrorf(urlStruct, "Unexpected value %s", hostPort)
			}
			host := fields[0]
			port := "80"
			if len(fields) == 2 {
				port = fields[1]
			}
			kubeProbe.TCPSocket = &v1.TCPSocketAction{
				Host: host,
				Port: intstr.IntOrString{
					StrVal: port,
				},
			}
		} else if urlStruct.Scheme == "HTTP" || urlStruct.Scheme == "HTTPS" {

			hostPort := urlStruct.Host
			fields := strings.Split(hostPort, ":")
			if len(fields) != 2 && len(fields) != 1 {
				return nil, util.TypeValueErrorf(urlStruct, "Unexpected value %s", hostPort)
			}
			host := fields[0]
			port := "80"
			if len(fields) == 2 {
				port = fields[1]
			}

			var scheme v1.URIScheme

			if strings.ToLower(urlStruct.Scheme) == "http" {
				scheme = v1.URISchemeHTTP
			} else if strings.ToLower(urlStruct.Scheme) == "https" {
				scheme = v1.URISchemeHTTPS
			} else {
				return nil, util.TypeValueErrorf(urlStruct, "Unexpected scheme %s", urlStruct.Scheme)
			}

			kubeProbe.HTTPGet = &v1.HTTPGetAction{
				Scheme: scheme,
				Path:   urlStruct.Path,
				Port: intstr.IntOrString{
					StrVal: port,
				},
				Host: host,
			}

			var headers []v1.HTTPHeader
			for i := range probe.Net.Headers {
				h := probe.Net.Headers[i]
				fields := strings.Split(h, ":")
				if len(fields) != 2 {
					return nil, util.TypeValueErrorf(h, "Unexpected value %s", h)
				}
				header := v1.HTTPHeader{
					Name:  fields[0],
					Value: fields[1],
				}
				headers = append(headers, header)
			}
			kubeProbe.HTTPGet.HTTPHeaders = headers
		} else {
			return nil, util.TypeValueErrorf(urlStruct, "Unexpected value %s", probe.Net.URL)
		}
	}
	return kubeProbe, nil
}

func revertResources(cpu *types.CPU, mem *types.Mem) (v1.ResourceRequirements, error) {
	limits := v1.ResourceList{}
	requests := v1.ResourceList{}
	requirements := v1.ResourceRequirements{
		Limits:   limits,
		Requests: requests,
	}

	if cpu != nil {
		if cpu.Min != "" {
			q, err := resource.ParseQuantity(cpu.Min)
			if err != nil {
				return requirements, err
			}
			requests[v1.ResourceCPU] = q
		}

		if cpu.Max != "" {
			q, err := resource.ParseQuantity(cpu.Max)
			if err != nil {
				return requirements, err
			}
			limits[v1.ResourceCPU] = q
		}
	}

	if mem != nil {
		if mem.Min != "" {
			q, err := resource.ParseQuantity(mem.Min)
			if err != nil {
				return requirements, err
			}
			requests[v1.ResourceMemory] = q
		}

		if mem.Max != "" {
			q, err := resource.ParseQuantity(mem.Max)
			if err != nil {
				return requirements, err
			}
			limits[v1.ResourceMemory] = q
		}
	}

	return requirements, nil
}

func revertEnv(envs []types.Env) ([]v1.EnvVar, []v1.EnvFromSource, error) {
	var envVars []v1.EnvVar
	var envsFromSource []v1.EnvFromSource

	for i := range envs {
		e := envs[i]
		if e.From == "" {
			fields := strings.Split(string(e.EnvStr), "=")
			if len(fields) != 2 {
				return nil, nil, util.TypeValueErrorf(e, "Unexpected value %s", string(e.EnvStr))
			}
			envVar := v1.EnvVar{
				Name:  fields[0],
				Value: fields[1],
			}
			envVars = append(envVars, envVar)
			continue
		}

		// ResourceFieldRef
		if strings.Index(e.From, "limits.") == 0 || strings.Index(e.From, "requests.") == 0 {
			envVar := v1.EnvVar{
				Name: string(e.EnvStr),
				ValueFrom: &v1.EnvVarSource{
					ResourceFieldRef: &v1.ResourceFieldSelector{
						Resource: e.From,
					},
				},
			}
			envVars = append(envVars, envVar)
			continue
		}

		// ConfigMapKeyRef or ConfigMapEnvSource
		if strings.Index(e.From, "config:") == 0 {
			fields := strings.Split(e.From, ":")
			if len(fields) == 3 {
				//ConfigMapKeyRef
				envVar := v1.EnvVar{
					Name: string(e.EnvStr),
					ValueFrom: &v1.EnvVarSource{
						ConfigMapKeyRef: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: fields[1],
							},
							Key:      fields[2],
							Optional: e.Required,
						},
					},
				}
				envVars = append(envVars, envVar)
			} else if len(fields) == 2 {
				//ConfigMapEnvSource
				envVarFromSrc := v1.EnvFromSource{
					Prefix: string(e.EnvStr),
					ConfigMapRef: &v1.ConfigMapEnvSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: fields[1],
						},
						Optional: e.Required,
					},
				}
				envsFromSource = append(envsFromSource, envVarFromSrc)
			} else {
				return nil, nil, util.TypeValueErrorf(e, "Unexpected value %s", e.From)
			}
			continue
		}

		// SecretKeyRef or SecretEnvSource
		if strings.Index(e.From, "secret:") == 0 {
			fields := strings.Split(e.From, ":")
			if len(fields) == 3 {
				//SecretKeyRef
				envVar := v1.EnvVar{
					Name: string(e.EnvStr),
					ValueFrom: &v1.EnvVarSource{
						SecretKeyRef: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: fields[1],
							},
							Key:      fields[2],
							Optional: e.Required,
						},
					},
				}
				envVars = append(envVars, envVar)
			} else if len(fields) == 2 {
				envVarFromSrc := v1.EnvFromSource{
					Prefix: string(e.EnvStr),
					SecretRef: &v1.SecretEnvSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: fields[1],
						},
						Optional: e.Required,
					},
				}
				envsFromSource = append(envsFromSource, envVarFromSrc)
			} else {
				return nil, nil, util.TypeValueErrorf(e, "Unexpected value %s", e.From)
			}
			continue
		}

		// FieldRef
		envVar := v1.EnvVar{
			Name: string(e.EnvStr),
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: e.From,
				},
			},
		}
		envVars = append(envVars, envVar)
	}

	return envVars, envsFromSource, nil
}

func revertExpose(ports []types.Port) ([]v1.ContainerPort, error) {
	var kubeContainerPorts []v1.ContainerPort
	for i := range ports {
		port := ports[i]
		kubePort := v1.ContainerPort{}

		kubePort.Name = port.Name
		protocol := v1.ProtocolTCP
		if port.Protocol == "UDP" {
			protocol = v1.ProtocolUDP
		}
		kubePort.Protocol = protocol
		fields := strings.Split(port.PortMap, ":")
		if len(fields) == 1 {
			// Then the value is container port
			containerPort, err := strconv.ParseInt(port.PortMap, 10, 32)
			if err != nil {
				return nil, err
			}
			kubePort.ContainerPort = int32(containerPort)
		} else if len(fields) == 2 {
			// Then the value is hostPort:containerport
			hostPort, err := strconv.ParseInt(fields[0], 10, 32)
			if err != nil {
				return nil, err
			}
			containerPort, err := strconv.ParseInt(fields[1], 10, 32)
			if err != nil {
				return nil, err
			}
			kubePort.ContainerPort = int32(containerPort)
			kubePort.HostPort = int32(hostPort)
		}

		kubeContainerPorts = append(kubeContainerPorts, kubePort)
	}
	return kubeContainerPorts, nil
}
