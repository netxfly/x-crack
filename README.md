
# 『安全开发教程』年轻人的第一款弱口令扫描器(x-crack)

## 概述

![白帽子安全开发实战](https://github.com/netxfly/sec-dev-in-action-src)第2章扫描器中的一个示例程序。


## BUILD

```
# linux
$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"

go build -ldflags "-w -s" -o xcrack_linux_amd64/xcrack

# windows
$env:CGO_ENABLED="0"
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -ldflags "-w -s" -o xcrack_windows_amd64.exe
```