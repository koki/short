package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodWrapper struct {
	Pod Pod `json:"pod"`
}

type Pod struct {
	Version                string            `json:"version,omitempty"`
	Cluster                string            `json:"cluster,omitempty"`
	Name                   string            `json:"name,omitempty"`
	Namespace              string            `json:"namespace,omitempty"`
	Labels                 map[string]string `json:"labels,omitempty"`
	Annotations            map[string]string `json:"annotations,omitempty"`
	Volumes                []Volume          `json:"volumes,omitempty"`
	Affinity               []Affinity        `json:"affinity,omitempty"`
	Containers             []Container       `json:"containers,omitempty"`
	InitContainers         []Container       `json:"init_containers,omitempty"`
	DNSPolicy              DNSPolicy         `json:"dns_policy,omitempty"`
	HostAliases            []string          `json:"host_aliases,omitempty"`
	HostMode               []HostMode        `json:"host_mode,omitempty"`
	Hostname               string            `json:"hostname,omitempty"`
	Registries             []string          `json:"registry_secrets,omitempty"`
	RestartPolicy          RestartPolicy     `json:"restart_policy,omitempty"`
	SchedulerName          string            `json:"scheduler_name,omitempty"`
	Account                string            `json:"account,omitempty"`
	Tolerations            []Toleration      `json:"tolerations,omitempty"`
	TerminationGracePeriod *int64            `json:"termination_grace_period,omitempty"`
	ActiveDeadline         *int64            `json:"active_deadline,omitempty"`
	Node                   string            `json:"node,omitempty"`
	Priority               *Priority         `json:"priority,omitempty"`
	Conditions             []PodCondition    `json:"condition,omitempty"`
	NodeIP                 string            `json:"node_ip,omitempty"`
	StartTime              *metav1.Time      `json:"start_time,omitempty"`
	Msg                    string            `json:"msg,omitempty"`
	Phase                  PodPhase          `json:"phase,omitempty"`
	IP                     string            `json:"ip,omitempty"`
	QOS                    PodQOSClass       `json:"qos,omitempty"`
	Reason                 string            `json:"reason,omitempty"`
	FSGID                  *int64            `json:"fs_gid,omitempty"`
	GIDs                   []int64           `json:"gids,omitempty"`
}

type Priority struct {
	Value *int32 `json:"value,omitempty"`
	Class string `json:"class,omitempty"`
}

type PodCondition struct {
	LastProbeTime      metav1.Time      `json:"last_probe_time,omitempty"`
	LastTransitionTime metav1.Time      `json:"last_transition_time,omitempty"`
	Msg                string           `json:"msg,omitempty"`
	Reason             string           `json:"reason,omitempty"`
	Status             ConditionStatus  `json:"status,omitempty"`
	Type               PodConditionType `json:"type,omitempty"`
}

type PodConditionType string

const (
	PodScheduled           PodConditionType = "PodScheduled"
	PodReady               PodConditionType = "Ready"
	PodInitialized         PodConditionType = "Initialized"
	PodReasonUnschedulable                  = "Unschedulable"
)

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

type DNSPolicy string

const (
	DNSClusterFirstWithHostNet DNSPolicy = "ClusterFirstWithHostNet"
	DNSClusterFirst            DNSPolicy = "ClusterFirst"
	DNSDefault                 DNSPolicy = "Default"
)

type HostMode string

const (
	HostModeNet HostMode = "net"
	HostModePID HostMode = "pid"
	HostModeIPC HostMode = "ipc"
)

type RestartPolicy string

const (
	RestartPolicyAlways    RestartPolicy = "Always"
	RestartPolicyOnFailure RestartPolicy = "OnFailure"
	RestartPolicyNever     RestartPolicy = "Never"
)

type PodPhase string

const (
	PodPending   PodPhase = "Pending"
	PodRunning   PodPhase = "Running"
	PodSucceeded PodPhase = "Succeeded"
	PodFailed    PodPhase = "Failed"
	PodUnknown   PodPhase = "Unknown"
)

type PodQOSClass string

const (
	PodQOSGuaranteed PodQOSClass = "Guaranteed"
	PodQOSBurstable  PodQOSClass = "Burstable"
	PodQOSBestEffort PodQOSClass = "BestEffort"
)
