/*
 * Author: fasion
 * Created time: 2019-08-23 09:34:47
 * Last Modified by: fasion
 * Last Modified time: 2019-09-18 16:59:08
 */

package sysinfo

import (
	"errors"
)

var NotImplementedError = errors.New("not implemented")

type Bios struct {
	Vendor string `bson:"Vendor"`
	Version string `bson:"Version"`
	ReleaseDate string `bson:"ReleaseDate"`
	Release string `bson:"Release"`
	FirmwareRelease string `bson:"FirmwareRelease"`
}

type System struct {
	Vendor string `bson:"Vendor"`
	Name string `bson:"Name"`
	Version string `bson:"Version"`

	SerialNumber string `bson:"SerialNumber"`
	Uuid string `bson:"Uuid"`
	SkuNumber string `bson:"SkuNumber"`

	Virtual bool `bson:"Virtual"`
}

type Os struct {
	Type string `bson:"Type"`
	Arch string `bson:"Arch"`

	Distro string `bson:"Distro"`
	DistroType string `bson:"DistroType"`
	DistroVersion string `bson:"DistroVersion"`

	KernelVersion string `bson:"KernelVersion"`
	KernelRevision string `bson:"KernelRevision"`
}

type OsDistro struct {
	Name string `bson:"Name"`
	Type string `bson:"Type"`
	Version string `bson:"Version"`
}

type ScsiDevice struct {
    Address string `bson:"Address"`

    Type string `bson:"Type"`
    Vendor string `bson:"Vendor"`
    Model string `bson:"Model"`
}

type PythonEnvironment struct {
	Installed bool
	Executable string
	Version string
}

type SysInfo struct {
	System System
	Bios Bios
	Os Os
	ScsiDevices []ScsiDevice

	Python2 PythonEnvironment
	Python3 PythonEnvironment
}
