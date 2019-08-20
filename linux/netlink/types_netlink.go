// +build ignore

/**
 * FileName:   types_netlink.go
 * Author:     Fasion Chan
 * @contact:   fasionchan@gmail.com
 * @version:   $Id$
 *
 * Description:
 *
 * Changelog:
 *
 **/

package netlink

/*
#include <sys/socket.h>
#include <linux/netlink.h>
*/
import "C"

const (
	AF_NETLINK          = C.AF_NETLINK
	SOCK_DGRAM          = C.SOCK_DGRAM
	SOCK_RAW            = C.SOCK_RAW

	NETLINK_ROUTE       = C.NETLINK_ROUTE
	//NETLINK_W1          = C.NETLINK_W1
	NETLINK_USERSOCK    = C.NETLINK_USERSOCK
	NETLINK_FIREWALL    = C.NETLINK_FIREWALL
	NETLINK_INET_DIAG   = C.NETLINK_INET_DIAG
	NETLINK_SOCK_DIAG   = C.NETLINK_SOCK_DIAG
	NETLINK_NFLOG       = C.NETLINK_NFLOG

	NLMSG_NOOP          = C.NLMSG_NOOP
	NLMSG_ERROR         = C.NLMSG_ERROR
	NLMSG_DONE          = C.NLMSG_DONE

	NLM_F_REQUEST       = C.NLM_F_REQUEST
	NLM_F_MULTI         = C.NLM_F_MULTI
	NLM_F_ACK           = C.NLM_F_ACK
	NLM_F_ECHO          = C.NLM_F_ECHO

	NLM_F_ROOT          = C.NLM_F_ROOT
	NLM_F_MATCH         = C.NLM_F_MATCH
	NLM_F_ATOMIC        = C.NLM_F_ATOMIC
	NLM_F_DUMP          = C.NLM_F_DUMP
)

type NlMsgHdr_c         C.struct_nlmsghdr
type NlMsgErr_c         C.struct_nlmsgerr
type SockAddrNl_c       C.struct_sockaddr_nl

const (
	SizeofNlMsgHdr_c        = C.sizeof_struct_nlmsghdr
	SizeofNlMsgErr_c        = C.sizeof_struct_nlmsgerr
	SizeofSockaddrNl_c      = C.sizeof_struct_sockaddr_nl
)
