/*
 * Author: fasion
 * Created time: 2019-08-19 13:45:59
 * Last Modified by: fasion
 * Last Modified time: 2019-09-10 19:25:56
 */

package perf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
	"unsafe"

	"github.com/fasionchan/libgo/arch"
)

const (
	MIB_TCP_STATE_CLOSED = 1
	MIB_TCP_STATE_LISTEN = 2
	MIB_TCP_STATE_SYN_SENT = 3
	MIB_TCP_STATE_SYN_RCVD = 4
	MIB_TCP_STATE_ESTAB = 5
	MIB_TCP_STATE_FIN_WAIT1 = 6
	MIB_TCP_STATE_FIN_WAIT2 = 7
	MIB_TCP_STATE_CLOSE_WAIT = 8
	MIB_TCP_STATE_CLOSING = 9
	MIB_TCP_STATE_LAST_ACK = 10
	MIB_TCP_STATE_TIME_WAIT = 11
	MIB_TCP_STATE_DELETE_TCB = 12
)

var _ = fmt.Println

const (
	NO_ERROR = 0
	ERROR_INSUFFICIENT_BUFFER = 122
)

var (
	procGetTcpTable = syscall.NewLazyDLL("iphlpapi.dll").NewProc("GetTcpTable")
	procGetTcp6Table = syscall.NewLazyDLL("iphlpapi.dll").NewProc("GetTcp6Table")
)

type MIB_TCPROW struct {
	State uint32
	LocalAddr [4]byte
	LocalPort uint32
	RemoteAddr [4]byte
	RemotePort uint32
}

func FetchTcpTable() ([]*MIB_TCPROW, error) {
	if err := procGetTcpTable.Find(); err != nil {
		return nil, err
	}

	var size uint32 = 0
	r1, _, err := procGetTcpTable.Call(
		uintptr(unsafe.Pointer(nil)),
		uintptr(unsafe.Pointer(&size)),
		uintptr(0),
	)
	if r1 != ERROR_INSUFFICIENT_BUFFER {
		return nil, err
	}

	buffer := make([]byte, size)

	r1, _, err = procGetTcpTable.Call(
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
		uintptr(0),
	)
	if r1 != NO_ERROR {
		return nil, err
	}

	reader := bytes.NewReader(buffer[:size])

	var n uint32
	if err := binary.Read(reader, arch.NativeEndian, &n); err != nil {
		return nil, err
	}

	rows := make([]*MIB_TCPROW, 0, n)
	for ; n > 0; n-- {
		r := new(MIB_TCPROW)
		if err := binary.Read(reader, arch.NativeEndian, r); err != nil {
			return nil, err
		}

		rows = append(rows, r)
	}

	return rows, nil
}
