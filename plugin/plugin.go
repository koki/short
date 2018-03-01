package plugin

/* A koki/short plugin that can install the short resources
 * onto a kubernetes cluster must satisfy the following interface
 */

const (
	PluginDir = ".short-plugins"
)

type Installer interface {
	// This is called for each and every resource that gets translated by the program
	Install(interface{}) error
}
