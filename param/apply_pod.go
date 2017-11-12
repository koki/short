package param

import (
	"fmt"
	"regexp"

	"github.com/koki/short/parser"
	"github.com/koki/short/types"
	"github.com/koki/short/util"
)

func ApplyPodParams(params map[string]interface{}, wrapper *types.PodWrapper) error {
	pod := &wrapper.Pod
	if name, ok := params["name"]; ok {
		if name, ok := name.(string); ok {
			pod.Name = name
			if pod.Labels != nil {
				pod.Labels["name"] = name
			} else {
				pod.Labels = map[string]string{
					"name": name,
				}
			}
		} else {
			return util.PrettyTypeError(params, "expected string for 'name'")
		}
	}

	err := applyVolumeMountParams(params, pod)
	if err != nil {
		return err
	}

	return nil
}

type volumeMountParam struct {
	ContainerName string
	Volumes       []interface{}
}

// containers.${CONTAINER_NAME}.volume_mounts:
func parseVolumeMountParam(key string, value interface{}) (*volumeMountParam, error) {
	re := regexp.MustCompile(`^containers\.([^\.]*)\.volume_mounts$`)

	matches := re.FindStringSubmatch(key)
	if len(matches) == 0 {
		return nil, fmt.Errorf("(%s) isn't a volume mount param", key)
	}

	vmParam := &volumeMountParam{
		ContainerName: matches[1],
	}

	var volumes []interface{}
	if values, ok := value.([]interface{}); ok {
		// value is multiple volumes
		volumes = values
	} else {
		// value is just one volume
		volumes = []interface{}{value}
	}

	vmParam.Volumes = make([]interface{}, len(volumes))
	for i, volume := range volumes {
		kokiVolume, err := parser.ParseKokiNativeObject(volume)
		if err != nil {
			return nil, err
		}
		vmParam.Volumes[i] = kokiVolume
	}

	return vmParam, nil
}

func applyVolumeMountParams(params map[string]interface{}, pod *types.Pod) error {
	for param, value := range params {
		vmParam, err := parseVolumeMountParam(param, value)
		if err != nil {
			continue
		}

		for _, volume := range vmParam.Volumes {
			// TODO: non-persistent volumes
			switch volume := volume.(type) {
			case *types.PersistentVolumeWrapper:
				err := applyPersistentVolumeMount(vmParam, volume, pod)
				if err != nil {
					return err
				}
			default:
				return util.PrettyTypeError(volume, "unsupported type for volume mount")
			}
		}
	}

	return nil
}

func applyPersistentVolumeMount(vmParam *volumeMountParam, volume *types.PersistentVolumeWrapper, pod *types.Pod) error {
	for i, container := range pod.Containers {
		if container.Name == vmParam.ContainerName {
			volumeName := volume.PersistentVolume.Name
			volumeMount := &types.VolumeMount{
				MountPath: "/" + volumeName,
				Store:     volumeName,
			}
			container.VolumeMounts = append(container.VolumeMounts, *volumeMount)
			pod.Containers[i] = container
		}
	}

	return nil
}
