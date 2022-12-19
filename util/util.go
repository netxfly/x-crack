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
	"gopkg.in/cheggaaa/pb.v2"
	"x-crack/logger"
	"x-crack/models"
	"x-crack/vars"

	"fmt"
	"net"
	"sync"
)

var (
	AliveAddr []models.IpAddr
	mutex     sync.Mutex
)

func init() {
	AliveAddr = make([]models.IpAddr, 0)
}

func CheckAlive(ipList []models.IpAddr) []models.IpAddr {
	logger.Log.Infoln("checking ip active")
	vars.ProcessBarActive = pb.StartNew(len(ipList))
	vars.ProcessBarActive.SetTemplate(`{{ rndcolor "Checking progress: " }} {{  percent . "[%.02f%%]" "[?]"| rndcolor}} {{ counters . "[%s/%s]" "[%s/?]" | rndcolor}} {{ bar . "「" "-" (rnd "ᗧ" "◔" "◕" "◷" ) "•" "」" | rndcolor}}  {{rtime . | rndcolor }}`)

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

type mm string

func (m *mm) Scan(state fmt.ScanState, verb rune) error {
	tok, err := state.Token(true, func(r rune) bool {
		// 默认string以空格分隔,我这里改为用逗号分隔
		var ret bool
		if r != ':' &&
			r != '/' {
			ret = true
		}

		return ret
	})
	if err != nil {
		return err
	}
	*m = mm(tok)
	return nil
}

func CheckAliveIpAddr(IP string) []models.IpAddr {
	logger.Log.Infoln("checking ip active")

	var IpAddr models.IpAddr
	// Fix err `unexpected EOF`, see: https://studygolang.com/topics/15855
	_, err := fmt.Sscanf(IP, "%s://%s:%d",
		(*mm)(&IpAddr.Protocol),
		(*mm)(&IpAddr.Ip),
		&IpAddr.Port,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(IpAddr)

	alive, ipAddr := check(IpAddr)

	if alive {
		AliveAddr = []models.IpAddr{
			ipAddr,
		}
	}

	return AliveAddr
}

func check(ipAddr models.IpAddr) (bool, models.IpAddr) {
	alive := false
	if vars.UdpProtocols[ipAddr.Protocol] {
		_, err := net.DialTimeout("udp", fmt.Sprintf("%v:%v", ipAddr.Ip, ipAddr.Port), vars.TimeOut)
		if err == nil {
			alive = true
		}
	} else {
		_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ipAddr.Ip, ipAddr.Port), vars.TimeOut)
		if err == nil {
			alive = true
		}
	}

	return alive, ipAddr
}

func SaveAddr(alive bool, ipAddr models.IpAddr) {
	if alive {
		mutex.Lock()
		AliveAddr = append(AliveAddr, ipAddr)
		mutex.Unlock()
	}
}
