package imports

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

// ReadYamls read a yaml file of potentially multiple documents.
func ReadYamls(filename string) ([]interface{}, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	yamls := SplitYaml(contents)
	objs := make([]interface{}, len(yamls))
	for ix, y := range yamls {
		obj := map[string]interface{}{}
		mapErr := yaml.Unmarshal(y, &obj)
		if mapErr == nil {
			objs[ix] = obj
			continue
		}

		str := ""
		strErr := yaml.Unmarshal(y, &str)
		if strErr == nil {
			objs[ix] = str
			continue
		}

		return nil, fmt.Errorf("(%s)\n(%s)", mapErr.Error(), strErr.Error())
	}

	return objs, nil
}

// SplitYaml split multi-document yaml file.
func SplitYaml(contents []byte) [][]byte {
	return bytes.Split(contents, []byte("\n---\n"))
}
