/*
 * Author: fasion
 * Created time: 2019-09-05 13:43:12
 * Last Modified by: fasion
 * Last Modified time: 2019-09-05 14:21:56
 */

package sysinfo

func FetchPythonEnvironment() (env *PythonEnvironment, err error) {
	for _, python := range []string{
		"python3",
		"python",
	} {
		if env, err = FetchPythonInfo(python); err == nil {
			return
		}
	}

	return
}
