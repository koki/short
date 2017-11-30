package parser

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/util/yaml"
)

// Parse reads input files and then returns a deserialized data structure
func Parse(filenames []string, useStdin bool) ([]map[string]interface{}, error) {
	glog.V(3).Info("validating input does not include both stdin and files")
	if len(filenames) > 0 && useStdin {
		return nil, fmt.Errorf("can only parse from either stdin or files")
	}

	var streams []io.ReadCloser

	if useStdin {
		glog.V(3).Info("reading data from stdin")
		streams = append(streams, os.Stdin)
	} else {
		glog.V(3).Info("reading data from input files")
		s, err := OpenStreamsFromFiles(filenames)
		if err != nil {
			return nil, err
		}
		streams = append(streams, s...)
	}

	glog.V(3).Info("decoding input data")
	return ParseStreams(streams)
}

//parses each stream into a go object and closes the stream once done
func ParseStreams(streams []io.ReadCloser) ([]map[string]interface{}, error) {
	structs := []map[string]interface{}{}

	for i := range streams {
		stream := streams[i]
		defer stream.Close()

		decoder := yaml.NewYAMLOrJSONDecoder(stream, 1024)

		var err error

		for err != io.EOF {
			into := map[string]interface{}{}
			err = decoder.Decode(&into)
			if err != nil && err != io.EOF {
				return nil, err
			}
			if err == nil {
				structs = append(structs, into)
			}
			// TBD: Add support for v1.List type by flattening it
			// and then converting them to individual objects
		}
	}
	return structs, nil
}
