/*
 * Author: fasion
 * Created time: 2019-08-29 14:29:27
 * Last Modified by: fasion
 * Last Modified time: 2019-08-29 16:32:10
 */

package sysfs

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fasionchan/osutil-go/linux/c"
)

const (
	ScsiDeviceRootPath = "/sys/class/scsi_device"
)

const (
	ScsiTypeDisk = ScsiType(c.TYPE_DISK)
	ScsiTypeRom = ScsiType(c.TYPE_ROM)
	ScsiTypeRaid = ScsiType(0x0c)
)

var (
	ScsiTypeNames = [0x100]string{
		c.TYPE_DISK: "Direct-access block device",
		c.TYPE_TAPE: "Sequential-access device",
		c.TYPE_ROM: "CD/DVD-ROM device",
		0x0c: "Storage array controller device",
	}

	ShortScsiTypeNames = [0x100]string{
		c.TYPE_DISK: "disk",
		c.TYPE_TAPE: "tape",
		c.TYPE_ROM: "cd/dvd",
		0x0c: "raid",
	}
)

type ScsiType uint8

func (self ScsiType) ShortString() (string) {
	return ShortScsiTypeNames[int(self)]
}

func (self ScsiType) String() (string) {
	return ScsiTypeNames[int(self)]
}

type ScsiDevice struct {
	Address string

	Type ScsiType
	Vendor string
	Model string
}

func FetchScsiDevices() ([]*ScsiDevice, error) {
	files, err := ioutil.ReadDir(ScsiDeviceRootPath)
	if err != nil {
		return nil, err
	}

	devices := []*ScsiDevice{}

	for _, f := range files {
		name := f.Name()

		path := filepath.Join(ScsiDeviceRootPath, name, "device")

		path, err := filepath.EvalSymlinks(path)
		if err != nil {
			continue
		}

		kobject := NewKobject(path)

		dt, err := kobject.FetchUintAttribute("type", 10, 8)
		if err != nil {
			continue
		}

		device := &ScsiDevice{
			Address: name,
			Type: ScsiType(dt),
			Vendor: strings.TrimSpace(kobject.SmartAttribute("vendor")),
			Model: strings.TrimSpace(kobject.SmartAttribute("model")),
		}

		devices = append(devices, device)
	}

	return devices, nil
}
