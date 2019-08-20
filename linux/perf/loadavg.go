/*
 * Author: fasion
 * Created time: 2019-06-18 16:41:52
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:38:35
 */

package perf

import (
    "github.com/fasionchan/osutil-go/linux/procfs"
)

type LoadAvgSampler struct {
    procfsFetcher *procfs.ProcfsFetcher
}

func NewLoadAvgSampler(procfsFetcher *procfs.ProcfsFetcher) (*LoadAvgSampler, error) {
    return &LoadAvgSampler{
        procfsFetcher: procfsFetcher,
    }, nil
}

func (self *LoadAvgSampler) Sample() (*procfs.LoadAvgSample, error) {
    loadAvg, err := self.procfsFetcher.FetchLoadAvg()
    return loadAvg, err
}
