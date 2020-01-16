/*
 * Author: fasion
 * Created time: 2019-05-17 17:47:42
 * Last Modified by: fasion
 * Last Modified time: 2019-12-26 15:51:24
 */

package sysinfo

import (
	"crypto/md5"

	"github.com/fasionchan/libgo/sorting"
)

func GetFingerprint() ([]byte, []NetworkInterfaceCard, error) {
	h := md5.New()

	nics, err := NetworkInterfaceCards()
	if err != nil {
		return nil, nil, err
	}

	nicsRefered := make([]NetworkInterfaceCard, 0, len(nics))
	macs := make([][]byte, 0, len(nics))
	for _, nic := range nics {
		if nic.Virtual {
			continue
		}

		nicsRefered = append(nicsRefered, nic)
		macs = append(macs, nic.HardwareAddr)
	}

	macs = sorting.ByteSlices(macs)
	for _, mac := range macs {
		h.Write(mac)
	}

	return h.Sum(nil), nicsRefered, nil
}
