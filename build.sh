#!/bin/bash

go build x-crack.go
go env -w CGO_ENABLED=0
go env -w GOOS=linux
go env -w GOARCH=amd64
go build -o build/x-crack_linux_amd64

go env -w CGO_ENABLED=0
go env -w GOOS=windows
go env -w GOARCH=amd64
go build -o build/x-crack_windows_amd64.exe

go env -w CGO_ENABLED=0
go env -w GOOS=darwin
go env -w GOARCH=arm64
go build -o build/xcrack_mac_arm
