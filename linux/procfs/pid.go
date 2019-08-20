/*
 * Author: fasion
 * Created time: 2019-06-27 18:27:39
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 10:53:53
 */

package procfs

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"
	"unsafe"

	"github.com/fasionchan/osutil-go/linux/c"
)

const (
	ProcPidAuxvPath = "/proc/%d/auxv"
	ProcPidStatPath = "/proc/%d/stat"
	ProcPidFdPath = "/proc/%d/fd"
)

var BadAuxvError = fmt.Errorf("bad auxv")

type PidStat struct {
	FetchTime time.Time
	Pid int
	Comm string
	State string
	Ppid int
	Pgrp int
	Session int
	TtyNr int
	TtyMajor int
	TtyMinor int
	Tpgid int
	Flags uint32
	Minflt uint64
	Cminflt uint64
	Majflt uint64
	Cmajflt uint64
	Utime uint64
	Stime uint64
	Cutime uint64
	Cstime uint64
	Priority uint64
	Nice uint64
	NumThreads uint64
	StartTime uint64
	Vsize uint64
	Rss uint64

}

func (self *PidStat) Sub(other *PidStat) (*PidStat) {
	dup := *self

	dup.Utime -= other.Utime
	dup.Stime -= other.Stime

	return &dup
}

func FetchPidStat(pid int) (*PidStat, error) {
	fetchTime := time.Now()

    content, err := ioutil.ReadFile(fmt.Sprintf(ProcPidStatPath, pid))
    if err != nil {
        return nil, err
	}

	stat := PidStat{
		FetchTime: fetchTime,
	}

	var ignoredInt int64
	reader := bytes.NewReader(content)
	_, err = fmt.Fscanf(reader, "%d %s %s %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d %d",
		&stat.Pid,
		&stat.Comm,
		&stat.State,
		&stat.Ppid,
		&stat.Pgrp,
		&stat.Session,
		&stat.TtyNr,
		&stat.Tpgid,
		&stat.Flags,
		&stat.Minflt,
		&stat.Cminflt,
		&stat.Majflt,
		&stat.Cmajflt,
		&stat.Utime,
		&stat.Stime,
		&stat.Cutime,
		&stat.Cstime,
		&stat.Priority,
		&stat.Nice,
		&stat.NumThreads,
		&ignoredInt,
		&stat.StartTime,
		&stat.Vsize,
		&stat.Rss,
	)
	if err != nil {
		return nil, err
	}

	return &stat, nil
}

type ElfAux struct {
	id c.UnsignedLong_c
	value c.UnsignedLong_c
}

type ElfAuxArray []ElfAux

func (self ElfAuxArray) Get(id c.UnsignedLong_c) (c.UnsignedLong_c, bool) {
	for _, aux := range self {
		if aux.id == id {
			return aux.value, true
		}
	}
	return 0, false
}

func (self ElfAuxArray) AT_CLKTCK() (c.UnsignedLong_c, bool) {
	return self.Get(c.AT_CLKTCK)
}

func FetchPidAuxv(pid int) (ElfAuxArray, error) {
    content, err := ioutil.ReadFile(fmt.Sprintf(ProcPidAuxvPath, pid))
    if err != nil {
        return nil, err
	}

	var aux ElfAux
	auxSize := unsafe.Sizeof(aux)

	auxv := make(ElfAuxArray, 0)

	for len(content) >= int(auxSize) {
		raw := content[:auxSize]
		auxv = append(auxv, *(*ElfAux)(unsafe.Pointer(&raw[0])))

		content = content[auxSize:]
	}

	if len(content) > 0 {
		return nil, BadAuxvError
	}

	return auxv, nil
}

func CountPidFds(pid int) (uint64, error) {
	files, err := ioutil.ReadDir(fmt.Sprintf(ProcPidFdPath, pid))
	return uint64(len(files)), err
}
