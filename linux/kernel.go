/*
 * Author: fasion
 * Created time: 2019-08-07 14:04:16
 * Last Modified by: fasion
 * Last Modified time: 2019-08-23 14:09:42
 */

package linux

import (
	"strconv"
	"strings"

	_syscall "github.com/fasionchan/osutil-go/syscall"
)

type KernelVersionNumber struct {
	KernelVersion int
	MajorRevision int
	MinorRevision int
}

func ParseKernelVersionNumber(s string) (*KernelVersionNumber, error) {
	var err error
	v := KernelVersionNumber{}

	parts := strings.Split(s, "-")

	fields := strings.Split(parts[0], ".")

	v.KernelVersion, err = strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}

	if len(fields) > 1 {
		v.MajorRevision, err = strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
	}

	if len(fields) > 2 {
		v.MinorRevision, err = strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}
	}

	return &v, nil
}

func MustKernelVersionNumber(s string) (*KernelVersionNumber) {
	v, _ := ParseKernelVersionNumber(s)
	return v
}

func (self KernelVersionNumber) Before(other KernelVersionNumber) (bool) {
	return self.KernelVersion < other.KernelVersion ||
		self.MajorRevision < other.MajorRevision ||
		self.MinorRevision < other.MinorRevision
}

func (self KernelVersionNumber) After(other KernelVersionNumber) (bool) {
	return self.KernelVersion > other.KernelVersion ||
		self.MajorRevision > other.MajorRevision ||
		self.MinorRevision > other.MinorRevision
}

func FetchKernelVersionNumber() (*KernelVersionNumber, error) {
	info, err := _syscall.Uname()
	if err != nil {
		return nil, err
	}

	return ParseKernelVersionNumber(info.Release)
}

type KernelVersion struct {
	KernelVersion int
	MajorRevision int
	MinorRevision int
	AbiNumber int
	UploadNumber int
	Flavour string
}

func FetchKernelVersion() (*KernelVersion, error) {
	info, err := _syscall.Uname()
	if err != nil {
		return nil, err
	}

	return ParseKernelVersion(info.Release)
}

func MustKernelVersion(s string) (*KernelVersion) {
	v, _ := ParseKernelVersion(s)
	return v
}

func ParseKernelVersion(s string) (*KernelVersion, error) {
	var err error
	v := KernelVersion{}

	parts := strings.Split(s, "-")

	fields := strings.Split(parts[0], ".")

	v.KernelVersion, err = strconv.Atoi(fields[0])
	if err != nil {
		return nil, err
	}

	if len(fields) > 1 {
		v.MajorRevision, err = strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
	}

	if len(fields) > 2 {
		v.MinorRevision, err = strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}
	}

	if len(parts) > 1 {
		fields := strings.Split(parts[1], ".")
		v.AbiNumber, err = strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}

		if len(fields) > 1 {
			v.UploadNumber, err = strconv.Atoi(fields[1])
			if err != nil {
				return nil, err
			}
		}
	}

	if len(parts) > 2 {
		v.Flavour = parts[2]
	}

	return &v, nil
}

func (self KernelVersion) Before(other KernelVersion) (bool) {
	return self.KernelVersion < other.KernelVersion ||
		self.MajorRevision < other.MajorRevision ||
		self.MinorRevision < other.MinorRevision ||
		self.AbiNumber < other.AbiNumber ||
		self.UploadNumber < other.UploadNumber
}

func (self KernelVersion) After(other KernelVersion) (bool) {
	return self.KernelVersion > other.KernelVersion ||
		self.MajorRevision > other.MajorRevision ||
		self.MinorRevision > other.MinorRevision ||
		self.AbiNumber > other.AbiNumber ||
		self.UploadNumber > other.UploadNumber
}
