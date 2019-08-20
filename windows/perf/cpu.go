/*
 * Author: fasion
 * Created time: 2019-08-20 09:39:14
 * Last Modified by: fasion
 * Last Modified time: 2019-08-20 14:20:50
 */

package perf

import (
	"time"
	"github.com/fasionchan/osutil-go/windows/wmi"
)

type CpuUtilMetric struct {
	Idle float64
	Interrupt float64
	Privileged float64
	User float64
}

type CpuUtilMetricSet struct {
	Whole CpuUtilMetric
	//PerCores CpuUtilMetric
}

type CpuUtilSample struct {
	FetchTime time.Time
	MetricSet CpuUtilMetricSet
}

type CpuUtilSampler struct {
	lastSample *wmi.ProcessorPerfRawSample
}

func NewCpuUtilSampler() (*CpuUtilSampler, error) {
	return &CpuUtilSampler{

	}, nil
}

func (self *CpuUtilSampler) Sample() (*CpuUtilSample, error) {
	sample, err := wmi.SampleProcessorPerfRaw()
	if err != nil {
		return nil, err
	}

	lastSample := self.lastSample
	self.lastSample = sample

	if lastSample == nil {
		return nil, nil
	}

	perf, ok := sample.Perfs["_Total"]
	if !ok {
		return nil, nil
	}

	lastPerf, ok := lastSample.Perfs["_Total"]
	if !ok {
		return nil, nil
	}

	delta := perf.Sub(lastPerf)
	total := delta.PercentUserTime + delta.PercentPrivilegedTime + delta.PercentInterruptTime + delta.PercentIdleTime

	return &CpuUtilSample{
		FetchTime: sample.FetchTime,
		MetricSet: CpuUtilMetricSet{
			Whole: CpuUtilMetric{
				User: float64(delta.PercentUserTime) / float64(total) * 100,
				Privileged: float64(delta.PercentPrivilegedTime) / float64(total) * 100,
				Interrupt: float64(delta.PercentInterruptTime) / float64(total) * 100,
				Idle: float64(delta.PercentIdleTime) / float64(total) * 100,
			},
		},
	}, nil
}
