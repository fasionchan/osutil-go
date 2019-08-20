// +build linux

/*
 * Author: fasion
 * Created time: 2019-05-22 08:40:43
 * Last Modified by: fasion
 * Last Modified time: 2019-05-22 08:54:39
 */

package dmi

import (
	"io/ioutil"
	"os"
)

const (
	DMI_PATH = "/sys/firmware/dmi/tables/DMI"
)

func FetchRawDMI() ([]byte, error) {
	_, err := os.Stat(DMI_PATH)
	if err == nil {
		return ioutil.ReadFile(DMI_PATH)
	}

	return nil, nil
}