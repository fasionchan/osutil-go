// +build ignore

/*
 * Author: fasion
 * Created time: 2019-07-12 16:48:51
 * Last Modified by: fasion
 * Last Modified time: 2019-08-13 10:00:54
 */

package netlink

/*
#include <linux/sock_diag.h>
#include <linux/inet_diag.h>
#include <netinet/in.h>
*/
import "C"

const (
    SOCK_DIAG_BY_FAMILY = C.SOCK_DIAG_BY_FAMILY

    INET_DIAG_MEMINFO   = C.INET_DIAG_MEMINFO
    INET_DIAG_INFO      = C.INET_DIAG_INFO
    INET_DIAG_CONG      = C.INET_DIAG_CONG
    INET_DIAG_TOS       = C.INET_DIAG_TOS
    INET_DIAG_TCLASS    = C.INET_DIAG_TCLASS
    INET_DIAG_SKMEMINFO = C.INET_DIAG_SKMEMINFO

    IPPROTO_TCP = C.IPPROTO_TCP

    TCPDIAG_GETSOCK = C.TCPDIAG_GETSOCK
)

type InetDiagReq_c          C.struct_inet_diag_req
type InetDiagMsg_c          C.struct_inet_diag_msg
type InetDiagSockId_c       C.struct_inet_diag_sockid
type InetDiagReqV2_c        C.struct_inet_diag_req_v2

const (
    SizeOfInetDiagReq_c = C.sizeof_struct_inet_diag_req
    SizeOfInetDiagMsg = C.sizeof_struct_inet_diag_msg
    SizeOfInetDiagReqV2_c = C.sizeof_struct_inet_diag_req_v2
)
