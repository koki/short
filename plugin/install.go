package plugin

import (
	"bytes"
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/golang/glog"
)

func Install(pluginName string, buf *bytes.Buffer) error {
	glog.V(3).Info("loading plugin ", pluginName)
	loadedPlugin, err := plugin.Open(filepath.Join(PluginDir, pluginName))
	if err != nil {
		glog.Errorf("Error loading path %s: [%v]", filepath.Join(PluginDir, pluginName), err)
		return err
	}

	installerPlugin, err := loadedPlugin.Lookup("KokiPlugin")
	if err != nil {
		glog.Errorf("Error looking up variable %v", err)
		return err
	}

	installer, ok := installerPlugin.(Installer)
	if !ok {
		return fmt.Errorf("Plugin is not of type installer")
	}

	glog.V(3).Info("Installing buffer ", buf.String())
	return installer.Install(buf)
}
