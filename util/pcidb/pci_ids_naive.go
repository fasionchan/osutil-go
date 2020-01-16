/*
 * Author: fasion
 * Created time: 2019-08-26 17:36:36
 * Last Modified by: fasion
 * Last Modified time: 2019-08-27 09:30:12
 */

package pcidb

import (
	"bufio"
	"io"
	"strings"
)

type PciId = uint16
type PciClassId = uint8

type IdKey = [4]string
type ClassKey = [3]string

type PciIdsQuerier struct {
	ids map[IdKey]string
	classes map[ClassKey]string
}

func (self *PciIdsQuerier) QueryProductInfo(vendorId, deviceId, subvendorId, subdeviceId string) (vendor, device, subsystem string) {
	key := IdKey{
		vendorId,
	}
	vendor = self.ids[key]

	if deviceId != "" {
		key[1] = deviceId
		device = self.ids[key]
	}

	if subvendorId != "" {
		key[2] = subvendorId
		key[3] = subdeviceId
		subsystem = self.ids[key]
	}

	return
}

func (self *PciIdsQuerier) QueryClassInfo(classId, subclassId, progIfId string) (class, subclass, progIf string) {
	key := ClassKey{
		classId,
	}
	class = self.classes[key]

	if subclassId != "" {
		key[1] = subclassId
		subclass = self.classes[key]
	}

	if progIfId != "" {
		key[2] = progIfId
		progIf = self.classes[key]
	}

	return
}

func (self *PciIdsQuerier) QueryClassInfoById(id string) (string, string, string) {
	if len(id) != 6 {
		return "", "", ""
	}
	return self.QueryClassInfo(id[:2], id[2:4], id[4:])
}

func PciIdsQuerierFromScanner(scanner *bufio.Scanner) (*PciIdsQuerier, error) {
	var vendor, device string
	var class, subclass string

	ids := make(map[[4]string]string)
	classes := make(map[[3]string]string)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, "\t\t") {
			if vendor != "" {
				subvendor := line[2:6]
				subdevice := line[7:11]
				name := strings.TrimSpace(line[13:])

				key := [4]string{
					vendor,
					device,
					subvendor,
					subdevice,
				}

				ids[key] = name
			} else if class != "" {
				progIf := line[2:4]
				name := strings.TrimSpace(line[6:])

				key := [3]string{
					class,
					subclass,
					progIf,
				}

				classes[key] = name
			}
		} else if strings.HasPrefix(line, "\t") {
			if vendor != "" {
				device := line[1:5]
				name := strings.TrimSpace(line[7:])

				key := [4]string{
					vendor,
					device,
				}

				ids[key] = name
			} else if class != "" {
				subclass := line[1:3]
				name := strings.TrimSpace(line[5:])

				key := [3]string{
					class,
					subclass,
				}

				classes[key] = name
			}
		} else if strings.HasPrefix(line, "C ") {
			vendor = ""
			class = line[2:4]
			name := line[6:]

			key := [3]string{
				class,
			}

			classes[key] = name
		} else {
			class = ""
			vendor = line[0:4]
			name := line[6:]

			key := [4]string{
				vendor,
			}

			ids[key] = name
		}
	}

	return &PciIdsQuerier{
		classes: classes,
		ids: ids,
	}, nil
}

func PciIdsQuerierFromReader(reader io.Reader) (*PciIdsQuerier, error) {
	return PciIdsQuerierFromScanner(bufio.NewScanner(reader))
}
