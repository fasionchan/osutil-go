/*
 * Author: fasion
 * Created time: 2019-08-20 09:23:12
 * Last Modified by: fasion
 * Last Modified time: 2019-08-20 09:35:48
 */

package perf

import (
	"time"

	"github.com/StackExchange/wmi"
)

type LoadMetric struct {
	Load float64
}

type LoadSample struct {
	FetchTime time.Time
	Metric LoadMetric
}

type LoadSampler struct {

}

func NewLoadSampler() (*LoadSampler, error) {
	return &LoadSampler{
	}, nil
}

func (self *LoadSampler) Sample() (*LoadSample, error) {
	ps := make([]*Win32_Processor, 0)
	query := "SELECT * FROM Win32_Processor"
	if err := wmi.QueryNamespace(query, &ps, `\root\CIMV2`); err != nil {
		return nil, err
	}

	var total float64
	for _, p := range ps {
		total += float64(p.LoadPercentage)
	}

	return &LoadSample{
		FetchTime: time.Now(),
		Metric: LoadMetric{
			Load: total / float64(len(ps)),
		},
	}, nil
}

type Win32_Processor struct {
	LoadPercentage uint16
}
