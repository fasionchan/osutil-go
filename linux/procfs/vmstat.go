/*
 * Author: fasion
 * Created time: 2019-06-18 16:14:10
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:44:34
 */

package procfs

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "strings"
    "time"
)

const (
    PROC_VMSTAT_PATH = "/proc/vmstat"
)

type VMStat struct {
    FetchTime time.Time
    Info map[string]uint64
}

func FetchVMStat() (*VMStat, error) {
    if false {
        fmt.Println()
    }

    fetchTime := time.Now()
    content, err := ioutil.ReadFile(PROC_VMSTAT_PATH)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(content), "\n")

    info := make(map[string]uint64)
    for _, line := range lines {
        if line == "" {
            continue
        }

        parts := strings.Fields(line)

        name := parts[0]
        value, err := strconv.ParseUint(parts[1], 10, 64)
        if err != nil {
            return nil, err
        }

        info[name] = value
    }

    return &VMStat{
        FetchTime: fetchTime,
        Info: info,
    }, nil
}
