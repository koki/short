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

func Convert_Koki_PersistentVolume_to_Kube_v1_PersistentVolume(pv *types.PersistentVolumeWrapper) (*v1.PersistentVolume, error) {
	var err error
	kubePV := &v1.PersistentVolume{}
	kokiPV := pv.PersistentVolume

	kubePV.Name = kokiPV.Name
	kubePV.Namespace = kokiPV.Namespace
	if len(kokiPV.Version) == 0 {
		kubePV.APIVersion = "v1"
	} else {
		kubePV.APIVersion = kokiPV.Version
	}
	kubePV.Kind = "PersistentVolume"
	kubePV.ClusterName = kokiPV.Cluster
	kubePV.Labels = kokiPV.Labels
	kubePV.Annotations = kokiPV.Annotations

	kubeSpec := &kubePV.Spec
	kubeSpec.Capacity = revertCapacity(kokiPV.Storage)
	kubeSpec.PersistentVolumeSource, err = revertPersistentVolumeSource(kokiPV.PersistentVolumeSource)
	if err != nil {
		return nil, err
	}
	if kokiPV.AccessModes != nil {
		kubeSpec.AccessModes = kokiPV.AccessModes.Modes
	}
	kubeSpec.ClaimRef = kokiPV.Claim
	kubeSpec.PersistentVolumeReclaimPolicy = revertReclaimPolicy(kokiPV.ReclaimPolicy)
	kubeSpec.StorageClassName = kokiPV.StorageClass
	if len(kokiPV.MountOptions) > 0 {
		kubeSpec.MountOptions = strings.Split(kokiPV.MountOptions, ",")
	}

	kubePV.Status, err = revertPersistentVolumeStatus(kokiPV.PersistentVolumeStatus)
	if err != nil {
		return nil, err
	}

	return kubePV, nil
}

func revertPersistentVolumeStatus(kokiStatus types.PersistentVolumeStatus) (v1.PersistentVolumeStatus, error) {
	phase, err := revertPersistentVolumePhase(kokiStatus.Phase)
	if err != nil {
		return v1.PersistentVolumeStatus{}, err
	}
	return v1.PersistentVolumeStatus{
		Phase:   phase,
		Message: kokiStatus.Message,
		Reason:  kokiStatus.Reason,
	}, nil
}

func revertPersistentVolumePhase(kokiPhase types.PersistentVolumePhase) (v1.PersistentVolumePhase, error) {
	if kokiPhase == "" {
		return "", nil
	}
	switch kokiPhase {
	case types.VolumePending:
		return v1.VolumePending, nil
	case types.VolumeAvailable:
		return v1.VolumeAvailable, nil
	case types.VolumeBound:
		return v1.VolumeBound, nil
	case types.VolumeReleased:
		return v1.VolumeReleased, nil
	case types.VolumeFailed:
		return v1.VolumeFailed, nil
	default:
		return v1.VolumeFailed, serrors.InvalidValueErrorf(kokiPhase, "unrecognized status (phase) for persistent volume")
	}
}

func revertSecretReference(kokiRef *types.SecretReference) *v1.SecretReference {
	if kokiRef == nil {
		return nil
	}

	return &v1.SecretReference{
		Namespace: kokiRef.Namespace,
		Name:      kokiRef.Name,
	}
}

func revertCephFSPersistentSecretFileOrRef(kokiSecret *types.CephFSPersistentSecretFileOrRef) (string, *v1.SecretReference) {
	if kokiSecret == nil {
		return "", nil
	}

	if len(kokiSecret.File) > 0 {
		return kokiSecret.File, nil
	}

	return "", revertSecretReference(kokiSecret.Ref)
}

func revertSecretRefToObjectRef(kokiRef *types.SecretReference) *v1.ObjectReference {
	if kokiRef == nil {
		return nil
	}
	return &v1.ObjectReference{
		Namespace: kokiRef.Namespace,
		Name:      kokiRef.Name,
	}
}

func revertPersistentVolumeSource(kokiSource types.PersistentVolumeSource) (v1.PersistentVolumeSource, error) {
	if kokiSource.GcePD != nil {
		return v1.PersistentVolumeSource{
			GCEPersistentDisk: revertGcePDVolume(kokiSource.GcePD),
		}, nil
	}
	if kokiSource.AwsEBS != nil {
		return v1.PersistentVolumeSource{
			AWSElasticBlockStore: revertAwsEBSVolume(kokiSource.AwsEBS),
		}, nil
	}
	if kokiSource.HostPath != nil {
		source, err := revertHostPathVolume(kokiSource.HostPath)
		if err != nil {
			return v1.PersistentVolumeSource{}, err
		}
		return v1.PersistentVolumeSource{
			HostPath: source,
		}, nil
	}
	if kokiSource.Glusterfs != nil {
		return v1.PersistentVolumeSource{
			Glusterfs: revertGlusterfsVolume(kokiSource.Glusterfs),
		}, nil
	}
	if kokiSource.NFS != nil {
		return v1.PersistentVolumeSource{
			NFS: revertNFSVolume(kokiSource.NFS),
		}, nil
	}
	if kokiSource.ISCSI != nil {
		source := kokiSource.ISCSI
		return v1.PersistentVolumeSource{
			ISCSI: &v1.ISCSIPersistentVolumeSource{
				TargetPortal:      source.TargetPortal,
				IQN:               source.IQN,
				Lun:               source.Lun,
				ISCSIInterface:    source.ISCSIInterface,
				FSType:            source.FSType,
				ReadOnly:          source.ReadOnly,
				Portals:           source.Portals,
				DiscoveryCHAPAuth: source.DiscoveryCHAPAuth,
				SessionCHAPAuth:   source.SessionCHAPAuth,
				SecretRef:         revertSecretReference(source.SecretRef),
				InitiatorName:     util.StringPtrOrNil(source.InitiatorName),
			},
		}, nil
	}
	if kokiSource.Cinder != nil {
		return v1.PersistentVolumeSource{
			Cinder: revertCinderVolume(kokiSource.Cinder),
		}, nil
	}
	if kokiSource.FibreChannel != nil {
		return v1.PersistentVolumeSource{
			FC: revertFibreChannelVolume(kokiSource.FibreChannel),
		}, nil
	}
	if kokiSource.Flocker != nil {
		return v1.PersistentVolumeSource{
			Flocker: revertFlockerVolume(kokiSource.Flocker),
		}, nil
	}
	if kokiSource.Flex != nil {
		flexVol, err := revertFlexPersistentVolume(kokiSource.Flex)
		return v1.PersistentVolumeSource{
			FlexVolume: flexVol,
		}, err
	}
	if kokiSource.Vsphere != nil {
		return v1.PersistentVolumeSource{
			VsphereVolume: revertVsphereVolume(kokiSource.Vsphere),
		}, nil
	}
	if kokiSource.Quobyte != nil {
		return v1.PersistentVolumeSource{
			Quobyte: revertQuobyteVolume(kokiSource.Quobyte),
		}, nil
	}
	if kokiSource.AzureDisk != nil {
		source, err := revertAzureDiskVolume(kokiSource.AzureDisk)
		if err != nil {
			return v1.PersistentVolumeSource{}, err
		}
		return v1.PersistentVolumeSource{
			AzureDisk: source,
		}, nil
	}
	if kokiSource.PhotonPD != nil {
		return v1.PersistentVolumeSource{
			PhotonPersistentDisk: revertPhotonPDVolume(kokiSource.PhotonPD),
		}, nil
	}
	if kokiSource.Portworx != nil {
		return v1.PersistentVolumeSource{
			PortworxVolume: revertPortworxVolume(kokiSource.Portworx),
		}, nil
	}
	if kokiSource.RBD != nil {
		source := kokiSource.RBD
		return v1.PersistentVolumeSource{
			RBD: &v1.RBDPersistentVolumeSource{
				CephMonitors: source.CephMonitors,
				RBDImage:     source.RBDImage,
				FSType:       source.FSType,
				RBDPool:      source.RBDPool,
				RadosUser:    source.RadosUser,
				Keyring:      source.Keyring,
				SecretRef:    revertSecretReference(source.SecretRef),
				ReadOnly:     source.ReadOnly,
			},
		}, nil
	}
	if kokiSource.CephFS != nil {
		source := kokiSource.CephFS
		secretFile, secretRef := revertCephFSPersistentSecretFileOrRef(source.SecretFileOrRef)
		return v1.PersistentVolumeSource{
			CephFS: &v1.CephFSPersistentVolumeSource{
				Monitors:   source.Monitors,
				Path:       source.Path,
				User:       source.User,
				SecretFile: secretFile,
				SecretRef:  secretRef,
				ReadOnly:   source.ReadOnly,
			},
		}, nil
	}
	if kokiSource.AzureFile != nil {
		source := kokiSource.AzureFile
		return v1.PersistentVolumeSource{
			AzureFile: &v1.AzureFilePersistentVolumeSource{
				SecretName:      source.Secret.Name,
				ShareName:       source.ShareName,
				ReadOnly:        source.ReadOnly,
				SecretNamespace: util.StringPtrOrNil(source.Secret.Namespace),
			},
		}, nil
	}
	if kokiSource.ScaleIO != nil {
		source := kokiSource.ScaleIO
		mode, err := revertScaleIOStorageMode(source.StorageMode)
		if err != nil {
			return v1.PersistentVolumeSource{}, serrors.ContextualizeErrorf(err, "ScaleIO storage mode")
		}
		return v1.PersistentVolumeSource{
			ScaleIO: &v1.ScaleIOPersistentVolumeSource{
				Gateway:          source.Gateway,
				System:           source.System,
				SecretRef:        revertSecretReference(&source.SecretRef),
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
	if kokiSource.Local != nil {
		source := kokiSource.Local
		return v1.PersistentVolumeSource{
			Local: &v1.LocalVolumeSource{
				Path: source.Path,
			},
		}, nil
	}
	if kokiSource.StorageOS != nil {
		source := kokiSource.StorageOS
		return v1.PersistentVolumeSource{
			StorageOS: &v1.StorageOSPersistentVolumeSource{
				VolumeName:      source.VolumeName,
				VolumeNamespace: source.VolumeNamespace,
				FSType:          source.FSType,
				ReadOnly:        source.ReadOnly,
				SecretRef:       revertSecretRefToObjectRef(source.SecretRef),
			},
		}, nil
	}
	if kokiSource.CSI != nil {
		source := kokiSource.CSI
		return v1.PersistentVolumeSource{
			CSI: &v1.CSIPersistentVolumeSource{
				Driver:       source.Driver,
				VolumeHandle: source.VolumeHandle,
				ReadOnly:     source.ReadOnly,
			},
		}, nil
	}

	return v1.PersistentVolumeSource{}, serrors.InvalidInstanceErrorf(kokiSource, "didn't find any supported volume source")
}

func revertFlexPersistentVolume(source *types.FlexVolume) (*v1.FlexPersistentVolumeSource, error) {
	secretRef, err := revertSecretRef(source.SecretRef)
	if err != nil {
		return nil, err
	}
	return &v1.FlexPersistentVolumeSource{
		Driver:    source.Driver,
		FSType:    source.FSType,
		SecretRef: secretRef,
		ReadOnly:  source.ReadOnly,
		Options:   source.Options,
	}, nil
}

func revertSecretRef(kokiRef string) (*v1.SecretReference, error) {
	if len(kokiRef) == 0 {
		return nil, nil
	}
	parts := strings.Split(kokiRef, "/")

	var name, namespace string

	name = parts[0]
	if len(parts) == 2 {
		namespace = parts[0]
		name = parts[1]
	}

	if len(parts) > 2 {
		return nil, fmt.Errorf("Unknown format for secret reference %s", kokiRef)
	}

	return &v1.SecretReference{
		Namespace: namespace,
		Name:      name,
	}, nil
}

func revertReclaimPolicy(kokiPolicy types.PersistentVolumeReclaimPolicy) v1.PersistentVolumeReclaimPolicy {
	return v1.PersistentVolumeReclaimPolicy(strings.Title(string(kokiPolicy)))
}

func revertCapacity(kokiStorage *resource.Quantity) v1.ResourceList {
	if kokiStorage == nil {
		return nil
	}

	kubeCapacity := v1.ResourceList{}
	kubeCapacity[v1.ResourceStorage] = *kokiStorage

	return kubeCapacity
}
