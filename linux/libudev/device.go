/*
 * Author: fasion
 * Created time: 2019-08-29 10:05:39
 * Last Modified by: fasion
 * Last Modified time: 2019-08-29 13:37:56
 */

package libudev

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var _ = fmt.Println

type Device struct {
	SysfsPath string

	RawAttrs map[string][]byte
	Attrs map[string]string
	UeventVars map[string]string

	Parent *Device
	Children []*Device
}

func ScanDeviceRawAttrs(sysfsPath string) (map[string][]byte, error) {
	files, err := ioutil.ReadDir(sysfsPath)
	if err != nil {
		return nil, err
	}

	attrs := map[string][]byte{}
	for _, f := range files {
		// skip directory
		if !f.Mode().IsRegular() {
			continue
		}

		// skip not readable
		if f.Mode().Perm() & 0400 == 0 {
			continue
		}

		name := f.Name()
		data, err := ioutil.ReadFile(filepath.Join(sysfsPath, name))
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}

			continue
		}

		attrs[name] = data
	}

	return attrs, nil
}

func ScanDeviceAttrs(sysfsPath string) (map[string]string, error) {
	rawAttrs, err := ScanDeviceRawAttrs(sysfsPath)
	if err != nil {
		return nil, err
	}

	return DeviceAttrsFromRaw(rawAttrs), nil
}

func DeviceAttrsFromRaw(rawAttrs map[string][]byte) (map[string]string) {
	attrs := map[string]string{}

	for name, value := range rawAttrs {
		attrs[name] = string(value)
	}

	return attrs
}

func ParseUeventVars(data string) (map[string]string) {
	vars := map[string]string{}
	for _, line := range strings.Split(strings.TrimSpace(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Split(line, "=")
		if len(fields) != 2 {
			continue
		}

		vars[fields[0]] = fields[1]
	}
	return vars
}

func FetchUdevInfo(dev string) (error) {
	return nil
}

func NewDevice(sysfsPath string) (*Device, error) {
	rawAttrs, err := ScanDeviceRawAttrs(sysfsPath)
	if err != nil {
		return nil, err
	}

	attrs := DeviceAttrsFromRaw(rawAttrs)

	dev, ok := attrs["dev"]
	if !ok {
		return nil, nil
	}

	FetchUdevInfo(dev)

	ueventVars := ParseUeventVars(attrs["uevent"])

	device := Device{
		SysfsPath: sysfsPath,
		RawAttrs: rawAttrs,
		Attrs: attrs,
		UeventVars: ueventVars,
	}

	return &device, nil
}
