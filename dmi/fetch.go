/*
 * Author: fasion
 * Created time: 2019-05-22 08:40:43
 * Last Modified by: fasion
 * Last Modified time: 2019-08-28 16:09:52
 */

package dmi

import (
	"github.com/digitalocean/go-smbios/smbios"
)

func FetchRawStructures() ([]*smbios.Structure, error) {
	rc, _, err := smbios.Stream()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	decoder := smbios.NewDecoder(rc)
	return decoder.Decode()
}

func FetchStructures() ([]*Structure, error) {
	raws, err := FetchRawStructures()
	if err != nil {
		return nil, err
	}

	return ParseRawStructures(raws)
}
