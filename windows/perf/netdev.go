/*
 * Author: fasion
 * Created time: 2019-08-20 17:27:26
 * Last Modified by: fasion
 * Last Modified time: 2019-08-21 10:49:54
 */

package perf

import (
	"fmt"
	"time"
	"github.com/fasionchan/osutil-go/windows/wmi"
)

var _ = fmt.Println

type NetDevMetrics struct {
	BitsReceived float64
	PacketsReceived float64
	ErrorsReceived float64

	BitsTransmitted float64
	PacketsTransmitted float64
	ErrorsTransmitted float64
}

type NetDevPerfSample struct {
	FetchTime time.Time
	Metric map[string]NetDevMetrics
}

type NetDevPerfSampler struct {
	lastSample *wmi.InterfacePerfRawSample
}

func NewNetDevPerfSampler() (*NetDevPerfSampler, error) {
	return &NetDevPerfSampler{

	}, nil
}

func (self *NetDevPerfSampler) Sample() (*NetDevPerfSample, error) {
	sample, err := wmi.SampleInterfacePerfRaw()
	if err != nil {
		return nil, err
	}

	if err = sample.FilterMisc(); err != nil {
		return nil, err
	}

	if err = sample.ToNetConnection(); err != nil {
		return nil, err
	}

	lastSample := self.lastSample
	self.lastSample = sample

	if lastSample == nil {
		return nil, nil
	}

	seconds := sample.FetchTime.Sub(lastSample.FetchTime).Seconds()

	name2LastPerf := make(map[string]*wmi.Win32_PerfRawData_Tcpip_NetworkInterface)
	for _, perf := range lastSample.Perfs {
		name2LastPerf[perf.Name] = perf
	}

	metric := make(map[string]NetDevMetrics)
	for _, perf := range sample.Perfs {
		lastPerf, ok := name2LastPerf[perf.Name]
		if !ok {
			continue
		}

		delta := perf.Sub(lastPerf)

		metric[perf.Name] = NetDevMetrics{
			BitsReceived: float64(delta.BytesReceivedPersec) / seconds * 8,
			BitsTransmitted: float64(delta.BytesSentPersec) / seconds * 8,
			PacketsReceived: float64(delta.PacketsReceivedPersec) / seconds * 8,
			PacketsTransmitted: float64(delta.PacketsSentPersec) / seconds * 8,
		}
	}

	return &NetDevPerfSample{
		FetchTime: sample.FetchTime,
		Metric: metric,
	}, nil
}
