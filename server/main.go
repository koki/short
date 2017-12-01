package main

import (
	"os"

	"github.com/koki/short/server/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
