/*
 * Author: fasion
 * Created time: 2019-06-18 16:41:52
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:38:51
 */

package perf

import (
    "time"
    "github.com/fasionchan/osutil-go/linux/procfs"
    unitutil "github.com/fasionchan/libgo/unit"
)

type MemInfo struct {
    Total       unitutil.Bytes
    Used        unitutil.Bytes
    Free        unitutil.Bytes
    Allocatable unitutil.Bytes
    Buffers     unitutil.Bytes
    Cached      unitutil.Bytes
    Slab        unitutil.Bytes
    Cache       unitutil.Bytes
    Active      unitutil.Bytes
    Inactive    unitutil.Bytes
    Available   unitutil.Bytes
}

type MemInfoSample struct {
    FetchTime   time.Time
    Info        MemInfo
}

type MemInfoSampler struct {
    procfsFetcher *procfs.ProcfsFetcher
}

func (self *MemInfoSampler) Sample() (*MemInfoSample, error) {
	info, err := self.procfsFetcher.FetchMemInfo()
	if err != nil {
		return nil, err
	}

	cache := info.Info["Cached"]
	if slab, ok := info.Info["SReclaimable"]; ok {
		cache += slab
	} else {
		cache += info.Info["Slab"]
	}

    allocatable := info.Info["MemFree"] + info.Info["Buffers"] + cache
    used := info.Info["MemTotal"] - allocatable

    sample := MemInfoSample{
        FetchTime: info.FetchTime,
        Info: MemInfo{
            Total: info.Info["MemTotal"],
            Used: used,
            Free: info.Info["MemFree"],
            Allocatable: allocatable,
            Buffers: info.Info["Buffers"],
            Cached: info.Info["Cached"],
            Slab: info.Info["Slab"],
            Cache: cache,
            Active: info.Info["Active"],
            Inactive: info.Info["Inactive"],
            Available: info.Info["MemAvailable"],
        },
    }

    return &sample, err
}

func NewMemInfoSampler(fetcher *procfs.ProcfsFetcher) (*MemInfoSampler, error) {
    return &MemInfoSampler{
        procfsFetcher: fetcher,
    }, nil
}
