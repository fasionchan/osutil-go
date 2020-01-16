/*
 * Author: fasion
 * Created time: 2019-06-28 10:02:33
 * Last Modified by: fasion
 * Last Modified time: 2019-12-04 16:40:15
 */

package perf

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"syscall"
	"time"

	"github.com/fasionchan/osutil-go/linux/procfs"
)

var _ = fmt.Println

type ProcessMetrics struct {
	CpuUtil       float64
	CpuUtilUser   float64
	CpuUtilSystem float64
	VirtualMemory float64
	RssMemory     float64
	OpenFiles     float64
}

type ProcessMetricSample struct {
	FetchTime time.Time
	LastTime  time.Time

	Pid     int
	Metrics ProcessMetrics
}

type ProcessMetricSampler struct {
	procfsFetcher *procfs.ProcfsFetcher

	lastPidStat *procfs.PidStat
}

func NewProcessMetricSampler(procfsFetcher *procfs.ProcfsFetcher) (*ProcessMetricSampler, error) {
	return &ProcessMetricSampler{
		procfsFetcher: procfsFetcher,
	}, nil
}

func (sampler *ProcessMetricSampler) DumpContext(w io.Writer) error {
	if err := json.NewEncoder(w).Encode(sampler.lastPidStat); err != nil {
		return err
	}

	return nil
}

func (sampler *ProcessMetricSampler) LoadContext(r io.Reader) error {
	lastPidStat := &procfs.PidStat{}
	if err := json.NewDecoder(r).Decode(lastPidStat); err != nil {
		return err
	}

	sampler.lastPidStat = lastPidStat

	return nil
}

func (sampler *ProcessMetricSampler) DumpContextToFile(path string) error {
	if f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644); err != nil {
		return err
	} else if err := sampler.DumpContext(f); err != nil {
		return err
	}

	return nil
}

func (sampler *ProcessMetricSampler) LoadContextFromFile(path string) error {
	if f, err := os.Open(path); err != nil {
		return err
	} else if err := sampler.LoadContext(f); err != nil {
		return err
	}

	return nil
}

func (self *ProcessMetricSampler) Sample(pid int) (*ProcessMetricSample, error) {
	pidAuxv, err := self.procfsFetcher.FetchPidAuxv(pid)
	if err != nil {
		return nil, err
	}

	clockTick, ok := pidAuxv.AT_CLKTCK()
	if !ok {
		clockTick = 100
	}

	pidStat, err := self.procfsFetcher.FetchPidStat(pid)
	if err != nil {
		return nil, err
	}

	lastPidStat := self.lastPidStat
	self.lastPidStat = pidStat

	if lastPidStat == nil {
		return nil, nil
	}
	if lastPidStat.Pid != pidStat.Pid {
		return nil, nil
	}

	pidStat = pidStat.Sub(lastPidStat)
	interval := (pidStat.FetchTime.Sub(lastPidStat.FetchTime)).Seconds()

	fds, err := procfs.CountPidFds(pid)
	if err != nil {
		return nil, err
	}

	sample := ProcessMetricSample{
		FetchTime: pidStat.FetchTime,
		LastTime:  lastPidStat.FetchTime,
		Pid:       pid,
		Metrics: ProcessMetrics{
			CpuUtil:       100 * float64(pidStat.Utime+pidStat.Stime) / float64(clockTick) / float64(interval),
			CpuUtilUser:   100 * float64(pidStat.Utime) / float64(clockTick) / float64(interval),
			CpuUtilSystem: 100 * float64(pidStat.Stime) / float64(clockTick) / float64(interval),
			VirtualMemory: float64(pidStat.Vsize),
			RssMemory:     float64(pidStat.Rss * uint64(syscall.Getpagesize())),
			OpenFiles:     float64(fds),
		},
	}

	return &sample, nil
}
