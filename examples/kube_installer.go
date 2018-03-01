package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
)

type kubeInstaller struct{}

var InstallerPlugin kubeInstaller

func (k *kubeInstaller) Install(input interface{}) error {
	buf, ok := input.(*bytes.Buffer)
	if !ok {
		return fmt.Errorf("input data is not of type bytes.Buffer")
	}

	cmd := exec.Command("kubectl", "create", "-f", "-")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		io.Copy(stdin, buf)
		stdin.Close()
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v\n%s", err, out)
	}

	fmt.Printf("%s", out)

	return nil
}
