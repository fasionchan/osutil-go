/*
 * Author: fasion
 * Created time: 2019-09-02 15:52:04
 * Last Modified by: fasion
 * Last Modified time: 2019-09-04 10:48:24
 */

package sysinfo

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"runtime"
	"strings"

	_syscall "github.com/fasionchan/osutil-go/syscall"
)

var _ = fmt.Println

const (
	LsbReleasePath = "/etc/lsb-release"
	CentOsReleasePath = "/etc/centos-release"
	RedhatReleasePath = "/etc/redhat-release"

	RedhatVersionRegexp = ` [0-9]+\.[0-9]+`
)

func FetchOs() (*Os, error) {
	os := Os{
		Type: runtime.GOOS,
		Arch: runtime.GOARCH,
	}

	if distro, err := FetchOsDistro(); err == nil {
		os.Distro = distro.Name
		os.DistroType = distro.Type
		os.DistroVersion = distro.Version

	}

	if uname, err := _syscall.Uname(); err == nil {
		os.KernelVersion = uname.Release
		os.KernelRevision = strings.Split(uname.Release, "-")[0]
	}

	return &os, nil
}

func FetchLsbRelease() (map[string]string, error) {
	datas := map[string]string{}

	flsb, err := os.Open(LsbReleasePath)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(flsb)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		fields := strings.SplitN(line, "=", 2)
		if len(fields) != 2 {
			continue
		}

		name := strings.TrimSpace(fields[0])
		value := strings.Trim(fields[1], `"`)

		datas[name] = strings.TrimSpace(value)
	}

	return datas, nil
}

func FetchOsDistro() (*OsDistro, error) {
	// ubuntu
	if lsb, err := FetchLsbRelease(); err == nil {
		return &OsDistro{
			Name: lsb["DISTRIB_DESCRIPTION"],
			Version: lsb["DISTRIB_RELEASE"],
			Type: strings.ToLower(lsb["DISTRIB_ID"]),
		}, nil
	}

	// centos
	if data, err := ioutil.ReadFile(CentOsReleasePath); err == nil {
		distro := strings.TrimSpace(string(data))

		var version string
		reg := regexp.MustCompile(RedhatVersionRegexp)
		matches := reg.FindAllString(distro, -1)
		if len(matches) == 1 {
			version = strings.TrimSpace(matches[0])
		}

		return &OsDistro{
			Name: distro,
			Type: "centos",
			Version: version,
		}, nil
	}

	// redhat
	if data, err := ioutil.ReadFile(RedhatReleasePath); err == nil {
		distro := strings.TrimSpace(string(data))

		var version string
		reg := regexp.MustCompile(RedhatVersionRegexp)
		matches := reg.FindAllString(distro, -1)
		if len(matches) == 1 {
			version = strings.TrimSpace(matches[0])
		}

		return &OsDistro{
			Name: distro,
			Type: "redhat",
			Version: version,
		}, nil
	}

	return &OsDistro{}, nil
}
