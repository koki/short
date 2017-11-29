package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GITCOMMIT = "HEAD"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version of short server",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("koki/short server: %s\n", GITCOMMIT)
		},
	}
)
