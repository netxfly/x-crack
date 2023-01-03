
# 『安全开发教程』年轻人的第一款弱口令扫描器(x-crack)

## 概述

![白帽子安全开发实战](https://github.com/netxfly/sec-dev-in-action-src)第2章扫描器中的一个示例程序。


## Supported Protocol 
* ftp
* mysql
* postgres
* redis
* smb
* snmp
* ssh
* pop3
* imap

考虑到跨平台的兼容性，linux需要开启cgo，暂不支持
* rdp

## Roadmap
1. 扩充支持的协议，计划增加： 
- rdp（需要依赖，暂时没找到好的替代办法）
- ldap（暂不开发）
- memcache（暂不开发）


2. 修改cli传参可直接输入用户名和密码，为类似hydra的形式

```
./x-crack scan -i ssh://127.0.0.1:22 -u root -P ./pass.dic
```




## ref
```
 go tool dist list
```


## Update
* 修改并发数，改为100
* 修改psql的名称，改为postgres
* 增加不同的协议



## BUILD
```
# static options
--extldflags -static
```

```
# linux
$env:CGO_ENABLED="0"
$env:GOOS="linux"
$env:GOARCH="amd64"

go build -ldflags "-w -s " -o build/xcrack_linux_amd64

# windows
$env:CGO_ENABLED="0"
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -ldflags "-w -s --extldflags  -static" -o build/xcrack_windows_amd64.exe
```