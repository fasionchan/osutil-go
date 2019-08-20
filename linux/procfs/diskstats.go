/*
 * Author: fasion
 * Created time: 2019-06-18 16:14:10
 * Last Modified by: fasion
 * Last Modified time: 2019-08-02 16:27:36
 */

package procfs

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "strings"
    "time"

    "github.com/fasionchan/libgo/parse"
)

var _ = fmt.Println

const (
    ProcDiskStatsPath = "/proc/diskstats"
    SectorBytes = 512
)

type DiskCounter struct {
    Reads uint64
    MergedReads uint64
    BytesRead uint64
    TimeReading uint64

    Writes uint64
    MergedWrites uint64
    BytesWritten uint64
    TimeWriting uint64

    CurrentIos uint64

    IoTime uint64
    WeightedIoTime uint64
}

func (self DiskCounter) IsZero() (bool) {
    return self.Reads == 0 &&
        self.Writes == 0 &&
        self.IoTime == 0
}

func (self DiskCounter) Sub(other DiskCounter) (DiskCounter) {
    return DiskCounter{
        Reads: self.Reads - other.Reads,
        MergedReads: self.MergedReads - other.MergedReads,
        BytesRead: self.BytesRead - other.BytesRead,
        TimeReading: self.TimeReading - other.TimeReading,

        Writes: self.Writes - other.Writes,
        MergedWrites: self.MergedWrites - other.MergedWrites,
        BytesWritten: self.BytesWritten - other.BytesWritten,
        TimeWriting: self.TimeWriting - other.TimeWriting,

        IoTime: self.IoTime - other.IoTime,
        WeightedIoTime: self.WeightedIoTime - other.WeightedIoTime,
    }
}

type DiskStat struct {
    Major int
    Minor int
    Name string
    Counter DiskCounter
}

func FilterMiscDiskStats(stats []DiskStat) ([]DiskStat) {
    results := make([]DiskStat, 0)

    lastName := ""
    for _, stat := range stats {
        // filter idle
        if stat.Counter.IsZero() {
            continue
        }

        // filter nonphysical
        if stat.Major == 1 ||
            stat.Major == 7 ||
            stat.Major == 11 {
            continue
        }

        // filter subdevice
        if lastName != "" && strings.HasPrefix(stat.Name, lastName) {
            continue
        }

        lastName = stat.Name
        results = append(results, stat)
    }

    return results
}

type DiskStats struct {
    FetchTime time.Time
    Stats []DiskStat
}

func (self *DiskStats) FilterMisc() {
    self.Stats = FilterMiscDiskStats(self.Stats)
}

func (self *DiskStats) NameMapping() (map[string]DiskStat) {
    n2d := make(map[string]DiskStat)
    for _, stat := range self.Stats {
        n2d[stat.Name] = stat
    }
    return n2d
}

func FetchDiskStats() (*DiskStats, error) {
    fetchTime := time.Now()
    content, err := ioutil.ReadFile(ProcDiskStatsPath)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(content), "\n")
    stats := make([]DiskStat, 0, len(lines))

    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) == 0 {
            continue
        }

        major, err := strconv.Atoi(fields[0])
        if err != nil {
            return nil, err
        }

        minor, err := strconv.Atoi(fields[1])
        if err != nil {
            return nil, err
        }

        values, err := parse.ParseUint64Slice(fields[3:], nil)
        if err != nil {
            return nil, err
        }

        if len(values) < 11 {
            continue
        }

        stats = append(stats, DiskStat{
            Major: major,
            Minor: minor,
            Name: fields[2],
            Counter: DiskCounter{
                Reads: values[0],
                MergedReads: values[1],
                BytesRead: values[2] * SectorBytes,
                TimeReading: values[3],

                Writes: values[4],
                MergedWrites: values[5],
                BytesWritten: values[6] * SectorBytes,
                TimeWriting: values[7],

                CurrentIos: values[8],

                IoTime: values[9],
                WeightedIoTime: values[10],
            },
        })
    }

    return &DiskStats{
        FetchTime: fetchTime,
        Stats: stats,
    }, nil
}
