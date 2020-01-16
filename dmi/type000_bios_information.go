/*
 * Author: fasion
 * Created time: 2019-09-02 18:48:44
 * Last Modified by: fasion
 * Last Modified time: 2019-09-10 19:24:39
 */

package dmi

import (
	"encoding/binary"
	"bytes"
	"fmt"

	"github.com/digitalocean/go-smbios/smbios"
	"github.com/fasionchan/libgo/arch"
)

var _ = fmt.Println

const (
	TypeBiosInformation = 0
)

type BiosInformationStructure struct {
	Vendor uint8
	Version uint8
	_ uint16
	ReleaseDate uint8
	RomSize uint8
	_ uint64
	_ uint16
	ReleaseMajor uint8
	ReleaseMinor uint8
	FirewareReleaseMajor uint8
	FirewareReleaseMinor uint8
	_ uint16
}

type BiosInformation struct {
	Vendor string
	Version string
	ReleaseDate string
	Release string
	FirmwareRelease string
}

func ToBiosInformation(raw *smbios.Structure) (*BiosInformation, error) {
	sl := len(raw.Formatted) + 4
	if sl < 0x12 {
		return nil, BadFormatError
	}

	var structure BiosInformationStructure
	reader := bytes.NewReader(padZeros(raw.Formatted, 0x1a - 4))
	if err := binary.Read(reader, arch.NativeEndian, &structure); err != nil {
		return nil, err
	}

	data := BiosInformation{

	}

	var ok bool
	if data.Vendor, ok = lookupString(structure.Vendor, raw.Strings); !ok {
		return nil, BadFormatError
	}
	if data.Version, ok = lookupString(structure.Version, raw.Strings); !ok {
		return nil, BadFormatError
	}
	if data.ReleaseDate, ok = lookupString(structure.ReleaseDate, raw.Strings); !ok {
		return nil, BadFormatError
	}

	if sl >= 0x16 {
		data.Release = fmt.Sprintf("%d.%d", structure.ReleaseMajor, structure.ReleaseMinor)
	}

	if sl >= 0x18 {
		data.FirmwareRelease = fmt.Sprintf("%d.%d", structure.FirewareReleaseMajor, structure.FirewareReleaseMinor)
	}


	return &data, nil
}
