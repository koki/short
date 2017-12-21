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
	"github.com/koki/short/parser"
)

var temporarilyIgnoredResourceIDs = map[string]bool{
	"../testdata/pods/pod_spec_with_volume_name": true,
}

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

func TestControllerRevision(t *testing.T) {
	err := testResource("controller_revisions", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestCRDs(t *testing.T) {
	err := testResource("crds", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestEvents(t *testing.T) {
	err := testResource("events", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestInitializerConfiguration(t *testing.T) {
	err := testResource("initializer_config", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPodDisruptionPolicy(t *testing.T) {
	err := testResource("pod_disruption_policy", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPriorityClass(t *testing.T) {
	err := testResource("priority_class", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPodPreset(t *testing.T) {
	err := testResource("pod_preset", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestPodSecurityPolicy(t *testing.T) {
	err := testResource("pod_security_policy", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

func TestLimitRange(t *testing.T) {
	err := testResource("limit_range", testFuncGenerator(t))
	if err != nil {
		t.Fatal(err)
	}
}

type filePair struct {
	kubeSpec   string
	kokiSpec   string
	rekubeSpec string
}

func objectsEqual(a, b interface{}, aBytes, bBytes []byte) bool {
	if reflect.DeepEqual(a, b) {
		return true
	}

	diff := pretty.Diff(a, b)
	if len(diff) == 0 {
		return true
	}

	aString := strings.Trim(string(aBytes), "\n")
	bString := strings.Trim(string(bBytes), "\n")

	return aString == bString
}

func objectDiffString(a, b interface{}, aBytes, bBytes []byte) string {
	diff := pretty.Diff(a, b)
	if len(diff) > 0 {
		return strings.Join(diff, "\n")
	}

	return fmt.Sprintf("{\n%s\n\n%s\n}", string(aBytes), string(bBytes))
}

func testFuncGenerator(t *testing.T) func(string, filePair) error {
	return func(path string, fp filePair) error {
		kokiFile, err := os.Open(fp.kokiSpec)
		if err != nil {
			t.Errorf("failed to open koki file for %s at %s", path, fp.kokiSpec)
			return err
		}
		kubeFile, err := os.Open(fp.kubeSpec)
		if err != nil {
			t.Errorf("failed to open kube file for %s at %s", path, fp.kubeSpec)
			return err
		}

		unconvertedKoki, err := ioutil.ReadAll(kokiFile)
		if err != nil {
			t.Errorf("failed to read koki at %s", path)
			return err
		}
		unconvertedKube, err := ioutil.ReadAll(kubeFile)
		if err != nil {
			t.Errorf("failed to read kube at %s", path)
			return err
		}

		err = testKubeToKoki(path, unconvertedKube, unconvertedKoki, t)
		if err != nil {
			t.Errorf("failed kube -> koki")
			return err
		}

		// Some test cases don't expect to round-trip exactly.
		// Those tests have a .rekube.yaml file.
		if len(fp.rekubeSpec) > 0 {
			kubeFile, err := os.Open(fp.rekubeSpec)
			if err != nil {
				t.Errorf("failed to open rekube file for %s at %s", path, fp.rekubeSpec)
				return err
			}
			unconvertedKube, err = ioutil.ReadAll(kubeFile)
			if err != nil {
				t.Errorf("failed to read rekube at %s", path)
				return err
			}
		}
		err = testKokiToKube(path, unconvertedKoki, unconvertedKube, t)
		if err != nil {
			t.Errorf("failed koki -> kube")
			return err
		}

		return nil
	}
}

func testKubeToKoki(path string, unconvertedKube, expectedKokiBytes []byte, t *testing.T) error {
	expectedKokis, err := parseKokiBytes(expectedKokiBytes)
	if err != nil {
		t.Errorf("couldn't parse expected koki: path %s err %v", path, err)
	}

	// convert kube to koki
	streams := []io.ReadCloser{ioutil.NopCloser(bytes.NewReader(unconvertedKube))}
	objs, err := parser.ParseStreams(streams)
	if err != nil {
		t.Errorf("path %s err %v", path, err)
		return err
	}

	kokis, err := client.ConvertKubeMaps(objs)
	if err != nil {
		t.Errorf("path %s err %v", path, err)
		return err
	}

	convertedKokiBuf := &bytes.Buffer{}
	err = client.WriteObjsToYamlStream(kokis, convertedKokiBuf)
	if err != nil {
		t.Errorf("path %s err %v", path, err)
		return err
	}
	convertedKokiBytes := convertedKokiBuf.Bytes()

	if !objectsEqual(kokis, expectedKokis, convertedKokiBytes, expectedKokiBytes) {
		t.Errorf("Failed to translate from Kube To Koki types. Resource Path=%s\n%s", path,
			objectDiffString(kokis, expectedKokis, convertedKokiBytes, expectedKokiBytes))
		return fmt.Errorf("failed to translate")
	}
	return nil
}

// Parse koki objects.
func parseKokiBytes(b []byte) ([]interface{}, error) {
	streams := []io.ReadCloser{ioutil.NopCloser(bytes.NewReader(b))}
	objs, err := parser.ParseStreams(streams)
	if err != nil {
		return nil, err
	}

	kokis := make([]interface{}, len(objs))
	for i, obj := range objs {
		kokis[i], err = parser.ParseKokiNativeObject(obj)
		if err != nil {
			return nil, err
		}
	}

	return kokis, nil
}

// Reformat the kube-native yaml by round-tripping it.
func roundTripKubeBytes(b []byte) ([]interface{}, []byte, error) {
	streams := []io.ReadCloser{ioutil.NopCloser(bytes.NewReader(b))}
	objs, err := parser.ParseStreams(streams)
	if err != nil {
		return nil, nil, err
	}

	kubes := make([]interface{}, len(objs))
	for i, obj := range objs {
		kubes[i], err = parser.ParseSingleKubeNative(obj)
		if err != nil {
			return nil, nil, err
		}
	}

	kubesBuf := &bytes.Buffer{}
	err = client.WriteObjsToYamlStream(kubes, kubesBuf)
	if err != nil {
		return nil, nil, err
	}
	return kubes, kubesBuf.Bytes(), nil
}

func testKokiToKube(path string, unconvertedKoki, unconvertedKubeRaw []byte, t *testing.T) error {
	// reformat "expected" kube yaml
	expectedKubes, expectedKubeBytes, err := roundTripKubeBytes(unconvertedKubeRaw)
	if err != nil {
		t.Errorf("couldn't reformat expected kube: path %s err %v", path, err)
	}

	// convert koki to kube
	streams := []io.ReadCloser{ioutil.NopCloser(bytes.NewReader(unconvertedKoki))}
	objs, err := parser.ParseStreams(streams)
	if err != nil {
		t.Errorf("path %s err %v", path, err)
		return err
	}

	kubes, err := client.ConvertKokiMaps(objs)
	if err != nil {
		t.Errorf("path %s err %v", path, err)
		return err
	}

	convertedKubeBuf := &bytes.Buffer{}
	err = client.WriteObjsToYamlStream(kubes, convertedKubeBuf)
	if err != nil {
		t.Errorf("path %s err %v", path, err)
		return err
	}
	convertedKubeBytes := convertedKubeBuf.Bytes()

	if !objectsEqual(kubes, expectedKubes, convertedKubeBytes, expectedKubeBytes) {
		t.Errorf("Failed to translate from Koki To Kube types. Resource Path=%s\n%s", path,
			objectDiffString(kubes, expectedKubes, convertedKubeBytes, expectedKubeBytes))
		return fmt.Errorf("failed to translate")
	}

	return nil
}

func filePairsForResource(resource string) (map[string]filePair, error) {
	filePairs := map[string]filePair{}
	root := fmt.Sprintf("../testdata/%s", resource)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".short.yaml") {
			resourceID := strings.TrimSuffix(path, ".short.yaml")
			fp := filePairs[resourceID]
			fp.kokiSpec = path
			filePairs[resourceID] = fp
		} else if strings.HasSuffix(path, ".rekube.yaml") {
			resourceID := strings.TrimSuffix(path, ".rekube.yaml")
			fp := filePairs[resourceID]
			fp.rekubeSpec = path
			filePairs[resourceID] = fp
		} else if strings.HasSuffix(path, ".yaml") {
			resourceID := strings.TrimSuffix(path, ".yaml")
			fp := filePairs[resourceID]
			fp.kubeSpec = path
			filePairs[resourceID] = fp
		} else if path == root {
			return nil
		} else {
			return fmt.Errorf("Unrecognized file %s", path)
		}
		return nil
	})

	return filePairs, err
}

func testResource(resource string, test func(string, filePair) error) error {
	filePairs, err := filePairsForResource(resource)
	if err != nil {
		return err
	}

	failCount := 0
	var lastError error
	for resourceID, files := range filePairs {
		if _, ok := temporarilyIgnoredResourceIDs[resourceID]; ok {
			continue
		}

		err := test(resourceID, files)
		if err != nil {
			failCount++
			lastError = err
		}
	}

	if lastError != nil {
		return fmt.Errorf("\n\nerror #%d: %s\n\n", failCount, lastError.Error())
	}

	return nil
}
