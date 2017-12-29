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
