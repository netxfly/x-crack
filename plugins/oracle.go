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

package plugins

//
//import (
//	"database/sql"
//	"fmt"
//	_ "github.com/godror/godror"
//	"x-crack/models"
//)
//
//func ScanOracle(service models.Service) (err error, result models.ScanResult) {
//	result.Service = service
//	dataSource := fmt.Sprintf(`user="%v" password="%v" connectString="%v:%v/%v"`,
//		service.Username,
//		service.Password,
//		service.Ip,
//		service.Port,
//		"oracle")
//	db, err := sql.Open("godror", dataSource)
//	if err != nil {
//		panic(err)
//	}
//	defer db.Close()
//	err = db.Ping()
//	if err != nil {
//		fmt.Println(err.Error())
//		panic(err)
//	}
//	fmt.Println("链接成功")
//	result.Result = true
//	return err, result
//}
