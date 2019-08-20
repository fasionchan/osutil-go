/*
 * Author: fasion
 * Created time: 2019-06-18 16:14:10
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:42:35
 */

package procfs

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "strings"
    "time"
)

var _ = fmt.Println

type LoadAvg struct {
    LoadAvgM1 float64
    LoadAvgM5 float64
    LoadAvgM15 float64
}

type LoadAvgSample struct {
    FetchTime time.Time

    LoadAvg LoadAvg

    TasksRunning uint64
    Tasks uint64

    LastPid uint64
}

const (
    PROC_LOADAVG_PATH = "/proc/loadavg"
)

func FetchLoadAvg() (*LoadAvgSample, error) {
    var err error

    fetchTime := time.Now()
    content, err := ioutil.ReadFile(PROC_LOADAVG_PATH)
    if err != nil {
        return nil, err
    }

    parts := strings.Split(string(content), " ")
    subparts := strings.Split(parts[3], "/")

    loadAvgM1, err := strconv.ParseFloat(parts[0], 64)
    if err != nil {
        return nil, err
    }

    loadAvgM5, err := strconv.ParseFloat(parts[1], 64)
    if err != nil {
        return nil, err
    }

    loadAvgM15, err := strconv.ParseFloat(parts[2], 64)
    if err != nil {
        return nil, err
    }

    tasksRunning, err := strconv.ParseUint(subparts[0], 10, 64)
    if err != nil {
        return nil, err
    }

    tasks, err := strconv.ParseUint(subparts[1], 10, 16)
    if err != nil {
        return nil, err
    }

    lastPid, err := strconv.ParseUint(strings.TrimSpace(parts[4]), 10, 64)
    if err != nil {
        return nil, err
    }

    return &LoadAvgSample{
        FetchTime: fetchTime,
        LoadAvg: LoadAvg{
            LoadAvgM1: loadAvgM1,
            LoadAvgM5: loadAvgM5,
            LoadAvgM15: loadAvgM15,
        },
        TasksRunning: tasksRunning,
        Tasks: tasks,
        LastPid: lastPid,
    }, nil
}
