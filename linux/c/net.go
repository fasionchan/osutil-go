/*
 * Author: fasion
 * Created time: 2019-08-07 10:21:08
 * Last Modified by: fasion
 * Last Modified time: 2019-08-07 10:34:58
 */

package c

var TcpStateNames = []string{
	TCP_ESTABLISHED: "Established",
	TCP_SYN_SENT: "SynSent",
	TCP_SYN_RECV: "SynRecv",
	TCP_FIN_WAIT1: "FinWait1",
	TCP_FIN_WAIT2: "FinWait2",
	TCP_TIME_WAIT: "TimeWait",
	TCP_CLOSE: "Close",
	TCP_CLOSE_WAIT: "CloseWait",
	TCP_LAST_ACK: "LastAck",
	TCP_LISTEN: "Listen",
	TCP_CLOSING: "Closing",
}

type TcpState uint8

func (self TcpState) String() (string) {
	if int(self) < len(TcpStateNames) {
		return TcpStateNames[self]
	}
	return "Illegal"
}
