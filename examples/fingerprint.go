/*
 * Author: fasion
 * Created time: 2019-05-21 17:30:51
 * Last Modified by: fasion
 * Last Modified time: 2019-05-21 17:33:31
 */

package main

import (
	"fmt"
	"github.com/fasionchan/osutil-go/device"
)

func main() {
	fp, err := device.GetFingerPrint()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%x\n", fp)
}