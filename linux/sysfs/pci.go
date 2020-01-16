/*
 * Author: fasion
 * Created time: 2019-08-27 08:12:28
 * Last Modified by: fasion
 * Last Modified time: 2019-08-27 10:12:52
 */

package sysfs

import (
	"fmt"

	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fasionchan/osutil-go/util/pcidb"
)

var _ = fmt.Println

const (
	PciDevicePath = "/sys/bus/pci/devices"
)

type PciProductInfo struct {
	VendorId string
	DeviceId string
	SubsystemVendorId string
	SubsystemDeviceId string

	Vendor string
	Device string
	Subsystem string
}

type PciClassInfo struct {
	ClassId string
	Class string
	Subclass string
	ProgIf string
}

type PciDevice struct {
	slot string
	path string

	querier *pcidb.PciIdsQuerier
}

func (self *PciDevice) readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(self.path, path))
}

func (self *PciDevice) Slot() (string) {
	return self.slot
}

func (self *PciDevice) ClassId() (string, error) {
	data, err := self.readFile("class")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data))[2:], nil
}

func (self *PciDevice) VendorId() (string, error) {
	data, err := self.readFile("vendor")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data))[2:], nil
}

func (self *PciDevice) DeviceId() (string, error) {
	data, err := self.readFile("device")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data))[2:], nil
}

func (self *PciDevice) SubsystemVendorId() (string, error) {
	data, err := self.readFile("subsystem_vendor")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data))[2:], nil
}

func (self *PciDevice) SubsystemDeviceId() (string, error) {
	data, err := self.readFile("subsystem_device")
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data))[2:], nil
}

func (self *PciDevice) ClassInfo() (*PciClassInfo, error) {
	classId, err := self.ClassId()
	if err != nil {
		return nil, err
	}

	class, subclass, progIf := self.querier.QueryClassInfoById(classId)

	return &PciClassInfo{
		ClassId: classId,
		Class: class,
		Subclass: subclass,
		ProgIf: progIf,
	}, nil
}

func (self *PciDevice) ProductInfo() (*PciProductInfo, error) {
	vendorId, err := self.VendorId()
	if err != nil {
		return nil, err
	}

	deviceId, err := self.DeviceId()
	if err != nil {
		return nil, err
	}

	subsystemVendorId, err := self.SubsystemVendorId()
	if err != nil {
		return nil, err
	}

	subsystemDeviceId, err := self.SubsystemDeviceId()
	if err != nil {
		return nil, err
	}

	var vendor, device, subsystem string
	if self.querier != nil {
		vendor, device, subsystem = self.querier.QueryProductInfo(vendorId, deviceId, subsystemVendorId, subsystemVendorId)
	}

	return &PciProductInfo{
		VendorId: vendorId,
		DeviceId: deviceId,
		SubsystemVendorId: subsystemVendorId,
		SubsystemDeviceId: subsystemDeviceId,

		Vendor: vendor,
		Device: device,
		Subsystem: subsystem,
	}, nil
}

func PciDeviceBySlot(slot string, querier *pcidb.PciIdsQuerier) (*PciDevice, error) {
	path := filepath.Join(PciDevicePath, slot)
	return &PciDevice{
		slot: slot,
		path: path,
		querier: querier,
	}, nil
}

func FetchPciDevices(querier *pcidb.PciIdsQuerier) ([]*PciDevice, error) {
	files, err := ioutil.ReadDir(PciDevicePath)
	if err != nil {
		return nil, err
	}

	devices := make([]*PciDevice, 0, len(files))
	for _, f := range files {
		device, err := PciDeviceBySlot(f.Name(), querier)
		if err != nil {
			return nil, err
		}

		devices = append(devices, device)
	}

	return devices, nil
}
