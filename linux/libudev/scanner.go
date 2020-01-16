/*
 * Author: fasion
 * Created time: 2019-08-29 10:03:47
 * Last Modified by: fasion
 * Last Modified time: 2019-08-29 10:43:38
 */

package libudev

import (
	// "bufio"
	"fmt"
	"os"
	"path/filepath"
)

const (
	SysfsDevicesPath = "/sys/devices"
)

type Scanner struct {
}

func NewScanner() (*Scanner) {
	return &Scanner{
	}
}

func (self *Scanner) Scan() ([]*Device, error) {
	devices := []*Device{}

	if err := filepath.Walk(SysfsDevicesPath, func(path string, info os.FileInfo, err error) (error) {
		if err != nil {
			fmt.Println(path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() != "uevent" {
			return nil
		}

		device, err := NewDevice(filepath.Dir(path))
		if err != nil {
			return err
		}
		if device == nil {
			return nil
		}

		devices = append(devices, device)

		return nil
	}); err != nil {
		return nil, err
	}

	return devices, nil
}
