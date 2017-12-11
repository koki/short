package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kr/pretty"

	"github.com/koki/short/imports"
	"github.com/koki/short/parser"
	"github.com/koki/short/util"
	"github.com/koki/short/yaml"
)

func TestImports(t *testing.T) {
	tests, err := getImportsTests()
	if err != nil {
		t.Fatal(err)
	}

	for importPath, rootPath := range tests.RootMap {
		testImportsRoot(importPath, rootPath, tests.PathMap, t)
	}
}

type ImportsTests struct {
	// RootMap ./resource_id.yaml -> ../testdata/imports/resource_id.after.yaml
	RootMap map[string]string

	// PathMap ./resource_id.yaml -> ../testdata/imports/resource_id.yaml
	PathMap map[string]string
}

func importPathForResourceID(resourceID string) string {
	return fmt.Sprintf("./%s.yaml", resourceID)
}

func getImportsTests() (ImportsTests, error) {
	tests := ImportsTests{
		RootMap: map[string]string{},
		PathMap: map[string]string{},
	}
	root := "../testdata/imports"

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		_, filename := filepath.Split(path)
		if strings.HasSuffix(filename, ".after.yaml") {
			// Set RootResourceIDs based on *.after.yaml files.
			resourceID := strings.TrimSuffix(filename, ".after.yaml")
			tests.RootMap[importPathForResourceID(resourceID)] = path
		} else if strings.HasSuffix(filename, ".yaml") {
			// Set the PathMap for input files (*.yaml)
			resourceID := strings.TrimSuffix(filename, ".yaml")
			tests.PathMap[importPathForResourceID(resourceID)] = path
		} else if path == root {
			return nil
		} else {
			return fmt.Errorf("Unrecognized file %s", path)
		}
		return nil
	})

	return tests, err
}

func testImportsRoot(importPath, rootPath string, pathMap map[string]string, t *testing.T) {
	evalContext := imports.EvalContext{
		RawToTyped: parser.ParseKokiNativeObject,
		ResolveImportPath: func(rootPath, importPath string) (string, error) {
			return pathMap[importPath], nil
		},
		ReadFromPath: imports.ReadFromLocalPath,
	}

	modules, err := evalContext.Parse(pathMap[importPath])
	if err != nil {
		t.Errorf("failed to parse file at %s:\n%s", importPath, util.PrettyError(err))
		return
	}
	if len(modules) != 1 {
		t.Errorf("expected 1 resource, got %d at %s", len(modules), importPath)
		return
	}
	module := modules[0]
	err = evalContext.EvaluateModule(&module, nil)
	if err != nil {
		t.Errorf("failed to evaluate module at %s:\n%s", importPath, util.PrettyError(err))
		return
	}

	rootFile, err := ioutil.ReadFile(rootPath)
	if err != nil {
		t.Errorf("failed to read file at %s:\n%s", rootPath, util.PrettyError(err))
	}

	resource := module.Export.Raw
	resourceFile, err := yaml.Marshal(resource)
	if err != nil {
		t.Errorf("at %s, couldn't marshal %s, error: %s", importPath, pretty.Sprint(resource), util.PrettyError(err))
		return
	}

	expected := strings.Trim(string(rootFile), "\n")
	actual := strings.Trim(string(resourceFile), "\n")
	if expected != actual {
		t.Errorf("unexpected result at %s - actual, expected:\n%s\n\n%s", importPath, actual, expected)
		return
	}
}
