package cmd

import (
	"flag"

	"github.com/golang/glog"
	"github.com/spf13/cobra"

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
	// output denotes the destination of the convereted data
	output string
)

func init() {
	// local flags to root command
	RootCmd.Flags().BoolVarP(&kubeNative, "kube-native", "k", false, "convert to kube-native syntax")
	RootCmd.Flags().StringSliceVarP(&filenames, "filenames", "f", nil, "path or url to input files to read manifests")
	RootCmd.Flags().StringVarP(&output, "output", "o", "", "output to filename instead of stdin")

	// parse the go default flagset to get flags for glog and other packages in future
	RootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	// defaulting this to true so that logs are printed to console
	flag.Set("logtostderr", "true")

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

	useStdin := false

	if len(args) == 1 && args[0] == "-" {
		glog.V(3).Info("using stdin for input data")
		useStdin = true
	}

	// parse input data from one of the sources - files or stdin
	glog.V(3).Info("parsing input data")
	if _, err := parser.Parse(filenames, useStdin); err != nil {
		return util.UsageErrorf(c.CommandPath(), err)
	}
	return nil
}
