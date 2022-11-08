package cmn

import (
	"bytes"
	"os/exec"
)

// 执行命令（Windows时为cmd，否则是bash）
func ExecCmd(command string) (stdout, stderr string, err error) {
	var out bytes.Buffer
	var errout bytes.Buffer

	cmd := exec.Command("/bin/bash", "-c", command)
	if IsWin() {
		cmd = exec.Command("cmd")
	}
	cmd.Stdout = &out
	cmd.Stderr = &errout
	err = cmd.Run()

	if err != nil {
		stderr = BytesToString(errout.Bytes())
	}
	stdout = BytesToString(out.Bytes())

	return
}
