package util

import (
	"os"

	"github.com/golang/glog"
)

func ExitWithErr(msg interface{}) {
	glog.Error(msg)
	os.Exit(1)
}
