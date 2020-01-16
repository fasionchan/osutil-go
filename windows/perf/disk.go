/*
 * Author: fasion
 * Created time: 2019-08-21 13:29:54
 * Last Modified by: fasion
 * Last Modified time: 2019-08-21 15:28:39
 */

package perf

import (
	"strings"
	"time"

	"github.com/fasionchan/osutil-go/windows/wmi"
)

type DiskMetric struct {
	Reads float64
	BytesRead float64

	Writes float64
	BytesWritten float64

	IoUtil float64
	QueueSize float64
}

type DiskPerfSample struct {
	FetchTime time.Time
	Metric map[string]DiskMetric
}

type DiskPerfSampler struct {
	lastSample *wmi.LogicalDiskPerfRawSample
}

func NewDiskPerfSampler() (*DiskPerfSampler, error) {
	return &DiskPerfSampler{
	}, nil
}

func (self *DiskPerfSampler) Sample() (*DiskPerfSample, error) {
	sample, err := wmi.SampleLogicalDiskPerfRaw()
	if err != nil {
		return nil, err
	}

	if len(sample.Perfs) == 0 {
		return nil, nil
	}

	lastSample := self.lastSample
	self.lastSample = sample

	if lastSample == nil {
		return nil, nil
	}

	seconds := sample.FetchTime.Sub(lastSample.FetchTime).Seconds()
	name2LastPerf := lastSample.PerfMapping()
	result := make(map[string]DiskMetric)

	for _, perf := range sample.Perfs {
		if perf.Name == "_Total" {
			continue
		}

		if !strings.HasSuffix(perf.Name, ":") {
			continue
		}

		lastPerf, ok := name2LastPerf[perf.Name]
		if !ok {
			continue
		}

		delta := perf.Sub(lastPerf)

		result[perf.Name] = DiskMetric{
			BytesRead: float64(delta.DiskReadBytesPerSec) / seconds,
			BytesWritten: float64(delta.DiskWriteBytesPerSec) / seconds,

			Reads: float64(delta.DiskReadsPerSec) / seconds,
			Writes: float64(delta.DiskWritesPerSec) / seconds,

			IoUtil: float64(delta.PercentDiskTime) / float64(delta.PercentDiskTime + delta.PercentIdleTime) * 100,
			QueueSize: float64(perf.CurrentDiskQueueLength),
		}
	}

	return &DiskPerfSample{
		FetchTime: sample.FetchTime,
		Metric: result,
	}, nil
}
