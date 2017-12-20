package converters

import (
	"k8s.io/api/core/v1"

	"github.com/koki/short/types"
	serrors "github.com/koki/structurederrors"
)

func Convert_Kube_Event_to_Koki(kube *v1.Event) (*types.EventWrapper, error) {
	var err error
	koki := &types.Event{}

	koki.Name = kube.Name
	koki.Namespace = kube.Namespace
	koki.Version = kube.APIVersion
	koki.Cluster = kube.ClusterName
	koki.Labels = kube.Labels
	koki.Annotations = kube.Annotations

	involvedObject, err := convertTarget(&kube.InvolvedObject)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "event involved object")
	}
	koki.InvolvedObject = *involvedObject
	koki.Reason = kube.Reason
	koki.Message = kube.Message
	koki.Host = kube.Source.Host
	koki.Component = kube.Source.Component
	koki.FirstTimestamp = kube.FirstTimestamp
	koki.LastTimestamp = kube.LastTimestamp
	koki.Count = kube.Count
	koki.Type = kube.Type
	koki.EventTime = kube.EventTime
	if kube.Series != nil {
		koki.SeriesCount = &kube.Series.Count
		koki.LastObservedTime = &kube.Series.LastObservedTime
		state, err := convertEventSeriesState(kube.Series.State)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "event series state")
		}
		koki.SeriesState = &state
	}

	koki.Action = kube.Action
	koki.Related, err = convertTarget(kube.Related)
	if err != nil {
		return nil, serrors.ContextualizeErrorf(err, "event related obj")
	}
	koki.ReportingController = kube.ReportingController
	koki.ReportingInstance = kube.ReportingInstance

	return &types.EventWrapper{
		Event: *koki,
	}, nil
}

func convertEventSeriesState(kube v1.EventSeriesState) (types.EventSeriesState, error) {
	switch kube {
	case v1.EventSeriesStateOngoing:
		return types.EventSeriesStateOngoing, nil
	case v1.EventSeriesStateFinished:
		return types.EventSeriesStateFinished, nil
	case v1.EventSeriesStateUnknown:
		return types.EventSeriesStateUnknown, nil
	default:
		return "", serrors.InvalidInstanceError(kube)
	}
}
