/*
 * Author: fasion
 * Created time: 2019-12-20 10:18:55
 * Last Modified by: fasion
 * Last Modified time: 2020-03-27 17:34:23
 */

package procfs

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"
)

const (
	PROC_NETTCP_PATH  = "/proc/net/tcp"
	PROC_NETTCP6_PATH = "/proc/net/tcp6"
)

var BadNetTcpFormatError = fmt.Errorf("bad /proc/net/tcp formart")

type NetTcpFields []string

func reverseByteSlice(bytes []byte) []byte {
	n := len(bytes)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}

func parseAddressPair(text string) (net.IP, int, error) {
	parts := strings.Split(strings.TrimSpace(text), ":")
	if len(parts) != 2 {
		return nil, 0, BadNetTcpFormatError
	}

	decoded, err := hex.DecodeString(parts[0])
	if err != nil {
		return nil, 0, BadNetTcpFormatError
	}
	// reverse for IPv4
	if len(decoded) == 4 {
		decoded = []byte{decoded[3], decoded[2], decoded[1], decoded[0]}
	}
	// reverse for IPv6
	if len(decoded) == 16 {
		if bytes.HasPrefix(decoded, []byte{0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 0, 0}) {
			decoded[12], decoded[15] = decoded[15], decoded[12]
			decoded[13], decoded[14] = decoded[14], decoded[13]
		}
	}

	port, err := strconv.ParseInt(parts[1], 16, 32)
	if err != nil {
		return nil, 0, BadNetTcpFormatError
	}

	return net.IP(decoded), int(port), nil
}

func (fields NetTcpFields) Local() (net.IP, int, error) {
	return parseAddressPair(fields[0])
}

func (fields NetTcpFields) Remote() (net.IP, int, error) {
	return parseAddressPair(fields[1])
}

func (fields NetTcpFields) State() (uint8, error) {
	state, err := strconv.ParseUint(strings.TrimSpace(fields[2]), 16, 8)
	if err != nil {
		return 0, err
	}

	return uint8(state), nil
}

type NetTcpScanner struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewNetTcpScanner(path string, bufferSize int) (*NetTcpScanner, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(file)

	buffer := make([]byte, 0, bufferSize)
	scanner.Buffer(buffer, 0)

	// skip the first line
	if scanner.Scan() {
		scanner.Text()
	}

	return &NetTcpScanner{
		file:    file,
		scanner: scanner,
	}, nil
}

func NewNetTcp4Scanner(bufferSize int) (*NetTcpScanner, error) {
	return NewNetTcpScanner(PROC_NETTCP_PATH, bufferSize)
}

func NewNetTcp6Scanner(bufferSize int) (*NetTcpScanner, error) {
	return NewNetTcpScanner(PROC_NETTCP6_PATH, bufferSize)
}

func NewNetTcpxScanner(family uint8, bufferSize int) (*NetTcpScanner, error) {
	switch family {
	case syscall.AF_INET6:
		return NewNetTcp6Scanner(bufferSize)
	default:
		return NewNetTcp4Scanner(bufferSize)
	}
}

func (scanner *NetTcpScanner) Close() error {
	return scanner.file.Close()
}

func (scanner *NetTcpScanner) Scan() bool {
	return scanner.scanner.Scan()
}

func (scanner *NetTcpScanner) Record() (string, NetTcpFields, error) {
	line := strings.TrimSpace(scanner.scanner.Text())
	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return "", nil, BadNetTcpFormatError
	}

	sl := strings.TrimSpace(parts[0])
	fields := strings.Fields(strings.TrimSpace(parts[1]))
	if len(fields) < 11 {
		return "", nil, BadNetTcpFormatError
	}

	return sl, NetTcpFields(fields), nil
}
