/*
 * Author: fasion
 * Created time: 2020-03-27 12:16:05
 * Last Modified by: fasion
 * Last Modified time: 2020-03-27 17:38:18
 */

package linux

import (
	"encoding/hex"
	"fmt"
	"net"
	"syscall"

	"github.com/fasionchan/libgo/arch"
	"github.com/fasionchan/osutil-go/linux/c"
	"github.com/fasionchan/osutil-go/linux/netlink"
	"github.com/fasionchan/osutil-go/linux/procfs"
)

const (
	DefaultNetlinkDiagBufferSize = 1024 * 1024
	DefaultProcfsTcpBufferSize   = 1024 * 1024
)

type TcpDiag struct {
	Family uint8

	LocalAddr net.IP
	LocalPort int

	PeerAddr net.IP
	PeerPort int

	State uint8
}

func (diag *TcpDiag) LocalAddrHex() string {
	return hex.EncodeToString([]byte(diag.LocalAddr))
}

func (diag *TcpDiag) PeerAddrHex() string {
	return hex.EncodeToString([]byte(diag.PeerAddr))
}

func (diag *TcpDiag) StateName() string {
	if int(diag.State) >= len(c.TcpStateNames) {
		return ""
	}

	return c.TcpStateNames[int(diag.State)]
}

type TcpDiagHandler func(*TcpDiag) bool

func IpFromUint32(family uint8, data [4]uint32) net.IP {
	var buf [16]byte
	arch.NativeEndian.PutUint32(buf[0:], data[0])
	arch.NativeEndian.PutUint32(buf[4:], data[1])
	arch.NativeEndian.PutUint32(buf[8:], data[2])
	arch.NativeEndian.PutUint32(buf[12:], data[3])

	if family == syscall.AF_INET {
		return net.IP(buf[:4])
	} else {
		return net.IP(buf[:])
	}
}

func ListTcpxDiagFromReceiver(receiver *netlink.InetDiagReceiver, family uint8, states uint32, handler TcpDiagHandler) error {
	for {
		diags, more, err := receiver.Receive()
		if err != nil {
			return err
		}

		for _, diag := range diags {
			if diag.Family != family {
				continue
			}

			if (1<<diag.State)&states == 0 {
				continue
			}

			if !handler(&TcpDiag{
				Family: diag.Family,

				LocalAddr: IpFromUint32(diag.Family, diag.Id.Src),
				LocalPort: int(arch.Ntohs(diag.Id.Sport)),

				PeerAddr: IpFromUint32(diag.Family, diag.Id.Dst),
				PeerPort: int(arch.Ntohs(diag.Id.Dport)),

				State: diag.State,
			}) {
				return nil
			}
		}

		if !more {
			break
		}
	}

	return nil
}

func ListTcpxDiagFromSockDiag(family uint8, states uint32, handler TcpDiagHandler, bufferSize int) error {
	if bufferSize == 0 {
		bufferSize = DefaultNetlinkDiagBufferSize
	}

	sockDiag, err := netlink.NewNetlinkSockDiag()
	if err != nil {
		return err
	}
	defer sockDiag.Close()

	// find tcp diags
	receiver, err := sockDiag.InetDiagRequest(netlink.InetDiagReqV2_c{
		Sdiag_family:   family,
		Sdiag_protocol: netlink.IPPROTO_TCP,
		Idiag_states:   uint32(states),
	}, bufferSize)
	if err != nil {
		return err
	}

	return ListTcpxDiagFromReceiver(receiver, family, states, handler)
}

func ListTcpxDiagFromInetDiag(family uint8, states uint32, handler TcpDiagHandler, bufferSize int) error {
	if bufferSize == 0 {
		bufferSize = DefaultNetlinkDiagBufferSize
	}

	inetDiag, err := netlink.NewNetlinkInetDiag()
	if err != nil {
		return err
	}
	defer inetDiag.Close()

	// find tcp diags
	receiver, err := inetDiag.InetDiagRequest(netlink.InetDiagReq_c{
		Family: family,
		States: uint32(states),
	}, bufferSize)
	if err != nil {
		return err
	}

	return ListTcpxDiagFromReceiver(receiver, family, states, handler)
}

func ListTcpxDiagFromNetlink(family uint8, states uint32, handler TcpDiagHandler, bufferSize int) error {
	kv, err := FetchKernelVersionNumber()
	if err != nil {
		return err
	}

	if !kv.Before(*MustKernelVersionNumber("3.3")) {
		return ListTcpxDiagFromSockDiag(family, states, handler, bufferSize)
	}

	if !kv.Before(*MustKernelVersionNumber("2.6.14")) {
		return ListTcpxDiagFromInetDiag(family, states, handler, bufferSize)
	}

	return fmt.Errorf("netlink is not supported by kernel")
}

func ListTcpxDiagFromProcfs(family uint8, states uint32, handler TcpDiagHandler, bufferSize int) error {
	if bufferSize == 0 {
		bufferSize = DefaultProcfsTcpBufferSize
	}

	scanner, err := procfs.NewNetTcpxScanner(family, bufferSize)
	if err != nil {
		return err
	}
	defer scanner.Close()

	for scanner.Scan() {
		_, fields, err := scanner.Record()
		if err != nil {
			return err
		}

		state, err := fields.State()
		if err != nil {
			return err
		}

		if (1<<state)&states == 0 {
			continue
		}

		laddr, lport, err := fields.Local()
		if err != nil {
			return err
		}

		raddr, rport, err := fields.Remote()
		if err != nil {
			return err
		}

		if !handler(&TcpDiag{
			Family: family,

			LocalAddr: laddr,
			LocalPort: lport,

			PeerAddr: raddr,
			PeerPort: rport,

			State: state,
		}) {
			return nil
		}
	}

	return nil
}

func ListTcpxDiag(family uint8, states uint32, handler TcpDiagHandler, bufferSize int) error {
	if err := ListTcpxDiagFromNetlink(family, states, handler, bufferSize); err != nil {
		return err
	}

	return ListTcpxDiagFromProcfs(family, states, handler, bufferSize)
}

func ListTcp4Diag(states uint32, handler TcpDiagHandler, bufferSize int) error {
	return ListTcpxDiag(syscall.AF_INET, states, handler, bufferSize)
}

func ListTcp6Diag(states uint32, handler TcpDiagHandler, bufferSize int) error {
	return ListTcpxDiag(syscall.AF_INET6, states, handler, bufferSize)
}

func ListTcpDiag(states uint32, handler TcpDiagHandler, bufferSize int) error {
	if err := ListTcpxDiag(syscall.AF_INET, states, handler, bufferSize); err != nil {
		return nil
	}

	return ListTcpxDiag(syscall.AF_INET6, states, handler, bufferSize)
}
