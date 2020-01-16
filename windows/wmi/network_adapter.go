/*
 * Author: fasion
 * Created time: 2019-08-21 09:11:09
 * Last Modified by: fasion
 * Last Modified time: 2019-12-26 16:45:24
 */

package wmi

import (
	"strings"

	"github.com/StackExchange/wmi"
)

type Win32_NetworkAdapter struct {
	Index           int
	Name            string
	NetConnectionID string
	InterfaceIndex  int
	// PhysicalAdapter bool
	PNPDeviceID   string
	AdapterType   string
	AdapterTypeId int
	MACAddress    string
}

func (self *Win32_NetworkAdapter) Virtual() bool {
	return !strings.HasPrefix(self.PNPDeviceID, `PCI\`)
}

func FetchNetworkAdapters() ([]*Win32_NetworkAdapter, error) {
	adapters := make([]*Win32_NetworkAdapter, 0)

	err := wmi.QueryNamespace("SELECT * FROM Win32_NetworkAdapter", &adapters, `\root\CIMV2`)
	if err != nil {
		return nil, err
	}

	return adapters, nil
}
