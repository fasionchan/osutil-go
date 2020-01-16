/*
 * Author: fasion
 * Created time: 2019-08-20 17:04:26
 * Last Modified by: fasion
 * Last Modified time: 2019-09-04 08:38:32
 */

package wmi

import (
	"fmt"
	"time"

	"github.com/StackExchange/wmi"
)

var _ = fmt.Println

type Win32_PerfRawData_Tcpip_NetworkInterface struct {
	Name string
	BytesReceivedPersec uint64
	BytesSentPersec uint64
	PacketsReceivedPersec uint64
	PacketsSentPersec uint64
}

func (self *Win32_PerfRawData_Tcpip_NetworkInterface) IsZero() (bool) {
	return self.BytesReceivedPersec == 0 &&
		self.BytesSentPersec == 0
}

func (self *Win32_PerfRawData_Tcpip_NetworkInterface) Sub(other *Win32_PerfRawData_Tcpip_NetworkInterface) (*Win32_PerfRawData_Tcpip_NetworkInterface) {
	return &Win32_PerfRawData_Tcpip_NetworkInterface{
		Name: self.Name,

		BytesReceivedPersec: self.BytesReceivedPersec - other.BytesReceivedPersec,
		BytesSentPersec: self.BytesSentPersec - other.BytesSentPersec,

		PacketsReceivedPersec: self.PacketsReceivedPersec - other.PacketsReceivedPersec,
		PacketsSentPersec: self.PacketsSentPersec - other.PacketsSentPersec,
	}
}

type InterfacePerfRawSample struct {
	FetchTime time.Time
	Perfs []*Win32_PerfRawData_Tcpip_NetworkInterface

	adapters []*Win32_NetworkAdapter
	name2adapter map[string]*Win32_NetworkAdapter
}

func SampleInterfacePerfRaw() (*InterfacePerfRawSample, error) {
	perfs := make([]*Win32_PerfRawData_Tcpip_NetworkInterface, 0)
	query := "SELECT * FROM Win32_PerfRawData_Tcpip_NetworkInterface"
	if err := wmi.QueryNamespace(query, &perfs, `\root\CIMV2`); err != nil {
		return nil, err
	}

	return &InterfacePerfRawSample{
		FetchTime: time.Now(),
		Perfs: perfs,
	}, nil
}

func (self *InterfacePerfRawSample) EnsureAdapters() (err error) {
	if self.adapters == nil {
		self.adapters, err = FetchNetworkAdapters()
		if err == nil {
			name2count := make(map[string]int)
			name2adapter := make(map[string]*Win32_NetworkAdapter)

			for _ , adapter := range self.adapters {
				count := name2count[adapter.Name] + 1
				name2count[adapter.Name] = count

				if count > 1 {
					name2adapter[fmt.Sprintf("%s _%d", adapter.Name, count)] = adapter
				} else {
					name2adapter[adapter.Name] = adapter
				}
			}

			self.name2adapter = name2adapter
		}
	}

	return
}

func (self *InterfacePerfRawSample) ToNetConnection() (error) {
	err := self.EnsureAdapters()
	if err != nil {
		return err
	}

	perfs := make([]*Win32_PerfRawData_Tcpip_NetworkInterface, 0, len(self.Perfs))
	for _, perf := range self.Perfs {
		adapter, ok := self.name2adapter[perf.Name]
		if !ok {
			continue
		}

		if adapter.NetConnectionID == "" {
			continue
		}

		perf.Name = adapter.NetConnectionID

		perfs = append(perfs, perf)
	}

	self.Perfs = perfs

	return nil
}

func (self *InterfacePerfRawSample) FilterMisc() (error) {
	err := self.EnsureAdapters()
	if err != nil {
		return err
	}

	perfs := make([]*Win32_PerfRawData_Tcpip_NetworkInterface, 0, len(self.Perfs))
	for _, perf := range self.Perfs {
		if perf.IsZero() {
			continue
		}

		adapter, ok := self.name2adapter[perf.Name]
		if !ok {
			continue
		}

		if adapter.Virtual() {
			continue
		}

		perfs = append(perfs, perf)
	}

	self.Perfs = perfs

	return nil
}
