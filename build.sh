#!/bin/bash

go build x-crack.go
mv x-crack x-crack_darwin_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build x-crack.go
mv x-crack x-crack_linux_amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build x-crack.go
mv x-crack.exe x-crack_windows_amd64.exe
go build x-crack.go
