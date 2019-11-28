#!/bin/sh

#For developer use only
#go get -u github.com/jteeuwen/go-bindata/...
rm resource.go -r
go-bindata -o=resource.go public views