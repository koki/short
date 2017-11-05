package main

import (
	"github.com/golang/glog"
	"github.com/koki/short/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		glog.Fatal(err)
	}
}
