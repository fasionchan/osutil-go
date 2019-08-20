#!/bin/sh

# Author: fasion
# Created time: 2019-06-28 14:11:44
# Last Modified by: fasion
# Last Modified time: 2019-08-07 08:50:36

generate() {
	name="$1"
	constraint="// +build $arch"

	(
		echo "$constraint"
		echo

		go tool cgo -godefs "${name}.go"
	) > "z${name}_${arch}.go"
}

(
	export arch="$(go version | cut -d ' ' -f 4 | cut -d / -f 2)"

	cd `dirname "$0"`

	generate "types_common"
	generate "types_elfauxv"
	generate "types_net"
)
