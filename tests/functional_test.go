package tests

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/kr/pretty"
	"github.com/udhos/equalfile"

	"github.com/koki/short/client"
	"github.com/koki/short/converter"
	"github.com/koki/short/parser"
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

func TestJobs(t *testing.T) {
	err := testResource("jobs", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDaemonSets(t *testing.T) {
	err := testResource("daemon_sets", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCronJobs(t *testing.T) {
	err := testResource("cron_jobs", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPVCs(t *testing.T) {
	err := testResource("pvcs", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestStatefulSets(t *testing.T) {
	err := testResource("stateful_sets", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestStorageClass(t *testing.T) {
	err := testResource("stateful_sets", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestConfigMap(t *testing.T) {
	err := testResource("config_maps", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestIngress(t *testing.T) {
	err := testResource("ingress", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

type filePair struct {
	kubeSpec *os.File
	kokiSpec *os.File
}

func objectDiffString(a, b interface{}) string {
	return strings.Join(pretty.Diff(a, b), "\n")
}

func testFuncGenerator(t *testing.T) func(string, filePair) error {
	return func(path string, fp filePair) error {
		if fp.kubeSpec == nil || fp.kokiSpec == nil {
			return nil
		}

		unconvertedKoki, err := ioutil.ReadAll(fp.kokiSpec)
		if err != nil {
			t.Errorf("failed to read koki at %s", path)
			return err
		}
		unconvertedKube, err := ioutil.ReadAll(fp.kubeSpec)
		if err != nil {
			t.Errorf("failed to read kube at %s", path)
			return err
		}

		err = testKubeToKoki(path, unconvertedKube, unconvertedKoki, t)
		if err != nil {
			t.Errorf("failed kube -> koki")
			return err
		}

		err = testKokiToKube(path, unconvertedKoki, unconvertedKube, t)
		if err != nil {
			t.Errorf("failed koki -> kube")
			return err
		}

		return nil
	}
}

func testKubeToKoki(path string, unconvertedKube, unconvertedKoki []byte, t *testing.T) error {
	// convert kube to koki
	streams := []io.ReadCloser{ioutil.NopCloser(bytes.NewReader(unconvertedKube))}
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

	convertedKoki := &bytes.Buffer{}

	err = client.WriteObjsToYamlStream(kokiObjs, convertedKoki)
	if err != nil {
		t.Errorf("path %s err %v", path, err)
		return err
	}

	// Extract the converted contents so we can output it if there's an error.
	convertedKokiString := convertedKoki.String()

	equal, err := cmp.CompareReader(bytes.NewBufferString(convertedKokiString), bytes.NewReader(unconvertedKoki))
	if err != nil {
		t.Errorf("path %s err %v\n%s\n\n%s", path, err, convertedKokiString, string(unconvertedKoki))
		return err
	}

	if !equal {
		t.Errorf("Failed to translate from Kube To Koki types. Resource Path=%s\n%s\n\n%s", path, convertedKokiString, string(unconvertedKoki))
		return nil
	}
	return nil
}

func testKokiToKube(path string, unconvertedKoki, unconvertedKube []byte, t *testing.T) error {
	// convert koki to kube
	streams := []io.ReadCloser{ioutil.NopCloser(bytes.NewReader(unconvertedKoki))}
	objs, err := parser.ParseStreams(streams)
	if err != nil {
		t.Errorf("path %s err %v", path, err)
		return err
	}

	kubeObjs, err := converter.ConvertToKubeNative(objs)
	if err != nil {
		t.Errorf("path %s err %v", path, err)
		return err
	}

	/*
		convertedKube := &bytes.Buffer{}

		err = client.WriteObjsToYamlStream(kubeObjs, convertedKube)
		if err != nil {
			t.Errorf("path %s err %v", path, err)
			return err
		}

		// Extract the converted contents so we can output it if there's an error.
		convertedKubeString := convertedKube.String()
	*/

	expectedStreams := []io.ReadCloser{ioutil.NopCloser(bytes.NewReader(unconvertedKube))}
	expectedObjs, err := parser.ParseStreams(expectedStreams)
	if err != nil {
		t.Errorf("error parsing expected kube objects")
		return err
	}
	if len(expectedObjs) != len(kubeObjs) {
		t.Errorf("different number of converted objects than expected: %d instead of %d", len(kubeObjs), len(expectedObjs))
		return nil
	}
	for i, obj := range expectedObjs {
		expectedKubeObj, err := parser.ParseSingleKubeNative(obj)
		if err != nil {
			t.Errorf("error parsing expected kube objects")
			return err
		}

		kubeObj := kubeObjs[i]
		if !reflect.DeepEqual(kubeObj, expectedKubeObj) {
			t.Errorf("Failed to translate from Koki To Kube types. Resource Path=%s\n%s", path, objectDiffString(kubeObj, expectedKubeObj))
			return nil
		}
	}

	return nil
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
