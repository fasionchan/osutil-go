/*
 * Author: fasion
 * Created time: 2019-08-07 10:36:13
 * Last Modified by: fasion
 * Last Modified time: 2019-12-20 12:52:16
 */

package perf

import (
	"fmt"
	"syscall"
	"time"

	"github.com/fasionchan/osutil-go/linux"
	"github.com/fasionchan/osutil-go/linux/c"
	"github.com/fasionchan/osutil-go/linux/netlink"
	"github.com/fasionchan/osutil-go/linux/procfs"
)

var _ = fmt.Println

type TcpStateStat struct {
	Established uint64
	SynSent     uint64
	SynRecv     uint64
	FinWait1    uint64
	FinWait2    uint64
	TimeWait    uint64
	Close       uint64
	CloseWait   uint64
	LastAck     uint64
	Listen      uint64
	Closing     uint64
}

type TcpConnSample struct {
	FetchTime time.Time
	StateStat TcpStateStat
}

type TcpConnSampler struct {
}

func NewTcpConnSampler() (*TcpConnSampler, error) {
	return &TcpConnSampler{}, nil
}

func (self *TcpConnSampler) SampleBySockDiag() (*TcpConnSample, error) {
	sockDiag, err := netlink.NewNetlinkSockDiag()
	if err != nil {
		return nil, err
	}
	defer sockDiag.Close()

	stateStat := make(map[uint8]uint64)

	// find ipv4 tcp diags
	receiver, err := sockDiag.InetDiagRequest(netlink.InetDiagReqV2_c{
		Sdiag_family:   syscall.AF_INET,
		Sdiag_protocol: netlink.IPPROTO_TCP,
		Idiag_states:   0xffffffff,
	}, 1024000)
	if err != nil {
		return nil, err
	}

	for {
		diags, more, err := receiver.Receive()
		if err != nil {
			return nil, err
		}

		for _, diag := range diags {
			stateStat[diag.State] += 1
		}

		if !more {
			break
		}
	}

	// find ipv6 tcp diags
	receiver, err = sockDiag.InetDiagRequest(netlink.InetDiagReqV2_c{
		Sdiag_family:   syscall.AF_INET6,
		Sdiag_protocol: netlink.IPPROTO_TCP,
		Idiag_states:   0xffffffff,
	}, 1024000)
	if err != nil {
		return nil, err
	}

	for {
		diags, more, err := receiver.Receive()
		if err != nil {
			return nil, err
		}

		for _, diag := range diags {
			stateStat[diag.State] += 1
		}

		if !more {
			break
		}
	}

	return &TcpConnSample{
		FetchTime: time.Now(),
		StateStat: TcpStateStat{
			Established: stateStat[c.TCP_ESTABLISHED],
			SynSent:     stateStat[c.TCP_SYN_SENT],
			SynRecv:     stateStat[c.TCP_SYN_RECV],
			FinWait1:    stateStat[c.TCP_FIN_WAIT1],
			FinWait2:    stateStat[c.TCP_FIN_WAIT2],
			TimeWait:    stateStat[c.TCP_TIME_WAIT],
			Close:       stateStat[c.TCP_CLOSE],
			CloseWait:   stateStat[c.TCP_CLOSE_WAIT],
			LastAck:     stateStat[c.TCP_LAST_ACK],
			Listen:      stateStat[c.TCP_LISTEN],
			Closing:     stateStat[c.TCP_CLOSING],
		},
	}, nil
}

func (self *TcpConnSampler) SampleByInetDiag() (*TcpConnSample, error) {
	inetDiag, err := netlink.NewNetlinkInetDiag()
	if err != nil {
		return nil, err
	}
	defer inetDiag.Close()

	stateStat := make(map[uint8]uint64)

	receiver, err := inetDiag.InetDiagRequest(netlink.InetDiagReq_c{
		Family: syscall.AF_INET,
		States: 0xffffffff,
	}, 1024000)
	if err != nil {
		return nil, err
	}

	for {
		diags, more, err := receiver.Receive()
		if err != nil {
			return nil, err
		}

		for _, diag := range diags {
			stateStat[diag.State] += 1
		}

		if !more {
			break
		}
	}

	return &TcpConnSample{
		FetchTime: time.Now(),
		StateStat: TcpStateStat{
			Established: stateStat[c.TCP_ESTABLISHED],
			SynSent:     stateStat[c.TCP_SYN_SENT],
			SynRecv:     stateStat[c.TCP_SYN_RECV],
			FinWait1:    stateStat[c.TCP_FIN_WAIT1],
			FinWait2:    stateStat[c.TCP_FIN_WAIT2],
			TimeWait:    stateStat[c.TCP_TIME_WAIT],
			Close:       stateStat[c.TCP_CLOSE],
			CloseWait:   stateStat[c.TCP_CLOSE_WAIT],
			LastAck:     stateStat[c.TCP_LAST_ACK],
			Listen:      stateStat[c.TCP_LISTEN],
			Closing:     stateStat[c.TCP_CLOSING],
		},
	}, nil
}

func (self *TcpConnSampler) SampleByNetlink() (*TcpConnSample, error) {
	kv, err := linux.FetchKernelVersionNumber()
	if err != nil {
		return nil, err
	}

	if !kv.Before(*linux.MustKernelVersionNumber("3.3")) {
		return self.SampleBySockDiag()
	}

	if !kv.Before(*linux.MustKernelVersionNumber("2.6.14")) {
		return self.SampleByInetDiag()
	}

	return nil, nil
}

func (self *TcpConnSampler) SampleByProcfs() (*TcpConnSample, error) {
	bufferSize := 1 * 1024 * 1024
	stateStat := make(map[string]uint64)

	scanner, err := procfs.NewNetTcp4Scanner(bufferSize)
	if err != nil {
		return nil, err
	}
	defer scanner.Close()

	for scanner.Scan() {
		_, fields, err := scanner.Record()
		if err != nil {
			return nil, err
		}

		stateStat[fields[2]] += 1
	}

	scanner6, err := procfs.NewNetTcp6Scanner(bufferSize)
	if err != nil {
		return nil, err
	}
	defer scanner6.Close()

	for scanner6.Scan() {
		_, fields, err := scanner6.Record()
		if err != nil {
			return nil, err
		}

		stateStat[fields[2]] += 1
	}

	return &TcpConnSample{
		FetchTime: time.Now(),
		StateStat: TcpStateStat{
			Established: stateStat[fmt.Sprintf("%02X", c.TCP_ESTABLISHED)],
			SynSent:     stateStat[fmt.Sprintf("%02X", c.TCP_SYN_SENT)],
			SynRecv:     stateStat[fmt.Sprintf("%02X", c.TCP_SYN_RECV)],
			FinWait1:    stateStat[fmt.Sprintf("%02X", c.TCP_FIN_WAIT1)],
			FinWait2:    stateStat[fmt.Sprintf("%02X", c.TCP_FIN_WAIT2)],
			TimeWait:    stateStat[fmt.Sprintf("%02X", c.TCP_TIME_WAIT)],
			Close:       stateStat[fmt.Sprintf("%02X", c.TCP_CLOSE)],
			CloseWait:   stateStat[fmt.Sprintf("%02X", c.TCP_CLOSE_WAIT)],
			LastAck:     stateStat[fmt.Sprintf("%02X", c.TCP_LAST_ACK)],
			Listen:      stateStat[fmt.Sprintf("%02X", c.TCP_LISTEN)],
			Closing:     stateStat[fmt.Sprintf("%02X", c.TCP_CLOSING)],
		},
	}, nil
}

func (self *TcpConnSampler) Sample() (*TcpConnSample, error) {
	return self.SampleByNetlink()
}
