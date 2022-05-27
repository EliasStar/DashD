#!/bin/sh

rm -rf build/
mkdir build/

go mod download
go build -o=build/DashD