/*
 * Author: fasion
 * Created time: 2019-06-18 16:14:10
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:43:40
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

const (
    PROC_MOUNTS_PATH = "/proc/mounts"
)

type MountPoint struct {
    Spec string
    Path string
    Type string
    Options []string
    Freq int
    PassNo int
}

type MountPoints struct {
    FetchTime time.Time
    Records []MountPoint
}

func FetchMountPoints() (*MountPoints, error) {
    fetchTime := time.Now()
    content, err := ioutil.ReadFile(PROC_MOUNTS_PATH)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(content), "\n")
    mounts := make([]MountPoint, 0, len(lines))

    for _, line := range lines {
        if line == "" {
            continue
        }

        fields := strings.Split(line, " ")

        freq, err := strconv.Atoi(fields[4])
        if err != nil {
            return nil, err
        }

        passno, err := strconv.Atoi(fields[5])
        if err != nil {
            return nil, err
        }

        mounts = append(mounts, MountPoint{
            Spec: fields[0],
            Path: fields[1],
            Type: fields[2],
            Options: strings.Split(fields[3], ","),
            Freq: freq,
            PassNo: passno,
        })
    }

    return &MountPoints{
        FetchTime: fetchTime,
        Records: mounts,
    }, nil
}
