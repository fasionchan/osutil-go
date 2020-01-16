/*
 * Author: fasion
 * Created time: 2019-08-21 13:09:01
 * Last Modified by: fasion
 * Last Modified time: 2019-08-21 13:37:37
 */

package wmi

import (
	"time"
	"github.com/StackExchange/wmi"
)

type Win32_PerfRawData_PerfDisk_LogicalDisk struct {
	Name string

	DiskReadBytesPerSec uint64
	DiskWriteBytesPerSec uint64

	DiskReadsPerSec uint64
	DiskWritesPerSec uint64

	PercentDiskReadTime uint64
	PercentDiskWriteTime uint64

	PercentDiskTime uint64
	PercentIdleTime uint64

	CurrentDiskQueueLength uint64
}

func (self *Win32_PerfRawData_PerfDisk_LogicalDisk) IsZero() (bool) {
	return self.DiskReadBytesPerSec == 0 &&
		self.DiskWriteBytesPerSec == 0
}

func (self *Win32_PerfRawData_PerfDisk_LogicalDisk) Sub(other *Win32_PerfRawData_PerfDisk_LogicalDisk) (*Win32_PerfRawData_PerfDisk_LogicalDisk) {
	return &Win32_PerfRawData_PerfDisk_LogicalDisk{
		Name: self.Name,

		DiskReadBytesPerSec: self.DiskReadBytesPerSec - other.DiskReadBytesPerSec,
		DiskWriteBytesPerSec: self.DiskWriteBytesPerSec - other.DiskWriteBytesPerSec,

		DiskReadsPerSec: self.DiskReadsPerSec - other.DiskReadsPerSec,
		DiskWritesPerSec: self.DiskWritesPerSec - other.DiskWritesPerSec,

		PercentDiskReadTime: self.PercentDiskReadTime - other.PercentDiskReadTime,
		PercentDiskWriteTime: self.PercentDiskWriteTime - other.PercentDiskWriteTime,

		PercentDiskTime: self.PercentDiskTime - other.PercentDiskTime,
		PercentIdleTime: self.PercentIdleTime - other.PercentIdleTime,

		CurrentDiskQueueLength: self.CurrentDiskQueueLength,
	}
}

type LogicalDiskPerfRawSample struct {
	FetchTime time.Time
	Perfs []*Win32_PerfRawData_PerfDisk_LogicalDisk
}

func SampleLogicalDiskPerfRaw() (*LogicalDiskPerfRawSample, error) {
	perfs := make([]*Win32_PerfRawData_PerfDisk_LogicalDisk, 0)
	query := "SELECT * FROM Win32_PerfRawData_PerfDisk_LogicalDisk"
	if err := wmi.QueryNamespace(query, &perfs, `\root\CIMV2`); err != nil {
		return nil, err
	}

	return &LogicalDiskPerfRawSample{
		FetchTime: time.Now(),
		Perfs: perfs,
	}, nil
}

func (self *LogicalDiskPerfRawSample) PerfMapping() (map[string]*Win32_PerfRawData_PerfDisk_LogicalDisk) {
	perfs := make(map[string]*Win32_PerfRawData_PerfDisk_LogicalDisk)
	for _, perf := range self.Perfs {
		perfs[perf.Name] = perf
	}
	return perfs
}
