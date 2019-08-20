/*
 * Author: fasion
 * Created time: 2019-08-16 14:57:13
 * Last Modified by: fasion
 * Last Modified time: 2019-08-16 19:58:03
 */

package perf

import (
	"fmt"
	"time"

	"github.com/StackExchange/wmi"
)

type Win32_PerfRawData_PerfProc_Process struct {
	IDProcess int

	Name string
	Caption string
	Description string

	PercentUserTime uint64
	PercentPrivilegedTime uint64
	PercentProcessorTime uint64

	VirtualBytes uint64
	WorkingSet uint64

	HandleCount uint64
}

func (self *Win32_PerfRawData_PerfProc_Process) Sub(other *Win32_PerfRawData_PerfProc_Process) (*Win32_PerfRawData_PerfProc_Process) {
	dup := *self

	dup.PercentPrivilegedTime -= other.PercentPrivilegedTime
	dup.PercentProcessorTime -= other.PercentProcessorTime
	dup.PercentUserTime -= other.PercentUserTime

	return &dup
}

type ProcessMetrics struct {
	CpuUtil float64
	CpuUtilUser float64
	CpuUtilSystem float64

	VirtualMemory float64
	RssMemory float64
	WorkingSetMemory float64

	Handles float64
}

type ProcessMetricSample struct {
	FetchTime time.Time
	LastTime time.Time

	Pid int
	Metrics ProcessMetrics
}

type ProcessMetricSampler struct {
	lastFetchTime time.Time
	lastPerf *Win32_PerfRawData_PerfProc_Process
}

func NewProcessMetricSampler() (*ProcessMetricSampler, error) {
	return &ProcessMetricSampler{
	}, nil
}

func (self *ProcessMetricSampler) Sample(pid int) (*ProcessMetricSample, error) {
	perfs := make([]*Win32_PerfRawData_PerfProc_Process, 0)
	query := fmt.Sprintf("SELECT * FROM Win32_PerfRawData_PerfProc_Process WHERE IDProcess=%d", pid)
	err := wmi.QueryNamespace(query, &perfs, `\root\CIMV2`)
	if err != nil {
		return nil, err
	}

	for _, perf := range perfs {
		if perf == nil {
			continue
		}

		if perf.IDProcess != pid {
			continue
		}

		fetchTime := time.Now()

		lastPerf := self.lastPerf
		lastFetchTime := self.lastFetchTime

		self.lastPerf = perf
		self.lastFetchTime = fetchTime

		interval := fetchTime.Sub(lastFetchTime).Seconds()

		if lastPerf == nil {
			continue
		}

		perf = perf.Sub(lastPerf)

		return &ProcessMetricSample{
			FetchTime: time.Now(),
			Pid: perf.IDProcess,
			Metrics: ProcessMetrics{
				CpuUtil: float64(perf.PercentProcessorTime) / 10000000 / interval * 100,
				CpuUtilUser: float64(perf.PercentUserTime) / 10000000 / interval * 100,
				CpuUtilSystem: float64(perf.PercentPrivilegedTime) / 10000000 / interval * 100,

				VirtualMemory: float64(perf.VirtualBytes),
				WorkingSetMemory: float64(perf.WorkingSet),
				RssMemory: float64(perf.WorkingSet),

				Handles: float64(perf.HandleCount),
			},
		}, nil
	}

	return nil, nil
}
