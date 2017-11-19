package types

import (
	"encoding/json"
	"reflect"

	apps "k8s.io/api/apps/v1beta2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/short/util"
)

type ReplicaSetWrapper struct {
	ReplicaSet ReplicaSet `json:"replica_set"`
}

type ReplicaSet struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	Replicas        *int32 `json:"replicas,omitempty"`
	MinReadySeconds int32  `json:"ready_seconds,omitempty"`

	Status *apps.ReplicaSetStatus `json:"status,omitempty"`

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

type RSTemplateMetadata struct {
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type RSSelector struct {
	Shorthand string
	Labels    map[string]string
}

func RSTemplateMetadataFromPod(pod *Pod) *RSTemplateMetadata {
	meta := RSTemplateMetadata{}
	meta.Cluster = pod.Cluster
	meta.Name = pod.Name
	meta.Namespace = pod.Namespace
	meta.Labels = pod.Labels
	meta.Annotations = pod.Annotations

	if reflect.DeepEqual(meta, RSTemplateMetadata{}) {
		return nil
	}

	return &meta
}

func (rs *ReplicaSet) SetTemplate(pod *Pod) {
	rs.TemplateMetadata = RSTemplateMetadataFromPod(pod)
	rs.Volumes = pod.Volumes
	rs.Affinity = pod.Affinity
	rs.Containers = pod.Containers
	rs.InitContainers = pod.InitContainers
	rs.DNSPolicy = pod.DNSPolicy
	rs.HostAliases = pod.HostAliases
	rs.HostMode = pod.HostMode
	rs.Hostname = pod.Hostname
	rs.Registries = pod.Registries
	rs.RestartPolicy = pod.RestartPolicy
	rs.SchedulerName = pod.SchedulerName
	rs.Account = pod.Account
	rs.Tolerations = pod.Tolerations
	rs.TerminationGracePeriod = pod.TerminationGracePeriod

	rs.ActiveDeadline = pod.ActiveDeadline
	rs.Node = pod.Node
	rs.Priority = pod.Priority
	rs.Conditions = pod.Conditions
	rs.NodeIP = pod.NodeIP
	rs.StartTime = pod.StartTime
	rs.Msg = pod.Msg
	rs.Phase = pod.Phase
	rs.IP = pod.IP
	rs.QOS = pod.QOS
	rs.Reason = pod.Reason
	rs.FSGID = pod.FSGID
	rs.GIDs = pod.GIDs
}

func (rs *ReplicaSet) GetTemplate() *Pod {
	pod := Pod{}

	if rs.TemplateMetadata != nil {
		pod.Cluster = rs.TemplateMetadata.Cluster
		pod.Name = rs.TemplateMetadata.Name
		pod.Namespace = rs.TemplateMetadata.Namespace
		pod.Labels = rs.TemplateMetadata.Labels
		pod.Annotations = rs.TemplateMetadata.Annotations
	}

	pod.Volumes = rs.Volumes
	pod.Affinity = rs.Affinity
	pod.Containers = rs.Containers
	pod.InitContainers = rs.InitContainers
	pod.DNSPolicy = rs.DNSPolicy
	pod.HostAliases = rs.HostAliases
	pod.HostMode = rs.HostMode
	pod.Hostname = rs.Hostname
	pod.Registries = rs.Registries
	pod.RestartPolicy = rs.RestartPolicy
	pod.SchedulerName = rs.SchedulerName
	pod.Account = rs.Account
	pod.Tolerations = rs.Tolerations
	pod.TerminationGracePeriod = rs.TerminationGracePeriod

	pod.ActiveDeadline = rs.ActiveDeadline
	pod.Node = rs.Node
	pod.Priority = rs.Priority
	pod.Conditions = rs.Conditions
	pod.NodeIP = rs.NodeIP
	pod.StartTime = rs.StartTime
	pod.Msg = rs.Msg
	pod.Phase = rs.Phase
	pod.IP = rs.IP
	pod.QOS = rs.QOS
	pod.Reason = rs.Reason
	pod.FSGID = rs.FSGID
	pod.GIDs = rs.GIDs

	return &pod
}

func (s *RSSelector) UnmarshalJSON(data []byte) error {
	var str string
	strErr := json.Unmarshal(data, &str)
	if strErr == nil {
		s.Shorthand = str
		return nil
	}

	labels := map[string]string{}
	dictErr := json.Unmarshal(data, &labels)
	if dictErr != nil {
		return util.InvalidValueForTypeErrorf(string(data), s, "couldn't parse JSON as string or dictionary: (%s), (%s)", strErr.Error(), dictErr.Error())
	}

	s.Labels = labels
	return nil
}

func (s RSSelector) MarshalJSON() ([]byte, error) {
	if len(s.Shorthand) > 0 {
		b, err := json.Marshal(s.Shorthand)
		if err != nil {
			return nil, util.InvalidInstanceErrorf(s, "couldn't marshal shorthand string to JSON: %s", err.Error())
		}

		return b, nil
	}

	b, err := json.Marshal(s.Labels)
	if err != nil {
		return nil, util.InvalidInstanceErrorf(s, "couldn't marshal labels dictionary to JSON: %s", err.Error())
	}

	return b, nil
}
