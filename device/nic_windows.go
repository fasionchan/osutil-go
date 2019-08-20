// +build windows

/*
 * Author: fasion
 * Created time: 2019-08-13 17:47:08
 * Last Modified by: fasion
 * Last Modified time: 2019-08-14 10:44:40
 */

package device

import (
	"net"
	"strings"
	"github.com/StackExchange/wmi"
)

type Win32_NetworkAdapter struct {
	Index int
	Name string
	InterfaceIndex int
	PhysicalAdapter bool
	PNPDeviceID string
	AdapterType string
	AdapterTypeId int
}

func NetworkInterfaceCards() ([]NetworkInterfaceCard, error) {
	adapters := make([]Win32_NetworkAdapter, 0)

	err := wmi.QueryNamespace("SELECT * FROM Win32_NetworkAdapter WHERE PhysicalAdapter=TRUE", &adapters, `\root\CIMV2`)
	if err != nil {
		return nil, err
	}

	index2adapter := make(map[int]Win32_NetworkAdapter)
	for _, adapter := range adapters {
		index2adapter[adapter.InterfaceIndex] = adapter
	}

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	nics := make([]NetworkInterfaceCard, 0, len(ifaces))

	for _, iface := range ifaces {
		nics = append(nics, NetworkInterfaceCard{
			Index: iface.Index,
			Name: iface.Name,
			MTU: iface.MTU,
			Virtual: !strings.HasPrefix(index2adapter[iface.Index].PNPDeviceID, `PCI\`),
			HardwareAddr: iface.HardwareAddr,
		})
	}

	return nics, nil
}
