package parser

import (
	"io"
	"os"

	"github.com/golang/glog"
)

func OpenStreamsFromFiles(filenames []string) ([]io.ReadCloser, error) {
	readers := []io.ReadCloser{}

	for i := range filenames {
		name := filenames[i]
		glog.V(5).Infof("opening file %s for reading", name)
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}

		readers = append(readers, f)
	}

	return readers, nil
}
