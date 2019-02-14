package converters

import (
	"net/url"
	"sort"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	intstr "k8s.io/apimachinery/pkg/util/intstr"

	"github.com/koki/short/converter/converters/affinity"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
	"github.com/koki/short/util/floatstr"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_Pod_to_Kube_v1_Pod(pod *types.PodWrapper) (*v1.Pod, error) {
	var err error
	kubePod := &v1.Pod{}
	kokiPod := pod.Pod

	if len(kokiPod.Version) == 0 {
		kubePod.APIVersion = "v1"
	} else {
		kubePod.APIVersion = kokiPod.Version
	}
	kubePod.Kind = "Pod"

	kubePod.ObjectMeta = revertPodObjectMeta(kokiPod.PodTemplateMeta)

	spec, err := revertPodSpec(kokiPod.PodTemplate)
	if err != nil {
		return nil, err
	}
	kubePod.Spec = *spec

	kubePod.Status = v1.PodStatus{}

	phase, err := revertPodPhase(kokiPod.Phase)
	if err != nil {
		return nil, err
	}
	kubePod.Status.Phase = phase
	kubePod.Status.Message = kokiPod.Msg
	kubePod.Status.Reason = kokiPod.Reason
	kubePod.Status.HostIP = kokiPod.NodeIP
	kubePod.Status.PodIP = kokiPod.IP

	qos, err := revertQOSClass(kokiPod.QOS)
	if err != nil {
		return nil, err
	}
	kubePod.Status.QOSClass = qos
	kubePod.Status.StartTime = kokiPod.StartTime

	conditions, err := revertPodConditions(kokiPod.Conditions)
	if err != nil {
		return nil, err
	}
	kubePod.Status.Conditions = conditions

	var initContainerStatuses []v1.ContainerStatus
	for i := range kokiPod.InitContainers {
		container := kokiPod.InitContainers[i]
		status, err := revertContainerStatus(container)
		if err != nil {
			return nil, err
		}
		initContainerStatuses = append(initContainerStatuses, status)
	}
	kubePod.Status.InitContainerStatuses = initContainerStatuses

	var containerStatuses []v1.ContainerStatus
	for i := range kokiPod.Containers {
		container := kokiPod.Containers[i]
		status, err := revertContainerStatus(container)
		if err != nil {
			return nil, err
		}
		if status.ContainerID != "" {
			containerStatuses = append(containerStatuses, status)
		}
	}
	kubePod.Status.ContainerStatuses = containerStatuses

	return kubePod, nil
}

func revertPodObjectMeta(kokiMeta types.PodTemplateMeta) metav1.ObjectMeta {
	var labels map[string]string
	if len(kokiMeta.Labels) > 0 {
		labels = kokiMeta.Labels
	}
	var annotations map[string]string
	if len(kokiMeta.Annotations) > 0 {
		annotations = kokiMeta.Annotations
	}
	return metav1.ObjectMeta{
		Name:        kokiMeta.Name,
		Namespace:   kokiMeta.Namespace,
		ClusterName: kokiMeta.Cluster,
		Labels:      labels,
		Annotations: annotations,
	}
}

func revertPodSpec(kokiPod types.PodTemplate) (*v1.PodSpec, error) {
	var err error
	spec := v1.PodSpec{}
	spec.Volumes, err = revertVolumes(kokiPod.Volumes)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "pod volumes")
	}
	fields := strings.SplitN(kokiPod.Hostname, ".", 2)
	if len(fields) == 1 {
		spec.Hostname = kokiPod.Hostname
	} else {
		spec.Hostname = fields[1]
		spec.Subdomain = fields[0]
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
	spec.InitContainers = initContainers

	var kubeContainers []v1.Container
	for i := range kokiPod.Containers {
		container := kokiPod.Containers[i]
		kubeContainer, err := revertKokiContainer(container)
		if err != nil {
			return nil, err
		}
		kubeContainers = append(kubeContainers, kubeContainer)
	}
	spec.Containers = kubeContainers

	hostAliases, err := revertHostAliases(kokiPod.HostAliases)
	if err != nil {
		return nil, err
	}
	spec.HostAliases = hostAliases

	restartPolicy, err := revertRestartPolicy(kokiPod.RestartPolicy)
	if err != nil {
		return nil, err
	}
	spec.RestartPolicy = restartPolicy

	affinity, err := affinity.Convert_Koki_Affinity_to_Kube_v1_Affinity(kokiPod.Affinity)
	if err != nil {
		return nil, err
	}
	spec.Affinity = affinity

	spec.TerminationGracePeriodSeconds = kokiPod.TerminationGracePeriod
	spec.ActiveDeadlineSeconds = kokiPod.ActiveDeadline

	dnsPolicy, err := revertDNSPolicy(kokiPod.DNSPolicy)
	if err != nil {
		return nil, err
	}
	spec.DNSPolicy = dnsPolicy

	serviceAccount, autoMount, err := revertServiceAccount(kokiPod.Account)
	if err != nil {
		return nil, err
	}
	spec.ServiceAccountName = serviceAccount
	spec.AutomountServiceAccountToken = autoMount
	spec.NodeName = kokiPod.Node

	net, pid, ipc, err := revertHostModes(kokiPod.HostMode)
	if err != nil {
		return nil, err
	}
	spec.HostNetwork = net
	spec.HostPID = pid
	spec.HostIPC = ipc
	spec.ImagePullSecrets = revertRegistries(kokiPod.Registries)
	spec.SchedulerName = kokiPod.SchedulerName

	tolerations, err := revertTolerations(kokiPod.Tolerations)
	if err != nil {
		return nil, err
	}
	spec.Tolerations = tolerations

	if kokiPod.FSGID != nil || kokiPod.GIDs != nil {
		spec.SecurityContext = &v1.PodSecurityContext{}
		spec.SecurityContext.FSGroup = kokiPod.FSGID
		spec.SecurityContext.SupplementalGroups = kokiPod.GIDs
	}

	if kokiPod.Priority != nil {
		spec.Priority = kokiPod.Priority.Value
		spec.PriorityClassName = kokiPod.Priority.Class
	}

	return &spec, nil
}

func revertVolumes(kokiVolumes map[string]types.Volume) ([]v1.Volume, error) {
	kubeVolumes := []v1.Volume{}
	names := []string{}
	for name, _ := range kokiVolumes {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		kokiVolume := kokiVolumes[name]
		kubeVolume, err := revertVolume(name, kokiVolume)
		if err != nil {
			return nil, err
		}
		kubeVolumes = append(kubeVolumes, *kubeVolume)
	}

	return kubeVolumes, nil
}

func revertStorageMedium(kokiMedium types.StorageMedium) (v1.StorageMedium, error) {
	switch kokiMedium {
	case types.StorageMediumDefault:
		return v1.StorageMediumDefault, nil
	case types.StorageMediumMemory:
		return v1.StorageMediumMemory, nil
	case types.StorageMediumHugePages:
		return v1.StorageMediumHugePages, nil
	default:
		return v1.StorageMediumDefault, serrors.InvalidValueErrorf(kokiMedium, "unrecognized storage medium")
	}
}

func revertHostPathType(kokiType types.HostPathType) (v1.HostPathType, error) {
	switch kokiType {
	case types.HostPathUnset:
		return v1.HostPathUnset, nil
	case types.HostPathDirectoryOrCreate:
		return v1.HostPathDirectoryOrCreate, nil
	case types.HostPathDirectory:
		return v1.HostPathDirectory, nil
	case types.HostPathFileOrCreate:
		return v1.HostPathFileOrCreate, nil
	case types.HostPathFile:
		return v1.HostPathFile, nil
	case types.HostPathSocket:
		return v1.HostPathSocket, nil
	case types.HostPathCharDev:
		return v1.HostPathCharDev, nil
	case types.HostPathBlockDev:
		return v1.HostPathBlockDev, nil
	default:
		return v1.HostPathUnset, serrors.InvalidValueErrorf(kokiType, "unrecognized host_path type")
	}
}

func revertAzureDiskKind(kokiKind *types.AzureDataDiskKind) (*v1.AzureDataDiskKind, error) {
	if kokiKind == nil {
		return nil, nil
	}

	var kind v1.AzureDataDiskKind
	switch *kokiKind {
	case types.AzureDedicatedBlobDisk:
		kind = v1.AzureDedicatedBlobDisk
	case types.AzureSharedBlobDisk:
		kind = v1.AzureSharedBlobDisk
	case types.AzureManagedDisk:
		kind = v1.AzureManagedDisk
	default:
		return nil, serrors.InvalidValueErrorf(kokiKind, "unrecognized kind")
	}

	return &kind, nil
}

func revertAzureDiskCachingMode(kokiMode *types.AzureDataDiskCachingMode) (*v1.AzureDataDiskCachingMode, error) {
	if kokiMode == nil {
		return nil, nil
	}

	var mode v1.AzureDataDiskCachingMode
	switch *kokiMode {
	case types.AzureDataDiskCachingNone:
		mode = v1.AzureDataDiskCachingNone
	case types.AzureDataDiskCachingReadOnly:
		mode = v1.AzureDataDiskCachingReadOnly
	case types.AzureDataDiskCachingReadWrite:
		mode = v1.AzureDataDiskCachingReadWrite
	default:
		return nil, serrors.InvalidValueErrorf(kokiMode, "unrecognized cache")
	}

	return &mode, nil
}

func revertCephFSSecretFileOrRef(kokiSecret *types.CephFSSecretFileOrRef) (string, *v1.LocalObjectReference) {
	if kokiSecret == nil {
		return "", nil
	}

	if len(kokiSecret.File) > 0 {
		return kokiSecret.File, nil
	}

	return "", &v1.LocalObjectReference{
		Name: kokiSecret.Ref,
	}
}

func revertLocalObjectRef(kokiRef string) *v1.LocalObjectReference {
	if len(kokiRef) == 0 {
		return nil
	}
	return &v1.LocalObjectReference{
		Name: kokiRef,
	}
}

func revertVsphereStoragePolicy(kokiPolicy *types.VsphereStoragePolicy) (name, ID string) {
	if kokiPolicy == nil {
		return "", ""
	}

	return kokiPolicy.Name, kokiPolicy.ID
}

func revertFileMode(kokiMode *types.FileMode) *int32 {
	if kokiMode == nil {
		return nil
	}

	return util.Int32Ptr(int32(*kokiMode))
}

func revertKeyToPathItems(kokiItems map[string]types.KeyAndMode) []v1.KeyToPath {
	if len(kokiItems) == 0 {
		return nil
	}

	kubeItems := []v1.KeyToPath{}
	for path, item := range kokiItems {
		kubeItems = append(kubeItems, v1.KeyToPath{
			Path: path,
			Key:  item.Key,
			Mode: revertFileMode(item.Mode),
		})
	}

	return kubeItems
}

func revertRequiredToOptional(required *bool) *bool {
	if required == nil {
		return nil
	}

	return util.BoolPtr(!*required)
}

func revertDownwardAPIVolumeFiles(kokiItems map[string]types.DownwardAPIVolumeFile) []v1.DownwardAPIVolumeFile {
	if len(kokiItems) == 0 {
		return nil
	}

	items := []v1.DownwardAPIVolumeFile{}
	paths := []string{}
	for path, _ := range kokiItems {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		kokiItem := kokiItems[path]
		items = append(items, v1.DownwardAPIVolumeFile{
			Path:             path,
			FieldRef:         revertObjectFieldRef(kokiItem.FieldRef),
			ResourceFieldRef: revertVolumeResourceFieldRef(kokiItem.ResourceFieldRef),
			Mode:             revertFileMode(kokiItem.Mode),
		})
	}

	return items
}

func revertObjectFieldRef(kokiRef *types.ObjectFieldSelector) *v1.ObjectFieldSelector {
	if kokiRef == nil {
		return nil
	}

	return &v1.ObjectFieldSelector{
		FieldPath:  kokiRef.FieldPath,
		APIVersion: kokiRef.APIVersion,
	}
}

func revertVolumeResourceFieldRef(kokiRef *types.VolumeResourceFieldSelector) *v1.ResourceFieldSelector {
	if kokiRef == nil {
		return nil
	}

	return &v1.ResourceFieldSelector{
		ContainerName: kokiRef.ContainerName,
		Resource:      kokiRef.Resource,
		Divisor:       kokiRef.Divisor,
	}
}

func revertVolumeProjections(kokiProjections []types.VolumeProjection) ([]v1.VolumeProjection, error) {
	if len(kokiProjections) == 0 {
		return nil, nil
	}

	projections := make([]v1.VolumeProjection, len(kokiProjections))
	for i, projection := range kokiProjections {
		secret, err := revertSecretProjection(projection.Secret)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "secret volume-projection #%d", i)
		}
		config, err := revertConfigMapProjection(projection.ConfigMap)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "config-map volume-projection #%d", i)
		}
		projections[i] = v1.VolumeProjection{
			Secret:      secret,
			DownwardAPI: revertDownwardAPIProjection(projection.DownwardAPI),
			ConfigMap:   config,
		}
	}

	return projections, nil
}

func revertSecretProjection(kokiProjection *types.SecretProjection) (*v1.SecretProjection, error) {
	if kokiProjection == nil {
		return nil, nil
	}

	ref := revertLocalObjectRef(kokiProjection.Name)
	if ref == nil {
		return nil, serrors.InvalidInstanceErrorf(kokiProjection, "secret name is missing")
	}
	return &v1.SecretProjection{
		LocalObjectReference: *ref,
		Items:                revertKeyToPathItems(kokiProjection.Items),
	}, nil
}

func revertConfigMapProjection(kokiProjection *types.ConfigMapProjection) (*v1.ConfigMapProjection, error) {
	if kokiProjection == nil {
		return nil, nil
	}

	ref := revertLocalObjectRef(kokiProjection.Name)
	if ref == nil {
		return nil, serrors.InvalidInstanceErrorf(kokiProjection, "config-map name is missing")
	}
	return &v1.ConfigMapProjection{
		LocalObjectReference: *ref,
		Items:                revertKeyToPathItems(kokiProjection.Items),
	}, nil
}

func revertDownwardAPIProjection(kokiProjection *types.DownwardAPIProjection) *v1.DownwardAPIProjection {
	if kokiProjection == nil {
		return nil
	}

	return &v1.DownwardAPIProjection{
		Items: revertDownwardAPIVolumeFiles(kokiProjection.Items),
	}
}

func revertGcePDVolume(source *types.GcePDVolume) *v1.GCEPersistentDiskVolumeSource {
	return &v1.GCEPersistentDiskVolumeSource{
		PDName:    source.PDName,
		FSType:    source.FSType,
		Partition: source.Partition,
		ReadOnly:  source.ReadOnly,
	}
}

func revertAwsEBSVolume(source *types.AwsEBSVolume) *v1.AWSElasticBlockStoreVolumeSource {
	return &v1.AWSElasticBlockStoreVolumeSource{
		VolumeID:  source.VolumeID,
		FSType:    source.FSType,
		Partition: source.Partition,
		ReadOnly:  source.ReadOnly,
	}
}

func revertHostPathVolume(source *types.HostPathVolume) (*v1.HostPathVolumeSource, error) {
	kubeType, err := revertHostPathType(source.Type)
	if err != nil {
		return nil, err
	}
	return &v1.HostPathVolumeSource{
		Path: source.Path,
		Type: &kubeType,
	}, nil
}

func revertGlusterfsVolume(source *types.GlusterfsVolume) *v1.GlusterfsVolumeSource {
	return &v1.GlusterfsVolumeSource{
		EndpointsName: source.EndpointsName,
		Path:          source.Path,
		ReadOnly:      source.ReadOnly,
	}
}

func revertNFSVolume(source *types.NFSVolume) *v1.NFSVolumeSource {
	return &v1.NFSVolumeSource{
		Server:   source.Server,
		Path:     source.Path,
		ReadOnly: source.ReadOnly,
	}
}

func revertISCSIVolume(source *types.ISCSIVolume) *v1.ISCSIVolumeSource {
	return &v1.ISCSIVolumeSource{
		TargetPortal:      source.TargetPortal,
		IQN:               source.IQN,
		Lun:               source.Lun,
		ISCSIInterface:    source.ISCSIInterface,
		FSType:            source.FSType,
		ReadOnly:          source.ReadOnly,
		Portals:           source.Portals,
		DiscoveryCHAPAuth: source.DiscoveryCHAPAuth,
		SessionCHAPAuth:   source.SessionCHAPAuth,
		SecretRef:         revertLocalObjectRef(source.SecretRef),
		InitiatorName:     util.StringPtrOrNil(source.InitiatorName),
	}
}

func revertCinderVolume(source *types.CinderVolume) *v1.CinderVolumeSource {
	return &v1.CinderVolumeSource{
		VolumeID: source.VolumeID,
		FSType:   source.FSType,
		ReadOnly: source.ReadOnly,
	}
}

func revertFibreChannelVolume(source *types.FibreChannelVolume) *v1.FCVolumeSource {
	return &v1.FCVolumeSource{
		TargetWWNs: source.TargetWWNs,
		Lun:        source.Lun,
		FSType:     source.FSType,
		ReadOnly:   source.ReadOnly,
		WWIDs:      source.WWIDs,
	}
}

func revertFlockerVolume(source *types.FlockerVolume) *v1.FlockerVolumeSource {
	return &v1.FlockerVolumeSource{
		DatasetUUID: source.DatasetUUID,
	}
}

func revertFlexVolume(source *types.FlexVolume) *v1.FlexVolumeSource {
	return &v1.FlexVolumeSource{
		Driver:    source.Driver,
		FSType:    source.FSType,
		SecretRef: revertLocalObjectRef(source.SecretRef),
		ReadOnly:  source.ReadOnly,
		Options:   source.Options,
	}
}

func revertVsphereVolume(source *types.VsphereVolume) *v1.VsphereVirtualDiskVolumeSource {
	storagePolicyName, storagePolicyID := revertVsphereStoragePolicy(source.StoragePolicy)
	return &v1.VsphereVirtualDiskVolumeSource{
		VolumePath:        source.VolumePath,
		FSType:            source.FSType,
		StoragePolicyName: storagePolicyName,
		StoragePolicyID:   storagePolicyID,
	}
}

func revertQuobyteVolume(source *types.QuobyteVolume) *v1.QuobyteVolumeSource {
	return &v1.QuobyteVolumeSource{
		Registry: source.Registry,
		Volume:   source.Volume,
		ReadOnly: source.ReadOnly,
		User:     source.User,
		Group:    source.Group,
	}
}

func revertScaleIOStorageMode(mode types.ScaleIOStorageMode) (string, error) {
	if len(mode) == 0 {
		return "", nil
	}

	switch mode {
	case types.ScaleIOStorageModeThick:
		return "ThickProvisioned", nil
	case types.ScaleIOStorageModeThin:
		return "ThinProvisioned", nil
	default:
		return "", serrors.InvalidInstanceError(mode)
	}
}

func revertAzureDiskVolume(source *types.AzureDiskVolume) (*v1.AzureDiskVolumeSource, error) {
	kind, err := revertAzureDiskKind(source.Kind)
	if err != nil {
		return nil, err
	}
	cachingMode, err := revertAzureDiskCachingMode(source.CachingMode)
	if err != nil {
		return nil, err
	}
	return &v1.AzureDiskVolumeSource{
		DiskName:    source.DiskName,
		DataDiskURI: source.DataDiskURI,
		FSType:      util.StringPtrOrNil(source.FSType),
		ReadOnly:    util.BoolPtrOrNil(source.ReadOnly),
		Kind:        kind,
		CachingMode: cachingMode,
	}, nil
}

func revertPhotonPDVolume(source *types.PhotonPDVolume) *v1.PhotonPersistentDiskVolumeSource {
	return &v1.PhotonPersistentDiskVolumeSource{
		PdID:   source.PdID,
		FSType: source.FSType,
	}
}

func revertPortworxVolume(source *types.PortworxVolume) *v1.PortworxVolumeSource {
	return &v1.PortworxVolumeSource{
		VolumeID: source.VolumeID,
		FSType:   source.FSType,
		ReadOnly: source.ReadOnly,
	}
}

func revertVolume(name string, kokiVolume types.Volume) (*v1.Volume, error) {
	if kokiVolume.EmptyDir != nil {
		medium, err := revertStorageMedium(kokiVolume.EmptyDir.Medium)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "volume (%s)", name)
		}
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{
					Medium:    medium,
					SizeLimit: kokiVolume.EmptyDir.SizeLimit,
				},
			},
		}, nil
	}
	if kokiVolume.HostPath != nil {
		source, err := revertHostPathVolume(kokiVolume.HostPath)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "volume (%s)", name)
		}
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				HostPath: source,
			},
		}, nil
	}
	if kokiVolume.GcePD != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				GCEPersistentDisk: revertGcePDVolume(kokiVolume.GcePD),
			},
		}, nil
	}
	if kokiVolume.AwsEBS != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				AWSElasticBlockStore: revertAwsEBSVolume(kokiVolume.AwsEBS),
			},
		}, nil
	}
	if kokiVolume.AzureDisk != nil {
		source, err := revertAzureDiskVolume(kokiVolume.AzureDisk)
		if err != nil {
			return nil, err
		}
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				AzureDisk: source,
			},
		}, nil
	}
	if kokiVolume.AzureFile != nil {
		source := kokiVolume.AzureFile
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				AzureFile: &v1.AzureFileVolumeSource{
					SecretName: source.SecretName,
					ShareName:  source.ShareName,
					ReadOnly:   source.ReadOnly,
				},
			},
		}, nil
	}
	if kokiVolume.CephFS != nil {
		source := kokiVolume.CephFS
		secretFile, secretRef := revertCephFSSecretFileOrRef(source.SecretFileOrRef)
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				CephFS: &v1.CephFSVolumeSource{
					Monitors:   source.Monitors,
					Path:       source.Path,
					User:       source.User,
					SecretFile: secretFile,
					SecretRef:  secretRef,
					ReadOnly:   source.ReadOnly,
				},
			},
		}, nil
	}
	if kokiVolume.Cinder != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				Cinder: revertCinderVolume(kokiVolume.Cinder),
			},
		}, nil
	}
	if kokiVolume.FibreChannel != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				FC: revertFibreChannelVolume(kokiVolume.FibreChannel),
			},
		}, nil
	}
	if kokiVolume.Flex != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				FlexVolume: revertFlexVolume(kokiVolume.Flex),
			},
		}, nil
	}
	if kokiVolume.Flocker != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				Flocker: revertFlockerVolume(kokiVolume.Flocker),
			},
		}, nil
	}
	if kokiVolume.Glusterfs != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				Glusterfs: revertGlusterfsVolume(kokiVolume.Glusterfs),
			},
		}, nil
	}
	if kokiVolume.ISCSI != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				ISCSI: revertISCSIVolume(kokiVolume.ISCSI),
			},
		}, nil
	}
	if kokiVolume.NFS != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				NFS: revertNFSVolume(kokiVolume.NFS),
			},
		}, nil
	}
	if kokiVolume.PhotonPD != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				PhotonPersistentDisk: revertPhotonPDVolume(kokiVolume.PhotonPD),
			},
		}, nil
	}
	if kokiVolume.Portworx != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				PortworxVolume: revertPortworxVolume(kokiVolume.Portworx),
			},
		}, nil

	}
	if kokiVolume.PVC != nil {
		source := kokiVolume.PVC
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: source.ClaimName,
					ReadOnly:  source.ReadOnly,
				},
			},
		}, nil
	}
	if kokiVolume.Quobyte != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				Quobyte: revertQuobyteVolume(kokiVolume.Quobyte),
			},
		}, nil
	}
	if kokiVolume.ScaleIO != nil {
		source := kokiVolume.ScaleIO
		mode, err := revertScaleIOStorageMode(source.StorageMode)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "ScaleIO storage mode")
		}
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				ScaleIO: &v1.ScaleIOVolumeSource{
					Gateway:          source.Gateway,
					System:           source.System,
					SecretRef:        revertLocalObjectRef(source.SecretRef),
					SSLEnabled:       source.SSLEnabled,
					ProtectionDomain: source.ProtectionDomain,
					StoragePool:      source.StoragePool,
					StorageMode:      mode,
					VolumeName:       source.VolumeName,
					FSType:           source.FSType,
					ReadOnly:         source.ReadOnly,
				},
			},
		}, nil
	}
	if kokiVolume.Vsphere != nil {
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				VsphereVolume: revertVsphereVolume(kokiVolume.Vsphere),
			},
		}, nil
	}
	if kokiVolume.ConfigMap != nil {
		source := kokiVolume.ConfigMap
		ref := revertLocalObjectRef(source.Name)
		if ref == nil {
			return nil, serrors.InvalidInstanceErrorf(source, "config name is required")
		}

		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				ConfigMap: &v1.ConfigMapVolumeSource{
					LocalObjectReference: *ref,
					Items:                revertKeyToPathItems(source.Items),
					DefaultMode:          revertFileMode(source.DefaultMode),
					Optional:             revertRequiredToOptional(source.Required),
				},
			},
		}, nil
	}
	if kokiVolume.Secret != nil {
		source := kokiVolume.Secret
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName:  source.SecretName,
					Items:       revertKeyToPathItems(source.Items),
					DefaultMode: revertFileMode(source.DefaultMode),
					Optional:    revertRequiredToOptional(source.Required),
				},
			},
		}, nil
	}
	if kokiVolume.DownwardAPI != nil {
		source := kokiVolume.DownwardAPI
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				DownwardAPI: &v1.DownwardAPIVolumeSource{
					Items:       revertDownwardAPIVolumeFiles(source.Items),
					DefaultMode: revertFileMode(source.DefaultMode),
				},
			},
		}, nil
	}
	if kokiVolume.Projected != nil {
		source := kokiVolume.Projected
		sources, err := revertVolumeProjections(source.Sources)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "volume (%s)", name)
		}
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				Projected: &v1.ProjectedVolumeSource{
					Sources:     sources,
					DefaultMode: revertFileMode(source.DefaultMode),
				},
			},
		}, nil
	}
	if kokiVolume.Git != nil {
		source := kokiVolume.Git
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				GitRepo: &v1.GitRepoVolumeSource{
					Repository: source.Repository,
					Revision:   source.Revision,
					Directory:  source.Directory,
				},
			},
		}, nil
	}
	if kokiVolume.RBD != nil {
		source := kokiVolume.RBD
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				RBD: &v1.RBDVolumeSource{
					CephMonitors: source.CephMonitors,
					RBDImage:     source.RBDImage,
					FSType:       source.FSType,
					RBDPool:      source.RBDPool,
					RadosUser:    source.RadosUser,
					Keyring:      source.Keyring,
					SecretRef:    revertLocalObjectRef(source.SecretRef),
					ReadOnly:     source.ReadOnly,
				},
			},
		}, nil
	}
	if kokiVolume.StorageOS != nil {
		source := kokiVolume.StorageOS
		return &v1.Volume{
			Name: name,
			VolumeSource: v1.VolumeSource{
				StorageOS: &v1.StorageOSVolumeSource{
					VolumeName:      source.VolumeName,
					VolumeNamespace: source.VolumeNamespace,
					FSType:          source.FSType,
					ReadOnly:        source.ReadOnly,
					SecretRef:       revertLocalObjectRef(source.SecretRef),
				},
			},
		}, nil
	}

	return nil, serrors.InvalidInstanceErrorf(kokiVolume, "empty volume definition")
}

func revertContainerStatus(container types.Container) (v1.ContainerStatus, error) {
	var status v1.ContainerStatus

	status.ContainerID = container.ContainerID
	status.ImageID = container.ImageID
	status.RestartCount = container.Restarts
	status.Ready = container.Ready
	status.State = revertContainerState(container.CurrentState)
	status.LastTerminationState = revertContainerState(container.LastState)

	return status, nil
}

func revertContainerState(state *types.ContainerState) v1.ContainerState {
	containerState := v1.ContainerState{}
	if state == nil {
		return containerState
	}

	if state.Waiting != nil {
		containerState.Waiting = &v1.ContainerStateWaiting{
			Reason:  state.Waiting.Reason,
			Message: state.Waiting.Msg,
		}
	}
	if state.Running != nil {
		containerState.Running = &v1.ContainerStateRunning{
			StartedAt: state.Running.StartTime,
		}
	}
	if state.Terminated != nil {
		containerState.Terminated = &v1.ContainerStateTerminated{
			StartedAt:  state.Terminated.StartTime,
			FinishedAt: state.Terminated.FinishTime,
			Reason:     state.Terminated.Reason,
			Message:    state.Terminated.Msg,
			Signal:     state.Terminated.Signal,
			ExitCode:   state.Terminated.ExitCode,
		}
	}
	return containerState
}

func revertQOSClass(class types.PodQOSClass) (v1.PodQOSClass, error) {
	if class == "" {
		return "", nil
	}
	if class == types.PodQOSGuaranteed {
		return v1.PodQOSGuaranteed, nil
	}
	if class == types.PodQOSBurstable {
		return v1.PodQOSBurstable, nil
	}
	if class == types.PodQOSBestEffort {
		return v1.PodQOSBestEffort, nil
	}
	return "", serrors.InvalidInstanceError(class)
}

func revertPodPhase(phase types.PodPhase) (v1.PodPhase, error) {
	if phase == "" {
		return "", nil
	}
	if phase == types.PodPending {
		return v1.PodPending, nil
	}
	if phase == types.PodRunning {
		return v1.PodRunning, nil
	}
	if phase == types.PodSucceeded {
		return v1.PodSucceeded, nil
	}
	if phase == types.PodFailed {
		return v1.PodFailed, nil
	}
	if phase == types.PodUnknown {
		return v1.PodUnknown, nil
	}
	return "", serrors.InvalidInstanceError(phase)
}

func revertPodConditions(conditions []types.PodCondition) ([]v1.PodCondition, error) {
	var kubeConditions []v1.PodCondition

	for i := range conditions {
		condition := conditions[i]

		kubeCondition := v1.PodCondition{
			LastProbeTime:      condition.LastProbeTime,
			LastTransitionTime: condition.LastTransitionTime,
			Message:            condition.Msg,
			Reason:             condition.Reason,
		}

		podConditionType, err := revertPodConditionType(condition.Type)
		if err != nil {
			return nil, err
		}
		kubeCondition.Type = podConditionType

		podConditionStatus, err := revertConditionStatus(condition.Status)
		if err != nil {
			return nil, err
		}
		kubeCondition.Status = podConditionStatus

		kubeConditions = append(kubeConditions, kubeCondition)
	}

	return kubeConditions, nil
}

func revertPodConditionType(typ types.PodConditionType) (v1.PodConditionType, error) {
	if typ == "" {
		return "", nil
	}
	if typ == types.PodScheduled {
		return v1.PodScheduled, nil
	}
	if typ == types.PodReady {
		return v1.PodReady, nil
	}
	if typ == types.PodInitialized {
		return v1.PodInitialized, nil
	}
	if typ == types.PodReasonUnschedulable {
		return v1.PodReasonUnschedulable, nil
	}
	return "", serrors.InvalidInstanceError(typ)
}

func revertConditionStatus(status types.ConditionStatus) (v1.ConditionStatus, error) {
	if status == "" {
		return "", nil
	}
	if status == types.ConditionTrue {
		return v1.ConditionTrue, nil
	}
	if status == types.ConditionFalse {
		return v1.ConditionFalse, nil
	}
	if status == types.ConditionUnknown {
		return v1.ConditionUnknown, nil
	}
	return "", serrors.InvalidInstanceError(status)

}

func revertTolerations(tolerations []types.Toleration) ([]v1.Toleration, error) {
	var kubeTolerations []v1.Toleration

	for i := range tolerations {
		toleration := tolerations[i]
		kubeToleration := v1.Toleration{
			TolerationSeconds: toleration.ExpiryAfter,
		}

		superFields := strings.Split(string(toleration.Selector), ":")
		switch len(superFields) {
		case 2:
			switch superFields[1] {
			case "NoSchedule":
				kubeToleration.Effect = v1.TaintEffectNoSchedule
			case "PreferNoSchedule":
				kubeToleration.Effect = v1.TaintEffectPreferNoSchedule
			case "NoExecute":
				kubeToleration.Effect = v1.TaintEffectNoExecute
			default:
				return nil, serrors.InvalidInstanceErrorf(toleration, "unexpected toleration selector")
			}
		case 1:
			// Do nothing
		default:
			return nil, serrors.InvalidInstanceErrorf(toleration, "unexpected toleration effect")
		}

		fields := strings.Split(superFields[0], "=")
		if len(fields) == 1 {
			if fields[0] == "*" {
				kubeToleration.Key = ""
			} else {
				kubeToleration.Key = fields[0]
			}
			kubeToleration.Operator = v1.TolerationOpExists
		} else if len(fields) == 2 {
			kubeToleration.Key = fields[0]
			kubeToleration.Operator = v1.TolerationOpEqual
			kubeToleration.Value = fields[1]
		} else {
			return nil, serrors.InvalidInstanceErrorf(toleration, "unexpected toleration selector")
		}

		kubeTolerations = append(kubeTolerations, kubeToleration)
	}

	return kubeTolerations, nil
}

func revertRegistries(registries []string) []v1.LocalObjectReference {
	var kubeRegistries []v1.LocalObjectReference

	for i := range registries {
		ref := v1.LocalObjectReference{
			Name: registries[i],
		}
		kubeRegistries = append(kubeRegistries, ref)
	}

	return kubeRegistries
}

func revertHostModes(modes []types.HostMode) (net bool, pid bool, ipc bool, err error) {
	for i := range modes {
		mode := modes[i]
		switch mode {
		case types.HostModeNet:
			net = true
		case types.HostModePID:
			pid = true
		case types.HostModeIPC:
			ipc = true
		default:
			return false, false, false, serrors.InvalidInstanceError(mode)
		}
	}

	return net, pid, ipc, nil
}

func revertServiceAccount(account string) (string, *bool, error) {
	if account == "" {
		return "", nil, nil
	}

	fields := strings.Split(account, ":")
	if len(fields) == 2 {
		if fields[1] == "auto" {
			return fields[0], util.BoolPtr(true), nil
		}
		return "", nil, serrors.InvalidValueErrorf(account, "unexpected service account automount value (%s)", fields[1])
	} else if len(fields) == 1 {
		return fields[0], nil, nil
	}

	return "", nil, serrors.InvalidValueErrorf(account, "unexpected service account automount value")
}

func revertDNSPolicy(dnsPolicy types.DNSPolicy) (v1.DNSPolicy, error) {
	if dnsPolicy == "" {
		return "", nil
	}
	if dnsPolicy == types.DNSClusterFirstWithHostNet {
		return v1.DNSClusterFirstWithHostNet, nil
	}
	if dnsPolicy == types.DNSClusterFirst {
		return v1.DNSClusterFirst, nil
	}
	if dnsPolicy == types.DNSDefault {
		return v1.DNSDefault, nil
	}
	return "", serrors.InvalidInstanceError(dnsPolicy)
}

func revertRestartPolicy(policy types.RestartPolicy) (v1.RestartPolicy, error) {
	if policy == "" {
		return "", nil
	}
	if policy == types.RestartPolicyAlways {
		return v1.RestartPolicyAlways, nil
	}
	if policy == types.RestartPolicyOnFailure {
		return v1.RestartPolicyOnFailure, nil
	}
	if policy == types.RestartPolicyNever {
		return v1.RestartPolicyNever, nil
	}
	return "", serrors.InvalidInstanceError(policy)
}

func revertHostAliases(aliases []string) ([]v1.HostAlias, error) {
	var hostAliases []v1.HostAlias
	for i := range aliases {
		alias := aliases[i]
		hostAlias := v1.HostAlias{}

		fields := strings.SplitN(alias, " ", 2)
		if len(fields) == 2 {
			hostAlias.IP = strings.TrimSpace(fields[0])
			hostNames := strings.Split(strings.TrimSpace(fields[1]), " ")
			for i := range hostNames {
				hostname := hostNames[i]
				if hostname != "" && hostname != " " {
					hostAlias.Hostnames = append(hostAlias.Hostnames, hostname)
				}
			}
		} else {
			return nil, serrors.InvalidValueForTypeErrorf(alias, hostAlias, "expected 2 space-separated values")
		}
		hostAliases = append(hostAliases, hostAlias)
	}
	return hostAliases, nil
}

func revertKokiContainer(container types.Container) (v1.Container, error) {
	kubeContainer := v1.Container{}

	kubeContainer.Name = container.Name
	kubeContainer.Args = revertContainerArgs(container.Args)
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

	lc, err := revertLifecycle(container.OnStart, container.PreStop)
	if err != nil {
		return v1.Container{}, err
	}
	kubeContainer.Lifecycle = lc

	sc, err := revertSecurityContext(container)
	if err != nil {
		return v1.Container{}, err
	}
	kubeContainer.SecurityContext = sc

	return kubeContainer, nil
}

func revertContainerArgs(kokiArgs []floatstr.FloatOrString) []string {
	if kokiArgs == nil {
		return nil
	}
	kubeArgs := make([]string, len(kokiArgs))
	for i, kokiArg := range kokiArgs {
		kubeArgs[i] = kokiArg.String()
	}

	return kubeArgs
}

func revertSecurityContext(container types.Container) (*v1.SecurityContext, error) {
	sc := &v1.SecurityContext{}

	var mark bool

	if container.Privileged != nil {
		sc.Privileged = container.Privileged
		mark = true
	}

	if container.AllowEscalation != nil {
		sc.AllowPrivilegeEscalation = container.AllowEscalation
		mark = true
	}

	if container.RO != nil || container.RW != nil {
		ro := util.FromBoolPtr(container.RO)
		rw := util.FromBoolPtr(container.RW)

		if !((!ro && rw) || (!rw && ro)) {
			return nil, serrors.InvalidInstanceErrorf(container, "conflicting value (Read Only) %v and (ReadWrite) %v", ro, rw)
		}

		sc.ReadOnlyRootFilesystem = &ro
		mark = true
	}

	if container.ForceNonRoot != nil {
		sc.RunAsNonRoot = container.ForceNonRoot
		mark = true
	}

	if container.UID != nil {
		sc.RunAsUser = container.UID
		mark = true
	}

	if container.GID != nil {
		sc.RunAsGroup = container.GID
		mark = true
	}

	if container.AddCapabilities != nil || container.DelCapabilities != nil {
		caps := &v1.Capabilities{}
		var capMark bool
		for i := range container.AddCapabilities {
			capability := container.AddCapabilities[i]
			caps.Add = append(caps.Add, v1.Capability(capability))
			capMark = true
		}

		for i := range container.DelCapabilities {
			capability := container.DelCapabilities[i]
			caps.Drop = append(caps.Drop, v1.Capability(capability))
			capMark = true
		}

		if capMark {
			sc.Capabilities = caps
			mark = true
		}
	}

	if container.SELinux != nil {
		sc.SELinuxOptions = &v1.SELinuxOptions{
			User:  container.SELinux.User,
			Role:  container.SELinux.Role,
			Type:  container.SELinux.Type,
			Level: container.SELinux.Level,
		}
		mark = true
	}

	if !mark {
		return nil, nil
	}
	return sc, nil
}

func revertLifecycle(onStart, preStop *types.Action) (*v1.Lifecycle, error) {
	var lc *v1.Lifecycle

	kubeOnStart, err := revertLifecycleAction(onStart)
	if err != nil {
		return nil, err
	}

	kubePreStop, err := revertLifecycleAction(preStop)
	if err != nil {
		return nil, err
	}

	if onStart != nil || preStop != nil {
		lc = &v1.Lifecycle{}
		lc.PostStart = kubeOnStart
		lc.PreStop = kubePreStop
	}

	return lc, nil
}

func revertLifecycleAction(action *types.Action) (*v1.Handler, error) {
	if action == nil {
		return nil, nil
	}

	handler := &v1.Handler{}

	if action.Command != nil {
		handler.Exec = &v1.ExecAction{
			Command: action.Command,
		}
	}

	if action.Net != nil {
		urlStruct, err := url.Parse(action.Net.URL)
		if err != nil {
			return nil, serrors.InvalidInstanceErrorf(action, "couldn't parse URL: %s", err)
		}
		var host string
		var port intstr.IntOrString

		hostPort := urlStruct.Host
		fields := strings.Split(hostPort, ":")
		if len(fields) == 2 {
			host = fields[0]
			port = intstr.FromString(fields[1])
		} else if len(fields) == 1 {
			host = hostPort
		} else {
			return nil, serrors.InvalidInstanceErrorf(action.Net, "unexpected HostPort %s", action.Net.URL)
		}

		if strings.ToUpper(urlStruct.Scheme) == "HTTP" || strings.ToUpper(urlStruct.Scheme) == "HTTPS" {
			var scheme v1.URIScheme
			if strings.ToUpper(urlStruct.Scheme) == "HTTP" {
				scheme = v1.URISchemeHTTP
			} else {
				scheme = v1.URISchemeHTTPS
			}

			path := urlStruct.Path

			var headers []v1.HTTPHeader
			for i := range action.Net.Headers {
				header := action.Net.Headers[i]
				fields := strings.Split(header, ":")
				if len(fields) != 2 {
					return nil, serrors.InvalidInstanceErrorf(action.Net, "unexpected HTTP Header %s", header)
				}
				kubeHeader := v1.HTTPHeader{
					Name:  fields[0],
					Value: fields[1],
				}
				headers = append(headers, kubeHeader)
			}

			handler.HTTPGet = &v1.HTTPGetAction{
				Scheme:      scheme,
				Path:        path,
				Port:        port,
				Host:        host,
				HTTPHeaders: headers,
			}
		} else if strings.ToUpper(urlStruct.Scheme) == "TCP" {
			handler.TCPSocket = &v1.TCPSocketAction{
				Host: host,
				Port: port,
			}
		} else {
			return nil, serrors.InvalidInstanceErrorf(action.Net, "unexpected URL Scheme %s", urlStruct.Scheme)
		}
	}

	return handler, nil
}

func revertVolumeMounts(mounts []types.VolumeMount) []v1.VolumeMount {
	var kubeMounts []v1.VolumeMount
	for i := range mounts {
		mount := mounts[i]
		kubeMount := v1.VolumeMount{}
		if mount.Propagation != nil {
			kubeMount.MountPropagation = revertMountPropagation(mount.Propagation)
		}
		kubeMount.MountPath = mount.MountPath

		fields := strings.Split(mount.Store, ":")
		if len(fields) == 1 {
			kubeMount.Name = mount.Store
		} else if len(fields) == 2 {
			kubeMount.Name = fields[0]
			if fields[1] == "ro" {
				kubeMount.ReadOnly = true
			} else {
				kubeMount.SubPath = fields[1]
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

func revertMountPropagation(prop *types.MountPropagation) *v1.MountPropagationMode {
	mode := v1.MountPropagationHostToContainer

	if *prop == types.MountPropagationBidirectional {
		mode = v1.MountPropagationBidirectional
	}

	if *prop == types.MountPropagationNone {
		mode = v1.MountPropagationNone
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
			return nil, serrors.InvalidInstanceContextErrorf(err, probe, "parsing URL")
		}
		fields := strings.Split(urlStruct.Host, ":")
		if len(fields) != 2 && len(fields) != 1 {
			return nil, serrors.InvalidInstanceErrorf(urlStruct, "unrecognized Probe Host")
		}
		host := fields[0]
		port := intstr.FromString("80")
		if len(fields) == 2 {
			port = intstr.Parse(fields[1])
		}
		if strings.ToUpper(urlStruct.Scheme) == "TCP" {
			kubeProbe.TCPSocket = &v1.TCPSocketAction{
				Host: host,
				Port: port,
			}
		} else if strings.ToUpper(urlStruct.Scheme) == "HTTP" || strings.ToUpper(urlStruct.Scheme) == "HTTPS" {
			var scheme v1.URIScheme

			if strings.ToLower(urlStruct.Scheme) == "http" {
				scheme = v1.URISchemeHTTP
			} else if strings.ToLower(urlStruct.Scheme) == "https" {
				scheme = v1.URISchemeHTTPS
			} else {
				return nil, serrors.InvalidInstanceErrorf(urlStruct, "unrecognized Probe URL Scheme")
			}

			kubeProbe.HTTPGet = &v1.HTTPGetAction{
				Scheme: scheme,
				Path:   urlStruct.RequestURI(),
				Host:   host,
				Port:   port,
			}

			var headers []v1.HTTPHeader
			for i := range probe.Net.Headers {
				h := probe.Net.Headers[i]
				fields := strings.Split(h, ":")
				if len(fields) != 2 {
					return nil, serrors.InvalidValueErrorf(h, "unrecognized Probe HTTPHeader")
				}
				header := v1.HTTPHeader{
					Name:  fields[0],
					Value: fields[1],
				}
				headers = append(headers, header)
			}
			kubeProbe.HTTPGet.HTTPHeaders = headers
		} else {
			return nil, serrors.InvalidInstanceErrorf(urlStruct, "unrecognized Probe URL")
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
				return requirements, serrors.InvalidInstanceErrorf(cpu, "couldn't parse min quantity: %s", err)
			}
			requests[v1.ResourceCPU] = q
		}

		if cpu.Max != "" {
			q, err := resource.ParseQuantity(cpu.Max)
			if err != nil {
				return requirements, serrors.InvalidInstanceErrorf(cpu, "couldn't parse max quantity: %s", err)
			}
			limits[v1.ResourceCPU] = q
		}
	}

	if mem != nil {
		if mem.Min != "" {
			q, err := resource.ParseQuantity(mem.Min)
			if err != nil {
				return requirements, serrors.InvalidInstanceErrorf(mem, "couldn't parse min quantity: %s", err)
			}
			requests[v1.ResourceMemory] = q
		}

		if mem.Max != "" {
			q, err := resource.ParseQuantity(mem.Max)
			if err != nil {
				return requirements, serrors.InvalidInstanceErrorf(mem, "couldn't parse max quantity: %s", err)
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
		if e.Type == types.EnvValEnvType {
			envVar := v1.EnvVar{
				Name:  e.Val.Key,
				Value: e.Val.Val,
			}
			envVars = append(envVars, envVar)
			continue
		}

		from := e.From

		// ResourceFieldRef
		if strings.Index(from.From, "limits.") == 0 || strings.Index(from.From, "requests.") == 0 {
			envVar := v1.EnvVar{
				Name: from.Key,
				ValueFrom: &v1.EnvVarSource{
					ResourceFieldRef: &v1.ResourceFieldSelector{
						Resource: from.From,
					},
				},
			}
			envVars = append(envVars, envVar)
			continue
		}

		// ConfigMapKeyRef or ConfigMapEnvSource
		if strings.Index(from.From, "config:") == 0 {
			fields := strings.Split(from.From, ":")
			if len(fields) == 3 {
				//ConfigMapKeyRef
				envVar := v1.EnvVar{
					Name: from.Key,
					ValueFrom: &v1.EnvVarSource{
						ConfigMapKeyRef: &v1.ConfigMapKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: fields[1],
							},
							Key:      fields[2],
							Optional: from.Optional(),
						},
					},
				}
				envVars = append(envVars, envVar)
			} else if len(fields) == 2 {
				//ConfigMapEnvSource
				envVarFromSrc := v1.EnvFromSource{
					Prefix: from.Key,
					ConfigMapRef: &v1.ConfigMapEnvSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: fields[1],
						},
						Optional: from.Optional(),
					},
				}
				envsFromSource = append(envsFromSource, envVarFromSrc)
			} else {
				return nil, nil, serrors.InvalidInstanceErrorf(e, "expected either one or two colon-separated values after 'config:'")
			}
			continue
		}

		// SecretKeyRef or SecretEnvSource
		if strings.Index(from.From, "secret:") == 0 {
			fields := strings.Split(from.From, ":")
			if len(fields) == 3 {
				//SecretKeyRef
				envVar := v1.EnvVar{
					Name: from.Key,
					ValueFrom: &v1.EnvVarSource{
						SecretKeyRef: &v1.SecretKeySelector{
							LocalObjectReference: v1.LocalObjectReference{
								Name: fields[1],
							},
							Key:      fields[2],
							Optional: from.Optional(),
						},
					},
				}
				envVars = append(envVars, envVar)
			} else if len(fields) == 2 {
				envVarFromSrc := v1.EnvFromSource{
					Prefix: from.Key,
					SecretRef: &v1.SecretEnvSource{
						LocalObjectReference: v1.LocalObjectReference{
							Name: fields[1],
						},
						Optional: from.Optional(),
					},
				}
				envsFromSource = append(envsFromSource, envVarFromSrc)
			} else {
				return nil, nil, serrors.InvalidInstanceErrorf(e, "expected either one or two colon-separated values after 'secret:'")
			}
			continue
		}

		// FieldRef
		envVar := v1.EnvVar{
			Name: from.Key,
			ValueFrom: &v1.EnvVarSource{
				FieldRef: &v1.ObjectFieldSelector{
					FieldPath: from.From,
				},
			},
		}
		envVars = append(envVars, envVar)
	}

	return envVars, envsFromSource, nil
}

func revertExpose(ports []types.Port) ([]v1.ContainerPort, error) {
	var err error
	var kubeContainerPorts []v1.ContainerPort
	for i := range ports {
		port := ports[i]
		kubePort := v1.ContainerPort{}

		kubePort.Name = port.Name
		kubePort.Protocol = revertProtocol(port.Protocol)

		kubePort.HostPort, err = port.HostPortInt()
		if err != nil {
			return nil, err
		}

		kubePort.ContainerPort, err = port.ContainerPortInt()
		if err != nil {
			return nil, err
		}

		kubeContainerPorts = append(kubeContainerPorts, kubePort)
	}
	return kubeContainerPorts, nil
}

func revertProtocol(kokiProtocol types.Protocol) v1.Protocol {
	return v1.Protocol(strings.ToUpper(string(kokiProtocol)))
}
