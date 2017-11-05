package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	GITCOMMIT = "HEAD"

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "prints the version of short",
		Run: func(*cobra.Command, []string) {
			fmt.Printf("koki/shorthand: %s\n", GITCOMMIT)
		},
	}
)
