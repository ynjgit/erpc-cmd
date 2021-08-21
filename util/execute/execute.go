package execute

import (
	"fmt"
	"os/exec"
)

// HasCmd ...
func HasCmd(cmd string) error {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return err
	}

	return nil
}

// RunCmd ...
func RunCmd(cmd string, args ...string) (output string, err error) {
	fmt.Println(cmd, args)
	c := exec.Command(cmd, args...)
	out, err := c.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
