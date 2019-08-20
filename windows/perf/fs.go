/*
 * Author: fasion
 * Created time: 2019-08-16 09:55:20
 * Last Modified by: fasion
 * Last Modified time: 2019-08-16 13:55:12
 */

package perf

import (
	"time"

	"github.com/StackExchange/wmi"

    unitUtil "github.com/fasionchan/libgo/unit"
)

type FileSystemUsage struct {
    TotalBytes          unitUtil.Bytes
    FreeBytes           unitUtil.Bytes
    AvailableBytes      unitUtil.Bytes
    UsedBytes           unitUtil.Bytes
    UsedBytesPercent    float64

    TotalFiles          uint64
    FreeFiles           uint64
    AvailableFiles      uint64
    UsedFiles           uint64
    UsedFilesPercent    float64
}

type FileSystemUsageSample struct {
    FetchTime   time.Time
    Usages       map[string]FileSystemUsage
}

type FileSystemUsageSampler struct {
}

func NewFileSystemUsageSampler() (*FileSystemUsageSampler, error) {
	return &FileSystemUsageSampler{
	}, nil
}

type Win32_LogicalDisk struct {
	Name string
	DriveType uint64
	Size uint64
	FreeSpace uint64
}

func (self *FileSystemUsageSampler) Sample() (*FileSystemUsageSample, error) {
	disks := make([]Win32_LogicalDisk, 0)
	err := wmi.QueryNamespace("SELECT * FROM Win32_LogicalDisk", &disks, `\root\CIMV2`)
	if err != nil {
		return nil, err
	}

	sample := FileSystemUsageSample{
		FetchTime: time.Now(),
		Usages: make(map[string]FileSystemUsage),
	}

	for _, disk := range disks {
		switch disk.DriveType {
		// local disk
		case 3:
		default:
			continue
		}

		if disk.Size == 0 {
			continue
		}

		usedBytes := disk.Size - disk.FreeSpace

		sample.Usages[disk.Name] = FileSystemUsage{
			TotalBytes: unitUtil.Bytes(disk.Size),
			FreeBytes: unitUtil.Bytes(disk.FreeSpace),
			AvailableBytes: unitUtil.Bytes(disk.FreeSpace),
			UsedBytes: unitUtil.Bytes(usedBytes),
			UsedBytesPercent: float64(usedBytes) / float64(disk.Size) * 100,
		}
	}

	return &sample, nil
}
