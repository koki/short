package tests

import (
	"bytes"
	"fmt"
	"io"
	_ "io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/koki/short/client"
	"github.com/koki/short/converter"
	"github.com/koki/short/parser"
	"github.com/udhos/equalfile"
)

var cmp *equalfile.Cmp

func TestMain(m *testing.M) {
	cmp = equalfile.New(nil, equalfile.Options{})
	os.Exit(m.Run())
}

func TestPods(t *testing.T) {
	err := testResource("pods", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeployments(t *testing.T) {
	err := testResource("deployments", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPersistentVolumes(t *testing.T) {
	err := testResource("persistent_volumes", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestReplicaSets(t *testing.T) {
	err := testResource("replica_sets", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestReplicationControllers(t *testing.T) {
	err := testResource("replication_controllers", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestServices(t *testing.T) {
	err := testResource("services", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

type filePair struct {
	kubeSpec *os.File
	kokiSpec *os.File
}

func testFuncGenerator(t *testing.T) func(string, filePair) error {
	return func(path string, fp filePair) error {
		if fp.kubeSpec == nil || fp.kokiSpec == nil {
			return nil
		}

		unconvertedKokiObj := &bytes.Buffer{}
		unconvertedKokiStream := TeeReader(fp.kokiSpec, unconvertedKokiObj)

		unconvertedKubeObj := &bytes.Buffer{}
		unconvertedKubeStream := TeeReader(fp.kubeSpec, unconvertedKubeObj)

		// convert kube to koki
		streams := []io.ReadCloser{unconvertedKubeStream}
		objs, err := parser.ParseStreams(streams)
		if err != nil {
			t.Errorf("path %s err %v", path, err)
			return err
		}

		kokiObjs, err := converter.ConvertToKokiNative(objs)
		if err != nil {
			t.Errorf("path %s err %v", path, err)
			return err
		}

		convertedKokiObj := &bytes.Buffer{}

		err = client.WriteObjsToYamlStream(kokiObjs, convertedKokiObj)
		if err != nil {
			t.Errorf("path %s err %v", path, err)
			return err
		}

		// convert koki to kube
		streams = []io.ReadCloser{unconvertedKokiStream}
		objs, err = parser.ParseStreams(streams)
		if err != nil {
			t.Errorf("path %s err %v", path, err)
			return err
		}

		kubeObjs, err := converter.ConvertToKubeNative(objs)
		if err != nil {
			t.Errorf("path %s err %v", path, err)
			return err
		}

		convertedKubeObj := &bytes.Buffer{}

		err = client.WriteObjsToYamlStream(kubeObjs, convertedKubeObj)
		if err != nil {
			t.Errorf("path %s err %v", path, err)
			return err
		}

		//cb, _ := ioutil.ReadAll(convertedKubeObj)
		//ucb, _ := ioutil.ReadAll(unconvertedKubeObj)

		//t.Fatalf("converted=%s \n unconverted=%s \n", string(cb), string(ucb))

		/*equal, err := cmp.CompareReader(convertedKubeObj, unconvertedKubeObj)
		if err != nil {
			t.Errorf("path %s err %v", path, err)
			return err
		}

		if !equal {
			t.Errorf("Failed to translate from Koki To Kube types. Resource Path=%s", path)
			return nil
		}*/

		/*if path == "../testdata/pods/pod_spec_with_affinity" {
			cb, _ := ioutil.ReadAll(convertedKokiObj)
			ucb, _ := ioutil.ReadAll(unconvertedKokiObj)

			t.Fatalf("\nconverted(%d)=\n%s\nunconverted(%d)=\n%s \n", len(cb), cb, len(ucb), ucb)
		}*/

		equal, err := cmp.CompareReader(convertedKokiObj, unconvertedKokiObj)
		if err != nil {
			t.Errorf("path %s err %v", path, err)
			return err
		}

		if !equal {
			t.Errorf("Failed to translate from Kube To Koki types. Resource Path=%s", path)
			return nil
		}
		return nil
	}
}

func testResource(resource string, test func(string, filePair) error) error {

	filePairs := map[string]filePair{}
	root := fmt.Sprintf("../testdata/%s", resource)

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".short.yaml") {
			resourceId := strings.TrimSuffix(path, ".short.yaml")
			fp := filePairs[resourceId]
			kokiSpec, err := os.Open(path)
			if err != nil {
				return err
			}
			fp.kokiSpec = kokiSpec
			filePairs[resourceId] = fp
			test(resourceId, fp)
		} else if strings.HasSuffix(path, ".yaml") {
			resourceId := strings.TrimSuffix(path, ".yaml")
			fp := filePairs[resourceId]
			kubeSpec, err := os.Open(path)
			if err != nil {
				return err
			}
			fp.kubeSpec = kubeSpec
			filePairs[resourceId] = fp
			test(resourceId, fp)
		} else if path == root {
			return nil
		} else {
			return fmt.Errorf("Unrecognized file %s", path)
		}
		return nil
	})
}

// TeeReader returns a Reader that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the write must complete before the read completes.
// Any error encountered while writing is reported as a read error.
func TeeReader(r io.ReadCloser, w io.Writer) io.ReadCloser {
	return &teeReader{r, w}
}

type teeReader struct {
	r io.ReadCloser
	w io.Writer
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

func (t *teeReader) Close() error {
	return t.r.Close()
}
