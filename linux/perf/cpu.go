/*
 * Author: fasion
 * Created time: 2019-06-18 15:56:49
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:36:39
 */

package perf

import (
    "fmt"
    "reflect"
    "time"
    "github.com/fasionchan/osutil-go/linux/procfs"
    mathutil "github.com/fasionchan/libgo/math"
)

var _ = fmt.Println

type CpuUsageMetrics struct {
    User float64
    Nice float64
    System float64
    Idle float64
    Iowait float64
    Irq float64
    Softirq float64
    Steal float64
    Guest float64
    GuestNice float64
}

func CalculateCpuUsage(counters []uint64, usage *CpuUsageMetrics) (*CpuUsageMetrics, error) {
    total := mathutil.SumUint64Slice(counters[:8])
    if total == 0 {
        total = 1
    }

    ref := reflect.ValueOf(usage)
    for i, value := range counters {
        ref.Elem().Field(i).Set(reflect.ValueOf(
            100 * float64(value) / float64(total),
        ))
    }

    return usage, nil
}

func NewCpuUsageMetrics(counters []uint64) (*CpuUsageMetrics, error) {
    usage := CpuUsageMetrics{}
    return CalculateCpuUsage(counters, &usage)
}

type CpuUsage struct {
    Usage CpuUsageMetrics
    PerCoreUsages []CpuUsageMetrics
}

type CpuUsageSample struct {
	FetchTime time.Time
	LastTime time.Time
	Usage CpuUsage
}

type CpuUsageSampler struct {
    procfsFetcher *procfs.ProcfsFetcher

    lastStatCounter *procfs.StatCounter
}

func NewCpuUsageSampler(procfsFetcher *procfs.ProcfsFetcher) (*CpuUsageSampler, error) {
    return &CpuUsageSampler{
        procfsFetcher: procfsFetcher,
    }, nil
}

func (self *CpuUsageSampler) Sample() (*CpuUsageSample, error) {
    statCounter, err := self.procfsFetcher.FetchStatCounter()
    if err != nil {
        return nil, err
    }

    lastStatCounter := self.lastStatCounter
    self.lastStatCounter = statCounter

    if lastStatCounter == nil {
        return nil, nil
    }

    delta := statCounter.Sub(lastStatCounter)

    usage := CpuUsage{
        PerCoreUsages: make([]CpuUsageMetrics, len(delta.Cpus)),
    }

    _, err = CalculateCpuUsage(delta.Cpu, &usage.Usage)
    if err != nil {
        return nil, err
    }

    for i, cpu := range delta.Cpus {
        _, err = CalculateCpuUsage(cpu, &usage.PerCoreUsages[i])
        if err != nil {
            return nil, nil
        }
    }

    return &CpuUsageSample{
		FetchTime: statCounter.FetchTime,
		LastTime: lastStatCounter.FetchTime,
		Usage: usage,
	}, nil
}
