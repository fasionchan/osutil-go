/*
 * Author: fasion
 * Created time: 2019-09-02 13:53:39
 * Last Modified by: fasion
 * Last Modified time: 2019-09-10 15:25:26
 */

package sysfs

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

const (
	SysfsHwmonPath = "/sys/class/hwmon"
)

type CoreTemperatureMetric struct {
	Critical float64
	Max float64
	Current float64
	Alarm bool
}

type CoreTemperatureData struct {
	Label string
	Metric CoreTemperatureMetric
}

func FetchCoreTemperatureData(path string, index string) (*CoreTemperatureData, error) {
	kobject := NewKobject(path)

	label, err := kobject.FetchAttribute(fmt.Sprintf("temp%s_label", index))
	if label == "" {
		return nil, err
	}

	max, err := kobject.FetchUintAttribute(fmt.Sprintf("temp%s_max", index), 10, 32)
	if err != nil {
		return nil, err
	}

	crit, err := kobject.FetchUintAttribute(fmt.Sprintf("temp%s_crit", index), 10, 32)
	if err != nil {
		return nil, err
	}

	input, err := kobject.FetchUintAttributeDeadline(fmt.Sprintf("temp%s_input", index), 10, 32, time.Second)
	if err != nil {
		return nil, err
	}

	return &CoreTemperatureData{
		Label: strings.TrimSpace(label),
		Metric: CoreTemperatureMetric{
			Critical: float64(crit) / 1000,
			Max: float64(max) / 1000,
			Current: float64(input) / 1000,
		},
	}, nil
}

func FetchCoreTemperatureDatas() ([]*CoreTemperatureData, error) {
	files, err := ioutil.ReadDir(SysfsHwmonPath)
	if err != nil {
		return nil, err
	}

	datas := []*CoreTemperatureData{}
	for _, f := range files {
		driverPath := filepath.Join(SysfsHwmonPath, f.Name(), "device/driver")
		driverPath, err = filepath.EvalSymlinks(driverPath)
		if err != nil {
			return nil, err
		}

		if filepath.Base(driverPath) != "coretemp" {
			continue
		}

		devicePath := filepath.Join(SysfsHwmonPath, f.Name(), "device")

		files, err := ioutil.ReadDir(devicePath)
		if err != nil {
			return nil, err
		}


		for _, f := range files {
			name := f.Name()

			if !strings.HasSuffix(name, "_label") {
				continue
			}

			index := name[4:len(name)-6]

			data, err := FetchCoreTemperatureData(devicePath, index)
			if err != nil {
				return nil, err
			}
			if data == nil {
				break
			}

			datas = append(datas, data)
		}
	}

	return datas, nil
}
