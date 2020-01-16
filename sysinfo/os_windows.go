/*
 * Author: fasion
 * Created time: 2019-09-03 10:56:19
 * Last Modified by: fasion
 * Last Modified time: 2019-09-04 10:48:51
 */

package sysinfo

import (
	"runtime"
	"strings"

	"github.com/fasionchan/osutil-go/windows/wmi"
)

const (
	KernelRevisionRegexp = ` [0-9]+\.[0-9]+`
)

var Sku2Name = map[int]string{
	10: "Enterprise Server Edition",
}

func FetchOs() (*Os, error) {
	winos, err := wmi.FetchWin32_OperatingSystem()
	if err != nil {
		return nil, err
	}

	os := Os{
		Type: runtime.GOOS,
		Arch: runtime.GOARCH,
		Distro: strings.TrimSpace(winos.Caption),
		DistroType: Sku2Name[winos.OperatingSystemSKU],
		KernelVersion: winos.Version,
		KernelRevision: strings.Join(strings.Split(winos.Version, ".")[:2], "."),
	}

	return &os, nil
}
