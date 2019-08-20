/*
 * Author: fasion
 * Created time: 2019-06-18 15:48:02
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:44:36
 */

package procfs

import (
    "fmt"
    "io/ioutil"
    "runtime"
    "strconv"
    "strings"
    "time"
    "github.com/fasionchan/libgo/parse"
    "github.com/fasionchan/libgo/copy"
    mathutil "github.com/fasionchan/libgo/math"
)

const (
    PROC_STAT_PATH = "/proc/stat"
)

var _ = fmt.Println

type StatCounter struct {
    FetchTime time.Time
    Cpu []uint64
    Cpus []([]uint64)
    Intr []uint64
    Ctxt uint64
    Btime uint64
    Processes uint64
    ProcsRunning uint64
    ProcsBlocked uint64
    SoftirqCounter []uint64
}

func (self *StatCounter) DeepCopy() (*StatCounter) {
    return &StatCounter{
        Cpu: copy.DeepCopyUint64Slice(self.Cpu),
        Cpus: copy.DeepCopyUint64Slice2D(self.Cpus),
        Intr: copy.DeepCopyUint64Slice(self.Intr),
        Ctxt: self.Ctxt,
        Btime: self.Btime,
        Processes: self.Processes,
        ProcsRunning: self.ProcsRunning,
        ProcsBlocked: self.ProcsBlocked,
        SoftirqCounter: copy.DeepCopyUint64Slice(self.SoftirqCounter),
    }
}

func (self *StatCounter) Sub(other *StatCounter) (*StatCounter) {
    dup := self.DeepCopy()

    mathutil.Uint64SliceSub(dup.Cpu, other.Cpu, true)
    mathutil.Uint64SliceSub2D(dup.Cpus, other.Cpus, true)
    mathutil.Uint64SliceSub(dup.Intr, other.Intr, true)
    dup.Ctxt -= other.Ctxt
    dup.Processes -= other.Processes
    mathutil.Uint64SliceSub(dup.SoftirqCounter, other.SoftirqCounter, true)

    return dup
}

func FetchStatCounter() (*StatCounter, error) {
    var err error

    cores := runtime.NumCPU()

    fetchTime := time.Now()
    content, err := ioutil.ReadFile(PROC_STAT_PATH)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(content), "\n")

    stat := StatCounter{
        FetchTime: fetchTime,
        Cpus: make([]([]uint64), 0, cores),
    }

    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) == 0 {
            continue
        }

        if fields[0] == "cpu" {
            stat.Cpu, err = parse.ParseUint64Slice(fields[1:], stat.Cpu)
        } else if strings.HasPrefix(fields[0], "cpu") {
            values, err := parse.ParseUint64Slice(fields[1:], nil)
            if err == nil {
                stat.Cpus = append(stat.Cpus, values)
            }
        }

        if fields[0] == "intr" {
            stat.Intr, err = parse.ParseUint64Slice(fields[1:], stat.Intr)
        }

        if fields[0] == "ctxt" {
            stat.Ctxt, err = strconv.ParseUint(fields[1], 10, 64)
        }

        if fields[0] == "btime" {
            stat.Btime, err = strconv.ParseUint(fields[1], 10, 64)
        }

        if fields[0] == "processes" {
            stat.Processes, err = strconv.ParseUint(fields[1], 10, 64)
        }

        if fields[0] == "procs_running" {
            stat.ProcsRunning, err = strconv.ParseUint(fields[1], 10, 64)
        }

        if fields[0] == "procs_blocked" {
            stat.ProcsBlocked, err = strconv.ParseUint(fields[1], 10, 64)
        }

        if fields[0] == "softirq" {
            stat.SoftirqCounter, err = parse.ParseUint64Slice(fields[1:], stat.SoftirqCounter)
        }

        if err != nil {
            return nil, err
        }
    }

    return &stat, nil
}
