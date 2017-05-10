#!/bin/bash
set -x
export PATH=$PATH:/usr/local/go/bin/
export GOPATH=`pwd`
export GOARCH=arm
export GOOS=linux

go build ./src/sar && adb push ./sar /sdcard/SAR/ && adb shell su -c /sdcard/SAR/sar
