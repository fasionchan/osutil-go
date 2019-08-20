/*
 * Author: fasion
 * Created time: 2019-08-06 15:40:57
 * Last Modified by: fasion
 * Last Modified time: 2019-08-20 16:36:20
 */

package netlink

import (
	"encoding/binary"
	"fmt"

	"github.com/fasionchan/libgo/encoding"
)

var _ = fmt.Println

type NetlinkSockDiag struct {
	*NetlinkSocket
}

func NewNetlinkSockDiag() (*NetlinkSockDiag, error) {
	s, err := NewNetlinkSocket(NETLINK_SOCK_DIAG)
	if err != nil {
		return nil, err
	}

	return &NetlinkSockDiag{
		NetlinkSocket: s,
	}, nil
}

func (self *NetlinkSockDiag) InetDiagRequest(req InetDiagReqV2_c, bufsize int) (*InetDiagReceiver, error) {
	hdr := NlMsgHdr_c{
		Len: SizeofNlMsgHdr_c + SizeOfInetDiagReqV2_c,
		Type: SOCK_DIAG_BY_FAMILY,
		Flags: NLM_F_REQUEST | NLM_F_DUMP,
	}

	data, err := encoding.MarshalBinary(binary.LittleEndian, hdr, req)
	if err != nil {
		return nil, err
	}

	_, err = self.SendBinaryRequest(data, 0)
	if err != nil {
		return nil, err
	}

	return &InetDiagReceiver{
		socket: self.NetlinkSocket,
		bufsize: bufsize,
	}, nil
}
