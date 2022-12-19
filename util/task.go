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

package util

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"gopkg.in/cheggaaa/pb.v2"

	"x-crack/logger"
	"x-crack/models"
	"x-crack/plugins"
	"x-crack/util/hash"
	"x-crack/vars"

	"fmt"
	"runtime"
	"strings"
	"sync"
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

func DistributionTask(tasks []models.Service) {
	totalTask := len(tasks)
	scanBatch := totalTask / vars.ScanNum
	logger.Log.Infoln("Start to scan")
	vars.ProgressBar = pb.StartNew(totalTask)
	vars.ProgressBar.SetTemplate(`{{ rndcolor "Scanning progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor }} {{rtime . | rndcolor}} `)

	for i := 0; i < scanBatch; i++ {
		curTasks := tasks[vars.ScanNum*i : vars.ScanNum*(i+1)]
		ExecuteTask(curTasks)
	}

	if totalTask%vars.ScanNum > 0 {
		lastTask := tasks[vars.ScanNum*scanBatch : totalTask]
		ExecuteTask(lastTask)
	}

	models.SavaResultToFile()
	models.ResultTotal()
	models.DumpToFile(vars.ResultFile)
}

func ExecuteTask(tasks []models.Service) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for _, task := range tasks {
		if vars.DebugMode {
			logger.Log.Debugf("checking: Ip: %v, Port: %v, [%v], UserName: %v, Password: %v", task.Ip, task.Port,
				task.Protocol, task.Username, task.Password)
		}

		var k string
		protocol := strings.ToUpper(task.Protocol)

		if protocol == "REDIS" || protocol == "FTP" || protocol == "SNMP" {
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

func RunTask(tasks []models.Service) {
	totalTask := len(tasks)
	vars.ProgressBar = pb.StartNew(totalTask)
	vars.ProgressBar.SetTemplate(`{{ rndcolor "Scanning progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor }} {{rtime . | rndcolor}} `)

	wg := &sync.WaitGroup{}

	// 创建一个buffer为vars.threadNum * 2的channel
	taskChan := make(chan models.Service, vars.ScanNum*2)

	// 创建vars.ThreadNum个协程
	for i := 0; i < vars.ScanNum; i++ {
		go crackPassword(taskChan, wg)
	}

	// 生产者，不断地往taskChan channel发送数据，直到channel阻塞
	for _, task := range tasks {
		wg.Add(1)
		taskChan <- task
	}

	close(taskChan)
	waitTimeout(wg, vars.TimeOut*2)
}

// 每个协程都从channel中读取数据后开始扫描并保存
func crackPassword(taskChan chan models.Service, wg *sync.WaitGroup) {
	for task := range taskChan {
		vars.ProgressBar.Increment()

		if vars.DebugMode {
			logger.Log.Debugf("checking: Ip: %v, Port: %v, [%v], UserName: %v, Password: %v, goroutineNum: %v", task.Ip, task.Port,
				task.Protocol, task.Username, task.Password, runtime.NumGoroutine())
		}

		var k string
		protocol := strings.ToUpper(task.Protocol)

		if protocol == "REDIS" || protocol == "FTP" || protocol == "SNMP" {
			k = fmt.Sprintf("%v-%v-%v", task.Ip, task.Port, task.Protocol)
		} else {
			k = fmt.Sprintf("%v-%v-%v", task.Ip, task.Port, task.Username)
		}

		h := hash.MakeTaskHash(k)
		if hash.CheckTashHash(h) {
			wg.Done()
			continue
		}

		fn := plugins.ScanFuncMap[protocol]
		models.SaveResult(fn(task))
		wg.Done()
	}
}

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

	if ctx.IsSet("ip") {
		vars.IP = ctx.String("ip")
	}

	if ctx.IsSet("username") {
		vars.USERNAME = ctx.String("username")
	}

	if ctx.IsSet("password") {
		vars.PASSWORD = ctx.String("password")
	}
	vars.StartTime = time.Now()
	userDict, uErr := ReadUserDict(vars.UserDict)
	if vars.USERNAME != "" {
		userDict = []string{vars.USERNAME}
	}
	passDict, pErr := ReadPasswordDict(vars.PassDict)
	if vars.PASSWORD != "" {
		passDict = []string{vars.PASSWORD}
	}
	ipList := ReadIpList(vars.IpList)
	var aliveIpList []models.IpAddr
	if vars.IP != "" {
		aliveIpList = CheckAliveIpAddr(vars.IP)
	} else {
		aliveIpList = CheckAlive(ipList)
	}
	if uErr == nil && pErr == nil {
		logger.Log.Printf("Got %v services to brute force", len(aliveIpList))
		logger.Log.Printf("Loaded %v usernames, %v passwords", len(userDict), len(passDict))
		tasks, _ := GenerateTask(aliveIpList, userDict, passDict)
		RunTask(tasks)
		// DistributionTask(tasks)
	}
	return err
}

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
