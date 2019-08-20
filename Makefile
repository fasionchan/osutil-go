# Author: fasion
# Created time: 2019-08-06 16:34:07
# Last Modified by: fasion
# Last Modified time: 2019-08-13 18:04:10

build-ss-linux:
	GOOS=linux CGO_ENABLED=0 go build github.com/fasionchan/osutil-go/cmd/ss

build-nics-windows:
	GOOS=windows GOARCH=386 CGO_ENABLED=0 go build github.com/fasionchan/osutil-go/cmd/nics
