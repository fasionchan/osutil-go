/*
 * Author: fasion
 * Created time: 2019-06-18 16:41:52
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:38:20
 */

package perf

import (
    "syscall"
    "time"

    "github.com/fasionchan/osutil-go/linux/procfs"
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

func FetchFileSystemUsage(path string) (*FileSystemUsage, error) {
    var stat syscall.Statfs_t

    err := syscall.Statfs(path, &stat)
    if err != nil {
        return nil, err
    }

    bsize := uint64(stat.Bsize)

    return &FileSystemUsage{
        TotalBytes:         unitUtil.Bytes(bsize * stat.Blocks),
        FreeBytes:          unitUtil.Bytes(bsize * stat.Bfree),
        AvailableBytes:     unitUtil.Bytes(bsize * stat.Bavail),
        UsedBytes:          unitUtil.Bytes(bsize * (stat.Blocks - stat.Bfree)),
        UsedBytesPercent:   float64(stat.Blocks-stat.Bfree) / float64(stat.Blocks-stat.Bfree+stat.Bavail) * 100,

        TotalFiles:         stat.Files,
        FreeFiles:          stat.Ffree,
        AvailableFiles:     stat.Ffree,
        UsedFiles:          stat.Files - stat.Ffree,
        UsedFilesPercent:   float64(stat.Files - stat.Ffree) / float64(stat.Files) * 100,
    }, nil
}

type FileSystemUsageSample struct {
    FetchTime   time.Time
    Usages       map[string]FileSystemUsage
}

type FileSystemUsageSampler struct {
    procfsFetcher *procfs.ProcfsFetcher
}

func (self *FileSystemUsageSampler) Sample() (*FileSystemUsageSample, error) {
	fetchTime := time.Now()

    sample := FileSystemUsageSample{
        FetchTime:  fetchTime,
        Usages:      make(map[string]FileSystemUsage),
    }

    infos, err := self.procfsFetcher.FetchMountInfos()
    if err != nil {
        return nil, err
    }

    for _, info := range procfs.FilterMountInfos(infos.Infos) {
        path := info.MountPoint

        stat, err := FetchFileSystemUsage(path)
        if err != nil {
            return nil, err
        }

        sample.Usages[path] = *stat
    }

    return &sample, nil
}

func NewFileSystemUsageSampler(fetcher *procfs.ProcfsFetcher) (*FileSystemUsageSampler, error) {
    return &FileSystemUsageSampler{
        procfsFetcher: fetcher,
    }, nil
}
