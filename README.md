
# 『安全开发教程』年轻人的第一款弱口令扫描器(x-crack)

## 概述

我们在做企业安全时，弱口令检测是系统/网络安全的最基础的部分之一，根据经验，经常会出现弱口令的服务如下：

- FTP
- SSH
- SMB
- MYSQL
- MSSQL
- POSTGRESQL
- REDIS
- MONGODB
- ELASTICSEARCH

那咱们就一起用GO来写一款常见服务的弱口令扫描器，且支持以插件的形式增加新的服务扫描模块。我们的教程暂定为只扫以上服务。

给扫描器启一个屌炸天的名字`x-crack`，在$GOPATH/src/中建立一个x-crack项目后开始撸码，不要给我说什么底层原理、框架内核，老夫敲代码就是一把梭。

开发完毕的项目地址为：[https://github.com/netxfly/x-crack](https://github.com/netxfly/x-crack)

[![asciicast](https://asciinema.org/a/3lKBxOMvevUZPR6DCSwsH4E3S.png)](https://asciinema.org/a/3lKBxOMvevUZPR6DCSwsH4E3S)

## 开工
### 数据结构定义

- 扫描模块的输入内容为为IP、端口及协议的列表，我们需要定义一个IpAddr的数据结构； 
- 每个服务的每次扫描需要传入的参数为IP、端口、协议、用户名和密码，需要定义一个Service结构来包括这些内容；
- 每条Service的记录在扫描模块进行尝试后，会得出扫描结果成功与否，我们再定义一个ScanResult数据结构。

按照开发规范，数据结构的定义统一放到models目录中，全部的数据结构定义如下：

```go

package models

type Service struct {
	Ip       string
	Port     int
	Protocol string
	Username string
	Password string
}

type ScanResult struct {
	Service Service
	Result  bool
}

type IpAddr struct {
	Ip       string
	Port     int
	Protocol string
}
```

### FTP扫描模块

go语言有现成的FTP模块，我们找一个star数最多的直接`go get`安装一下即可使用了：
```bash
go get -u github.com/jlaffaye/ftp
``` 
我们把所有的扫描模块放到`plugins`目录中，FTP协议的扫描插件如下所示：
```go

package plugins

import (
	"github.com/jlaffaye/ftp"

	"x-crack/models"
	"x-crack/vars"


	"fmt"
)

func ScanFtp(s models.Service) (err error, result models.ScanResult) {
	result.Service = s
	conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", s.Ip, s.Port), vars.TimeOut)
	if err == nil {
		err = conn.Login(s.Username, s.Password)
		if err == nil {
			defer conn.Logout()
			result.Result = true
		}
	}
	return err, result
}
```

每个连接需要设置超时时间，防止因网络问题导致的阻塞，我们打算通过程序的命令行来控制超时时间，所以定义了一个全局变量TimeOut。
放在vars模块中的原因是防止放在这个模块中后会和其他模块互相调用导致的循环import

写代码虽然可以一把梭，但是不能等着洋洋洒洒地把几万行都写完再运行，比如我们的目标是造一辆豪车，不能等着所有零件设计好，都装上去再发动车测试，正确的开发流程是把写边测，不要等轮子造出来，而是在螺丝、齿轮阶段就测试。

以下为FTP扫描插件这个齿轮的测试代码及结果。

```go
package plugins_test

import (
	"x-crack/models"
	"x-crack/plugins"

	"testing"
)

func TestScanFtp(t *testing.T) {
	s := models.Service{Ip: "127.0.0.1", Port: 21, Protocol: "ftp", Username: "ftp", Password: "ftp"}
	t.Log(plugins.ScanFtp(s))
}
```

测试结果满足预期，说明我们这个零件不是次品，可以继续再造其他零件了。

```bash
$ go test -v plugins/ftp_test.go
=== RUN   TestScanFtp
--- PASS: TestScanFtp (0.00s)
	ftp_test.go:36: dial tcp 127.0.0.1:21: getsockopt: connection refused {{127.0.0.1 21 ftp ftp ftp} false}
PASS
ok  	command-line-arguments	0.025s
```

### SSH扫描模块

go的标准库中自带了ssh包，直接调用即可，完整代码如下：

```go

package plugins

import (
	"golang.org/x/crypto/ssh"

	"x-crack/models"
	"x-crack/vars"

	"fmt"
	"net"
)

func ScanSsh(s models.Service) (err error, result models.ScanResult) {
	result.Service = s
	config := &ssh.ClientConfig{
		User: s.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		Timeout: vars.TimeOut,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", s.Ip, s.Port), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		errRet := session.Run("echo xsec")
		if err == nil && errRet == nil {
			defer session.Close()
			result.Result = true
		}
	}
	return err, result
}
 
```
同样，每个子模块写好后都需要先用go test跑一下看是否满足预期，测试代码如下：
```go

package plugins_test

import (
	"x-crack/models"
	"x-crack/plugins"

	"testing"
)

func TestScanSsh(t *testing.T) {
	s := models.Service{Ip: "127.0.0.1", Port: 22, Username: "root", Password: "123456", Protocol: "ssh"}
	t.Log(plugins.ScanSsh(s))
} 
```

测试结果如下：

```go
$ go test -v plugins/ssh_test.go
=== RUN   TestScanSsh
--- PASS: TestScanSsh (0.00s)
	ssh_test.go:36: dial tcp 127.0.0.1:22: getsockopt: connection refused {{127.0.0.1 22 ssh root 123456} false}
PASS
ok  	command-line-arguments	0.026s
```

### SMB扫描模块
SMB弱口令的扫描插件，我们使用了`github.com/stacktitan/smb/smb`包，同样直接`go get`安装一下即可拿来使用。
代码如下：

```go

package plugins

import (
	"github.com/stacktitan/smb/smb"

	"x-crack/models"
)

func ScanSmb(s models.Service) (err error, result models.ScanResult) {
	result.Service = s
	options := smb.Options{
		Host:        s.Ip,
		Port:        s.Port,
		User:        s.Username,
		Password:    s.Password,
		Domain:      "",
		Workstation: "",
	}

	session, err := smb.NewSession(options, false)
	if err == nil {
		session.Close()
		if session.IsAuthenticated {
			result.Result = true
		}
	}
	return err, result
}

``` 

同样也先写测试用例来测试一下，测试代码如下：
```go
package plugins_test

import (
	"x-crack/models"
	"x-crack/plugins"

	"testing"
)

func TestScanSmb(t *testing.T) {
	s := models.Service{Ip: "share.xsec.io", Port: 445, Protocol: "smb", Username: "xsec", Password: "fsafffdsfdsa"}
	t.Log(plugins.ScanSmb(s))
}
```

测试结果：

```bash
hartnett at hartnettdeMacBook-Pro in /data/code/golang/src/x-crack (master)
$ go test -v plugins/smb_test.go
=== RUN   TestScanSmb
--- PASS: TestScanSmb (0.04s)
	smb_test.go:36: NT Status Error: Logon failed
		 {{share.xsec.io 445 smb xsec fsafffdsfdsa} false}
PASS
ok  	command-line-arguments	0.069s
```

### MYSQL、MSSQL和POSTGRESQL扫描模块

MYSQL、MSSQL和POSTGRESQL的扫描模块，我使用了第三方的ORM `xorm`，当然也可以直接使用原生的sql driver来实现，我们这里图方便用`xorm`一把梭了。
对于`xorm`来说，这3个扫描插件的实现方法大同小异，为了节约篇幅，咱们只看mysql扫描插件的实现，其他2个插件可以参考github中的完整源码。
首先还是先`go get`要用到的包:

```bash
go get github.com/netxfly/mysql
go get github.com/go-xorm/xorm
github.com/go-xorm/core
```

接下来我们把需要验证的IP、port、username、password组成datasource传递给xorm，完整代码如下：

```go
package plugins

import (
	_ "github.com/netxfly/mysql"
	"github.com/go-xorm/xorm"
	"github.com/go-xorm/core"

	"x-crack/models"

	"fmt"
)

func ScanMysql(service models.Service) (err error, result models.ScanResult) {
	result.Service = service

	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", service.Username,
		service.Password, service.Ip, service.Port, "mysql")
	Engine, err := xorm.NewEngine("mysql", dataSourceName)

	if err == nil {
		Engine.SetLogLevel(core.LOG_OFF)
		// fix "[mysql] packets.go:33: unexpected EOF" error
		Engine.SetMaxIdleConns(0)
		// Engine.SetConnMaxLifetime(time.Second * 30)
		defer Engine.Close()
		err = Engine.Ping()
		if err == nil {
			result.Result = true
		}
	}
	return err, result
}
```
眼尖的同学也许发现了，上面 `github.com/netxfly/mysql` 这个mysql包是放在笔者的github下的，这是为什么呢？

因为直接用mysql这个包的话，在扫描的过程中会遇到`[mysql] packets.go:33: unexpected EOF" error`的异常输出，影响了我们程序在扫描过程中输出UI的美观性，这对于帅气的我是无法接受的，通过设置参数的方法无法解决，最后只好直接fork了一份mysql的包，把打印这个异常的语句注释掉再提交上去直接使用了。

测试代码：

```go

package plugins_test

import (
	"testing"

	"x-crack/plugins"
	"x-crack/models"
)

func TestScanMysql(t *testing.T) {
	service := models.Service{Ip: "10.10.10.10", Port: 3306, Protocol: "mysql", Username: "root", Password: "123456"}
	t.Log(plugins.ScanMysql(service))
}
```

测试结果：

```bash
go test -v plugins/mysql_test.go
=== RUN   TestScanMysql
--- PASS: TestScanMysql (0.02s)
	mysql_test.go:36: Error 1045: Access denied for user 'root'@'10.10.10.100' (using password: YES) {{10.10.10.10 3306 mysql root 123456} false}
PASS
ok  	command-line-arguments	0.041s
```

### Redis扫描模块

`go get`安装第三方包`github.com/go-redis/redis`，完整代码如下：

```go

package plugins

import (
	"github.com/go-redis/redis"

	"x-crack/models"
	"x-crack/vars"

	"fmt"
)

func ScanRedis(s models.Service) (err error, result models.ScanResult) {
	result.Service = s
	opt := redis.Options{Addr: fmt.Sprintf("%v:%v", s.Ip, s.Port),
		Password: s.Password, DB: 0, DialTimeout: vars.TimeOut}
	client := redis.NewClient(&opt)
	defer client.Close()
	_, err = client.Ping().Result()
	if err == nil {
		result.Result = true
	}
	return err, result
}

```
测试代码：
```go

package plugins_test

import (
	"x-crack/models"
	"x-crack/plugins"

	"testing"
)

func TestScanRedis(t *testing.T) {
	s := models.Service{Ip: "127.0.0.1", Port: 6379, Password: "test"}
	t.Log(plugins.ScanRedis(s))
}
```
测试结果：

```bash
go test -v plugins/redis_test.go
=== RUN   TestScanRedis
--- PASS: TestScanRedis (0.00s)
	redis_test.go:36: dial tcp 127.0.0.1:6379: getsockopt: connection refused {{127.0.0.1 6379   test} false}
PASS
ok  	command-line-arguments	0.025s
```

### MONGODB扫描模块

mongodb扫描模块依赖mgo包，可用`go get`合令直接安装。

```bash
go get gopkg.in/mgo.v2
```

完整代码：
```go

package plugins

import (
	"gopkg.in/mgo.v2"

	"x-crack/models"
	"x-crack/vars"

	"fmt"
)

func ScanMongodb(s models.Service) (err error, result models.ScanResult) {
	result.Service = s
	url := fmt.Sprintf("mongodb://%v:%v@%v:%v/%v", s.Username, s.Password, s.Ip, s.Port, "test")
	session, err := mgo.DialWithTimeout(url, vars.TimeOut)

	if err == nil {
		defer session.Close()
		err = session.Ping()
		if err == nil {
			result.Result = true
		}
	}

	return err, result
}
```

测试结果：

```bash
go test -v plugins/mongodb_test.go
=== RUN   TestScanMongodb
--- PASS: TestScanMongodb (3.53s)
	mongodb_test.go:36: no reachable servers {{127.0.0.1 27017 mongodb test test} false}
PASS
ok  	command-line-arguments	3.558s
```

### ELASTICSEARCH扫描模块
ELASTICSEARCH扫描插件依赖第三方包`gopkg.in/olivere/elastic.v3`，同样也是直接`go get`安装。
完整代码如下：

```go
package plugins

import (
	"gopkg.in/olivere/elastic.v3"

	"x-crack/models"
	
	"fmt"
)

func ScanElastic(s models.Service) (err error, result models.ScanResult) {
	result.Service = s
	client, err := elastic.NewClient(elastic.SetURL(fmt.Sprintf("http://%v:%v", s.Ip, s.Port)),
		elastic.SetMaxRetries(3),
		elastic.SetBasicAuth(s.Username, s.Password),
	)
	if err == nil {
		_, _, err = client.Ping(fmt.Sprintf("http://%v:%v", s.Ip, s.Port)).Do()
		if err == nil {
			result.Result = true
		}
	}
	return err, result
}

```
测试代码：

```go
package plugins_test

import (
	"x-crack/models"
	"x-crack/plugins"

	"testing"
)

func TestScanElastic(t *testing.T) {
	s := models.Service{Ip: "127.0.0.1", Port: 9200, Protocol: "elastic", Username: "root", Password: "123456"}
	t.Log(plugins.ScanElastic(s))
}

```
测试结果如下：

```bash
go test -v plugins/elastic_test.go
=== RUN   TestScanElastic
--- PASS: TestScanElastic (5.02s)
	elastic_test.go:36: no Elasticsearch node available {{127.0.0.1 9200 elastic root 123456} false}
PASS
ok  	command-line-arguments	5.061s
```

### 扫描模块插件化

前面我们写好的扫描插件的函数原始是一致，我们可以将这组函数放到一个map中，在扫描的过程中自动化根据不同的协议调用不同的扫描插件。

以后新加的扫描插件，可以按这种方法直接注册。

```go

package plugins

import (
	"x-crack/models"
)

type ScanFunc func(service models.Service) (err error, result models.ScanResult)

var (
	ScanFuncMap map[string]ScanFunc
)

func init() {
	ScanFuncMap = make(map[string]ScanFunc)
	ScanFuncMap["FTP"] = ScanFtp
	ScanFuncMap["SSH"] = ScanSsh
	ScanFuncMap["SMB"] = ScanSmb
	ScanFuncMap["MSSQL"] = ScanMssql
	ScanFuncMap["MYSQL"] = ScanMysql
	ScanFuncMap["POSTGRESQL"] = ScanPostgres
	ScanFuncMap["REDIS"] = ScanRedis
	ScanFuncMap["ELASTICSEARCH"] = ScanElastic
	ScanFuncMap["MONGODB"] = ScanMongodb
}

```

## 扫描任务调度

前面我们写好了一些常见服务的弱口令扫描插件，也测试通过了。
接下来我们需要实现从命令行参数传递iplist、用户名字典和密码字典进去，并读取相应的信息进行扫描调度的功能，细分一下，需要做以下几件事：

- 读取iplist列表
- 读取用户名字典
- 读取密码字典
- 生成扫描任务
- 扫描任务调度
- 扫描任务执行
- 扫描结果保存
- 命令行调用外壳

### 读取ip\用户名和密码字典

该模块主要用了标准库中的`bufio`包，逐行读取文件，进行过滤后直接生成相应的slice。其中iplist支持以下格式：
```bash
127.0.0.1:3306|mysql
8.8.8.8:22
9.9.9.9:6379
108.61.223.105:2222|ssh
```

对于标准的端口，程序可以自动判断其协议，对于非标准端口的协议，需要在后面加一个字段标注一下协议。

为了防止咱们的程序被脚本小子们滥用，老夫就不提供端口扫描、协议识别等功能了，安全工程师们可以把自己公司的端口扫描器产出的结果丢到这个里面来扫。

```go

package util

import (
	"x-crack/models"
	"x-crack/logger"
	"x-crack/vars"

	"os"
	"bufio"
	"strings"
	"strconv"
)

func ReadIpList(fileName string) (ipList []models.IpAddr) {
	ipListFile, err := os.Open(fileName)
	if err != nil {
		logger.Log.Fatalf("Open ip List file err, %v", err)
	}

	defer ipListFile.Close()

	scanner := bufio.NewScanner(ipListFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		ipPort := strings.TrimSpace(scanner.Text())
		t := strings.Split(ipPort, ":")
		ip := t[0]
		portProtocol := t[1]
		tmpPort := strings.Split(portProtocol, "|")
		// ip列表中指定了端口对应的服务
		if len(tmpPort) == 2 {
			port, _ := strconv.Atoi(tmpPort[0])
			protocol := strings.ToUpper(tmpPort[1])
			if vars.SupportProtocols[protocol] {
				addr := models.IpAddr{Ip: ip, Port: port, Protocol: protocol}
				ipList = append(ipList, addr)
			} else {
				logger.Log.Infof("Not support %v, ignore: %v:%v", protocol, ip, port)
			}
		} else {
			// 通过端口查服务
			port, err := strconv.Atoi(tmpPort[0])
			if err == nil {
				protocol, ok := vars.PortNames[port]
				if ok && vars.SupportProtocols[protocol] {
					addr := models.IpAddr{Ip: ip, Port: port, Protocol: protocol}
					ipList = append(ipList, addr)
				}
			}
		}

	}

	return ipList
}

func ReadUserDict(userDict string) (users []string, err error) {
	file, err := os.Open(userDict)
	if err != nil {
		logger.Log.Fatalf("Open user dict file err, %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			users = append(users, user)
		}
	}
	return users, err
}

func ReadPasswordDict(passDict string) (password []string, err error) {
	file, err := os.Open(passDict)
	if err != nil {
		logger.Log.Fatalf("Open password dict file err, %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		passwd := strings.TrimSpace(scanner.Text())
		if passwd != "" {
			password = append(password, passwd)
		}
	}
	password = append(password, "")
	return password, err
}

```
IP列表、用户名字典与密码字典读取的测试代码：

```go

package util_test

import (
	"x-crack/util"

	"testing"
)

func TestReadIpList(t *testing.T) {
	ipList := "/tmp/iplist.txt"
	t.Log(util.ReadIpList(ipList))
}

func TestReadUserDict(t *testing.T) {
	userDict := "/tmp/user.dic"
	t.Log(util.ReadUserDict(userDict))
}

func TestReadPasswordDict(t *testing.T) {
	passDict := "/tmp/pass.dic"
	t.Log(util.ReadPasswordDict(passDict))
}

```
这个模块的测试结果如下：

```bash
go test -v util/file_test.go
=== RUN   TestReadIpList
--- PASS: TestReadIpList (0.00s)
	file_test.go:35: [{127.0.0.1 3306 MYSQL} {8.8.8.8 22 SSH} {9.9.9.9 6379 REDIS} {108.61.223.105 2222 SSH}]
=== RUN   TestReadUserDict
--- PASS: TestReadUserDict (0.00s)
	file_test.go:40: [root admin test guest info adm mysql user administrator ftp sa] <nil>
=== RUN   TestReadPasswordDict
--- PASS: TestReadPasswordDict (0.00s)
	file_test.go:45: [1314520520 135246 135246789 135792468 1357924680 147258369 1472583690 1qaz2wsx 5201314 54321 55555 654321 789456123 88888 888888 88888888 987654321 9876543210 ^%$#@~! a123123 a123456 a12345678 a123456789 aa123456 aa123456789 aaa123456 aaaaa aaaaaa aaaaaaaa abc123 abc123456 abc123456789 abcd123 abcd1234 abcd123456 admin admin888 ] <nil>
PASS
ok  	command-line-arguments	0.022s
```

其中iplist在加载的过程中不是无脑全部读进去的，在正式扫描前会先过滤一次，把不通的ip和端口对剔除掉，以免影响扫描效率，代码如下：
```go

package util

import (
	"gopkg.in/cheggaaa/pb.v2"

	"x-crack/models"
	"x-crack/logger"
	"x-crack/vars"

	"net"
	"sync"
	"fmt"
)

var (
	AliveAddr []models.IpAddr
	mutex     sync.Mutex
)

func init() {
	AliveAddr = make([]models.IpAddr, 0)
}

func CheckAlive(ipList []models.IpAddr) ([]models.IpAddr) {
	logger.Log.Infoln("checking ip active")
	
	var wg sync.WaitGroup
	wg.Add(len(ipList))

	for _, addr := range ipList {
		go func(addr models.IpAddr) {
			defer wg.Done()
			SaveAddr(check(addr))
		}(addr)
	}
	wg.Wait()
	vars.ProcessBarActive.Finish()

	return AliveAddr
}

func check(ipAddr models.IpAddr) (bool, models.IpAddr) {
	alive := false
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ipAddr.Ip, ipAddr.Port), vars.TimeOut)
	if err == nil {
		alive = true
	}
	vars.ProcessBarActive.Increment()
	return alive, ipAddr
}

func SaveAddr(alive bool, ipAddr models.IpAddr) {
	if alive {
		mutex.Lock()
		AliveAddr = append(AliveAddr, ipAddr)
		mutex.Unlock()
	}
}

```

通过标准端口查询对应服务的功能在vars包中定义了，为了避免多个包之间的循环导入，我们把所有的全局变量都集中到了一个独立的vars包中。

`PortNames`map为标准端口对应的服务，在加了新的扫描插件后，也需要更新这个map的内容。

```go

package vars

import (
	"github.com/patrickmn/go-cache"

	"gopkg.in/cheggaaa/pb.v2"

	"sync"
	"time"
	"strings"
)

var (
	IpList     = "iplist.txt"
	ResultFile = "x_crack.txt"

	UserDict = "user.dic"
	PassDict = "pass.dic"

	TimeOut = 3 * time.Second
	ScanNum = 5000

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
		445:   "SMB",
		1433:  "MSSQL",
		3306:  "MYSQL",
		5432:  "POSTGRESQL",
		6379:  "REDIS",
		9200:  "ELASTICSEARCH",
		27017: "MONGODB",
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

```

### 任务调度
任务调度模块包含了生成扫描任务，按指定的协程数分发和执行扫描任务的功能。

```go

package util

import (
	"github.com/sirupsen/logrus"

	"gopkg.in/cheggaaa/pb.v2"

	"x-crack/models"
	"x-crack/logger"
	"x-crack/vars"
	"x-crack/util/hash"
	"x-crack/plugins"

	"sync"
	"strings"
	"fmt"
	"time"
)

func GenerateTask(ipList []models.IpAddr, users []string, passwords []string) (tasks []models.Service, taskNum int) {
	tasks = make([]models.Service, 0)

	for _, user := range users {
		for _, password := range passwords {
			for _, addr := range ipList {
				service := models.Service{Ip: addr.Ip, Port: addr.Port, Protocol: addr.Protocol, Username: user, Password: password}
				tasks = append(tasks, service)
			}
		}
	}

	return tasks, len(tasks)
}

func DistributionTask(tasks []models.Service) () {
	totalTask := len(tasks)
	scanBatch := totalTask / vars.ScanNum
	logger.Log.Infoln("Start to scan")
	
	for i := 0; i < scanBatch; i++ {
		curTasks := tasks[vars.ScanNum*i:vars.ScanNum*(i+1)]
		ExecuteTask(curTasks)
	}

	if totalTask%vars.ScanNum > 0 {
		lastTask := tasks[vars.ScanNum*scanBatch:totalTask]
		ExecuteTask(lastTask)
	}

	models.SavaResultToFile()
	models.ResultTotal()
	models.DumpToFile(vars.ResultFile)
}

func ExecuteTask(tasks []models.Service) () {
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for _, task := range tasks {
		if vars.DebugMode {
			logger.Log.Debugf("checking: Ip: %v, Port: %v, [%v], UserName: %v, Password: %v", task.Ip, task.Port,
				task.Protocol, task.Username, task.Password)
		}

		var k string
		protocol := strings.ToUpper(task.Protocol)

		if protocol == "REDIS" || protocol == "FTP" {
			k = fmt.Sprintf("%v-%v-%v", task.Ip, task.Port, task.Protocol)
		} else {
			k = fmt.Sprintf("%v-%v-%v", task.Ip, task.Port, task.Username)
		}

		h := hash.MakeTaskHash(k)
		if hash.CheckTashHash(h) {
			wg.Done()
			continue
		}

		go func(task models.Service, protocol string) {
			defer wg.Done()
			fn := plugins.ScanFuncMap[protocol]
			models.SaveResult(fn(task))
		}(task, protocol)

		vars.ProgressBar.Increment()
	}
	waitTimeout(&wg, vars.TimeOut)
}

```

个别扫描插件没有指定超时时间的功能，所以我们额外为所有的扫描插件都提供了一个超时函数，防止个别协程被阻塞，影响了扫描器整体的速度。
```go
// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

```

任务调度模块的测试代码如下：

```go

package util_test

import (
	"x-crack/util"

	"testing"
)

func TestGenerateTask(t *testing.T) {
	ipList := "/tmp/iplist.txt"
	userDic := "/tmp/user.dic"
	passDic := "/tmp/pass.dic"

	users, _ := util.ReadUserDict(userDic)
	passwords, _ := util.ReadPasswordDict(passDic)

	t.Log(util.GenerateTask(util.ReadIpList(ipList), users, passwords))
}

func TestDistributionTask(t *testing.T) {
	ipList := "/tmp/iplist.txt"
	userDic := "/tmp/user.dic"
	passDic := "/tmp/pass.dic"

	users, _ := util.ReadUserDict(userDic)
	passwords, _ := util.ReadPasswordDict(passDic)

	tasks, _ := util.GenerateTask(util.ReadIpList(ipList), users, passwords)
	util.DistributionTask(tasks)
}

```
测试结果如下：
```bash
$ go test -v util/task_test.go
=== RUN   TestGenerateTask
--- PASS: TestGenerateTask (0.00s)
	task_test.go:41: [{127.0.0.1 3306 MYSQL root admin} {8.8.8.8 22 SSH root admin} {9.9.9.9 6379 REDIS root admin} {108.61.223.105 2222 SSH root admin} {127.0.0.1 3306 MYSQL root admin888} {8.8.8.8 22 SSH root admin888} {9.9.9.9 6379 REDIS root admin888} {108.61.223.105 2222 SSH root admin888} {127.0.0.1 3306 MYSQL root 123456} {8.8.8.8 22 SSH root 123456} {9.9.9.9 6379 REDIS root 123456} {108.61.223.105 2222 SSH root 123456} {127.0.0.1 3306 MYSQL root } {8.8.8.8 22 SSH root } {9.9.9.9 6379 REDIS root } {108.61.223.105 2222 SSH root } {127.0.0.1 3306 MYSQL admin admin} {8.8.8.8 22 SSH admin admin} {9.9.9.9 6379 REDIS admin admin} {108.61.223.105 2222 SSH admin admin} {127.0.0.1 3306 MYSQL admin admin888} {8.8.8.8 22 SSH admin admin888} {9.9.9.9 6379 REDIS admin admin888} {108.61.223.105 2222 SSH admin admin888} {127.0.0.1 3306 MYSQL admin 123456} {8.8.8.8 22 SSH admin 123456} {9.9.9.9 6379 REDIS admin 123456} {108.61.223.105 2222 SSH admin 123456} {127.0.0.1 3306 MYSQL admin } {8.8.8.8 22 SSH admin } {9.9.9.9 6379 REDIS admin } {108.61.223.105 2222 SSH admin } {127.0.0.1 3306 MYSQL test admin} {8.8.8.8 22 SSH test admin} {9.9.9.9 6379 REDIS test admin} {108.61.223.105 2222 SSH test admin} {127.0.0.1 3306 MYSQL test admin888} {8.8.8.8 22 SSH test admin888} {9.9.9.9 6379 REDIS test admin888} {108.61.223.105 2222 SSH test admin888} {127.0.0.1 3306 MYSQL test 123456} {8.8.8.8 22 SSH test 123456} {9.9.9.9 6379 REDIS test 123456} {108.61.223.105 2222 SSH test 123456} {127.0.0.1 3306 MYSQL test } {8.8.8.8 22 SSH test } {9.9.9.9 6379 REDIS test } {108.61.223.105 2222 SSH test } {127.0.0.1 3306 MYSQL guest admin} {8.8.8.8 22 SSH guest admin} {9.9.9.9 6379 REDIS guest admin} {108.61.223.105 2222 SSH guest admin} {127.0.0.1 3306 MYSQL guest admin888} {8.8.8.8 22 SSH guest admin888} {9.9.9.9 6379 REDIS guest admin888} {108.61.223.105 2222 SSH guest admin888} {127.0.0.1 3306 MYSQL guest 123456} {8.8.8.8 22 SSH guest 123456} {9.9.9.9 6379 REDIS guest 123456} {108.61.223.105 2222 SSH guest 123456} {127.0.0.1 3306 MYSQL guest } {8.8.8.8 22 SSH guest } {9.9.9.9 6379 REDIS guest } {108.61.223.105 2222 SSH guest } {127.0.0.1 3306 MYSQL info admin} {8.8.8.8 22 SSH info admin} {9.9.9.9 6379 REDIS info admin} {108.61.223.105 2222 SSH info admin} {127.0.0.1 3306 MYSQL info admin888} {8.8.8.8 22 SSH info admin888} {9.9.9.9 6379 REDIS info admin888} {108.61.223.105 2222 SSH info admin888} {127.0.0.1 3306 MYSQL info 123456} {8.8.8.8 22 SSH info 123456} {9.9.9.9 6379 REDIS info 123456} {108.61.223.105 2222 SSH info 123456} {127.0.0.1 3306 MYSQL info } {8.8.8.8 22 SSH info } {9.9.9.9 6379 REDIS info } {108.61.223.105 2222 SSH info } {127.0.0.1 3306 MYSQL adm admin} {8.8.8.8 22 SSH adm admin} {9.9.9.9 6379 REDIS adm admin} {108.61.223.105 2222 SSH adm admin} {127.0.0.1 3306 MYSQL adm admin888} {8.8.8.8 22 SSH adm admin888} {9.9.9.9 6379 REDIS adm admin888} {108.61.223.105 2222 SSH adm admin888} {127.0.0.1 3306 MYSQL adm 123456} {8.8.8.8 22 SSH adm 123456} {9.9.9.9 6379 REDIS adm 123456} {108.61.223.105 2222 SSH adm 123456} {127.0.0.1 3306 MYSQL adm } {8.8.8.8 22 SSH adm } {9.9.9.9 6379 REDIS adm } {108.61.223.105 2222 SSH adm } {127.0.0.1 3306 MYSQL mysql admin} {8.8.8.8 22 SSH mysql admin} {9.9.9.9 6379 REDIS mysql admin} {108.61.223.105 2222 SSH mysql admin} {127.0.0.1 3306 MYSQL mysql admin888} {8.8.8.8 22 SSH mysql admin888} {9.9.9.9 6379 REDIS mysql admin888} {108.61.223.105 2222 SSH mysql admin888} {127.0.0.1 3306 MYSQL mysql 123456} {8.8.8.8 22 SSH mysql 123456} {9.9.9.9 6379 REDIS mysql 123456} {108.61.223.105 2222 SSH mysql 123456} {127.0.0.1 3306 MYSQL mysql } {8.8.8.8 22 SSH mysql } {9.9.9.9 6379 REDIS mysql } {108.61.223.105 2222 SSH mysql } {127.0.0.1 3306 MYSQL user admin} {8.8.8.8 22 SSH user admin} {9.9.9.9 6379 REDIS user admin} {108.61.223.105 2222 SSH user admin} {127.0.0.1 3306 MYSQL user admin888} {8.8.8.8 22 SSH user admin888} {9.9.9.9 6379 REDIS user admin888} {108.61.223.105 2222 SSH user admin888} {127.0.0.1 3306 MYSQL user 123456} {8.8.8.8 22 SSH user 123456} {9.9.9.9 6379 REDIS user 123456} {108.61.223.105 2222 SSH user 123456} {127.0.0.1 3306 MYSQL user } {8.8.8.8 22 SSH user } {9.9.9.9 6379 REDIS user } {108.61.223.105 2222 SSH user } {127.0.0.1 3306 MYSQL administrator admin} {8.8.8.8 22 SSH administrator admin} {9.9.9.9 6379 REDIS administrator admin} {108.61.223.105 2222 SSH administrator admin} {127.0.0.1 3306 MYSQL administrator admin888} {8.8.8.8 22 SSH administrator admin888} {9.9.9.9 6379 REDIS administrator admin888} {108.61.223.105 2222 SSH administrator admin888} {127.0.0.1 3306 MYSQL administrator 123456} {8.8.8.8 22 SSH administrator 123456} {9.9.9.9 6379 REDIS administrator 123456} {108.61.223.105 2222 SSH administrator 123456} {127.0.0.1 3306 MYSQL administrator } {8.8.8.8 22 SSH administrator } {9.9.9.9 6379 REDIS administrator } {108.61.223.105 2222 SSH administrator } {127.0.0.1 3306 MYSQL ftp admin} {8.8.8.8 22 SSH ftp admin} {9.9.9.9 6379 REDIS ftp admin} {108.61.223.105 2222 SSH ftp admin} {127.0.0.1 3306 MYSQL ftp admin888} {8.8.8.8 22 SSH ftp admin888} {9.9.9.9 6379 REDIS ftp admin888} {108.61.223.105 2222 SSH ftp admin888} {127.0.0.1 3306 MYSQL ftp 123456} {8.8.8.8 22 SSH ftp 123456} {9.9.9.9 6379 REDIS ftp 123456} {108.61.223.105 2222 SSH ftp 123456} {127.0.0.1 3306 MYSQL ftp } {8.8.8.8 22 SSH ftp } {9.9.9.9 6379 REDIS ftp } {108.61.223.105 2222 SSH ftp } {127.0.0.1 3306 MYSQL sa admin} {8.8.8.8 22 SSH sa admin} {9.9.9.9 6379 REDIS sa admin} {108.61.223.105 2222 SSH sa admin} {127.0.0.1 3306 MYSQL sa admin888} {8.8.8.8 22 SSH sa admin888} {9.9.9.9 6379 REDIS sa admin888} {108.61.223.105 2222 SSH sa admin888} {127.0.0.1 3306 MYSQL sa 123456} {8.8.8.8 22 SSH sa 123456} {9.9.9.9 6379 REDIS sa 123456} {108.61.223.105 2222 SSH sa 123456} {127.0.0.1 3306 MYSQL sa } {8.8.8.8 22 SSH sa } {9.9.9.9 6379 REDIS sa } {108.61.223.105 2222 SSH sa }] 176
=== RUN   TestDistributionTask
[0000]  INFO xsec crack: Start to scan
[0003]  INFO xsec crack: Finshed scan, total result: 0, used time: 2562047h47m16.854775807s
--- PASS: TestDistributionTask (3.01s)
PASS
ok  	command-line-arguments	3.035s
```

到此为止，我们的扫描器的核心部件已经造好了，接下来需要给扫描器上个高上大的命令行调用的外壳就大功造成了。

### 命令行模块

命令行控制模块，我们单独定义了一个cmd包，依赖第三方包`github.com/urfave/cli`。

我们在cmd模块中定义了扫描和扫描结果导出为txt文件2个命令及一系统全局选项。

```go

package cmd

import (
	"github.com/urfave/cli"

	"x-crack/util"
	"x-crack/models"
)

var Scan = cli.Command{
	Name:        "scan",
	Usage:       "start to crack weak password",
	Description: "start to crack weak password",
	Action:      util.Scan,
	Flags: []cli.Flag{
		boolFlag("debug, d", "debug mode"),
		intFlag("timeout, t", 5, "timeout"),
		intFlag("scan_num, n", 5000, "thread num"),
		stringFlag("ip_list, i", "iplist.txt", "iplist"),
		stringFlag("user_dict, u", "user.dic", "user dict"),
		stringFlag("pass_dict, p", "pass.dic", "password dict"),
	},
}

var Dump = cli.Command{
	Name:        "dump",
	Usage:       "dump result to a text file",
	Description: "dump result to a text file",
	Action:      models.Dump,
	Flags: []cli.Flag{
		stringFlag("outfile, o", "x_crack.txt", "scan result file"),
	},
}

func stringFlag(name, value, usage string) cli.StringFlag {
	return cli.StringFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

func boolFlag(name, usage string) cli.BoolFlag {
	return cli.BoolFlag{
		Name:  name,
		Usage: usage,
	}
}

func intFlag(name string, value int, usage string) cli.IntFlag {
	return cli.IntFlag{
		Name:  name,
		Value: value,
		Usage: usage,
	}
}

```

然后再回到`x-crack/util`包为我们的scan command模块专门写一个Action，如下：

```go

func Scan(ctx *cli.Context) (err error) {
	if ctx.IsSet("debug") {
		vars.DebugMode = ctx.Bool("debug")
	}

	if vars.DebugMode {
		logger.Log.Level = logrus.DebugLevel
	}

	if ctx.IsSet("timeout") {
		vars.TimeOut = time.Duration(ctx.Int("timeout")) * time.Second
	}

	if ctx.IsSet("scan_num") {
		vars.ScanNum = ctx.Int("scan_num")
	}

	if ctx.IsSet("ip_list") {
		vars.IpList = ctx.String("ip_list")
	}

	if ctx.IsSet("user_dict") {
		vars.UserDict = ctx.String("user_dict")
	}

	if ctx.IsSet("pass_dict") {
		vars.PassDict = ctx.String("pass_dict")
	}

	if ctx.IsSet("outfile") {
		vars.ResultFile = ctx.String("outfile")
	}

	vars.StartTime = time.Now()

	userDict, uErr := ReadUserDict(vars.UserDict)
	passDict, pErr := ReadPasswordDict(vars.PassDict)
	ipList := ReadIpList(vars.IpList)
	aliveIpList := CheckAlive(ipList)
	if uErr == nil && pErr == nil {
		tasks, _ := GenerateTask(aliveIpList, userDict, passDict)
		DistributionTask(tasks)
	}
	return err
}

```

然后再到`x-crack/models`中为dump命令写一个Action，如下：

```go

package models

import (
	"github.com/patrickmn/go-cache"
	"github.com/urfave/cli"

	"x-crack/vars"
	"x-crack/logger"
	"x-crack/util/hash"

	"encoding/gob"
	"time"
	"fmt"
	"os"
	"strings"
)

func init() {
	gob.Register(Service{})
	gob.Register(ScanResult{})
}

func SaveResult(err error, result ScanResult) {
	if err == nil && result.Result {
		var k string
		protocol := strings.ToUpper(result.Service.Protocol)

		if protocol == "REDIS" || protocol == "FTP" {
			k = fmt.Sprintf("%v-%v-%v", result.Service.Ip, result.Service.Port, result.Service.Protocol)
		} else {
			k = fmt.Sprintf("%v-%v-%v", result.Service.Ip, result.Service.Port, result.Service.Username)
		}

		h := hash.MakeTaskHash(k)
		hash.SetTaskHask(h)

		_, found := vars.CacheService.Get(k)
		if !found {
			logger.Log.Infof("Ip: %v, Port: %v, Protocol: [%v], Username: %v, Password: %v", result.Service.Ip,
				result.Service.Port, result.Service.Protocol, result.Service.Username, result.Service.Password)
		}
		vars.CacheService.Set(k, result, cache.NoExpiration)
	}
}

func SavaResultToFile() (error) {
	return vars.CacheService.SaveFile("x_crack.db")
}

func CacheStatus() (count int, items map[string]cache.Item) {
	count = vars.CacheService.ItemCount()
	items = vars.CacheService.Items()
	return count, items
}

func ResultTotal() {
	vars.ProgressBar.Finish()
	logger.Log.Info(fmt.Sprintf("Finshed scan, total result: %v, used time: %v",
		vars.CacheService.ItemCount(),
		time.Since(vars.StartTime)))
}

func LoadResultFromFile() {
	vars.CacheService.LoadFile("x_crack.db")
	vars.ProgressBar.Finish()
	logger.Log.Info(fmt.Sprintf("Finshed scan, total result: %v", vars.CacheService.ItemCount()))
}

func Dump(ctx *cli.Context) (err error) {
	LoadResultFromFile()

	err = DumpToFile(vars.ResultFile)
	if err != nil {
		logger.Log.Fatalf("Dump result to file err, Err: %v", err)
	}
	return err
}

func DumpToFile(filename string) (err error) {
	file, err := os.Create(filename)
	if err == nil {
		_, items := CacheStatus()
		for _, v := range items {
			result := v.Object.(ScanResult)
			file.WriteString(fmt.Sprintf("%v:%v|%v,%v:%v\n", result.Service.Ip, result.Service.Port,
				result.Service.Protocol, result.Service.Username, result.Service.Password))
		}
	}
	return err
}

```

最后给IP\port过滤与任务扫描模块加上一个骚气的进度条，我们的扫描器就算大功告成了。

`x-crack/util/util.go`的代码片段：

```go

func CheckAlive(ipList []models.IpAddr) ([]models.IpAddr) {
	logger.Log.Infoln("checking ip active")
	vars.ProcessBarActive = pb.StartNew(len(ipList))
	vars.ProcessBarActive.SetTemplate(`{{ rndcolor "Checking progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor}}  {{rtime . | rndcolor }}`)
....
```

`x-crack/util/task.go`的代码片断：
```go
func DistributionTask(tasks []models.Service) () {
	totalTask := len(tasks)
	scanBatch := totalTask / vars.ScanNum
	logger.Log.Infoln("Start to scan")
	vars.ProgressBar = pb.StartNew(totalTask)
	vars.ProgressBar.SetTemplate(`{{ rndcolor "Scanning progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor }} {{rtime . | rndcolor}} `)
...
```

扫描器代码中还有些细节没有在教程中详细说，有兴趣的同学可以思考下以下问题，然后再结合代码看看老夫的实现方式：
1. 扫到一个弱口令后，如何取消相同IP\port和用户名请求，避免扫描效率低下
1. 对于FTP匿名访问，如何只记录一个密码，而不是把所有用户名都记录下来
1. 对于Redis这种没有用户名的服务，如何只记录一次密码，而不是记录所有的所有用户及正常的密码的组合
1. 对于不支持设置超时的扫描插件，如何统一设置超时时间

## 扫描器测试

到现在为止，我们的扫描器已经大功告成了，可以编译出来运行一下看看效果了。以下脚本可一键同时编译出mac、linux和Windows平台的可执行文件（笔者的开发环境为MAC）
```bash
#!/bin/bash

go build x-crack.go
mv x-crack x-crack_darwin_amd64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build x-crack.go
mv x-crack x-crack_linux_amd64
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build x-crack.go
mv x-crack.exe x-crack_windows_amd64.exe
go build x-crack.go
```


### 使用参数
```bash
hartnett at hartnettdeMacBook-Pro in /data/code/golang/src/x-crack (master)
$ ./x-crack
NAME:
   x-crack - Weak password scanner, Support: FTP/SSH/MSSQL/MYSQL/PostGreSQL/REDIS/ElasticSearch/MONGODB

USAGE:
   x-crack [global options] command [command options] [arguments...]

VERSION:
   20171227

AUTHOR(S):
   netxfly <x@xsec.io>

COMMANDS:
     scan     start to crack weak password
     dump     dump result to a text file
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug, -d                  debug mode
   --timeout value, -t value    timeout (default: 5)
   --scan_num value, -n value   thread num (default: 5000)
   --ip_list value, -i value    iplist (default: "iplist.txt")
   --user_dict value, -u value  user dict (default: "user.dic")
   --pass_dict value, -p value  password dict (default: "pass.dic")
   --outfile value, -o value    scan result file (default: "x_crack.txt")
   --help, -h                   show help
   --version, -v                print the version
```

### 使用截图
![](https://docs.xsec.io/images/x-crack/x-crack001.png)

![](https://docs.xsec.io/images/x-crack/x-crack002.png)