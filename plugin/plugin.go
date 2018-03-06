package plugin

import (
	"bytes"
	"context"
)

/* A koki/short plugin that can install the short resources
 * onto a kubernetes cluster must satisfy the following interface
 */

// This is called for each and every resource that gets translated by the program
type Installer interface {
	Install(*bytes.Buffer) error
}

// This is called with every input file when this plugin is activated
type Admitter interface {
	Admit(context.Context, interface{}) (interface{}, error)
}

// The data structure of the config file that persists plugin info
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

// This context is sent to the admitter plugin. The plugin can use this
// context to obtain information about the resource and the command which
// triggered the plugin to be called
type AdmitterContext struct {
	PluginName   string
	Filename     string
	KubeNative   bool
	ResourceType string
}
