/*
 * Author: fasion
 * Created time: 2019-08-26 14:43:29
 * Last Modified by: fasion
 * Last Modified time: 2019-08-26 15:08:03
 */

package linux

import (
	"syscall"
)

type Ethtool struct {
	socket int
}

func NewEthtool() (*Ethtool, error) {
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		return nil, err
	}

	return &Ethtool{
		socket: socket,
	}, nil
}

func (self *Ethtool) Close() (err error) {
	return syscall.Close(self.socket)
}
