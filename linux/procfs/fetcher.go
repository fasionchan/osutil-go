/*
 * Author: fasion
 * Created time: 2019-06-18 16:00:16
 * Last Modified by: fasion
 * Last Modified time: 2019-08-01 14:45:14
 */

package procfs

import (
    "sync"
)

type Object interface {}
type Handler func() (Object, error)

type SubFetcher struct {
    handler Handler

    mutex sync.RWMutex
    data Object
}

func (self *SubFetcher) Reset() {
    self.data = nil
}

func (self *SubFetcher) Fetch() (Object, error) {
    self.mutex.RLock()

    if self.data != nil {
		defer self.mutex.RUnlock()
		return self.data, nil
    }

    self.mutex.RUnlock()

    self.mutex.Lock()
    defer self.mutex.Unlock()

    if self.data != nil {
        return self.data, nil
    }

    var err error
    self.data, err = self.handler()
    if err != nil {
        return nil, err
    }

    return self.data, nil
}

type ProcfsFetcher struct {
    diskStatFetcher SubFetcher
    loadAvgFetcher SubFetcher
    meminfoFetcher SubFetcher
    mountPointFetcher SubFetcher
    mountInfoFetcher SubFetcher
    netDevCounterFetcher SubFetcher
    softirqCounterFetcher SubFetcher
    statCounterFetcher SubFetcher
	vmStatFetcher SubFetcher
	pidStats map[int]*PidStat
}

func (self *ProcfsFetcher) ResetAll() {
    self.diskStatFetcher.Reset()
    self.loadAvgFetcher.Reset()
    self.meminfoFetcher.Reset()
    self.mountPointFetcher.Reset()
    self.mountInfoFetcher.Reset()
    self.netDevCounterFetcher.Reset()
    self.softirqCounterFetcher.Reset()
    self.statCounterFetcher.Reset()
	self.vmStatFetcher.Reset()
	self.pidStats = make(map[int]*PidStat)
}

func NewProcfsFetcher() (*ProcfsFetcher, error) {
    return &ProcfsFetcher{
        diskStatFetcher: SubFetcher{
            handler: func() (Object, error) {
                return FetchDiskStats()
            },
        },

        loadAvgFetcher: SubFetcher{
            handler: func() (Object, error) {
                return FetchLoadAvg()
            },
        },

        meminfoFetcher: SubFetcher{
            handler: func() (Object, error) {
                return FetchMemInfo()
            },
        },

        mountPointFetcher: SubFetcher{
            handler: func() (Object, error) {
                return FetchMountPoints()
            },
        },

        mountInfoFetcher: SubFetcher{
            handler: func() (Object, error) {
                return FetchMountInfos()
            },
        },

        netDevCounterFetcher: SubFetcher{
            handler: func() (Object, error) {
                return FetchNetDevCounters()
            },
        },

        softirqCounterFetcher: SubFetcher{
            handler: func() (Object, error) {
                return FetchSoftirqCounters()
            },
        },

        statCounterFetcher: SubFetcher{
            handler: func() (Object, error) {
                return FetchStatCounter()
            },
        },

        vmStatFetcher: SubFetcher{
            handler: func() (Object, error) {
                return FetchVMStat()
            },
		},
		pidStats: make(map[int]*PidStat),
    }, nil
}

func (self *ProcfsFetcher) FetchDiskStats() (*DiskStats, error) {
    data, err := self.diskStatFetcher.Fetch()
    if data == nil {
        return nil, err
    }

    return data.(*DiskStats), err
}

func (self *ProcfsFetcher) FetchLoadAvg() (*LoadAvgSample, error) {
    data, err := self.loadAvgFetcher.Fetch()
    if data == nil {
        return nil, err
    }

    return data.(*LoadAvgSample), err
}

func (self *ProcfsFetcher) FetchMemInfo() (*MemInfo, error) {
    data, err := self.meminfoFetcher.Fetch()
    if data == nil {
        return nil, err
    }

    return data.(*MemInfo), err
}

func (self *ProcfsFetcher) FetchMountPoints() ([]MountPoint, error) {
    data, err := self.mountPointFetcher.Fetch()
    if data == nil {
        return nil, err
    }

    return data.([]MountPoint), err
}

func (self *ProcfsFetcher) FetchMountInfos() (*MountInfos, error) {
    data, err := self.mountInfoFetcher.Fetch()
    if data == nil {
        return nil, err
    }

    return data.(*MountInfos), err
}

func (self *ProcfsFetcher) FetchNetDevCounters() (*NetDevCounters, error) {
    data, err := self.netDevCounterFetcher.Fetch()
    if data == nil {
        return nil, err
    }

    return data.(*NetDevCounters), err
}

func (self *ProcfsFetcher) FetchSoftirqCounters() (*SoftirqCounter, error) {
    data, err := self.softirqCounterFetcher.Fetch()
    if data == nil {
        return nil, err
    }

    return data.(*SoftirqCounter), err
}

func (self *ProcfsFetcher) FetchStatCounter() (*StatCounter, error) {
	data, err := self.statCounterFetcher.Fetch()
    if data == nil {
        return nil, err
    }

    return data.(*StatCounter), err
}

func (self *ProcfsFetcher) FetchVMStat() (*VMStat, error) {
    data, err := self.vmStatFetcher.Fetch()
    if data == nil {
        return nil, err
    }

    return data.(*VMStat), err
}

func (self *ProcfsFetcher) FetchPidStat(pid int) (*PidStat, error) {
	stat, ok := self.pidStats[pid]
	if ok {
		return stat, nil
	}

	stat, err := FetchPidStat(pid)
	if err != nil {
		return nil, err
	}

	self.pidStats[pid] = stat

	return stat, nil
}

func (self *ProcfsFetcher) FetchPidAuxv(pid int) (ElfAuxArray, error) {
	return FetchPidAuxv(pid)
}
