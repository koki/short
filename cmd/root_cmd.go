package cmd

import (
	"flag"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use:   "short",
		Short: "Manageable Kubernetes manifests using koki/short",
		Long: `Short converts the api-friendly kubernetes manifests into ops-friendly syntax. 
		
Full documentation available at https://docs.koki.io/short
`,
		RunE: short,
		Example: `
#Find the shorthand representation of kubernetes objects
short man pod							#opens man page for pod syntax and conversion
short man deployment						#opens man page for deployment syntax and conversion

#Convert existing kubernetes manifestes to shorthand format
short -f pod.yaml						#converts the contents of pod.yaml to shorthand syntax

#Stream in manifest files
cat pod.yaml | short -						#reads from stdin instead of from an input file

#Read from url
short -f http://spec.com/pod.yaml				#reads from URL instead of from an input file

#Convert shorthand file to native syntax
short --kube-native -f pod_short.yaml				#converts the contents of pod_short.yaml to native syntax
short -k -f pod_short.yaml					#converts the contents of pod_short.yaml to native syntax

#Output to file 
short -f pod.yaml -o pod_short.yaml				#prints the output into file named pod_short.yaml
`,
	}

	//verbosity denotes the amount of information output as log info messages
	verbosity int
	//kubeNative denotes that the conversion must output in kubernetes native syntax
	kubeNative bool
	//filenames holds the input files that are to be converted to shorthand or kuberenetes native syntax
	filenames []string
	//output denotes the destination of the convereted data
	output string
)

func init() {
	//local flags to root command
	RootCmd.Flags().BoolVarP(&kubeNative, "kube-native", "k", false, "convert to kube-native syntax")
	RootCmd.Flags().StringSliceVarP(&filenames, "filenames", "f", nil, "path or url to input files to read manifests")
	RootCmd.Flags().StringVarP(&output, "output", "o", "", "output to filename instead of stdin")

	//parse the go default flagset to get flags for glog and other packages in future
	RootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)

	//defaulting this to true so that logs are printed to console
	flag.Set("logtostderr", "true")

	RootCmd.AddCommand(versionCmd)
	//RootCmd.AddCommand(manCommand)
}

func short(c *cobra.Command, args []string) error {
	return nil
}
