/*
 * Author: fasion
 * Created time: 2019-08-13 18:14:47
 * Last Modified by: fasion
 * Last Modified time: 2019-08-23 13:34:36
 */

package sysinfo

import (
	"net"
)

type NetworkInterfaceCard struct {
    Index           int
    Name            string
	MTU             int
	Virtual 		bool
	HardwareAddr 	net.HardwareAddr
}
