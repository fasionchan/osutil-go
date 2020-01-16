// +build windows

/*
 * Author: fasion
 * Created time: 2019-08-13 17:47:08
 * Last Modified by: fasion
 * Last Modified time: 2019-12-26 16:59:42
 */

package sysinfo

import (
	"fmt"
	"net"
	"strings"

	"github.com/fasionchan/osutil-go/windows/wmi"
)

var _ = fmt.Println

func NetworkInterfaceCards() ([]NetworkInterfaceCard, error) {
	adapters, err := wmi.FetchNetworkAdapters()
	if err != nil {
		return nil, err
	}

	nics := make([]NetworkInterfaceCard, 0, len(adapters))

	for _, adapter := range adapters {
		var mac net.HardwareAddr
		address := strings.TrimSpace(string(adapter.MACAddress))
		if len(address) == 17 {
			mac, err = net.ParseMAC(address)
			if err != nil {
				return nil, err
			}
		}

		nics = append(nics, NetworkInterfaceCard{
			Index:        adapter.InterfaceIndex,
			Name:         adapter.NetConnectionID,
			MTU:          -1,
			Virtual:      adapter.Virtual(),
			HardwareAddr: mac,
		})
	}

	return nics, nil
}
