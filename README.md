
# 『安全开发教程』年轻人的第一款弱口令扫描器(x-crack)

## 概述

![白帽子安全开发实战](https://github.com/netxfly/sec-dev-in-action-src)第2章扫描器中的一个示例程序。


## Supported Protocol 
* ftp
* mysql, (TODO)mysql db_name -> mysql
* postgres
* redis
* smb
* snmp
* ssh
* pop3
* imap
* 

## Roadmap
1. 扩充支持的协议，计划增加： 
- rdp（需要依赖，暂时没找到好的替代办法）
- x ldap（暂不开发）
- memcache（暂不开发）


2. 修改cli传参可直接输入用户名和密码，为类似hydra的形式

```
./x-crack scan -i ssh://127.0.0.1:22 -u root -p 123456
```




## ref
```
 go tool dist list
```
