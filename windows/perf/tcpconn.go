/*
 * Author: fasion
 * Created time: 2019-08-19 16:07:16
 * Last Modified by: fasion
 * Last Modified time: 2019-08-20 09:22:22
 */

package perf

import (
	"time"
)

type TcpStateStat struct {
	Close uint64
	Listen uint64
	SynSent uint64
	SynRecv uint64
	Established uint64
	FinWait1 uint64
	FinWait2 uint64
	CloseWait uint64
	Closing uint64
	LastAck uint64
	TimeWait uint64
	DeleteTcb uint64
}

type TcpStateStatSample struct {
	FetchTime time.Time
	StateStat TcpStateStat
}

type TcpStateStatSampler struct {

}

func NewTcpStateStatSampler() (*TcpStateStatSampler, error) {
	return &TcpStateStatSampler{}, nil
}

func (self *TcpStateStatSampler) Sample() (*TcpStateStatSample, error) {
	rows, err := FetchTcpTable()
	if err != nil {
		return nil, err
	}

	stateStat := make([]uint64, 0x100)
	for _, row := range rows {
		stateStat[row.State & 0xff]++
	}

	return &TcpStateStatSample{
		FetchTime: time.Now(),
		StateStat: TcpStateStat{
			Close: stateStat[MIB_TCP_STATE_CLOSED],
			Listen: stateStat[MIB_TCP_STATE_LISTEN],
			SynSent: stateStat[MIB_TCP_STATE_SYN_SENT],
			SynRecv: stateStat[MIB_TCP_STATE_SYN_RCVD],
			Established: stateStat[MIB_TCP_STATE_ESTAB],
			FinWait1: stateStat[MIB_TCP_STATE_FIN_WAIT1],
			FinWait2: stateStat[MIB_TCP_STATE_FIN_WAIT2],
			CloseWait: stateStat[MIB_TCP_STATE_CLOSE_WAIT],
			Closing: stateStat[MIB_TCP_STATE_CLOSING],
			LastAck: stateStat[MIB_TCP_STATE_LAST_ACK],
			TimeWait: stateStat[MIB_TCP_STATE_TIME_WAIT],
			DeleteTcb: stateStat[MIB_TCP_STATE_DELETE_TCB],
		},
	}, nil
}
