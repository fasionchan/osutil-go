#!/bin/bash
# FileName:   godefs.sh
# Author:     Fasion Chan
# @contact:   fasionchan@gmail.com
# @version:   $Id$
#
# Description:
#
# Changelog:
#
#

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

    generate "types_netlink"
    generate "types_sock_diag"
)
