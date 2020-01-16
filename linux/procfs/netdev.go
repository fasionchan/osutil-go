/*
 * Author: fasion
 * Created time: 2019-06-18 16:14:10
 * Last Modified by: fasion
 * Last Modified time: 2019-12-20 10:19:18
 */

package procfs

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/fasionchan/libgo/parse"
	"github.com/fasionchan/libgo/unit"
	"github.com/fasionchan/osutil-go/sysinfo"
)

const (
	PROC_NETDEV_PATH = "/proc/net/dev"
)

type CommonCounter struct {
	Bytes       unit.Bytes
	Packets     uint64
	Errors      uint64
	Drops       uint64
	Compresseds uint64
}

type ReceiveCounter struct {
	Bytes         unit.Bytes
	Packets       uint64
	Errors        uint64
	Drops         uint64
	FIFOErrors    uint64
	FramingErrors uint64
	Compresseds   uint64
	Multicasts    uint64
}

func (self ReceiveCounter) Sub(other ReceiveCounter) ReceiveCounter {
	return ReceiveCounter{
		Bytes:         self.Bytes - other.Bytes,
		Packets:       self.Packets - other.Packets,
		Errors:        self.Errors - other.Errors,
		Drops:         self.Drops - other.Drops,
		FIFOErrors:    self.FIFOErrors - other.FIFOErrors,
		FramingErrors: self.FramingErrors - other.FramingErrors,
		Compresseds:   self.Compresseds - other.Compresseds,
		Multicasts:    self.Multicasts - other.Multicasts,
	}
}

type TransmitCounter struct {
	Bytes         unit.Bytes
	Packets       uint64
	Errors        uint64
	Drops         uint64
	FIFOErrors    uint64
	Collisions    uint64
	CarrierLosses uint64
	Compresseds   uint64
}

func (self TransmitCounter) Sub(other TransmitCounter) TransmitCounter {
	return TransmitCounter{
		Bytes:         self.Bytes - other.Bytes,
		Packets:       self.Packets - other.Packets,
		Errors:        self.Errors - other.Errors,
		Drops:         self.Drops - other.Drops,
		FIFOErrors:    self.FIFOErrors - other.FIFOErrors,
		Collisions:    self.Collisions - other.Collisions,
		CarrierLosses: self.CarrierLosses - other.CarrierLosses,
		Compresseds:   self.Compresseds - other.Compresseds,
	}
}

type NetDevCounter struct {
	Name     string
	Receive  ReceiveCounter
	Transmit TransmitCounter
}

func (self NetDevCounter) IsZero() bool {
	return self.Receive.Bytes == 0 &&
		self.Transmit.Bytes == 0
}

type NetDevCounters struct {
	FetchTime time.Time
	Records   []NetDevCounter
}

func (self *NetDevCounters) FilterMisc() {
	records := make([]NetDevCounter, 0)
	for _, record := range self.Records {
		if record.IsZero() {
			continue
		}

		if virtual, err := sysinfo.IsVirtualNic(record.Name); err != nil || virtual {
			continue
		}

		records = append(records, record)
	}

	self.Records = records
}

func FetchNetDevCounters() (*NetDevCounters, error) {
	if false {
		fmt.Println()
	}

	fetchTime := time.Now()
	content, err := ioutil.ReadFile(PROC_NETDEV_PATH)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")

	devs := make([]NetDevCounter, 0, len(lines)-3)

	for _, line := range lines[2 : len(lines)-1] {
		parts := strings.SplitN(line, ":", 2)
		name := strings.TrimSpace(parts[0])

		fields := strings.Fields(strings.TrimSpace(parts[1]))
		values, err := parse.ParseUint64Slice(fields, nil)
		if err != nil {
			return nil, err
		}

		devs = append(devs, NetDevCounter{
			Name: name,
			Receive: ReceiveCounter{
				Bytes:         unit.Bytes(values[0]),
				Packets:       values[1],
				Errors:        values[2],
				Drops:         values[3],
				FIFOErrors:    values[4],
				FramingErrors: values[5],
				Compresseds:   values[6],
				Multicasts:    values[7],
			},
			Transmit: TransmitCounter{
				Bytes:         unit.Bytes(values[8]),
				Packets:       values[9],
				Errors:        values[10],
				Drops:         values[11],
				FIFOErrors:    values[12],
				Collisions:    values[13],
				CarrierLosses: values[14],
				Compresseds:   values[15],
			},
		})
	}

	return &NetDevCounters{
		FetchTime: fetchTime,
		Records:   devs,
	}, nil
}
