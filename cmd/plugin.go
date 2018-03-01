package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

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

	rmPluginCmd = &cobra.Command{
		Use:   "rm",
		Short: "Remove plugin",
		RunE: func(c *cobra.Command, args []string) error {
			return removePlugin(c, args)
		},
	}

	// absolute path to the plugin in the local file system
	absPluginPath string
)

func init() {
	os.MkdirAll(plugin.PluginDir, 0755)

	pluginCmd.AddCommand(lsPluginCmd)
	pluginCmd.AddCommand(installPluginCmd)
	pluginCmd.AddCommand(rmPluginCmd)

	installPluginCmd.Flags().StringVarP(&absPluginPath, "path", "p", "", "absolute path to the plugin in the local file system")
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
		return fmt.Errorf("Plugin %s already exists", name)
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
	return nil
}

func listPlugins(c *cobra.Command, args []string) error {
	files, err := ioutil.ReadDir(plugin.PluginDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	return nil
}

func removePlugin(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("No plugin specified to delete")
	}

	return os.Remove(filepath.Join(plugin.PluginDir, args[0]))
}
