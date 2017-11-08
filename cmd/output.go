package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/koki/short/util/fs"
)

func writerForOutputFlag(output string) (out io.Writer, useYaml bool, err error) {
	flg := strings.ToLower(output)
	if flg == "yaml" {
		out = os.Stdout
		useYaml = true
		return
	} else if flg == "json" {
		out = os.Stdout
		useYaml = false
		return
	}

	ext := strings.ToLower(filepath.Ext(output))
	useYaml = ext != ".json"
	err = fs.MkdirP(output)
	if err != nil {
		glog.Error("Failed creating parent dir for (%s)", output)
		return
	}

	out, err = os.Create(output)
	if err != nil {
		glog.Error("Failed creating file (%s)", output)
		return
	}

	return
}

func writeAsMultiDoc(marshal func(interface{}) ([]byte, error), out io.Writer, objs []interface{}) error {
	first := true
	for _, obj := range objs {
		if !first {
			_, err := out.Write([]byte("---\n"))
			if err != nil {
				glog.Error("Failed writing separator: ---")
				return err
			}
		}

		y, err := marshal(obj)
		if err != nil {
			glog.Error("Failed marshaling: %#v", obj)
			return err
		}

		_, err = out.Write(y)
		if err != nil {
			glog.Error("Failed writing: %#v", y)
			return err
		}

		first = false
	}

	return nil
}
