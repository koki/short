package types

import (
	apps "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type DeploymentWrapper struct {
	Deployment Deployment `json:"deployment"`
}

type Deployment struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Replicas       *int32              `json:"replicas,omitempty"`
	Recreate       bool                `json:"recreate,omitempty"`
	MaxUnavailable *intstr.IntOrString `json:"max_unavailable,omitempty"`
	MaxSurge       *intstr.IntOrString `json:"max_extra,omitempty"`

	MinReadySeconds         int32  `json:"min_ready,omitempty"`
	RevisionHistoryLimit    *int32 `json:"max_revs,omitempty"`
	Paused                  bool   `json:"paused,omitempty"`
	ProgressDeadlineSeconds *int32 `json:"progress_deadline,omitempty"`

	Status *apps.DeploymentStatus `json:"status,omitempty"`

	// Selector in ReplicaSet can express more complex rules than just matching
	// pod labels, so it needs its own field (unlike in ReplicationController).
	// Leaving it blank has the same effect as omitting Selector in RC.
	Selector *RSSelector `json:"selector,omitempty"`

	// Template fields
	TemplateMetadata       *RSTemplateMetadata `json:"pod_meta,omitempty"`
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

func (d *Deployment) SetTemplate(pod *Pod) {
	d.TemplateMetadata = RSTemplateMetadataFromPod(pod)
	d.Volumes = pod.Volumes
	d.Affinity = pod.Affinity
	d.Containers = pod.Containers
	d.InitContainers = pod.InitContainers
	d.DNSPolicy = pod.DNSPolicy
	d.HostAliases = pod.HostAliases
	d.HostMode = pod.HostMode
	d.Hostname = pod.Hostname
	d.Registries = pod.Registries
	d.RestartPolicy = pod.RestartPolicy
	d.SchedulerName = pod.SchedulerName
	d.Account = pod.Account
	d.Tolerations = pod.Tolerations
	d.TerminationGracePeriod = pod.TerminationGracePeriod

	d.ActiveDeadline = pod.ActiveDeadline
	d.Node = pod.Node
	d.Priority = pod.Priority
	d.Conditions = pod.Conditions
	d.NodeIP = pod.NodeIP
	d.StartTime = pod.StartTime
	d.Msg = pod.Msg
	d.Phase = pod.Phase
	d.IP = pod.IP
	d.QOS = pod.QOS
	d.Reason = pod.Reason
	d.FSGID = pod.FSGID
	d.GIDs = pod.GIDs
}

func (d *Deployment) GetTemplate() *Pod {
	pod := Pod{}

	if d.TemplateMetadata != nil {
		pod.Cluster = d.TemplateMetadata.Cluster
		pod.Name = d.TemplateMetadata.Name
		pod.Namespace = d.TemplateMetadata.Namespace
		pod.Labels = d.TemplateMetadata.Labels
		pod.Annotations = d.TemplateMetadata.Annotations
	}

	pod.Volumes = d.Volumes
	pod.Affinity = d.Affinity
	pod.Containers = d.Containers
	pod.InitContainers = d.InitContainers
	pod.DNSPolicy = d.DNSPolicy
	pod.HostAliases = d.HostAliases
	pod.HostMode = d.HostMode
	pod.Hostname = d.Hostname
	pod.Registries = d.Registries
	pod.RestartPolicy = d.RestartPolicy
	pod.SchedulerName = d.SchedulerName
	pod.Account = d.Account
	pod.Tolerations = d.Tolerations
	pod.TerminationGracePeriod = d.TerminationGracePeriod

	pod.ActiveDeadline = d.ActiveDeadline
	pod.Node = d.Node
	pod.Priority = d.Priority
	pod.Conditions = d.Conditions
	pod.NodeIP = d.NodeIP
	pod.StartTime = d.StartTime
	pod.Msg = d.Msg
	pod.Phase = d.Phase
	pod.IP = d.IP
	pod.QOS = d.QOS
	pod.Reason = d.Reason
	pod.FSGID = d.FSGID
	pod.GIDs = d.GIDs

	return &pod
}
