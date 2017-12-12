package parser

import (
	"io"
	"os"

	"github.com/golang/glog"

	serrors "github.com/koki/structurederrors"
)

func OpenStreamsFromFiles(filenames []string) ([]io.ReadCloser, error) {
	readers := []io.ReadCloser{}

	for _, name := range filenames {
		glog.V(5).Infof("opening file %s for reading", name)
		f, err := os.Open(name)
		if err != nil {
			return nil, serrors.ContextualizeErrorf(err, "opening file %s", name)
		}

		readers = append(readers, f)
	}

	return readers, nil
}
