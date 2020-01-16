/*
 * Author: fasion
 * Created time: 2019-09-04 14:29:09
 * Last Modified by: fasion
 * Last Modified time: 2019-09-04 14:49:03
 */

package sysinfo

import (
	"github.com/fasionchan/osutil-go/linux/sysfs"
)

func FetchScsiDevices() ([]ScsiDevice, error) {
	devices, err := sysfs.FetchScsiDevices()
	if err != nil {
		return nil, err
	}

	result := []ScsiDevice{}
	for _, device := range devices {
		result = append(result, ScsiDevice{
			Address: device.Address,
			Type: device.Type.String(),
			Vendor: device.Vendor,
			Model: device.Model,
		})
	}

	return result, nil
}
