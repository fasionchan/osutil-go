/*
 * Author: fasion
 * Created time: 2019-06-18 16:14:10
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:44:04
 */

package procfs

import (
    //"encoding/json"
    "fmt"
    "io/ioutil"
    "runtime"
    "strconv"
    "strings"
    "time"
)

const (
    PROC_SOFTIRQS_PATH = "/proc/softirqs"
)

type SoftirqCounter struct {
    FetchTime time.Time
    Cores int
    Names []string
    CounterArray []([]uint64)
    CpuCounters []uint64
}

func (self *SoftirqCounter) Sub(other *SoftirqCounter) (*SoftirqCounter) {
    dup := SoftirqCounter{
        Cores: self.Cores,
        Names: self.Names,
        CounterArray: make([]([]uint64), 0, len(self.CounterArray)),
        CpuCounters: make([]uint64, len(self.CpuCounters)),
    }

    for i, values := range other.CounterArray {
        dup.CounterArray = append(
            dup.CounterArray,
            make([]uint64, len(self.CounterArray[0])),
        )

        copy(dup.CounterArray[i], self.CounterArray[i])

        for j, value := range values {
            dup.CounterArray[i][j] -= value
        }
    }

    copy(dup.CpuCounters, self.CpuCounters)

    for i, value := range other.CpuCounters {
        dup.CpuCounters[i] -= value
    }

    return &dup
}

func FetchSoftirqCounters() (*SoftirqCounter, error) {
    cores := runtime.NumCPU()

    fetchTime := time.Now()
    content, err := ioutil.ReadFile(PROC_SOFTIRQS_PATH)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(content), "\n")

    counter := SoftirqCounter{
        FetchTime: fetchTime,
        Cores: cores,
        Names: make([]string, 0, len(lines)-1),
        CounterArray: make([]([]uint64), 0, 5),
        CpuCounters: make([]uint64, cores),
    }

    for i, line := range lines {
        if i == 0 {
            continue
        }

        fields := strings.Fields(line)
        if len(fields) == 0 {
            continue
        }

        if false {
            fmt.Println(len(fields), fields)
        }

        name := fields[0]
        name = name[:len(name)-1]
        counter.Names = append(counter.Names, name)

        values := make([]uint64, 0, cores)
        for core, text := range fields[1:] {
            value, err := strconv.ParseUint(text, 10, 64)
            if err != nil {
                return nil, err
            }

            values = append(values, value)
            counter.CpuCounters[core] += value
        }

        counter.CounterArray = append(counter.CounterArray, values)
    }

    return &counter, nil
}
