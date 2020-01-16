// +build linux

/*
 * Author: fasion
 * Created time: 2019-05-17 17:44:50
 * Last Modified by: fasion
 * Last Modified time: 2019-12-19 12:26:16
 */

package sysinfo

import (
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	NicPath       = "/sys/class/net"
	VirtualPrefix = "/sys/devices/virtual"
)

func IsVirtualNicByPath(path string) (bool, error) {
	// has virtual prefix, then it is virtual device
	if strings.HasPrefix(path, VirtualPrefix) {
		return true, nil
	}

	// device path
	devicePath, err := filepath.Abs(
		filepath.Join(path, "device"),
	)
	if err != nil {
		return false, err
	}

	// check device path
	if _, err := os.Stat(devicePath); err != nil {
		// if device path not exists, it is virtual device
		if os.IsNotExist(err) {
			return true, nil
		}

		return false, err
	}

	return false, nil
}

func IsVirtualNic(nic string) (bool, error) {
	path, err := filepath.Abs(
		filepath.Join(NicPath, nic),
	)
	if err != nil {
		return false, err
	}

	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false, err
	}

	return IsVirtualNicByPath(realPath)
}

func NetworkInterfaceCards() ([]NetworkInterfaceCard, error) {
	files, err := ioutil.ReadDir(NicPath)
	if err != nil {
		return nil, err
	}

	nics := make([]NetworkInterfaceCard, 0, len(files))

	for _, file := range files {
		path, err := filepath.Abs(
			filepath.Join(NicPath, file.Name()),
		)
		if err != nil {
			return nil, err
		}

		realPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			return nil, err
		}

		name := file.Name()
		virtual, err := IsVirtualNicByPath(realPath)
		if err != nil {
			return nil, err
		}

		rawAddress, err := ioutil.ReadFile(filepath.Join(realPath, "address"))
		if err != nil {
			return nil, err
		}

		var mac net.HardwareAddr
		address := strings.TrimSpace(string(rawAddress))
		if len(address) == 17 {
			mac, err = net.ParseMAC(address)
			if err != nil {
				return nil, err
			}
		}

		nics = append(nics, NetworkInterfaceCard{
			Name:         name,
			Virtual:      virtual,
			HardwareAddr: mac,
		})
	}

	return nics, nil
}
