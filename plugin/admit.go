package plugin

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/golang/glog"
)

func Admit(pluginName string, filename string, datas []map[string]interface{}, toKube bool, cache map[string]interface{}) ([]map[string]interface{}, error) {
	glog.V(3).Info("loading plugin ", pluginName)
	loadedPlugin, err := plugin.Open(filepath.Join(PluginDir, pluginName))
	if err != nil {
		glog.Errorf("Error loading path %s: [%v]", filepath.Join(PluginDir, pluginName), err)
		return nil, err
	}

	admitterPlugin, err := loadedPlugin.Lookup("KokiPlugin")
	if err != nil {
		glog.Errorf("Error looking up variable %v", err)
		return nil, err
	}

	admitter, ok := admitterPlugin.(Admitter)
	if !ok {
		return nil, fmt.Errorf("Plugin is not of type admitter")
	}

	glog.V(3).Infof("Admitting resources from filename %s ", filename)
	result, err := admitter.Admit(filename, datas, toKube, cache)
	if err != nil {
		return nil, err
	}

	toRet := map[string]interface{}{}
	toRets := []map[string]interface{}{}

	b, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &toRet)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling admitter output to map[string]interface{} %v", err)
	}
	toRets = append(toRets, toRet)

	return toRets, nil
}
