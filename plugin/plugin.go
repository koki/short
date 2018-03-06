package plugin

import (
	"bytes"
	"fmt"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/golang/glog"
	"github.com/spf13/viper"

	serrors "github.com/koki/structurederrors"
)

/* A koki/short plugin that can install the short resources
 * onto a kubernetes cluster must satisfy the following interface
 */

var (
	config      map[string]PluginConfig
	pluginMutex sync.Mutex
)

func init() {
	config = map[string]PluginConfig{}

	viper.SetConfigFile(filepath.Join(PluginDir, ConfigFile))
	viper.SetDefault("plugins", config)
	if err := viper.ReadInConfig(); err != nil {
		viper.SafeWriteConfig()
	} else {
		if err := viper.UnmarshalKey("plugins", &config); err != nil {
			panic(fmt.Errorf("Invalid config file: %v", err))
		}
	}
}

const (
	ConfigFile = "config.yaml"
	PluginDir  = ".short-plugins"
)

type Installer interface {
	// This is called for each and every resource that gets translated by the program
	Install(*bytes.Buffer) error
}

type Admitter interface {
	// This is called with every input file when this plugin is activated
	Admit(string, []map[string]interface{}, bool, map[string]interface{}) (interface{}, error)
}

type PluginConfig struct {
	Installer *InstallerConfig
	Admitter  *AdmitterConfig
}

type InstallerConfig struct {
	Enabled bool
	Active  bool
}

type AdmitterConfig struct {
	Enabled bool
	Active  bool
}

func RegisterPlugin(pluginName string, activate bool) error {
	newPlugin, err := plugin.Open(filepath.Join(PluginDir, pluginName))
	if err != nil {
		return err
	}

	installer := false
	admitter := false

	kokiPlugin, err := newPlugin.Lookup("KokiPlugin")
	if err != nil {
		return err
	}

	if _, ok := kokiPlugin.(Installer); ok {
		installer = true
	}

	if _, ok := kokiPlugin.(Admitter); ok {
		admitter = true
	}

	if !installer && !admitter {
		return fmt.Errorf("Plugin is not an Installer or Admitter")
	}

	// start critical section
	pluginMutex.Lock()

	config[pluginName] = PluginConfig{
		Installer: &InstallerConfig{
			Enabled: installer,
		},
		Admitter: &AdmitterConfig{
			Enabled: admitter,
		},
	}

	viper.Set("plugins", config)
	err = viper.WriteConfig()

	// end critical section
	pluginMutex.Unlock()

	if err != nil {
		return err
	}

	if activate {
		return ActivatePlugin(pluginName)
	}
	return nil
}

// Active plugins are called by short during its lifecycle
func ActivatePlugin(pluginName string) error {
	// start critical section
	pluginMutex.Lock()

	toActivate, ok := config[pluginName]
	if !ok {
		return fmt.Errorf("Plugin %s is not installed", pluginName)
	}

	if toActivate.Installer != nil {
		toActivate.Installer.Active = true
	}
	if toActivate.Admitter != nil {
		toActivate.Admitter.Active = true
	}

	viper.Set("plugins", config)
	err := viper.WriteConfig()

	// end critical section
	pluginMutex.Unlock()

	return err
}

// Deactivate plugin disabled a plugin from being called during the lifecycle of short
func DeactivatePlugin(pluginName string) error {
	// start critical section
	pluginMutex.Lock()

	toActivate, ok := config[pluginName]
	if !ok {
		return fmt.Errorf("Plugin %s is not installed", pluginName)
	}

	if toActivate.Installer != nil {
		toActivate.Installer.Active = false
	}
	if toActivate.Admitter != nil {
		toActivate.Admitter.Active = false
	}

	viper.Set("plugins", config)
	err := viper.WriteConfig()

	// end critical section
	pluginMutex.Unlock()

	return err
}

// callers of this function SHOULD not update the map
func ListPlugins() map[string]PluginConfig {
	return config
}

func RemovePlugin(pluginName string) error {
	var err error

	// start critical section
	pluginMutex.Lock()

	// delete the entry from map
	delete(config, pluginName)
	viper.Set("plugins", config)
	err = viper.WriteConfig()

	// end critical section
	pluginMutex.Unlock()

	return err
}

func RunInstallers(buf *bytes.Buffer) error {
	for k, v := range config {
		if v.Installer != nil {
			if v.Installer.Active == true {
				glog.V(3).Infof("Installing using plugin %s", k)
				err := Install(k, buf)
				if err != nil {
					return serrors.ContextualizeErrorf(err, "Error installing resources with plugin %s", k)
				}
			}
		}
	}
	return nil
}

func RunAdmitters(filename string, data []map[string]interface{}, toKubernetes bool, cache map[string]interface{}) ([]map[string]interface{}, error) {
	filteredDatas := []map[string]interface{}{}
	for k, v := range config {
		if v.Admitter != nil {
			if v.Admitter.Active == true {
				glog.V(3).Infof("Filtering resources using pluging %s", k)
				filteredData, err := Admit(k, filename, data, toKubernetes, cache)
				if err != nil {
					return nil, serrors.ContextualizeErrorf(err, "Error admitting resources with plugin %s", k)
				}
				glog.Errorf("filtered data %+v", filteredData)
				filteredDatas = append(filteredDatas, filteredData...)
			}
		}
	}
	return filteredDatas, nil
}
