package converters

import (
	"strings"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/koki/short/types"
	"github.com/koki/short/util"
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

	if kokiPV.Status != nil {
		kubePV.Status = *kokiPV.Status
	}

	return kubePV, nil
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
		return v1.PersistentVolumeSource{
			ISCSI: revertISCSIVolume(kokiSource.ISCSI),
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
		return v1.PersistentVolumeSource{
			FlexVolume: revertFlexVolume(kokiSource.Flex),
		}, nil
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

	return v1.PersistentVolumeSource{}, util.InvalidInstanceErrorf(kokiSource, "didn't find any supported volume source")
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
