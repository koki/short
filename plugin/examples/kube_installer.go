package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/koki/short/types"
)

var KokiPlugin kokiPlugin

type kokiPlugin struct{}

func (k *kokiPlugin) Admit(filename string, data []map[string]interface{}, toKube bool, cache map[string]interface{}) (interface{}, error) {
	if !toKube {
		return data, nil
	}

	for _, resource := range data {
		for k := range resource {
			if k == "config_map" {
				cm := resource[k].(map[string]interface{})
				if _, ok := cm["data"]; !ok {
					continue
				}
				resourceData := cm["data"].(map[string]interface{})
				configMapName := ""
				configMapMountPath := ""
				appName := ""
				image := ""
				port := ""
				cmd := ""

				if cmName, ok := resourceData["config_map"]; ok {
					configMapName = cmName.(string)
				} else {
					continue
				}
				if mountPath, ok := resourceData["config_map_mount_path"]; ok {
					configMapMountPath = mountPath.(string)
				} else {
					continue
				}
				if appNameInterface, ok := resourceData["name"]; ok {
					appName = appNameInterface.(string)
				} else {
					continue
				}
				if imageInterface, ok := resourceData["image"]; ok {
					image = imageInterface.(string)
				} else {
					continue
				}
				if portInterface, ok := resourceData["port"]; ok {
					port = portInterface.(string)
				} else {
					continue
				}
				if cmdInterface, ok := resourceData["cmd"]; ok {
					cmd = cmdInterface.(string)
				} else {
					continue
				}

				return generatePodSpec(configMapName, configMapMountPath, appName, image, port, []string{cmd})
			}
		}
	}
	return data, nil
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

func generatePodSpec(configMapName, configMapMountPath, appName, image, port string, cmd []string) (interface{}, error) {
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
