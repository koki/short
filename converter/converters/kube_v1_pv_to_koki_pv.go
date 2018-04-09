package converters

import (
	"fmt"
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_v1_PersistentVolume_to_Koki_PersistentVolume(kubePV *v1.PersistentVolume) (*types.PersistentVolumeWrapper, error) {
	var err error
	kokiPV := &types.PersistentVolume{}

	kokiPV.Name = kubePV.Name
	kokiPV.Namespace = kubePV.Namespace
	kokiPV.Version = kubePV.APIVersion
	kokiPV.Cluster = kubePV.ClusterName
	kokiPV.Labels = kubePV.Labels
	kokiPV.Annotations = kubePV.Annotations

	kubeSpec := kubePV.Spec
	kokiPV.Storage, err = convertCapacity(kubeSpec.Capacity)
	if err != nil {
		return nil, err
	}

	kokiPV.PersistentVolumeSource, err = convertPersistentVolumeSource(kubeSpec.PersistentVolumeSource)
	if err != nil {
		return nil, err
	}
	if len(kubeSpec.AccessModes) > 0 {
		kokiPV.AccessModes = &types.AccessModes{
			Modes: kubeSpec.AccessModes,
		}
	}
	kokiPV.Claim = kubeSpec.ClaimRef
	kokiPV.ReclaimPolicy = convertReclaimPolicy(kubeSpec.PersistentVolumeReclaimPolicy)
	kokiPV.StorageClass = kubeSpec.StorageClassName
	if len(kubeSpec.MountOptions) > 0 {
		kokiPV.MountOptions = strings.Join(kubeSpec.MountOptions, ",")
	}

	kokiPV.PersistentVolumeStatus, err = convertPersistentVolumeStatus(kubePV.Status)

	return &types.PersistentVolumeWrapper{
		PersistentVolume: *kokiPV,
	}, nil
}

func convertPersistentVolumeStatus(kubeStatus v1.PersistentVolumeStatus) (types.PersistentVolumeStatus, error) {
	phase, err := convertPersistentVolumePhase(kubeStatus.Phase)
	if err != nil {
		return types.PersistentVolumeStatus{}, err
	}
	return types.PersistentVolumeStatus{
		Phase:   phase,
		Message: kubeStatus.Message,
		Reason:  kubeStatus.Reason,
	}, nil
}

func convertPersistentVolumePhase(kubePhase v1.PersistentVolumePhase) (types.PersistentVolumePhase, error) {
	switch kubePhase {
	case v1.VolumePending:
		return types.VolumePending, nil
	case v1.VolumeAvailable:
		return types.VolumeAvailable, nil
	case v1.VolumeBound:
		return types.VolumeBound, nil
	case v1.VolumeReleased:
		return types.VolumeReleased, nil
	case v1.VolumeFailed:
		return types.VolumeFailed, nil
	default:
		return types.VolumeFailed, serrors.InvalidValueErrorf(kubePhase, "unrecognized status (phase) for persistent volume")
	}
}

func convertSecretReference(kubeRef *v1.SecretReference) *types.SecretReference {
	if kubeRef == nil {
		return nil
	}

	return &types.SecretReference{
		Namespace: kubeRef.Namespace,
		Name:      kubeRef.Name,
	}
}

func convertCephFSPersistentSecretFileOrRef(kubeFile string, kubeRef *v1.SecretReference) *types.CephFSPersistentSecretFileOrRef {
	if len(kubeFile) > 0 {
		return &types.CephFSPersistentSecretFileOrRef{
			File: kubeFile,
		}
	}

	if kubeRef != nil {
		return &types.CephFSPersistentSecretFileOrRef{
			Ref: convertSecretReference(kubeRef),
		}
	}

	return nil
}

func convertObjectRefToSecretRef(kubeRef *v1.ObjectReference) *types.SecretReference {
	if kubeRef == nil {
		return nil
	}

	return &types.SecretReference{
		Namespace: kubeRef.Namespace,
		Name:      kubeRef.Name,
	}
}

func convertPersistentVolumeSource(kubeSource v1.PersistentVolumeSource) (types.PersistentVolumeSource, error) {
	if kubeSource.GCEPersistentDisk != nil {
		return types.PersistentVolumeSource{
			GcePD: convertGcePDVolume(kubeSource.GCEPersistentDisk),
		}, nil
	}
	if kubeSource.AWSElasticBlockStore != nil {
		return types.PersistentVolumeSource{
			AwsEBS: convertAwsEBSVolume(kubeSource.AWSElasticBlockStore),
		}, nil
	}
	if kubeSource.HostPath != nil {
		source, err := convertHostPathVolume(kubeSource.HostPath)
		if err != nil {
			return types.PersistentVolumeSource{}, err
		}
		return types.PersistentVolumeSource{
			HostPath: source,
		}, nil
	}
	if kubeSource.Glusterfs != nil {
		return types.PersistentVolumeSource{
			Glusterfs: convertGlusterfsVolume(kubeSource.Glusterfs),
		}, nil
	}
	if kubeSource.NFS != nil {
		return types.PersistentVolumeSource{
			NFS: convertNFSVolume(kubeSource.NFS),
		}, nil
	}
	if kubeSource.ISCSI != nil {
		source := kubeSource.ISCSI
		return types.PersistentVolumeSource{
			ISCSI: &types.ISCSIPersistentVolume{
				TargetPortal:      source.TargetPortal,
				IQN:               source.IQN,
				Lun:               source.Lun,
				ISCSIInterface:    source.ISCSIInterface,
				FSType:            source.FSType,
				ReadOnly:          source.ReadOnly,
				Portals:           source.Portals,
				DiscoveryCHAPAuth: source.DiscoveryCHAPAuth,
				SessionCHAPAuth:   source.SessionCHAPAuth,
				SecretRef:         convertSecretReference(source.SecretRef),
				InitiatorName:     util.FromStringPtr(source.InitiatorName),
			},
		}, nil
	}
	if kubeSource.Cinder != nil {
		return types.PersistentVolumeSource{
			Cinder: convertCinderVolume(kubeSource.Cinder),
		}, nil
	}
	if kubeSource.FC != nil {
		return types.PersistentVolumeSource{
			FibreChannel: convertFibreChannelVolume(kubeSource.FC),
		}, nil
	}
	if kubeSource.Flocker != nil {
		return types.PersistentVolumeSource{
			Flocker: convertFlockerVolume(kubeSource.Flocker),
		}, nil
	}
	if kubeSource.FlexVolume != nil {
		return types.PersistentVolumeSource{
			Flex: convertFlexPersistentVolume(kubeSource.FlexVolume),
		}, nil
	}
	if kubeSource.VsphereVolume != nil {
		return types.PersistentVolumeSource{
			Vsphere: convertVsphereVolume(kubeSource.VsphereVolume),
		}, nil
	}
	if kubeSource.Quobyte != nil {
		return types.PersistentVolumeSource{
			Quobyte: convertQuobyteVolume(kubeSource.Quobyte),
		}, nil
	}
	if kubeSource.AzureDisk != nil {
		source, err := convertAzureDiskVolume(kubeSource.AzureDisk)
		if err != nil {
			return types.PersistentVolumeSource{}, err
		}
		return types.PersistentVolumeSource{
			AzureDisk: source,
		}, nil
	}
	if kubeSource.PhotonPersistentDisk != nil {
		return types.PersistentVolumeSource{
			PhotonPD: convertPhotonPDVolume(kubeSource.PhotonPersistentDisk),
		}, nil
	}
	if kubeSource.PortworxVolume != nil {
		return types.PersistentVolumeSource{
			Portworx: convertPortworxVolume(kubeSource.PortworxVolume),
		}, nil
	}
	if kubeSource.RBD != nil {
		source := kubeSource.RBD
		return types.PersistentVolumeSource{
			RBD: &types.RBDPersistentVolume{
				CephMonitors: source.CephMonitors,
				RBDImage:     source.RBDImage,
				FSType:       source.FSType,
				RBDPool:      source.RBDPool,
				RadosUser:    source.RadosUser,
				Keyring:      source.Keyring,
				SecretRef:    convertSecretReference(source.SecretRef),
				ReadOnly:     source.ReadOnly,
			},
		}, nil
	}
	if kubeSource.CephFS != nil {
		source := kubeSource.CephFS
		secretFileOrRef := convertCephFSPersistentSecretFileOrRef(source.SecretFile, source.SecretRef)
		return types.PersistentVolumeSource{
			CephFS: &types.CephFSPersistentVolume{
				Monitors:        source.Monitors,
				Path:            source.Path,
				User:            source.User,
				SecretFileOrRef: secretFileOrRef,
				ReadOnly:        source.ReadOnly,
			},
		}, nil
	}
	if kubeSource.AzureFile != nil {
		source := kubeSource.AzureFile
		return types.PersistentVolumeSource{
			AzureFile: &types.AzureFilePersistentVolume{
				Secret: types.SecretReference{
					Name:      source.SecretName,
					Namespace: util.FromStringPtr(source.SecretNamespace),
				},
				ShareName: source.ShareName,
				ReadOnly:  source.ReadOnly,
			},
		}, nil
	}
	if kubeSource.ScaleIO != nil {
		source := kubeSource.ScaleIO
		mode, err := convertScaleIOStorageMode(source.StorageMode)
		if err != nil {
			return types.PersistentVolumeSource{}, err
		}
		secret := convertSecretReference(source.SecretRef)
		if secret == nil {
			return types.PersistentVolumeSource{}, serrors.InvalidInstanceErrorf(source, "secret is required for ScaleIO volume")
		}
		return types.PersistentVolumeSource{
			ScaleIO: &types.ScaleIOPersistentVolume{
				Gateway:          source.Gateway,
				System:           source.System,
				SecretRef:        *secret,
				SSLEnabled:       source.SSLEnabled,
				ProtectionDomain: source.ProtectionDomain,
				StoragePool:      source.StoragePool,
				StorageMode:      mode,
				VolumeName:       source.VolumeName,
				FSType:           source.FSType,
				ReadOnly:         source.ReadOnly,
			},
		}, nil
	}
	if kubeSource.Local != nil {
		source := kubeSource.Local
		return types.PersistentVolumeSource{
			Local: &types.LocalVolume{
				Path: source.Path,
			},
		}, nil
	}
	if kubeSource.StorageOS != nil {
		source := kubeSource.StorageOS
		return types.PersistentVolumeSource{
			StorageOS: &types.StorageOSPersistentVolume{
				VolumeName:      source.VolumeName,
				VolumeNamespace: source.VolumeNamespace,
				FSType:          source.FSType,
				ReadOnly:        source.ReadOnly,
				SecretRef:       convertObjectRefToSecretRef(source.SecretRef),
			},
		}, nil
	}
	if kubeSource.CSI != nil {
		source := kubeSource.CSI
		return types.PersistentVolumeSource{
			CSI: &types.CSIPersistentVolume{
				Driver:       source.Driver,
				VolumeHandle: source.VolumeHandle,
				ReadOnly:     source.ReadOnly,
			},
		}, nil
	}

	return types.PersistentVolumeSource{}, serrors.InvalidInstanceErrorf(kubeSource, "didn't find any supported volume source")
}

func convertFlexPersistentVolume(source *v1.FlexPersistentVolumeSource) *types.FlexVolume {
	return &types.FlexVolume{
		Driver:    source.Driver,
		FSType:    source.FSType,
		SecretRef: convertSecretRef(source.SecretRef),
		ReadOnly:  source.ReadOnly,
		Options:   source.Options,
	}
}

func convertSecretRef(kubeRef *v1.SecretReference) string {
	if kubeRef == nil {
		return ""
	}
	if len(kubeRef.Namespace) == 0 {
		return kubeRef.Name
	}

	return fmt.Sprintf("%s/%s", kubeRef.Namespace, kubeRef.Name)
}

func convertReclaimPolicy(kubePolicy v1.PersistentVolumeReclaimPolicy) types.PersistentVolumeReclaimPolicy {
	return types.PersistentVolumeReclaimPolicy(strings.ToLower(string(kubePolicy)))
}

func convertCapacity(kubeCapacity v1.ResourceList) (*resource.Quantity, error) {
	if len(kubeCapacity) == 0 {
		return nil, nil
	}

	for res, quantity := range kubeCapacity {
		if res == v1.ResourceStorage {
			return &quantity, nil
		}
	}

	return nil, serrors.InvalidInstanceErrorf(kubeCapacity, "only supports Storage resource")
}
