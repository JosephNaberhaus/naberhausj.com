package image

import (
	"bytes"
	"fmt"
	"os/exec"
)

func runWithErrorChecking(cmd *exec.Cmd) ([]byte, error) {
	stdOut := new(bytes.Buffer)
	cmd.Stdout = stdOut
	stdErr := new(bytes.Buffer)
	cmd.Stderr = stdErr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("error when running \"%s\": %w\n%s", cmd.String(), err, stdErr.String())
	}

	return stdOut.Bytes(), nil
}
