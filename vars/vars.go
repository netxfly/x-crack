/*

Copyright (c) 2017 xsec.io

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package vars

import (
	"github.com/patrickmn/go-cache"

	"gopkg.in/cheggaaa/pb.v2"

	"strings"
	"sync"
	"time"
)

var (
	IpList     = "../iplist.txt"
	ResultFile = "../results.txt"

	UserDictFile = "../user.dic"
	PassDictFile = "../pass.dic"

	USER = "root"
	PASS = "123456"

	USERNAME = ""
	PASSWORD = ""
	IP       = "" // ssh://127.0.0.1

	TimeOut = 3 * time.Second
	ScanNum = 100

	DebugMode bool

	StartTime time.Time

	ProgressBar      *pb.ProgressBar
	ProcessBarActive *pb.ProgressBar
)

var (
	CacheService *cache.Cache
	Mutex        sync.Mutex

	PortNames = map[int]string{
		21:    "FTP",
		22:    "SSH",
		161:   "SNMP",
		445:   "SMB",
		1433:  "MSSQL",
		3306:  "MYSQL",
		3389:  "RDP",
		5432:  "POSTGRES",
		6379:  "REDIS",
		9200:  "ELASTICSEARCH",
		27017: "MONGODB",
		110:   "POP3",
		143:   "IMAP",
		389:   "LDAP",
	}

	UdpProtocols = map[string]bool{
		"SNMP": true,
	}

	// 标记特定服务的特定用户是否破解成功，成功的话不再尝试破解该用户
	SuccessHash map[string]bool

	SupportProtocols map[string]bool
)

func init() {
	SuccessHash = make(map[string]bool)
	CacheService = cache.New(cache.NoExpiration, cache.DefaultExpiration)

	SupportProtocols = make(map[string]bool)
	for _, proto := range PortNames {
		SupportProtocols[strings.ToUpper(proto)] = true
	}

}
