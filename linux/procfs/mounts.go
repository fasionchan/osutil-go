/*
 * Author: fasion
 * Created time: 2019-06-18 16:14:10
 * Last Modified by: fasion
 * Last Modified time: 2019-12-23 15:04:38
 */

package procfs

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var _ = fmt.Println

const (
	ProcMountsPath = "/proc/mounts"
	EtcMtabPath    = "/etc/mtab"
)

type MountPoint struct {
	Spec    string
	Path    string
	Type    string
	Options []string
	Freq    int
	PassNo  int
}

func MountPoints2MountInfos(mountpoints []MountPoint) []MountInfo {
	mountinfos := make([]MountInfo, 0, len(mountpoints))
	for _, mountpoint := range mountpoints {
		var major, minor int
		var stat syscall.Stat_t
		if err := syscall.Stat(mountpoint.Spec, &stat); err == nil {
			major = int(stat.Rdev >> 8)
			minor = int(stat.Rdev & 0xff)
		}

		mountinfos = append(mountinfos, MountInfo{
			MountPoint:     mountpoint.Path,
			FileSystemType: mountpoint.Type,
			Major:          major,
			Minor:          minor,
		})
	}

	return mountinfos
}

type MountPoints struct {
	FetchTime time.Time
	Records   []MountPoint
}

func FetchMountPointsByPath(path string) (*MountPoints, error) {
	fetchTime := time.Now()
	content, err := ioutil.ReadFile(path)
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
			Spec:    fields[0],
			Path:    fields[1],
			Type:    fields[2],
			Options: strings.Split(fields[3], ","),
			Freq:    freq,
			PassNo:  passno,
		})
	}

	return &MountPoints{
		FetchTime: fetchTime,
		Records:   mounts,
	}, nil
}

func FetchMountPoints() (*MountPoints, error) {
	return FetchMountPointsByPath(ProcMountsPath)
}

func FetchMountPointsFromEtcMtab() (*MountPoints, error) {
	return FetchMountPointsByPath(EtcMtabPath)
}
