/*
 * Author: fasion
 * Created time: 2019-08-28 15:24:37
 * Last Modified by: fasion
 * Last Modified time: 2019-09-02 19:02:55
 */

package dmi

import (
	"fmt"

	"github.com/digitalocean/go-smbios/smbios"
)

var (
	BadFormatError = fmt.Errorf("bad structure format")
)

type Structure struct {
	Raw *smbios.Structure
	data interface{}
}

func (self *Structure) Data() (interface{}) {
	return self.data
}

func ParseRawStructures(raws []*smbios.Structure) ([]*Structure, error) {
	results := make([]*Structure, 0)
	for _, raw := range raws {
		switch raw.Header.Type {
		case TypeSystemInformation:
			if data, err := ToSystemInformation(raw); err != nil {
				return nil, err
			} else {
				results = append(results, &Structure{
					Raw: raw,
					data: data,
				})
			}
		case TypeBiosInformation:
			if data, err := ToBiosInformation(raw); err != nil {
				return nil, err
			} else {
				results = append(results, &Structure{
					Raw: raw,
					data: data,
				})
			}
		}
	}

	return results, nil
}

func padZeros(data []byte, size int) ([]byte) {
	n := len(data)
	if n < size {
		return append(data, make([]byte, size-n)...)
	}
	return data
}

func lookupString(i uint8, strs []string) (string, bool) {
	if i > uint8(len(strs)) {
		return "", false
	}

	if i == 0 {
		return "", true
	}

	return strs[i-1], true
}
