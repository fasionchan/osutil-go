/*
 * Author: fasion
 * Created time: 2019-06-18 16:41:52
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:41:44
 */

package perf

import (
    "fmt"
    "time"
    "github.com/fasionchan/osutil-go/linux/procfs"
)

var _ = fmt.Println

type TaskStat struct {
    ContextSwitches float64
    Forks float64
    ProcsRunning uint64
    ProcsBlocked uint64
}

type TaskStatSample struct {
	FetchTime time.Time
	LastTime time.Time
    Stat TaskStat
}

type TaskStatSampler struct {
    procfsFetcher *procfs.ProcfsFetcher

    lastStat *procfs.StatCounter
}

func (self *TaskStatSampler) Sample() (*TaskStatSample, error) {
    stat, err := self.procfsFetcher.FetchStatCounter()
    if err != nil {
        return nil, err
    }

    lastStat := self.lastStat
    self.lastStat = stat

    if lastStat == nil {
        return nil, err
    }

    interval := float64(stat.FetchTime.Sub(lastStat.FetchTime)) / float64(time.Second)

    taskStat := TaskStat{
        ContextSwitches: float64(stat.Ctxt-lastStat.Ctxt) / interval,
        Forks: float64(stat.Processes-lastStat.Processes) / interval,
        ProcsRunning: stat.ProcsRunning,
        ProcsBlocked: stat.ProcsBlocked,
    }

	return &TaskStatSample{
		FetchTime: stat.FetchTime,
		LastTime: lastStat.FetchTime,
		Stat: taskStat,
	}, nil
}

func NewTaskStatSampler(procfsFetcher *procfs.ProcfsFetcher) (*TaskStatSampler, error) {
    var err error
    if procfsFetcher == nil {
        procfsFetcher, err = procfs.NewProcfsFetcher()
        if err != nil {
            return nil, err
        }
    }

    return &TaskStatSampler{
        procfsFetcher: procfsFetcher,
    }, nil
}
