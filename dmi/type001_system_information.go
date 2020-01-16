/*
 * Author: fasion
 * Created time: 2019-08-28 17:07:04
 * Last Modified by: fasion
 * Last Modified time: 2019-09-10 19:25:04
 */

package dmi

import (
	"encoding/binary"
	"bytes"
	"fmt"

	"github.com/digitalocean/go-smbios/smbios"
	"github.com/google/uuid"
	"github.com/fasionchan/libgo/arch"
)

var _ = fmt.Println

const (
	TypeSystemInformation = 1
)

type SystemInformationStructure struct {
	Manufacturer uint8
	ProductName uint8
	Version uint8
	SerialNumber uint8
	Uuid [16]byte
	WakeupType uint8
	SkuNumber uint8
	Family uint8
}

type SystemInformation struct {
	Manufacturer string
	ProductName string
	Version string
	SerialNumber string
	Uuid uuid.UUID
	WakeupType uint8
	SkuNumber string
	Family string
}

func ToSystemInformation(raw *smbios.Structure) (*SystemInformation, error) {
	sl := len(raw.Formatted) + 4
	if sl < 0x08 {
		return nil, BadFormatError
	}

	var structure SystemInformationStructure
	reader := bytes.NewReader(padZeros(raw.Formatted, 0x1b - 4))
	if err := binary.Read(reader, arch.NativeEndian, &structure); err != nil {
		return nil, err
	}

	data := SystemInformation{}

	var ok bool
	if data.Manufacturer, ok = lookupString(structure.Manufacturer, raw.Strings); !ok {
		return nil, BadFormatError
	}
	if data.ProductName, ok = lookupString(structure.ProductName, raw.Strings); !ok {
		return nil, BadFormatError
	}
	if data.Version, ok = lookupString(structure.Version, raw.Strings); !ok {
		return nil, BadFormatError
	}
	if data.SerialNumber, ok = lookupString(structure.SerialNumber, raw.Strings); !ok {
		return nil, BadFormatError
	}

	if sl < 0x19 {
		return nil, nil
	}

	data.Uuid = uuid.UUID(structure.Uuid)

	if sl < 0x1b {
		return nil, nil
	}

	if data.SkuNumber, ok = lookupString(structure.SkuNumber, raw.Strings); !ok {
		return nil, BadFormatError
	}
	if data.Family, ok = lookupString(structure.Family, raw.Strings); !ok {
		return nil, BadFormatError
	}

	return &data, nil
}
