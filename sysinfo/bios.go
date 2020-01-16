/*
 * Author: fasion
 * Created time: 2019-09-03 08:58:55
 * Last Modified by: fasion
 * Last Modified time: 2019-09-03 09:04:05
 */

package sysinfo

import (
	"fmt"
	"github.com/fasionchan/osutil-go/dmi"
)

func FetchBios() (*Bios, error) {
	structures, err := dmi.FetchStructures()
	if err != nil {
		return nil, err
	}

	var info *dmi.BiosInformation
	for _, structure := range structures {
		if structure.Raw.Header.Type == dmi.TypeBiosInformation {
			info, _ = structure.Data().(*dmi.BiosInformation)
			break
		}
	}

	if info == nil {
		return nil, fmt.Errorf("bios information bios missing")
	}

	bios := Bios{
		Vendor: info.Vendor,
		Version: info.Version,
		ReleaseDate: info.ReleaseDate,
		Release: info.Release,
		FirmwareRelease: info.FirmwareRelease,
	}

	return &bios, nil
}
