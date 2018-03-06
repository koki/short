package cmd

import (
	"fmt"
	"io"
	_ "io/ioutil"
	"os"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/koki/short/plugin"
)

var (
	pluginCmd = &cobra.Command{
		Use:          "plugin",
		Short:        "Manage short plugins",
		SilenceUsage: true,
		Example: `
 # Install a new plugin named "plugin-name"
 short plugin install plugin-name --path=/absolute/path/to/plugin/binary

 # List installed plugins
 short plugin ls 

 # Remove installed plugin named "plugin-name"
 short rm plugin plugin-name

 # activate plugin
 short plugin activate plugin-name

 # deactivate plugin
 short plugin deactivate plugin-name
`,
	}

	lsPluginCmd = &cobra.Command{
		Use:   "ls",
		Short: "List installed plugins",
		RunE: func(c *cobra.Command, args []string) error {
			return listPlugins(c, args)
		},
	}

	installPluginCmd = &cobra.Command{
		Use:   "install",
		Short: "Install plugin",
		RunE: func(c *cobra.Command, args []string) error {
			return installPlugin(c, args)
		},
	}

	activatePluginCmd = &cobra.Command{
		Use:   "activate",
		Short: "Activate Plugin",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Plugin name not specified")
			}

			name := args[0]
			return togglePlugin(name, true)
		},
	}

	deactivatePluginCmd = &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate Plugin",
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Plugin name not specified")
			}

			name := args[0]
			return togglePlugin(name, false)
		},
	}

	rmPluginCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove plugin",
		RunE: func(c *cobra.Command, args []string) error {
			return removePlugin(c, args)
		},
	}

	// absolute path to the plugin in the local file system
	absPluginPath string

	// activate a plugin while installing
	activate bool
)

func init() {
	os.MkdirAll(plugin.PluginDir, 0755)

	pluginCmd.AddCommand(lsPluginCmd)
	pluginCmd.AddCommand(installPluginCmd)
	pluginCmd.AddCommand(rmPluginCmd)
	pluginCmd.AddCommand(activatePluginCmd)
	pluginCmd.AddCommand(deactivatePluginCmd)

	installPluginCmd.Flags().StringVarP(&absPluginPath, "path", "p", "", "absolute path to the plugin in the local file system")
	installPluginCmd.Flags().BoolVarP(&activate, "activate", "a", false, "activate the plugin as well if installation succeeded")
}

func togglePlugin(pluginName string, toggleType bool) error {
	if toggleType == true {
		return plugin.ActivatePlugin(pluginName)
	}
	return plugin.DeactivatePlugin(pluginName)
}

func installPlugin(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("Plugin name not specified")
	}

	name := args[0]

	if absPluginPath == "" {
		return fmt.Errorf("Plugin path not specified. Use flag --path")
	}

	in, err := os.Open(absPluginPath)
	if err != nil {
		return err
	}
	defer in.Close()

	if _, err := os.Stat(filepath.Join(plugin.PluginDir, name)); err == nil {
		glog.V(3).Infof("Plugin name %s already exists. Reloading...", name)
	}

	out, err := os.Create(filepath.Join(plugin.PluginDir, name))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return plugin.RegisterPlugin(name, activate)
}

func listPlugins(c *cobra.Command, args []string) error {
	plugins := plugin.ListPlugins()

	fmt.Println("PLUGIN NAME \t\t ADMITTER \t INSTALLER \t ACTIVE")

	for k, v := range plugins {
		installer := "false"
		admitter := "false"

		if v.Admitter != nil {
			if v.Admitter.Enabled == true {
				admitter = "true"
				if v.Admitter.Active == true {
					admitter = admitter + "*"
				}
			}
		}

		if v.Installer != nil {
			if v.Installer.Enabled == true {
				installer = "true"
				if v.Installer.Active == true {
					installer = installer + "*"
				}
			}
		}

		active := "false"

		if admitter == "true*" && installer == "true*" {
			active = "true"
			admitter = "true"
			installer = "true"
		}

		fmt.Printf("%-8s \t\t %-6s \t %-6s \t %-6s\n", k, admitter, installer, active)
	}

	return nil
}

func removePlugin(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No plugin specified to delete")
	}

	if err := os.Remove(filepath.Join(plugin.PluginDir, args[0])); err != nil {
		return err
	}

	return plugin.RemovePlugin(args[0])
}
