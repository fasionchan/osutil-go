/*
 * Author: fasion
 * Created time: 2019-08-20 08:32:00
 * Last Modified by: fasion
 * Last Modified time: 2019-08-20 09:21:01
 */

package perf

import (
	"time"

	"github.com/StackExchange/wmi"

    unitutil "github.com/fasionchan/libgo/unit"
)

type MemoryMetric struct {
	Total unitutil.Bytes
	Free unitutil.Bytes
	Used unitutil.Bytes
}

type MemorySample struct {
	FetchTime time.Time
	Metric MemoryMetric
}

type MemorySampler struct {
}

func NewMemorySampler() (*MemorySampler, error) {
	return &MemorySampler{}, nil
}

func (self *MemorySampler) Sample() (*MemorySample, error) {
	var query string

	csSlice := make([]*Win32_ComputerSystem, 0)
	query = "SELECT * FROM Win32_ComputerSystem"
	if err := wmi.QueryNamespace(query, &csSlice, `\root\CIMV2`); err != nil {
		return nil, err
	}

	if len(csSlice) == 0 {
		return nil, nil
	}

	cs := csSlice[0]

	mpSlice := make([]*Win32_PerfRawData_PerfOS_Memory, 0)
	query = "SELECT * FROM Win32_PerfRawData_PerfOS_Memory"
	if err := wmi.QueryNamespace(query, &mpSlice, `\root\CIMV2`); err != nil {
		return nil, err
	}

	if len(mpSlice) == 0 {
		return nil, nil
	}

	mp := mpSlice[0]

	return &MemorySample{
		FetchTime: time.Now(),
		Metric: MemoryMetric{
			Total: unitutil.Bytes(cs.TotalPhysicalMemory),
			Free: unitutil.Bytes(mp.AvailableBytes),
			Used: unitutil.Bytes(cs.TotalPhysicalMemory-mp.AvailableBytes),
		},
	}, nil
}

type Win32_ComputerSystem struct {
	TotalPhysicalMemory uint64
}

type Win32_PerfRawData_PerfOS_Memory struct {
	AvailableBytes uint64
}
