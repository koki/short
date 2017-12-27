package types

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EventWrapper struct {
	Event `json:"event"`
}

type Event struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	// The object that this event is about.
	InvolvedObject ObjectReference `json:"involved"`

	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`

	// Source::EventSource
	Component string `json:"component,omitempty"`
	Host      string `json:"host,omitempty"`

	// first recorded
	FirstTimestamp metav1.Time `json:"first_recorded,omitempty"`
	// last recorded
	LastTimestamp metav1.Time `json:"last_recorded,omitempty"`
	Count         int32       `json:"count,omitempty"`
	Type          string      `json:"type,omitempty"`
	// first observed
	EventTime metav1.MicroTime `json:"first_observed,omitempty"`

	// Series::*EventSeries
	// Number of occurrences in this series up to the last heartbeat time
	SeriesCount *int32 `json:"series_count,omitempty"`
	// Time of the last occurence observed
	LastObservedTime *metav1.MicroTime `json:"last_observed,omitempty"`
	SeriesState      *EventSeriesState `json:"series_state,omitempty"`

	Action  string           `json:"action,omitempty"`
	Related *ObjectReference `json:"related,omitempty"`
	// controller name
	ReportingController string `json:"reporter"`
	// controller instance ID
	ReportingInstance string `json:"reporter_id"`
}

type EventSeriesState string

const (
	EventSeriesStateOngoing  EventSeriesState = "ongoing"
	EventSeriesStateFinished EventSeriesState = "finished"
	EventSeriesStateUnknown  EventSeriesState = "unknown"
)
