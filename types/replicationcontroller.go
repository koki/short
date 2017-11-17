package types

import (
	"reflect"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ReplicationControllerWrapper struct {
	ReplicationController ReplicationController `json:"replication_controller"`
}

type ReplicationController struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Replicas        *int32 `json:"replicas,omitempty"`
	MinReadySeconds int32  `json:"ready_seconds,omitempty"`

	Status *v1.ReplicationControllerStatus `json:"status,omitempty"`

	// Selector and the Template's Labels are expected to be equal
	// if both exist, so we standardize on using the Template's labels.
	Selector map[string]string `json:"selector,omitempty"`

	// Template fields
	TemplateMetadata       *RCTemplateMetadata `json:"pod_meta,omitempty"`
	Volumes                []Volume            `json:"volumes,omitempty"`
	Affinity               []Affinity          `json:"affinity,omitempty"`
	Containers             []Container         `json:"containers,omitempty"`
	InitContainers         []Container         `json:"init_containers,omitempty"`
	DNSPolicy              DNSPolicy           `json:"dns_policy,omitempty"`
	HostAliases            []string            `json:"host_aliases,omitempty"`
	HostMode               []HostMode          `json:"host_mode,omitempty"`
	Hostname               string              `json:"hostname,omitempty"`
	Registries             []string            `json:"registry_secrets,omitempty"`
	RestartPolicy          RestartPolicy       `json:"restart_policy,omitempty"`
	SchedulerName          string              `json:"scheduler_name,omitempty"`
	Account                string              `json:"account,omitempty"`
	Tolerations            []Toleration        `json:"tolerations,omitempty"`
	TerminationGracePeriod *int64              `json:"termination_grace_period,omitempty"`
	ActiveDeadline         *int64              `json:"active_deadline,omitempty"`
	Node                   string              `json:"node,omitempty"`
	Priority               *Priority           `json:"priority,omitempty"`
	Conditions             []PodCondition      `json:"condition,omitempty"`
	NodeIP                 string              `json:"node_ip,omitempty"`
	StartTime              *metav1.Time        `json:"start_time,omitempty"`
	Msg                    string              `json:"msg,omitempty"`
	Phase                  PodPhase            `json:"phase,omitempty"`
	IP                     string              `json:"ip,omitempty"`
	QOS                    PodQOSClass         `json:"qos,omitempty"`
	Reason                 string              `json:"reason,omitempty"`
	FSGID                  *int64              `json:"fs_gid,omitempty"`
	GIDs                   []int64             `json:"gids,omitempty"`
}

type RCTemplateMetadata struct {
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

func RCTemplateMetadataFromPod(pod *Pod) *RCTemplateMetadata {
	meta := RCTemplateMetadata{}
	meta.Cluster = pod.Cluster
	meta.Name = pod.Name
	meta.Namespace = pod.Namespace
	meta.Annotations = pod.Annotations

	if reflect.DeepEqual(meta, RCTemplateMetadata{}) {
		return nil
	}

	return &meta
}

func (rc *ReplicationController) SetTemplate(pod *Pod) {
	rc.TemplateMetadata = RCTemplateMetadataFromPod(pod)

	rc.Selector = pod.Labels

	rc.Volumes = pod.Volumes
	rc.Affinity = pod.Affinity
	rc.Containers = pod.Containers
	rc.InitContainers = pod.InitContainers
	rc.DNSPolicy = pod.DNSPolicy
	rc.HostAliases = pod.HostAliases
	rc.HostMode = pod.HostMode
	rc.Hostname = pod.Hostname
	rc.Registries = pod.Registries
	rc.RestartPolicy = pod.RestartPolicy
	rc.SchedulerName = pod.SchedulerName
	rc.Account = pod.Account
	rc.Tolerations = pod.Tolerations
	rc.TerminationGracePeriod = pod.TerminationGracePeriod

	rc.ActiveDeadline = pod.ActiveDeadline
	rc.Node = pod.Node
	rc.Priority = pod.Priority
	rc.Conditions = pod.Conditions
	rc.NodeIP = pod.NodeIP
	rc.StartTime = pod.StartTime
	rc.Msg = pod.Msg
	rc.Phase = pod.Phase
	rc.IP = pod.IP
	rc.QOS = pod.QOS
	rc.Reason = pod.Reason
	rc.FSGID = pod.FSGID
	rc.GIDs = pod.GIDs
}

func (rc *ReplicationController) GetTemplate() *Pod {
	pod := Pod{}

	if rc.TemplateMetadata != nil {
		pod.Cluster = rc.TemplateMetadata.Cluster
		pod.Name = rc.TemplateMetadata.Name
		pod.Namespace = rc.TemplateMetadata.Namespace
		pod.Annotations = rc.TemplateMetadata.Annotations
	}

	pod.Labels = rc.Selector

	pod.Volumes = rc.Volumes
	pod.Affinity = rc.Affinity
	pod.Containers = rc.Containers
	pod.InitContainers = rc.InitContainers
	pod.DNSPolicy = rc.DNSPolicy
	pod.HostAliases = rc.HostAliases
	pod.HostMode = rc.HostMode
	pod.Hostname = rc.Hostname
	pod.Registries = rc.Registries
	pod.RestartPolicy = rc.RestartPolicy
	pod.SchedulerName = rc.SchedulerName
	pod.Account = rc.Account
	pod.Tolerations = rc.Tolerations
	pod.TerminationGracePeriod = rc.TerminationGracePeriod

	pod.ActiveDeadline = rc.ActiveDeadline
	pod.Node = rc.Node
	pod.Priority = rc.Priority
	pod.Conditions = rc.Conditions
	pod.NodeIP = rc.NodeIP
	pod.StartTime = rc.StartTime
	pod.Msg = rc.Msg
	pod.Phase = rc.Phase
	pod.IP = rc.IP
	pod.QOS = rc.QOS
	pod.Reason = rc.Reason
	pod.FSGID = rc.FSGID
	pod.GIDs = rc.GIDs

	return &pod
}
