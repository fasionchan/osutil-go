/*
 * Author: fasion
 * Created time: 2019-08-01 14:45:31
 * Last Modified by: fasion
 * Last Modified time: 2019-08-09 16:16:53
 */

package perf

import (
	"fmt"
    "time"
    "github.com/fasionchan/osutil-go/linux/procfs"
)

var _ = fmt.Println

type DiskPerf struct {
	Reads float64
	MergedReads float64
	BytesRead float64
	ReadUtil float64

	Writes float64
	MergedWrites float64
	BytesWritten float64
	WriteUtil float64

	CurrentIos float64

	IoUtil float64
	QueueSize float64
}

type DiskPerfSample struct {
	FetchTime time.Time
	Perfs map[string]DiskPerf
}

type DiskPerfSampler struct {
	procfsFetcher *procfs.ProcfsFetcher

	lastStats *procfs.DiskStats
}

func NewDiskPerfSampler(procfsFetcher *procfs.ProcfsFetcher) (*DiskPerfSampler, error) {
	return &DiskPerfSampler{
		procfsFetcher: procfsFetcher,
	}, nil
}

func (self *DiskPerfSampler) Sample() (*DiskPerfSample, error) {
	stats, err := self.procfsFetcher.FetchDiskStats()
	if err != nil {
		return nil, err
	}

	stats.FilterMisc()

	lastStats := self.lastStats
	self.lastStats = stats

	if lastStats == nil {
		return nil, nil
	}

	seconds := stats.FetchTime.Sub(lastStats.FetchTime).Seconds()

	name2Perf := make(map[string]DiskPerf)
	name2LastStat := lastStats.NameMapping()

	for _, stat := range stats.Stats {
		lastStat, ok := name2LastStat[stat.Name]
		if !ok {
			continue
		}

		delta := stat.Counter.Sub(lastStat.Counter)

		perf := DiskPerf{
			Reads: float64(delta.Reads) / seconds,
			MergedReads: float64(delta.MergedReads) / seconds,
			BytesRead: float64(delta.BytesRead) / seconds,
			ReadUtil: float64(delta.TimeReading) / 1000 / seconds * 100,

			Writes: float64(delta.Writes) / seconds,
			MergedWrites: float64(delta.MergedWrites) / seconds,
			BytesWritten: float64(delta.BytesWritten) / seconds,
			WriteUtil: float64(delta.TimeWriting) / 1000 / seconds * 100,

			CurrentIos: float64(delta.CurrentIos),

			IoUtil: float64(delta.IoTime) / 1000 / seconds * 100,
			QueueSize: float64(delta.WeightedIoTime) / 1000 / seconds,
		}

		name2Perf[stat.Name] = perf
	}

	return &DiskPerfSample{
		FetchTime: stats.FetchTime,
		Perfs: name2Perf,
	}, nil
}
