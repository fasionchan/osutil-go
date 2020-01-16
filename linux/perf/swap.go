/*
 * Author: fasion
 * Created time: 2019-09-12 08:57:35
 * Last Modified by: fasion
 * Last Modified time: 2019-09-12 09:24:34
 */

package perf

import (
    "time"
    "github.com/fasionchan/osutil-go/linux/procfs"
    unitutil "github.com/fasionchan/libgo/unit"
)

type SwapInfo struct {
    Total       unitutil.Bytes
    Used        unitutil.Bytes
	Free        unitutil.Bytes
}

type SwapInfoSample struct {
    FetchTime   time.Time
    Info        SwapInfo
}

type SwapInfoSampler struct {
    procfsFetcher *procfs.ProcfsFetcher
}

func NewSwapInfoSampler(fetcher *procfs.ProcfsFetcher) (*SwapInfoSampler, error) {
    return &SwapInfoSampler{
        procfsFetcher: fetcher,
    }, nil
}

func (self *SwapInfoSampler) Sample() (*SwapInfoSample, error) {
	info, err := self.procfsFetcher.FetchMemInfo()
	if err != nil {
		return nil, err
	}

    sample := SwapInfoSample{
        FetchTime: info.FetchTime,
        Info: SwapInfo{
            Total: info.Info["SwapTotal"],
            Used: info.Info["SwapTotal"] - info.Info["SwapFree"],
			Free: info.Info["SwapFree"],
		},
	}

	return &sample, nil
}
