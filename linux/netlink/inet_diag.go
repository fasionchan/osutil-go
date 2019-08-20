/*
 * Author: fasion
 * Created time: 2019-08-06 15:41:04
 * Last Modified by: fasion
 * Last Modified time: 2019-08-20 16:36:29
 */

package netlink

import (
	"encoding/binary"
	"fmt"

	"github.com/fasionchan/libgo/encoding"
)

var _ = fmt.Println

type NetlinkInetDiag struct {
	*NetlinkSocket
}

func NewNetlinkInetDiag() (*NetlinkInetDiag, error) {
	s, err := NewNetlinkSocket(NETLINK_INET_DIAG)
	if err != nil {
		return nil, err
	}

	return &NetlinkInetDiag{
		NetlinkSocket: s,
	}, nil
}

func (self *NetlinkInetDiag) InetDiagRequest(req InetDiagReq_c, bufsize int) (*InetDiagReceiver, error) {
	hdr := NlMsgHdr_c{
		Len: SizeofNlMsgHdr_c + SizeOfInetDiagReq_c,
		Type: TCPDIAG_GETSOCK,
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


type InetDiagReceiver struct {
	bufsize int
	socket *NetlinkSocket
}

func (self InetDiagReceiver) Receive() ([]*InetDiagMsg_c, bool, error) {
	msgs, err := self.socket.ReceiveMessages(self.bufsize, 0)
	if err != nil {
		return nil, false, err
	}

	diags := make([]*InetDiagMsg_c, 0, len(msgs))

	for _, msg := range msgs {
		if (msg.Header.Type == NLMSG_DONE) {
			return diags, false, nil
		}

		if (msg.Header.Len < SizeofNlMsgHdr_c + SizeOfInetDiagMsg) {
			continue
		}

		diag := new(InetDiagMsg_c)
		err := encoding.UnmarshalBinary(binary.BigEndian, msg.Data, diag)
		if err != nil {
			return nil, false, err
		}

		diags = append(diags, diag)
	}

	return diags, true, nil
}
