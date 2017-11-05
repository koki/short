package parser

import (
	"fmt"

	"github.com/golang/glog"
)

// Parse reads input files and then returns a deserialized data structure
func Parse(filenames []string, useStdin bool) (interface{}, error) {
	glog.V(3).Info("validating input does not include both stdin and files")
	if len(filenames) > 0 && useStdin {
		return nil, fmt.Errorf("can only parse from either stdin or files")
	}

	if useStdin {
		return nil, nil
	}
	return nil, nil
}
