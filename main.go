package main

import (
	"os"

	"github.com/golang/glog"
	"github.com/koki/short/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		glog.Error(err)
		os.Exit(1)
	}
}
