/*
 * Author: fasion
 * Created time: 2019-06-18 16:41:52
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:41:21
 */

package perf

import (
    //"encoding/json"
    "fmt"
    "time"
    "github.com/fasionchan/osutil-go/linux/procfs"
    mathutil "github.com/fasionchan/libgo/math"
)

type SoftirqStat struct {
    Cpu []float64
    Irq map[string]float64
}

type SoftirqStatSample struct {
	FetchTime time.Time
	LastTime time.Time
	Stat SoftirqStat
}

type SoftirqStatSampler struct {
    fetcher *procfs.ProcfsFetcher
    lastTime time.Time
    lastCounter *procfs.SoftirqCounter
}

func (self *SoftirqStatSampler) Sample() (*SoftirqStatSample, error) {
    counter, err := self.fetcher.FetchSoftirqCounters()
    if err != nil {
        return nil, err
    }

    lastCounter := self.lastCounter
    lastTime := lastCounter.FetchTime

	self.lastCounter = counter
	fetchTime := counter.FetchTime

    if lastCounter == nil {
        return nil, nil
    }

    interval := fetchTime.Sub(lastTime).Seconds()
    delta := counter.Sub(lastCounter)

    stat := SoftirqStat{
        Irq: make(map[string]float64),
        Cpu: make([]float64, delta.Cores),
    }

    for i, name := range delta.Names {
        stat.Irq[name] = float64(mathutil.SumUint64Slice(delta.CounterArray[i])) / interval

        for j, value := range delta.CounterArray[i] {
            stat.Cpu[j] += float64(value)
        }

        for j:=0; j<delta.Cores; j++ {
            stat.Cpu[j] /= interval
        }
    }

	return &SoftirqStatSample{
		FetchTime: fetchTime,
		LastTime: lastTime,
		Stat: stat,
	}, nil
}

func NewSoftirqSampler(fetcher *procfs.ProcfsFetcher) (*SoftirqStatSampler, error) {
    if false {
        fmt.Println()
    }

    var err error
    if fetcher == nil {
        fetcher, err = procfs.NewProcfsFetcher()
        if err != nil {
            return nil, err
        }
    }

    sampler := SoftirqStatSampler{
        fetcher: fetcher,
        lastCounter: nil,
    }

    return &sampler, nil
}
