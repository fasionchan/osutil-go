/*
 * Author: fasion
 * Created time: 2019-06-18 16:14:10
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 17:43:03
 */

package procfs

import (
    "fmt"
    "io/ioutil"
    "sort"
    "strconv"
    "strings"
    "time"
)

var _ = fmt.Println

const (
    PROC_MOUNTINFO_PATH = "/proc/self/mountinfo"
)

type MountInfo struct {
    Id              int
    ParentId        int
    Major           int
    Minor           int
    Root            string
    MountPoint      string
    MountOptions    string
    Optional        string
    FileSystemType  string
    Source          string
    SuperOptions    string
}

type MountInfos struct {
    FetchTime time.Time
    Infos []MountInfo
}

func ParseMountInfoFirstHalf(text string, info *MountInfo) (error) {
    var err error

    fields := strings.Split(text, " ")
    nFields := len(fields)

    info.Id, err = strconv.Atoi(fields[0])
    if err != nil {
        return err
    }

    info.ParentId, err = strconv.Atoi(fields[1])
    if err != nil {
        return err
    }

    parts := strings.Split(fields[2], ":")

    info.Major, err = strconv.Atoi(parts[0])
    if err != nil {
        return err
    }

    info.Minor, err = strconv.Atoi(parts[1])
    if err != nil {
        return err
    }

    info.Root = fields[3]
    info.MountPoint = fields[4]

    if nFields > 5 {
        info.MountOptions = fields[5]
    }

    if nFields > 6 {
        info.Optional = fields[6]
    }

    return nil
}

func ParseMountInfoSecondHalf(text string, info *MountInfo) (error) {
        fields := strings.Split(text, " ")
        nFields := len(fields)

        if nFields > 0 {
            info.FileSystemType = fields[0]
        }

        if nFields > 1 {
            info.Source = fields[1]
        }

        if nFields > 2 {
            info.SuperOptions = fields[2]
        }

        return nil
}

func FetchMountInfos() (*MountInfos, error) {
    var err error

    fetchTime := time.Now()

    content, err := ioutil.ReadFile(PROC_MOUNTINFO_PATH)
    if err != nil {
        return nil, err
    }

    lines := strings.Split(string(content), "\n")
    infos := make([]MountInfo, 0, len(lines))

    for _, line := range lines {
        if line == "" {
            continue
        }

        info := MountInfo{}

        parts := strings.Split(line, " - ")

        err := ParseMountInfoFirstHalf(parts[0], &info)
        if err != nil {
            return nil, err
        }

        err = ParseMountInfoSecondHalf(parts[1], &info)
        if err != nil {
            return nil, err
        }

        infos = append(infos, info)
    }

    return &MountInfos{
        fetchTime,
        infos,
    }, nil
}

type ByDeviceMountPoint []MountInfo

func (self ByDeviceMountPoint) Len() (int) { return len(self); }
func (self ByDeviceMountPoint) Swap(i, j int) { self[i], self[j] = self[j], self[i] }
func (self ByDeviceMountPoint) Less(i, j int) (bool) {
    if self[i].Major < self[j].Major {
        return true
    } else if self[i].Major > self[j].Major {
        return false
    }

    if self[i].Minor < self[j].Minor {
        return true
    } else if self[i].Minor > self[j].Minor {
        return false
    }

    return self[i].MountPoint < self[j].MountPoint
}

func FilterMountInfos(infos []MountInfo) ([]MountInfo) {
    sort.Sort(ByDeviceMountPoint(infos))
    remains := make([]MountInfo, 0, len(infos))

    for _, info := range infos {
        nRemains := len(remains)
        if nRemains > 0 {
            last := &remains[nRemains-1]
            if info.Major == last.Major && info.Minor == last.Minor {
                continue
            }
        }

		if info.Major == 0 ||
			info.Major == 7 ||
			info.Major == 144 ||
			info.Major == 145 ||
			info.Major == 146 {
            continue
        }

        if info.Major >= 240 && info.Major <= 254 {
        }

        remains = append(remains, info)
    }

    return remains
}
