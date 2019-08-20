// +build linux

/*
 * Author: fasion
 * Created time: 2019-05-17 17:44:50
 * Last Modified by: fasion
 * Last Modified time: 2019-08-14 10:53:27
 */

package device

import (
    "io/ioutil"
    "net"
	"path/filepath"
	"strings"
)

const (
	NicPath = "/sys/class/net"
	VirtualPrefix = "/sys/devices/virtual"
)

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

	return strings.HasPrefix(realPath, VirtualPrefix), nil
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
		virtual := strings.HasPrefix(realPath, VirtualPrefix)

		address, err := ioutil.ReadFile(filepath.Join(realPath, "address"))
		if err != nil {
			return nil, err
		}

		mac, err := net.ParseMAC(strings.TrimSpace(string(address)))
		if err != nil {
			return nil, err
		}

		nics = append(nics, NetworkInterfaceCard{
			Name: 			name,
			Virtual:		virtual,
			HardwareAddr: 	mac,
		})
	}

	return nics, nil
}
