/*
 * Author: fasion
 * Created time: 2019-09-03 13:46:11
 * Last Modified by: fasion
 * Last Modified time: 2019-09-03 14:02:27
 */

package wmi

import (
	"github.com/StackExchange/wmi"
)

type Win32_OperatingSystem struct {
	Caption string
	ProductType int
	OperatingSystemSKU int
	Version string
}

func FetchWin32_OperatingSystem() (*Win32_OperatingSystem, error) {
	records := make([]*Win32_OperatingSystem, 0)

	err := wmi.QueryNamespace("SELECT * FROM Win32_OperatingSystem", &records, `\root\CIMV2`)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, nil
	}

	return records[0], nil
}
