package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Koki_Event_to_Kube(wrapper *types.EventWrapper) (*v1.Event, error) {
	var err error
	kube := &v1.Event{}
	koki := wrapper.Event

	kube.Name = koki.Name
	kube.Namespace = koki.Namespace
	if len(koki.Version) == 0 {
		kube.APIVersion = "v1"
	} else {
		kube.APIVersion = koki.Version
	}
	kube.Kind = "Event"
	kube.ClusterName = koki.Cluster
	kube.Labels = koki.Labels
	kube.Annotations = koki.Annotations

	involvedObject, err := revertTarget(&koki.InvolvedObject)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "event involved object")
	}
	kube.InvolvedObject = *involvedObject
	kube.Reason = koki.Reason
	kube.Message = koki.Message
	kube.Source.Host = koki.Host
	kube.Source.Component = koki.Component
	kube.FirstTimestamp = koki.FirstTimestamp
	kube.LastTimestamp = koki.LastTimestamp
	kube.Count = koki.Count
	kube.Type = koki.Type
	kube.EventTime = koki.EventTime
	kube.Series, err = revertEventSeries(koki)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "event series")
	}

	kube.Action = koki.Action
	kube.Related, err = revertTarget(koki.Related)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "event related obj")
	}
	kube.ReportingController = koki.ReportingController
	kube.ReportingInstance = koki.ReportingInstance

	return kube, nil
}

func revertEventSeries(koki types.Event) (*v1.EventSeries, error) {
	var err error
	nonempty := false
	series := v1.EventSeries{}
	if koki.SeriesCount != nil {
		series.Count = *koki.SeriesCount
		nonempty = true
	}

	if koki.LastObservedTime != nil {
		series.LastObservedTime = *koki.LastObservedTime
		nonempty = true
	}

	if koki.SeriesState != nil {
		series.State, err = revertEventSeriesState(*koki.SeriesState)
		nonempty = true

		if err != nil {
			return nil, err
		}
	}

	if nonempty {
		return &series, nil
	}

	return nil, nil
}

func revertEventSeriesState(koki types.EventSeriesState) (v1.EventSeriesState, error) {
	switch koki {
	case types.EventSeriesStateOngoing:
		return v1.EventSeriesStateOngoing, nil
	case types.EventSeriesStateFinished:
		return v1.EventSeriesStateFinished, nil
	case types.EventSeriesStateUnknown:
		return v1.EventSeriesStateUnknown, nil
	default:
		return "", serrors.InvalidInstanceError(koki)
	}
}
