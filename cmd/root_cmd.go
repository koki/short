package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/koki/short/converter"
	"github.com/koki/short/parser"
	"github.com/koki/short/util"
)

var (
	RootCmd = &cobra.Command{
		Use:   "short",
		Short: "Manageable Kubernetes manifests using koki/short",
		Long: `Short converts the api-friendly kubernetes manifests into ops-friendly syntax. 
		
Full documentation available at https://docs.koki.io/short
`,
		RunE:         short,
		SilenceUsage: true,
		Example: `
  # Find the shorthand representation of kubernetes objects
  short man pod
  short man deployment

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
  short -f pod.yaml -o pod_short.yaml	
`,
	}

	// kubeNative denotes that the conversion must output in kubernetes native syntax
	kubeNative bool
	// filenames holds the input files that are to be converted to shorthand or kuberenetes native syntax
	filenames []string
	// output denotes the destination of the converted data
	output string
	// silent denotes that the conversion output should not be printed to stdout
	silent bool
)

func init() {
	// local flags to root command
	RootCmd.Flags().BoolVarP(&kubeNative, "kube-native", "k", false, "convert to kube-native syntax")
	RootCmd.Flags().StringSliceVarP(&filenames, "filenames", "f", nil, "path or url to input files to read manifests")
	RootCmd.Flags().StringVarP(&output, "output", "o", "yaml", "output format (yaml*|json)")
	RootCmd.Flags().BoolVarP(&silent, "silent", "s", false, "silence output to stdout")

	// parse the go default flagset to get flags for glog and other packages in future
	RootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// defaulting this to true so that logs are printed to console
	flag.Set("logtostderr", "true")

	//suppress the incorrect prefix in glog output
	flag.CommandLine.Parse([]string{})

	RootCmd.AddCommand(versionCmd)
	//RootCmd.AddCommand(manCommand)
}

func short(c *cobra.Command, args []string) error {
	// validate that the user used the command correctly
	glog.V(3).Infof("validating command %q", args)

	if len(args) != 0 {
		if len(args) == 1 {
			if args[0] != "-" {
				//this check ensures that we do not have any dangling args at the end
				return util.UsageErrorf(c.CommandPath(), "unexpected value [%s]", args[0])
			} else if len(filenames) > 0 {
				//if '-' is specified at the end, then we expect the value to be streamed in
				return util.UsageErrorf(c.CommandPath(), "unexpected value [%s]", args[0])
			}
		} else { //more than one dangling arg left. Abort!
			return util.UsageErrorf(c.CommandPath(), "unexpected values %q", args)
		}
	}

	if strings.ToLower(output) != "yaml" && strings.ToLower(output) != "json" {
		return util.UsageErrorf("unexpected value %s for -o --output", output)
	}

	useStdin := false

	if len(args) == 1 && args[0] == "-" {
		glog.V(3).Info("using stdin for input data")
		useStdin = true
	}

	// parse input data from one of the sources - files or stdin
	glog.V(3).Info("parsing input data")
	data, err := parser.Parse(filenames, useStdin)
	if err != nil {
		return err
	}

	var convertedData interface{}

	if kubeNative {
		glog.V(3).Info("converting input to kubernetes native syntax")
		convertedData, err = converter.ConvertToKubeNative(data)
		if err != nil {
			return err
		}
	} else {
		glog.V(3).Info("converting input to koki native syntax")
		convertedData, err = converter.ConvertToKokiNative(data)
		if err != nil {
			return err
		}
	}

	if silent {
		return nil
	}

	var outputData []byte

	if strings.ToLower(output) == "yaml" {
		glog.V(3).Info("marshalling converted data into yaml")
		outputData, err = yaml.Marshal(convertedData)
		if err != nil {
			return err
		}
	} else {
		glog.V(3).Info("marshalling converted data into json")
		outputData, err = json.MarshalIndent(convertedData, " ", "  ")
		if err != nil {
			return err
		}
	}

	fmt.Printf("%s\n", outputData)

	return nil
}
