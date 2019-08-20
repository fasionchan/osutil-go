/*
 * Author: fasion
 * Created time: 2019-08-20 09:54:44
 * Last Modified by: fasion
 * Last Modified time: 2019-08-20 10:18:15
 */

package wmi

import (
	"time"

	"github.com/StackExchange/wmi"
)

type Win32_PerfRawData_PerfOS_Processor struct {
	Name string
	PercentIdleTime uint64
	PercentInterruptTime uint64
	PercentPrivilegedTime uint64
	PercentUserTime uint64
}

func (self *Win32_PerfRawData_PerfOS_Processor) Sub(other *Win32_PerfRawData_PerfOS_Processor) (*Win32_PerfRawData_PerfOS_Processor) {
	return &Win32_PerfRawData_PerfOS_Processor{
		PercentIdleTime: self.PercentIdleTime - other.PercentIdleTime,
		PercentInterruptTime: self.PercentInterruptTime - other.PercentInterruptTime,
		PercentPrivilegedTime: self.PercentPrivilegedTime - other.PercentPrivilegedTime,
		PercentUserTime: self.PercentUserTime - other.PercentUserTime,
	}
}

type ProcessorPerfRawSample struct {
	FetchTime time.Time
	Perfs map[string]*Win32_PerfRawData_PerfOS_Processor
}

func SampleProcessorPerfRaw() (*ProcessorPerfRawSample, error) {
	perfs := make([]*Win32_PerfRawData_PerfOS_Processor, 0)
	query := "SELECT * FROM Win32_PerfRawData_PerfOS_Processor"
	if err := wmi.QueryNamespace(query, &perfs, `\root\CIMV2`); err != nil {
		return nil, err
	}

	mapping := make(map[string]*Win32_PerfRawData_PerfOS_Processor)
	for _, perf := range perfs {
		mapping[perf.Name] = perf
	}

	return &ProcessorPerfRawSample{
		FetchTime: time.Now(),
		Perfs: mapping,
	}, nil
}
