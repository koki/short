package types

import (
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/koki/json"
	serrors "github.com/koki/structurederrors"
)

type HorizontalPodAutoscalerWrapper struct {
	HPA HorizontalPodAutoscaler `json:"hpa"`
}

type HorizontalPodAutoscaler struct {
	Version     string            `json:"version,omitempty"`
	Cluster     string            `json:"cluster,omitempty"`
	Name        string            `json:"name,omitempty"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`

	HorizontalPodAutoscalerSpec `json:",inline"`

	HorizontalPodAutoscalerStatus `json:",inline"`
}

type HorizontalPodAutoscalerSpec struct {
	ScaleTargetRef                 CrossVersionObjectReference `json:"ref"`
	MinReplicas                    *int32                      `json:"min,omitempty"`
	MaxReplicas                    int32                       `json:"max"`
	TargetCPUUtilizationPercentage *int32                      `json:"percent_cpu,omitempty"`
}

// current status of a horizontal pod autoscaler
type HorizontalPodAutoscalerStatus struct {
	ObservedGeneration              *int64       `json:"generation_observed,omitempty"`
	LastScaleTime                   *metav1.Time `json:"last_scaling,omitempty"`
	CurrentReplicas                 int32        `json:"current,omitempty"`
	DesiredReplicas                 int32        `json:"desired,omitempty"`
	CurrentCPUUtilizationPercentage *int32       `json:"current_percent_cpu,omitempty"`
}

type CrossVersionObjectReference struct {
	Kind       string
	Name       string
	APIVersion string
}

func (r CrossVersionObjectReference) VersionKind() string {
	if len(r.APIVersion) > 0 {
		return fmt.Sprintf("%s.%s", r.APIVersion, r.Kind)
	}

	return r.Kind
}

func (r CrossVersionObjectReference) MarshalJSON() ([]byte, error) {
	str := fmt.Sprintf("%s:%s", r.VersionKind(), r.Name)
	return json.Marshal(str)
}

func (r *CrossVersionObjectReference) UnmarshalJSON(data []byte) error {
	var str string
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	segments := strings.Split(str, ":")
	if len(segments) != 2 {
		return serrors.InvalidValueForTypeErrorf(str, r, "expected 'version.kind:name' OR 'kind:name'")
	}

	r.Name = segments[1]

	splitAt := strings.LastIndex(segments[0], ".")
	if splitAt >= 0 {
		r.APIVersion = segments[0][:splitAt]
		r.Kind = segments[0][splitAt+1:]
	} else {
		r.Kind = segments[0]
	}

	return nil
}
