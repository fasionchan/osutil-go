/*
 * Author: fasion
 * Created time: 2019-08-01 13:18:05
 * Last Modified by: fasion
 * Last Modified time: 2019-08-20 11:25:49
 */

package perf

import (
	"fmt"
    "time"
    "github.com/fasionchan/osutil-go/linux/procfs"
)

var _ = fmt.Println

type NetDevMetrics struct {
	BitsReceived float64
	PacketsReceived float64
	ErrorsReceived float64
	DropsReceived float64

	BitsTransmitted float64
	PacketsTransmitted float64
	ErrorsTransmitted float64
	DropsTransmitted float64
}

type NetDevPerfSample struct {
	FetchTime time.Time
	Perfs map[string]NetDevMetrics
}


type NetDevPerfSampler struct {
	procfsFetcher *procfs.ProcfsFetcher

	lastCounters *procfs.NetDevCounters
}

func NewNetDevPerfSampler(fetcher *procfs.ProcfsFetcher) (*NetDevPerfSampler, error) {
	return &NetDevPerfSampler{
		procfsFetcher: fetcher,
	}, nil
}

func (self *NetDevPerfSampler) Sample() (*NetDevPerfSample, error) {
	counters, err := self.procfsFetcher.FetchNetDevCounters()
	if err != nil {
		return nil, err
	}

	counters.FilterMisc()

	lastCounters := self.lastCounters
	self.lastCounters = counters

	if lastCounters == nil {
		return nil, nil
	}

	seconds := counters.FetchTime.Sub(lastCounters.FetchTime).Seconds()

	name2LastCounter := make(map[string]procfs.NetDevCounter)
	for _, lastCounter := range lastCounters.Records {
		name2LastCounter[lastCounter.Name] = lastCounter
	}

	perfs := make(map[string]NetDevMetrics)

	for _, counter := range counters.Records {
		lastCounter, ok := name2LastCounter[counter.Name]
		if !ok {
			continue
		}

		rDelta := counter.Receive.Sub(lastCounter.Receive)
		tDelta := counter.Transmit.Sub(lastCounter.Transmit)

		perfs[counter.Name] = NetDevMetrics{
			BitsReceived: float64(rDelta.Bytes.ToBit()) / seconds,
			BitsTransmitted: float64(tDelta.Bytes.ToBit()) / seconds,
		}
	}

	return &NetDevPerfSample{
		FetchTime: counters.FetchTime,
		Perfs: perfs,
	}, nil
}
