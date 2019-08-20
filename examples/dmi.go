/*
 * Author: fasion
 * Created time: 2019-05-22 08:51:51
 * Last Modified by: fasion
 * Last Modified time: 2019-05-22 08:53:16
 */

package main

import (
	"fmt"

	"github.com/fasionchan/osutil-go/dmi"
)

func main() {
	fmt.Println(dmi.FetchRawDMI())
}