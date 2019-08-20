/*
 * Author: fasion
 * Created time: 2019-05-17 18:04:59
 * Last Modified by: fasion
 * Last Modified time: 2019-05-21 17:21:13
 */

package main

import (
	"fmt"
	"github.com/fasionchan/osutil-go/device"
)

func main() {
	nics, err := device.NetworkInterfaceCards()
	if err != nil {
		fmt.Println(err)
	}

	for _, nic := range nics {
		fmt.Println(nic)
	}

	fmt.Println(device.GetFingerPrint())
}