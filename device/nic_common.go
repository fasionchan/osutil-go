/*
 * Author: fasion
 * Created time: 2019-08-13 18:14:47
 * Last Modified by: fasion
 * Last Modified time: 2019-08-14 10:38:56
 */

package device

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
