package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/koki/short/plugin"
	"github.com/koki/short/types"
)

var KokiPlugin kokiPlugin

type kokiPlugin struct{}

func (k *kokiPlugin) Admit(ctx context.Context, resource interface{}) (interface{}, error) {
	cfg := ctx.Value("config").(*plugin.AdmitterContext)
	if cfg == nil {
		return nil, fmt.Errorf("Empty context provided")
	}

	if cfg.KubeNative == false {
		return resource, nil
	}

	if cfg.ResourceType == "config_map" {
		configMap := resource.(*types.ConfigMapWrapper)

		data := configMap.ConfigMap.Data

		return generatePerceptorSpec(configMap.ConfigMap.Name, data["config_map_mount_path"], data["name"], data["image"], data["port"], []string{data["cmd"]})
	}
	return resource, nil
}

func (k *kokiPlugin) Install(buf *bytes.Buffer) error {
	cmd := exec.Command("kubectl", "create", "-f", "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		io.Copy(stdin, buf)
		stdin.Close()
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n%s", err, out)
	}

	fmt.Printf("%s", out)

	return nil
}

func generatePerceptorSpec(configMapName, configMapMountPath, appName, image, port string, cmd []string) (interface{}, error) {
	container := types.Container{
		Name:    appName,
		Image:   image,
		Pull:    types.PullAlways,
		Command: cmd,
		Expose: []types.Port{
			{
				ContainerPort: port,
				Protocol:      types.ProtocolTCP,
			},
		},
		CPU: &types.CPU{
			Min: "500m",
		},
		Mem: &types.Mem{
			Min: "2Gi",
		},
		VolumeMounts: []types.VolumeMount{
			{
				Store:     configMapName,
				MountPath: configMapMountPath,
			},
		},
	}

	volume := map[string]types.Volume{
		configMapName: types.Volume{
			ConfigMap: &types.ConfigMapVolume{
				Name: configMapName,
			},
		},
	}

	rc := types.ReplicationControllerWrapper{
		types.ReplicationController{
			Name: appName,
			PodTemplate: types.PodTemplate{
				Volumes: volume,
				Containers: []types.Container{
					container,
				},
			},
			Selector: map[string]string{
				"name": appName,
			},
		},
	}

	return rc, nil
}
