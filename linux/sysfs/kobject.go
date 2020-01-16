/*
 * Author: fasion
 * Created time: 2019-08-29 14:51:32
 * Last Modified by: fasion
 * Last Modified time: 2019-09-10 15:24:17
 */

package sysfs

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Kobject struct {
	SysfsPath string
}

func NewKobject(sysfsPath string) (*Kobject) {
	return &Kobject{
		SysfsPath: sysfsPath,
	}
}

func (self *Kobject) FetchAttributeRaw(name string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(self.SysfsPath, name))
}

func (self *Kobject) FetchAttributeRawDeadline(name string, d time.Duration) ([]byte, error) {
	path := filepath.Join(self.SysfsPath, name)

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	f.SetReadDeadline(time.Now().Add(d))

	var buf bytes.Buffer
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (self *Kobject) AttributeRaw(name string) ([]byte) {
	data, _ := self.FetchAttributeRaw(name)
	return data
}

func (self *Kobject) FetchAttribute(name string) (string, error) {
	raw, err := ioutil.ReadFile(filepath.Join(self.SysfsPath, name))
	if err != nil {
		return "", err
	}
	return string(raw), err
}

func (self *Kobject) FetchAttributeDeadline(name string, d time.Duration) (string, error) {
	raw, err := self.FetchAttributeRawDeadline(name, d)
	if err != nil {
		return "", err
	}
	return string(raw), err
}

func (self *Kobject) FetchUintAttribute(name string, base int, bits int) (uint64, error) {
	data, err := self.FetchAttribute(name)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(strings.TrimSpace(data), base, bits)
}

func (self *Kobject) FetchUintAttributeDeadline(name string, base, bits int, d time.Duration) (uint64, error) {
	data, err := self.FetchAttributeDeadline(name, d)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(strings.TrimSpace(data), base, bits)
}

func (self *Kobject) Attribute(name string) (string) {
	data, _ := self.FetchAttribute(name)
	return data
}

func (self *Kobject) SmartAttribute(name string) (string) {
	data, _ := self.FetchAttribute(name)
	return strings.TrimSpace(data)
}
