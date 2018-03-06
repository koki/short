package plugin

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/golang/glog"
	"github.com/spf13/viper"

	serrors "github.com/koki/structurederrors"
)

const (
	ConfigFile = "config.yaml"
	PluginDir  = ".short-plugins"
)

var (
	config      map[string]PluginConfig
	pluginMutex sync.Mutex
)

func init() {
	config = map[string]PluginConfig{}

	viper.SetConfigFile(filepath.Join(PluginDir, ConfigFile))
	viper.SetDefault("plugins", config)

	err := viper.ReadInConfig()
	if err != nil {
		// config doesn't already exist or read error. Create if it doens't exist.
		viper.SafeWriteConfig()
		return
	}

	// read existing config
	err = viper.UnmarshalKey("plugins", &config)
	if err != nil {
		panic(fmt.Errorf("Invalid config file: %v", err))
	}
}

func getActiveAdmitters() []string {
	activeAdmitters := []string{}
	for k, v := range config {
		if v.Admitter != nil {
			if v.Admitter.Active == true {
				activeAdmitters = append(activeAdmitters, k)
			}
		}
	}
	return activeAdmitters
}

func getActiveInstallers() []string {
	activeInstallers := []string{}
	for k, v := range config {
		if v.Installer != nil {
			if v.Installer.Active == true {
				activeInstallers = append(activeInstallers, k)
			}
		}
	}
	return activeInstallers
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

func RunAdmitters(ctx context.Context, objects []map[string]interface{}) ([]map[string]interface{}, error) {
	if len(getActiveAdmitters()) == 0 {
		return objects, nil
	}

	filtered := []map[string]interface{}{}
	cfg := ctx.Value("config").(*AdmitterContext)
	if cfg == nil {
		return nil, fmt.Errorf("Empty admitter config")
	}

	for _, admitter := range getActiveAdmitters() {
		cfg.PluginName = admitter
		admissionCtx := context.WithValue(ctx, "config", cfg)
		filteredData, err := Admit(admissionCtx, objects)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "Error admitting resources with plugin %s", admitter)
		}
		filtered = append(filtered, filteredData...)
	}
	return filtered, nil
}
