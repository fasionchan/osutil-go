/*
 * Author: fasion
 * Created time: 2019-06-18 16:14:10
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:42:48
 */

package procfs

import (
    "fmt"
    "io/ioutil"
    "strconv"
    "strings"
    "time"

    unitutil "github.com/fasionchan/libgo/unit"
)

var _ = fmt.Println

const (
    PROC_MEMINFO_PATH = "/proc/meminfo"
)

var units = map[string]uint64{
    "kB": unitutil.KiB,
}

type MemInfo struct {
    FetchTime time.Time

    Info map[string]unitutil.Bytes
}

func FetchMemInfo() (*MemInfo, error) {
    fetchTime := time.Now()
    content, err := ioutil.ReadFile(PROC_MEMINFO_PATH)
    if err != nil {
        return nil, err
    }

    info := make(map[string]unitutil.Bytes)

    lines := strings.Split(string(content), "\n")
    for _, line := range lines {
        fields := strings.Fields(line)
        if len(fields) == 0 {
            continue
        }

        name := fields[0][:len(fields[0])-1]

        value, err := strconv.ParseUint(fields[1], 0, 64)
        if err != nil {
            return nil, err
        }

        var unit uint64 = 1
        if len(fields) > 2 {
            var ok bool
            unit, ok = units[fields[2]]
            if !ok {
                return nil ,nil
            }
        }

        info[name] = unitutil.Bytes(value * unit)
    }

    return &MemInfo{
        FetchTime: fetchTime,
        Info: info,
    }, nil
}
