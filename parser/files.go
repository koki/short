package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/glog"
)

func OpenStreamsFromFiles(filenames []string) ([]io.ReadCloser, error) {
	readers := []io.ReadCloser{}

	for _, name := range filenames {
		glog.V(5).Infof("opening file %s for reading", name)
		f, err := os.Open(name)
		if err != nil {
			return nil, fmt.Errorf("failed opening file (%s): %s", name, err.Error())
		}

		readers = append(readers, f)
	}

	return readers, nil
}
