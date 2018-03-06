package plugin

import (
	"context"
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/golang/glog"

	"github.com/koki/json"
	"github.com/koki/short/parser"
)

func Admit(ctx context.Context, objects []map[string]interface{}) ([]map[string]interface{}, error) {
	cfg := ctx.Value("config").(*AdmitterContext)
	if cfg == nil {
		return nil, fmt.Errorf("empty config for admitter")
	}

	glog.V(3).Info("loading plugin ", cfg.PluginName)
	loadedPlugin, err := plugin.Open(filepath.Join(PluginDir, cfg.PluginName))
	if err != nil {
		glog.Errorf("Error loading path %s: [%v]", filepath.Join(PluginDir, cfg.PluginName), err)
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

	glog.V(3).Infof("Admitting resources from filename %s ", cfg.Filename)

	filteredObjs := []map[string]interface{}{}
	for i := range objects {
		untypedObject := objects[i]

		var object interface{}
		var err error

		if cfg.KubeNative {
			object, err = parser.ParseKokiNativeObject(untypedObject)

			if len(untypedObject) != 1 {
				return nil, fmt.Errorf("Invalid input data")
			}

			for k := range untypedObject {
				cfg.ResourceType = k
			}

		} else {
			object, err = parser.ParseSingleKubeNative(untypedObject)
		}

		if err != nil {
			return nil, err
		}

		admissionCtx := context.WithValue(ctx, "config", cfg)
		result, err := admitter.Admit(admissionCtx, object)
		if err != nil {
			return nil, err
		}

		b, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}

		filteredObj := map[string]interface{}{}
		err = json.Unmarshal(b, &filteredObj)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling admitter output to map[string]interface{} %v", err)
		}

		filteredObjs = append(filteredObjs, filteredObj)
	}

	return filteredObjs, nil
}
