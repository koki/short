package cmd

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/koki/short/client"
	"github.com/koki/short/parser"
	"github.com/koki/short/plugin"
	serrors "github.com/koki/structurederrors"
)

var (
	RootCmd = &cobra.Command{
		Use:   "short",
		Short: "Manageable Kubernetes manifests using koki/short",
		Long: `Short converts the api-friendly kubernetes manifests into ops-friendly syntax.

Full documentation available at https://docs.koki.io/short
`,
		RunE: func(c *cobra.Command, args []string) error {
			err := short(c, args)
			if err != nil && !verboseErrors {
				fmt.Fprintln(os.Stderr, "Use flag '--verbose-errors' for more detailed error info.")
			}

			if err != nil {
				return fmt.Errorf(serrors.PrettyError(err))
			}

			return nil
		},
		SilenceUsage: true,
		Example: `

  # Convert existing kubernetes manifestes to shorthand format
  short -f pod.yaml

  # Stream in manifest files
  cat pod.yaml | short -

  # Read from url
  short -f http://spec.com/pod.yaml

  # Convert shorthand file to native syntax
  short --kube-native -f pod_short.yaml
  short -k -f pod_short.yaml

  # Output to file
  short -f pod.yaml > pod_short.yaml

  # Output as yaml* or json
  short -f pod.yaml -o json
`,
	}

	// kubeNative denotes that the conversion must output in kubernetes native syntax
	kubeNative bool
	// filenames holds the input files that are to be converted to shorthand or kuberenetes native syntax
	filenames []string
	// output denotes the destination of the converted data
	output string
	// dryRun denotes that none of the activate installed should be invoked
	dryRun bool
	// verboseErrors denotes that error messages should contain full information instead of just a summary
	verboseErrors bool
	// debugImportsDepth is the number of levels of imports to output debug info for
	debugImportsDepth int
)

const (
	// default value for debugImportsDepth
	defaultDebugImportsDepth = 0
)

func init() {
	// local flags to root command
	RootCmd.Flags().BoolVarP(&kubeNative, "kube-native", "k", false, "convert to kube-native syntax")
	RootCmd.Flags().StringSliceVarP(&filenames, "filenames", "f", nil, "path or url to input files to read manifests")
	RootCmd.Flags().StringVarP(&output, "output", "o", "yaml", "output format (yaml*|json)")
	RootCmd.Flags().BoolVarP(&dryRun, "dry-run", "r", false, "do not invoke any installers")
	RootCmd.Flags().BoolVarP(&verboseErrors, "verbose-errors", "", false, "include more information in errors")
	RootCmd.Flags().IntVarP(&debugImportsDepth, "debug-imports-depth", "", defaultDebugImportsDepth, "how many levels of imports to output debug info for")

	// parse the go default flagset to get flags for glog and other packages in future
	RootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// defaulting this to true so that logs are printed to console
	flag.Set("logtostderr", "true")

	//suppress the incorrect prefix in glog output
	flag.CommandLine.Parse([]string{})

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(pluginCmd)
}

func short(c *cobra.Command, args []string) error {
	var err error
	serrors.SetVerboseErrors(verboseErrors)
	// validate that the user used the command correctly
	glog.V(3).Infof("validating command %q", args)

	if len(args) != 0 {
		if len(args) == 1 {
			if args[0] != "-" {
				//this check ensures that we do not have any dangling args at the end
				return serrors.UsageErrorf(c.CommandPath(), "unexpected value [%s]", args[0])
			} else if len(filenames) > 0 {
				//if '-' is specified at the end, then we expect the value to be streamed in
				return serrors.UsageErrorf(c.CommandPath(), "unexpected value [%s]", args[0])
			}
		} else { //more than one dangling arg left. Abort!
			return serrors.UsageErrorf(c.CommandPath(), "unexpected values %q", args)
		}
	}

	if strings.ToLower(output) != "yaml" && strings.ToLower(output) != "json" {
		return serrors.UsageErrorf("unexpected value %s for -o --output", output)
	}

	useStdin := false
	if len(args) == 1 && args[0] == "-" {
		glog.V(3).Info("using stdin for input data")
		useStdin = true
	}

	var convertedData []interface{}
	if !useStdin && kubeNative {
		// Imports are only supported for normal files in koki syntax.
		kokiModules, err := loadKokiFiles(filenames)
		if err != nil {
			return err
		}

		convertedData, err = convertKokiModules(kokiModules)
		if err != nil {
			return err
		}
	} else {
		// parse input data from one of the sources - files or stdin
		glog.V(3).Info("parsing input data")
		fileDatas := map[string][]map[string]interface{}{}
		if useStdin {
			fileDatas["stdin"], err = parser.Parse(nil, true)
			if err != nil {
				return fmt.Errorf("parsing stdin: %s", err.Error())
			}
		} else {
			for _, filename := range filenames {
				fileDatas[filename], err = parser.Parse([]string{filename}, false)
				if err != nil {
					return fmt.Errorf("parsing %s: %s", filename, err.Error())
				}
			}
		}

		i := 0
		convertedData = []interface{}{}

		for filename, unfilteredData := range fileDatas {
			if err != nil {
				return err
			}

			config := &plugin.AdmitterContext{
				Filename:   filename,
				KubeNative: kubeNative,
			}

			// using context because it allows for easy copying of
			// common info and ability to add arbitrary info
			// per loop
			ctx := context.WithValue(context.Background(), "config", config)
			ctx = context.WithValue(ctx, "index", i)

			data, err := plugin.RunAdmitters(ctx, unfilteredData)
			if err != nil {
				return err
			}

			if kubeNative {
				glog.V(3).Info("converting input to kubernetes native syntax")
				objs, err := client.ConvertKokiMaps(data)
				if err != nil {
					return fmt.Errorf("converting %s: %s", filename, err.Error())
				}
				convertedData = append(convertedData, objs...)
			} else {
				glog.V(3).Info("converting input to koki native syntax")
				objs, err := client.ConvertKubeMaps(data)
				if err != nil {
					return fmt.Errorf("converting %s: %s", filename, err.Error())
				}
				convertedData = append(convertedData, objs...)
			}
			i = i + 1
		}
	}

	buf := &bytes.Buffer{}
	if strings.ToLower(output) == "yaml" {
		glog.V(3).Info("marshalling converted data into yaml")
		err = client.WriteObjsToYamlStream(convertedData, buf)
		if err != nil {
			return err
		}
	} else {
		glog.V(3).Info("marshalling converted data into json")
		err = client.WriteObjsToJSONStream(convertedData, buf)
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s\n", buf.String())

	if dryRun {
		return nil
	}

	return plugin.RunInstallers(buf)
}
