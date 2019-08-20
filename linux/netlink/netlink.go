/*
 * Author: fasion
 * Created time: 2019-08-06 15:33:44
 * Last Modified by: fasion
 * Last Modified time: 2019-08-20 16:35:40
 */

package netlink

import (
	"fmt"
	"syscall"
)

var _ = fmt.Println

func NetlinkSocketFd(proto int) (fd int, err error) {
	return syscall.Socket(AF_NETLINK, SOCK_DGRAM, proto)
}

type NetlinkSocket struct {
	fd int
}

func NewNetlinkSocket(proto int) (*NetlinkSocket, error) {
	fd, err := NetlinkSocketFd(proto)
	if err != nil {
		return nil, err
	}

	return &NetlinkSocket{
		fd: fd,
	}, nil
}

func (self *NetlinkSocket) SendBinaryRequest(req []byte, flags int) (int, error) {
	addr := syscall.SockaddrNetlink{
		Family: syscall.AF_NETLINK,
	}

	return syscall.SendmsgN(self.fd, req, nil, &addr, flags)
}

func (self *NetlinkSocket) ReceiveMessages(bufsize, flags int) ([]syscall.NetlinkMessage, error) {
	buf := make([]byte, bufsize)
	n, _, _, _, err := syscall.Recvmsg(self.fd, buf, nil, flags)
	if err != nil {
		return nil, err
	}

	return syscall.ParseNetlinkMessage(buf[:n])
}

func (self *NetlinkSocket) Close() (error) {
	fd := self.fd
	self.fd = -1
	return syscall.Close(fd)
}
