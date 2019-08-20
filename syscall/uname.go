// +build linux

/*
 * Author: fasion
 * Created time: 2019-07-12 15:14:15
 * Last Modified by: fasion
 * Last Modified time: 2019-07-12 15:48:09
 */

package syscall

import (
	"syscall"
	"strings"
)

type UnameInfo struct {
	Sysname string
	Nodename string
	Release string
	Version string
	Machine string
	Domainname string
}

func StringFromInt8Slice(slice []int8) string {
	buf := make([]byte, 0)
	for _, value := range slice {
		buf = append(buf, byte(value))
	}
	return strings.TrimRight(string(buf), "\x00")
}

func Uname() (*UnameInfo, error) {
	var data syscall.Utsname

	err := syscall.Uname(&data)
	if err != nil {
		return nil, err
	}

	var info UnameInfo

	info.Sysname = StringFromInt8Slice(data.Sysname[:])
	info.Nodename = StringFromInt8Slice(data.Nodename[:])
	info.Release = StringFromInt8Slice(data.Release[:])
	info.Version = StringFromInt8Slice(data.Version[:])
	info.Machine = StringFromInt8Slice(data.Machine[:])
	info.Domainname = StringFromInt8Slice(data.Domainname[:])

	return &info, nil
}
