/*
 * Author: fasion
 * Created time: 2019-09-02 16:37:46
 * Last Modified by: fasion
 * Last Modified time: 2019-09-03 14:20:37
 */

package sysinfo

import (
	"fmt"
	"strings"
	"github.com/fasionchan/osutil-go/dmi"
)

var VirtualVendors = []string{
	"innotek gmbh",
	"vmware",
	"openstack",
}

func FetchSystem() (*System, error) {
	structures, err := dmi.FetchStructures()
	if err != nil {
		return nil, err
	}

	var info *dmi.SystemInformation
	for _, structure := range structures {
		if structure.Raw.Header.Type == dmi.TypeSystemInformation {
			info, _ = structure.Data().(*dmi.SystemInformation)
			break
		}
	}

	if info == nil {
		return nil, fmt.Errorf("system information bios missing")
	}

	system := System{
		Vendor: info.Manufacturer,
		Name: info.ProductName,
		Version: info.Version,
		SerialNumber: info.SerialNumber,
		Uuid: info.Uuid.String(),
		SkuNumber: info.SkuNumber,
	}

	vendor := strings.ToLower(info.Manufacturer)
	for _, vv := range VirtualVendors {
		if strings.Contains(vendor, vv) {
			system.Virtual = true
		}
	}

	return &system, nil
}
