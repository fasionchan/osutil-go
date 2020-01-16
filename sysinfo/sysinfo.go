/*
 * Author: fasion
 * Created time: 2019-09-02 16:00:27
 * Last Modified by: fasion
 * Last Modified time: 2019-09-05 14:57:48
 */

package sysinfo

import (
	"fmt"
	"strings"
)

var _ = fmt.Println

func FetchSysInfo() (*SysInfo, error) {
	system, err := FetchSystem()
	if err != nil {
		return nil, err
	}

	bios, err := FetchBios()
	if err != nil {
		return nil, err
	}

	os, err := FetchOs()
	if err != nil {
		return nil, err
	}

	scsi, err := FetchScsiDevices()
	if err != nil {
		return nil, err
	}

	sysinfo := SysInfo{
		System: *system,
		Bios: *bios,
		Os: *os,
		ScsiDevices: scsi,
	}

	python2, err := FetchPythonInfo("python2")
	if err == nil {
		sysinfo.Python2 = *python2
	}

	python3, err := FetchPythonInfo("python3")
	if err == nil {
		sysinfo.Python3 = *python3
	}

	if !sysinfo.Python2.Installed || !sysinfo.Python3.Installed {
		python, err := FetchPythonInfo("python")
		if err == nil {
			if strings.HasPrefix(python.Version, "3.") {
				if !sysinfo.Python3.Installed {
					sysinfo.Python3 = *python
				}
			} else if strings.HasPrefix(python.Version, "2.") {
				if !sysinfo.Python2.Installed {
					sysinfo.Python2 = *python
				}
			}
		}
	}

	return &sysinfo, nil
}
