package converters

import (
	"github.com/koki/short/types"

	"k8s.io/api/core/v1"
)

func Convert_Kube_v1_Service_to_Koki_Service(service *v1.Service) (*types.Service, error) {
	return nil, nil
}
