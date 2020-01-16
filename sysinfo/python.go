/*
 * Author: fasion
 * Created time: 2019-09-05 13:45:36
 * Last Modified by: fasion
 * Last Modified time: 2019-12-16 16:13:05
 */

package sysinfo

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

var _ = fmt.Println

const (
	ShowPythonEnvCode = `
import json
import platform
import sys

def main():
	print(json.dumps({
		'Installed': True,
		'Executable': sys.executable,
		'Version': platform.python_version(),
	}))

if __name__ == '__main__':
	main()
	`
)

func FetchPythonInfo(python string) (*PythonEnvironment, error) {
	output, err := exec.Command(python, "-c", ShowPythonEnvCode).CombinedOutput()
	if err != nil {
		return nil, err
	}

	var env PythonEnvironment
	if err = json.Unmarshal(output, &env); err != nil {
		return nil, err
	}

	env.Installed = true

	return &env, nil
}
