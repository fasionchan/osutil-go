/*
 * Author: fasion
 * Created time: 2019-08-26 17:04:19
 * Last Modified by: fasion
 * Last Modified time: 2019-08-26 17:33:34
 */

package pcidb

type PciVendor struct {
	id string
	name string
	devices map[string]*PciDevice
}

type PciDevice struct {
	id string
	name string
	subsytems []*PciSubsystem
}

type PciSubsystem struct {
	vendorId string
	deviceId string
	name string
}

type PciClass struct {
	id string
	name string
	subclasses map[string]PciSubclass
}

type PciSubclass struct {
	id string
	name string
}
